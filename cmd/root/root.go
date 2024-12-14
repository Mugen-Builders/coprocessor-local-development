package root

import (
	"context"
	"log/slog"
	"os"

	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/coprocessor-local-solver/configs"
	"github.com/henriquemarlon/coprocessor-local-solver/internal/evm_reader"
	"github.com/henriquemarlon/coprocessor-local-solver/internal/node_reader"
	"github.com/henriquemarlon/coprocessor-local-solver/pkg/coprocessor_contracts"
	"github.com/henriquemarlon/coprocessor-local-solver/pkg/rollups_contracts"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for events and process execute outputs",
	Long:  `This command listens for events from the input box and processes the notices associated with them.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration file", "error", err)
		os.Exit(1)
	}

	ethClient, opts, err := configs.SetupTransactor(cfg)
	if err != nil {
		slog.Error("Failed to setup transactor", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	gqlClient := genqlient.NewClient("http://127.0.0.1:8080/graphql", nil)

	inputChan := make(chan rollups_contracts.IInputBoxInputAdded)
	go func() {
		ethClientWs, err := configs.SetupTransactorWS(cfg)
		if err != nil {
			slog.Error("Failed to setup transactor WS", "error", err)
			os.Exit(1)
		}
		reader, err := evm_reader.NewEVMReader(
			ethClientWs,
			common.HexToAddress(cfg.INPUT_BOX_ADDRESS),
		)
		if err != nil {
			slog.Error("Failed to create EVM Reader", "error", err)
			os.Exit(1)
		}
		err = reader.GetInputAddedEvents(inputChan)
		if err != nil {
			slog.Error("Failed to get 'InputAdded' events from IInputBox contract", "error", err)
			os.Exit(1)
		}
	}()

	for {
		select {
		case event := <-inputChan:
			reader := node_reader.NewNodeReader(gqlClient)
			outputs, err := reader.GetNoticesByInputIndex(ctx, int(event.InputIndex.Int64()))
			if err != nil {
				slog.Error("Failed to get notices by input index", "error", err)
				os.Exit(1)
			}
			instance, err := coprocessor_contracts.NewCoprocessorContracts(
				common.HexToAddress(cfg.COPROCESSOR_CALLER_MOCK_ADDRESS),
				ethClient,
			)
			if err != nil {
				slog.Error("Failed to create Coprocessor instance", "error", err)
				os.Exit(1)
			}

			if len(outputs) == 0 {
				slog.Info("No outputs to process", "inputIndex", event.InputIndex)
				continue
			}

			slog.Info("Calling CoprocessorCallerMock", "outputs", outputs)

			tx, err := instance.CoprocessorCallbackOutputsOnly(
				opts,
				[32]byte{},
				[32]byte{},
				outputs,
			)
			if err != nil {
				slog.Error("Failed to call CoprocessorCallbackOutputsOnly", "error", err)
				os.Exit(1)
			}
			slog.Info("CoprocessorCallbackOutputsOnly called", "tx", tx.Hash().Hex())
		case <-ctx.Done():
			slog.Info("Context done, shutting down.")
			return
		}
	}
}
