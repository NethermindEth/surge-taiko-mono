package eventiterator

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	chainIterator "github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/chain_iterator"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

// EndBatchProposedEventIterFunc ends the current iteration.
type EndBatchProposedEventIterFunc func()

// OnBatchProposedEvent represents the callback function which will be called
// when a Pacaya or Shasta proposal event is iterated.
type OnBatchProposedEvent func(
	context.Context,
	metadata.TaikoProposalMetaData,
	EndBatchProposedEventIterFunc,
) error

// BatchProposedIterator iterates the emitted Pacaya or Shasta proposal events in the chain,
// with the awareness of reorganization.
type BatchProposedIterator struct {
	blockBatchIterator *chainIterator.BlockBatchIterator
	isEnd              bool
}

// BatchProposedIteratorConfig represents the configs of a BatchProposed event iterator.
type BatchProposedIteratorConfig struct {
	RpcClient             *rpc.Client
	MaxBlocksReadPerEpoch *uint64
	StartHeight           *big.Int
	EndHeight             *big.Int
	OnBatchProposedEvent  OnBatchProposedEvent
	BlockConfirmations    *uint64
	Fork                  string
}

// NewBatchProposedIterator creates a new instance of BatchProposed event iterator.
func NewBatchProposedIterator(ctx context.Context, cfg *BatchProposedIteratorConfig) (*BatchProposedIterator, error) {
	if cfg.RpcClient == nil || cfg.RpcClient.L1 == nil {
		return nil, errors.New("invalid RPC client")
	}

	if cfg.OnBatchProposedEvent == nil {
		return nil, errors.New("invalid callback")
	}

	iterator := new(BatchProposedIterator)

	// Initialize the inner block iterator.
	blockIterator, err := chainIterator.NewBlockBatchIterator(ctx, &chainIterator.BlockBatchIteratorConfig{
		Client:                cfg.RpcClient.L1,
		MaxBlocksReadPerEpoch: cfg.MaxBlocksReadPerEpoch,
		StartHeight:           cfg.StartHeight,
		EndHeight:             cfg.EndHeight,
		BlockConfirmations:    cfg.BlockConfirmations,
		OnBlocks: assembleBatchProposedIteratorCallback(
			cfg.RpcClient,
			cfg.OnBatchProposedEvent,
			iterator,
			cfg.Fork,
		),
	})
	if err != nil {
		return nil, err
	}

	iterator.blockBatchIterator = blockIterator

	return iterator, nil
}

// Iter iterates the given chain between the given start and end heights,
// will call the callback when a BatchProposed event is iterated.
func (i *BatchProposedIterator) Iter() error {
	return i.blockBatchIterator.Iter()
}

// end ends the current iteration.
func (i *BatchProposedIterator) end() {
	i.isEnd = true
}

// assembleBatchProposedIteratorCallback assembles the callback which will be used
// by a event iterator's inner block iterator. Only the event loop for the active
// fork is executed.
func assembleBatchProposedIteratorCallback(
	rpcClient *rpc.Client,
	callback OnBatchProposedEvent,
	eventIter *BatchProposedIterator,
	fork string,
) chainIterator.OnBlocksFunc {
	return func(
		ctx context.Context,
		start, end *types.Header,
		updateCurrentFunc chainIterator.UpdateCurrentFunc,
		endFunc chainIterator.EndIterFunc,
	) error {
		endHeight := end.Number.Uint64()

		switch fork {
		case "pacaya":
			return iteratePacayaEvents(ctx, rpcClient, callback, eventIter, start, endHeight, updateCurrentFunc, endFunc)
		case "shasta":
			return iterateShastaEvents(ctx, rpcClient, callback, eventIter, start, endHeight, updateCurrentFunc, endFunc)
		case "realtime":
			return iterateRealTimeEvents(ctx, rpcClient, callback, eventIter, start, endHeight, updateCurrentFunc, endFunc)
		default:
			return fmt.Errorf("unknown fork %q", fork)
		}
	}
}

// iteratePacayaEvents filters and processes Pacaya BatchProposed events.
func iteratePacayaEvents(
	ctx context.Context,
	rpcClient *rpc.Client,
	callback OnBatchProposedEvent,
	eventIter *BatchProposedIterator,
	start *types.Header,
	endHeight uint64,
	updateCurrentFunc chainIterator.UpdateCurrentFunc,
	endFunc chainIterator.EndIterFunc,
) error {
	iter, err := rpcClient.PacayaClients.TaikoInbox.FilterBatchProposed(
		&bind.FilterOpts{Start: start.Number.Uint64(), End: &endHeight, Context: ctx},
	)
	if err != nil {
		return err
	}
	defer iter.Close()

	var lastBatchID uint64
	for iter.Next() {
		event := iter.Event
		log.Debug("Processing BatchProposed event", "batch", event.Meta.BatchId, "l1BlockHeight", event.Raw.BlockNumber)

		if lastBatchID != 0 && event.Meta.BatchId != lastBatchID+1 {
			return fmt.Errorf(
				"BatchProposed event is not continuous, lastScannedBatchID: %d, currentScannedBatchID: %d",
				lastBatchID, event.Meta.BatchId,
			)
		}

		if err := callback(ctx, metadata.NewTaikoDataBlockMetadataPacaya(event), eventIter.end); err != nil {
			return err
		}

		if eventIter.isEnd {
			endFunc()
			return nil
		}

		current, err := rpcClient.L1.HeaderByHash(ctx, event.Raw.BlockHash)
		if err != nil {
			return err
		}

		lastBatchID = event.Meta.BatchId
		updateCurrentFunc(current)
	}

	return iter.Error()
}

// iterateShastaEvents filters and processes Shasta Proposed events.
func iterateShastaEvents(
	ctx context.Context,
	rpcClient *rpc.Client,
	callback OnBatchProposedEvent,
	eventIter *BatchProposedIterator,
	start *types.Header,
	endHeight uint64,
	updateCurrentFunc chainIterator.UpdateCurrentFunc,
	endFunc chainIterator.EndIterFunc,
) error {
	iter, err := rpcClient.ShastaClients.Inbox.FilterProposed(
		&bind.FilterOpts{Start: start.Number.Uint64(), End: &endHeight, Context: ctx}, nil, nil,
	)
	if err != nil {
		return err
	}
	defer iter.Close()

	var lastProposalID uint64
	for iter.Next() {
		event := iter.Event

		header, err := rpcClient.L1.HeaderByHash(ctx, event.Raw.BlockHash)
		if err != nil {
			return fmt.Errorf("failed to fetch L1 block header: %w", err)
		}

		proposedEventPayload := metadata.NewTaikoProposalMetadataShasta(event, header.Time)
		proposalID := proposedEventPayload.Shasta().GetEventData().Id.Uint64()
		log.Debug("Processing Proposed event", "proposalID", proposalID, "l1BlockHeight", event.Raw.BlockNumber)

		if lastProposalID != 0 && proposalID != lastProposalID+1 {
			return fmt.Errorf(
				"Proposed event is not continuous, lastScannedProposalID: %d, currentScannedProposalID: %d",
				lastProposalID, proposalID,
			)
		}

		if err := callback(ctx, proposedEventPayload, eventIter.end); err != nil {
			return err
		}

		if eventIter.isEnd {
			endFunc()
			return nil
		}

		current, err := rpcClient.L1.HeaderByHash(ctx, event.Raw.BlockHash)
		if err != nil {
			return err
		}

		lastProposalID = proposalID
		updateCurrentFunc(current)
	}

	return iter.Error()
}

// iterateRealTimeEvents filters and processes RealTime ProposedAndProved events.
func iterateRealTimeEvents(
	ctx context.Context,
	rpcClient *rpc.Client,
	callback OnBatchProposedEvent,
	eventIter *BatchProposedIterator,
	start *types.Header,
	endHeight uint64,
	updateCurrentFunc chainIterator.UpdateCurrentFunc,
	endFunc chainIterator.EndIterFunc,
) error {
	iter, err := rpcClient.RealTimeClients.Inbox.FilterProposedAndProved(
		&bind.FilterOpts{Start: start.Number.Uint64(), End: &endHeight, Context: ctx},
		nil,
	)
	if err != nil {
		return err
	}
	defer iter.Close()

	var lastProposalHash [32]byte
	for iter.Next() {
		event := iter.Event

		header, err := rpcClient.L1.HeaderByHash(ctx, event.Raw.BlockHash)
		if err != nil {
			return fmt.Errorf("failed to fetch L1 block header: %w", err)
		}

		proposalMetadata := metadata.NewTaikoProposalMetadataRealTime(event, header.Time)
		log.Debug("Processing ProposedAndProved event",
			"proposalHash", common.Hash(event.ProposalHash),
			"parentProposalHash", common.Hash(event.ParentProposalHash),
			"l1BlockHeight", event.Raw.BlockNumber,
		)

		if lastProposalHash != ([32]byte{}) && event.ParentProposalHash != lastProposalHash {
			return fmt.Errorf(
				"ProposedAndProved event hash chain is not continuous, lastHash: %s, parentHash: %s",
				common.Hash(lastProposalHash),
				common.Hash(event.ParentProposalHash),
			)
		}

		if err := callback(ctx, proposalMetadata, eventIter.end); err != nil {
			return err
		}

		if eventIter.isEnd {
			endFunc()
			return nil
		}

		current, err := rpcClient.L1.HeaderByHash(ctx, event.Raw.BlockHash)
		if err != nil {
			return err
		}

		lastProposalHash = event.ProposalHash
		updateCurrentFunc(current)
	}

	return iter.Error()
}
