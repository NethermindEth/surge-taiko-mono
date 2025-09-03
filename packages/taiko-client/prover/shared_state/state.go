package state

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"

	taikoTypes "github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/types"
)

// SharedState represents the internal state of a prover.
type SharedState struct {
	lastHandledBatchID      atomic.Uint64
	l1Current               atomic.Value
	batchesRollbackedRanges atomic.Value
}

// New creates a new prover shared state instance.
func New() *SharedState {
	return new(SharedState)
}

// GetLastHandledBatchID returns the last handled batch ID.
func (s *SharedState) GetLastHandledBatchID() uint64 {
	return s.lastHandledBatchID.Load()
}

// SetLastHandledBatchID sets the last handled batch ID.
func (s *SharedState) SetLastHandledBatchID(batchID uint64) {
	s.lastHandledBatchID.Store(batchID)
}

// GetL1Current returns the current L1 header cursor.
func (s *SharedState) GetL1Current() *types.Header {
	if val := s.l1Current.Load(); val != nil {
		return val.(*types.Header)
	}
	return nil
}

// SetL1Current sets the current L1 header cursor.
func (s *SharedState) SetL1Current(header *types.Header) {
	s.l1Current.Store(header)
}

// GetBatchesRollbackedRanges returns the batches rollbacked ranges.
func (s *SharedState) GetBatchesRollbackedRanges() taikoTypes.BatchesRollbackedRanges {
	if val := s.batchesRollbackedRanges.Load(); val != nil {
		return val.(taikoTypes.BatchesRollbackedRanges)
	}
	return nil
}

// SetBatchesRollbackedRanges sets the batches rollbacked ranges.
func (s *SharedState) SetBatchesRollbackedRanges(ranges taikoTypes.BatchesRollbackedRanges) {
	s.batchesRollbackedRanges.Store(ranges)
}
