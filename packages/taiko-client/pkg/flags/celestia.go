package flags

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/celestiaorg/go-square/v2/share"
	"github.com/urfave/cli/v2"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// InitProposerCelestiaConfigsFromCli initializes the Celestia RPC client configs from the proposer command line flags.
func InitProposerCelestiaConfigsFromCli(c *cli.Context) (*rpc.CelestiaConfig, error) {
	namespaceValueString := strings.Replace(c.String(flags.CelestiaNamespace.Name), "0x", "", -1)
	namespaceValue, err := hex.DecodeString(namespaceValueString)
	if err != nil {
		return nil, fmt.Errorf("invalid Celestia namespace: %s", c.String(flags.CelestiaNamespace.Name))
	}

	namespace, err := share.NewV0Namespace(namespaceValue)
	if err != nil {
		return nil, fmt.Errorf("invalid Celestia namespace: %s", c.String(flags.CelestiaNamespace.Name))
	}

	return &rpc.CelestiaConfig{
		Enabled:   c.Bool(flags.CelestiaEnabled.Name),
		Endpoint:  c.String(flags.CelestiaEndpoint.Name),
		AuthToken: c.String(flags.CelestiaAuthToken.Name),
		Namespace: &namespace,
	}, nil
}

// InitProposerCelestiaConfigsFromCli initializes the Celestia RPC client configs from the driver command line flags.
func InitDriverCelestiaConfigsFromCli(c *cli.Context) (*rpc.CelestiaConfig, error) {
	return &rpc.CelestiaConfig{
		Enabled:   c.Bool(flags.CelestiaEnabled.Name),
		Endpoint:  c.String(flags.CelestiaEndpoint.Name),
		AuthToken: c.String(flags.CelestiaAuthToken.Name),
	}, nil
}
