package metadata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	surgeBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/surge"
)

// Ensure TaikoProposalMetadataSurge implements TaikoBlockMetaData.
var _ TaikoProposalMetaData = (*TaikoProposalMetadataSurge)(nil)

// TaikoProposalMetadataSurge is the metadata of a Surge Taiko blocks batch.
type TaikoProposalMetadataSurge struct {
	*surgeBindings.SurgeInboxClientProposed
	timestamp uint64
}

// NewTaikoProposalMetadataSurge creates a new instance of TaikoProposalMetadataSurge
// from the SurgeInboxClient.Proposed event.
func NewTaikoProposalMetadataSurge(
	e *surgeBindings.SurgeInboxClientProposed,
	timestamp uint64,
) *TaikoProposalMetadataSurge {
	return &TaikoProposalMetadataSurge{
		SurgeInboxClientProposed: e,
		timestamp:                timestamp,
	}
}

// Pacaya implements TaikoProposalMetaData interface.
func (m *TaikoProposalMetadataSurge) Pacaya() TaikoBatchMetaDataPacaya {
	return nil
}

// IsPacaya implements TaikoProposalMetaData interface.
func (m *TaikoProposalMetadataSurge) IsPacaya() bool {
	return false
}

// Shasta implements TaikoProposalMetaData interface.
func (m *TaikoProposalMetadataSurge) Shasta() TaikoProposalMetaDataShasta {
	return nil
}

// IsShasta implements TaikoProposalMetaData interface.
func (m *TaikoProposalMetadataSurge) IsShasta() bool {
	return false
}

// Surge implements TaikoProposalMetaData interface.
func (m *TaikoProposalMetadataSurge) Surge() TaikoProposalMetaDataSurge {
	return m
}

// IsSurge implements TaikoProposalMetaData interface.
func (m *TaikoProposalMetadataSurge) IsSurge() bool {
	return true
}

// GetRawBlockHeight returns the raw L1 block height.
func (m *TaikoProposalMetadataSurge) GetRawBlockHeight() *big.Int {
	return new(big.Int).SetUint64(m.Raw.BlockNumber)
}

// GetRawBlockHash returns the raw L1 block hash.
func (m *TaikoProposalMetadataSurge) GetRawBlockHash() common.Hash {
	return m.Raw.BlockHash
}

// GetTxIndex returns the transaction index.
func (m *TaikoProposalMetadataSurge) GetTxIndex() uint {
	return m.Raw.TxIndex
}

// GetTxHash returns the transaction hash.
func (m *TaikoProposalMetadataSurge) GetTxHash() common.Hash {
	return m.Raw.TxHash
}

// GetProposer returns the proposer of this batch.
func (m *TaikoProposalMetadataSurge) GetProposer() common.Address {
	return m.Proposer
}

// GetCoinbase returns block coinbase. Sets it to common.Address{}, because we need to fetch the value from blob.
func (m *TaikoProposalMetadataSurge) GetCoinbase() common.Address {
	return common.Address{}
}

func (m *TaikoProposalMetadataSurge) GetLog() *types.Log {
	return &m.Raw
}

// GetBlobHashes returns blob hashes in this proposal.
func (m *TaikoProposalMetadataSurge) GetBlobHashes(idx int) []common.Hash {
	var blobHashes []common.Hash
	if len(m.Sources) <= idx {
		return blobHashes
	}
	for _, hash := range m.Sources[idx].BlobSlice.BlobHashes {
		blobHashes = append(blobHashes, hash)
	}
	return blobHashes
}

// GetBlobTimestamp returns the timestamp of the blob slice in this proposal.
func (m *TaikoProposalMetadataSurge) GetBlobTimestamp(idx int) uint64 {
	if len(m.Sources) <= idx {
		return 0
	}
	return m.Sources[idx].BlobSlice.Timestamp.Uint64()
}

// GetProposalID returns proposal ID.
func (m *TaikoProposalMetadataSurge) GetProposalID() *big.Int {
	return m.Id
}

// GetEventData returns the underlying event data.
func (m *TaikoProposalMetadataSurge) GetEventData() *surgeBindings.SurgeInboxClientProposed {
	return m.SurgeInboxClientProposed
}

// GetTimestamp returns the timestamp of the proposal.
func (m *TaikoProposalMetadataSurge) GetTimestamp() uint64 {
	return m.timestamp
}
