package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type QueueProposalRequestBody struct {
	TxDest common.Address `json:"txDest"`
	TxData hexutil.Bytes  `json:"txData"`
	TxList hexutil.Bytes  `json:"txList"`
}
