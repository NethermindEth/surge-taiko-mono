package producer

import (
	"context"
	"fmt"
	"math/big"
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
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		proofs     = make([]encoding.SubProof, 0, len(c.Producers))
		errorsChan = make(chan error, len(c.Producers))
	)

	taskCtx, taskCtxCancel := context.WithCancel(ctx)
	defer taskCtxCancel()

	for i, producer := range c.Producers {
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

			if uint8(len(proofs)) < c.RequiredProofs {
				proofs = append(
					proofs,
					encoding.SubProof{
						Proof:    proofWithHeader.Proof,
						Verifier: verifier,
					},
				)
			}

			if uint8(len(proofs)) == c.RequiredProofs {
				taskCtxCancel()
			}
		}(i, producer, verifier)
	}

	wg.Wait()

	if uint8(len(proofs)) < c.RequiredProofs {
		var errMsgs []string

		errMsgs = append(
			errMsgs,
			fmt.Sprintf("not enough proofs collected: required %d, got %d", c.RequiredProofs, len(proofs)),
		)

		for err := range errorsChan {
			errMsgs = append(errMsgs, err.Error())
		}

		return nil, fmt.Errorf("combined proof production failed: %s", strings.Join(errMsgs, "; "))
	}

	combinedProof, err := encoding.EncodeSubProofs(proofs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode sub proofs: %w", err)
	}

	log.Info(
		"Combined proofs generated",
		"blockID", blockID,
		"time", time.Since(requestAt),
		"producer", "CombinedProducer",
	)

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
