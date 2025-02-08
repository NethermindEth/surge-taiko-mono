package proposer

import (
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/stretchr/testify/require"
)

func TestCalculateLocally(t *testing.T) {
	// Set up TxCostCalculator with mock values
	calculator := &TxCostCalculator{
		gasPrice:    big.NewInt(100),  // 100 wei
		blobBaseFee: big.NewInt(1000), // 1000 wei
		method:      LocalCalculationMethod,
	}

	// Create test inputs
	txCandidate := &txmgr.TxCandidate{
		Blobs: []*eth.Blob{
			{}, // Empty blob for testing
			{}, // Second empty blob
		},
	}

	txLists := [][]byte{
		[]byte("test data 1"), // 11 bytes
		[]byte("test data 2"), // 11 bytes
	}

	callDataGasUsage := uint64(1000)

	cost, err := calculator.calculateLocally(txCandidate, txLists, callDataGasUsage)

	// blobCost = 2 * 1000 = 2000
	// keccakGas = 72
	// totalSize = 22
	// callDataOverhead = 1000 - 384 - 72 = 544
	// overheadCost = 544 * 100 = 54400
	// totalCost = 2000 + 54400 = 56400

	require.NoError(t, err)
	require.Equal(t, big.NewInt(56400), cost)
}

func TestCalculateKeccakTxListsGasUsage(t *testing.T) {
	txLists := [][]byte{
		[]byte("test data 1"), // 11 bytes
		[]byte("test data 2"), // 11 bytes
	}

	keccakGas, totalSize := keccakTxListsGasUsage(txLists)

	require.Equal(t, uint64(72), keccakGas)
	require.Equal(t, uint64(22), totalSize)
}

func TestCalculateCallDataCostOfTxLists(t *testing.T) {
	cost := callDataCostOfTxLists(uint64(22))

	require.Equal(t, uint64(384), cost)
}
