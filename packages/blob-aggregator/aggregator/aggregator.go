package aggregator

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	txmgrMetrics "github.com/ethereum-optimism/optimism/op-service/txmgr/metrics"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/bindings"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue"
	taikoEncoding "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/urfave/cli/v2"
)

type Aggregator struct {
	aggregatorPrivateKey *ecdsa.PrivateKey
	txmgr                txmgr.TxManager

	queue queue.Queue

	batch              []*queue.Message
	batchTxListSize    uint64
	minTxBatchDataSize uint64

	proposalCh chan queue.Message

	minBlobFillupPercentage uint64

	wg  sync.WaitGroup
	ctx context.Context
}

func (agg *Aggregator) InitFromCli(ctx context.Context, c *cli.Context) error {
	cfg, err := NewConfigFromCliContext(c)
	if err != nil {
		return err
	}

	return InitFromConfig(ctx, agg, cfg)
}

func InitFromConfig(ctx context.Context, agg *Aggregator, cfg *Config) (err error) {
	q, err := cfg.OpenQueueFunc()
	if err != nil {
		return err
	}

	agg.aggregatorPrivateKey = cfg.L1AggregatorPrivKey
	agg.txmgr, err = txmgr.NewSimpleTxManager("aggregator", log.Root(), new(txmgrMetrics.NoopTxMetrics), *cfg.TxMgrConfig)
	if err != nil {
		return err
	}

	agg.queue = q
	agg.batch = []*queue.Message{}
	agg.minTxBatchDataSize = eth.BlobSize * cfg.MinAggregatedBlobs
	agg.minBlobFillupPercentage = cfg.MinBlobsFillupPercentage // Todo: this is currently unused
	agg.proposalCh = make(chan queue.Message)
	agg.ctx = ctx

	return nil
}

func (agg *Aggregator) Name() string {
	return "aggregator"
}

func (agg *Aggregator) Close(ctx context.Context) {
	agg.wg.Wait()
}

func (agg *Aggregator) Start() error {
	go func() {
		if err := backoff.Retry(func() error {
			slog.Info("attempting backoff queue subscription")
			if err := agg.queue.Subscribe(agg.ctx, agg.proposalCh, &agg.wg); err != nil {
				slog.Error("aggregator queue subscription error", "err", err.Error())
				return err
			}
			return nil
		}, backoff.WithContext(backoff.NewConstantBackOff(1*time.Second), agg.ctx)); err != nil {
			slog.Error("queue subscribe backoff retry error", "err", err.Error())
		}
	}()

	go agg.eventLoop()

	return nil
}

func (agg *Aggregator) eventLoop() {
	agg.wg.Add(1)
	defer agg.wg.Done()

	for {
		select {
		case <-agg.ctx.Done():
			return
		case msg := <-agg.proposalCh:
			err := agg.addProposalToBatch(&msg)
			if err != nil {
				slog.Error("error processing proposal", "error", err.Error())
			}
		}
	}
}

func (agg *Aggregator) addProposalToBatch(msg *queue.Message) error {
	slog.Info("Received proposal", "proposal", msg.Proposal)

	agg.batch = append(agg.batch, msg)
	agg.batchTxListSize += uint64(len(msg.Proposal.TxList))

	if agg.batchTxListSize >= agg.minTxBatchDataSize {
		if err := agg.processBatch(); err != nil {
			return err
		}

		agg.batch = []*queue.Message{}
		agg.batchTxListSize = 0
	}

	return nil
}

func (agg *Aggregator) processBatch() error {
	blobParamsList := agg.buildBlobParamsForBatchedProposal()

	// Function parameter for batched proposal via the Minimal Batcher
	calls, err := agg.buildParamsForBatchedProposal(blobParamsList)
	if err != nil {
		return err
	}

	// Full calldata for the batched proposal
	executeBatchCalldata, err := agg.buildExecuteBatchCalldata(calls)
	if err != nil {
		return err
	}

	// Shared blob for the batched proposal
	blobs, err := agg.buildBlobs(agg.batch)
	if err != nil {
		return err
	}

	to := crypto.PubkeyToAddress(agg.aggregatorPrivateKey.PublicKey)
	candidate := txmgr.TxCandidate{
		TxData:   executeBatchCalldata,
		Blobs:    blobs,
		To:       &to,
		GasLimit: 0,
	}

	receipt, err := agg.txmgr.Send(agg.ctx, candidate)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		slog.Warn("Transaction reverted", "txHash", hex.EncodeToString(receipt.TxHash.Bytes()))
		return errors.New("transaction reverted")
	}

	return nil
}

func (agg *Aggregator) buildParamsForBatchedProposal(blobParamsList []*taikoEncoding.BlobParams) ([]bindings.Call, error) {
	var (
		encodedProposalParams []byte
		proposalCalldata      []byte
		calls                 []bindings.Call
		err                   error
	)

	for i := 0; i < len(agg.batch); i++ {
		if encodedProposalParams, err = taikoEncoding.EncodeBatchParamsWithForcedInclusion(
			&agg.batch[i].Proposal.ForcedInclusionParams,
			&taikoEncoding.BatchParams{
				Proposer:                 crypto.PubkeyToAddress(agg.aggregatorPrivateKey.PublicKey),
				Coinbase:                 agg.batch[i].Proposal.Coinbase,
				RevertIfNotFirstProposal: agg.batch[i].Proposal.RevertIfNotFirstProposal,
				BlobParams:               *blobParamsList[i],
				Blocks:                   agg.batch[i].Proposal.Blocks,
			}); err != nil {
			return nil, err
		}

		if proposalCalldata, err = taikoEncoding.TaikoWrapperABI.Pack("proposeBatch", encodedProposalParams, []byte{}); err != nil {
			return nil, err
		}

		calls = append(calls, bindings.Call{
			Target: agg.batch[i].Proposal.Inbox,
			Value:  big.NewInt(0),
			Data:   proposalCalldata,
		})
	}

	return calls, nil
}

func (agg *Aggregator) buildExecuteBatchCalldata(calls []bindings.Call) ([]byte, error) {
	executeBatchParams, err := encoding.EncodeExecuteBatchInput(calls)
	if err != nil {
		return nil, err
	}

	executeBatchCalldata, err := encoding.MinimalBatcherABI.Pack("executeBatch", executeBatchParams)
	if err != nil {
		return nil, err
	}

	return executeBatchCalldata, nil
}
