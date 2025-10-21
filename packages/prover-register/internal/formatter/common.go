package formatter

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-tpm-tools/proto/attest"
	tpmproto "github.com/google/go-tpm-tools/proto/tpm"
	"github.com/google/go-tpm/legacy/tpm2"
)

const (
	TPMPCRCount = 24
	HexPrefix   = "0x"
)

// HexBytes is a byte slice that marshals/unmarshals as a hex string
type HexBytes []byte

func (h HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(HexPrefix + hex.EncodeToString(h))
}

func (h *HexBytes) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return fmt.Errorf("unmarshal hex string: %w", err)
	}

	hexStr = strings.TrimPrefix(hexStr, HexPrefix)

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("decode hex string: %w", err)
	}

	*h = decoded
	return nil
}

// HexBytes32 is a 32-byte array that marshals as a hex string
type HexBytes32 [32]byte

func (h HexBytes32) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s%064x", HexPrefix, h))
}

// InstanceInfo contains attestation report and runtime data
type InstanceInfo struct {
	AttestationReport string `json:"attestationReport"`
	RuntimeData       string `json:"runtimeData"`
}

func (i *InstanceInfo) DecodeAttestationReport() ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(i.AttestationReport)
	if err != nil {
		return nil, fmt.Errorf("decode attestation report: %w", err)
	}
	return data, nil
}

func (i *InstanceInfo) DecodeRuntimeData() ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(i.RuntimeData)
	if err != nil {
		return nil, fmt.Errorf("decode runtime data: %w", err)
	}
	return data, nil
}

// AttestationDocument contains the attestation data from Azure
type AttestationDocument struct {
	Attestation  *attest.Attestation
	InstanceInfo []byte
	UserData     string
}

func (a *AttestationDocument) DecodeUserData() ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(a.UserData)
	if err != nil {
		return nil, fmt.Errorf("decode user data: %w", err)
	}
	return data, nil
}

// PCRValue represents a single PCR value with its index
type PCRValue struct {
	Index uint8      `json:"index"`
	Value HexBytes32 `json:"value"`
}

// TPMQuoteData contains TPM quote information
type TPMQuoteData struct {
	Quote        HexBytes                `json:"quote"`
	RsaSignature HexBytes                `json:"rsaSignature"`
	Pcrs         [TPMPCRCount]HexBytes32 `json:"pcrs"`
}

// RuntimeDataInfo contains runtime data and HCL AK public key
type RuntimeDataInfo struct {
	Raw      HexBytes `json:"raw"`
	HclAkPub struct {
		ExponentRaw uint32   `json:"exponentRaw"`
		ModulusRaw  HexBytes `json:"modulusRaw"`
	} `json:"hclAkPub"`
}

// processedAttestationData holds processed attestation information
type processedAttestationData struct {
	hclAkPub    *tpm2.Public
	sha256Quote *tpmproto.Quote
	pcrs        [TPMPCRCount]HexBytes32
	decodedSig  *tpm2.Signature
}

// decodedAdditionalData holds decoded additional data
type decodedAdditionalData struct {
	attestationReport []byte
	runtimeData       []byte
	userData          []byte
	runtimeDataHash   [32]byte
}

// TDXTrustedParams represents the trusted parameters for raw TDX verification
// This matches the TrustedParams struct in TdxVerifier.sol
type TDXTrustedParams struct {
	TeeTcbSvn [16]byte // TEE TCB SVN
	MrSeam    []byte   // MR_SEAM measurement
	MrTd      []byte   // MR_TD measurement
	RtMr0     []byte   // Runtime Measurement Register 0
	RtMr1     []byte   // Runtime Measurement Register 1
	RtMr2     []byte   // Runtime Measurement Register 2
	RtMr3     []byte   // Runtime Measurement Register 3
}

// AzureTDXTrustedParams represents the trusted parameters for Azure TDX verification
// This matches the TrustedParams struct in AzureTdxVerifier.sol
type AzureTDXTrustedParams struct {
	TeeTcbSvn [16]byte // TEE TCB SVN
	PcrBitmap uint32   // Bitmap of PCRs to verify (usually 0xffffff for all 24)
	MrSeam    []byte   // MR_SEAM measurement
	MrTd      []byte   // MR_TD measurement
	Pcrs      [][]byte // Array of expected PCR values based on bitmap
}
