package handler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/metrics"
)

// BlockVerifiedEventHandler is responsible for handling the BlockVerified event.
type BlockVerifiedEventHandler struct {
}

// NewBlockVerifiedEventHandler creates a new BlockVerifiedEventHandler instance.
func NewBlockVerifiedEventHandler() *BlockVerifiedEventHandler {
	return &BlockVerifiedEventHandler{}
}

// Handle handles the BlockVerified event.
func (h *BlockVerifiedEventHandler) Handle(e *bindings.TaikoL1ClientBlockVerifiedV2) {
	metrics.DriverL2VerifiedHeightGauge.Set(float64(e.BlockId.Uint64()))
	metrics.ProverLatestVerifiedIDGauge.Set(float64(e.BlockId.Uint64()))

	log.Info(
		"New verified block",
		"blockID", e.BlockId,
		"hash", common.BytesToHash(e.BlockHash[:]),
		"prover", e.Prover,
	)
}
