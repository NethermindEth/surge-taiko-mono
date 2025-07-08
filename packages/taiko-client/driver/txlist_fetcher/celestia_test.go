package txlistfetcher

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

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	txListDecompressor "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/txlist_decompressor"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

func TestCelestiaFetchPacaya(t *testing.T) {
	if shouldSkipCelestiaTests() {
		t.Skip("Skipping as Celestia is not enabled in the test context.")
	}

	celestiaBlobHeight := getTestCelestiaBlobHeight()
	if celestiaBlobHeight <= 0 {
		t.Skip("Skipping as Celestia blob for the test was not provided.")
	}

	txListFetcherCelestia := NewCelestiaFetcher(
		&rpc.Client{
			CelestiaDA: newTestCelestiaClient(t),
		},
	)

	txListDecompressor := txListDecompressor.NewTxListDecompressor(
		uint64(240_000_000),
		rpc.BlockMaxTxListBytes,
	)

	namespace, err := getTestCelestiaNamespace()
	require.Nil(t, err)

	metadata := metadata.NewTaikoDataBlockMetadataPacaya(
		&pacaya.TaikoInboxClientBatchProposed{
			Info: pacaya.ITaikoInboxBatchInfo{
				CelestiaBlobParams: pacaya.ITaikoInboxCelestiaBlobParams{
					Height:    celestiaBlobHeight,
					Namespace: namespace.Bytes(),
				},
				BlobByteOffset: 0,
				BlobByteSize:   getTestCelestiaBlobByteSize(),
			},
		},
	)
	meta := metadata.Pacaya()

	txListBytes, err := txListFetcherCelestia.FetchPacaya(context.Background(), meta)
	require.Nil(t, err)

	allTxs := txListDecompressor.TryDecompress(txListBytes, meta.GetCelestiaBlobsHeight() > 0 || len(meta.GetBlobHashes()) != 0)

	require.Greater(t, allTxs.Len(), 0)
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

func getTestCelestiaBlobHeight() uint64 {
	if celestiaBlobHeight, err := strconv.ParseUint(os.Getenv("CELESTIA_BLOB_HEIGHT"), 10, 64); err == nil {
		return celestiaBlobHeight
	}

	return 0
}

func getTestCelestiaBlobByteSize() uint32 {
	if celestiaBlobByteSize, err := strconv.ParseUint(os.Getenv("CELESTIA_BLOB_BYTE_SIZE"), 10, 32); err == nil {
		return uint32(celestiaBlobByteSize)
	}

	return 0
}

func newTestCelestiaClient(t *testing.T) *rpc.CelestiaClient {

	client, err := rpc.NewCelestiaClient(context.Background(), &rpc.CelestiaConfig{
		Enabled:   true,
		Endpoint:  os.Getenv("CELESTIA_ENDPOINT"),
		AuthToken: os.Getenv("CELESTIA_AUTH_TOKEN"),
		Namespace: &share.Namespace{},
	}, 5*time.Second)

	require.Nil(t, err)
	require.NotNil(t, client)

	return client
}
