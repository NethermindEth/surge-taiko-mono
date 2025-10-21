package formatter

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
)

type SGXFormatter struct {
	log *logger.Logger
}

func NewSGXFormatter(log *logger.Logger) *SGXFormatter {
	return &SGXFormatter{log: log}
}

type SGXProcessedData struct {
	MrEnclave string `json:"mr_enclave"`
	MrSigner  string `json:"mr_signer"`
	Quote     string `json:"quote"`
}

func (f *SGXFormatter) ProcessGuestData(guestData map[string]interface{}) (*SGXProcessedData, error) {
	f.log.Info("processing SGX guest data")

	mrEnclave, ok := guestData["mr_enclave"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'mr_enclave' field in guest data")
	}

	mrSigner, ok := guestData["mr_signer"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'mr_signer' field in guest data")
	}

	quote, ok := guestData["quote"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'quote' field in guest data")
	}

	// Ensure hex strings have 0x prefix
	if !strings.HasPrefix(mrEnclave, "0x") {
		mrEnclave = "0x" + mrEnclave
	}
	if !strings.HasPrefix(mrSigner, "0x") {
		mrSigner = "0x" + mrSigner
	}
	if !strings.HasPrefix(quote, "0x") {
		quote = "0x" + quote
	}

	// Validate hex encoding
	if _, err := hex.DecodeString(strings.TrimPrefix(mrEnclave, "0x")); err != nil {
		return nil, fmt.Errorf("invalid mr_enclave hex: %w", err)
	}
	if _, err := hex.DecodeString(strings.TrimPrefix(mrSigner, "0x")); err != nil {
		return nil, fmt.Errorf("invalid mr_signer hex: %w", err)
	}
	if _, err := hex.DecodeString(strings.TrimPrefix(quote, "0x")); err != nil {
		return nil, fmt.Errorf("invalid quote hex: %w", err)
	}

	processedData := &SGXProcessedData{
		MrEnclave: mrEnclave,
		MrSigner:  mrSigner,
		Quote:     quote,
	}

	f.log.Info("SGX data processed successfully")
	return processedData, nil
}
