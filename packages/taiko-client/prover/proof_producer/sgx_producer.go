package producer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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

// A type to track proof state for each block
type sgxProofStatus int

const (
	sgxProofStatusNew sgxProofStatus = iota
	sgxProofStatusInProgress
	sgxProofStatusDone
)

// Store the final proof or error
type sgxProofCache struct {
	status  sgxProofStatus
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

	bid := blockID.Uint64()

	s.cacheMutex.Lock()
	cache, exists := s.proofCache[bid]
	if !exists {
		cache = &sgxProofCache{status: sgxProofStatusNew}
		s.proofCache[bid] = cache
	}

	switch cache.status {
	case sgxProofStatusDone:
		proof := cache.proof
		lastErr := cache.lastErr
		s.cacheMutex.Unlock()

		if lastErr != nil {
			log.Info("SGX proof generation failed", "blockID", bid, "error", lastErr)
			return nil, lastErr
		}

		log.Info("SGX proof generation completed", "blockID", bid, "proof", proof)
		// Return the cached proof
		return &ProofWithHeader{
			BlockID: blockID,
			Header:  header,
			Meta:    meta,
			Proof:   proof,
			Opts:    opts,
			Tier:    s.Tier(),
		}, nil

	case sgxProofStatusNew, sgxProofStatusInProgress:
		log.Info("SGX proof generation in progress", "blockID", bid)

		// Mark as in-progress, then generate the proof
		cache.status = sgxProofStatusInProgress
	}
	s.cacheMutex.Unlock()

	// We'll do the actual proof generation outside the lock
	proofBytes, err := s.callProverDaemon(ctx, opts, requestAt)

	// Re-lock to update the cache
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if err != nil {
		// Received error, so mark as in-progress and let it retry
		cache.lastErr = err
		if errors.Is(err, errProofGenerating) {
			log.Info("SGX proof generation in progress, received normal error", "blockID", bid, "error", err)
			return nil, err
		} else {
			log.Info("SGX proof generation in progress, received unknown error", "blockID", bid, "error", err)
			return nil, err
		}
	}

	// If we got a valid proof:
	log.Info("SGX proof generation completed", "blockID", bid, "proof", proofBytes)
	cache.status = sgxProofStatusDone
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
