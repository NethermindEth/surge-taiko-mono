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

	// Dual ZKVM configuration
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
		"Request dual ZKVM proofs from raiko-host service",
		"proposalID", proposalID,
		"zkvm1Endpoint", s.RaikoZKVMEndpoint1,
		"zkvm1Type", s.ZKVMProofType1,
		"zkvm2Endpoint", s.RaikoZKVMEndpoint2,
		"zkvm2Type", s.ZKVMProofType2,
		"time", time.Since(requestAt),
		"dummy", s.Dummy,
	)

	var (
		zkvm1Proof []byte
		zkvm2Proof []byte
		g          = new(errgroup.Group)
		zkvm1Err   error
		zkvm2Err   error
	)

	// Request first ZKVM proof in parallel
	g.Go(func() error {
		if s.Dummy {
			if resp, err := s.DummyProofProducer.RequestProof(ctx, opts, proposalID, meta, requestAt); err != nil {
				zkvm1Err = fmt.Errorf("zkvm1 (%s): %w", s.ZKVMProofType1, err)
				return zkvm1Err
			} else {
				zkvm1Proof = resp.Proof
			}
		} else {
			resp, err := s.requestBatchProof(
				ctx,
				[]ProofRequestOptions{opts},
				[]metadata.TaikoProposalMetaData{meta},
				false,
				s.ZKVMProofType1,
				s.RaikoZKVMEndpoint1,
				requestAt,
				opts.IsZKVMProof1Generated(),
			)
			if err != nil {
				zkvm1Err = fmt.Errorf("zkvm1 (%s): %w", s.ZKVMProofType1, err)
				return zkvm1Err
			} else {
				if resp.ProofType != s.ZKVMProofType1 {
					log.Error(
						"ZKVM1 proof type mismatch",
						"expected", s.ZKVMProofType1,
						"got", resp.ProofType,
					)
				}

				// Note: we mark the `IsZKVMProof1Generated` with true to record if it is first time generated
				if opts.IsShasta() {
					opts.ShastaOptions().ZKVMProof1Generated = true
				} else {
					opts.PacayaOptions().ZKVMProof1Generated = true
				}
				// Note: Since the single sp1 proof from raiko is null, we need to ignore the case.
				if ProofTypeZKSP1 != resp.ProofType {
					zkvm1Proof = common.Hex2Bytes(resp.Data.Proof.Proof[2:])
				}
			}
		}
		return nil
	})

	// Request second ZKVM proof in parallel
	g.Go(func() error {
		if s.Dummy {
			if resp, err := s.DummyProofProducer.RequestProof(ctx, opts, proposalID, meta, requestAt); err != nil {
				zkvm2Err = fmt.Errorf("zkvm2 (%s): %w", s.ZKVMProofType2, err)
				return zkvm2Err
			} else {
				zkvm2Proof = resp.Proof
			}
		} else {
			resp, err := s.requestBatchProof(
				ctx,
				[]ProofRequestOptions{opts},
				[]metadata.TaikoProposalMetaData{meta},
				false,
				s.ZKVMProofType2,
				s.RaikoZKVMEndpoint2,
				requestAt,
				opts.IsZKVMProof2Generated(),
			)
			if err != nil {
				zkvm2Err = fmt.Errorf("zkvm2 (%s): %w", s.ZKVMProofType2, err)
				return zkvm2Err
			} else {
				if resp.ProofType != s.ZKVMProofType2 {
					log.Error(
						"ZKVM2 proof type mismatch",
						"expected", s.ZKVMProofType2,
						"got", resp.ProofType,
					)
				}

				// Note: we mark the `IsZKVMProof2Generated` with true to record if it is first time generated
				if opts.IsShasta() {
					opts.ShastaOptions().ZKVMProof2Generated = true
				} else {
					opts.PacayaOptions().ZKVMProof2Generated = true
				}
				// Note: Since the single sp1 proof from raiko is null, we need to ignore the case.
				if ProofTypeZKSP1 != resp.ProofType {
					zkvm2Proof = common.Hex2Bytes(resp.Data.Proof.Proof[2:])
				}
			}
		}
		return nil
	})

	// Wait for both proofs to complete
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get dual ZKVM proofs: %w and %w and %w", err, zkvm1Err, zkvm2Err)
	}

	return &ProofResponse{
		BatchID:    proposalID,
		Meta:       meta,
		Proof1:     zkvm1Proof,
		ProofType1: s.ZKVMProofType1,
		Proof2:     zkvm2Proof,
		ProofType2: s.ZKVMProofType2,
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
	isShasta := items[0].Meta.IsShasta()
	var (
		verifierID1 uint8
		verifierID2 uint8
		verifier1   common.Address
		verifier2   common.Address
		exist       bool
	)
	if isShasta {
		if verifierID1, exist = s.VerifierIDs[proofType1]; !exist {
			return nil, fmt.Errorf("unknown proof type from raiko %s", proofType1)
		}
		if verifierID2, exist = s.VerifierIDs[proofType2]; !exist {
			return nil, fmt.Errorf("unknown proof type from raiko %s", proofType2)
		}
	} else {
		if verifier1, exist = s.Verifiers[proofType1]; !exist {
			return nil, fmt.Errorf("unknown proof type from raiko %s", proofType1)
		}
		if verifier2, exist = s.Verifiers[proofType2]; !exist {
			return nil, fmt.Errorf("unknown proof type from raiko %s", proofType2)
		}
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
		batchProofs1 []byte
		batchProofs2 []byte
		batchIDs     = make([]*big.Int, 0, len(items))
		opts         = make([]ProofRequestOptions, 0, len(items))
		metas        = make([]metadata.TaikoProposalMetaData, 0, len(items))
		g            = new(errgroup.Group)
	)
	for _, item := range items {
		opts = append(opts, item.Opts)
		metas = append(metas, item.Meta)
		batchIDs = append(batchIDs, item.Meta.GetProposalID())
	}

	// Request first batch proof aggregation in parallel
	g.Go(func() error {
		if s.Dummy {
			resp, _ := s.DummyProofProducer.RequestBatchProofs(items, proofType1)
			batchProofs1 = resp.BatchProof1
		} else {
			if resp, err := s.requestBatchProof(
				ctx,
				opts,
				metas,
				true,
				proofType1,
				s.RaikoZKVMEndpoint1,
				requestAt,
				items[0].Opts.IsZKVMProof1AggregationGenerated(),
			); err != nil {
				return fmt.Errorf("failed to aggregate proof type 1 (%s): %w", proofType1, err)
			} else {
				// Note: we mark the `IsZKVMProof1AggregationGenerated` in the first item with true
				// to record if it is first time generated
				if items[0].Opts.IsShasta() {
					items[0].Opts.ShastaOptions().ZKVMProof1AggregationGenerated = true
				} else {
					items[0].Opts.PacayaOptions().ZKVMProof1AggregationGenerated = true
				}
				batchProofs1 = common.Hex2Bytes(resp.Data.Proof.Proof[2:])
			}
		}
		return nil
	})

	// Request second batch proof aggregation in parallel
	g.Go(func() error {
		if s.Dummy {
			resp, _ := s.DummyProofProducer.RequestBatchProofs(items, proofType2)
			batchProofs2 = resp.BatchProof1
		} else {
			if resp, err := s.requestBatchProof(
				ctx,
				opts,
				metas,
				true,
				proofType2,
				s.RaikoZKVMEndpoint2,
				requestAt,
				items[0].Opts.IsZKVMProof2AggregationGenerated(),
			); err != nil {
				return fmt.Errorf("failed to aggregate proof type 2 (%s): %w", proofType2, err)
			} else {
				// Note: we mark the `IsZKVMProof2AggregationGenerated` in the first item with true
				// to record if it is first time generated
				if items[0].Opts.IsShasta() {
					items[0].Opts.ShastaOptions().ZKVMProof2AggregationGenerated = true
				} else {
					items[0].Opts.PacayaOptions().ZKVMProof2AggregationGenerated = true
				}
				batchProofs2 = common.Hex2Bytes(resp.Data.Proof.Proof[2:])
			}
		}
		return nil
	})

	// Wait for both aggregations to complete
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to aggregate batch proofs: %w", err)
	}

	return &BatchProofs{
		ProofResponses: items,
		BatchIDs:       batchIDs,

		BatchProof1: batchProofs1,
		ProofType1:  proofType1,
		Verifier1:   verifier1,
		VerifierID1: verifierID1,

		BatchProof2: batchProofs2,
		ProofType2:  proofType2,
		Verifier2:   verifier2,
		VerifierID2: verifierID2,
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
				ProposalId:             meta.Shasta().GetEventData().Id,
				L1InclusionBlockNumber: meta.GetRawBlockHeight(),
				L2BlockNumbers:         opts[i].ShastaOptions().L2BlockNums,
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
