package evm_reader

import (
	"context"
	"encoding/hex"
	"log/slog"
	"strings"
	"sync"

	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/coprocessor_contracts"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TaskReader struct {
	Client                 *ethclient.Client
	MockCoprocessorAddress common.Address
	mu                     sync.Mutex
	processedTransactions  map[common.Hash]bool // Tracks processed transaction hashes
}

func NewTaskReader(client *ethclient.Client, mockCoprocessorAddress common.Address) *TaskReader {
	configs.ConfigureLog(slog.LevelInfo)
	return &TaskReader{
		Client:                 client,
		MockCoprocessorAddress: mockCoprocessorAddress,
		processedTransactions:  make(map[common.Hash]bool),
	}
}

func (r *TaskReader) GetTaskIssuedEvents(ctx context.Context, out chan<- coprocessor_contracts.MockCoprocessorTaskIssued) error {
	// Define the filter query for logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{r.MockCoprocessorAddress},
	}

	logsChan := make(chan types.Log)
	subscription, err := r.Client.SubscribeFilterLogs(ctx, query, logsChan)
	if err != nil {
		slog.Error("Failed to subscribe to logs", "error", err)
		return err
	}
	defer subscription.Unsubscribe()

	// Parse the contract ABI
	contractAbi, err := abi.JSON(strings.NewReader(string(coprocessor_contracts.MockCoprocessorMetaData.ABI)))
	if err != nil {
		slog.Error("Failed to parse contract ABI", "error", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-subscription.Err():
			slog.Error("Subscription error", "error", err)
			return err
		case vLog := <-logsChan:
			txHash := vLog.TxHash
			r.mu.Lock()
			if r.processedTransactions[txHash] {
				slog.Info("Skipping already processed transaction", "txHash", txHash.Hex())
				r.mu.Unlock()
				continue
			}
			r.processedTransactions[txHash] = true
			r.mu.Unlock()

			event := coprocessor_contracts.MockCoprocessorTaskIssued{}
			err := contractAbi.UnpackIntoInterface(&event, "TaskIssued", vLog.Data)
			if err != nil {
				slog.Error("Failed to unpack event data", "error", err)
				continue
			}

			select {
			case out <- event:
				slog.Info("Receiving TaskIssued", "machineHash", common.Hash(event.MachineHash), "input", hex.EncodeToString(event.Input), "callback", event.Callback)
			case <-ctx.Done():
				return nil
			}
		}
	}
}
