package builder

import (
	"context"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/celestiaorg/go-square/v2/share"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"

	"github.com/stretchr/testify/require"
)

func TestCelestiaBuildPacaya(t *testing.T) {
	if shouldSkipCelestiaTests() {
		t.Skip("Skipping as Celestia is not enabled in the test context.")
	}

	proposerPrivateKey, _ := crypto.GenerateKey()

	celestiaTxBuilder := NewCelestiaTransactionBuilder(
		&rpc.Client{
			CelestiaDA: newTestCelestiaClient(t),
		},
		proposerPrivateKey,
		common.Address{},
		common.Address{},
		common.Address{},
		common.Address{},
		0,
		config.NewChainConfig(
			common.Big0,
			0,
			0,
		),
		false,
	)

	var txsToPropose []types.Transactions
	for i := 0; i < 100; i++ {
		txsToPropose = append(txsToPropose, []*types.Transaction{types.NewTransaction(
			uint64(i),
			common.Address{},
			common.Big0,
			0,
			common.Big0,
			nil,
		)})
	}

	_, err := celestiaTxBuilder.BuildPacaya(context.Background(), txsToPropose, nil, nil, common.Hash{})

	require.Nil(t, err)
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

func newTestCelestiaClient(t *testing.T) *rpc.CelestiaClient {
	namespace, err := getTestCelestiaNamespace()

	require.Nil(t, err)

	client, err := rpc.NewCelestiaClient(context.Background(), &rpc.CelestiaConfig{
		Enabled:   true,
		Endpoint:  os.Getenv("CELESTIA_ENDPOINT"),
		AuthToken: os.Getenv("CELESTIA_AUTH_TOKEN"),
		Namespace: namespace,
	}, 5*time.Second)

	require.Nil(t, err)
	require.NotNil(t, client)

	return client
}
