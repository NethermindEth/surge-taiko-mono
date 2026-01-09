package eventiterator

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	shastaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/shasta"
	chainIterator "github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/chain_iterator"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// EndProposalsRollbackedEventIterFunc ends the current iteration.
type EndProposalsRollbackedEventIterFunc func()

// OnProposalsRollbackedEvent represents the callback function which will be called when a RollbackInbox.Rollbacked
// event is iterated.
type OnProposalsRollbackedEvent func(
	context.Context,
	*shastaBindings.RollbackInboxRollbacked,
	EndProposalsRollbackedEventIterFunc,
) error

// ProposalsRollbackedIterator iterates the emitted RollbackInbox.Rollbacked events in the chain,
// with the awareness of reorganization.
type ProposalsRollbackedIterator struct {
	ctx                context.Context
	rollbackInbox      *shastaBindings.RollbackInbox
	blockBatchIterator *chainIterator.BlockBatchIterator
	isEnd              bool
}

// ProposalsRollbackedIteratorConfig represents the configs of a ProposalsRollbacked event iterator.
type ProposalsRollbackedIteratorConfig struct {
	Client                     *rpc.EthClient
	RollbackInbox              *shastaBindings.RollbackInbox
	MaxBlocksReadPerEpoch      *uint64
	StartHeight                *big.Int
	EndHeight                  *big.Int
	OnProposalsRollbackedEvent OnProposalsRollbackedEvent
	BlockConfirmations         *uint64
}

// NewProposalsRollbackedIterator creates a new instance of ProposalsRollbacked event iterator.
func NewProposalsRollbackedIterator(
	ctx context.Context,
	cfg *ProposalsRollbackedIteratorConfig,
) (*ProposalsRollbackedIterator, error) {
	if cfg.OnProposalsRollbackedEvent == nil {
		return nil, errors.New("invalid callback")
	}

	iterator := &ProposalsRollbackedIterator{ctx: ctx, rollbackInbox: cfg.RollbackInbox}

	// Initialize the inner block iterator.
	blockIterator, err := chainIterator.NewBlockBatchIterator(ctx, &chainIterator.BlockBatchIteratorConfig{
		Client:                cfg.Client,
		MaxBlocksReadPerEpoch: cfg.MaxBlocksReadPerEpoch,
		StartHeight:           cfg.StartHeight,
		EndHeight:             cfg.EndHeight,
		BlockConfirmations:    cfg.BlockConfirmations,
		OnBlocks: assembleProposalsRollbackedIteratorCallback(
			cfg.Client,
			cfg.RollbackInbox,
			cfg.OnProposalsRollbackedEvent,
			iterator,
		),
	})
	if err != nil {
		return nil, err
	}

	iterator.blockBatchIterator = blockIterator

	return iterator, nil
}

// Iter iterates the given chain between the given start and end heights,
// will call the callback when a ProposalsRollbacked event is iterated.
func (i *ProposalsRollbackedIterator) Iter() error {
	return i.blockBatchIterator.Iter()
}

// end ends the current iteration.
func (i *ProposalsRollbackedIterator) end() {
	i.isEnd = true
}

// assembleProposalsRollbackedIteratorCallback assembles the callback which will be used
// by a event iterator's inner block iterator.
func assembleProposalsRollbackedIteratorCallback(
	client *rpc.EthClient,
	rollbackInbox *shastaBindings.RollbackInbox,
	callback OnProposalsRollbackedEvent,
	eventIter *ProposalsRollbackedIterator,
) chainIterator.OnBlocksFunc {
	return func(
		ctx context.Context,
		start, end *types.Header,
		updateCurrentFunc chainIterator.UpdateCurrentFunc,
		endFunc chainIterator.EndIterFunc,
	) error {
		var (
			endHeight = end.Number.Uint64()
		)

		// Iterate the Rollbacked events.
		iterRollbackInbox, err := rollbackInbox.FilterRollbacked(
			&bind.FilterOpts{Start: start.Number.Uint64(), End: &endHeight, Context: ctx},
		)
		if err != nil {
			log.Error("Failed to filter Rollbacked events", "error", err)
			return err
		}
		defer iterRollbackInbox.Close()

		for iterRollbackInbox.Next() {
			event := iterRollbackInbox.Event
			log.Debug("Processing Rollbacked event",
				"start proposal id", event.FirstProposalId,
				"end proposal id", event.LastProposalId,
				"l1BlockHeight", event.Raw.BlockNumber,
			)

			if err := callback(ctx, event, eventIter.end); err != nil {
				log.Warn("Error while processing Rollbacked events, keep retrying", "error", err)
				return err
			}

			if eventIter.isEnd {
				log.Debug("RollbackedIterator is ended", "start", start.Number, "end", endHeight)
				endFunc()
				return nil
			}

			current, err := client.HeaderByHash(ctx, event.Raw.BlockHash)
			if err != nil {
				return err
			}

			log.Debug("Updating current block cursor for processing Rollbacked events", "block", current.Number)

			updateCurrentFunc(current)
		}

		// Check if there is any error during the iteration.
		if iterRollbackInbox.Error() != nil {
			log.Error("Error while iterating Rollbacked events", "error", iterRollbackInbox.Error())
			return iterRollbackInbox.Error()
		}

		return nil
	}
}
