package input_reader

import (
	"context"
	"log/slog"
	"math/big"
	"sync"

	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/rollups_contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type InputReader struct {
	Client          *ethclient.Client
	InputBoxAddress common.Address
	mu              sync.Mutex
	lastIndex       *big.Int
}

func NewInputReader(client *ethclient.Client) (*InputReader, error) {
	configs.ConfigureLog(slog.LevelInfo)
	return &InputReader{
		Client:          client,
		InputBoxAddress: common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768"),
		lastIndex:       big.NewInt(-1),
	}, nil
}

func (r *InputReader) GetInputAddedEvents(ctx context.Context, out chan<- rollups_contracts.IInputBoxInputAdded) error {
	iinputBox, err := rollups_contracts.NewIInputBox(r.InputBoxAddress, r.Client)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			itr, err := iinputBox.FilterInputAdded(nil, nil, nil)
			if err != nil {
				slog.Error("failed to filter events", "error", err)
				return err
			}
			defer itr.Close()

			for itr.Next() {
				event := itr.Event
				r.mu.Lock()
				if r.lastIndex.Cmp(event.InputIndex) < 0 {
					r.lastIndex.Set(event.InputIndex)
					out <- *event
				}
				r.mu.Unlock()
			}

			if err := itr.Error(); err != nil {
				slog.Error("error iterating through events", "error", err)
				return err
			}
		}
	}
}
