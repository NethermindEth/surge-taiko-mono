package prover

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/metrics"
	eventIterator "github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/chain_iterator/event_iterator"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	handler "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/event_handler"
	guardianProverHeartbeater "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/guardian_prover_heartbeater"
	proofProducer "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_producer"
	proofSubmitter "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_submitter"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_submitter/transaction"
	state "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/shared_state"
)

// Prover keeps trying to prove newly proposed blocks.
type Prover struct {
	// Configurations
	cfg     *Config
	backoff backoff.BackOffContext

	// Clients
	rpc *rpc.Client

	// Guardian prover related
	guardianProverHeartbeater guardianProverHeartbeater.BlockSenderHeartbeater

	// Contract configurations
	protocolConfig *bindings.TaikoDataConfig

	// States
	sharedState *state.SharedState

	// Event handlers
	blockProposedHandler       handler.BlockProposedHandler
	blockVerifiedHandler       handler.BlockVerifiedHandler
	transitionContestedHandler handler.TransitionContestedHandler
	transitionProvedHandler    handler.TransitionProvedHandler
	assignmentExpiredHandler   handler.AssignmentExpiredHandler

	// Proof submitters
	proofSubmitters []proofSubmitter.Submitter
	proofContester  proofSubmitter.Contester

	assignmentExpiredCh chan metadata.TaikoBlockMetaData
	proveNotify         chan struct{}

	// Proof related channels
	proofSubmissionCh chan *proofProducer.ProofRequestBody
	proofContestCh    chan *proofProducer.ContestRequestBody
	proofGenerationCh chan *proofProducer.ProofWithHeader

	// Transactions manager
	txmgr        *txmgr.SimpleTxManager
	privateTxmgr *txmgr.SimpleTxManager

	ctx context.Context
	wg  sync.WaitGroup
}

// InitFromCli initializes the given prover instance based on the command line flags.
func (p *Prover) InitFromCli(ctx context.Context, c *cli.Context) error {
	cfg, err := NewConfigFromCliContext(c)
	if err != nil {
		return err
	}

	return InitFromConfig(ctx, p, cfg, nil, nil)
}

// InitFromConfig initializes the prover instance based on the given configurations.
func InitFromConfig(
	ctx context.Context,
	p *Prover, cfg *Config,
	txMgr *txmgr.SimpleTxManager,
	privateTxMgr *txmgr.SimpleTxManager,
) (err error) {
	log.Debug("Initializing prover from config")
	p.cfg = cfg
	p.ctx = ctx
	// Initialize state which will be shared by event handlers.
	p.sharedState = state.New()
	p.backoff = backoff.WithContext(
		backoff.WithMaxRetries(
			backoff.NewConstantBackOff(p.cfg.BackOffRetryInterval),
			p.cfg.BackOffMaxRetries,
		),
		p.ctx,
	)

	// Clients
	log.Debug("Initializing RPC client",
		"l1Endpoint", cfg.L1WsEndpoint,
		"l2Endpoint", cfg.L2WsEndpoint)
	if p.rpc, err = rpc.NewClient(p.ctx, &rpc.ClientConfig{
		L1Endpoint:                    cfg.L1WsEndpoint,
		L2Endpoint:                    cfg.L2WsEndpoint,
		TaikoL1Address:                cfg.TaikoL1Address,
		TaikoL2Address:                cfg.TaikoL2Address,
		TaikoTokenAddress:             cfg.TaikoTokenAddress,
		ProverSetAddress:              cfg.ProverSetAddress,
		GuardianProverMinorityAddress: cfg.GuardianProverMinorityAddress,
		GuardianProverMajorityAddress: cfg.GuardianProverMajorityAddress,
		Timeout:                       cfg.RPCTimeout,
	}); err != nil {
		return err
	}

	// Configs
	log.Debug("Getting protocol config")
	p.protocolConfig = encoding.GetProtocolConfig(p.rpc.L2.ChainID.Uint64())
	log.Info("Protocol configs", "configs", p.protocolConfig)

	log.Debug("Initializing channels", "bufferSize", p.protocolConfig.BlockMaxProposals)
	chBufferSize := p.protocolConfig.BlockMaxProposals
	p.proofGenerationCh = make(chan *proofProducer.ProofWithHeader, chBufferSize)
	p.assignmentExpiredCh = make(chan metadata.TaikoBlockMetaData, chBufferSize)
	p.proofSubmissionCh = make(chan *proofProducer.ProofRequestBody, p.cfg.Capacity)
	p.proofContestCh = make(chan *proofProducer.ContestRequestBody, p.cfg.Capacity)
	p.proveNotify = make(chan struct{}, 1)

	log.Debug("Initializing L1 current cursor", "startingBlockID", cfg.StartingBlockID)
	if err := p.initL1Current(cfg.StartingBlockID); err != nil {
		return fmt.Errorf("initialize L1 current cursor error: %w", err)
	}

	// Protocol proof tiers
	log.Debug("Getting protocol tiers")
	tiers, err := p.rpc.GetTiers(ctx)
	if err != nil {
		return err
	}
	p.sharedState.SetTiers(tiers)
	log.Debug("Protocol tiers set", "tiers", tiers)

	log.Debug("Initializing transaction builder")
	txBuilder := transaction.NewProveBlockTxBuilder(
		p.rpc,
		p.cfg.TaikoL1Address,
		p.cfg.ProverSetAddress,
		p.cfg.GuardianProverMajorityAddress,
		p.cfg.GuardianProverMinorityAddress,
	)

	log.Debug("Setting up transaction managers")
	if txMgr != nil {
		p.txmgr = txMgr
		log.Debug("Using provided transaction manager")
	} else {
		log.Debug("Creating new transaction manager")
		if p.txmgr, err = txmgr.NewSimpleTxManager(
			"prover",
			log.Root(),
			&metrics.TxMgrMetrics,
			*cfg.TxmgrConfigs,
		); err != nil {
			return err
		}
	}

	if privateTxMgr != nil {
		log.Debug("Using provided private transaction manager")
		p.privateTxmgr = privateTxMgr
	} else {
		log.Debug("Checking private transaction manager config")
		if cfg.PrivateTxmgrConfigs != nil && len(cfg.PrivateTxmgrConfigs.L1RPCURL) > 0 {
			log.Debug("Creating new private transaction manager")
			if p.privateTxmgr, err = txmgr.NewSimpleTxManager(
				"privateMempoolProver",
				log.Root(),
				&metrics.TxMgrMetrics,
				*cfg.PrivateTxmgrConfigs,
			); err != nil {
				return err
			}
		} else {
			log.Debug("No private transaction manager configured")
			p.privateTxmgr = nil
		}
	}

	// Proof submitters
	log.Debug("Initializing proof submitters")
	if err := p.initProofSubmitters(txBuilder, tiers); err != nil {
		return err
	}

	// Proof contester
	log.Debug("Initializing proof contester")
	p.proofContester = proofSubmitter.NewProofContester(
		p.rpc,
		p.cfg.ProveBlockGasLimit,
		p.txmgr,
		p.privateTxmgr,
		p.cfg.ProverSetAddress,
		p.cfg.Graffiti,
		txBuilder,
	)

	// Initialize event handlers.
	log.Debug("Initializing event handlers")
	if err := p.initEventHandlers(); err != nil {
		return err
	}

	log.Debug("Prover initialization completed successfully")
	return nil
}

// Start starts the main loop of the L2 block prover.
func (p *Prover) Start() error {
	log.Info("Starting prover")
	// 1. Set approval amount for the contracts.
	for _, contract := range []common.Address{p.cfg.TaikoL1Address} {
		if err := p.setApprovalAmount(p.ctx, contract); err != nil {
			log.Crit("Failed to set approval amount", "contract", contract, "error", err)
		}
	}

	// 2. Start the main event loop of the prover.
	go p.eventLoop()

	return nil
}

// eventLoop starts the main loop of Taiko prover.
func (p *Prover) eventLoop() {
	p.wg.Add(1)
	defer p.wg.Done()

	// reqProving requests performing a proving operation, won't block
	// if we are already proving.
	reqProving := func() {
		select {
		case p.proveNotify <- struct{}{}:
		default:
		}
	}
	// Call reqProving() right away to catch up with the latest state.
	reqProving()

	// If there is too many (TaikoData.Config.blockMaxProposals) pending blocks in TaikoL1 contract, there will be no new
	// BlockProposed event temporarily, so except the BlockProposed subscription, we need another trigger to start
	// fetching the proposed blocks.
	forceProvingTicker := time.NewTicker(15 * time.Second)
	defer forceProvingTicker.Stop()

	// Channels
	chBufferSize := p.protocolConfig.BlockMaxProposals
	blockProposedCh := make(chan *bindings.TaikoL1ClientBlockProposed, chBufferSize)
	blockVerifiedCh := make(chan *bindings.TaikoL1ClientBlockVerified, chBufferSize)
	transitionProvedCh := make(chan *bindings.TaikoL1ClientTransitionProved, chBufferSize)
	transitionContestedCh := make(chan *bindings.TaikoL1ClientTransitionContested, chBufferSize)
	blockProposedV2Ch := make(chan *bindings.TaikoL1ClientBlockProposedV2, chBufferSize)
	blockVerifiedV2Ch := make(chan *bindings.TaikoL1ClientBlockVerifiedV2, chBufferSize)
	transitionProvedV2Ch := make(chan *bindings.TaikoL1ClientTransitionProvedV2, chBufferSize)
	transitionContestedV2Ch := make(chan *bindings.TaikoL1ClientTransitionContestedV2, chBufferSize)
	// Subscriptions
	blockProposedSub := rpc.SubscribeBlockProposed(p.rpc.TaikoL1, blockProposedCh)
	blockVerifiedSub := rpc.SubscribeBlockVerified(p.rpc.TaikoL1, blockVerifiedCh)
	transitionProvedSub := rpc.SubscribeTransitionProved(p.rpc.TaikoL1, transitionProvedCh)
	transitionContestedSub := rpc.SubscribeTransitionContested(p.rpc.TaikoL1, transitionContestedCh)
	blockProposedV2Sub := rpc.SubscribeBlockProposedV2(p.rpc.TaikoL1, blockProposedV2Ch)
	blockVerifiedV2Sub := rpc.SubscribeBlockVerifiedV2(p.rpc.TaikoL1, blockVerifiedV2Ch)
	transitionProvedV2Sub := rpc.SubscribeTransitionProvedV2(p.rpc.TaikoL1, transitionProvedV2Ch)
	transitionContestedV2Sub := rpc.SubscribeTransitionContestedV2(p.rpc.TaikoL1, transitionContestedV2Ch)
	defer func() {
		blockProposedSub.Unsubscribe()
		blockVerifiedSub.Unsubscribe()
		transitionProvedSub.Unsubscribe()
		transitionContestedSub.Unsubscribe()
		blockProposedV2Sub.Unsubscribe()
		blockVerifiedV2Sub.Unsubscribe()
		transitionProvedV2Sub.Unsubscribe()
		transitionContestedV2Sub.Unsubscribe()
	}()

	for {
		select {
		case <-p.ctx.Done():
			return
		case req := <-p.proofContestCh:
			p.withRetry(func() error { return p.contestProofOp(req) })
		case proofWithHeader := <-p.proofGenerationCh:
			p.withRetry(func() error { return p.submitProofOp(proofWithHeader) })
		case req := <-p.proofSubmissionCh:
			p.withRetry(func() error { return p.requestProofOp(req.Meta, req.Tier) })
		case <-p.proveNotify:
			if err := p.proveOp(); err != nil {
				log.Error("Prove new blocks error", "error", err)
			}
		case e := <-blockVerifiedCh:
			p.blockVerifiedHandler.Handle(encoding.BlockVerifiedEventToV2(e))
		case e := <-transitionProvedCh:
			p.withRetry(func() error {
				blockInfo, err := p.rpc.GetL2BlockInfo(p.ctx, e.BlockId)
				if err != nil {
					return err
				}
				return p.transitionProvedHandler.Handle(p.ctx, encoding.TransitionProvedEventToV2(e, blockInfo.ProposedIn))
			})
		case e := <-transitionContestedCh:
			p.withRetry(func() error {
				blockInfo, err := p.rpc.GetL2BlockInfo(p.ctx, e.BlockId)
				if err != nil {
					return err
				}
				return p.transitionContestedHandler.Handle(
					p.ctx,
					encoding.TransitionContestedEventToV2(e, blockInfo.ProposedIn),
				)
			})
		case e := <-blockVerifiedV2Ch:
			p.blockVerifiedHandler.Handle(e)
		case e := <-transitionProvedV2Ch:
			p.withRetry(func() error {
				return p.transitionProvedHandler.Handle(p.ctx, e)
			})
		case e := <-transitionContestedV2Ch:
			p.withRetry(func() error {
				return p.transitionContestedHandler.Handle(p.ctx, e)
			})
		case m := <-p.assignmentExpiredCh:
			p.withRetry(func() error { return p.assignmentExpiredHandler.Handle(p.ctx, m) })
		case <-blockProposedCh:
			reqProving()
		case <-blockProposedV2Ch:
			reqProving()
		case <-forceProvingTicker.C:
			reqProving()
		}
	}
}

// Close closes the prover instance.
func (p *Prover) Close(_ context.Context) {
	p.wg.Wait()
}

// proveOp iterates through BlockProposed events.
func (p *Prover) proveOp() error {
	iter, err := eventIterator.NewBlockProposedIterator(p.ctx, &eventIterator.BlockProposedIteratorConfig{
		Client:               p.rpc.L1,
		TaikoL1:              p.rpc.TaikoL1,
		StartHeight:          new(big.Int).SetUint64(p.sharedState.GetL1Current().Number.Uint64()),
		OnBlockProposedEvent: p.blockProposedHandler.Handle,
		BlockConfirmations:   &p.cfg.BlockConfirmations,
	})
	if err != nil {
		log.Error("Failed to start event iterator", "event", "BlockProposed", "error", err)
		return err
	}

	return iter.Iter()
}

// contestProofOp performs a proof contest operation.
func (p *Prover) contestProofOp(req *proofProducer.ContestRequestBody) error {
	if err := p.proofContester.SubmitContest(
		p.ctx,
		req.BlockID,
		req.ProposedIn,
		req.ParentHash,
		req.Meta,
		req.Tier,
	); err != nil {
		if strings.Contains(err.Error(), vm.ErrExecutionReverted.Error()) {
			log.Error(
				"Proof contest submission reverted",
				"blockID", req.BlockID,
				"minTier", req.Meta.GetMinTier(),
				"error", err,
			)
			return nil
		}
		log.Error(
			"Request new proof contest error",
			"blockID", req.BlockID,
			"minTier", req.Meta.GetMinTier(),
			"error", err,
		)
		return err
	}

	return nil
}

// requestProofOp requests a new proof generation operation.
func (p *Prover) requestProofOp(meta metadata.TaikoBlockMetaData, minTier uint16) error {
	if submitter := p.selectSubmitter(minTier); submitter != nil {
		if err := submitter.RequestProof(p.ctx, meta); err != nil {
			log.Error("Request new proof error", "blockID", meta.GetBlockID(), "minTier", meta.GetMinTier(), "error", err)
			return err
		}

		return nil
	}

	log.Error("Failed to find proof submitter", "blockID", meta.GetBlockID(), "minTier", minTier)
	return nil
}

// submitProofOp performs a proof submission operation.
func (p *Prover) submitProofOp(proofWithHeader *proofProducer.ProofWithHeader) error {
	submitter := p.getSubmitterByTier(proofWithHeader.Tier)
	if submitter == nil {
		return nil
	}

	if err := submitter.SubmitProof(p.ctx, proofWithHeader); err != nil {
		if strings.Contains(err.Error(), vm.ErrExecutionReverted.Error()) {
			log.Error(
				"Proof submission reverted",
				"blockID", proofWithHeader.BlockID,
				"minTier", proofWithHeader.Meta.GetMinTier(),
				"error", err,
			)
			return nil
		}
		log.Error(
			"Submit proof error",
			"blockID", proofWithHeader.BlockID,
			"minTier", proofWithHeader.Meta.GetMinTier(),
			"error", err,
		)
		return err
	}

	return nil
}

// Name returns the application name.
func (p *Prover) Name() string {
	return "prover"
}

// selectSubmitter returns the proof submitter with the given minTier.
func (p *Prover) selectSubmitter(minTier uint16) proofSubmitter.Submitter {
	for _, s := range p.proofSubmitters {
		if s.Tier() >= minTier {
			log.Debug("Proof submitter selected", "tier", s.Tier(), "minTier", minTier)
			return s
		}
	}

	log.Warn("No proof producer / submitter found for the given minTier", "minTier", minTier)

	return nil
}

// getSubmitterByTier returns the proof submitter with the given tier.
func (p *Prover) getSubmitterByTier(tier uint16) proofSubmitter.Submitter {
	for _, s := range p.proofSubmitters {
		if s.Tier() == tier {
			return s
		}
	}

	log.Warn("No proof producer / submitter found for the given tier", "tier", tier)

	return nil
}

// IsGuardianProver returns true if the current prover is a guardian prover.
func (p *Prover) IsGuardianProver() bool {
	return p.cfg.GuardianProverMajorityAddress != common.Address{}
}

// ProverAddress returns the current prover account address.
func (p *Prover) ProverAddress() common.Address {
	return p.txmgr.From()
}

// withRetry retries the given function with prover backoff policy.
func (p *Prover) withRetry(f func() error) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := backoff.Retry(f, p.backoff); err != nil {
			log.Error("Operation failed", "error", err)
		}
	}()
}
