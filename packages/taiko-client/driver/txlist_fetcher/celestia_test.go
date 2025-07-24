package txlistfetcher

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/celestiaorg/go-square/v2/share"
	"github.com/stretchr/testify/require"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	txListDecompressor "github.com/taikoxyz/taiko-mono/packages/taiko-client/driver/txlist_decompressor"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/testutils"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

func TestCelestiaFetchPacaya(t *testing.T) {
	if testutils.ShouldSkipCelestiaTests() {
		t.Skip("Skipping as Celestia is not enabled in the test context.")
	}

	celestiaBlobHeight := getTestCelestiaBlobHeight()
	if celestiaBlobHeight <= 0 {
		t.Skip("Skipping as Celestia blob for the test was not provided.")
	}

	txListFetcherCelestia := NewCelestiaFetcher(
		&rpc.Client{
			CelestiaDA: testutils.NewTestCelestiaClient(t, &share.Namespace{}),
		},
	)

	txListDecompressor := txListDecompressor.NewTxListDecompressor(
		uint64(240_000_000),
		rpc.BlockMaxTxListBytes,
	)

	namespace, err := testutils.GetTestCelestiaNamespace()
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

	blobUsed := meta.GetCelestiaBlobsHeight() > 0 || len(meta.GetBlobHashes()) != 0
	allTxs := txListDecompressor.TryDecompress(txListBytes, blobUsed)

	require.Greater(t, allTxs.Len(), 0)
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
