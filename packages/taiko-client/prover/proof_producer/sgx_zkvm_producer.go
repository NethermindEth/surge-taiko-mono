package producer

import (
	"bytes"
	"context"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/metrics"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// SGXProofProducer generates a SGX proof for the given block.
type SGXAndZkVMProofProducer struct {
	RaikoHostEndpoint    string // a proverd RPC endpoint
	SgxProofType         string // SGX Proof type
	ZKProofType          string // ZK Proof type
	JWT                  string // JWT provided by Raiko
	Dummy                bool
	Risc0VerifierAddress common.Address
	SgxVerifierAddress   common.Address
	RaikoRequestTimeout  time.Duration
	DummyProofProducer
}

// RequestProof implements the ProofProducer interface.
func (s *SGXAndZkVMProofProducer) RequestProof(
	ctx context.Context,
	opts *ProofRequestOptions,
	blockID *big.Int,
	meta metadata.TaikoBlockMetaData,
	header *types.Header,
	requestAt time.Time,
) (*ProofWithHeader, error) {
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

	proof, err := s.callProverDaemon(ctx, opts, requestAt)
	if err != nil {
		return nil, err
	}

	metrics.ProverSgxProofGeneratedCounter.Add(1)

	return &ProofWithHeader{
		BlockID: blockID,
		Header:  header,
		Meta:    meta,
		Proof:   proof,
		Opts:    opts,
		Tier:    s.Tier(),
	}, nil
}

func (s *SGXAndZkVMProofProducer) RequestCancel(
	ctx context.Context,
	opts *ProofRequestOptions,
) error {
	reqBody := RaikoRequestProofBody{
		Type:     s.ZKProofType,
		Block:    opts.BlockID,
		Prover:   opts.ProverAddress.Hex()[2:],
		Graffiti: opts.Graffiti,
		RISC0: &RISC0RequestProofBodyParam{
			Bonsai:       true,
			Snark:        true,
			Profile:      false,
			ExecutionPo2: big.NewInt(20),
		},
	}

	client := &http.Client{}

	jsonValue, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.RaikoHostEndpoint+"/v2/proof/cancel",
		bytes.NewBuffer(jsonValue),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(s.JWT) > 0 {
		req.Header.Set("Authorization", "Bearer "+base64.StdEncoding.EncodeToString([]byte(s.JWT)))
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to cancel requesting proof, statusCode: %d", res.StatusCode)
	}

	return nil
}

// callProverDaemon keeps polling the proverd service to get the requested proof.
func (s *SGXAndZkVMProofProducer) callProverDaemon(
	ctx context.Context,
	opts *ProofRequestOptions,
	requestAt time.Time,
) ([]byte, error) {
	var (
		proof []byte
	)

	// SGX

	sgxCtx, SgxCancel := rpc.CtxWithTimeoutOrDefault(ctx, s.RaikoRequestTimeout)
	defer SgxCancel()

	sgxOutput, err := s.requestSgxProof(sgxCtx, opts)
	if err != nil {
		log.Error("Failed to request proof", "height", opts.BlockID, "error", err, "endpoint", s.RaikoHostEndpoint)
		return nil, err
	}

	if sgxOutput == nil {
		log.Info(
			"Proof generating",
			"height", opts.BlockID,
			"time", time.Since(requestAt),
			"producer", "SGXProofProducer",
		)
		return nil, errProofGenerating
	}

	var sgxProof []byte

	// Raiko returns "" as proof when proof type is native,
	// so we just convert "" to bytes
	if s.SgxProofType == ProofTypeCPU {
		sgxProof = common.Hex2Bytes(sgxOutput.Data.Proof)
	} else {
		if len(sgxOutput.Data.Proof) == 0 {
			return nil, errEmptyProof
		}
		sgxProof = common.Hex2Bytes(sgxOutput.Data.Proof[2:])
	}

	log.Info(
		"SGX Proof generated",
		"height", opts.BlockID,
		"time", time.Since(requestAt),
		"producer", "SGXAndZkVMProofProducer",
	)

	// ZkVM

	zkCtx, zkCancel := rpc.CtxWithTimeoutOrDefault(ctx, s.RaikoRequestTimeout)
	defer zkCancel()

	zkOutput, err := s.requestZkProof(zkCtx, opts)
	if err != nil {
		log.Error("Failed to request proof", "height", opts.BlockID, "error", err, "endpoint", s.RaikoHostEndpoint)
		return nil, err
	}

	if zkOutput.Data.Status == ErrProofInProgress.Error() {
		return nil, ErrProofInProgress
	}
	if zkOutput.Data.Status == StatusRegistered {
		return nil, ErrRetry
	}

	var zkProof []byte

	if len(zkOutput.Data.Proof.Proof) == 0 {
		return nil, errEmptyProof
	}
	zkProof = common.Hex2Bytes(zkOutput.Data.Proof.Proof[2:])
	log.Info(
		"ZkVM Proof generated",
		"height", opts.BlockID,
		"time", time.Since(requestAt),
		"producer", "SGXAndZkVMProofProducer",
	)

	proof, err = encoding.EncodeSubProofs([]encoding.SubProof{
		{
			Verifier: s.Risc0VerifierAddress,
			Proof:    zkProof,
		},
		{
			Verifier: s.SgxVerifierAddress,
			Proof:    sgxProof,
		},
	})
	if err != nil {
		log.Error("Failed to encode sub proofs", "height", opts.BlockID, "error", err, "endpoint", s.RaikoHostEndpoint)
		return nil, err
	}

	return proof, nil
}

// requestSgxProof sends a RPC request to proverd to try to get the requested SGX proof.
func (s *SGXAndZkVMProofProducer) requestSgxProof(
	ctx context.Context,
	opts *ProofRequestOptions,
) (*RaikoRequestProofBodyResponse, error) {
	reqBody := RaikoRequestProofBody{
		Type:     s.SgxProofType,
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

// requestZkProof sends a RPC request to proverd to try to get the requested proof.
func (s *SGXAndZkVMProofProducer) requestZkProof(
	ctx context.Context,
	opts *ProofRequestOptions,
) (*RaikoRequestProofBodyResponseV2, error) {
	var reqBody RaikoRequestProofBody
	switch s.ZKProofType {
	case ZKProofTypeSP1:
		reqBody = RaikoRequestProofBody{
			Type:     s.ZKProofType,
			Block:    opts.BlockID,
			Prover:   opts.ProverAddress.Hex()[2:],
			Graffiti: opts.Graffiti,
			SP1: &SP1RequestProofBodyParam{
				Recursion: "plonk",
				Prover:    "network",
			},
		}
	default:
		reqBody = RaikoRequestProofBody{
			Type:     s.ZKProofType,
			Block:    opts.BlockID,
			Prover:   opts.ProverAddress.Hex()[2:],
			Graffiti: opts.Graffiti,
			RISC0: &RISC0RequestProofBodyParam{
				Bonsai:       true,
				Snark:        true,
				Profile:      false,
				ExecutionPo2: big.NewInt(20),
			},
		}
	}

	client := &http.Client{}

	jsonValue, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.RaikoHostEndpoint+"/v2/proof", bytes.NewBuffer(jsonValue))
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
		"zkType", s.ZKProofType,
		"output", string(resBytes),
	)
	var output RaikoRequestProofBodyResponseV2
	if err := json.Unmarshal(resBytes, &output); err != nil {
		return nil, err
	}

	if len(output.ErrorMessage) > 0 {
		return nil, fmt.Errorf("failed to get proof, msg: %s", output.ErrorMessage)
	}

	return &output, nil
}

// Tier implements the ProofProducer interface.
func (s *SGXAndZkVMProofProducer) Tier() uint16 {
	return encoding.TierSgxAndZkVMID
}
