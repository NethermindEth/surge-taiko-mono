package proposer

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/utils"
)

// CostEstimator is an interface for estimating the L1 cost of proposing a batch.
type CostEstimator interface {
	estimateL1Cost(ctx context.Context, candidate *txmgr.TxCandidate) (*big.Int, error)
}

// isProfitable checks if proposing the given transaction batch is profitable.
// It performs profitability checks and can adjust the base fee and filter transactions if needed.
// Returns (isProfitable bool, baseFeeAdjusted bool, error)
func (p *Proposer) isProfitable(
	ctx context.Context,
	txBatch *[]types.Transactions,
	l2BaseFee **big.Int,
	candidate *txmgr.TxCandidate,
	txs *uint64,
) (bool, bool, error) {
	estimatedCost, err := p.costEstimator.estimateL1Cost(ctx, candidate)
	if err != nil {
		return false, false, fmt.Errorf("failed to estimate L1 cost: %w", err)
	}

	// Compute collected fees for the original batch/base fee
	originalCollectedFees := p.computeL2Fees(*txBatch, *l2BaseFee)
	collectedFees := originalCollectedFees

	log.Info("Profitability check (standard base fee)",
		"estimatedCost", utils.WeiToEther(estimatedCost),
		"collectedFees", utils.WeiToEther(originalCollectedFees),
		"l2BaseFee", utils.WeiToEther(*l2BaseFee),
		"numBatches", len(*txBatch),
		"numTransactions", *txs,
	)

	// Try to find a better base fee threshold
	optimalBaseFee, optimalCollectedFees := p.findOptimalBaseFeeThreshold(*txBatch)
	baseFeeAdjusted := false

	if optimalBaseFee != nil && optimalCollectedFees.Cmp(originalCollectedFees) > 0 {
		// Filter transactions with the optimal base fee
		filteredTxBatch := p.filterTxsByBaseFee(*txBatch, optimalBaseFee)
		filteredTxCount := countTxsInBatch(filteredTxBatch)

		// Calculate improvement percentage using big.Float to avoid int64 overflow
		optimalFloat := new(big.Float).SetInt(optimalCollectedFees)
		originalFloat := new(big.Float).SetInt(originalCollectedFees)
		diff := new(big.Float).Sub(optimalFloat, originalFloat)
		improvementPct, _ := new(big.Float).Quo(
			new(big.Float).Mul(diff, big.NewFloat(100)),
			originalFloat,
		).Float64()

		log.Info("Found better base fee threshold using optimized algorithm",
			"optimalBaseFee", utils.WeiToEther(optimalBaseFee),
			"optimalCollectedFees", utils.WeiToEther(optimalCollectedFees),
			"originalCollectedFees", utils.WeiToEther(originalCollectedFees),
			"filteredTxCount", filteredTxCount,
			"originalTxCount", *txs,
			"improvement", fmt.Sprintf("%.2f%%", improvementPct),
		)

		log.Info("Selecting adjusted batch/baseFee as it yields higher collected fees",
			"optimalCollectedFees", utils.WeiToEther(optimalCollectedFees),
			"estimatedCost", utils.WeiToEther(estimatedCost),
			"optimalBaseFee", utils.WeiToEther(optimalBaseFee),
			"originalCollectedFees", utils.WeiToEther(originalCollectedFees),
			"originalL2BaseFee", utils.WeiToEther(*l2BaseFee),
			"originalTxCount", *txs,
			"filteredTxCount", filteredTxCount,
		)

		// Apply the optimal configuration
		*txBatch = filteredTxBatch
		*l2BaseFee = optimalBaseFee
		*txs = filteredTxCount
		collectedFees = optimalCollectedFees
		baseFeeAdjusted = true
	}

	// Final profitability decision
	isProfitable := collectedFees.Cmp(estimatedCost) >= 0
	return isProfitable, baseFeeAdjusted, nil
}

// estimateL1Cost estimates the cost of proposing the L2 batch to L1.
// It considers both blob-based and calldata-based posting, along with proving costs.
func (p *Proposer) estimateL1Cost(
	ctx context.Context,
	candidate *txmgr.TxCandidate,
) (*big.Int, error) {
	// Fetch the latest L1 base fee
	feeHistory, err := p.rpc.L1.FeeHistory(ctx, 1, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get L1 base fee: %w", err)
	}

	if len(feeHistory.BaseFee) == 0 {
		return nil, fmt.Errorf("no base fee data available")
	}
	l1BaseFee := feeHistory.BaseFee[len(feeHistory.BaseFee)-1]

	blobBaseFee := new(big.Int)
	costWithBlobs := new(big.Int)
	costWithCalldata := new(big.Int)
	totalCost := new(big.Int)

	// If blobs are used, calculate batch posting cost with blobs
	if len(candidate.Blobs) > 0 {
		blobBaseFee, err = p.rpc.L1.BlobBaseFee(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get L1 blob base fee: %w", err)
		}

		costWithBlobs = new(big.Int).Mul(
			new(big.Int).SetUint64(p.BatchPostingGasWithBlobs),
			l1BaseFee,
		)

		costOfBlobs := new(big.Int).Mul(
			blobBaseFee,
			big.NewInt(eth.BlobSize*int64(len(candidate.Blobs))),
		)

		costWithBlobs = new(big.Int).Add(
			costWithBlobs,
			costOfBlobs,
		)
		totalCost = costWithBlobs
	} else {
		// Calculate batch posting cost with calldata
		costWithCalldata = new(big.Int).Mul(
			big.NewInt(int64(p.BatchPostingGasWithCalldata)),
			l1BaseFee,
		)
		totalCost = costWithCalldata
	}

	// Add proving and proof posting cost
	totalCost.Add(totalCost, p.ProvingCostPerL2Batch)
	proofPostingCost := new(big.Int).Mul(
		big.NewInt(int64(p.ProofPostingGas)),
		l1BaseFee,
	)
	totalCost = new(big.Int).Add(totalCost, proofPostingCost)

	log.Info("L1 cost estimation",
		"l1BaseFee", utils.WeiToEther(l1BaseFee),
		"costWithCalldata", utils.WeiToEther(costWithCalldata),
		"costWithBlobs", utils.WeiToEther(costWithBlobs),
		"blobBaseFee", utils.WeiToEther(blobBaseFee),
		"proofPostingCost", utils.WeiToEther(proofPostingCost),
		"provingCostPerL2Batch", utils.WeiToEther(p.ProvingCostPerL2Batch),
		"totalCost", utils.WeiToEther(totalCost),
	)

	return totalCost, nil
}

// computeL2Fees computes the total fees collected from a batch of transactions.
// It calculates the base fee share that goes to the proposer plus the full priority fee (tip).
func (p *Proposer) computeL2Fees(txBatch []types.Transactions, l2BaseFee *big.Int) *big.Int {
	baseFeeForProposer := p.getPercentageFromBaseFeeToTheProposer(l2BaseFee)

	collectedFees := new(big.Int)
	for _, txs := range txBatch {
		for _, tx := range txs {
			gasConsumed := big.NewInt(int64(tx.Gas()))

			// Base fee portion for proposer
			baseFeeRevenue := new(big.Int).Mul(gasConsumed, baseFeeForProposer)

			// Priority fee (tip) goes entirely to the proposer
			// For EIP-1559 txs: min(GasTipCap, GasFeeCap - baseFee)
			// For legacy txs: GasPrice - baseFee (if positive)
			priorityFee := p.calculatePriorityFee(tx, l2BaseFee)
			tipRevenue := new(big.Int).Mul(gasConsumed, priorityFee)

			// Total revenue = base fee share + tips
			txRevenue := new(big.Int).Add(baseFeeRevenue, tipRevenue)
			collectedFees.Add(collectedFees, txRevenue)
		}
	}

	return collectedFees
}

// calculatePriorityFee calculates the priority fee (tip) for a transaction.
// Returns the tip per gas that the proposer receives.
func (p *Proposer) calculatePriorityFee(tx *types.Transaction, baseFee *big.Int) *big.Int {
	// For dynamic fee transactions (EIP-1559)
	if tx.GasTipCap() != nil && tx.GasFeeCap() != nil {
		// Effective tip = min(maxPriorityFeePerGas, maxFeePerGas - baseFee)
		maxTip := tx.GasTipCap()
		maxFeeMinusBase := new(big.Int).Sub(tx.GasFeeCap(), baseFee)

		// Return the minimum of the two (or 0 if negative)
		if maxFeeMinusBase.Sign() <= 0 {
			return big.NewInt(0)
		}
		if maxTip.Cmp(maxFeeMinusBase) < 0 {
			return new(big.Int).Set(maxTip)
		}
		return maxFeeMinusBase
	}

	// For legacy transactions
	if tx.GasPrice() != nil {
		tip := new(big.Int).Sub(tx.GasPrice(), baseFee)
		if tip.Sign() > 0 {
			return tip
		}
	}

	return big.NewInt(0)
}

// getPercentageFromBaseFeeToTheProposer calculates what percentage of the base fee goes to the proposer.
// Uses the protocol configuration's SharingPctg to determine the split.
func (p *Proposer) getPercentageFromBaseFeeToTheProposer(num *big.Int) *big.Int {
	if p.protocolConfigs.BaseFeeConfig().SharingPctg == 0 {
		return big.NewInt(0)
	}

	result := new(big.Int).Mul(num, big.NewInt(int64(p.protocolConfigs.BaseFeeConfig().SharingPctg)))
	return new(big.Int).Div(result, big.NewInt(100))
}

// txWithBaseFee is a helper struct to sort transactions by their base fee
type txWithBaseFee struct {
	tx      *types.Transaction
	baseFee *big.Int
	gas     uint64
}

// findOptimalBaseFeeThreshold finds the optimal base fee threshold that maximizes collected fees.
// This algorithm has complexity:
// - O(n) to flatten and extract base fees
// - O(n log n) to sort transactions by base fee
// - O(n^2) to find the optimal threshold
// Total: O(n^2) where n is the number of transactions
// Returns the optimal base fee and the collected fees at that threshold.
func (p *Proposer) findOptimalBaseFeeThreshold(txBatch []types.Transactions) (*big.Int, *big.Int) {
	// Flatten all transactions into a single slice with their base fees
	var allTxs []txWithBaseFee
	for _, txs := range txBatch {
		for _, tx := range txs {
			txBaseFee := tx.GasFeeCap()
			if txBaseFee == nil {
				txBaseFee = tx.GasPrice()
			}
			if txBaseFee != nil {
				allTxs = append(allTxs, txWithBaseFee{
					tx:      tx,
					baseFee: new(big.Int).Set(txBaseFee),
					gas:     tx.Gas(),
				})
			}
		}
	}

	if len(allTxs) == 0 {
		return nil, big.NewInt(0)
	}

	// Sort transactions by base fee in descending order (highest first)
	sort.Slice(allTxs, func(i, j int) bool {
		return allTxs[i].baseFee.Cmp(allTxs[j].baseFee) > 0
	})

	baseFeeSharePct := p.protocolConfigs.BaseFeeConfig().SharingPctg
	var bestBaseFee *big.Int
	bestCollectedFees := big.NewInt(0)
	cumulativeGas := uint64(0)

	// Now iterate through sorted transactions to find optimal threshold
	// Key insight: if we set base fee to tx[i].baseFee, we include tx[0..i]
	// Collected fees = (baseFee * proposerSharePct + avgTip) * sum(gas[0..i])
	for i := 0; i < len(allTxs); i++ {
		cumulativeGas += allTxs[i].gas
		// Calculate collected fees if we use this base fee as threshold
		// All transactions from 0 to i have baseFee >= currentBaseFee
		currentBaseFee := allTxs[i].baseFee

		// Calculate total revenue from base fees (proposer share)
		baseFeeForProposer := new(big.Int).Mul(
			currentBaseFee,
			big.NewInt(int64(baseFeeSharePct)),
		)
		baseFeeForProposer.Div(baseFeeForProposer, big.NewInt(100))
		baseFeeRevenue := new(big.Int).Mul(
			baseFeeForProposer,
			new(big.Int).SetUint64(cumulativeGas),
		)

		// Calculate total revenue from tips for all included transactions
		tipRevenue := big.NewInt(0)
		for j := 0; j <= i; j++ {
			txTip := p.calculatePriorityFee(allTxs[j].tx, currentBaseFee)
			txTipRevenue := new(big.Int).Mul(txTip, new(big.Int).SetUint64(allTxs[j].gas))
			tipRevenue.Add(tipRevenue, txTipRevenue)
		}

		// Total collected fees = base fee revenue + tip revenue
		collectedFees := new(big.Int).Add(baseFeeRevenue, tipRevenue)

		if collectedFees.Cmp(bestCollectedFees) > 0 {
			bestCollectedFees.Set(collectedFees)
			bestBaseFee = new(big.Int).Set(currentBaseFee)
		}

		// Log progress periodically
		if i%1000 == 0 || i == len(allTxs)-1 {
			log.Debug("Searching for optimal base fee",
				"progress", fmt.Sprintf("%d/%d", i+1, len(allTxs)),
				"currentBaseFee", utils.WeiToEther(currentBaseFee),
				"cumulativeGas", cumulativeGas,
				"baseFeeRevenue", utils.WeiToEther(baseFeeRevenue),
				"tipRevenue", utils.WeiToEther(tipRevenue),
				"collectedFees", utils.WeiToEther(collectedFees),
				"bestSoFar", utils.WeiToEther(bestCollectedFees),
			)
		}
	}

	return bestBaseFee, bestCollectedFees
}

// filterTxsByBaseFee filters transactions that have a GasFeeCap >= the specified base fee.
// For legacy transactions, it uses GasPrice instead of GasFeeCap.
// Returns the filtered batch maintaining the batch structure but removing empty batches.
func (p *Proposer) filterTxsByBaseFee(txBatch []types.Transactions, minBaseFee *big.Int) []types.Transactions {
	filteredBatch := make([]types.Transactions, 0, len(txBatch))

	for _, txs := range txBatch {
		filteredTxs := make(types.Transactions, 0, len(txs))

		for _, tx := range txs {
			// Get the GasFeeCap which represents the maximum base fee the transaction is willing to pay
			txBaseFee := tx.GasFeeCap()
			if txBaseFee == nil {
				// For legacy transactions, use GasPrice
				txBaseFee = tx.GasPrice()
			}

			// Include transaction if it meets the minimum base fee
			if txBaseFee != nil && txBaseFee.Cmp(minBaseFee) >= 0 {
				filteredTxs = append(filteredTxs, tx)
			}
		}

		// Only add non-empty transaction lists to the batch
		if len(filteredTxs) > 0 {
			filteredBatch = append(filteredBatch, filteredTxs)
		}
	}

	return filteredBatch
}

// countTxsInBatch counts the total number of transactions in a batch.
func countTxsInBatch(txBatch []types.Transactions) uint64 {
	var count uint64
	for _, txs := range txBatch {
		count += uint64(len(txs))
	}
	return count
}
