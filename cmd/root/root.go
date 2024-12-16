package root

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/coprocessor-local-solver/configs"
	"github.com/henriquemarlon/coprocessor-local-solver/internal/evm_reader"
	"github.com/henriquemarlon/coprocessor-local-solver/internal/node_reader"
	"github.com/henriquemarlon/coprocessor-local-solver/pkg/coprocessor_contracts"
	"github.com/henriquemarlon/coprocessor-local-solver/pkg/rollups_contracts"
	"github.com/spf13/cobra"
)

var (
	cfg                      *configs.CFG
	coprocessorCallerAddress string
	Cmd                      = &cobra.Command{
		Use:   "mugen",
		Short: "Mugen Builders tool for executing coprocessor outputs",
		Long:  "This tool listens for input events, processes the associated notices, and executes them on-chain.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.COPROCESSOR_CALLER_MOCK_ADDRESS = coprocessorCallerAddress
			run()
		},
	}
)

var startupMessage = `
Mugen Builders <> Cartesi Coprocessor local development tool started
GraphQL polling pointing to GRAPHQL_URL
CoprocessorCaller address MOCK_ADDRESS

Press Ctrl+C to stop the application.
`

func init() {
	var err error
	cfg, err = configs.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	Cmd.Flags().StringVar(&coprocessorCallerAddress, "coprocessor-caller-address", "", "Address of the coprocessor caller")
	Cmd.MarkFlagRequired("coprocessor-caller-address")
}

func run() {
	message := strings.ReplaceAll(startupMessage, "GRAPHQL_URL", cfg.GRAPHQL_URL)
	message = strings.ReplaceAll(message, "MOCK_ADDRESS", cfg.COPROCESSOR_CALLER_MOCK_ADDRESS)
	fmt.Println(message)

	ethClient, opts, err := configs.SetupTransactor(cfg)
	if err != nil {
		slog.Error("failed to set up transactor", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gqlClient := genqlient.NewClient(cfg.GRAPHQL_URL, nil)
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
