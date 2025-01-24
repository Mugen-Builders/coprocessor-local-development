package root

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/coprocessor_contracts"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/rollups_contracts"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/tools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"
)

//go:embed anvil_state.json
var devnetState []byte

const stateFileName = "anvil_state.json"

var (
	configPath string
	Cmd        = &cobra.Command{
		Use:   "nonodox",
		Short: "nonodox is a tool for local development of Cartesi coprocessor",
		Long:  "This tool listens for input events, processes notices, and vouchers linked to an input, and executes them on-chain (anvil).",
		Run: func(cmd *cobra.Command, args []string) {
			if err := configs.LoadConfig(configPath); err != nil {
				slog.Error("Failed to load configuration file", "error", err)
				os.Exit(1)
			}
			run()
		},
	}
)

var startupMessage = `
NonodoX local development tool started

Anvil running at ANVIL_HTTP_URL
GraphQL polling at http://localhost:8080/graphql
CoprocessorAdapter address COPROCESSOR_ADAPTER_CONTRACT_ADDRESS
Coprocessor application machine hash COPROCESSOR_MACHINE_HASH

Press Ctrl+C to stop the application.
`

func init() {
	Cmd.Flags().StringVar(&configPath, "config", "", "Path to the configuration file (required)")
	if err := Cmd.MarkFlagRequired("config"); err != nil {
		os.Exit(1)
	}

	Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		configs.ConfigureLog(slog.LevelInfo)
	}
}

func run() {
	_ = exec.Command("pkill", "nonodo").Run()
	_ = exec.Command("pkill", "anvil").Run()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	group, ctx := errgroup.WithContext(ctx)

	dir, err := makeStateTemp()
	if err != nil {
		slog.Error("Failed to create temp state directory", "error", err)
		os.Exit(1)
	}
	defer removeTemp(dir)

	anvil := exec.CommandContext(
		ctx,
		"anvil",
		"--load-state",
		path.Join(dir, stateFileName),
	)
	anvil.Stdout = os.Stdout
	anvil.Stderr = os.Stderr
	group.Go(func() error {
		return anvil.Run()
	})

	nonodo := exec.CommandContext(
		ctx,
		"nonodo",
		"--disable-devnet",
		"--rpc-url",
		"ws://localhost:8545",
		"--contracts-input-box-block",
		"7",
	)
	notifyWriter := tools.NewNotifyWriter(io.Discard, "nonodo: ready")
	nonodo.Stdout = notifyWriter
	group.Go(func() error {
		return nonodo.Run()
	})

	go func() {
		<-signalChan
		slog.Info("Received termination signal, shutting down gracefully...")
		cancel()
		_ = anvil.Process.Kill()
		_ = nonodo.Process.Kill()
	}()

	select {
	case <-notifyWriter.Ready():
		message := strings.NewReplacer(
			"ANVIL_HTTP_URL", os.Getenv("ANVIL_HTTP_URL"),
			"COPROCESSOR_ADAPTER_CONTRACT_ADDRESS", os.Getenv("COPROCESSOR_ADAPTER_CONTRACT_ADDRESS"),
			"COPROCESSOR_MACHINE_HASH", os.Getenv("COPROCESSOR_MACHINE_HASH"),
		).Replace(startupMessage)
		fmt.Println(message)
	case <-ctx.Done():
		slog.Error("Context canceled before nonodo became ready")
		return
	}

	ethClient, opts, err := configs.SetupTransactor()
	if err != nil {
		slog.Error("Failed to set up transactor", "error", err)
		os.Exit(1)
	}

	outputsChan := make(chan struct {
		Event          rollups_contracts.IInputBoxInputAdded
		MachineOutputs [][]byte
	})
	inputAddedChan := make(chan rollups_contracts.IInputBoxInputAdded)
	taskIssuedChan := make(chan coprocessor_contracts.MockCoprocessorTaskIssued)

	go func() {
		reader, err := NewTaskReader()
		if err != nil {
			slog.Error("Failed to set up task reader", "error", err)
			os.Exit(1)
		}

		if err := reader.GetTaskIssuedEvents(ctx, taskIssuedChan); err != nil {
			slog.Error("Error fetching tasks", "error", err)
			os.Exit(1)
		}

		for {
			select {
			case input := <-taskIssuedChan:
				instance, err := rollups_contracts.NewIInputBox(
					common.HexToAddress("0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE"),
					ethClient,
				)
				if err != nil {
					slog.Error("Failed to create inputBox instance", "error", err)
					os.Exit(1)
				}

				tx, err := instance.AddInput(
					opts,
					common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"),
					input.Input,
				)
				if err != nil {
					slog.Error("Failed to call issueTask function", "error", err, "input", input)
					os.Exit(1)
				}

				slog.Info("Input added", "dapp", common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"), "tx", tx.Hash().Hex())
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		reader, err := NewInputReader()
		if err != nil {
			slog.Error("Failed to set up inputBox reader", "error", err)
			os.Exit(1)
		}

		if err := reader.GetInputAddedEvents(ctx, inputAddedChan); err != nil {
			slog.Error("Error fetching inputs", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		findOutputsByIdUseCase, err := NewFindOutputsByIdUseCase("http://localhost:8080/graphql", nil)
		if err != nil {
			slog.Error("Failed to setup node reader", "error", err)
			os.Exit(1)
		}

		for event := range inputAddedChan {
			outputs, err := findOutputsByIdUseCase.Execute(ctx, int(event.InputIndex.Int64()))
			if err != nil {
				slog.Error("Failed to find outputs by ID", "error", err)
				os.Exit(1)
			}

			if len(outputs) > 0 {
				outputsChan <- struct {
					Event          rollups_contracts.IInputBoxInputAdded
					MachineOutputs [][]byte
				}{Event: event, MachineOutputs: outputs}
			}
		}
	}()

	for {
		select {
		case output := <-outputsChan:
			instance, err := coprocessor_contracts.NewCoprocessorContracts(
				common.HexToAddress(os.Getenv("COPROCESSOR_ADAPTER_CONTRACT_ADDRESS")),
				ethClient,
			)
			if err != nil {
				slog.Error("Failed to create coprocessor instance", "error", err)
				return
			}

			hash := sha3.NewLegacyKeccak256()
			hash.Write(output.Event.Input)
			tx, err := instance.CoprocessorCallbackOutputsOnly(
				opts,
				common.HexToHash(os.Getenv("COPROCESSOR_MACHINE_HASH")),
				common.HexToHash(fmt.Sprintf("0x%x", hash.Sum(nil))),
				output.MachineOutputs,
			)
			if err != nil {
				slog.Error("Failed to call coprocessor callback function", "error", err, "outputs", output.MachineOutputs)
				return
			}

			slog.Info("Outputs executed", "tx", tx.Hash().Hex(), "outputs", output.MachineOutputs)
		case <-ctx.Done():
			return
		}
	}
}

func makeStateTemp() (string, error) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("anvil: failed to create temp dir: %w", err)
	}
	stateFile := path.Join(tempDir, stateFileName)
	const permissions = 0644
	err = os.WriteFile(stateFile, devnetState, permissions)
	if err != nil {
		return "", fmt.Errorf("anvil: failed to write state file: %w", err)
	}
	return tempDir, nil
}

func removeTemp(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		slog.Warn("anvil: failed to remove temp file", "error", err)
	}
}
