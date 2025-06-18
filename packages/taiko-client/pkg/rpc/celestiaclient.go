package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/go-square/v2/share"
)

// CelestiaConfig contains all configs which will be used to initializing a Celestia RPC client.
type CelestiaConfig struct {
	Enabled   bool
	Endpoint  string
	AuthToken string
	Namespace *share.Namespace
}

// CelestiaClient is a client for the Celestia node.
type CelestiaClient struct {
	Endpoint  string
	AuthToken string
	Namespace *share.Namespace

	Timeout time.Duration
}

// NewCelestiaClient creates a new CelestiaClient.
func NewCelestiaClient(ctx context.Context, cfg *CelestiaConfig, timeout time.Duration) (*CelestiaClient, error) {
	if cfg.Endpoint == "" || cfg.AuthToken == "" || cfg.Namespace == nil {
		return nil, errors.New("endpoint, authentication token, or namespace is empty")
	}

	var timeoutVal = defaultTimeout
	if timeout != 0 {
		timeoutVal = timeout
	}

	client, err := client.NewClient(ctx, cfg.Endpoint, cfg.AuthToken)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Get network head to verify connectivity
	if _, err := client.Header.NetworkHead(ctx); err != nil {
		return nil, err
	}

	return &CelestiaClient{
		Endpoint:  cfg.Endpoint,
		AuthToken: cfg.AuthToken,
		Namespace: cfg.Namespace,
		Timeout:   timeoutVal,
	}, nil
}

func (c *CelestiaClient) CheckBalance(ctx context.Context) (bool, error) {
	ctxWithTimeout, cancel := CtxWithTimeoutOrDefault(ctx, c.Timeout)
	defer cancel()

	client, err := client.NewClient(ctxWithTimeout, c.Endpoint, c.AuthToken)
	if err != nil {
		return false, err
	}
	defer client.Close()

	balance, err := client.State.Balance(ctx)
	if err != nil {
		return false, err
	}

	return balance.Amount > 0, nil
}
