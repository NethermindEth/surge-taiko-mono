package flags

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
)

var (
	Endpoint          = "http://localhost:26658"
	AuthToken         = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiXX0.cSrJjpfUdTNFtzGho69V0D_8kyECn9Mzv8ghJSpKRDE"
	Namespace         = "0xDEADBEEF"
	ExpectedNamespace = "00000000000000000000000000000000000000000000000000deadbeef"
)

func TestInitCelestiaConfigsFromCliCelestiaDisabled(t *testing.T) {
	app := cli.NewApp()
	app.Flags = append(app.Flags, flags.CelestiaFlags...)

	app.Action = func(cliCtx *cli.Context) error {
		celestiaConfigs, err := InitCelestiaConfigsFromCli(cliCtx)

		require.NoError(t, err)
		require.Equal(t, false, celestiaConfigs.Enabled)

		return nil
	}

	app.Run([]string{
		"TestNewConfigFromCliContext",
	})
}

func TestInitCelestiaConfigsFromCliCelestiaEnabled(t *testing.T) {
	app := cli.NewApp()
	app.Flags = append(app.Flags, flags.CelestiaFlags...)

	app.Action = func(cliCtx *cli.Context) error {
		celestiaConfigs, err := InitCelestiaConfigsFromCli(cliCtx)

		require.NoError(t, err)
		require.Equal(t, true, celestiaConfigs.Enabled)
		require.Equal(t, Endpoint, celestiaConfigs.Endpoint)
		require.Equal(t, AuthToken, celestiaConfigs.AuthToken)
		require.Equal(t, ExpectedNamespace, celestiaConfigs.Namespace.String())

		return nil
	}

	app.Run([]string{
		"TestNewConfigFromCliContext",
		"--" + flags.CelestiaEnabled.Name,
		"--" + flags.CelestiaEndpoint.Name, Endpoint,
		"--" + flags.CelestiaAuthToken.Name, AuthToken,
		"--" + flags.CelestiaNamespace.Name, Namespace,
	})
}
