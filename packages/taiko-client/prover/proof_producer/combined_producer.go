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

type BlockProofState struct {
	verifiedTiers []uint16
	proofs        []encoding.SubProof
}

type ProofStateManager struct {
	mu     sync.Mutex
	states map[uint64]*BlockProofState
}

func NewProofStateManager() *ProofStateManager {
	return &ProofStateManager{
		states: make(map[uint64]*BlockProofState),
	}
}

func (m *ProofStateManager) create(blockID *big.Int) {
	blockIDUint64 := blockID.Uint64()

	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.states[blockIDUint64]
	if !ok {
		state = &BlockProofState{
			verifiedTiers: []uint16{},
			proofs:        []encoding.SubProof{},
		}
		m.states[blockIDUint64] = state
	}
}

func (m *ProofStateManager) containsTier(blockID *big.Int, tier uint16) bool {
	blockIDUint64 := blockID.Uint64()

	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.states[blockIDUint64]
	if !ok {
		return false
	}

	return slices.Contains(state.verifiedTiers, tier)
}

// addTierAndProof marks the tier as verified and adds the subproof to the block's state if
// the state has not yet collected enough proofs. It returns whether the required number
// of proofs has now been reached.
func (m *ProofStateManager) addTierAndProof(
	blockID *big.Int,
	tier uint16,
	subProof encoding.SubProof,
	requiredProofs uint8,
) (reachedRequired bool) {
	blockIDUint64 := blockID.Uint64()

	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.states[blockIDUint64]
	if !ok {
		// create if it doesn't exist
		state = &BlockProofState{
			verifiedTiers: []uint16{},
			proofs:        []encoding.SubProof{},
		}
		m.states[blockIDUint64] = state
	}

	// Record that we've verified this tier
	state.verifiedTiers = append(state.verifiedTiers, tier)

	// Only add to the subproofs if we haven't reached requiredProofs count yet
	if uint8(len(state.proofs)) < requiredProofs {
		state.proofs = append(state.proofs, subProof)
	}

	// Return true if we've now collected enough proofs
	return uint8(len(state.proofs)) == requiredProofs
}

func (m *ProofStateManager) currentProofCount(blockID *big.Int) int {
	blockIDUint64 := blockID.Uint64()

	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.states[blockIDUint64]
	if !ok {
		return 0
	}
	return len(state.proofs)
}

func (m *ProofStateManager) encodeSubProofs(blockID *big.Int) ([]byte, error) {
	blockIDUint64 := blockID.Uint64()

	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.states[blockIDUint64]
	if !ok {
		return nil, fmt.Errorf("block proof state not found for blockID: %d", blockIDUint64)
	}

	return encoding.EncodeSubProofs(state.proofs)
}

// cleanOldProofStates removes proof states for blocks older than `blockHistoryLength` blocks.
func (m *ProofStateManager) cleanOldProofStates(latestBlockID *big.Int, blockHistoryLength uint64) {
	blockID := latestBlockID.Uint64()

	m.mu.Lock()
	defer m.mu.Unlock()

	threshold := blockID - blockHistoryLength
	for key := range m.states {
		if key < threshold {
			delete(m.states, key)
		}
	}
}

// CombinedProducer generates proofs from multiple producers in parallel and combines them.
type CombinedProducer struct {
	ProofTier      uint16
	RequiredProofs uint8
	Producers      []ProofProducer
	Verifiers      []common.Address

	// Thread-safe manager for block proof states
	ProofStates *ProofStateManager
}

const (
	// BlockHistoryLength represents the number of blocks to keep in history of proof states.
	BlockHistoryLength = 256
)

// RequestProof implements the ProofProducer interface.
func (c *CombinedProducer) RequestProof(
	ctx context.Context,
	opts *ProofRequestOptions,
	blockID *big.Int,
	meta metadata.TaikoBlockMetaData,
	header *types.Header,
	requestAt time.Time,
) (*ProofWithHeader, error) {
	log.Debug("CombinedProducer: RequestProof", "blockID", blockID, "Producers num", len(c.Producers))

	var (
		wg         sync.WaitGroup
		errorsChan = make(chan error, len(c.Producers))
	)

	// We create the proof state to track progress
	c.ProofStates.create(blockID)

	taskCtx, taskCtxCancel := context.WithCancel(ctx)
	defer taskCtxCancel()

	for i, producer := range c.Producers {
		tier := producer.Tier()

		if c.ProofStates.containsTier(blockID, tier) {
			log.Debug("Skipping producer, proof already verified", "tier", tier)
			continue
		}

		log.Debug("Adding proof producer", "blockID", blockID, "tier", tier)
		verifier := c.Verifiers[i]

		wg.Add(1)
		go func(idx int, p ProofProducer, v common.Address) {
			defer wg.Done()

			proofWithHeader, err := p.RequestProof(taskCtx, opts, blockID, meta, header, requestAt)
			if err != nil {
				errorsChan <- fmt.Errorf("producer %d error: %w", idx, err)
				return
			}

			reachedRequired := c.ProofStates.addTierAndProof(blockID, p.Tier(), encoding.SubProof{
				Proof:    proofWithHeader.Proof,
				Verifier: v,
			}, c.RequiredProofs)

			if reachedRequired {
				taskCtxCancel()
			}
		}(i, producer, verifier)
	}

	wg.Wait()
	close(errorsChan)

	currentProofCount := c.ProofStates.currentProofCount(blockID)
	if uint8(currentProofCount) < c.RequiredProofs {
		// Not enough proofs, gather errors
		var errMsgs []string
		errMsgs = append(errMsgs,
			fmt.Sprintf("not enough proofs collected: required %d, got %d", c.RequiredProofs, currentProofCount),
		)
		for err := range errorsChan {
			errMsgs = append(errMsgs, err.Error())
		}
		return nil, fmt.Errorf("combined proof production failed: %s", strings.Join(errMsgs, "; "))
	}

	combinedProof, err := c.ProofStates.encodeSubProofs(blockID)
	if err != nil {
		return nil, fmt.Errorf("failed to encode sub proofs: %w", err)
	}

	log.Info(
		"Combined proofs generated",
		"blockID", blockID,
		"time", time.Since(requestAt),
		"producer", "CombinedProducer",
	)

	c.ProofStates.cleanOldProofStates(blockID, BlockHistoryLength)

	return &ProofWithHeader{
		BlockID: blockID,
		Header:  header,
		Meta:    meta,
		Proof:   combinedProof,
		Opts:    opts,
		Tier:    c.Tier(),
	}, nil
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
