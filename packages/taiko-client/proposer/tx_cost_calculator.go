package proposer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// BlobCostCalculationMethod represents the method used for transaction cost calculation
type BlobCostCalculationMethod int

const (
	// EthEstimateMethod uses the standard way to calculate transaction costs
	EthEstimateMethod BlobCostCalculationMethod = iota
	// LocalCalculationMethod uses an optimized way to calculate transaction costs
	LocalCalculationMethod
)

// TxCostCalculator handles the calculation of transaction costs
type TxCostCalculator struct {
	l1Client        *rpc.EthClient
	proposerAddress common.Address
	ctx             context.Context
	method          BlobCostCalculationMethod
	gasPrice        *big.Int
	blobBaseFee     *big.Int
}

// NewTxCostCalculator creates a new TxCostCalculator instance
func NewTxCostCalculator(
	ctx context.Context,
	l1Client *rpc.EthClient,
	proposerAddress common.Address,
	method BlobCostCalculationMethod,
) (*TxCostCalculator, error) {
	log.Debug("NewTxCostCalculator", "method", method)
	gasPrice, err := l1Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewTxCostCalculator: failed to get gas price: %w", err)
	}
	blobBaseFee, err := l1Client.BlobBaseFee(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewTxCostCalculator: failed to get blob base fee: %w", err)
	}

	return &TxCostCalculator{
		l1Client:        l1Client,
		proposerAddress: proposerAddress,
		ctx:             ctx,
		method:          method,
		gasPrice:        gasPrice,
		blobBaseFee:     blobBaseFee,
	}, nil
}

// CalculateTxCost calculates the total cost of a transaction
func (c *TxCostCalculator) CalculateTxCost(
	txCandidate *txmgr.TxCandidate,
	txLists [][]byte,
	callDataGasUsage uint64,
) (*big.Int, error) {
	if c.method == LocalCalculationMethod {
		return c.calculateLocally(txCandidate, txLists, callDataGasUsage)
	}
	return c.calculateWithEthEstimate(txCandidate)
}

// calculateWithEthEstimate calculates transaction cost using eth_estimateGas
func (c *TxCostCalculator) calculateWithEthEstimate(txCandidate *txmgr.TxCandidate) (*big.Int, error) {
	log.Debug("get tx cost with eth_estimateGas")
	blobTxCost, _, err := c.GetBlobTransactionCost(txCandidate)
	if err != nil {
		return nil, err
	}

	blobCost, err := getBlobCost(txCandidate.Blobs, c.blobBaseFee)
	if err != nil {
		return nil, err
	}

	return new(big.Int).Add(blobTxCost, blobCost), nil
}

// calculateLocally calculates transaction cost using the optimized method
func (c *TxCostCalculator) calculateLocally(
	txCandidate *txmgr.TxCandidate,
	txLists [][]byte,
	callDataGasUsage uint64,
) (*big.Int, error) {
	log.Debug("get tx cost locally")
	blobCost, err := getBlobCost(txCandidate.Blobs, c.blobBaseFee)
	if err != nil {
		return nil, err
	}

	keccakGas := keccakTxListsGasUsage(txLists)
	callDataOverhead := callDataGasUsage - callDataCostOfTxLists(txLists) - keccakGas

	totalBlobCost := blobCost.Add(blobCost, new(big.Int).Mul(new(big.Int).SetUint64(callDataOverhead), c.gasPrice))
	return totalBlobCost, nil
}

func keccakTxListsGasUsage(txLists [][]byte) uint64 {
	gasUsage := uint64(0)
	totalSize := uint64(0)
	for _, txList := range txLists {
		totalSize += uint64(len(txList))
	}
	minimumWordSize := (totalSize + uint64(31)) / uint64(32)
	staticGas := uint64(30)
	dynamicGas := uint64(6) * minimumWordSize
	gasUsage += staticGas + dynamicGas

	return gasUsage
}

func callDataCostOfTxLists(txLists [][]byte) uint64 {
	var cost uint64
	for _, txList := range txLists {
		for _, b := range txList {
			if b == 0 {
				cost += 4
			} else {
				cost += 16
			}
		}
	}
	return cost
}

// GetTransactionCost calculates the cost of a transaction
func (c *TxCostCalculator) GetCallDataTransactionCost(
	txCandidate *txmgr.TxCandidate,
) (*big.Int, uint64, error) {
	log.Debug("GetCallDataTransactionCost")

	msg := ethereum.CallMsg{
		From:  c.proposerAddress,
		To:    txCandidate.To,
		Data:  txCandidate.TxData,
		Gas:   0,
		Value: nil,
	}

	estimatedGasUsage, err := c.l1Client.EstimateGas(c.ctx, msg)
	if err != nil {
		log.Info("GetCallDataTransactionCost: estimate gas ethereum.CallMsg", "from", msg.From,
			"to", msg.To, "Gas", msg.Gas, "Value", msg.Value)
		return nil, 0, fmt.Errorf("GetTransactionCost: failed to estimate gas: %w", err)
	}

	log.Debug("GetCallDataTransactionCost", "estimatedGasUsage", estimatedGasUsage)

	return new(big.Int).Mul(c.gasPrice, new(big.Int).SetUint64(estimatedGasUsage)), estimatedGasUsage, nil
}

// GetTransactionCost calculates the cost of a transaction
func (c *TxCostCalculator) GetBlobTransactionCost(
	txCandidate *txmgr.TxCandidate,
) (*big.Int, uint64, error) {
	hexData := hex.EncodeToString(txCandidate.TxData)
	log.Debug("GetBlobTransactionCost", "blobBaseFee", c.blobBaseFee, "txCandidate.TxData hex", hexData)

	blobHashes, err := calculateBlobHashes(txCandidate.Blobs)
	if err != nil {
		return nil, 0, fmt.Errorf("GetBlobTransactionCost: failed to calculate blob hashes: %w", err)
	}
	msg := ethereum.CallMsg{
		From:          c.proposerAddress,
		To:            txCandidate.To,
		Data:          txCandidate.TxData,
		Gas:           0,
		Value:         nil,
		BlobGasFeeCap: c.blobBaseFee,
		BlobHashes:    blobHashes,
	}

	estimatedGasUsage, err := c.l1Client.EstimateGas(c.ctx, msg)
	if err != nil {
		log.Info("GetTransactionCost: estimate gas ethereum.CallMsg", "from", msg.From,
			"to", msg.To, "Gas", msg.Gas, "Value", msg.Value, "BlobGasFeeCap", msg.BlobGasFeeCap, "BlobHashes", msg.BlobHashes)
		return nil, 0, fmt.Errorf("GetTransactionCost: failed to estimate gas: %w", err)
	}

	log.Debug("GetBlobTransactionCost", "estimatedGasUsage", estimatedGasUsage)

	return new(big.Int).Mul(c.gasPrice, new(big.Int).SetUint64(estimatedGasUsage)), estimatedGasUsage, nil
}

func calculateBlobHashes(blobs []*eth.Blob) ([]common.Hash, error) {
	log.Debug("Calculating blob hashes")

	var blobHashes []common.Hash
	for _, blob := range blobs {
		commitment, err := blob.ComputeKZGCommitment()
		if err != nil {
			return nil, err
		}
		blobHash := kzg4844.CalcBlobHashV1(sha256.New(), &commitment)
		blobHashes = append(blobHashes, blobHash)
	}
	return blobHashes, nil
}

func getBlobCost(blobs []*eth.Blob, blobBaseFee *big.Int) (*big.Int, error) {
	// Each blob costs exactly 131072 (0x20000) blob gas
	blobGasPerBlob := uint64(131072)
	totalBlobGas := blobGasPerBlob * uint64(len(blobs))

	// Total cost is blob gas * blob base fee
	return new(big.Int).Mul(
		new(big.Int).SetUint64(totalBlobGas),
		blobBaseFee,
	), nil
}
