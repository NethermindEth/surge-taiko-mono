package surge

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// L1Origin represents the L1 origin information for a block.
type L1Origin struct {
	BlockID            *big.Int    `json:"blockID" rlp:"optional"`
	L2BlockHash        common.Hash `json:"l2BlockHash" rlp:"optional"`
	L1BlockHeight      *big.Int    `json:"l1BlockHeight" rlp:"optional"`
	L1BlockHash        common.Hash `json:"l1BlockHash" rlp:"optional"`
	BatchID            *big.Int    `json:"batchID" rlp:"optional"`
	BuildPayloadArgsID [8]byte     `json:"buildPayloadArgsID" rlp:"optional"`
}
