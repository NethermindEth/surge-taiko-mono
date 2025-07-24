package builder

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/celestia"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/utils"
)

// CelestiaTransactionBuilder is responsible for building a TaikoInbox.proposeBatch transaction with txList
// bytes saved in Celestia.
type CelestiaTransactionBuilder struct {
	rpc                     *rpc.Client
	proposerPrivateKey      *ecdsa.PrivateKey
	taikoInboxAddress       common.Address
	taikoWrapperAddress     common.Address
	proverSetAddress        common.Address
	l2SuggestedFeeRecipient common.Address
	gasLimit                uint64
	chainConfig             *config.ChainConfig
	revertProtectionEnabled bool
}

// NewCelestiaTransactionBuilder creates a new CelestiaTransactionBuilder instance based on giving configurations.
func NewCelestiaTransactionBuilder(
	rpc *rpc.Client,
	proposerPrivateKey *ecdsa.PrivateKey,
	taikoInboxAddress common.Address,
	taikoWrapperAddress common.Address,
	proverSetAddress common.Address,
	l2SuggestedFeeRecipient common.Address,
	gasLimit uint64,
	chainConfig *config.ChainConfig,
	revertProtectionEnabled bool,
) *CelestiaTransactionBuilder {
	return &CelestiaTransactionBuilder{
		rpc,
		proposerPrivateKey,
		taikoInboxAddress,
		taikoWrapperAddress,
		proverSetAddress,
		l2SuggestedFeeRecipient,
		gasLimit,
		chainConfig,
		revertProtectionEnabled,
	}
}

// BuildPacaya implements the ProposeBatchTransactionBuilder interface.
func (b *CelestiaTransactionBuilder) BuildPacaya(
	ctx context.Context,
	txBatch []types.Transactions,
	forcedInclusion *pacayaBindings.IForcedInclusionStoreForcedInclusion,
	minTxsPerForcedInclusion *big.Int,
	parentMetahash common.Hash,
	baseFee *big.Int,
) (*txmgr.TxCandidate, error) {
	// ABI encode the TaikoWrapper.proposeBatch / ProverSet.proposeBatch parameters.
	var (
		to                    = &b.taikoWrapperAddress
		proposer              = crypto.PubkeyToAddress(b.proposerPrivateKey.PublicKey)
		data                  []byte
		encodedParams         []byte
		blockParams           []pacayaBindings.ITaikoInboxBlockParams
		forcedInclusionParams *encoding.BatchParams
		allTxs                types.Transactions
	)

	if b.proverSetAddress != rpc.ZeroAddress {
		to = &b.proverSetAddress
		proposer = b.proverSetAddress
	}

	if forcedInclusion != nil {
		blobParams, blockParams := buildParamsForForcedInclusion(forcedInclusion, minTxsPerForcedInclusion)
		forcedInclusionParams = &encoding.BatchParams{
			Proposer:                 proposer,
			Coinbase:                 b.l2SuggestedFeeRecipient,
			RevertIfNotFirstProposal: b.revertProtectionEnabled,
			BlobParams:               *blobParams,
			Blocks:                   blockParams,
			BaseFee:                  baseFee,
		}
	}

	for _, txs := range txBatch {
		allTxs = append(allTxs, txs...)
		blockParams = append(blockParams, pacayaBindings.ITaikoInboxBlockParams{
			NumTransactions: uint16(len(txs)),
			TimeShift:       0,
			SignalSlots:     make([][32]byte, 0),
		})
	}

	txListsBytes, err := utils.EncodeAndCompressTxList(allTxs)
	if err != nil {
		return nil, err
	}

	celestiaBlobs, err := b.splitToCelestiaBlobs(txListsBytes)
	if err != nil {
		return nil, err
	}

	celestiaHeight, err := b.rpc.CelestiaDA.Submit(ctx, celestiaBlobs)
	if err != nil {
		return nil, err
	}

	params := &encoding.BatchParams{
		Proposer:                 proposer,
		Coinbase:                 b.l2SuggestedFeeRecipient,
		RevertIfNotFirstProposal: b.revertProtectionEnabled,
		BlobParams: encoding.BlobParams{
			ByteOffset: 0,
			ByteSize:   uint32(len(txListsBytes)),
		},
		CelestiaBlobParams: encoding.CelestiaBlobParams{
			Height:    celestiaHeight,
			Namespace: b.rpc.CelestiaDA.Namespace.Bytes(),
		},
		Blocks:  blockParams,
		BaseFee: baseFee,
	}

	if b.revertProtectionEnabled {
		if forcedInclusionParams != nil {
			forcedInclusionParams.ParentMetaHash = parentMetahash
		} else {
			params.ParentMetaHash = parentMetahash
		}
	}

	if encodedParams, err = encoding.EncodeBatchParamsWithForcedInclusion(forcedInclusionParams, params); err != nil {
		return nil, err
	}

	if b.proverSetAddress != rpc.ZeroAddress {
		if data, err = encoding.ProverSetPacayaABI.Pack("proposeBatch", encodedParams, []byte{}); err != nil {
			return nil, err
		}
	} else {
		if data, err = encoding.TaikoWrapperABI.Pack("proposeBatch", encodedParams, []byte{}); err != nil {
			return nil, err
		}
	}

	return &txmgr.TxCandidate{
		TxData:   data,
		Blobs:    nil,
		To:       to,
		GasLimit: b.gasLimit,
	}, nil
}

// splitToCelestiaBlobs splits the txListBytes into multiple Celestia blobs.
func (b *CelestiaTransactionBuilder) splitToCelestiaBlobs(txListBytes []byte) ([]*celestia.Blob, error) {
	var blobs []*celestia.Blob
	for start := 0; start < len(txListBytes); start += rpc.AdvisableCelestiaBlobSize {
		end := start + rpc.AdvisableCelestiaBlobSize
		if end > len(txListBytes) {
			end = len(txListBytes)
		}

		blob, err := celestia.NewBlobV0(b.rpc.CelestiaDA.Namespace, txListBytes[start:end])
		if err != nil {
			return nil, err
		}

		blobs = append(blobs, blob)
	}

	return blobs, nil
}
