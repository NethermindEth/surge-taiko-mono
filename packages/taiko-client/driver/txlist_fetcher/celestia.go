package txlistfetcher

import (
	"context"
	"errors"

	"github.com/celestiaorg/go-square/v2/share"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// CelestiaFetcher is responsible for fetching  txList blob from Celestia.
type CelestiaFetcher struct {
	rpc *rpc.Client
}

// NewCelestiaFetcher creates a new CelestiaFetcher instance based on the given rpc client.
func NewCelestiaFetcher(rpc *rpc.Client) *CelestiaFetcher {
	return &CelestiaFetcher{rpc: rpc}
}

// FetchPacaya implements the TxListFetcher interface.
func (d *CelestiaFetcher) FetchPacaya(ctx context.Context, meta metadata.TaikoBatchMetaDataPacaya) ([]byte, error) {
	if meta.GetCelestiaBlobsHeight() <= 0 {
		return nil, errors.New("celestia is not used")
	}

	namespace, err := share.NewNamespaceFromBytes(meta.GetCelestiaBlobsNamespace())
	if err != nil {
		return nil, err
	}

	celestiaBlobs, err := d.rpc.CelestiaDA.GetAll(ctx, meta.GetCelestiaBlobsHeight(), namespace)
	if err != nil {
		return nil, err
	}

	if len(celestiaBlobs) == 0 {
		return nil, errors.New("celestia blobs not found")
	}

	var txListBytes []byte
	for _, celestiaBlob := range celestiaBlobs {
		txListBytes = append(txListBytes, celestiaBlob.Data()...)
	}

	return sliceTxList(meta.GetBatchID(), txListBytes, meta.GetTxListOffset(), meta.GetTxListSize())
}
