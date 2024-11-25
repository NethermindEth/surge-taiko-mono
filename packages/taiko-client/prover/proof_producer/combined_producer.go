package producer

import (
	"context"
	"fmt"
	"math/big"
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
	producers []ProofProducer
	verifiers []common.Address
	tier      uint16
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
		proofChan  = make(chan []byte, len(c.producers))
		errChan    = make(chan error, len(c.producers))
		proofs     = make([][]byte, len(c.producers))
		subProofs  = make([]encoding.SubProof, len(c.producers))
		finalError error
	)

	for i, producer := range c.producers {
		wg.Add(1)
		go func(idx int, p ProofProducer) {
			defer wg.Done()

			proofWithHeader, err := p.RequestProof(ctx, opts, blockID, meta, header, requestAt)
			if err != nil {
				errChan <- fmt.Errorf("producer %d error: %w", idx, err)
				return
			}

			proofChan <- proofWithHeader.Proof
		}(i, producer)
	}

	wg.Wait()
	close(proofChan)
	close(errChan)

	for err := range errChan {
		if finalError == nil {
			finalError = err
		}
	}

	if finalError != nil {
		return nil, finalError
	}

	i := 0
	for proof := range proofChan {
		proofs[i] = proof
		i++
	}

	for i := range proofs {
		subProofs[i] = encoding.SubProof{
			Verifier: c.verifiers[i],
			Proof:    proofs[i],
		}
	}

	combinedProof, err := encoding.EncodeSubProofs(subProofs)
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
	for _, producer := range c.producers {
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
	return c.tier
}
