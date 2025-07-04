package txlistfetcher

import (
	"context"
	"errors"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// CelestiaFetcher is responsible for fetching  txList blob from Celestia.
type CelestiaFetcher struct {
	rpc *rpc.Client
}

// NewCelestiaFetcher creates a new CelestiaFetcher instance based on the given rpc client.
func NewCelestiaFetcher(rpc *rpc.Client) *CalldataFetcher {
	return &CalldataFetcher{rpc: rpc}
}

// FetchPacaya implements the TxListFetcher interface.
func (d *CelestiaFetcher) FetchPacaya(ctx context.Context, meta metadata.TaikoBatchMetaDataPacaya) ([]byte, error) {
	return nil, errors.New("[Not Implemented] CelestiaFetcher.FetchPacaya")
}
