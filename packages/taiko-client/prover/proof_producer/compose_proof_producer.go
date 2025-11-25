package producer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// RaikoBatches represents the JSON body of RaikoRequestProofBodyV3Pacaya's `Batches` field.
type RaikoBatches struct {
	BatchID                *big.Int `json:"batch_id"`
	L1InclusionBlockNumber *big.Int `json:"l1_inclusion_block_number"`
}

// RaikoRequestProofBodyV3Pacaya represents the JSON body for requesting the proof.
type RaikoRequestProofBodyV3Pacaya struct {
	Batches   []*RaikoBatches `json:"batches"`
	Prover    string          `json:"prover"`
	Aggregate bool            `json:"aggregate"`
	Type      ProofType       `json:"proof_type"`
}

type RaikoCheckpoint struct {
	BlockNum  *big.Int `json:"block_number"`
	BlockHash string   `json:"block_hash"`
	StateRoot string   `json:"state_root"`
}

// RaikoProposals represents the JSON body of RaikoRequestProofBodyV3Shasta's `Proposals` field.
type RaikoProposals struct {
	ProposalId             *big.Int         `json:"proposal_id"`
	L1InclusionBlockNumber *big.Int         `json:"l1_inclusion_block_number"`
	L2BlockNumbers         []*big.Int       `json:"l2_block_numbers"`
	DesignatedProver       string           `json:"designated_prover"`
	ParentTransitionHash   string           `json:"parent_transition_hash"`
	Checkpoint             *RaikoCheckpoint `json:"checkpoint"`
	LastAnchorBlockNumber  *big.Int         `json:"last_anchor_block_number"`
}

// RaikoRequestProofBodyV3Shasta represents the JSON body for requesting the proof.
type RaikoRequestProofBodyV3Shasta struct {
	Proposals []*RaikoProposals `json:"proposals"`
	Prover    string            `json:"prover"`
	Aggregate bool              `json:"aggregate"`
	Type      ProofType         `json:"proof_type"`
}

// ComposeProofProducer generates a compose proof for the given block.
type ComposeProofProducer struct {
	// We use Verifiers for Pacaya proof
	Verifiers map[ProofType]common.Address
	// We use VerifierIDs for Shasta proof
	VerifierIDs map[ProofType]uint8

	// Pacaya legacy fields (kept for compilation, not used)
	RaikoHostEndpoint string                // No-op: kept for Pacaya compatibility
	SgxGethProducer   *SgxGethProofProducer // No-op: kept for Pacaya compatibility
	ProofType         ProofType             // No-op: kept for Pacaya compatibility

	// Dual ZKVM configuration (active)
	RaikoZKVMEndpoint1 string
	RaikoZKVMEndpoint2 string
	ZKVMProofType1     ProofType
	ZKVMProofType2     ProofType

	RaikoRequestTimeout time.Duration
	ApiKey              string // ApiKey provided by Raiko
	Dummy               bool
	DummyProofProducer
}

// RequestProof implements the ProofProducer interface.
// Makes two parallel ZK proof requests to different ZKVM endpoints.
func (s *ComposeProofProducer) RequestProof(
	ctx context.Context,
	opts ProofRequestOptions,
	proposalID *big.Int,
	meta metadata.TaikoProposalMetaData,
	requestAt time.Time,
) (*ProofResponse, error) {
	log.Info(
		"Requesting dual ZKVM proofs",
		"proposalID", proposalID,
		"zkvm1Type", s.ZKVMProofType1,
		"zkvm2Type", s.ZKVMProofType2,
		"time", time.Since(requestAt),
	)

	var (
		g             = new(errgroup.Group)
		zkvm1Response *RaikoRequestProofBodyResponseV2
		zkvm2Response *RaikoRequestProofBodyResponseV2
		zkvm1Err      error
		zkvm2Err      error
	)

	// Request first ZKVM proof in parallel
	g.Go(func() error {
		resp, err := s.requestBatchProof(
			ctx,
			[]ProofRequestOptions{opts},
			[]metadata.TaikoProposalMetaData{meta},
			false,
			s.ZKVMProofType1,
			s.RaikoZKVMEndpoint1,
			requestAt,
			false,
		)
		if err != nil {
			zkvm1Err = fmt.Errorf("zkvm1 (%s): %w", s.ZKVMProofType1, err)
			return zkvm1Err
		}
		zkvm1Response = resp
		return nil
	})

	// Request second ZKVM proof in parallel
	g.Go(func() error {
		resp, err := s.requestBatchProof(
			ctx,
			[]ProofRequestOptions{opts},
			[]metadata.TaikoProposalMetaData{meta},
			false,
			s.ZKVMProofType2,
			s.RaikoZKVMEndpoint2,
			requestAt,
			false,
		)
		if err != nil {
			zkvm2Err = fmt.Errorf("zkvm2 (%s): %w", s.ZKVMProofType2, err)
			return zkvm2Err
		}
		zkvm2Response = resp
		return nil
	})

	// Wait for both proofs to complete
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get dual ZKVM proofs: %w", err)
	}

	log.Info(
		"Both ZKVM proofs generated successfully",
		"proposalID", proposalID,
		"zkvm1Type", zkvm1Response.ProofType,
		"zkvm2Type", zkvm2Response.ProofType,
		"time", time.Since(requestAt),
	)

	// Note: Since the single sp1 proof from raiko is null, we need to ignore the case.
	var proof1, proof2 []byte
	if ProofTypeZKSP1 != zkvm1Response.ProofType {
		proof1 = common.Hex2Bytes(zkvm1Response.Data.Proof.Proof[2:])
	}
	if ProofTypeZKSP1 != zkvm2Response.ProofType {
		proof2 = common.Hex2Bytes(zkvm2Response.Data.Proof.Proof[2:])
	}

	return &ProofResponse{
		BatchID:    proposalID,
		Meta:       meta,
		Proof1:     proof1,
		ProofType1: zkvm1Response.ProofType,
		Proof2:     proof2,
		ProofType2: zkvm2Response.ProofType,
		Opts:       opts,
	}, nil
}

// Aggregate implements the ProofProducer interface to aggregate a batch of proofs.
// Makes two parallel aggregation requests for both proof types.
func (s *ComposeProofProducer) Aggregate(
	ctx context.Context,
	items []*ProofResponse,
	requestAt time.Time,
) (*BatchProofs, error) {
	if len(items) == 0 {
		return nil, ErrInvalidLength
	}

	proofType1 := items[0].ProofType1
	proofType2 := items[0].ProofType2

	// Get verifier IDs for both proof types
	verifierID1, exist1 := s.VerifierIDs[proofType1]
	if !exist1 {
		return nil, fmt.Errorf("unknown proof type 1 from raiko: %s", proofType1)
	}
	verifierID2, exist2 := s.VerifierIDs[proofType2]
	if !exist2 {
		return nil, fmt.Errorf("unknown proof type 2 from raiko: %s", proofType2)
	}

	log.Info(
		"Aggregate batch proofs from raiko-host service",
		"proofType1", proofType1,
		"proofType2", proofType2,
		"batchSize", len(items),
		"firstID", items[0].BatchID,
		"lastID", items[len(items)-1].BatchID,
		"time", time.Since(requestAt),
	)

	var (
		batchIDs = make([]*big.Int, 0, len(items))
		opts     = make([]ProofRequestOptions, 0, len(items))
		metas    = make([]metadata.TaikoProposalMetaData, 0, len(items))
	)
	for _, item := range items {
		opts = append(opts, item.Opts)
		metas = append(metas, item.Meta)
		batchIDs = append(batchIDs, item.Meta.GetProposalID())
	}

	var (
		g              = new(errgroup.Group)
		batchProof1Res *RaikoRequestProofBodyResponseV2
		batchProof2Res *RaikoRequestProofBodyResponseV2
	)

	// Request first batch proof aggregation in parallel
	g.Go(func() error {
		resp, err := s.requestBatchProof(
			ctx,
			opts,
			metas,
			true,
			proofType1,
			s.RaikoZKVMEndpoint1,
			requestAt,
			false,
		)
		if err != nil {
			return fmt.Errorf("failed to aggregate proof type 1 (%s): %w", proofType1, err)
		}
		batchProof1Res = resp
		return nil
	})

	// Request second batch proof aggregation in parallel
	g.Go(func() error {
		resp, err := s.requestBatchProof(
			ctx,
			opts,
			metas,
			true,
			proofType2,
			s.RaikoZKVMEndpoint2,
			requestAt,
			false,
		)
		if err != nil {
			return fmt.Errorf("failed to aggregate proof type 2 (%s): %w", proofType2, err)
		}
		batchProof2Res = resp
		return nil
	})

	// Wait for both aggregations to complete
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to aggregate batch proofs: %w", err)
	}

	log.Info(
		"Both batch proofs aggregated successfully",
		"proofType1", batchProof1Res.ProofType,
		"proofType2", batchProof2Res.ProofType,
		"batchSize", len(items),
		"time", time.Since(requestAt),
	)

	batchProof1 := common.Hex2Bytes(batchProof1Res.Data.Proof.Proof[2:])
	batchProof2 := common.Hex2Bytes(batchProof2Res.Data.Proof.Proof[2:])

	return &BatchProofs{
		ProofResponses: items,
		BatchProof1:    batchProof1,
		ProofType1:     proofType1,
		VerifierID1:    verifierID1,
		BatchProof2:    batchProof2,
		ProofType2:     proofType2,
		VerifierID2:    verifierID2,
		BatchIDs:       batchIDs,
	}, nil
}

// requestBatchProof poll the proof aggregation service to get the aggregated proof.
func (s *ComposeProofProducer) requestBatchProof(
	ctx context.Context,
	opts []ProofRequestOptions,
	metas []metadata.TaikoProposalMetaData,
	isAggregation bool,
	proofType ProofType,
	endpoint string,
	requestAt time.Time,
	alreadyGenerated bool,
) (*RaikoRequestProofBodyResponseV2, error) {
	ctx, cancel := rpc.CtxWithTimeoutOrDefault(ctx, s.RaikoRequestTimeout)
	defer cancel()
	if len(opts) == 0 || len(opts) != len(metas) {
		return nil, ErrInvalidLength
	}
	var (
		output     *RaikoRequestProofBodyResponseV2
		err        error
		batches    = make([]*RaikoBatches, 0, len(opts))
		proposals  = make([]*RaikoProposals, 0, len(opts))
		start, end *big.Int
	)

	if metas[0].IsShasta() {
		for i, meta := range metas {
			proposals = append(proposals, &RaikoProposals{
				ProposalId:             meta.Shasta().GetProposal().Id,
				L1InclusionBlockNumber: meta.GetRawBlockHeight(),
				L2BlockNumbers:         opts[i].ShastaOptions().L2BlockNums,
				DesignatedProver:       opts[i].ShastaOptions().DesignatedProver.Hex()[2:],
				ParentTransitionHash:   opts[i].ShastaOptions().ParentTransitionHash.Hex()[2:],
				Checkpoint: &RaikoCheckpoint{
					BlockNum:  opts[i].ShastaOptions().Checkpoint.BlockNumber,
					BlockHash: common.BytesToHash(opts[i].ShastaOptions().Checkpoint.BlockHash[:]).Hex()[2:],
					StateRoot: common.BytesToHash(opts[i].ShastaOptions().Checkpoint.StateRoot[:]).Hex()[2:],
				},
				LastAnchorBlockNumber: opts[i].ShastaOptions().LastAnchorBlockNumber,
			})
		}
		output, err = requestHTTPProof[RaikoRequestProofBodyV3Shasta, RaikoRequestProofBodyResponseV2](
			ctx,
			endpoint+"/v3/proof/batch/shasta",
			s.ApiKey,
			RaikoRequestProofBodyV3Shasta{
				Type:      proofType,
				Proposals: proposals,
				Prover:    opts[0].GetProverAddress().Hex()[2:],
				Aggregate: isAggregation,
			},
		)
		start, end = proposals[0].ProposalId, proposals[len(proposals)-1].ProposalId
	} else {
		for _, meta := range metas {
			batches = append(batches, &RaikoBatches{
				BatchID:                meta.Pacaya().GetBatchID(),
				L1InclusionBlockNumber: meta.GetRawBlockHeight(),
			})
		}
		output, err = requestHTTPProof[RaikoRequestProofBodyV3Pacaya, RaikoRequestProofBodyResponseV2](
			ctx,
			endpoint+"/v3/proof/batch",
			s.ApiKey,
			RaikoRequestProofBodyV3Pacaya{
				Type:      proofType,
				Batches:   batches,
				Prover:    opts[0].GetProverAddress().Hex()[2:],
				Aggregate: isAggregation,
			},
		)
		start, end = batches[0].BatchID, batches[len(batches)-1].BatchID
	}
	if err != nil {
		return nil, err
	}

	if err := output.Validate(); err != nil {
		log.Debug(
			"Proof output validation result",
			"start", start,
			"end", end,
			"proofType", output.ProofType,
			"err", err,
		)
		return nil, fmt.Errorf("invalid Raiko response(start: %d, end: %d): %w",
			start,
			end,
			err,
		)
	}

	if !alreadyGenerated {
		proofType = output.ProofType
		log.Info(
			"Batch proof generated",
			"isAggregation", isAggregation,
			"proofType", proofType,
			"start", start,
			"end", end,
			"time", time.Since(requestAt),
		)
		// Update metrics.
		updateProvingMetrics(proofType, requestAt, isAggregation)
	}

	return output, nil
}
