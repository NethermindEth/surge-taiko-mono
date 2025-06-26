package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/celestiaorg/go-square/v2/share"
	"github.com/filecoin-project/go-jsonrpc"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/celestia"
)

const (
	// https://docs.celestia.org/how-to-guides/submit-data#maximum-blob-size
	AdvisableCelestiaBlobSize = 500000
)

var (
	// Minimum gas price (0.2 TIA --> 200000 utia) https://docs.celestia.org/how-to-guides/submit-data#fee-market-and-mempool
	MinimumGasPrice = big.NewInt(200000)
)

type Blob struct {
}

// CelestiaConfig contains all configs which will be used to initializing a Celestia RPC client.
type CelestiaConfig struct {
	Enabled   bool
	Endpoint  string
	AuthToken string
	Namespace *share.Namespace
}

// CelestiaClient is a client for the Celestia node.
type CelestiaClient struct {
	Endpoint   string
	AuthHeader http.Header
	Namespace  *share.Namespace

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

	authHeader := http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", cfg.AuthToken)}}

	client := celestia.CelestiaHeaderHandler{}
	closer, err := jsonrpc.NewClient(ctx, cfg.Endpoint, celestia.CelestiaHeaderNamespace, &client, authHeader)
	if err != nil {
		return nil, err
	}
	defer closer()

	// Get network head to verify connectivity
	if err := client.NetworkHead(ctx); err != nil {
		return nil, err
	}

	return &CelestiaClient{
		Endpoint:   cfg.Endpoint,
		AuthHeader: authHeader,
		Namespace:  cfg.Namespace,
		Timeout:    timeoutVal,
	}, nil
}

func (c *CelestiaClient) CheckBalance(ctx context.Context) (bool, error) {
	ctxWithTimeout, cancel := CtxWithTimeoutOrDefault(ctx, c.Timeout)
	defer cancel()

	client := celestia.CelestiaStateHandler{}
	closer, err := jsonrpc.NewClient(ctxWithTimeout, c.Endpoint, celestia.CelestiaStateNamespace, &client, c.AuthHeader)
	if err != nil {
		return false, err
	}
	defer closer()

	balance, err := client.Balance(ctx)
	if err != nil {
		return false, err
	}

	amount, success := new(big.Int).SetString(balance.Amount, 0)
	return success && (amount.Cmp(MinimumGasPrice) > 0), nil
}

func (c *CelestiaClient) Submit(ctx context.Context, blobs []*Blob) (uint64, error) {
	// TODO: Resolved the celestia-node dependencies issues or write our own minimalistic client
	/*
		ctxWithTimeout, cancel := CtxWithTimeoutOrDefault(ctx, c.Timeout)
		defer cancel()

		client, err := client.NewClient(ctxWithTimeout, c.Endpoint, c.AuthToken)
		if err != nil {
			return 0, err
		}
		defer client.Close()

		options := state.NewTxConfig()

		height, err := client.Blob.Submit(ctxWithTimeout, blobs, options)
		if err != nil {
			return 0, err
		}

		return height, nil
	*/
	return 0, nil
}
