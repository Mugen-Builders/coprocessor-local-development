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
	Short: "listen for events and process execute outputs",
	Long:  `this command listens for events from the input box and processes the notices associated with them.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ethClient, opts, err := configs.SetupTransactor(cfg)
	if err != nil {
		slog.Error("failed to set up transactor", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gqlClient := genqlient.NewClient("http://127.0.0.1:8080/graphql", nil)
	inputChan := make(chan rollups_contracts.IInputBoxInputAdded)

	go func() {
		ethClientWs, err := configs.SetupTransactorWS(cfg)
		if err != nil {
			slog.Error("failed to set up ws transactor", "error", err)
			os.Exit(1)
		}
		reader, err := evm_reader.NewEVMReader(ethClientWs, common.HexToAddress(cfg.INPUT_BOX_ADDRESS))
		if err != nil {
			slog.Error("failed to create evm reader", "error", err)
			os.Exit(1)
		}
		if err := reader.GetInputAddedEvents(ctx, inputChan); err != nil {
			slog.Error("error fetching inputadded events", "error", err)
			os.Exit(1)
		}
	}()

	for {
		select {
		case event := <-inputChan:
			outputs, err := node_reader.NewNodeReader(gqlClient).GetNoticesByInputIndex(ctx, int(event.InputIndex.Int64()))
			if err != nil {
				slog.Error("failed to get notices by input index", "error", err)
				os.Exit(1)
			}

			if len(outputs) == 0 {
				slog.Info("no outputs to process", "inputIndex", event.InputIndex)
				continue
			}

			slog.Info("processing outputs", "inputIndex", event.InputIndex, "outputs", outputs)

			instance, err := coprocessor_contracts.NewCoprocessorContracts(common.HexToAddress(cfg.COPROCESSOR_CALLER_MOCK_ADDRESS), ethClient)
			if err != nil {
				slog.Error("failed to create coprocessor instance", "error", err)
				os.Exit(1)
			}

			tx, err := instance.CoprocessorCallbackOutputsOnly(opts, [32]byte{}, [32]byte{}, outputs)
			if err != nil {
				slog.Error("failed to call coprocessorcallbackoutputsonly", "error", err)
				os.Exit(1)
			}

			slog.Info("output executed", "tx", tx.Hash().Hex(), "inputIndex", event.InputIndex)
		case <-ctx.Done():
			slog.Info("shutting down due to canceled context")
			return
		}
	}
}