package producer

import (
	"context"
	"errors"
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

// ComposeProofProducer generates a compose proof for the given block.
type ComposeProofProducer struct {
	Verifiers map[ProofType]common.Address

	RaikoSGXHostEndpoint      string
	RaikoTDXHostEndpoint      string
	RaikoAzureTDXHostEndpoint string
	RaikoZKVMHostEndpoint     string

	RaikoRequestTimeout time.Duration
	JWT                 string // JWT provided by Raiko
	ProofType           ProofType
	Dummy               bool
	DummyProofProducer
}

// RequestProof implements the ProofProducer interface.
func (s *ComposeProofProducer) RequestProof(
	ctx context.Context,
	opts ProofRequestOptions,
	batchID *big.Int,
	meta metadata.TaikoProposalMetaData,
	requestAt time.Time,
) (*ProofResponse, error) {
	if !meta.IsPacaya() {
		return nil, fmt.Errorf("current proposal (%d) is not a Pacaya proposal", batchID)
	}

	log.Info(
		"Request SGX + TDX + ZK proofs from raiko-host service",
		"batchID", batchID,
		"coinbase", meta.Pacaya().GetCoinbase(),
		"time", time.Since(requestAt),
	)

	var (
		proof        []byte
		zkProofType  ProofType
		sgxProofType ProofType
		tdxProofType ProofType
		batches      = []*RaikoBatches{{BatchID: batchID, L1InclusionBlockNumber: meta.GetRawBlockHeight()}}
		g            = new(errgroup.Group)
	)

	g.Go(func() error {
		// SGX (any) proof request raiko-host service - can be SGX Reth or SGX Geth
		if s.Dummy {
			log.Debug("Dummy proof producer requested SGX proof", "batchID", batchID)

			// The following line is a no-op; this is just to showcase the dummy proof producer
			_, _ = s.DummyProofProducer.RequestProof(opts, batchID, meta, requestAt)
			sgxProofType = ProofTypeSgx // Default to SGX for dummy
			return nil
		}

		resp, err := s.requestBatchProof(
			ctx,
			batches,
			opts.GetProverAddress(),
			false,
			ProofTypeSgxAny,
			requestAt,
			opts.PacayaOptions().IsRethSGXProofGenerated,
		)
		if err != nil {
			return err
		}

		sgxProofType = resp.ProofType

		// Note: we mark the `IsRethSGXProofGenerated` with true to record if it is first time generated
		opts.PacayaOptions().IsRethSGXProofGenerated = true
		return nil
	})
	g.Go(func() error {
		// TDX (any) proof request raiko-host service - can be either TDX or Azure TDX
		if s.Dummy {
			log.Debug("Dummy proof producer requested TDX proof", "batchID", batchID)

			// The following line is a no-op; this is just to showcase the dummy proof producer
			_, _ = s.DummyProofProducer.RequestProof(opts, batchID, meta, requestAt)
			tdxProofType = ProofTypeTdx // Default to TDX for dummy
			return nil
		}

		resp, err := s.requestBatchProof(
			ctx,
			batches,
			opts.GetProverAddress(),
			false,
			ProofTypeTdxAny,
			requestAt,
			opts.PacayaOptions().IsNethermindTdxProofGenerated,
		)
		if err != nil {
			return err
		}

		tdxProofType = resp.ProofType

		// Note: we mark the `IsNethermindTdxProofGenerated` with true to record if it is first time generated
		if resp.ProofType == ProofTypeAzureTdx {
			opts.PacayaOptions().IsNethermindAzureTdxProofGenerated = true
		} else {
			opts.PacayaOptions().IsNethermindTdxProofGenerated = true
		}
		return nil
	})
	g.Go(func() error {
		// ZK proof request raiko-host service
		if s.Dummy {
			log.Debug("Dummy proof producer requested ZK proof", "batchID", batchID)

			// For the dummy proof producer, we just use the sp1 proof type (as zk_any would break the logic down the line)
			zkProofType = ProofTypeZKSP1
			if resp, err := s.DummyProofProducer.RequestProof(opts, batchID, meta, requestAt); err != nil {
				return err
			} else {
				proof = resp.Proof
			}

			return nil
		}

		resp, err := s.requestBatchProof(
			ctx,
			batches,
			opts.GetProverAddress(),
			false,
			ProofTypeZKAny,
			requestAt,
			opts.PacayaOptions().IsRethZKProofGenerated,
		)
		if err != nil {
			return err
		}

		zkProofType = resp.ProofType
		// Note: we mark the `IsRethZKProofGenerated` with true to record if it is first time generated
		opts.PacayaOptions().IsRethZKProofGenerated = true
		// Note: Since the single sp1 proof from raiko is null, we need to ignore the case.
		if ProofTypeZKSP1 != zkProofType {
			proof = common.Hex2Bytes(resp.Data.Proof.Proof[2:])
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get batches proofs: %w", err)
	}

	return &ProofResponse{
		BatchID:      batchID,
		Meta:         meta,
		Proof:        proof,
		Opts:         opts,
		ProofType:    zkProofType, // Backward compatibility
		ZKProofType:  zkProofType,
		SGXProofType: sgxProofType,
		TDXProofType: tdxProofType,
	}, nil
}

// Aggregate implements the ProofProducer interface to aggregate a batch of proofs.
func (s *ComposeProofProducer) Aggregate(
	ctx context.Context,
	items []*ProofResponse,
	requestAt time.Time,
) (*BatchProofs, error) {
	if len(items) == 0 {
		return nil, ErrInvalidLength
	}
	// TODO(@jmadibekov): manually test the scenario when risc0 and sp1 proofs are mixed in the same group of batches.
	firstItem := items[0]
	zkProofType := firstItem.ZKProofType
	if zkProofType == "" {
		zkProofType = firstItem.ProofType // Backward compatibility
	}
	verifier, exist := s.Verifiers[zkProofType]
	if !exist {
		return nil, fmt.Errorf("unknown proof type from raiko %s", zkProofType)
	}
	log.Info(
		"Aggregate batch proofs from raiko-host service",
		"zkProofType", zkProofType,
		"batchSize", len(items),
		"firstID", firstItem.BatchID,
		"lastID", items[len(items)-1].BatchID,
		"time", time.Since(requestAt),
	)
	var (
		sgxBatchProofs      []byte
		tdxBatchProofs      []byte
		azureTdxBatchProofs []byte
		batchProofs         []byte
		batches             = make([]*RaikoBatches, 0, len(items))
		batchIDs            = make([]*big.Int, 0, len(items))
		g                   = new(errgroup.Group)
	)
	for _, item := range items {
		batches = append(batches, &RaikoBatches{
			BatchID:                item.Meta.Pacaya().GetBatchID(),
			L1InclusionBlockNumber: item.Meta.GetRawBlockHeight(),
		})
		batchIDs = append(batchIDs, item.Meta.Pacaya().GetBatchID())
	}
	g.Go(func() error {
		if s.Dummy {
			log.Debug("Dummy proof producer requested SGX batch proof aggregation", "batchSize", len(items))

			resp, _ := s.DummyProofProducer.RequestBatchProofs(items, ProofTypeSgx)
			sgxBatchProofs = resp.BatchProof
			return nil
		}

		resp, err := s.requestBatchProof(
			ctx,
			batches,
			firstItem.Opts.GetProverAddress(),
			true,
			ProofTypeSgx,
			requestAt,
			firstItem.Opts.PacayaOptions().IsRethSGXProofAggregationGenerated,
		)
		if err != nil {
			return err
		}

		// Note: we mark the `IsRethSGXProofAggregationGenerated` in the first item with true
		// to record if it is first time generated
		firstItem.Opts.PacayaOptions().IsRethSGXProofAggregationGenerated = true
		sgxBatchProofs = common.Hex2Bytes(resp.Data.Proof.Proof[2:])

		return nil
	})
	g.Go(func() error {
		if s.Dummy {
			log.Debug("Dummy proof producer requested TDX batch proof aggregation", "batchSize", len(items))

			resp, _ := s.DummyProofProducer.RequestBatchProofs(items, ProofTypeTdx)
			tdxBatchProofs = resp.BatchProof
			return nil
		}

		resp, err := s.requestBatchProof(
			ctx,
			batches,
			firstItem.Opts.GetProverAddress(),
			true,
			ProofTypeTdxAny,
			requestAt,
			firstItem.Opts.PacayaOptions().IsNethermindTdxProofAggregationGenerated,
		)
		if err != nil {
			return err
		}

		// Note: we mark the `IsNethermindTdxProofAggregationGenerated` in the first item with true
		// to record if it is first time generated
		firstItem.Opts.PacayaOptions().IsNethermindTdxProofAggregationGenerated = true
		// Determine which TDX batch proofs to use based on the response proof type
		if resp.ProofType == ProofTypeAzureTdx {
			azureTdxBatchProofs = common.Hex2Bytes(resp.Data.Proof.Proof[2:])
		} else {
			tdxBatchProofs = common.Hex2Bytes(resp.Data.Proof.Proof[2:])
		}

		return nil
	})
	g.Go(func() error {
		if s.Dummy {
			log.Debug("Dummy proof producer requested ZK batch proof aggregation", "batchSize", len(items))

			resp, _ := s.DummyProofProducer.RequestBatchProofs(items, zkProofType)
			batchProofs = resp.BatchProof
			return nil
		}

		resp, err := s.requestBatchProof(
			ctx,
			batches,
			firstItem.Opts.GetProverAddress(),
			true,
			zkProofType,
			requestAt,
			firstItem.Opts.PacayaOptions().IsRethZKProofAggregationGenerated,
		)
		if err != nil {
			return err
		}

		// Note: we mark the `IsRethZKProofAggregationGenerated` in the first item with true
		// to record if it is first time generated
		firstItem.Opts.PacayaOptions().IsRethZKProofAggregationGenerated = true
		batchProofs = common.Hex2Bytes(resp.Data.Proof.Proof[2:])

		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get batches proofs: %w", err)
	}

	return &BatchProofs{
		ProofResponses:        items,
		BatchProof:            batchProofs,
		BatchIDs:              batchIDs,
		ProofType:             zkProofType,
		Verifier:              verifier,
		SgxBatchProof:         sgxBatchProofs,
		SgxProofVerifier:      s.Verifiers[ProofTypeSgx],
		TdxBatchProof:         tdxBatchProofs,
		TdxProofVerifier:      s.Verifiers[ProofTypeTdx],
		AzureTdxBatchProof:    azureTdxBatchProofs,
		AzureTdxProofVerifier: s.Verifiers[ProofTypeAzureTdx],
	}, nil
}

// requestBatchProof poll the proof aggregation service to get the aggregated proof.
func (s *ComposeProofProducer) requestBatchProof(
	ctx context.Context,
	batches []*RaikoBatches,
	proverAddress common.Address,
	isAggregation bool,
	proofType ProofType,
	requestAt time.Time,
	alreadyGenerated bool,
) (*RaikoRequestProofBodyResponseV2, error) {
	ctx, cancel := rpc.CtxWithTimeoutOrDefault(ctx, s.RaikoRequestTimeout)
	defer cancel()

	var endpoints []string
	switch proofType {
	case ProofTypeSgx, ProofTypeSgxAny:
		endpoints = append(endpoints, s.RaikoSGXHostEndpoint)
	case ProofTypeTdx:
		endpoints = append(endpoints, s.RaikoTDXHostEndpoint)
	case ProofTypeAzureTdx:
		endpoints = append(endpoints, s.RaikoAzureTDXHostEndpoint)
	case ProofTypeTdxAny:
		endpoints = append(endpoints, s.RaikoTDXHostEndpoint, s.RaikoAzureTDXHostEndpoint)
	case ProofTypeZKAny, ProofTypeZKR0, ProofTypeZKSP1:
		endpoints = append(endpoints, s.RaikoZKVMHostEndpoint)
	default:
		return nil, fmt.Errorf("unexpected proof type: %s", proofType)
	}

	var errs []error

	for _, endpoint := range endpoints {
		log.Debug(
			"Making HTTP request to raiko",
			"endpoint", endpoint+"/v3/proof/batch",
			"request", RaikoRequestProofBodyV3Pacaya{
				Type:      proofType,
				Batches:   batches,
				Prover:    proverAddress.Hex()[2:],
				Aggregate: isAggregation,
			},
		)

		output, err := requestHTTPProof[RaikoRequestProofBodyV3Pacaya, RaikoRequestProofBodyResponseV2](
			ctx,
			endpoint+"/v3/proof/batch",
			s.JWT,
			RaikoRequestProofBodyV3Pacaya{
				Type:      proofType,
				Batches:   batches,
				Prover:    proverAddress.Hex()[2:],
				Aggregate: isAggregation,
			},
		)
		if err != nil {
			log.Debug(
				"Error making HTTP request to raiko",
				"endpoint", endpoint+"/v3/proof/batch",
				"request", RaikoRequestProofBodyV3Pacaya{
					Type:      proofType,
					Batches:   batches,
					Prover:    proverAddress.Hex()[2:],
					Aggregate: isAggregation,
				},
				"error", err,
			)
			errs = append(errs, err)
			continue
		}

		if err := output.Validate(); err != nil {
			log.Debug(
				"Proof output validation result",
				"start", batches[0].BatchID,
				"end", batches[len(batches)-1].BatchID,
				"proofType", output.ProofType,
				"err", err,
			)
			errs = append(errs, fmt.Errorf("invalid Raiko response(start: %d, end: %d): %w",
				batches[0].BatchID,
				batches[len(batches)-1].BatchID,
				err,
			))
			continue
		}

		if !alreadyGenerated {
			proofType = output.ProofType
			log.Info(
				"Batch proof generated",
				"isAggregation", isAggregation,
				"proofType", proofType,
				"start", batches[0].BatchID,
				"end", batches[len(batches)-1].BatchID,
				"time", time.Since(requestAt),
			)
			// Update metrics.
			updateProvingMetrics(proofType, requestAt, isAggregation)
		}

		return output, nil
	}

	return nil, fmt.Errorf("failed to get batch proof from proof type %s endpoints: %w", proofType, errors.Join(errs...))
}
