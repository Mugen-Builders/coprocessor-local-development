package node_reader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	genqlient "github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/coprocessor-local-solver/configs"
)

var (
	ErrNoNoticesFound  = errors.New("no notices found")
	ErrNoVouchersFound = errors.New("no vouchers found")
)

type NodeReader struct {
	Client genqlient.Client
}

func NewNodeReader(client genqlient.Client) *NodeReader {
	configs.ConfigureLog(slog.LevelInfo)
	return &NodeReader{
		Client: client,
	}
}

func (r *NodeReader) GetNoticesByInputIndex(ctx context.Context, index int) ([][]byte, error) {
	err := waitForInput(ctx, r.Client, index)
	if err != nil {
		slog.Error("failed to wait for input", "error", err)
		os.Exit(1)
	}

	res, err := getNoticesByInputIndex(ctx, r.Client, index)
	if err != nil {
		return nil, err
	}

	if len(res.Input.Notices.Edges) == 0 {
		return nil, ErrNoNoticesFound
	}

	abiJSON := `[{"inputs":[{"internalType":"bytes","name":"payload","type":"bytes"}],"name":"Notice","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		slog.Error("failed to parse abi", "error", err)
		os.Exit(1)
	}

	outputs := make([][]byte, len(res.Input.Notices.Edges))
	for i, edge := range res.Input.Notices.Edges {
		payload, err := abiInterface.Pack("Notice", common.Hex2Bytes(edge.Node.Payload[2:]))
		if err != nil {
			return nil, err
		}
		outputs[i] = payload
	}
	return outputs, nil
}

func (r *NodeReader) GetVouchersByInputIndex(ctx context.Context, index int) ([][]byte, error) {
	err := waitForInput(ctx, r.Client, index)
	if err != nil {
		slog.Error("failed to wait for input", "error", err)
		os.Exit(1)
	}
	res, err := getVouchersByInputIndex(ctx, r.Client, index)
	if err != nil {
		return nil, err
	}
	if len(res.Input.Vouchers.Edges) == 0 {
		return nil, ErrNoVouchersFound
	}

	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create address type: %w", err)
	}

	bytesType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create bytes type: %w", err)
	}

	args := abi.Arguments{
		{
			Type: addressType,
		},
		{
			Type: bytesType,
		},
	}

	outputs := make([][]byte, len(res.Input.Vouchers.Edges))
	for i, edge := range res.Input.Vouchers.Edges {
		payload, err := args.Pack(common.HexToAddress(edge.Node.Destination), common.Hex2Bytes(edge.Node.Payload[2:]))
		if err != nil {
			return nil, err
		}
		outputs[i] = payload
	}

	return outputs, nil
}

func waitForInput(ctx context.Context, client genqlient.Client, index int) error {
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
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
