//go:build wireinject
// +build wireinject

package root

import (
	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/configs"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/infra/input_reader"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/infra/node_reader"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/usecase"
	"github.com/google/wire"
)

var setTransactorProvider = wire.NewSet(
	configs.SetupTransactorWS,
	input_reader.NewInputReader,
)

var setFindOutputsByIdUseCase = wire.NewSet(
	genqlient.NewClient,
	node_reader.NewNodeReader,
	wire.Bind(new(node_reader.NodeReaderRepository), new(*node_reader.NodeReader)),
	usecase.NewFindOutputsByIdUseCase,
)

func NewTransactor() (*input_reader.InputReader, error) {
	wire.Build(setTransactorProvider)
	return nil, nil
}

func NewFindOutputsByIdUseCase(graphqlUrl string, httpClient genqlient.Doer) (*usecase.FindOutputsByIdUseCase, error) {
	wire.Build(setFindOutputsByIdUseCase)
	return nil, nil
}
