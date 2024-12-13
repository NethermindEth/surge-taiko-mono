package producer

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
)

// CombinedProducer generates proofs from multiple producers in parallel and combines them.
type CombinedProducer struct {
	ProofTier      uint16
	RequiredProofs uint8
	Producers      []ProofProducer
	Verifiers      []common.Address
	// Map blockID to its proof state
	ProofStates map[*big.Int]BlockProofState
}

type BlockProofState struct {
	verifiedTiers []uint16
	proofs        []encoding.SubProof
}

// RequestProof implements the ProofProducer interface.
func (c *CombinedProducer) RequestProof(
	ctx context.Context,
	opts *ProofRequestOptions,
	blockID *big.Int,
	meta metadata.TaikoBlockMetaData,
	header *types.Header,
	requestAt time.Time,
) (*ProofWithHeader, error) {
	log.Debug("CombinedProducer: RequestProof", "blockID", blockID)
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		errorsChan = make(chan error, len(c.Producers))
	)

	// Get or initialize proof state
	proofState, ok := c.ProofStates[blockID]
	if !ok {
		proofState = BlockProofState{
			verifiedTiers: []uint16{},
			proofs:        []encoding.SubProof{},
		}
		c.ProofStates[blockID] = proofState
	}

	taskCtx, taskCtxCancel := context.WithCancel(ctx)
	defer taskCtxCancel()

	for i, producer := range c.Producers {
		if slices.Contains(proofState.verifiedTiers, producer.Tier()) {
			log.Debug("Skipping producer, proof already verified", "tier", producer.Tier())
			continue
		}

		verifier := c.Verifiers[i]

		wg.Add(1)
		go func(idx int, p ProofProducer, verifier common.Address) {
			defer wg.Done()

			proofWithHeader, err := p.RequestProof(taskCtx, opts, blockID, meta, header, requestAt)
			if err != nil {
				errorsChan <- fmt.Errorf("producer %d error: %w", idx, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			proofState.verifiedTiers = append(proofState.verifiedTiers, p.Tier())
			if uint8(len(proofState.proofs)) < c.RequiredProofs {
				proofState.proofs = append(
					proofState.proofs,
					encoding.SubProof{
						Proof:    proofWithHeader.Proof,
						Verifier: verifier,
					},
				)
			}

			if uint8(len(proofState.proofs)) == c.RequiredProofs {
				taskCtxCancel()
			}
		}(i, producer, verifier)
	}

	wg.Wait()

	if uint8(len(proofState.proofs)) < c.RequiredProofs {
		var errMsgs []string

		errMsgs = append(
			errMsgs,
			fmt.Sprintf("not enough proofs collected: required %d, got %d", c.RequiredProofs, len(proofState.proofs)),
		)

		close(errorsChan)
		for err := range errorsChan {
			errMsgs = append(errMsgs, err.Error())
		}

		return nil, fmt.Errorf("combined proof production failed: %s", strings.Join(errMsgs, "; "))
	}

	combinedProof, err := encoding.EncodeSubProofs(proofState.proofs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode sub proofs: %w", err)
	}

	log.Info(
		"Combined proofs generated",
		"blockID", blockID,
		"time", time.Since(requestAt),
		"producer", "CombinedProducer",
	)

	c.CleanOldProofStates(blockID, &mu)

	return &ProofWithHeader{
		BlockID: blockID,
		Header:  header,
		Meta:    meta,
		Proof:   combinedProof,
		Opts:    opts,
		Tier:    c.Tier(),
	}, nil
}

// CleanOldProofStates removes proof states for blocks older than 256 blocks.
func (c *CombinedProducer) CleanOldProofStates(latestBlockID *big.Int, mutex *sync.Mutex) {
	if len(c.ProofStates) == 0 {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	delete(c.ProofStates, latestBlockID)

	threshold := new(big.Int).Sub(latestBlockID, big.NewInt(256))
	for blockID := range c.ProofStates {
		if blockID.Cmp(threshold) < 0 {
			delete(c.ProofStates, blockID)
		}
	}
}

// RequestCancel implements the ProofProducer interface.
func (c *CombinedProducer) RequestCancel(
	ctx context.Context,
	opts *ProofRequestOptions,
) error {
	var finalError error
	for _, producer := range c.Producers {
		if err := producer.RequestCancel(ctx, opts); err != nil {
			if finalError == nil {
				finalError = err
			}
		}
	}
	return finalError
}

// Tier implements the ProofProducer interface.
func (c *CombinedProducer) Tier() uint16 {
	return c.ProofTier
}
