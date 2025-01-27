package usecase

import (
	"context"
	"errors"

	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/Mugen-Builders/cartesi-coprocessor-nonodox/internal/infra/node_reader"
)

var (
	ErrNoOutputsFound = errors.New("no outputs found")
)

type FindOutputsByIdUseCase struct {
	Client               genqlient.Client
	NodeReaderRepository node_reader.NodeReaderRepository
}

func NewFindOutputsByIdUseCase(client genqlient.Client, nodeReaderRepository node_reader.NodeReaderRepository) *FindOutputsByIdUseCase {
	return &FindOutputsByIdUseCase{
		Client:               client,
		NodeReaderRepository: nodeReaderRepository,
	}
}

func (r *FindOutputsByIdUseCase) Execute(ctx context.Context, index int) ([][]byte, error) {
	var outputs [][]byte
	outputs, err := r.NodeReaderRepository.GetOutputsByInputIndex(ctx, index)
	if err != nil {
		if errors.Is(err, node_reader.ErrNoNoticesFound) {
			return nil, ErrNoOutputsFound
		}
		return nil, err
	}
	return outputs, nil
}
