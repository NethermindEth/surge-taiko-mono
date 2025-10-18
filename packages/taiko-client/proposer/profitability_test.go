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

// testProposer embeds Proposer and allows tests to override estimated cost
type testProposer struct {
	Proposer
	estimatedCost *big.Int
}

// estimateL2Cost overrides the Proposer method for tests to avoid RPC calls.
func (p *testProposer) estimateL2Cost(ctx context.Context, candidate *txmgr.TxCandidate) (*big.Int, error) {
	if p.estimatedCost != nil {
		return new(big.Int).Set(p.estimatedCost), nil
	}
	return big.NewInt(0), nil
}

// isProfitable mirrors the production logic but uses p.estimatedCost instead of calling RPC
func (p *testProposer) isProfitable(ctx context.Context, txBatch *[]types.Transactions, l2BaseFee **big.Int, candidate *txmgr.TxCandidate, txs *uint64) (bool, bool, error) {
	estimatedCost := big.NewInt(0)
	if p.estimatedCost != nil {
		estimatedCost = new(big.Int).Set(p.estimatedCost)
	}

	// Compute collected fees for the original batch/base fee
	originalCollectedFees := p.computeL2Fees(*txBatch, *l2BaseFee)

	// We'll keep track of the best option found (highest collected fees)
	bestCollectedFees := new(big.Int).Set(originalCollectedFees)
	bestTxBatch := *txBatch
	bestL2BaseFee := new(big.Int).Set(*l2BaseFee)
	bestTxs := *txs
	bestAdjusted := false

	// If there are no candidate transactions with a higher base fee, we only consider the original
	highestTxBaseFee := p.findHighestBaseFeeInBatch(*txBatch)
	if highestTxBaseFee == nil || highestTxBaseFee.Cmp(*l2BaseFee) <= 0 {
		// No higher base fee transactions found; decide profitability based on originalCollectedFees
		isProfitable := originalCollectedFees.Cmp(estimatedCost) >= 0
		return isProfitable, false, nil
	}

	// Try different percentage thresholds: 50%, then 75%
	percentages := []int64{50, 75}

	for _, percentage := range percentages {
		adjustedBaseFee := new(big.Int).Mul(highestTxBaseFee, big.NewInt(percentage))
		adjustedBaseFee = new(big.Int).Div(adjustedBaseFee, big.NewInt(100))

		filteredTxBatch := p.filterTxsByBaseFee(*txBatch, adjustedBaseFee)
		if len(filteredTxBatch) == 0 {
			continue
		}

		collectedFeesAdjusted := p.computeL2Fees(filteredTxBatch, adjustedBaseFee)

		if collectedFeesAdjusted.Cmp(bestCollectedFees) > 0 {
			bestCollectedFees.Set(collectedFeesAdjusted)
			bestTxBatch = filteredTxBatch
			bestL2BaseFee.Set(adjustedBaseFee)
			bestTxs = countTxsInBatch(filteredTxBatch)
			bestAdjusted = true
		}
	}

	if bestAdjusted {
		*txBatch = bestTxBatch
		*l2BaseFee = bestL2BaseFee
		*txs = bestTxs
	}

	isProfitable := bestCollectedFees.Cmp(estimatedCost) >= 0
	return isProfitable, bestAdjusted, nil
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

// TestIsProfitableReferenceModification tests that the refactored isProfitable
// correctly modifies the references when using adjusted base fee
func TestIsProfitableReferenceModification(t *testing.T) {
	// Create a test proposer (minimal setup)
	p := &Proposer{}

	t.Run("Modifies references when using 50% threshold", func(t *testing.T) {
		// Create test transactions
		testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

		// Create transactions with different gas fees
		tx1 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     0,
			GasFeeCap: big.NewInt(10_000_000_000), // 10 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		tx2 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     1,
			GasFeeCap: big.NewInt(100_000_000_000), // 100 gwei (highest)
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		tx3 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     2,
			GasFeeCap: big.NewInt(20_000_000_000), // 20 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		// Create batch
		txBatch := []types.Transactions{{tx1, tx2, tx3}}
		l2BaseFee := big.NewInt(5_000_000_000) // 5 gwei
		txCount := uint64(3)

		// Store original values
		originalBaseFee := new(big.Int).Set(l2BaseFee)
		originalTxCount := txCount

		// Test: find highest base fee
		highestFee := p.findHighestBaseFeeInBatch(txBatch)
		if highestFee.Cmp(big.NewInt(100_000_000_000)) != 0 {
			t.Errorf("Expected highest fee to be 100 gwei, got %v", highestFee)
		}

		// Test: calculate 50% threshold
		threshold50 := new(big.Int).Mul(highestFee, big.NewInt(50))
		threshold50 = new(big.Int).Div(threshold50, big.NewInt(100))
		expectedThreshold50 := big.NewInt(50_000_000_000) // 50 gwei

		if threshold50.Cmp(expectedThreshold50) != 0 {
			t.Errorf("Expected 50%% threshold to be 50 gwei, got %v", threshold50)
		}

		// Test: filter by 50% threshold
		filtered := p.filterTxsByBaseFee(txBatch, threshold50)

		// Only tx2 (100 gwei) should pass
		if len(filtered) != 1 {
			t.Errorf("Expected 1 filtered batch, got %d", len(filtered))
		}
		if len(filtered[0]) != 1 {
			t.Errorf("Expected 1 transaction in filtered batch, got %d", len(filtered[0]))
		}

		// Verify the filtered transaction is tx2
		if filtered[0][0].GasFeeCap().Cmp(big.NewInt(100_000_000_000)) != 0 {
			t.Errorf("Expected filtered transaction to have 100 gwei gas fee cap, got %v", filtered[0][0].GasFeeCap())
		}

		// Test: count transactions
		count := countTxsInBatch(filtered)
		if count != 1 {
			t.Errorf("Expected count to be 1, got %d", count)
		}

		t.Logf("✓ Original base fee: %v gwei", new(big.Int).Div(originalBaseFee, big.NewInt(1_000_000_000)))
		t.Logf("✓ Highest transaction fee: %v gwei", new(big.Int).Div(highestFee, big.NewInt(1_000_000_000)))
		t.Logf("✓ 50%% threshold: %v gwei", new(big.Int).Div(threshold50, big.NewInt(1_000_000_000)))
		t.Logf("✓ Original transaction count: %d", originalTxCount)
		t.Logf("✓ Filtered transaction count: %d", count)
		t.Logf("✓ Transactions filtered from %d to %d", originalTxCount, count)
	})

	t.Run("Tries 75% if 50% doesn't have transactions", func(t *testing.T) {
		testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

		// Create transactions where only 75% threshold will include them
		tx1 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     0,
			GasFeeCap: big.NewInt(80_000_000_000), // 80 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		tx2 := types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     1,
			GasFeeCap: big.NewInt(100_000_000_000), // 100 gwei (highest)
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		// Create batch
		txBatch := []types.Transactions{{tx1, tx2}}

		// Test 50% threshold (should exclude both since 50 gwei < 80 gwei)
		highestFee := p.findHighestBaseFeeInBatch(txBatch)
		threshold50 := new(big.Int).Mul(highestFee, big.NewInt(50))
		threshold50 = new(big.Int).Div(threshold50, big.NewInt(100))

		filtered50 := p.filterTxsByBaseFee(txBatch, threshold50)
		count50 := countTxsInBatch(filtered50)

		// Both transactions should pass 50% (50 gwei)
		if count50 != 2 {
			t.Errorf("Expected 2 transactions to pass 50%% threshold, got %d", count50)
		}

		// Test 75% threshold
		threshold75 := new(big.Int).Mul(highestFee, big.NewInt(75))
		threshold75 = new(big.Int).Div(threshold75, big.NewInt(100))

		filtered75 := p.filterTxsByBaseFee(txBatch, threshold75)
		count75 := countTxsInBatch(filtered75)

		// Both transactions (80 and 100 gwei) should pass 75% threshold (75 gwei)
		if count75 != 2 {
			t.Errorf("Expected 2 transactions to pass 75%% threshold, got %d", count75)
		}

		t.Logf("Highest fee: %v gwei", new(big.Int).Div(highestFee, big.NewInt(1_000_000_000)))
		t.Logf("50%% threshold: %v gwei (passes: %d txs)", new(big.Int).Div(threshold50, big.NewInt(1_000_000_000)), count50)
		t.Logf("75%% threshold: %v gwei (passes: %d txs)", new(big.Int).Div(threshold75, big.NewInt(1_000_000_000)), count75)
	})
}

// TestIsProfitableSelection tests that isProfitable picks the highest-paying
// batch (original or adjusted) and updates references accordingly.
func TestIsProfitableSelection(t *testing.T) {
	// Use package-level testProposer and set estimatedCost for this test
	tp := &testProposer{estimatedCost: big.NewInt(100_000_000_000_000)}
	tp.protocolConfigs = &dummyProtocolConfigs{}

	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	// Create three transactions: one low, one very high, one medium
	txLow := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     0,
		GasFeeCap: big.NewInt(5_000_000_000), // 5 gwei
		GasTipCap: big.NewInt(1_000_000_000),
		Gas:       21000,
		To:        &testAddr,
		Value:     common.Big0,
	})

	txHigh := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     1,
		GasFeeCap: big.NewInt(200_000_000_000), // 200 gwei (highest)
		GasTipCap: big.NewInt(1_000_000_000),
		Gas:       21000,
		To:        &testAddr,
		Value:     common.Big0,
	})

	txMid := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     2,
		GasFeeCap: big.NewInt(30_000_000_000), // 30 gwei
		GasTipCap: big.NewInt(1_000_000_000),
		Gas:       21000,
		To:        &testAddr,
		Value:     common.Big0,
	})

	// Batch with all three transactions
	txBatch := []types.Transactions{{txLow, txHigh, txMid}}
	l2BaseFee := big.NewInt(10_000_000_000) // 10 gwei
	txCount := uint64(3)

	// Wrap in txmgr candidate (not used by our fake estimate but needed for signature)
	var candidate txmgr.TxCandidate

	// Call isProfitable and expect it to select the adjusted batch (50% of highest=100 gwei)
	profitable, adjusted, err := tp.isProfitable(context.Background(), &txBatch, &l2BaseFee, &candidate, &txCount)
	if err != nil {
		t.Fatalf("isProfitable returned error: %v", err)
	}

	if !profitable {
		t.Fatalf("Expected batch to be profitable with selected best option")
	}

	if !adjusted {
		t.Fatalf("Expected adjusted to be true since filtered batch should yield higher fees")
	}

	// After adjustment, only txHigh should remain (200 gwei >= 100 gwei threshold)
	if len(txBatch) != 1 || len(txBatch[0]) != 1 {
		t.Fatalf("Expected filtered batch to contain only the highest tx, got %v", txBatch)
	}

	if txBatch[0][0].GasFeeCap().Cmp(big.NewInt(200_000_000_000)) != 0 {
		t.Fatalf("Expected remaining tx to be the 200 gwei tx, got %v", txBatch[0][0].GasFeeCap())
	}
}

// TestIsProfitableKeepsOriginal ensures that if original collected fees are
// greater than any adjusted subset, the original batch is kept and adjusted=false
func TestIsProfitableKeepsOriginal(t *testing.T) {
	tp := &testProposer{estimatedCost: big.NewInt(1)}
	tp.protocolConfigs = &dummyProtocolConfigs{}
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := big.NewInt(1)

	// Create many medium-paying transactions and one very high-paying tx.
	var txs types.Transactions
	for i := 0; i < 10; i++ {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     uint64(i),
			GasFeeCap: big.NewInt(11_000_000_000), // 11 gwei
			GasTipCap: big.NewInt(1_000_000_000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})
		txs = append(txs, tx)
	}

	// One very high-paying transaction
	highTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     11,
		GasFeeCap: big.NewInt(25_000_000_000), // 25 gwei (highest)
		GasTipCap: big.NewInt(1_000_000_000),
		Gas:       21000,
		To:        &testAddr,
		Value:     common.Big0,
	})
	txs = append(txs, highTx)

	txBatch := []types.Transactions{txs}
	l2BaseFee := big.NewInt(10_000_000_000) // 10 gwei
	txCount := uint64(len(txs))
	var candidate txmgr.TxCandidate

	profitable, adjusted, err := tp.isProfitable(context.Background(), &txBatch, &l2BaseFee, &candidate, &txCount)
	if err != nil {
		t.Fatalf("isProfitable returned error: %v", err)
	}

	if !profitable {
		t.Fatalf("Expected original batch to be profitable")
	}

	if adjusted {
		t.Fatalf("Expected adjusted to be false since original is best")
	}

	// Ensure original batch remains intact (all transactions kept)
	if len(txBatch) != 1 || len(txBatch[0]) != int(txCount) {
		t.Fatalf("Expected original batch to be kept with %d txs, got %v", txCount, txBatch)
	}
}

// TestFindHighestBaseFeeInBatch tests finding the highest base fee in a batch of transactions
func TestFindHighestBaseFeeInBatch(t *testing.T) {
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
		highestFee := p.findHighestBaseFeeInBatch(emptyBatch)
		if highestFee != nil {
			t.Errorf("Empty batch should return nil, got %v", highestFee)
		}
	})

	t.Run("SingleTransaction", func(t *testing.T) {
		gasFeeCap := big.NewInt(1000000000) // 1 gwei
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			GasFeeCap: gasFeeCap,
			GasTipCap: big.NewInt(100000000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})
		batch := []types.Transactions{{tx}}
		highestFee := p.findHighestBaseFeeInBatch(batch)
		if highestFee == nil {
			t.Error("Expected non-nil highest fee")
			return
		}
		if highestFee.Cmp(gasFeeCap) != 0 {
			t.Errorf("Expected %v, got %v", gasFeeCap, highestFee)
		}
	})

	t.Run("MultipleTransactions", func(t *testing.T) {
		gasFees := []*big.Int{
			big.NewInt(1000000000), // 1 gwei
			big.NewInt(5000000000), // 5 gwei (highest)
			big.NewInt(2000000000), // 2 gwei
			big.NewInt(3000000000), // 3 gwei
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
		highestFee := p.findHighestBaseFeeInBatch(batch)
		if highestFee == nil {
			t.Error("Expected non-nil highest fee")
			return
		}
		if highestFee.Int64() != 5000000000 {
			t.Errorf("Expected 5000000000, got %v", highestFee.Int64())
		}
	})

	t.Run("LegacyTransactions", func(t *testing.T) {
		signer := types.LatestSignerForChainID(chainID)

		legacyTx := types.NewTx(&types.LegacyTx{
			Nonce:    0,
			GasPrice: big.NewInt(3000000000), // 3 gwei
			Gas:      21000,
			To:       &testAddr,
			Value:    common.Big0,
		})
		signedLegacyTx, err := types.SignTx(legacyTx, signer, privateKey)
		if err != nil {
			t.Fatalf("Failed to sign legacy tx: %v", err)
		}

		dynamicTx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     1,
			GasFeeCap: big.NewInt(2000000000), // 2 gwei
			GasTipCap: big.NewInt(100000000),
			Gas:       21000,
			To:        &testAddr,
			Value:     common.Big0,
		})

		batch := []types.Transactions{{signedLegacyTx, dynamicTx}}
		highestFee := p.findHighestBaseFeeInBatch(batch)
		if highestFee == nil {
			t.Error("Expected non-nil highest fee")
			return
		}
		if highestFee.Int64() != 3000000000 {
			t.Errorf("Expected legacy transaction's GasPrice (3000000000), got %v", highestFee.Int64())
		}
	})

	t.Run("MultipleBatches", func(t *testing.T) {
		batch1 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     0,
				GasFeeCap: big.NewInt(2000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}

		batch2 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     1,
				GasFeeCap: big.NewInt(7000000000), // Highest
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}

		batch3 := types.Transactions{
			types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainID,
				Nonce:     2,
				GasFeeCap: big.NewInt(4000000000),
				GasTipCap: big.NewInt(100000000),
				Gas:       21000,
				To:        &testAddr,
				Value:     common.Big0,
			}),
		}

		batches := []types.Transactions{batch1, batch2, batch3}
		highestFee := p.findHighestBaseFeeInBatch(batches)
		if highestFee == nil {
			t.Error("Expected non-nil highest fee")
			return
		}
		if highestFee.Int64() != 7000000000 {
			t.Errorf("Expected 7000000000, got %v", highestFee.Int64())
		}
	})
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
