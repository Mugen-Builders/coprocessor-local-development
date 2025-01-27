//go:build wireinject
// +build wireinject

package root

import (
	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/infra/evm_reader"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/infra/node_reader"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/usecase"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/wire"
)

var setInputReader = wire.NewSet(
	configs.SetupTransactorWS,
	evm_reader.NewInputReader,
)

var setTaskReader = wire.NewSet(
	configs.SetupTransactorWS,
	evm_reader.NewTaskReader,
)

var setFindOutputsByIdUseCase = wire.NewSet(
	genqlient.NewClient,
	node_reader.NewNodeReader,
	wire.Bind(new(node_reader.NodeReaderRepository), new(*node_reader.NodeReader)),
	usecase.NewFindOutputsByIdUseCase,
)

func NewInputReader() (*evm_reader.InputReader, error) {
	wire.Build(setInputReader)
	return nil, nil
}

func NewTaskReader(mockCoprocessorAddress common.Address) (*evm_reader.TaskReader, error) {
	wire.Build(setTaskReader)
	return nil, nil
}

func NewFindOutputsByIdUseCase(graphqlUrl string, httpClient genqlient.Doer) (*usecase.FindOutputsByIdUseCase, error) {
	wire.Build(setFindOutputsByIdUseCase)
	return nil, nil
}
