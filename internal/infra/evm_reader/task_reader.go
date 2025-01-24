package evm_reader

import (
	"context"
	"log/slog"
	"math/big"
	"sync"

	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/coprocessor_contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TaskReader struct {
	Client                 *ethclient.Client
	MockCoprocessorAddress common.Address
	mu                     sync.Mutex
	lastIndex              *big.Int
}

func NewTaskReader(client *ethclient.Client) *TaskReader {
	configs.ConfigureLog(slog.LevelInfo)
	return &TaskReader{
		Client:                 client,
		MockCoprocessorAddress: common.HexToAddress("0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE"),
		lastIndex:              big.NewInt(-1),
	}
}

func (r *TaskReader) GetTaskIssuedEvents(ctx context.Context, out chan<- coprocessor_contracts.MockCoprocessorTaskIssued) error {
	mockCoprocessor, err := coprocessor_contracts.NewMockCoprocessor(r.MockCoprocessorAddress, r.Client)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			itr, err := mockCoprocessor.FilterTaskIssued(nil)
			if err != nil {
				slog.Error("failed to filter events", "error", err)
				return err
			}
			defer itr.Close()

			for itr.Next() {
				event := itr.Event
				r.mu.Lock()
				if r.lastIndex.Cmp(big.NewInt(int64(event.Raw.Index))) < 0 {
					r.lastIndex.Set(big.NewInt(int64(event.Raw.Index)))
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
