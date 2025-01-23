package usecase

import (
	"context"
	"errors"
	"log/slog"

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

	notices, err := r.NodeReaderRepository.GetNoticesByInputIndex(ctx, index)
	if err != nil {
		return nil, err
	} else {
		outputs = append(outputs, notices...)
	}

	vouchers, err := r.NodeReaderRepository.GetVouchersByInputIndex(ctx, index)
	if err != nil {
		slog.Error("failed to get vouchers by input index", "error", err)
	} else {
		outputs = append(outputs, vouchers...)
	}

	if len(outputs) == 0 {
		return nil, ErrNoOutputsFound
	}

	return outputs, nil
}
