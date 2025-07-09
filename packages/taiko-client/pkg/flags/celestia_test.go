package flags

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
)

var (
	Endpoint          = "http://localhost:26658"
	AuthToken         = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiXX0"
	Namespace         = "0xDEADBEEF"
	ExpectedNamespace = "00000000000000000000000000000000000000000000000000deadbeef"
)

func TestInitProposerCelestiaConfigsFromCliCelestiaDisabled(t *testing.T) {
	if isCelestiaEnabledInEnv() {
		t.Skip("Skipping as Celestia is enabled at the environment level.")
	}

	app := cli.NewApp()
	app.Flags = append(app.Flags, flags.CelestiaProposerFlags...)

	app.Action = func(cliCtx *cli.Context) error {
		celestiaConfigs, err := InitProposerCelestiaConfigsFromCli(cliCtx)

		require.NoError(t, err)
		require.Equal(t, false, celestiaConfigs.Enabled)

		return nil
	}

	require.Nil(t, app.Run([]string{
		"TestNewConfigFromCliContext",
	}))
}

func TestProposerInitCelestiaConfigsFromCliCelestiaEnabled(t *testing.T) {
	app := cli.NewApp()
	app.Flags = append(app.Flags, flags.CelestiaProposerFlags...)

	app.Action = func(cliCtx *cli.Context) error {
		celestiaConfigs, err := InitProposerCelestiaConfigsFromCli(cliCtx)

		require.NoError(t, err)
		require.Equal(t, true, celestiaConfigs.Enabled)
		require.Equal(t, Endpoint, celestiaConfigs.Endpoint)
		require.Equal(t, AuthToken, celestiaConfigs.AuthToken)
		require.Equal(t, ExpectedNamespace, celestiaConfigs.Namespace.String())

		return nil
	}

	require.Nil(t, app.Run([]string{
		"TestNewConfigFromCliContext",
		"--" + flags.CelestiaEnabled.Name,
		"--" + flags.CelestiaEndpoint.Name, Endpoint,
		"--" + flags.CelestiaAuthToken.Name, AuthToken,
		"--" + flags.CelestiaNamespace.Name, Namespace,
	}))
}

func TestInitDriverCelestiaConfigsFromCliCelestiaDisabled(t *testing.T) {
	if isCelestiaEnabledInEnv() {
		t.Skip("Skipping as Celestia is enabled at the environment level.")
	}

	app := cli.NewApp()
	app.Flags = append(app.Flags, flags.CelestiaDriverFlags...)

	app.Action = func(cliCtx *cli.Context) error {
		celestiaConfigs, err := InitDriverCelestiaConfigsFromCli(cliCtx)

		require.NoError(t, err)
		require.Equal(t, false, celestiaConfigs.Enabled)

		return nil
	}

	require.Nil(t, app.Run([]string{
		"TestNewConfigFromCliContext",
	}))
}

func TestDriverInitCelestiaConfigsFromCliCelestiaEnabled(t *testing.T) {
	app := cli.NewApp()
	app.Flags = append(app.Flags, flags.CelestiaDriverFlags...)

	app.Action = func(cliCtx *cli.Context) error {
		celestiaConfigs, err := InitDriverCelestiaConfigsFromCli(cliCtx)

		require.NoError(t, err)
		require.Equal(t, true, celestiaConfigs.Enabled)
		require.Equal(t, Endpoint, celestiaConfigs.Endpoint)
		require.Equal(t, AuthToken, celestiaConfigs.AuthToken)

		return nil
	}

	require.Nil(t, app.Run([]string{
		"TestNewConfigFromCliContext",
		"--" + flags.CelestiaEnabled.Name,
		"--" + flags.CelestiaEndpoint.Name, Endpoint,
		"--" + flags.CelestiaAuthToken.Name, AuthToken,
	}))
}

func isCelestiaEnabledInEnv() bool {
	if celestiaEnabled, err := strconv.ParseBool(os.Getenv("CELESTIA_ENABLED")); err == nil {
		return celestiaEnabled
	}

	return false
}
