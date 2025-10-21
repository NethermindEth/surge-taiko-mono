package formatter

import (
	"fmt"
	"strings"

	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
)

type SP1Formatter struct {
	log *logger.Logger
}

func NewSP1Formatter(log *logger.Logger) *SP1Formatter {
	return &SP1Formatter{log: log}
}

type SP1ProcessedData struct {
	AggregationProgramHash string `json:"aggregation_program_hash"`
	BlockProgramHash       string `json:"block_program_hash"`
}

func (f *SP1Formatter) ProcessGuestData(guestData map[string]interface{}) (*SP1ProcessedData, error) {
	f.log.Info("processing SP1 guest data")

	// Extract SP1 specific data
	sp1Data, ok := guestData["sp1"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'sp1' field in guest data")
	}

	aggregationHash, ok := sp1Data["aggregation_program_hash"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'aggregation_program_hash' field in SP1 data")
	}

	blockHash, ok := sp1Data["block_program_hash"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'block_program_hash' field in SP1 data")
	}

	// Ensure hex strings have 0x prefix
	if !strings.HasPrefix(aggregationHash, "0x") {
		aggregationHash = "0x" + aggregationHash
	}
	if !strings.HasPrefix(blockHash, "0x") {
		blockHash = "0x" + blockHash
	}

	processedData := &SP1ProcessedData{
		AggregationProgramHash: aggregationHash,
		BlockProgramHash:       blockHash,
	}

	f.log.Info("SP1 data processed successfully")
	return processedData, nil
}
