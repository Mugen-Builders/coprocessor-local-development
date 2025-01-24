package root

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
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

var (
	configPath string
	cfg        *configs.Config
	Cmd        = &cobra.Command{
		Use:   "nonodox",
		Short: "nonodox is a tool for local development of Cartesi coprocessor",
		Long:  "This tool listens for input events, processes notices, and vouchers linked to an input, and executes them on-chain (anvil).",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			cfg, err = configs.LoadConfig(configPath)
			if err != nil {
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
Coprocessor Adapter contract address COPROCESSOR_ADAPTER_ADDRESS
Coprocessor application machine hash COPROCESSOR_MACHINE_HASH

Press Ctrl+C to stop the application.
`

func init() {
	Cmd.Flags().StringVar(&configPath, "config", "", "Path to the configuration file (required)")
	if err := Cmd.MarkFlagRequired("config"); err != nil {
		slog.Error("Failed to mark flag as required", "error", err)
		os.Exit(1)
	}

	Cmd.PreRun = func(cmd *cobra.Command, args []string) {
		configs.ConfigureLog(slog.LevelInfo)
	}
}

func run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		slog.Info("Received termination signal, shutting down gracefully...")
		cancel()
		_ = exec.Command("pkill", "nonodo").Run()
	}()

	group, ctx := errgroup.WithContext(ctx)
	nonodo := exec.CommandContext(
		ctx,
		"nonodo",
		"--disable-devnet",
		"--rpc-url",
		cfg.AnvilWsURL,
		"--contracts-input-box-block",
		cfg.AnvilInputBoxBlock,
	)
	notifyWriter := tools.NewNotifyWriter(io.Discard, "nonodo: ready")
	nonodo.Stdout = notifyWriter
	group.Go(nonodo.Run)

	select {
	case <-notifyWriter.Ready():
		message := strings.NewReplacer(
			"ANVIL_HTTP_URL", cfg.AnvilHttpURL,
			"COPROCESSOR_MACHINE_HASH", cfg.CoprocessorMachineHash,
			"COPROCESSOR_ADAPTER_ADDRESS", cfg.CoprocessorAdapterAddress,
		).Replace(startupMessage)
		fmt.Println(message)
	case <-ctx.Done():
		slog.Error("Context canceled before nonodo became ready, make sure that all environment variables are set correctly")
		return
	}

	ethClient, opts, err := configs.SetupTransactor()
	if err != nil {
		slog.Error("Failed to set up transactor", "error", err)
		os.Exit(1)
	}

	inputsHash := make(map[common.Hash]bool)
	chann1 := make(chan coprocessor_contracts.MockCoprocessorTaskIssued)
	chann2 := make(chan rollups_contracts.IInputBoxInputAdded)
	chann3 := make(chan struct {
		Event       rollups_contracts.IInputBoxInputAdded
		PayloadHash common.Hash
	})
	chann4 := make(chan struct {
		Event          rollups_contracts.IInputBoxInputAdded
		PayloadHash    common.Hash
		MachineOutputs [][]byte
	})

	go func() {
		reader, err := NewTaskReader()
		if err != nil {
			slog.Error("Failed to set up task reader", "error", err)
			os.Exit(1)
		}

		if err := reader.GetTaskIssuedEvents(ctx, chann1); err != nil {
			slog.Error("Error fetching tasks", "error", err)
			os.Exit(1)
		}

		for {
			select {
			case task := <-chann1:
				instance, err := rollups_contracts.NewIInputBox(
					common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768"),
					ethClient,
				)
				if err != nil {
					slog.Error("Failed to create inputBox instance", "error", err)
					os.Exit(1)
				}

				hash := sha3.NewLegacyKeccak256()
				hash.Write(task.Input)
				inputsHash[common.HexToHash(fmt.Sprintf("0x%x", hash.Sum(nil)))] = true

				tx, err := instance.AddInput(
					opts,
					common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"),
					task.Input,
				)
				if err != nil {
					slog.Error("Failed to call addInput function", "error", err, "input", task.Input)
					os.Exit(1)
				}
				slog.Info("Input added", "dapp", "0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e", "tx", tx.Hash().Hex())
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
		if err := reader.GetInputAddedEvents(ctx, chann2); err != nil {
			slog.Error("Error fetching inputs", "error", err)
			os.Exit(1)
		}
		for {
			select {
			case event := <-chann2:
				hash := sha3.NewLegacyKeccak256()
				hash.Write(event.Input)
				payloadHash := common.HexToHash(fmt.Sprintf("0x%x", hash.Sum(nil)))
				if _, ok := inputsHash[payloadHash]; !ok {
					slog.Error("Input not found", "input", event.Input)
				}

				chann3 <- struct {
					Event       rollups_contracts.IInputBoxInputAdded
					PayloadHash common.Hash
				}{Event: event, PayloadHash: payloadHash}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		findOutputsByIdUseCase, err := NewFindOutputsByIdUseCase("http://localhost:8080/graphql", nil)
		if err != nil {
			slog.Error("Failed to setup node reader", "error", err)
			os.Exit(1)
		}

		for data := range chann3 {
			outputs, err := findOutputsByIdUseCase.Execute(ctx, int(data.Event.InputIndex.Int64()))
			if err != nil {
				slog.Error("Failed to find outputs by ID", "error", err)
				os.Exit(1)
			}

			if len(outputs) > 0 {
				chann4 <- struct {
					Event          rollups_contracts.IInputBoxInputAdded
					PayloadHash    common.Hash
					MachineOutputs [][]byte
				}{Event: data.Event, PayloadHash: data.PayloadHash, MachineOutputs: outputs}
			}
		}
	}()

	for {
		select {
		case output := <-chann4:
			instance, err := coprocessor_contracts.NewCoprocessorContracts(
				common.HexToAddress(cfg.CoprocessorAdapterAddress),
				ethClient,
			)
			if err != nil {
				slog.Error("Failed to create coprocessor instance", "error", err)
				return
			}

			tx, err := instance.CoprocessorCallbackOutputsOnly(
				opts,
				common.HexToHash(cfg.CoprocessorMachineHash),
				output.PayloadHash,
				output.MachineOutputs,
			)
			if err != nil {
				slog.Error("Failed to call coprocessor callback function", "error", err, "outputs", output.MachineOutputs)
				return
			}
			slog.Info("Outputs executed", "payload hash", output.PayloadHash, "tx", tx.Hash().Hex(), "outputs", output.MachineOutputs)
		case <-ctx.Done():
			return
		}
	}
}