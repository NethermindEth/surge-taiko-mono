package surge

import (
	"math/big"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// This is an extended version of engine.PayloadAttributes that uses surge.L1Origin instead of rawdb.L1Origin.
type PayloadAttributes struct {
	Timestamp             uint64
	Random                common.Hash
	SuggestedFeeRecipient common.Address
	Withdrawals           []*types.Withdrawal
	BlockMetadata         *engine.BlockMetadata
	BaseFeePerGas         *big.Int
	L1Origin              *L1Origin
}
