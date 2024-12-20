package rpc

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
)

var (
	testAddress = common.HexToAddress("0x98f86166571FE624778203d87A8eD6fd84695B79")
)

func TestL2AccountNonce(t *testing.T) {
	client := newTestClientWithTimeout(t)

	nonce, err := client.L2AccountNonce(context.Background(), testAddress, common.Big0)

	require.Nil(t, err)
	require.Zero(t, nonce)
}

func TestGetGenesisL1Header(t *testing.T) {
	client := newTestClient(t)

	header, err := client.GetGenesisL1Header(context.Background())

	require.Nil(t, err)
	require.NotZero(t, header.Number.Uint64())
}

func TestLatestL2KnownL1Header(t *testing.T) {
	client := newTestClient(t)

	header, err := client.LatestL2KnownL1Header(context.Background())

	require.Nil(t, err)
	require.NotZero(t, header.Number.Uint64())
}

func TestL2ParentByBlockId(t *testing.T) {
	client := newTestClient(t)

	header, err := client.L2ParentByBlockID(context.Background(), common.Big1)
	require.Nil(t, err)
	require.Zero(t, header.Number.Uint64())

	_, err = client.L2ParentByBlockID(context.Background(), common.Big2)
	require.NotNil(t, err)
}

func TestL2ExecutionEngineSyncProgress(t *testing.T) {
	client := newTestClient(t)

	progress, err := client.L2ExecutionEngineSyncProgress(context.Background())
	require.Nil(t, err)
	require.NotNil(t, progress)
}

func TestGetProtocolStateVariables(t *testing.T) {
	client := newTestClient(t)
	_, err := client.GetProtocolStateVariables(nil)
	require.Nil(t, err)
}

func TestWaitTillL2ExecutionEngineSyncedNewClient(t *testing.T) {
	client := newTestClient(t)
	err := client.WaitTillL2ExecutionEngineSynced(context.Background())
	require.Nil(t, err)
}

func TestGetSyncedL1SnippetFromAnchor(t *testing.T) {
	client := newTestClient(t)

	l1BlockHash := randomHash()
	l1StateRoot := randomHash()
	l1Height := randomHash().Big().Uint64()
	parentGasUsed := uint32(randomHash().Big().Uint64())

	testAddrPrivKey, err := crypto.ToECDSA(common.Hex2Bytes(encoding.GoldenTouchPrivKey))
	require.Nil(t, err)

	opts, err := bind.NewKeyedTransactorWithChainID(testAddrPrivKey, client.L2.ChainID)
	require.Nil(t, err)

	opts.NoSend = true
	opts.GasLimit = 1_000_000

	tx, err := client.TaikoL2.Anchor(opts, l1BlockHash, l1StateRoot, l1Height, parentGasUsed)
	require.Nil(t, err)

	syncedL1StateRoot,
		syncedL1Height,
		syncedParentGasUsed,
		err := client.getSyncedL1SnippetFromAnchor(tx)
	require.Nil(t, err)
	require.Equal(t, l1StateRoot, syncedL1StateRoot)
	require.Equal(t, l1Height, syncedL1Height)
	require.Equal(t, parentGasUsed, syncedParentGasUsed)
}

func TestWaitTillL2ExecutionEngineSyncedContextErr(t *testing.T) {
	client := newTestClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.WaitTillL2ExecutionEngineSynced(ctx)
	require.ErrorContains(t, err, "context canceled")
}

// randomHash generates a random blob of data and returns it as a hash.
func randomHash() common.Hash {
	var hash common.Hash
	if n, err := rand.Read(hash[:]); n != common.HashLength || err != nil {
		panic(err)
	}
	return hash
}

func TestGetTiers(t *testing.T) {
	client := newTestClient(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tiers, err := client.GetTiers(ctx)

	require.Nil(t, err)
	require.NotNil(t, tiers)
	require.Greater(t, len(tiers), 0, "should have at least one tier")

	// Check for tier with ID 1100 (TIER_TWO_OF_THREE)
	found := false
	for _, tier := range tiers {
		if tier.ID == 1100 {
			found = true
			break
		}
	}
	require.True(t, found, "tier with ID 1100 (TIER_TWO_OF_THREE) should exist")
}
