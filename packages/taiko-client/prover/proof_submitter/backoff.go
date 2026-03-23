package submitter

import (
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	proofProducer "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_producer"
)

// newRaikoBackOff creates an exponential backoff for raiko proof polling.
// Starts at initialInterval and grows up to maxInterval on consecutive hard errors.
func newRaikoBackOff(initialInterval, maxInterval time.Duration) *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = initialInterval
	bo.MaxInterval = maxInterval
	bo.Multiplier = 2.0
	bo.RandomizationFactor = 0.3
	bo.MaxElapsedTime = 0
	bo.Reset()
	return bo
}

// isSoftRaikoError returns true for expected raiko states indicating healthy operation
// (proof still being generated). These reset the backoff so polling stays fast
// once raiko recovers from a hard failure.
func isSoftRaikoError(err error) bool {
	return errors.Is(err, proofProducer.ErrProofInProgress) ||
		errors.Is(err, proofProducer.ErrRetry)
}
