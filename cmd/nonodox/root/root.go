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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
		_ = exec.Command("pkill", "nonodo").Run()
	}()

	group, ctx := errgroup.WithContext(ctx)
	nonodo := exec.CommandContext(ctx, "nonodo")
	notifyWriter := tools.NewNotifyWriter(io.Discard, "nonodo: ready")
	nonodo.Stdout = notifyWriter
	group.Go(nonodo.Run)

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
	inputBoxChan := make(chan rollups_contracts.IInputBoxInputAdded)
	coprocessorInputChan := make(chan rollups_contracts.IInputBoxInputAdded)

	go func() {
		reader, err := NewTransactor()
		if err != nil {
			slog.Error("Failed to set up transactor", "error", err)
			os.Exit(1)
		}

		if err := reader.GetInputAddedEvents(ctx, inputBoxChan); err != nil {
			slog.Error("Error fetching inputs", "error", err)
			os.Exit(1)
		}

		for {
			select {
			case input := <-inputBoxChan:
				instance, err := coprocessor_contracts.NewCoprocessorContracts(
					common.HexToAddress(os.Getenv("COPROCESSOR_ADAPTER_CONTRACT_ADDRESS")),
					ethClient,
				)
				if err != nil {
					slog.Error("Failed to create coprocessor instance", "error", err)
					os.Exit(1)
				}

				_, err = instance.CallCoprocessor(opts, input.Input)
				if err != nil {
					slog.Error("Failed to call coprocessor function", "error", err, "input", input)
					os.Exit(1)
				}

				slog.Info("New input", "dapp", input.Dapp, "index", input.InputIndex, "sender", input.Sender)
				coprocessorInputChan <- input
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

		for event := range coprocessorInputChan {
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