package rpc

import (
	"github.com/celestiaorg/go-square/v2/share"
)

// CelestiaConfig contains all configs which will be used to initializing a Celestia RPC client.
type CelestiaConfig struct {
	Enabled   bool
	Endpoint  string
	AuthToken string
	Namespace share.Namespace
}
