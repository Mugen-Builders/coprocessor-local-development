package evm_reader

import (
	"context"
	"log/slog"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/henriquemarlon/coprocessor-local-solver/configs"
	"github.com/henriquemarlon/coprocessor-local-solver/pkg/rollups_contracts"
)

type EVMReader struct {
	Client          *ethclient.Client
	InputBoxAddress common.Address
	mu              sync.Mutex
	lastIndex       *big.Int
}

func NewEVMReader(client *ethclient.Client, inputBoxAddress common.Address) (*EVMReader, error) {
	configs.ConfigureLog(slog.LevelInfo)
	return &EVMReader{
		Client:          client,
		InputBoxAddress: inputBoxAddress,
		lastIndex:       big.NewInt(-1),
	}, nil
}

func (r *EVMReader) GetInputAddedEvents(ctx context.Context, out chan<- rollups_contracts.IInputBoxInputAdded) error {
	iinputBox, err := rollups_contracts.NewIInputBox(r.InputBoxAddress, r.Client)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled, stopping event monitoring")
			return nil
		default:
			itr, err := iinputBox.FilterInputAdded(nil, nil, nil)
			if err != nil {
				slog.Error("failed to filter inputadded events", "error", err)
				return err
			}
			defer itr.Close()

			for itr.Next() {
				event := itr.Event
				r.mu.Lock()
				if r.lastIndex.Cmp(event.InputIndex) < 0 {
					r.lastIndex.Set(event.InputIndex)
					out <- *event
					slog.Info("new inputadded event sent to the channel", "dapp", event.Dapp, "index", event.InputIndex, "sender", event.Sender)
				}
				r.mu.Unlock()
			}
			if err := itr.Error(); err != nil {
				slog.Error("error iterating inputadded events", "error", err)
				return err
			}
		}
	}
}
