package evm_reader

import (
	"context"
	"log/slog"
	"math/big"
	"strings"
	"sync"

	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/pkg/rollups_contracts"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type InputReader struct {
	Client          *ethclient.Client
	InputBoxAddress common.Address
	mu              sync.Mutex
	processedEvents map[common.Hash]bool // Tracks processed event hashes
}

func NewInputReader(client *ethclient.Client) (*InputReader, error) {
	configs.ConfigureLog(slog.LevelInfo)
	return &InputReader{
		Client:          client,
		InputBoxAddress: common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768"),
		processedEvents: make(map[common.Hash]bool),
	}, nil
}

func (r *InputReader) GetInputAddedEvents(ctx context.Context, out chan<- rollups_contracts.IInputBoxInputAdded) error {
	// Define the filter query for logs
	query := ethereum.FilterQuery{
		Addresses: []common.Address{r.InputBoxAddress},
	}

	logsChan := make(chan types.Log)
	subscription, err := r.Client.SubscribeFilterLogs(ctx, query, logsChan)
	if err != nil {
		slog.Error("Failed to subscribe to logs", "error", err)
		return err
	}
	defer subscription.Unsubscribe()

	// Parse the contract ABI
	contractAbi, err := abi.JSON(strings.NewReader(string(rollups_contracts.IInputBoxMetaData.ABI)))
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
			eventHash := vLog.TxHash
			r.mu.Lock()
			if r.processedEvents[eventHash] {
				slog.Info("Skipping already processed event", "txHash", eventHash.Hex())
				r.mu.Unlock()
				continue
			}
			r.processedEvents[eventHash] = true
			r.mu.Unlock()

			event := rollups_contracts.IInputBoxInputAdded{
				Dapp:       common.HexToAddress(vLog.Topics[1].Hex()),
				InputIndex: new(big.Int).SetBytes(vLog.Topics[2].Bytes()),
			}

			err := contractAbi.UnpackIntoInterface(&event, "InputAdded", vLog.Data)
			if err != nil {
				slog.Error("Failed to unpack event data", "error", err)
				continue
			}

			select {
			case out <- event:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
