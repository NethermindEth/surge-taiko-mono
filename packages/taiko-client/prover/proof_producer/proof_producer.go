package producer

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
)

var (
	ErrEmptyProof      = errors.New("proof is empty")
	ErrInvalidLength   = errors.New("invalid items length")
	ErrProofInProgress = errors.New("work_in_progress")
	ErrRetry           = errors.New("retry")
	ErrZkAnyNotDrawn   = errors.New("zk_any_not_drawn")
	StatusRegistered   = "registered"
)

// ProofRequestBody represents a request body to generate a proof.
type ProofRequestBody struct {
	Meta metadata.TaikoProposalMetaData
}

// ProofResponse represents a response of a dual proof request.
type ProofResponse struct {
	BatchID    *big.Int
	Meta       metadata.TaikoProposalMetaData
	Proof1     []byte
	ProofType1 ProofType
	Proof2     []byte
	ProofType2 ProofType
	Opts       ProofRequestOptions
	// Legacy fields (for Pacaya compatibility)
	Proof     []byte
	ProofType ProofType
}

// BatchProofs represents a response of a dual batch proof request.
type BatchProofs struct {
	ProofResponses []*ProofResponse
	BatchIDs       []*big.Int

	BatchProof1 []byte
	ProofType1  ProofType
	Verifier1   common.Address
	VerifierID1 uint8

	BatchProof2 []byte
	ProofType2  ProofType
	Verifier2   common.Address
	VerifierID2 uint8
}

// ProofProducer is an interface that contains all methods to generate a proof.
type ProofProducer interface {
	RequestProof(
		ctx context.Context,
		opts ProofRequestOptions,
		batchID *big.Int,
		meta metadata.TaikoProposalMetaData,
		requestAt time.Time,
	) (*ProofResponse, error)
	Aggregate(ctx context.Context, items []*ProofResponse, requestAt time.Time) (*BatchProofs, error)
}
