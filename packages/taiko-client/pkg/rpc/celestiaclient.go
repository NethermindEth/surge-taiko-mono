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
	// https://docs.celestia.org/how-to-guides/submit-data#fee-market-and-mempool
	// https://docs.celestia.org/learn/tia
	MinimalBalancee = big.NewInt(celestia.DefaultMaxGasPrice * 1000000)
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
	Endpoint   string
	AuthHeader http.Header
	Namespace  share.Namespace

	Timeout time.Duration
}

// NewCelestiaClient creates a new CelestiaClient.
func NewCelestiaClient(ctx context.Context, cfg *CelestiaConfig, timeout time.Duration) (*CelestiaClient, error) {
	if cfg.Endpoint == "" || cfg.AuthToken == "" {
		return nil, errors.New("endpoint, authentication token is empty")
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
		Namespace:  *cfg.Namespace,
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
	return success && (amount.Cmp(MinimalBalancee) > 0), nil
}

func (c *CelestiaClient) Submit(ctx context.Context, blobs []*celestia.Blob) (uint64, error) {
	ctxWithTimeout, cancel := CtxWithTimeoutOrDefault(ctx, c.Timeout)
	defer cancel()

	client := celestia.CelestiaBlobHandler{}
	closer, err := jsonrpc.NewClient(ctxWithTimeout, c.Endpoint, celestia.CelestiaBlobNamespace, &client, c.AuthHeader)
	if err != nil {
		return 0, err
	}
	defer closer()

	options := celestia.NewSubmitOptions()

	height, err := client.Submit(ctxWithTimeout, blobs, options)
	if err != nil {
		return 0, err
	}

	return height, nil
}

func (c *CelestiaClient) GetAll(ctx context.Context, height uint64, namespace share.Namespace) ([]*celestia.Blob, error) {
	ctxWithTimeout, cancel := CtxWithTimeoutOrDefault(ctx, c.Timeout)
	defer cancel()

	client := celestia.CelestiaBlobHandler{}
	closer, err := jsonrpc.NewClient(ctxWithTimeout, c.Endpoint, celestia.CelestiaBlobNamespace, &client, c.AuthHeader)
	if err != nil {
		return nil, err
	}
	defer closer()

	blobs, err := client.GetAll(ctxWithTimeout, height, []share.Namespace{namespace})
	if err != nil {
		return nil, err
	}

	return blobs, nil
}
