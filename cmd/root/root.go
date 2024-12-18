package root

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
GraphQL polling at GRAPHQL_URL
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
	outputChan := make(chan [][]byte)

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
			slog.Error("error fetching inputs", "error", err)
			os.Exit(1)
		}
	}()

	go func(ctx context.Context, gqlClient genqlient.Client, inputChan <-chan rollups_contracts.IInputBoxInputAdded, outputsChan chan<- [][]byte) {
		for event := range inputChan {
			var outputs [][]byte

			notices, err := node_reader.NewNodeReader(gqlClient).GetNoticesByInputIndex(ctx, int(event.InputIndex.Int64()))
			if err != nil {
				slog.Error("failed to get notices by input index", "error", err)
			} else if len(notices) > 0 {
				outputs = append(outputs, notices...)
			} else {
				slog.Warn("no notices to process", "inputIndex", event.InputIndex)
			}

			vouchers, err := node_reader.NewNodeReader(gqlClient).GetVouchersByInputIndex(ctx, int(event.InputIndex.Int64()))
			if err != nil {
				slog.Error("failed to get vouchers by input index", "error", err)
			} else if len(vouchers) > 0 {
				outputs = append(outputs, vouchers...)
			} else {
				slog.Warn("no vouchers to process", "inputIndex", event.InputIndex)
			}

			if len(outputs) > 0 {
				outputsChan <- outputs
			}
		}
	}(ctx, gqlClient, inputChan, outputChan)

	for {
		select {
		case outputs := <-outputChan:
			handleOutput(ethClient, opts, outputs)
		case <-ctx.Done():
			slog.Info("shutting down due to canceled context")
			return
		}
	}
}

func handleOutput(ethClient *ethclient.Client, opts *bind.TransactOpts, outputs [][]byte) {
	instance, err := coprocessor_contracts.NewCoprocessorContracts(common.HexToAddress(cfg.COPROCESSOR_CALLER_MOCK_ADDRESS), ethClient)
	if err != nil {
		slog.Error("failed to create coprocessor instance", "error", err)
		return
	}
	tx, err := instance.CoprocessorCallbackOutputsOnly(opts, [32]byte{}, [32]byte{}, outputs)
	if err != nil {
		slog.Error("failed to call coprocessor callback function", "error", err, "outputs", outputs)
		return
	}
	slog.Info("outputs executed", "tx", tx.Hash().Hex(), "outputs", outputs)
}
