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
	"sync"
	"syscall"
	"time"

	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/usecase"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/coprocessor_contracts"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/rollups_contracts"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/tools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	configPath string
	cfg        *configs.Config
	Cmd        = &cobra.Command{
		Use:   "nonodox",
		Short: "nonodox is a tool for local development of Cartesi coprocessor applications",
		Long:  "This tool listens for input events, processes notices linked to an input, and executes them on-chain (anvil).",
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

var startupMessage = `
NoNodoX local development tool started
GraphQL server pooling at GRAPHQL_URL

Press Ctrl+C to stop the application.
`

func run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if _, err := exec.LookPath("nonodo"); err != nil {
		slog.Error("nonodo binary not available, please install this and try again", "error", err)
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	group, ctx := errgroup.WithContext(ctx)

	nonodo := exec.CommandContext(
		ctx,
		"nonodo",
		"--disable-devnet",
		"--rpc-url", cfg.AnvilWsURL,
		"--contracts-input-box-block", cfg.AnvilInputBoxBlock,
	)
	notifyWriter := tools.NewNotifyWriter(io.Discard, "nonodo: ready")
	nonodo.Stdout = notifyWriter

	errCh := make(chan error, 1)

	go func() {
		if err := nonodo.Run(); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-notifyWriter.Ready():
		time.Sleep(2 * time.Second)
		message := strings.NewReplacer(
			"GRAPHQL_URL", cfg.GraphQLURL,
		).Replace(startupMessage)
		fmt.Println(message)
	case err := <-errCh:
		if err != nil {
			slog.Error("Error running nonodo", "error", err)
			cancel()
		}
	case <-ctx.Done():
		slog.Info("Context canceled before nonodo became ready")
	}

	chann1 := make(chan coprocessor_contracts.MockCoprocessorTaskIssued, 100)
	chann2 := make(chan rollups_contracts.IInputBoxInputAdded, 100)
	chann3 := make(chan struct {
		Event       rollups_contracts.IInputBoxInputAdded
		PayloadHash common.Hash
		MachineHash common.Hash
		Callback    common.Address
	}, 100)
	chann4 := make(chan struct {
		Event       rollups_contracts.IInputBoxInputAdded
		PayloadHash common.Hash
		MachineHash common.Hash
		Callback    common.Address
		Outputs     [][]byte
	}, 100)

	defer close(chann1)
	defer close(chann2)
	defer close(chann3)
	defer close(chann4)

	var inputsHash sync.Map

	ethClient, opts, err := configs.SetupTransactor()
	if err != nil {
		slog.Error("Failed to set up transactor", "error", err)
		os.Exit(1)
	}

	group.Go(func() error {
		reader, err := NewTaskReader(common.HexToAddress(cfg.CoprocessorMockAddress))
		if err != nil {
			slog.Error("Failed to create task reader", "error", err)
			cancel()
			return err
		}
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				if err := reader.GetTaskIssuedEvents(ctx, chann1); err != nil {
					slog.Error("Error reading tasks", "error", err)
					cancel()
					return err
				}
			}
		}
	})

	group.Go(func() error {
		instance, err := rollups_contracts.NewIInputBox(common.HexToAddress(cfg.InputBoxAddress), ethClient)
		if err != nil {
			slog.Error("Failed to create inputBox instance", "error", err)
			return err
		}
		for {
			select {
			case task := <-chann1:
				inputsHash.Store(crypto.Keccak256Hash(task.Input), struct {
					MachineHash common.Hash
					Callback    common.Address
				}{
					MachineHash: task.MachineHash,
					Callback:    task.Callback,
				})

				_, err := instance.AddInput(opts, common.HexToAddress(cfg.DappAddress), task.Input)
				if err != nil {
					slog.Error("Failed to call addInput function", "error", err)
					continue
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	group.Go(func() error {
		reader, err := NewInputReader()
		if err != nil {
			slog.Error("Failed to create input reader", "error", err)
			cancel()
			return err
		}
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				if err := reader.GetInputAddedEvents(ctx, chann2); err != nil {
					slog.Error("Error reading inputs", "error", err)
					cancel()
					return err
				}
			}
		}
	})

	group.Go(func() error {
		for {
			select {
			case event := <-chann2:
				rawData, ok := inputsHash.Load(crypto.Keccak256Hash(event.Input))
				if !ok {
					slog.Warn("Input not found", "input", crypto.Keccak256Hash(event.Input))
				}

				data, ok := rawData.(struct {
					MachineHash common.Hash
					Callback    common.Address
				})
				if !ok {
					slog.Error("Failed to cast input data")
				}

				chann3 <- struct {
					Event       rollups_contracts.IInputBoxInputAdded
					PayloadHash common.Hash
					MachineHash common.Hash
					Callback    common.Address
				}{
					Event:       event,
					PayloadHash: crypto.Keccak256Hash(event.Input),
					MachineHash: data.MachineHash,
					Callback:    data.Callback,
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	group.Go(func() error {
		findOutputsByIdUseCase, err := NewFindOutputsByIdUseCase(cfg.GraphQLURL, nil)
		if err != nil {
			slog.Error("Failed to setup node reader", "error", err)
			return err
		}
		for {
			select {
			case data := <-chann3:
				time.Sleep(1 * time.Second)

				outputs, err := findOutputsByIdUseCase.Execute(ctx, int(data.Event.InputIndex.Int64()))
				if err != nil {
					if err == usecase.ErrNoOutputsFound {
						slog.Warn("No outputs found", "inputIndex", data.Event.InputIndex)
						continue
					}
					slog.Error("Failed to find outputs by ID", "error", err)
				}

				if len(outputs) > 0 {
					chann4 <- struct {
						Event       rollups_contracts.IInputBoxInputAdded
						PayloadHash common.Hash
						MachineHash common.Hash
						Callback    common.Address
						Outputs     [][]byte
					}{
						Event:       data.Event,
						PayloadHash: data.PayloadHash,
						MachineHash: data.MachineHash,
						Callback:    data.Callback,
						Outputs:     outputs,
					}
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	group.Go(func() error {
		instance, err := coprocessor_contracts.NewMockCoprocessor(common.HexToAddress(cfg.CoprocessorMockAddress), ethClient)
		if err != nil {
			slog.Error("Failed to create coprocessor instance", "error", err)
			return err
		}
		for {
			select {
			case output := <-chann4:
				tx, err := instance.SolverCallbackOutputsOnly(
					opts,
					output.MachineHash,
					output.PayloadHash,
					output.Outputs,
					output.Callback,
				)
				if err != nil {
					if strings.Contains(err.Error(), "execution reverted") {
						slog.Warn("Execution reverted, make sure the CoprocessorAdapter address is correct")
					}
					slog.Error("Failed to call coprocessor callback", "error", err)
					continue
				}
				slog.Info("Outputs executed", "payload hash", output.PayloadHash, "tx", tx.Hash().Hex())
			case <-ctx.Done():
				return nil
			}
		}
	})

	if err := group.Wait(); err != nil {
		slog.Error("Execution error", "error", err)
	}
}
