package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
)

type QueueProposalRequestBody struct {
	Inbox                    common.Address                  `json:"inbox"`
	Coinbase                 common.Address                  `json:"coinbase"`
	RevertIfNotFirstProposal bool                            `json:"revertIfNotFirstProposal"`
	Blocks                   []pacaya.ITaikoInboxBlockParams `json:"blocks"`
	TxList                   hexutil.Bytes                   `json:"txList"`
	ForcedInclusionParams    encoding.BatchParams            `json:"forcedInclusionParams"`
}
