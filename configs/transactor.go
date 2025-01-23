package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SetupTransactor() (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(os.Getenv("ANVIL_HTTP_URL"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to blockchain: %v", err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain id: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(os.Getenv("ANVIL_PRIVATE_KEY"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load private key: %v", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	return client, opts, nil
}

func SetupTransactorWS() (*ethclient.Client, error) {
	client, err := ethclient.Dial(os.Getenv("ANVIL_WS_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %v", err)
	}
	return client, nil
}
