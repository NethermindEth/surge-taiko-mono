package prover

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	proofProducer "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_producer"
)

const ProofStatusComplete = "complete"

// ProverAPIServer represents a prover API server instance.
type ProverAPIServer struct {
	echo *echo.Echo
	// Proof producers
	baseLevelProofProducer proofProducer.ProofProducer
	zkvmProofProducer      proofProducer.ProofProducer
	// Prover address
	proverAddress common.Address
}

// ProofRequestCheckpoint represents the checkpoint data in a proof request.
type ProofRequestCheckpoint struct {
	BlockNumber uint64      `json:"blockNumber"`
	BlockHash   common.Hash `json:"blockHash"`
	StateRoot   common.Hash `json:"stateRoot"`
}

// ProofRequest represents the request body for proof requests (limp mode).
type ProofRequest struct {
	ProposalID            uint64                    `json:"proposalId"`
	L1LimpData            *proofProducer.L1LimpData `json:"l1LimpData"`
	L2BlockNumbers        []uint64                  `json:"l2BlockNumbers"`
	Checkpoint            *ProofRequestCheckpoint   `json:"checkpoint"`
	DesignatedProver      common.Address            `json:"designatedProver"`
	LastAnchorBlockNumber uint64                    `json:"lastAnchorBlockNumber"`
}

// ProofResponse represents the response body for proof requests.
// The caller should poll until Status is "complete".
// Status values: "complete", "work_in_progress", "registered", or error message.
type ProofResponse struct {
	Status    string                  `json:"status"`
	Proof     []byte                  `json:"proof,omitempty"`
	ProofType proofProducer.ProofType `json:"proofType,omitempty"`
}

// NewProverAPIServer creates a new prover API server instance.
func NewProverAPIServer(
	cors string,
	jwtSecret []byte,
	baseLevelProofProducer proofProducer.ProofProducer,
	zkvmProofProducer proofProducer.ProofProducer,
	proverAddress common.Address,
) *ProverAPIServer {
	server := &ProverAPIServer{
		echo:                   echo.New(),
		baseLevelProofProducer: baseLevelProofProducer,
		zkvmProofProducer:      zkvmProofProducer,
		proverAddress:          proverAddress,
	}

	server.echo.HideBanner = true
	server.configureMiddleware([]string{cors})
	server.configureRoutes()
	if jwtSecret != nil {
		server.echo.Use(echojwt.JWT(jwtSecret))
	}

	return server
}

// configureMiddleware configures the server middlewares.
func (s *ProverAPIServer) configureMiddleware(corsOrigins []string) {
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowCredentials: true,
	}))
}

// configureRoutes contains all routes which will be used by the HTTP server.
func (s *ProverAPIServer) configureRoutes() {
	s.echo.GET("/", s.HealthCheck)
	s.echo.GET("/healthz", s.HealthCheck)
	s.echo.POST("/proof", s.RequestProof)
}

// Start starts the HTTP server.
func (s *ProverAPIServer) Start(port uint64) error {
	return s.echo.Start(fmt.Sprintf(":%v", port))
}

// Shutdown shuts down the HTTP server.
func (s *ProverAPIServer) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

// HealthCheck is the health check endpoint.
func (s *ProverAPIServer) HealthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// RequestProof handles proof requests for limp mode.
// This follows the same pattern as the existing proof request flow:
// - Makes a single request to the proof producer
// - Returns immediately with the current status
// - Caller is responsible for polling until proof is complete
func (s *ProverAPIServer) RequestProof(c echo.Context) error {
	ctx := c.Request().Context()

	var req ProofRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ProofResponse{
			Status: fmt.Sprintf("invalid request body: %v", err),
		})
	}

	if err := s.validateProofRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ProofResponse{
			Status: err.Error(),
		})
	}

	log.Debug(
		"Received proof request",
		"proposalID", req.ProposalID,
		"l2BlockCount", len(req.L2BlockNumbers),
	)

	// Build proof request options
	l2BlockNums := make([]*big.Int, len(req.L2BlockNumbers))
	for i, num := range req.L2BlockNumbers {
		l2BlockNums[i] = new(big.Int).SetUint64(num)
	}

	opts := &proofProducer.ProofRequestOptionsShasta{
		ProposalID:    new(big.Int).SetUint64(req.ProposalID),
		ProverAddress: s.proverAddress,
		L2BlockNums:   l2BlockNums,
		Checkpoint: &proofProducer.Checkpoint{
			BlockNumber: new(big.Int).SetUint64(req.Checkpoint.BlockNumber),
			BlockHash:   req.Checkpoint.BlockHash,
			StateRoot:   req.Checkpoint.StateRoot,
		},
		DesignatedProver:      req.DesignatedProver,
		LastAnchorBlockNumber: new(big.Int).SetUint64(req.LastAnchorBlockNumber),
		L1LimpData:            req.L1LimpData,
		Headers: []*types.Header{{
			Number: new(big.Int).SetUint64(req.Checkpoint.BlockNumber),
		}},
	}

	// Select proof producer (prefer ZK if available)
	producer := s.baseLevelProofProducer
	if s.zkvmProofProducer != nil {
		producer = s.zkvmProofProducer
	}

	// Make a single request - caller handles polling
	proofResp, err := producer.RequestProof(
		ctx,
		opts,
		new(big.Int).SetUint64(req.ProposalID),
		nil, // No metadata for limp mode
		time.Now(),
	)

	// Handle proof-in-progress and retry errors
	if err != nil {
		// Return the error string as the status, e.g., "work_in_progress",
		// "registered", or "zk_any_not_drawn".
		if errors.Is(err, proofProducer.ErrProofInProgress) ||
			errors.Is(err, proofProducer.ErrRetry) ||
			errors.Is(err, proofProducer.ErrZkAnyNotDrawn) {
			return c.JSON(http.StatusOK, ProofResponse{
				Status: err.Error(),
			})
		}
		log.Error("Proof request failed", "proposalID", req.ProposalID, "error", err)
		return c.JSON(http.StatusInternalServerError, ProofResponse{
			Status: err.Error(),
		})
	}

	log.Info(
		"Proof generated",
		"proposalID", req.ProposalID,
		"proofType", proofResp.ProofType,
	)

	return c.JSON(http.StatusOK, ProofResponse{
		Status:    ProofStatusComplete,
		Proof:     proofResp.Proof,
		ProofType: proofResp.ProofType,
	})
}

// validateProofRequest validates the proof request.
func (s *ProverAPIServer) validateProofRequest(req *ProofRequest) error {
	if req.ProposalID == 0 {
		return errors.New("proposalId is required")
	}
	if req.L1LimpData == nil {
		return errors.New("l1LimpData is required")
	}
	if req.L1LimpData.ExpectedProposeEvent.Id == "" {
		return errors.New("l1LimpData.expected_propose_event.id is required")
	}
	if len(req.L1LimpData.Blobs.Data) == 0 {
		return errors.New("l1LimpData.blobs.data is required and cannot be empty")
	}
	if len(req.L2BlockNumbers) == 0 {
		return errors.New("l2BlockNumbers is required")
	}
	if req.Checkpoint == nil {
		return errors.New("checkpoint is required")
	}
	if req.DesignatedProver == (common.Address{}) {
		return errors.New("designatedProver is required")
	}
	return nil
}
