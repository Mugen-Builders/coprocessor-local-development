package node_reader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNoNoticesFound  = errors.New("no notices found")
	ErrNoVouchersFound = errors.New("no vouchers found")
)

type NodeReaderRepository interface {
	GetOutputsByInputIndex(ctx context.Context, index int) ([][]byte, error)
}

type NodeReader struct {
	Client genqlient.Client
}

func NewNodeReader(client genqlient.Client) *NodeReader {
	return &NodeReader{
		Client: client,
	}
}

func (r *NodeReader) GetOutputsByInputIndex(ctx context.Context, index int) ([][]byte, error) {
	err := _waitForInput(ctx, r.Client, index)
	if err != nil {
		return nil, err
	}

	res, err := getOutputsByInputIndex(ctx, r.Client, index)
	if err != nil {
		return nil, err
	}

	if len(res.Input.Notices.Edges) == 0 {
		return nil, ErrNoNoticesFound
	}
	slog.Info("Notices", "count", len(res.Input.Notices.Edges), "inputIndex", index, "payload", res.Input.Notices.Edges[0].Node.Payload)

	outputs := make([][]byte, len(res.Input.Notices.Edges))

	noticeAbiJSON := `[{"inputs":[{"internalType":"bytes","name":"payload","type":"bytes"}],"name":"Notice","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	noticeAbiInterface, err := abi.JSON(strings.NewReader(noticeAbiJSON))
	if err != nil {
		return nil, err
	}

	for i, edge := range res.Input.Notices.Edges {
		notice, err := noticeAbiInterface.Pack("Notice", common.Hex2Bytes(edge.Node.Payload[2:]))
		if err != nil {
			return nil, err
		}
		outputs[i] = notice
	}

	return outputs, nil
}

func _waitForInput(ctx context.Context, client genqlient.Client, index int) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		result, err := getInputStatus(ctx, client, index)
		if err != nil && !strings.Contains(err.Error(), "input not found") {
			return fmt.Errorf("failed to get input status: %w", err)
		}
		if result.Input.Status == CompletionStatusAccepted {
			return nil
		}
		if result.Input.Status == CompletionStatusRejected {
			return fmt.Errorf("input rejected")
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
