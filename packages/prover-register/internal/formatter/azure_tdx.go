package formatter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-tdx-guest/abi"
	"github.com/google/go-tdx-guest/proto/tdx"
	tpmproto "github.com/google/go-tpm-tools/proto/tpm"
	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/logger"
)

type AzureTDXFormatter struct {
	log *logger.Logger
}

func NewAzureTDXFormatter(log *logger.Logger) *AzureTDXFormatter {
	return &AzureTDXFormatter{log: log}
}

type AzureTDXInputData struct {
	RawQuote json.RawMessage `json:"rawQuote"`
	Nonce    HexBytes        `json:"nonce"`
}

type AzureTDXProcessedData struct {
	AttestationDocument struct {
		Attestation struct {
			TpmQuote TPMQuoteData `json:"tpmQuote"`
		} `json:"attestation"`
		InstanceInfo struct {
			AttestationReport HexBytes        `json:"attestationReport"`
			RuntimeData       RuntimeDataInfo `json:"runtimeData"`
		} `json:"instanceInfo"`
		UserData HexBytes `json:"userData"`
	} `json:"attestationDocument"`
	Pcrs           []PCRValue `json:"pcrs"`
	Nonce          HexBytes   `json:"nonce"`
	AdditionalData struct {
		RuntimeDataHash HexBytes32 `json:"runtimeDataHash"`
	} `json:"additionalData"`
}

func (f *AzureTDXFormatter) ProcessGuestData(guestData map[string]interface{}) (*AzureTDXProcessedData, error) {
	f.log.Info("processing Azure TDX guest data")

	// Extract the quote field - for Azure TDX, this is a hex-encoded JSON
	quoteStr, ok := guestData["quote"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'quote' field in guest data")
	}

	nonceStr, _ := guestData["nonce"].(string)

	// Decode the hex-encoded JSON quote
	quoteBytesHex := strings.TrimPrefix(quoteStr, "0x")
	quoteBytes, err := hex.DecodeString(quoteBytesHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode quote hex: %w", err)
	}

	// Parse the JSON quote
	var attestationDoc AttestationDocument
	if err := json.Unmarshal(quoteBytes, &attestationDoc); err != nil {
		return nil, fmt.Errorf("failed to parse attestation document JSON: %w", err)
	}

	// Create input data structure
	inputData := &AzureTDXInputData{
		RawQuote: json.RawMessage(quoteBytes),
	}

	if nonceStr != "" {
		nonceBytes, _ := hex.DecodeString(strings.TrimPrefix(nonceStr, "0x"))
		inputData.Nonce = HexBytes(nonceBytes)
	}

	// Process the attestation using the Azure TDX format
	return f.formatAttestation(inputData, &attestationDoc)
}

func (f *AzureTDXFormatter) formatAttestation(inputData *AzureTDXInputData, doc *AttestationDocument) (*AzureTDXProcessedData, error) {
	f.log.Info("formatting Azure TDX attestation")

	if doc.Attestation == nil {
		return nil, fmt.Errorf("attestation is nil")
	}

	// Parse instance info
	var instanceInfo InstanceInfo
	if err := json.Unmarshal(doc.InstanceInfo, &instanceInfo); err != nil {
		return nil, fmt.Errorf("parse instance info: %w", err)
	}

	// Process attestation data
	processedData, err := f.processAttestationData(doc, &instanceInfo)
	if err != nil {
		return nil, err
	}

	// Decode additional data
	decodedData, err := f.decodeAdditionalData(inputData, doc, &instanceInfo)
	if err != nil {
		return nil, err
	}

	// Build output
	output := f.buildOutput(inputData, processedData, decodedData)
	f.log.Info("Azure TDX attestation formatted successfully")
	return output, nil
}

func (f *AzureTDXFormatter) processAttestationData(doc *AttestationDocument, instanceInfo *InstanceInfo) (*processedAttestationData, error) {
	f.log.Debug("processing Azure TDX attestation data")

	// Decode HCL attestation key
	hclAkPub, err := tpm2.DecodePublic(doc.Attestation.AkPub)
	if err != nil {
		f.log.Error("failed to decode HCL attestation key", "error", err)
		return nil, fmt.Errorf("decode HCL attestation key: %w", err)
	}

	// Find SHA256 quote
	sha256Quote := f.findSHA256Quote(doc.Attestation.Quotes)
	if sha256Quote == nil {
		f.log.Error("no SHA256 quote found")
		return nil, fmt.Errorf("no SHA256 quote found")
	}
	f.log.Debug("found SHA256 quote")

	// Decode attestation report
	decodedAttestationReport, err := instanceInfo.DecodeAttestationReport()
	if err != nil {
		f.log.Error("failed to decode attestation report", "error", err)
		return nil, fmt.Errorf("decode attestation report: %w", err)
	}

	// Convert quote to proto
	quoteProto, err := abi.QuoteToProto(decodedAttestationReport)
	if err != nil {
		f.log.Error("failed to convert quote to proto", "error", err)
		return nil, fmt.Errorf("convert quote to proto: %w", err)
	}
	f.log.Debug("converted quote to proto", "quote", quoteProto)

	// Extract PCRs
	pcrs, err := f.extractPCRs(sha256Quote)
	if err != nil {
		return nil, err
	}
	f.log.Debug("extracted PCR values", "count", TPMPCRCount)

	// Decode TPM signature
	decodedSig, err := tpm2.DecodeSignature(bytes.NewBuffer(sha256Quote.RawSig))
	if err != nil {
		f.log.Error("failed to decode TPM signature", "error", err)
		return nil, fmt.Errorf("decode TPM signature: %w", err)
	}

	f.log.Info("successfully processed Azure TDX attestation data")
	return &processedAttestationData{
		hclAkPub:    &hclAkPub,
		sha256Quote: sha256Quote,
		pcrs:        pcrs,
		decodedSig:  decodedSig,
	}, nil
}

func (f *AzureTDXFormatter) findSHA256Quote(quotes []*tpmproto.Quote) *tpmproto.Quote {
	for _, quote := range quotes {
		if quote != nil && quote.Pcrs != nil && quote.Pcrs.Hash == tpmproto.HashAlgo_SHA256 {
			return quote
		}
	}
	return nil
}

func (f *AzureTDXFormatter) extractPCRs(quote *tpmproto.Quote) ([TPMPCRCount]HexBytes32, error) {
	var pcrs [TPMPCRCount]HexBytes32

	if quote.Pcrs == nil || quote.Pcrs.Pcrs == nil {
		return pcrs, fmt.Errorf("PCRs not found in quote")
	}

	for i, pcr := range quote.Pcrs.Pcrs {
		if i >= TPMPCRCount {
			break
		}
		copy(pcrs[i][:], pcr[:])
	}

	return pcrs, nil
}

func (f *AzureTDXFormatter) decodeAdditionalData(inputData *AzureTDXInputData, doc *AttestationDocument, instanceInfo *InstanceInfo) (*decodedAdditionalData, error) {
	f.log.Debug("decoding additional data fields")

	attestationReport, err := instanceInfo.DecodeAttestationReport()
	if err != nil {
		return nil, err
	}
	f.log.Debug("decoded attestation report", "size", len(attestationReport))

	runtimeData, err := instanceInfo.DecodeRuntimeData()
	if err != nil {
		return nil, err
	}
	f.log.Debug("decoded runtime data", "size", len(runtimeData))

	userData, err := doc.DecodeUserData()
	if err != nil {
		return nil, err
	}
	f.log.Debug("decoded user data", "size", len(userData))

	runtimeDataHash := sha256.Sum256(runtimeData)
	f.log.Debug("computed runtime data hash", "hash", hex.EncodeToString(runtimeDataHash[:]))

	return &decodedAdditionalData{
		attestationReport: attestationReport,
		runtimeData:       runtimeData,
		userData:          userData,
		runtimeDataHash:   runtimeDataHash,
	}, nil
}

func (f *AzureTDXFormatter) buildOutput(inputData *AzureTDXInputData, processed *processedAttestationData, decoded *decodedAdditionalData) *AzureTDXProcessedData {
	output := &AzureTDXProcessedData{}

	// Build TPM quote data
	output.AttestationDocument.Attestation.TpmQuote.Quote = HexBytes(processed.sha256Quote.Quote)
	output.AttestationDocument.Attestation.TpmQuote.RsaSignature = HexBytes(processed.decodedSig.RSA.Signature)
	output.AttestationDocument.Attestation.TpmQuote.Pcrs = processed.pcrs

	// Build instance info
	output.AttestationDocument.InstanceInfo.AttestationReport = HexBytes(decoded.attestationReport)
	output.AttestationDocument.InstanceInfo.RuntimeData.Raw = HexBytes(decoded.runtimeData)
	output.AttestationDocument.InstanceInfo.RuntimeData.HclAkPub.ExponentRaw = processed.hclAkPub.RSAParameters.ExponentRaw
	output.AttestationDocument.InstanceInfo.RuntimeData.HclAkPub.ModulusRaw = HexBytes(processed.hclAkPub.RSAParameters.ModulusRaw)

	output.AttestationDocument.UserData = HexBytes(decoded.userData)

	// Build PCR list
	output.Pcrs = f.buildPCRList(processed.pcrs)

	// Set nonce and additional data
	output.Nonce = inputData.Nonce
	output.AdditionalData.RuntimeDataHash = HexBytes32(decoded.runtimeDataHash)

	return output
}

func (f *AzureTDXFormatter) buildPCRList(pcrs [TPMPCRCount]HexBytes32) []PCRValue {
	pcrList := make([]PCRValue, 0, TPMPCRCount)

	for i, pcr := range pcrs {
		pcrList = append(pcrList, PCRValue{
			Index: uint8(i),
			Value: pcr,
		})
	}

	return pcrList
}

// ExtractTrustedParams extracts the trusted parameters from Azure TDX attestation data
// These parameters are used for the setTrustedParams call in the AzureTdxVerifier contract
func (f *AzureTDXFormatter) ExtractTrustedParams(data *AzureTDXProcessedData) (*AzureTDXTrustedParams, error) {
	f.log.Info("extracting trusted params from Azure TDX attestation")

	// Parse the attestation report to extract TDX-specific fields
	attestationReport := data.AttestationDocument.InstanceInfo.AttestationReport
	if len(attestationReport) == 0 {
		return nil, fmt.Errorf("attestation report is empty")
	}

	// Convert the attestation report (TDX quote) to Proto structure to extract measurements
	quoteProto, err := abi.QuoteToProto(attestationReport)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TDX quote: %w", err)
	}

	// Cast to QuoteV4 which contains the actual TDX measurements
	quoteV4, ok := quoteProto.(*tdx.QuoteV4)
	if !ok {
		return nil, fmt.Errorf("expected QuoteV4, got %T", quoteProto)
	}

	// Extract measurements from the TD quote body
	var teeTcbSvn [16]byte
	var mrSeam []byte
	var mrTd []byte

	if quoteV4.TdQuoteBody != nil {
		// TEE TCB SVN is in the TD quote body (16 bytes)
		if len(quoteV4.TdQuoteBody.TeeTcbSvn) == 16 {
			copy(teeTcbSvn[:], quoteV4.TdQuoteBody.TeeTcbSvn)
		}

		// MR_SEAM is in the TD quote body (48 bytes)
		if len(quoteV4.TdQuoteBody.MrSeam) == 48 {
			mrSeam = make([]byte, 48)
			copy(mrSeam, quoteV4.TdQuoteBody.MrSeam)
		}

		// MR_TD is in the TD quote body (48 bytes)
		if len(quoteV4.TdQuoteBody.MrTd) == 48 {
			mrTd = make([]byte, 48)
			copy(mrTd, quoteV4.TdQuoteBody.MrTd)
		}
	} else {
		return nil, fmt.Errorf("TDQuoteBody is nil")
	}

	// Convert PCR values from our processed data to byte arrays
	var pcrs [][]byte
	for _, pcr := range data.Pcrs {
		pcrBytes := make([]byte, 32)
		copy(pcrBytes, pcr.Value[:])
		pcrs = append(pcrs, pcrBytes)
	}

	params := &AzureTDXTrustedParams{
		TeeTcbSvn: teeTcbSvn,
		PcrBitmap: 0xffffff, // All 24 PCRs
		MrSeam:    mrSeam,
		MrTd:      mrTd,
		Pcrs:      pcrs,
	}

	f.log.Info("extracted trusted params successfully",
		"teeTcbSvn", hex.EncodeToString(teeTcbSvn[:]),
		"mrSeam", hex.EncodeToString(mrSeam),
		"mrTd", hex.EncodeToString(mrTd),
		"pcrCount", len(pcrs))
	return params, nil
}
