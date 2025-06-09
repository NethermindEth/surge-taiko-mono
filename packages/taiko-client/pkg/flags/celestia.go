package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// InitCelestiaConfigsFromCli initializes the Celestia RPC client configs from the command line flags.
func InitCelestiaConfigsFromCli(c *cli.Context) *rpc.CelestiaConfig {
	return &rpc.CelestiaConfig{
		Enabled:   c.Bool(flags.CelestiaEnabled.Name),
		Endpoint:  c.String(flags.CelestiaEndpoint.Name),
		AuthToken: c.String(flags.CelestiaAuthToken.Name),
		Namespace: c.String(flags.CelestiaNamespace.Name),
	}
}
