package rpc

import (
	"context"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/celestiaorg/go-square/v2/share"

	"github.com/stretchr/testify/require"
)

func TestNewCelestiaClient(t *testing.T) {
	if shouldSkipCelestiaTests() {
		t.Skip("Skipping as Celestia is not enabled in the test context.")
	}

	newTestCelestiaClient(t)
}

func shouldSkipCelestiaTests() bool {
	if celestiaEnabled, err := strconv.ParseBool(os.Getenv("CELESTIA_ENABLED")); err == nil {
		return !celestiaEnabled
	}

	return true
}

func getTestCelestiaNamespace() (*share.Namespace, error) {
	namespaceValueString := strings.Replace(os.Getenv("CELESTIA_NAMESPACE"), "0x", "", -1)
	namespaceValue, err := hex.DecodeString(namespaceValueString)
	if err != nil {
		return nil, err
	}

	namespace, err := share.NewV0Namespace(namespaceValue)
	if err != nil {
		return nil, err
	}

	return &namespace, nil
}

func newTestCelestiaClient(t *testing.T) *CelestiaClient {
	namespace, err := getTestCelestiaNamespace()

	require.Nil(t, err)

	client, err := NewCelestiaClient(context.Background(), &CelestiaConfig{
		Enabled:   true,
		Endpoint:  os.Getenv("CELESTIA_ENDPOINT"),
		AuthToken: os.Getenv("CELESTIA_AUTH_TOKEN"),
		Namespace: namespace,
	}, 5*time.Second)

	require.Nil(t, err)
	require.NotNil(t, client)

	return client
}
