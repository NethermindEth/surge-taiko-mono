package rpc

import (
	"errors"

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
}

// NewCelestiaClient creates a new CelestiaClient.
func NewCelestiaClient(cfg *CelestiaConfig) (*CelestiaClient, error) {
	if cfg.Endpoint == "" || cfg.AuthToken == "" || cfg.Namespace == nil {
		return nil, errors.New("endpoint, authentication token, or namespace is empty")
	}

	return &CelestiaClient{
		Endpoint:  cfg.Endpoint,
		AuthToken: cfg.AuthToken,
		Namespace: cfg.Namespace,
	}, nil
}
