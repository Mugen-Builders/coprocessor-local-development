package evm_reader

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/henriquemarlon/coprocessor-local-solver/configs"
	"github.com/henriquemarlon/coprocessor-local-solver/pkg/rollups_contracts"
)

type EVMReader struct {
	Client          *ethclient.Client
	InputBoxAddress common.Address
}

func NewEVMReader(client *ethclient.Client, inputBoxAddress common.Address) (*EVMReader, error) {
	return &EVMReader{
		Client:          client,
		InputBoxAddress: inputBoxAddress,
	}, nil
}

func (r *EVMReader) GetInputAddedEvents(out chan<- rollups_contracts.IInputBoxInputAdded) error {
	configs.ConfigureLog(slog.LevelInfo)
	iinputBox, err := rollups_contracts.NewIInputBox(r.InputBoxAddress, r.Client)
	if err != nil {
		return err
	}

	itr, err := iinputBox.FilterInputAdded(nil, nil, nil)
	if err != nil {
		return err
	}
	defer itr.Close()

	var events []rollups_contracts.IInputBoxInputAdded
	for itr.Next() {
		inputAddedEvent := itr.Event
		events = append(events, *inputAddedEvent)
		slog.Info("Received input added event", "dapp", inputAddedEvent.Dapp, "index", inputAddedEvent.InputIndex, "sender", inputAddedEvent.Sender, "input", common.Bytes2Hex(inputAddedEvent.Input))
	}
	if err := itr.Error(); err != nil {
		return err
	}
	for _, event := range events {
		out <- event
	}
	return nil
}
