package builder

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/config"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// CelestiaTransactionBuilder is responsible for building a TaikoInbox.proposeBatch transaction with txList bytes saved in Celestia.
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
) (*txmgr.TxCandidate, error) {
	return nil, errors.New("not implemented")
}
