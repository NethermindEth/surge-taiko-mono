package builder

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/testutils"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"

	"github.com/stretchr/testify/require"
)

func TestCelestiaBuildPacaya(t *testing.T) {
	if testutils.ShouldSkipCelestiaTests() {
		t.Skip("Skipping as Celestia is not enabled in the test context.")
	}

	namespace, err := testutils.GetTestCelestiaNamespace()
	require.Nil(t, err)

	proposerPrivateKey, _ := crypto.GenerateKey()

	celestiaTxBuilder := NewCelestiaTransactionBuilder(
		&rpc.Client{
			CelestiaDA: testutils.NewTestCelestiaClient(t, namespace),
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

	_, err = celestiaTxBuilder.BuildPacaya(context.Background(), txsToPropose, nil, nil, common.Hash{}, common.Big0)

	require.Nil(t, err)
}
