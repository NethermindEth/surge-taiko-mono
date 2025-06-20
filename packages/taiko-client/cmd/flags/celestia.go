package flags

import (
	"github.com/urfave/cli/v2"
)

var (
	celestiaCategory = "CELESTIA DATA AVAILABILITY LAYER"
)

// Flags used by Celestia as as an alternative DA layer.
var (
	CelestiaEnabled = &cli.BoolFlag{
		Name:     "celestia.enabled",
		Usage:    "Enable Celestia as an alternative DA layer",
		Category: celestiaCategory,
		Value:    false,
		EnvVars:  []string{"CELESTIA_ENABLED"},
	}
	CelestiaEndpoint = &cli.StringFlag{
		Name:     "celestia.endpoint",
		Usage:    "RPC endpoint of a Celestia node",
		Category: celestiaCategory,
		EnvVars:  []string{"CELESTIA_ENDPOINT"},
	}
	CelestiaAuthToken = &cli.StringFlag{
		Name:     "celestia.authToken",
		Usage:    "Authentication token of a Celestia node",
		Category: celestiaCategory,
		EnvVars:  []string{"CELESTIA_AUTH_TOKEN"},
	}
	CelestiaNamespace = &cli.StringFlag{
		Name:     "celestia.Namespace",
		Usage:    "Namespace (in hex format) of the Celestia-based alternative DA layer",
		Category: celestiaCategory,
		EnvVars:  []string{"CELESTIA_NAMESPACE"},
	}
)

var CelestiaFlags = []cli.Flag{
	CelestiaEnabled,
	CelestiaEndpoint,
	CelestiaAuthToken,
	CelestiaNamespace,
}
