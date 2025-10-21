package formatter

import (
	"fmt"
	"strings"

	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
)

type RISC0Formatter struct {
	log *logger.Logger
}

func NewRISC0Formatter(log *logger.Logger) *RISC0Formatter {
	return &RISC0Formatter{log: log}
}

type RISC0ProcessedData struct {
	AggregationProgramHash string `json:"aggregation_program_hash"`
	BlockProgramHash       string `json:"block_program_hash"`
}

func (f *RISC0Formatter) ProcessGuestData(guestData map[string]interface{}) (*RISC0ProcessedData, error) {
	f.log.Info("processing RISC0 guest data")

	// Extract RISC0 specific data
	risc0Data, ok := guestData["risc0"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'risc0' field in guest data")
	}

	aggregationHash, ok := risc0Data["aggregation_program_hash"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'aggregation_program_hash' field in RISC0 data")
	}

	blockHash, ok := risc0Data["block_program_hash"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'block_program_hash' field in RISC0 data")
	}

	// Ensure hex strings have 0x prefix
	if !strings.HasPrefix(aggregationHash, "0x") {
		aggregationHash = "0x" + aggregationHash
	}
	if !strings.HasPrefix(blockHash, "0x") {
		blockHash = "0x" + blockHash
	}

	processedData := &RISC0ProcessedData{
		AggregationProgramHash: aggregationHash,
		BlockProgramHash:       blockHash,
	}

	f.log.Info("RISC0 data processed successfully")
	return processedData, nil
}
