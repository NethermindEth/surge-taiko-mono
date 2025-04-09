package producer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

var (
	errProofInProgress = fmt.Errorf("sgx proof is still being generated")
)

// subProofStatus is just an internal state for each block’s proof.
type subProofStatus int

const (
	subProofStatusNew subProofStatus = iota
	subProofStatusInProgress
	subProofStatusDone
)

// sgxProofCache holds the cached proof/result for one block.
type sgxProofCache struct {
	status  subProofStatus
	proof   []byte
	lastErr error
}

const (
	ProofTypeSgx = "sgx"
	ProofTypeCPU = "native"
)

// SGXProofProducer generates a SGX proof for the given block.
type SGXProofProducer struct {
	RaikoHostEndpoint   string // a proverd RPC endpoint
	ProofType           string // Proof type
	JWT                 string // JWT provided by Raiko
	Dummy               bool
	RaikoRequestTimeout time.Duration
	DummyProofProducer

	// 1) Add a concurrency-safe map from blockID -> sgxProofCache
	proofCache map[uint64]*sgxProofCache
	cacheMutex sync.Mutex
}

// RaikoRequestProofBody represents the JSON body for requesting the proof.
type RaikoRequestProofBody struct {
	Block    *big.Int                    `json:"block_number"`
	Prover   string                      `json:"prover"`
	Graffiti string                      `json:"graffiti"`
	Type     string                      `json:"proof_type"`
	SGX      *SGXRequestProofBodyParam   `json:"sgx"`
	RISC0    *RISC0RequestProofBodyParam `json:"risc0"`
	SP1      *SP1RequestProofBodyParam   `json:"sp1"`
}

// SGXRequestProofBodyParam represents the JSON body of RaikoRequestProofBody's `sgx` field.
type SGXRequestProofBodyParam struct {
	Setup     bool `json:"setup"`
	Bootstrap bool `json:"bootstrap"`
	Prove     bool `json:"prove"`
}

// RISC0RequestProofBodyParam represents the JSON body of RaikoRequestProofBody's `risc0` field.
type RISC0RequestProofBodyParam struct {
	Bonsai       bool     `json:"bonsai"`
	Snark        bool     `json:"snark"`
	Profile      bool     `json:"profile"`
	ExecutionPo2 *big.Int `json:"execution_po2"`
}

// SP1RequestProofBodyParam represents the JSON body of RaikoRequestProofBody's `sp1` field.
type SP1RequestProofBodyParam struct {
	Recursion string `json:"recursion"`
	Prover    string `json:"prover"`
}

// RaikoRequestProofBodyResponse represents the JSON body of the response of the proof requests.
type RaikoRequestProofBodyResponse struct {
	Data         *RaikoProofData `json:"data"`
	ErrorMessage string          `json:"message"`
}

type RaikoProofData struct {
	Proof  string `json:"proof"` //nolint:revive,stylecheck
	Status string `json:"status"`
}

// RequestProof implements the ProofProducer interface.
func (s *SGXProofProducer) RequestProof(
	ctx context.Context,
	opts *ProofRequestOptions,
	blockID *big.Int,
	meta metadata.TaikoBlockMetaData,
	header *types.Header,
	requestAt time.Time,
) (*ProofWithHeader, error) {
	if s.proofCache == nil {
		s.proofCache = make(map[uint64]*sgxProofCache)
	}

	log.Info(
		"Request sgx proof from raiko-host service",
		"blockID", blockID,
		"coinbase", meta.GetCoinbase(),
		"height", header.Number,
		"hash", header.Hash(),
	)

	if s.Dummy {
		return s.DummyProofProducer.RequestProof(opts, blockID, meta, header, s.Tier(), requestAt)
	}

	// 2) Check or create the cache entry for this block
	bid := blockID.Uint64()

	s.cacheMutex.Lock()
	cache, ok := s.proofCache[bid]
	if !ok {
		cache = &sgxProofCache{status: subProofStatusNew}
		s.proofCache[bid] = cache
	}

	log.Info("============== sgx_producer.go: proof cache", bid, cache.status)

	if cache.status == subProofStatusNew {
		log.Info("================= hey hey sgx_producer.go: proof new", bid)
	}

	switch cache.status {

	// case subProofStatusInProgress:
	// 	// This block is still generating a proof
	// 	s.cacheMutex.Unlock()

	// 	log.Info("================= sgx_producer.go: proof in progress", bid)

	// 	return nil, errProofInProgress

	case subProofStatusDone:
		// We’ve already finalized a result for this block. Return the proof or the error we cached.
		proof := cache.proof
		lastErr := cache.lastErr
		s.cacheMutex.Unlock()

		log.Info("================= sgx_producer.go: proof done", bid)

		if lastErr != nil {
			// e.g., if we ended in an error last time, we return that error
			return nil, lastErr
		}

		// Otherwise, return the cached proof
		return &ProofWithHeader{
			BlockID: blockID,
			Header:  header,
			Meta:    meta,
			Proof:   proof,
			Opts:    opts,
			Tier:    s.Tier(),
		}, nil

	case subProofStatusNew, subProofStatusInProgress:
		log.Info("================= sgx_producer.go: proof new or in progress", bid)

		// Mark as in-progress, then generate the proof
		cache.status = subProofStatusInProgress
		s.cacheMutex.Unlock()

		// We'll do the actual proof generation outside the lock
		proofBytes, err := s.callProverDaemon(ctx, opts, requestAt)

		// Re-lock to update the cache
		s.cacheMutex.Lock()
		defer s.cacheMutex.Unlock()

		if err != nil {
			cache.lastErr = err
			if err == errProofGenerating {
				log.Info("======================= at sgx_producer.go received proof generating error", err)
				return nil, err
			} else {
				log.Info("======================= at sgx_producer.go received bad bad error", err)
				// cache.status = subProofStatusDone
				return nil, err
			}
		}

		// If we got a valid proof:
		cache.status = subProofStatusDone
		cache.proof = proofBytes
		cache.lastErr = nil

		// Return final result
		return &ProofWithHeader{
			BlockID: blockID,
			Header:  header,
			Meta:    meta,
			Proof:   proofBytes,
			Opts:    opts,
			Tier:    s.Tier(),
		}, nil

	default:
		log.Info("================= sgx_producer.go: unhandled status", bid, cache.status)
	}

	// Should never happen, but just in case:
	s.cacheMutex.Unlock()
	return nil, fmt.Errorf("unhandled status for block %d", bid)
}

func (s *SGXProofProducer) RequestCancel(
	_ context.Context,
	_ *ProofRequestOptions,
) error {
	return nil
}

// callProverDaemon keeps polling the proverd service to get the requested proof.
func (s *SGXProofProducer) callProverDaemon(
	ctx context.Context,
	opts *ProofRequestOptions,
	requestAt time.Time,
) ([]byte, error) {
	var (
		proof []byte
	)

	ctx, cancel := rpc.CtxWithTimeoutOrDefault(ctx, s.RaikoRequestTimeout)
	defer cancel()

	output, err := s.requestProof(ctx, opts)
	if err != nil {
		log.Error("Failed to request proof", "height", opts.BlockID, "error", err, "endpoint", s.RaikoHostEndpoint)
		return nil, err
	}

	if output == nil {
		log.Info(
			"Proof generating",
			"height", opts.BlockID,
			"time", time.Since(requestAt),
			"producer", "SGXProofProducer",
		)
		return nil, errProofGenerating
	}

	// Raiko returns "" as proof when proof type is native,
	// so we just convert "" to bytes
	if s.ProofType == ProofTypeCPU {
		proof = common.Hex2Bytes(output.Data.Proof)
	} else {
		if len(output.Data.Proof) == 0 {
			return nil, errEmptyProof
		}
		proof = common.Hex2Bytes(output.Data.Proof[2:])
	}

	log.Info(
		"Proof generated",
		"height", opts.BlockID,
		"time", time.Since(requestAt),
		"producer", "SGXProofProducer",
	)

	return proof, nil
}

// requestProof sends a RPC request to proverd to try to get the requested proof.
func (s *SGXProofProducer) requestProof(
	ctx context.Context,
	opts *ProofRequestOptions,
) (*RaikoRequestProofBodyResponse, error) {
	reqBody := RaikoRequestProofBody{
		Type:     s.ProofType,
		Block:    opts.BlockID,
		Prover:   opts.ProverAddress.Hex()[2:],
		Graffiti: opts.Graffiti,
		SGX: &SGXRequestProofBodyParam{
			Setup:     false,
			Bootstrap: false,
			Prove:     true,
		},
	}

	client := &http.Client{}

	jsonValue, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.RaikoHostEndpoint+"/v1/proof", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(s.JWT) > 0 {
		req.Header.Set("Authorization", "Bearer "+base64.StdEncoding.EncodeToString([]byte(s.JWT)))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request proof, id: %d, statusCode: %d", opts.BlockID, res.StatusCode)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	log.Debug(
		"Proof generation output",
		"blockID", opts.BlockID,
		"zkType", "sgx",
		"output", string(resBytes),
	)

	var output RaikoRequestProofBodyResponse
	if err := json.Unmarshal(resBytes, &output); err != nil {
		return nil, err
	}

	if len(output.ErrorMessage) > 0 {
		return nil, fmt.Errorf("failed to get proof, msg: %s", output.ErrorMessage)
	}

	return &output, nil
}

// Tier implements the ProofProducer interface.
func (s *SGXProofProducer) Tier() uint16 {
	return encoding.TierSgxID
}
