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
		[]byte("test\x00data\x001"), // 11 bytes with null bytes
		[]byte("test\x00data\x002"), // 11 bytes with null bytes
	}

	callDataGasUsage := uint64(1000)

	cost, err := calculator.calculateLocally(txCandidate, txLists, callDataGasUsage)

	// blobCost = 2 * 1000 * 131072 = 262144000
	// keccakGas = 36
	// totalSize = 22
	// callDataOverhead = 1000 - 304 - 36 = 660
	// overheadCost = 660 * 100 = 66000
	// totalCost = 262144000 + 66000 = 262210000

	require.NoError(t, err)
	require.Equal(t, big.NewInt(262210000), cost)
}

func TestCalculateKeccakTxListsGasUsage(t *testing.T) {
	txLists := [][]byte{
		[]byte("test\x00data\x001"), // 11 bytes with null bytes
		[]byte("test\x00data\x002"), // 11 bytes with null bytes
	}

	keccakGas := keccakTxListsGasUsage(txLists)

	require.Equal(t, uint64(36), keccakGas)
}

func TestCalculateCallDataCostOfTxLists(t *testing.T) {
	txLists := [][]byte{
		[]byte("test\x00data\x001"), // 11 bytes with null bytes
		[]byte("test\x00data\x002"), // 11 bytes with null bytes
	}

	cost := callDataCostOfTxLists(txLists)

	require.Equal(t, uint64(304), cost)
}
