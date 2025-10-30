package contract

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/taikoxyz/taiko-mono/packages/prover-register/internal/formatter"
	"go.uber.org/zap"
)

type Client struct {
	ethClient       *ethclient.Client
	verifierAddress common.Address
	privateKey      *ecdsa.PrivateKey
	log             *zap.SugaredLogger
	chainID         *big.Int
	address         common.Address
	dryRun          bool
	dryRunAsOwner   bool
}

func NewClient(ethClient *ethclient.Client, verifierAddress common.Address, privateKey *ecdsa.PrivateKey, log *zap.SugaredLogger, dryRun bool, dryRunAsOwner bool) (*Client, error) {
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	var address common.Address
	if privateKey != nil {
		address = crypto.PubkeyToAddress(privateKey.PublicKey)
	}

	return &Client{
		ethClient:       ethClient,
		verifierAddress: verifierAddress,
		dryRun:          dryRun,
		dryRunAsOwner:   dryRunAsOwner,
		privateKey:      privateKey,
		log:             log,
		chainID:         chainID,
		address:         address,
	}, nil
}

func (c *Client) TrustCollateral(ctx context.Context, verifierType string, data interface{}) error {
	switch strings.ToLower(verifierType) {
	case "tdx":
		return c.trustTDXCollateral(ctx, data.(*formatter.TDXProcessedData))
	case "azure-tdx":
		return c.trustAzureTDXCollateral(ctx, data.(*formatter.AzureTDXProcessedData))
	case "sgx":
		return c.trustSGXCollateral(ctx, data.(*formatter.SGXProcessedData))
	case "sp1":
		return c.trustSP1Collateral(ctx, data.(*formatter.SP1ProcessedData))
	case "risc0":
		return c.trustRISC0Collateral(ctx, data.(*formatter.RISC0ProcessedData))
	default:
		return fmt.Errorf("unsupported verifier type: %s", verifierType)
	}
}

func (c *Client) RegisterInstance(ctx context.Context, verifierType string, data interface{}) error {
	switch strings.ToLower(verifierType) {
	case "tdx":
		return c.registerTDXInstance(ctx, data.(*formatter.TDXProcessedData))
	case "azure-tdx":
		return c.registerAzureTDXInstance(ctx, data.(*formatter.AzureTDXProcessedData))
	case "sgx":
		return c.registerSGXInstance(ctx, data.(*formatter.SGXProcessedData))
	case "sp1":
		return c.registerSP1Instance(ctx, data.(*formatter.SP1ProcessedData))
	case "risc0":
		return c.registerRISC0Instance(ctx, data.(*formatter.RISC0ProcessedData))
	default:
		return fmt.Errorf("unsupported verifier type: %s", verifierType)
	}
}

// TDX specific methods
func (c *Client) trustTDXCollateral(ctx context.Context, data *formatter.TDXProcessedData) error {
	if c.dryRun {
		c.log.Info("[DRY RUN] Would trust raw TDX collateral")
	} else {
		c.log.Info("trusting raw TDX collateral")
	}

	// Parse TdxVerifier ABI for setTrustedParams function
	// This matches the TrustedParams struct in TdxVerifier.sol
	abiJSON := `[{
		"name": "setTrustedParams",
		"type": "function",
		"inputs": [
			{"name": "_index", "type": "uint256"},
			{"name": "_params", "type": "tuple", "components": [
				{"name": "teeTcbSvn", "type": "bytes16"},
				{"name": "mrSeam", "type": "bytes"},
				{"name": "mrTd", "type": "bytes"},
				{"name": "rtMr0", "type": "bytes"},
				{"name": "rtMr1", "type": "bytes"},
				{"name": "rtMr2", "type": "bytes"},
				{"name": "rtMr3", "type": "bytes"}
			]}
		]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Extract trusted parameters from the processed data
	tdxFormatter := formatter.NewTDXFormatter(c.log)
	trustedParamsData, err := tdxFormatter.ExtractTrustedParams(data)
	if err != nil {
		return fmt.Errorf("failed to extract trusted params: %w", err)
	}

	// For raw TDX, match the TrustedParams structure in TdxVerifier.sol
	trustedParams := struct {
		TeeTcbSvn [16]byte
		MrSeam    []byte
		MrTd      []byte
		RtMr0     []byte
		RtMr1     []byte
		RtMr2     []byte
		RtMr3     []byte
	}{
		TeeTcbSvn: trustedParamsData.TeeTcbSvn,
		MrSeam:    trustedParamsData.MrSeam,
		MrTd:      trustedParamsData.MrTd,
		RtMr0:     trustedParamsData.RtMr0,
		RtMr1:     trustedParamsData.RtMr1,
		RtMr2:     trustedParamsData.RtMr2,
		RtMr3:     trustedParamsData.RtMr3,
	}

	// Pack the function call data
	callData, err := parsedABI.Pack("setTrustedParams", big.NewInt(0), trustedParams)
	if err != nil {
		return fmt.Errorf("failed to pack call data: %w", err)
	}

	// Send transaction
	txHash, err := c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	if c.dryRun {
		c.log.Info("[DRY RUN] Raw TDX collateral would be trusted")
	} else {
		c.log.Infow("raw TDX collateral trusted", "tx", txHash.Hex())
	}
	return nil
}

// Azure TDX specific methods
func (c *Client) trustAzureTDXCollateral(ctx context.Context, data *formatter.AzureTDXProcessedData) error {
	if c.dryRun {
		c.log.Info("[DRY RUN] Would trust Azure TDX collateral")
	} else {
		c.log.Info("trusting Azure TDX collateral")
	}

	// Parse AzureTdxVerifier ABI for setTrustedParams function
	// This matches the TrustedParams struct in AzureTdxVerifier.sol
	abiJSON := `[{
		"name": "setTrustedParams",
		"type": "function",
		"inputs": [
			{"name": "index", "type": "uint256"},
			{"name": "_params", "type": "tuple", "components": [
				{"name": "teeTcbSvn", "type": "bytes16"},
				{"name": "pcrBitmap", "type": "uint24"},
				{"name": "mrSeam", "type": "bytes"},
				{"name": "mrTd", "type": "bytes"},
				{"name": "pcrs", "type": "bytes32[]"}
			]}
		]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Extract trusted parameters from the processed data
	azureFormatter := formatter.NewAzureTDXFormatter(c.log)
	trustedParamsData, err := azureFormatter.ExtractTrustedParams(data)
	if err != nil {
		return fmt.Errorf("failed to extract trusted params: %w", err)
	}

	// Convert PCR bytes to bytes32 array for Solidity
	var pcrBytes32 [][32]byte
	for _, pcr := range trustedParamsData.Pcrs {
		var b32 [32]byte
		copy(b32[:], pcr)
		pcrBytes32 = append(pcrBytes32, b32)
	}

	// For Azure TDX, use the contract-specific structure
	trustedParams := struct {
		TeeTcbSvn [16]byte
		PcrBitmap *big.Int // uint24 in contract
		MrSeam    []byte
		MrTd      []byte
		Pcrs      [][32]byte
	}{
		TeeTcbSvn: trustedParamsData.TeeTcbSvn,
		PcrBitmap: big.NewInt(int64(trustedParamsData.PcrBitmap)),
		MrSeam:    trustedParamsData.MrSeam,
		MrTd:      trustedParamsData.MrTd,
		Pcrs:      pcrBytes32,
	}

	// Pack the function call data
	callData, err := parsedABI.Pack("setTrustedParams", big.NewInt(0), trustedParams)
	if err != nil {
		return fmt.Errorf("failed to pack call data: %w", err)
	}

	// Send transaction
	txHash, err := c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	if c.dryRun {
		c.log.Info("[DRY RUN] Azure TDX collateral would be trusted")
	} else {
		c.log.Infow("Azure TDX collateral trusted", "tx", txHash.Hex())
	}
	return nil
}

func (c *Client) registerTDXInstance(ctx context.Context, data *formatter.TDXProcessedData) error {
	if c.dryRun {
		c.log.Info("[DRY RUN] Would register raw TDX instance")
	} else {
		c.log.Info("registering raw TDX instance")
	}

	// Parse TDX verifier ABI for registerInstance function
	abiJSON := `[{
		"name": "registerInstance",
		"type": "function",
		"inputs": [
			{"name": "_params", "type": "tuple", "components": [
				{"name": "quote", "type": "bytes"},
				{"name": "userData", "type": "bytes"},
				{"name": "nonce", "type": "bytes"}
			]}
		],
		"outputs": [{"name": "", "type": "uint256"}]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// For raw TDX, use the raw quote directly
	verifyParams := struct {
		Quote    []byte
		UserData []byte
		Nonce    []byte
	}{
		Quote:    data.Quote,
		UserData: []byte(data.PublicKey), // Use public key as user data for raw TDX
		Nonce:    data.Nonce,
	}

	// Pack the function call data
	callData, err := parsedABI.Pack("registerInstance", verifyParams)
	if err != nil {
		return fmt.Errorf("failed to pack call data: %w", err)
	}

	// Send transaction
	txHash, err := c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	if c.dryRun {
		c.log.Info("[DRY RUN] Raw TDX instance would be registered")
	} else {
		c.log.Infow("raw TDX instance registered", "tx", txHash.Hex())
	}
	return nil
}

func (c *Client) registerAzureTDXInstance(ctx context.Context, data *formatter.AzureTDXProcessedData) error {
	if c.dryRun {
		c.log.Info("[DRY RUN] Would register Azure TDX instance")
	} else {
		c.log.Info("registering Azure TDX instance")
	}

	// Parse AzureTdxVerifier ABI for registerInstance function
	// This takes a trustedParamsIdx and an AzureTDX.VerifyParams struct
	abiJSON := `[{
		"name": "registerInstance",
		"type": "function",
		"inputs": [
			{"name": "_trustedParamsIdx", "type": "uint256"},
			{"name": "_attestation", "type": "tuple", "components": [
				{"name": "attestationDocument", "type": "tuple", "components": [
					{"name": "attestation", "type": "tuple", "components": [
						{"name": "tpmQuote", "type": "tuple", "components": [
							{"name": "quote", "type": "bytes"},
							{"name": "rsaSignature", "type": "bytes"},
							{"name": "pcrs", "type": "bytes32[24]"}
						]}
					]},
					{"name": "instanceInfo", "type": "tuple", "components": [
						{"name": "attestationReport", "type": "bytes"},
						{"name": "runtimeData", "type": "tuple", "components": [
							{"name": "raw", "type": "bytes"},
							{"name": "hclAkPub", "type": "tuple", "components": [
								{"name": "exponentRaw", "type": "uint24"},
								{"name": "modulusRaw", "type": "bytes"}
							]}
						]}
					]},
					{"name": "userData", "type": "bytes"}
				]},
				{"name": "pcrs", "type": "tuple[]", "components": [
					{"name": "index", "type": "uint256"},
					{"name": "digest", "type": "bytes32"}
				]},
				{"name": "nonce", "type": "bytes"}
			]}
		],
		"outputs": [{"name": "", "type": "uint256"}]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Build the PCR array for the contract
	type PCREntry struct {
		Index  *big.Int
		Digest [32]byte
	}
	var pcrEntries []PCREntry
	for _, pcr := range data.Pcrs {
		pcrEntries = append(pcrEntries, PCREntry{
			Index:  big.NewInt(int64(pcr.Index)),
			Digest: pcr.Value,
		})
	}

	// Build the complete VerifyParams structure matching AzureTDX.sol
	verifyParams := struct {
		AttestationDocument struct {
			Attestation struct {
				TpmQuote struct {
					Quote        []byte
					RsaSignature []byte
					Pcrs         [24][32]byte
				}
			}
			InstanceInfo struct {
				AttestationReport []byte
				RuntimeData       struct {
					Raw      []byte
					HclAkPub struct {
						ExponentRaw *big.Int
						ModulusRaw  []byte
					}
				}
			}
			UserData []byte
		}
		Pcrs  []PCREntry
		Nonce []byte
	}{}

	// Fill in the structure from our processed data
	verifyParams.AttestationDocument.Attestation.TpmQuote.Quote = data.AttestationDocument.Attestation.TpmQuote.Quote
	verifyParams.AttestationDocument.Attestation.TpmQuote.RsaSignature = data.AttestationDocument.Attestation.TpmQuote.RsaSignature
	// Convert PCRs from HexBytes32 to [32]byte
	for i, pcr := range data.AttestationDocument.Attestation.TpmQuote.Pcrs {
		verifyParams.AttestationDocument.Attestation.TpmQuote.Pcrs[i] = [32]byte(pcr)
	}

	verifyParams.AttestationDocument.InstanceInfo.AttestationReport = data.AttestationDocument.InstanceInfo.AttestationReport
	verifyParams.AttestationDocument.InstanceInfo.RuntimeData.Raw = data.AttestationDocument.InstanceInfo.RuntimeData.Raw
	verifyParams.AttestationDocument.InstanceInfo.RuntimeData.HclAkPub.ExponentRaw = big.NewInt(int64(data.AttestationDocument.InstanceInfo.RuntimeData.HclAkPub.ExponentRaw))
	verifyParams.AttestationDocument.InstanceInfo.RuntimeData.HclAkPub.ModulusRaw = data.AttestationDocument.InstanceInfo.RuntimeData.HclAkPub.ModulusRaw

	verifyParams.AttestationDocument.UserData = data.AttestationDocument.UserData
	verifyParams.Pcrs = pcrEntries
	verifyParams.Nonce = data.Nonce

	// Use trusted params index 0 by default (can be made configurable)
	trustedParamsIdx := big.NewInt(0)

	// Pack the function call data
	callData, err := parsedABI.Pack("registerInstance", trustedParamsIdx, verifyParams)
	if err != nil {
		return fmt.Errorf("failed to pack call data: %w", err)
	}

	// Send transaction
	txHash, err := c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	if c.dryRun {
		c.log.Info("[DRY RUN] Azure TDX instance would be registered")
	} else {
		c.log.Infow("Azure TDX instance registered", "tx", txHash.Hex())
	}
	return nil
}

// SGX specific methods
func (c *Client) trustSGXCollateral(ctx context.Context, data *formatter.SGXProcessedData) error {
	c.log.Infow("SGX doesn't require explicit collateral trust", "mr_enclave", data.MrEnclave)
	// SGX verifier doesn't have a separate trust collateral step
	// The trust is established through the attestation process
	return nil
}

func (c *Client) registerSGXInstance(ctx context.Context, data *formatter.SGXProcessedData) error {
	c.log.Info("registering SGX instance")

	// For SGX, we need to call the registerInstance function with the quote
	// The quote contains the attestation data that will be verified
	abiJSON := `[{
		"name": "registerInstance",
		"type": "function",
		"inputs": [
			{"name": "_attestation", "type": "bytes"}
		],
		"outputs": [{"name": "", "type": "uint256"}]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Decode quote hex string
	quoteBytes, err := hex.DecodeString(strings.TrimPrefix(data.Quote, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode quote hex: %w", err)
	}

	// Pack the function call data
	callData, err := parsedABI.Pack("registerInstance", quoteBytes)
	if err != nil {
		return fmt.Errorf("failed to pack call data: %w", err)
	}

	// Send transaction
	txHash, err := c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	c.log.Infow("SGX instance registered", "tx", txHash.Hex())
	return nil
}

// SP1 specific methods
func (c *Client) trustSP1Collateral(ctx context.Context, data *formatter.SP1ProcessedData) error {
	c.log.Info("trusting SP1 collateral")

	// Parse SP1 verifier ABI for setProgramTrusted function
	abiJSON := `[{
		"name": "setProgramTrusted",
		"type": "function",
		"inputs": [
			{"name": "_programVKey", "type": "bytes32"},
			{"name": "_trusted", "type": "bool"}
		]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Convert program hashes to bytes32
	aggregationHash := common.HexToHash(data.AggregationProgramHash)
	blockHash := common.HexToHash(data.BlockProgramHash)

	// Trust aggregation program
	callData, err := parsedABI.Pack("setProgramTrusted", aggregationHash, true)
	if err != nil {
		return fmt.Errorf("failed to pack aggregation call data: %w", err)
	}

	txHash, err := c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to trust aggregation program: %w", err)
	}
	c.log.Infow("SP1 aggregation program trusted", "tx", txHash.Hex(), "hash", data.AggregationProgramHash)

	// Trust block program
	callData, err = parsedABI.Pack("setProgramTrusted", blockHash, true)
	if err != nil {
		return fmt.Errorf("failed to pack block call data: %w", err)
	}

	txHash, err = c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to trust block program: %w", err)
	}
	c.log.Infow("SP1 block program trusted", "tx", txHash.Hex(), "hash", data.BlockProgramHash)

	return nil
}

func (c *Client) registerSP1Instance(ctx context.Context, data *formatter.SP1ProcessedData) error {
	c.log.Info("SP1 doesn't require instance registration")
	// SP1 verifier doesn't have instance registration
	// It only needs program verification keys to be trusted
	return nil
}

// RISC0 specific methods
func (c *Client) trustRISC0Collateral(ctx context.Context, data *formatter.RISC0ProcessedData) error {
	if c.dryRun {
		c.log.Info("[DRY RUN] Would trust RISC0 collateral")
	} else {
		c.log.Info("trusting RISC0 collateral")
	}

	// Similar to SP1, RISC0 needs to trust program verification keys
	abiJSON := `[{
		"name": "setImageIdTrusted",
		"type": "function",
		"inputs": [
			{"name": "_imageId", "type": "bytes32"},
			{"name": "_trusted", "type": "bool"}
		]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Convert program hashes to bytes32
	aggregationHash := common.HexToHash(data.AggregationProgramHash)
	blockHash := common.HexToHash(data.BlockProgramHash)

	// Trust aggregation image
	callData, err := parsedABI.Pack("setImageIdTrusted", aggregationHash, true)
	if err != nil {
		return fmt.Errorf("failed to pack aggregation call data: %w", err)
	}

	txHash, err := c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to trust aggregation image: %w", err)
	}
	c.log.Infow("RISC0 aggregation image trusted", "tx", txHash.Hex(), "hash", data.AggregationProgramHash)

	// Trust block image
	callData, err = parsedABI.Pack("setImageIdTrusted", blockHash, true)
	if err != nil {
		return fmt.Errorf("failed to pack block call data: %w", err)
	}

	txHash, err = c.sendTransaction(ctx, callData)
	if err != nil {
		return fmt.Errorf("failed to trust block image: %w", err)
	}
	c.log.Infow("RISC0 block image trusted", "tx", txHash.Hex(), "hash", data.BlockProgramHash)

	return nil
}

func (c *Client) registerRISC0Instance(ctx context.Context, data *formatter.RISC0ProcessedData) error {
	c.log.Info("RISC0 doesn't require instance registration")
	// RISC0 verifier doesn't have instance registration
	// It only needs image IDs to be trusted
	return nil
}

// Helper method to send transactions
func (c *Client) sendTransaction(ctx context.Context, callData []byte) (common.Hash, error) {
	// If dry run, just simulate the transaction
	if c.dryRun {
		return c.simulateTransaction(ctx, callData)
	}

	nonce, err := c.ethClient.PendingNonceAt(ctx, c.address)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Estimate gas
	msg := ethereum.CallMsg{
		From: c.address,
		To:   &c.verifierAddress,
		Data: callData,
	}

	gasLimit, err := c.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Create transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &c.verifierAddress,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     callData,
	})

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.chainID), c.privateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = c.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	c.log.Infow("transaction sent", "hash", signedTx.Hash().Hex())

	// Wait for receipt
	receipt, err := c.waitForReceipt(ctx, signedTx.Hash())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get receipt: %w", err)
	}

	if receipt.Status == 0 {
		return common.Hash{}, fmt.Errorf("transaction failed")
	}

	return signedTx.Hash(), nil
}

// simulateTransaction simulates a transaction using eth_call
func (c *Client) simulateTransaction(ctx context.Context, callData []byte) (common.Hash, error) {
	c.log.Info("[DRY RUN] Simulating transaction as contract owner")

	var msg ethereum.CallMsg
	if c.dryRunAsOwner {
		ownerAddress, err := c.getContractOwner(ctx)
		if err != nil {
			c.log.Warnw("[DRY RUN] Could not fetch contract owner, using zero address", "error", err)
			ownerAddress = common.Address{}
		} else {
			c.log.Infow("[DRY RUN] Simulating as owner", "owner", ownerAddress.Hex())
		}
		msg = ethereum.CallMsg{
			From: ownerAddress,
			To:   &c.verifierAddress,
			Data: callData,
		}
	} else {
		msg = ethereum.CallMsg{
			From: c.address,
			To:   &c.verifierAddress,
			Data: callData,
		}
	}

	// Try to estimate gas (this also validates the transaction)
	gasLimit, err := c.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		// If gas estimation fails, try the call anyway to get more details
		c.log.Warnw("[DRY RUN] Gas estimation failed, attempting call anyway", "error", err)
		gasLimit = uint64(0)
	}

	// Perform the call to check if it would succeed
	result, err := c.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		c.log.Errorw("[DRY RUN] Call failed", "error", err, "from", msg.From.Hex())
		return common.Hash{}, fmt.Errorf("[DRY RUN] call failed (simulating as %s): %w", msg.From.Hex(), err)
	}

	// Log simulation results
	c.log.Infow("[DRY RUN] Transaction simulation successful",
		"simulatedFrom", msg.From.Hex(),
		"estimatedGas", gasLimit,
		"callDataSize", len(callData),
		"resultSize", len(result))

	// Return a fake hash for dry run
	fakeHash := common.BytesToHash([]byte("dry-run-simulation"))
	return fakeHash, nil
}

// getContractOwner attempts to fetch the owner of the contract
func (c *Client) getContractOwner(ctx context.Context) (common.Address, error) {
	// Try to call owner() function (common pattern for Ownable contracts)
	ownerABI := `[{"name":"owner","type":"function","inputs":[],"outputs":[{"name":"","type":"address"}]}]`

	parsedABI, err := abi.JSON(strings.NewReader(ownerABI))
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to parse owner ABI: %w", err)
	}

	callData, err := parsedABI.Pack("owner")
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to pack owner call: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.verifierAddress,
		Data: callData,
	}

	result, err := c.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to call owner function: %w", err)
	}

	if len(result) < 32 {
		return common.Address{}, fmt.Errorf("invalid owner response")
	}

	// Extract address from the result (last 20 bytes of the 32-byte response)
	var owner common.Address
	copy(owner[:], result[12:32])

	return owner, nil
}

func (c *Client) waitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		// Wait a bit before retrying
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue
		}
	}
}
