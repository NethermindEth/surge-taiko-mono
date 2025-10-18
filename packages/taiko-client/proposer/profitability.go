package proposer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/utils"
)

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
	estimatedCost, err := p.estimateL2Cost(ctx, candidate)
	if err != nil {
		return false, false, fmt.Errorf("failed to estimate L2 cost: %w", err)
	}

	// First, check profitability with the standard L2 base fee
	collectedFees := p.computeL2Fees(*txBatch, *l2BaseFee)
	isProfitable := collectedFees.Cmp(estimatedCost) >= 0

	log.Info("Profitability check (standard base fee)",
		"estimatedCost", utils.WeiToEther(estimatedCost),
		"collectedFees", utils.WeiToEther(collectedFees),
		"isProfitable", isProfitable,
		"l2BaseFee", utils.WeiToEther(*l2BaseFee),
		"numBatches", len(*txBatch),
		"numTransactions", *txs,
	)

	// If profitable with standard base fee, return early
	if isProfitable {
		return true, false, nil
	}

	// If not profitable with standard base fee, try with adjusted percentages of highest transaction base fee
	highestTxBaseFee := p.findHighestBaseFeeInBatch(*txBatch)

	if highestTxBaseFee == nil || highestTxBaseFee.Cmp(*l2BaseFee) <= 0 {
		// No higher base fee transactions found, return not profitable
		return false, false, nil
	}

	// Try different percentage thresholds: 50%, then 75%
	percentages := []int64{50, 75}

	for _, percentage := range percentages {
		// Calculate the adjusted base fee
		adjustedBaseFee := new(big.Int).Mul(highestTxBaseFee, big.NewInt(percentage))
		adjustedBaseFee = new(big.Int).Div(adjustedBaseFee, big.NewInt(100))

		// Filter transactions that meet the adjusted base fee
		filteredTxBatch := p.filterTxsByBaseFee(*txBatch, adjustedBaseFee)

		if len(filteredTxBatch) == 0 {
			log.Info("No transactions meet the adjusted base fee threshold",
				"percentage", percentage,
				"adjustedBaseFee", utils.WeiToEther(adjustedBaseFee),
			)
			continue
		}

		// Recalculate collected fees with filtered transactions and adjusted base fee
		collectedFeesAdjusted := p.computeL2Fees(filteredTxBatch, adjustedBaseFee)
		isProfitableAdjusted := collectedFeesAdjusted.Cmp(estimatedCost) >= 0

		log.Info("Profitability check (adjusted with highest tx base fee)",
			"estimatedCost", utils.WeiToEther(estimatedCost),
			"collectedFeesAdjusted", utils.WeiToEther(collectedFeesAdjusted),
			"isProfitableAdjusted", isProfitableAdjusted,
			"highestTxBaseFee", utils.WeiToEther(highestTxBaseFee),
			"percentage", percentage,
			"adjustedBaseFee", utils.WeiToEther(adjustedBaseFee),
			"originalTxCount", *txs,
			"filteredTxCount", countTxsInBatch(filteredTxBatch),
		)

		if isProfitableAdjusted {
			log.Info("Proposing with adjusted base fee due to high-paying transactions",
				"percentage", percentage,
			)
			// Modify the references in-place
			*txBatch = filteredTxBatch
			*l2BaseFee = adjustedBaseFee
			*txs = countTxsInBatch(filteredTxBatch)
			return true, true, nil
		}
	}

	return false, false, nil
}

// estimateL2Cost estimates the cost of proposing the L2 batch to L1.
// It considers both blob-based and calldata-based posting, along with proving costs.
func (p *Proposer) estimateL2Cost(
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

	log.Info("L2 cost estimation",
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
// It calculates the base fee share that goes to the proposer.
func (p *Proposer) computeL2Fees(txBatch []types.Transactions, l2BaseFee *big.Int) *big.Int {
	baseFeeForProposer := p.getPercentageFromBaseFeeToTheProposer(l2BaseFee)

	collectedFees := new(big.Int)
	for _, txs := range txBatch {
		for _, tx := range txs {
			gasConsumed := big.NewInt(int64(tx.Gas()))
			expectedFee := new(big.Int).Mul(gasConsumed, baseFeeForProposer)
			collectedFees.Add(collectedFees, expectedFee)
		}
	}

	return collectedFees
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

// findHighestBaseFeeInBatch finds the highest base fee (GasFeeCap) from all transactions in the batch.
// For legacy transactions, it uses GasPrice instead of GasFeeCap.
func (p *Proposer) findHighestBaseFeeInBatch(txBatch []types.Transactions) *big.Int {
	var highestBaseFee *big.Int

	for _, txs := range txBatch {
		for _, tx := range txs {
			// Get the GasFeeCap which represents the maximum base fee the transaction is willing to pay
			txBaseFee := tx.GasFeeCap()
			if txBaseFee == nil {
				// For legacy transactions, use GasPrice
				txBaseFee = tx.GasPrice()
			}

			if txBaseFee != nil {
				if highestBaseFee == nil || txBaseFee.Cmp(highestBaseFee) > 0 {
					highestBaseFee = new(big.Int).Set(txBaseFee)
				}
			}
		}
	}

	return highestBaseFee
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
