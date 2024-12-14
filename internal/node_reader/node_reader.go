package node_reader

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
	"github.com/henriquemarlon/coprocessor-local-solver/configs"
	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common"
)

type NodeReader struct {
	Client genqlient.Client
}

func NewNodeReader(client genqlient.Client) *NodeReader {
	return &NodeReader{
		Client: client,
	}
}

func (r *NodeReader) GetNoticesByInputIndex(ctx context.Context, index int) ([][]byte, error) {
	configs.ConfigureLog(slog.LevelInfo)
	err := waitForInput(ctx, r.Client, index)
	if err != nil {
		slog.Error("Failed to wait for input", "error", err)
		os.Exit(1)
	}

	res, err := getNoticesByInputIndex(ctx, r.Client, index)
	if err != nil {
		return nil, err
	}
	outputs := make([][]byte, len(res.Input.Notices.Edges))
	for _, edge := range res.Input.Notices.Edges {
		noticeBytesPayload := common.Hex2Bytes(edge.Node.GetPayload())
		outputs = append(outputs, noticeBytesPayload)
		slog.Info("Received notice", "index", edge.Node.GetIndex(), "payload", edge.Node.GetPayload())
	}
	return outputs, nil
}

func waitForInput(ctx context.Context, client genqlient.Client, index int) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		result, err := getInputStatus(ctx, client, index)
		if err != nil && !strings.Contains(err.Error(), "input not found") {
			return fmt.Errorf("failed to get input status: %w", err)
		}
		if result.Input.Status == CompletionStatusAccepted {
			return nil
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
