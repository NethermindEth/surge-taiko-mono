package formatter

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/go-tdx-guest/abi"
	"github.com/google/go-tdx-guest/proto/tdx"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
)

type TDXFormatter struct {
	log *logger.Logger
}

func NewTDXFormatter(log *logger.Logger) *TDXFormatter {
	return &TDXFormatter{log: log}
}

// TDXProcessedData represents processed data for regular TDX (not Azure TDX)
type TDXProcessedData struct {
	Quote      HexBytes       `json:"quote"`
	Nonce      HexBytes       `json:"nonce"`
	IssuerType string         `json:"issuerType"`
	PublicKey  string         `json:"publicKey"`
	Metadata   map[string]any `json:"metadata"`
}

func (f *TDXFormatter) ProcessGuestData(guestData map[string]interface{}) (*TDXProcessedData, error) {
	f.log.Info("processing TDX guest data")

	// For regular TDX, we just extract the raw quote and metadata
	quoteStr, ok := guestData["quote"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'quote' field in guest data")
	}

	nonceStr, _ := guestData["nonce"].(string)
	issuerType, _ := guestData["issuer_type"].(string)
	publicKey, _ := guestData["public_key"].(string)
	metadata, _ := guestData["metadata"].(map[string]interface{})

	// Ensure quote has 0x prefix
	if !strings.HasPrefix(quoteStr, "0x") {
		quoteStr = "0x" + quoteStr
	}

	// Decode quote from hex
	quoteBytes, err := hex.DecodeString(strings.TrimPrefix(quoteStr, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode quote hex: %w", err)
	}

	var nonceBytes []byte
	if nonceStr != "" {
		nonceBytes, _ = hex.DecodeString(strings.TrimPrefix(nonceStr, "0x"))
	}

	processedData := &TDXProcessedData{
		Quote:      HexBytes(quoteBytes),
		Nonce:      HexBytes(nonceBytes),
		IssuerType: issuerType,
		PublicKey:  publicKey,
		Metadata:   metadata,
	}

	f.log.Info("TDX data processed successfully")
	return processedData, nil
}

// ExtractTrustedParams extracts the trusted parameters from regular TDX attestation data
func (f *TDXFormatter) ExtractTrustedParams(data *TDXProcessedData) (*TDXTrustedParams, error) {
	f.log.Info("extracting trusted params from raw TDX attestation")

	// Parse the raw TDX quote to extract trusted parameters
	quoteProto, err := abi.QuoteToProto(data.Quote)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TDX quote: %w", err)
	}

	// Cast to QuoteV4 which contains the actual TDX measurements
	quoteV4, ok := quoteProto.(*tdx.QuoteV4)
	if !ok {
		return nil, fmt.Errorf("expected QuoteV4, got %T", quoteProto)
	}

	// Extract measurements from the TD quote body
	if quoteV4.TdQuoteBody == nil {
		return nil, fmt.Errorf("TDQuoteBody is nil")
	}

	body := quoteV4.TdQuoteBody

	// Extract TEE TCB SVN
	var teeTcbSvn [16]byte
	if len(body.TeeTcbSvn) == 16 {
		copy(teeTcbSvn[:], body.TeeTcbSvn)
	}

	// Extract individual RTMRs (Runtime Measurement Registers)
	// TDX has 4 RTMRs, each 48 bytes
	var rtMr0, rtMr1, rtMr2, rtMr3 []byte
	if len(body.Rtmrs) >= 1 && len(body.Rtmrs[0]) == 48 {
		rtMr0 = make([]byte, 48)
		copy(rtMr0, body.Rtmrs[0])
	}
	if len(body.Rtmrs) >= 2 && len(body.Rtmrs[1]) == 48 {
		rtMr1 = make([]byte, 48)
		copy(rtMr1, body.Rtmrs[1])
	}
	if len(body.Rtmrs) >= 3 && len(body.Rtmrs[2]) == 48 {
		rtMr2 = make([]byte, 48)
		copy(rtMr2, body.Rtmrs[2])
	}
	if len(body.Rtmrs) >= 4 && len(body.Rtmrs[3]) == 48 {
		rtMr3 = make([]byte, 48)
		copy(rtMr3, body.Rtmrs[3])
	}

	params := &TDXTrustedParams{
		TeeTcbSvn: teeTcbSvn,
		MrSeam:    body.MrSeam,
		MrTd:      body.MrTd,
		RtMr0:     rtMr0,
		RtMr1:     rtMr1,
		RtMr2:     rtMr2,
		RtMr3:     rtMr3,
	}

	f.log.Info("extracted trusted params successfully",
		"teeTcbSvn", hex.EncodeToString(teeTcbSvn[:]),
		"mrSeam", hex.EncodeToString(params.MrSeam),
		"mrTd", hex.EncodeToString(params.MrTd),
		"rtMr0", hex.EncodeToString(params.RtMr0),
		"rtMr1", hex.EncodeToString(params.RtMr1),
		"rtMr2", hex.EncodeToString(params.RtMr2),
		"rtMr3", hex.EncodeToString(params.RtMr3))

	return params, nil
}
