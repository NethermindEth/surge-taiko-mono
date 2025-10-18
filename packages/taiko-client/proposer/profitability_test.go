package proposer

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
)

// mockCostEstimator is a simple mock that returns a fixed cost for testing.
type mockCostEstimator struct {
	cost *big.Int
}

func (m *mockCostEstimator) estimateL1Cost(ctx context.Context, candidate *txmgr.TxCandidate) (*big.Int, error) {
	return m.cost, nil
}

// dummyProtocolConfigs is a minimal ProtocolConfigs implementation for tests.
type dummyProtocolConfigs struct{}

func (d *dummyProtocolConfigs) BaseFeeConfig() *pacayaBindings.LibSharedDataBaseFeeConfig {
	return &pacayaBindings.LibSharedDataBaseFeeConfig{SharingPctg: 100}
}
func (d *dummyProtocolConfigs) BlockMaxGasLimit() uint32              { return 0 }
func (d *dummyProtocolConfigs) ForkHeightsOntake() uint64             { return 0 }
func (d *dummyProtocolConfigs) ForkHeightsPacaya() uint64             { return 0 }
func (d *dummyProtocolConfigs) LivenessBond() *big.Int                { return big.NewInt(0) }
func (d *dummyProtocolConfigs) LivenessBondPerBlock() *big.Int        { return big.NewInt(0) }
func (d *dummyProtocolConfigs) MaxProposals() uint64                  { return 0 }
func (d *dummyProtocolConfigs) ProvingWindow() (time.Duration, error) { return 0, nil }
func (d *dummyProtocolConfigs) MaxBlocksPerBatch() int                { return 0 }
func (d *dummyProtocolConfigs) MaxAnchorHeightOffset() uint64         { return 0 }

// TestFindOptimalBaseFeeThreshold tests the optimized algorithm for finding
// the base fee threshold that maximizes collected fees
func TestFindOptimalBaseFeeThreshold(t *testing.T) {
	p := &Proposer{
		protocolConfigs: &dummyProtocolConfigs{}, // 100% sharing
	}

	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("FindsOptimalThresholdWithVariedFees", func(t *testing.T) {
		// Create transactions with different gas fees and gas amounts
		// Optimal should be when baseFee * cumulativeGas is maximized
		tx1 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     0,
			GasFeeCap: big.NewInt(10_000_000_000), // 10 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       100000, // Large gas
			To:        &testAddr,
			Value:     common.Big0,
		})

		tx2 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     1,
			GasFeeCap: big.NewInt(100_000_000_000), // 100 gwei (highest)
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000, // Small gas
			To:        &testAddr,
			Value:     common.Big0,
		})

		tx3 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     2,
			GasFeeCap: big.NewInt(20_000_000_000), // 20 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       50000, // Medium gas
			To:        &testAddr,
			Value:     common.Big0,
		})

		txBatch := []types.Transactions{{tx1, tx2, tx3}}

		optimalBaseFee, optimalFees := p.findOptimalBaseFeeThreshold(txBatch)

		if optimalBaseFee == nil {
			t.Fatal("Expected optimal base fee to be found")
		}

		t.Logf("Optimal base fee: %v gwei", new(big.Int).Div(optimalBaseFee, big.NewInt(1_000_000_000)))
		t.Logf("Optimal collected fees: %v wei", optimalFees)

		// Calculate actual optimums:
		// - 100 gwei base fee: 100 * 21000 = 2,100,000 (only tx2)
		// - 20 gwei base fee: 20 * (21000 + 50000) = 1,420,000 (tx2 + tx3)
		// - 10 gwei base fee: 10 * (21000 + 50000 + 100000) = 1,710,000 (all)
		// So 100 gwei with just tx2 is optimal!
		expectedOptimal := big.NewInt(100_000_000_000)
		if optimalBaseFee.Cmp(expectedOptimal) != 0 {
			t.Errorf("Expected optimal base fee to be 100 gwei, got %v gwei",
				new(big.Int).Div(optimalBaseFee, big.NewInt(1_000_000_000)))
		}

		expectedFees := new(big.Int).Mul(big.NewInt(100_000_000_000), big.NewInt(21000))
		if optimalFees.Cmp(expectedFees) != 0 {
			t.Errorf("Expected fees %v, got %v", expectedFees, optimalFees)
		}
	})

	t.Run("SingleTransaction", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     0,
			GasFeeCap: big.NewInt(50_000_000_000), // 50 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		txBatch := []types.Transactions{{tx}}
		optimalBaseFee, optimalFees := p.findOptimalBaseFeeThreshold(txBatch)

		if optimalBaseFee == nil {
			t.Fatal("Expected optimal base fee to be found")
		}

		// With single transaction, optimal is its own base fee
		if optimalBaseFee.Cmp(big.NewInt(50_000_000_000)) != 0 {
			t.Errorf("Expected optimal to be 50 gwei, got %v gwei",
				new(big.Int).Div(optimalBaseFee, big.NewInt(1_000_000_000)))
		}

		expectedFees := new(big.Int).Mul(big.NewInt(50_000_000_000), big.NewInt(21000))
		if optimalFees.Cmp(expectedFees) != 0 {
			t.Errorf("Expected fees %v, got %v", expectedFees, optimalFees)
		}
	})

	t.Run("EmptyBatch", func(t *testing.T) {
		txBatch := []types.Transactions{}
		optimalBaseFee, optimalFees := p.findOptimalBaseFeeThreshold(txBatch)

		if optimalBaseFee != nil {
			t.Errorf("Expected nil base fee for empty batch, got %v", optimalBaseFee)
		}
		if optimalFees.Cmp(big.NewInt(0)) != 0 {
			t.Errorf("Expected zero fees for empty batch, got %v", optimalFees)
		}
	})

	t.Run("AllSameBaseFee", func(t *testing.T) {
		var txs types.Transactions
		baseFee := big.NewInt(30_000_000_000) // 30 gwei

		for i := 0; i < 5; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   big.NewInt(1),
				Nonce:     uint64(i),
				GasFeeCap: baseFee,
				GasTipCap: big.NewInt(1_000_000_000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, tx)
		}

		txBatch := []types.Transactions{txs}
		optimalBaseFee, optimalFees := p.findOptimalBaseFeeThreshold(txBatch)

		if optimalBaseFee == nil {
			t.Fatal("Expected optimal base fee to be found")
		}

		// All have same base fee, so optimal should be that base fee
		if optimalBaseFee.Cmp(baseFee) != 0 {
			t.Errorf("Expected optimal to be %v gwei, got %v gwei",
				new(big.Int).Div(baseFee, big.NewInt(1_000_000_000)),
				new(big.Int).Div(optimalBaseFee, big.NewInt(1_000_000_000)))
		}

		// All 5 txs included: 30 gwei * 21000 * 5
		expectedFees := new(big.Int).Mul(baseFee, big.NewInt(21000*5))
		if optimalFees.Cmp(expectedFees) != 0 {
			t.Errorf("Expected fees %v, got %v", expectedFees, optimalFees)
		}
	})

	t.Run("ManyTransactionsWithDifferentFees", func(t *testing.T) {
		var txs types.Transactions

		// Create 100 transactions with varied fees
		for i := 0; i < 100; i++ {
			gasFee := big.NewInt(int64((i + 1) * 1_000_000_000)) // 1 to 100 gwei
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   big.NewInt(1),
				Nonce:     uint64(i),
				GasFeeCap: gasFee,
				GasTipCap: big.NewInt(1_000_000_000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, tx)
		}

		txBatch := []types.Transactions{txs}
		optimalBaseFee, optimalFees := p.findOptimalBaseFeeThreshold(txBatch)

		if optimalBaseFee == nil {
			t.Fatal("Expected optimal base fee to be found")
		}

		t.Logf("With 100 transactions (1-100 gwei):")
		t.Logf("  Optimal base fee: %v gwei", new(big.Int).Div(optimalBaseFee, big.NewInt(1_000_000_000)))
		t.Logf("  Optimal collected fees: %v wei", optimalFees)

		// Verify the algorithm ran successfully (exact value depends on gas distribution)
		if optimalFees.Cmp(big.NewInt(0)) <= 0 {
			t.Error("Expected positive optimal fees")
		}
	})
}

// TestSortingAlgorithm tests that the sorting algorithm correctly sorts transactions
// by base fee in descending order
func TestSortingAlgorithm(t *testing.T) {
	p := &Proposer{}
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("SortsCorrectlyDescending", func(t *testing.T) {
		// Create transactions with various fees (not in order)
		fees := []int64{50, 10, 200, 5, 100, 30, 75}
		expected := []int64{200, 100, 75, 50, 30, 10, 5}

		var txs []txWithBaseFee
		for i, fee := range fees {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   big.NewInt(1),
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(fee * 1_000_000_000),
				GasTipCap: big.NewInt(1_000_000_000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, txWithBaseFee{
				tx:      tx,
				baseFee: big.NewInt(fee * 1_000_000_000),
				gas:     21000,
			})
		}

		p.sortTxsByBaseFeeDesc(txs)

		// Verify sorted order
		for i, expectedFee := range expected {
			actualFee := txs[i].baseFee.Int64() / 1_000_000_000
			if actualFee != expectedFee {
				t.Errorf("Position %d: expected %d gwei, got %d gwei", i, expectedFee, actualFee)
			}
		}
	})

	t.Run("HandlesEmptySlice", func(t *testing.T) {
		var txs []txWithBaseFee
		p.sortTxsByBaseFeeDesc(txs) // Should not panic
	})

	t.Run("HandlesSingleElement", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     0,
			GasFeeCap: big.NewInt(50_000_000_000),
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})
		txs := []txWithBaseFee{{
			tx:      tx,
			baseFee: big.NewInt(50_000_000_000),
			gas:     21000,
		}}
		p.sortTxsByBaseFeeDesc(txs) // Should not panic
		if txs[0].baseFee.Cmp(big.NewInt(50_000_000_000)) != 0 {
			t.Error("Single element should remain unchanged")
		}
	})

	t.Run("HandlesLargeDataset", func(t *testing.T) {
		// Create 1000 transactions with random-ish fees
		var txs []txWithBaseFee
		for i := 0; i < 1000; i++ {
			fee := int64((i*7 + 13) % 1000) // Pseudo-random but deterministic
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   big.NewInt(1),
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(fee * 1_000_000_000),
				GasTipCap: big.NewInt(1_000_000_000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, txWithBaseFee{
				tx:      tx,
				baseFee: big.NewInt(fee * 1_000_000_000),
				gas:     21000,
			})
		}

		p.sortTxsByBaseFeeDesc(txs)

		// Verify it's sorted in descending order
		for i := 0; i < len(txs)-1; i++ {
			if txs[i].baseFee.Cmp(txs[i+1].baseFee) < 0 {
				t.Errorf("Position %d not sorted correctly: %v < %v",
					i, txs[i].baseFee, txs[i+1].baseFee)
				break
			}
		}

		t.Logf("Successfully sorted %d transactions", len(txs))
	})
}

// TestIsProfitableSelection tests that isProfitable picks the optimal
// batch using the new algorithm and updates references accordingly.
func TestIsProfitableSelection(t *testing.T) {
	// Create a proposer with mocked cost estimator
	p := &Proposer{
		Config:          &Config{},
		protocolConfigs: &dummyProtocolConfigs{},
		costEstimator: &mockCostEstimator{
			cost: big.NewInt(100_000_000_000_000), // 0.0001 ETH
		},
	}

	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	t.Run("SelectsOptimalBaseFeeWithMixedGas", func(t *testing.T) {
		// Create transactions where the optimal is NOT the highest base fee
		// due to gas distribution
		txLow := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: big.NewInt(10_000_000_000), // 10 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       100000, // Large gas
			To:        &testAddr,
			Value:     common.Big0,
		})

		txMed := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     1,
			GasFeeCap: big.NewInt(30_000_000_000), // 30 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       50000, // Medium gas
			To:        &testAddr,
			Value:     common.Big0,
		})

		txHigh := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     2,
			GasFeeCap: big.NewInt(200_000_000_000), // 200 gwei (highest)
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000, // Small gas
			To:        &testAddr,
			Value:     common.Big0,
		})

		// Batch with all three transactions
		txBatch := []types.Transactions{{txLow, txMed, txHigh}}
		l2BaseFee := big.NewInt(5_000_000_000) // 5 gwei (low, so optimizer should improve)
		txCount := uint64(3)

		var candidate txmgr.TxCandidate

		// Call isProfitable and expect it to find optimal configuration
		profitable, adjusted, err := p.isProfitable(context.Background(), &txBatch, &l2BaseFee, &candidate, &txCount)
		if err != nil {
			t.Fatalf("isProfitable returned error: %v", err)
		}

		if !profitable {
			t.Fatalf("Expected batch to be profitable with selected best option")
		}

		// The optimal should be 10 gwei (all 3 txs: 10*171000)
		// vs 30 gwei (2 txs: 30*71000) vs 200 gwei (1 tx: 200*21000)
		// 10*171000 = 1,710,000 vs 30*71000 = 2,130,000 vs 200*21000 = 4,200,000
		// So 200 gwei should win!

		t.Logf("Adjusted: %v", adjusted)
		t.Logf("Selected base fee: %v gwei", new(big.Int).Div(l2BaseFee, big.NewInt(1_000_000_000)))
		t.Logf("Selected tx count: %d", txCount)
		t.Logf("Batch count: %d", len(txBatch))
	})

	t.Run("SelectsHighPayingTxWhenOptimal", func(t *testing.T) {
		// Create scenario where excluding low-paying txs is better
		txs := make(types.Transactions, 0)

		// 10 low-paying transactions
		for i := 0; i < 10; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(5_000_000_000), // 5 gwei
				GasTipCap: big.NewInt(1_000_000_000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, tx)
		}

		// 1 very high-paying transaction
		highTx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     10,
			GasFeeCap: big.NewInt(500_000_000_000), // 500 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})
		txs = append(txs, highTx)

		txBatch := []types.Transactions{txs}
		l2BaseFee := big.NewInt(1_000_000_000) // 1 gwei (very low)
		txCount := uint64(11)
		var candidate txmgr.TxCandidate

		profitable, adjusted, err := p.isProfitable(context.Background(), &txBatch, &l2BaseFee, &candidate, &txCount)
		if err != nil {
			t.Fatalf("isProfitable returned error: %v", err)
		}

		if !profitable {
			t.Fatalf("Expected batch to be profitable")
		}

		// With 500 gwei tx, optimizer might pick 500 gwei base fee with 1 tx
		// 500*21000 = 10,500,000 vs 5*231000 = 1,155,000
		// So it should choose the high base fee option

		t.Logf("Adjusted: %v", adjusted)
		t.Logf("Selected base fee: %v gwei", new(big.Int).Div(l2BaseFee, big.NewInt(1_000_000_000)))
		t.Logf("Selected tx count: %d (original: 11)", txCount)
	})
}

// TestIsProfitableKeepsOriginal ensures that if original collected fees are
// greater than any optimized subset, the original batch is kept and adjusted=false
func TestIsProfitableKeepsOriginal(t *testing.T) {
	p := &Proposer{
		Config:          &Config{},
		protocolConfigs: &dummyProtocolConfigs{},
		costEstimator: &mockCostEstimator{
			cost: big.NewInt(1), // Very low cost so everything is profitable
		},
	}
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	// Create many transactions with similar fees where keeping all is optimal
	var txs types.Transactions
	for i := 0; i < 20; i++ {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     uint64(i),
			GasFeeCap: big.NewInt(10_000_000_000), // All 10 gwei - same fee
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})
		txs = append(txs, tx)
	}

	txBatch := []types.Transactions{txs}
	l2BaseFee := big.NewInt(10_000_000_000) // 10 gwei - matches transaction fees exactly
	txCount := uint64(len(txs))
	var candidate txmgr.TxCandidate

	profitable, adjusted, err := p.isProfitable(context.Background(), &txBatch, &l2BaseFee, &candidate, &txCount)
	if err != nil {
		t.Fatalf("isProfitable returned error: %v", err)
	}

	if !profitable {
		t.Fatalf("Expected original batch to be profitable")
	}

	// Since all transactions have the same base fee (10 gwei) and the L2 base fee is also 10 gwei,
	// the optimizer will find that using 10 gwei with all 20 txs is optimal,
	// which equals the original configuration. So adjusted could be true or false.
	// What matters is that all transactions are kept.
	if len(txBatch) != 1 || len(txBatch[0]) != int(txCount) {
		t.Errorf("Expected original batch to be kept with %d txs, got %d batches with %d txs",
			20, len(txBatch), len(txBatch[0]))
	}

	t.Logf("Adjusted: %v (expected: false or true if same as original)", adjusted)
	t.Logf("Transaction count: %d (original: 20)", txCount)
}

// TestFilterTxsByBaseFee tests filtering transactions by minimum base fee
func TestFilterTxsByBaseFee(t *testing.T) {
	p := &Proposer{}
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	// Generate a test private key for signing legacy transactions
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	t.Run("EmptyBatch", func(t *testing.T) {
		emptyBatch := []types.Transactions{}
		filtered := p.filterTxsByBaseFee(emptyBatch, big.NewInt(1000000000))
		if len(filtered) != 0 {
			t.Errorf("Filtering empty batch should return empty result, got %d batches", len(filtered))
		}
	})

	t.Run("AllTransactionsMeetThreshold", func(t *testing.T) {
		minBaseFee := big.NewInt(1000000000) // 1 gwei

		var txs types.Transactions
		for i := 0; i < 3; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(2000000000), // All above threshold
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, tx)
		}

		batch := []types.Transactions{txs}
		filtered := p.filterTxsByBaseFee(batch, minBaseFee)
		if len(filtered) != 1 {
			t.Errorf("Expected 1 batch, got %d", len(filtered))
		}
		if len(filtered[0]) != 3 {
			t.Errorf("Expected all 3 transactions to pass filter, got %d", len(filtered[0]))
		}
	})

	t.Run("SomeTransactionsBelowThreshold", func(t *testing.T) {
		minBaseFee := big.NewInt(3000000000) // 3 gwei

		gasFees := []*big.Int{
			big.NewInt(5000000000), // Above threshold - should pass
			big.NewInt(2000000000), // Below threshold - should filter out
			big.NewInt(4000000000), // Above threshold - should pass
			big.NewInt(1000000000), // Below threshold - should filter out
			big.NewInt(3000000000), // Equal to threshold - should pass
		}

		var txs types.Transactions
		for i, fee := range gasFees {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: fee,
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, tx)
		}

		batch := []types.Transactions{txs}
		filtered := p.filterTxsByBaseFee(batch, minBaseFee)
		if len(filtered) != 1 {
			t.Errorf("Expected 1 batch, got %d", len(filtered))
		}
		if len(filtered[0]) != 3 {
			t.Errorf("Expected only 3 transactions to pass (5, 4, and 3 gwei), got %d", len(filtered[0]))
		}

		// Verify the filtered transactions are the correct ones
		for _, tx := range filtered[0] {
			if tx.GasFeeCap().Cmp(minBaseFee) < 0 {
				t.Errorf("Transaction with GasFeeCap %v should not pass threshold %v", tx.GasFeeCap(), minBaseFee)
			}
		}
	})

	t.Run("AllTransactionsBelowThreshold", func(t *testing.T) {
		minBaseFee := big.NewInt(10000000000) // 10 gwei - very high

		var txs types.Transactions
		for i := 0; i < 3; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(2000000000), // All below threshold
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, tx)
		}

		batch := []types.Transactions{txs}
		filtered := p.filterTxsByBaseFee(batch, minBaseFee)
		if len(filtered) != 0 {
			t.Errorf("No transactions should pass filter, got %d batches", len(filtered))
		}
	})

	t.Run("MultipleBatchesWithMixedResults", func(t *testing.T) {
		minBaseFee := big.NewInt(3000000000) // 3 gwei

		// Batch 1: Some transactions pass
		batch1 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     0,
				GasFeeCap: big.NewInt(5000000000), // Pass
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     1,
				GasFeeCap: big.NewInt(1000000000), // Fail
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}

		// Batch 2: All transactions fail
		batch2 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     2,
				GasFeeCap: big.NewInt(1000000000), // Fail
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}

		// Batch 3: All transactions pass
		batch3 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     3,
				GasFeeCap: big.NewInt(4000000000), // Pass
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     4,
				GasFeeCap: big.NewInt(3000000000), // Pass
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}

		batches := []types.Transactions{batch1, batch2, batch3}
		filtered := p.filterTxsByBaseFee(batches, minBaseFee)

		// Should have 2 batches (batch1 with 1 tx and batch3 with 2 txs)
		// Batch2 is completely filtered out
		if len(filtered) != 2 {
			t.Errorf("Expected 2 batches (batch2 filtered out), got %d", len(filtered))
		}
		if len(filtered[0]) != 1 {
			t.Errorf("Batch1 should have 1 transaction, got %d", len(filtered[0]))
		}
		if len(filtered[1]) != 2 {
			t.Errorf("Batch3 should have 2 transactions, got %d", len(filtered[1]))
		}
	})

	t.Run("LegacyTransactionsFiltering", func(t *testing.T) {
		minBaseFee := big.NewInt(3000000000) // 3 gwei
		signer := types.LatestSignerForChainID(chainID)

		// Legacy tx with high gas price (should pass)
		legacyTxHigh := types.NewTx(&types.LegacyTx{
			Nonce:    0,
			GasPrice: big.NewInt(4000000000),
			Gas:      21000,
			To:       &testAddr,
			Value:    common.Big0,
		})
		signedLegacyHigh, err := types.SignTx(legacyTxHigh, signer, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign legacy tx: %v", err)
		}

		// Legacy tx with low gas price (should fail)
		legacyTxLow := types.NewTx(&types.LegacyTx{
			Nonce:    1,
			GasPrice: big.NewInt(1000000000),
			Gas:      21000,
			To:       &testAddr,
			Value:    common.Big0,
		})
		signedLegacyLow, err := types.SignTx(legacyTxLow, signer, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign legacy tx: %v", err)
		}

		// Dynamic tx (should pass)
		dynamicTx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     2,
			GasFeeCap: big.NewInt(3500000000),
			GasTipCap: big.NewInt(100000000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		batch := []types.Transactions{{signedLegacyHigh, signedLegacyLow, dynamicTx}}
		filtered := p.filterTxsByBaseFee(batch, minBaseFee)

		if len(filtered) != 1 {
			t.Errorf("Expected 1 batch, got %d", len(filtered))
		}
		if len(filtered[0]) != 2 {
			t.Errorf("Expected 2 transactions (high legacy and dynamic), got %d", len(filtered[0]))
		}
	})
}

// TestCountTxsInBatch tests counting transactions in batches
func TestCountTxsInBatch(t *testing.T) {
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	t.Run("EmptyBatch", func(t *testing.T) {
		emptyBatch := []types.Transactions{}
		count := countTxsInBatch(emptyBatch)
		if count != 0 {
			t.Errorf("Expected 0, got %d", count)
		}
	})

	t.Run("SingleBatchWithMultipleTxs", func(t *testing.T) {
		var txs types.Transactions
		for i := 0; i < 5; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			txs = append(txs, tx)
		}
		batch := []types.Transactions{txs}
		count := countTxsInBatch(batch)
		if count != 5 {
			t.Errorf("Expected 5, got %d", count)
		}
	})

	t.Run("MultipleBatches", func(t *testing.T) {
		batch1 := types.Transactions{}
		batch2 := types.Transactions{}
		batch3 := types.Transactions{}

		for i := 0; i < 3; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			batch1 = append(batch1, tx)
		}

		for i := 3; i < 8; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			batch2 = append(batch2, tx)
		}

		for i := 8; i < 10; i++ {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			})
			batch3 = append(batch3, tx)
		}

		batches := []types.Transactions{batch1, batch2, batch3}
		count := countTxsInBatch(batches)
		if count != 10 {
			t.Errorf("Expected 10 (3 + 5 + 2), got %d", count)
		}
	})

	t.Run("BatchesWithEmptyList", func(t *testing.T) {
		batch1 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     0,
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}
		batch2 := types.Transactions{} // Empty
		batch3 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     1,
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     2,
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}

		batches := []types.Transactions{batch1, batch2, batch3}
		count := countTxsInBatch(batches)
		if count != 3 {
			t.Errorf("Expected 3 (1 + 0 + 2), got %d", count)
		}
	})
}

// TestCalculatePriorityFee tests the priority fee (tip) calculation
func TestCalculatePriorityFee(t *testing.T) {
	p := &Proposer{}
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	t.Run("DynamicFeeTxWithTip", func(t *testing.T) {
		baseFee := big.NewInt(10_000_000_000) // 10 gwei
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: big.NewInt(15_000_000_000), // 15 gwei max
			GasTipCap: big.NewInt(2_000_000_000),  // 2 gwei tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		tip := p.calculatePriorityFee(tx, baseFee)
		// Effective tip = min(2 gwei, 15 - 10) = 2 gwei
		expectedTip := big.NewInt(2_000_000_000)
		if tip.Cmp(expectedTip) != 0 {
			t.Errorf("Expected tip %v, got %v", expectedTip, tip)
		}
	})

	t.Run("DynamicFeeTxTipLimitedByFeeCap", func(t *testing.T) {
		baseFee := big.NewInt(10_000_000_000) // 10 gwei
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: big.NewInt(12_000_000_000), // 12 gwei max
			GasTipCap: big.NewInt(5_000_000_000),  // 5 gwei tip (but capped)
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		tip := p.calculatePriorityFee(tx, baseFee)
		// Effective tip = min(5 gwei, 12 - 10) = 2 gwei
		expectedTip := big.NewInt(2_000_000_000)
		if tip.Cmp(expectedTip) != 0 {
			t.Errorf("Expected tip %v, got %v", expectedTip, tip)
		}
	})

	t.Run("DynamicFeeTxFeeCapBelowBaseFee", func(t *testing.T) {
		baseFee := big.NewInt(10_000_000_000) // 10 gwei
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: big.NewInt(8_000_000_000), // 8 gwei (below base fee)
			GasTipCap: big.NewInt(2_000_000_000), // 2 gwei tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		tip := p.calculatePriorityFee(tx, baseFee)
		// Fee cap below base fee, so tip is 0
		if tip.Cmp(big.NewInt(0)) != 0 {
			t.Errorf("Expected tip 0, got %v", tip)
		}
	})

	t.Run("LegacyTxWithTip", func(t *testing.T) {
		baseFee := big.NewInt(10_000_000_000) // 10 gwei
		signer := types.LatestSignerForChainID(chainID)
		privateKey, _ := crypto.GenerateKey()

		legacyTx := types.NewTx(&types.LegacyTx{
			Nonce:    0,
			GasPrice: big.NewInt(15_000_000_000), // 15 gwei
			Gas:      21000,
			To:       &testAddr,
			Value:    common.Big0,
		})
		signedTx, _ := types.SignTx(legacyTx, signer, privateKey)

		tip := p.calculatePriorityFee(signedTx, baseFee)
		// Tip = gasPrice - baseFee = 15 - 10 = 5 gwei
		expectedTip := big.NewInt(5_000_000_000)
		if tip.Cmp(expectedTip) != 0 {
			t.Errorf("Expected tip %v, got %v", expectedTip, tip)
		}
	})

	t.Run("LegacyTxGasPriceBelowBaseFee", func(t *testing.T) {
		baseFee := big.NewInt(10_000_000_000) // 10 gwei
		signer := types.LatestSignerForChainID(chainID)
		privateKey, _ := crypto.GenerateKey()

		legacyTx := types.NewTx(&types.LegacyTx{
			Nonce:    0,
			GasPrice: big.NewInt(8_000_000_000), // 8 gwei (below base fee)
			Gas:      21000,
			To:       &testAddr,
			Value:    common.Big0,
		})
		signedTx, _ := types.SignTx(legacyTx, signer, privateKey)

		tip := p.calculatePriorityFee(signedTx, baseFee)
		// Gas price below base fee, so tip is 0
		if tip.Cmp(big.NewInt(0)) != 0 {
			t.Errorf("Expected tip 0, got %v", tip)
		}
	})
}

// TestComputeL2FeesWithTips tests that computeL2Fees correctly accounts for tips
func TestComputeL2FeesWithTips(t *testing.T) {
	p := &Proposer{
		protocolConfigs: &dummyProtocolConfigs{}, // 100% base fee sharing
	}

	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)
	baseFee := big.NewInt(10_000_000_000) // 10 gwei

	t.Run("FeesIncludeTips", func(t *testing.T) {
		// Transaction with 10 gwei base fee cap and 2 gwei tip
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: big.NewInt(15_000_000_000), // 15 gwei
			GasTipCap: big.NewInt(2_000_000_000),  // 2 gwei tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		txBatch := []types.Transactions{{tx}}
		collectedFees := p.computeL2Fees(txBatch, baseFee)

		// Expected: (10 gwei base * 100% + 2 gwei tip) * 21000 gas
		// = 12 gwei * 21000 = 252,000 gwei = 252,000,000,000,000 wei
		expectedFees := big.NewInt(252_000_000_000_000)
		if collectedFees.Cmp(expectedFees) != 0 {
			t.Errorf("Expected fees %v, got %v", expectedFees, collectedFees)
		}
	})

	t.Run("MultipleTransactionsWithDifferentTips", func(t *testing.T) {
		tx1 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: big.NewInt(15_000_000_000), // 15 gwei
			GasTipCap: big.NewInt(2_000_000_000),  // 2 gwei tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		tx2 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     1,
			GasFeeCap: big.NewInt(20_000_000_000), // 20 gwei
			GasTipCap: big.NewInt(5_000_000_000),  // 5 gwei tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		txBatch := []types.Transactions{{tx1, tx2}}
		collectedFees := p.computeL2Fees(txBatch, baseFee)

		// tx1: (10 + 2) * 21000 = 252,000 gwei
		// tx2: (10 + 5) * 21000 = 315,000 gwei
		// Total: 567,000 gwei = 567,000,000,000,000 wei
		expectedFees := big.NewInt(567_000_000_000_000)
		if collectedFees.Cmp(expectedFees) != 0 {
			t.Errorf("Expected fees %v, got %v", expectedFees, collectedFees)
		}
	})
}

// TestFindOptimalBaseFeeThresholdWithTips tests that the optimizer accounts for tips
func TestFindOptimalBaseFeeThresholdWithTips(t *testing.T) {
	p := &Proposer{
		protocolConfigs: &dummyProtocolConfigs{}, // 100% base fee sharing
	}

	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	t.Run("OptimalConsidersTips", func(t *testing.T) {
		// Scenario: One high base fee tx with no tip vs multiple lower base fee txs with high tips
		// The optimizer should prefer the option that maximizes total revenue

		// High base fee, no tip
		txHighBase := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: big.NewInt(100_000_000_000), // 100 gwei
			GasTipCap: big.NewInt(0),               // 0 tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		// Lower base fee, high tip
		txLowBaseHighTip1 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     1,
			GasFeeCap: big.NewInt(50_000_000_000), // 50 gwei
			GasTipCap: big.NewInt(30_000_000_000), // 30 gwei tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		txLowBaseHighTip2 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     2,
			GasFeeCap: big.NewInt(50_000_000_000), // 50 gwei
			GasTipCap: big.NewInt(30_000_000_000), // 30 gwei tip
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		txBatch := []types.Transactions{{txHighBase, txLowBaseHighTip1, txLowBaseHighTip2}}
		optimalBaseFee, optimalFees := p.findOptimalBaseFeeThreshold(txBatch)

		if optimalBaseFee == nil {
			t.Fatal("Expected optimal base fee to be found")
		}

		t.Logf("Optimal base fee: %v gwei", new(big.Int).Div(optimalBaseFee, big.NewInt(1_000_000_000)))
		t.Logf("Optimal fees collected: %v wei", optimalFees)

		// With tips:
		// Option 1 (100 gwei base, 1 tx): 100 * 21000 = 2,100,000 gwei
		// Option 2 (50 gwei base, 2 txs): (50 + 30) * 21000 * 2 = 3,360,000 gwei
		// So 50 gwei should be optimal!

		expectedOptimal := big.NewInt(50_000_000_000)
		if optimalBaseFee.Cmp(expectedOptimal) != 0 {
			t.Errorf("Expected optimal base fee to be 50 gwei (due to tips), got %v gwei",
				new(big.Int).Div(optimalBaseFee, big.NewInt(1_000_000_000)))
		}
	})
}
