package prover

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/suite"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/driver"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/metrics"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/testutils"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/jwt"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/proposer"
	proofProducer "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_producer"
)

type ProverTestSuite struct {
	testutils.ClientTestSuite
	p        *Prover
	cancel   context.CancelFunc
	d        *driver.Driver
	proposer *proposer.Proposer
	txmgr    *txmgr.SimpleTxManager
}

func (s *ProverTestSuite) SetupTest() {
	s.ClientTestSuite.SetupTest()

	// Init prover
	var (
		l1ProverPrivKey = s.KeyFromEnv("L1_PROVER_PRIVATE_KEY")
		err             error
	)

	s.txmgr, err = txmgr.NewSimpleTxManager(
		"prover_test",
		log.Root(),
		&metrics.TxMgrMetrics,
		txmgr.CLIConfig{
			L1RPCURL:                  os.Getenv("L1_WS"),
			NumConfirmations:          0,
			SafeAbortNonceTooLowCount: txmgr.DefaultBatcherFlagValues.SafeAbortNonceTooLowCount,
			PrivateKey:                common.Bytes2Hex(crypto.FromECDSA(l1ProverPrivKey)),
			FeeLimitMultiplier:        txmgr.DefaultBatcherFlagValues.FeeLimitMultiplier,
			FeeLimitThresholdGwei:     txmgr.DefaultBatcherFlagValues.FeeLimitThresholdGwei,
			MinBaseFeeGwei:            txmgr.DefaultBatcherFlagValues.MinBaseFeeGwei,
			MinTipCapGwei:             txmgr.DefaultBatcherFlagValues.MinTipCapGwei,
			ResubmissionTimeout:       txmgr.DefaultBatcherFlagValues.ResubmissionTimeout,
			ReceiptQueryInterval:      1 * time.Second,
			NetworkTimeout:            txmgr.DefaultBatcherFlagValues.NetworkTimeout,
			TxSendTimeout:             txmgr.DefaultBatcherFlagValues.TxSendTimeout,
			TxNotInMempoolTimeout:     txmgr.DefaultBatcherFlagValues.TxNotInMempoolTimeout,
		},
	)
	s.Nil(err)

	ctx, cancel := context.WithCancel(context.Background())
	s.initProver(ctx, l1ProverPrivKey)
	s.cancel = cancel

	// Init driver
	jwtSecret, err := jwt.ParseSecretFromFile(os.Getenv("JWT_SECRET"))
	s.Nil(err)
	s.NotEmpty(jwtSecret)

	d := new(driver.Driver)
	s.Nil(d.InitFromConfig(context.Background(), &driver.Config{
		ClientConfig: &rpc.ClientConfig{
			L1Endpoint:         os.Getenv("L1_WS"),
			L2Endpoint:         os.Getenv("L2_WS"),
			L2EngineEndpoint:   os.Getenv("L2_AUTH"),
			TaikoInboxAddress:  common.HexToAddress(os.Getenv("TAIKO_INBOX")),
			TaikoAnchorAddress: common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
			JwtSecret:          string(jwtSecret),
		},
		BlobServerEndpoint: s.BlobServer.URL(),
	}))
	s.d = d

	// Init proposer
	var (
		l1ProposerPrivKey = s.KeyFromEnv("L1_PROVER_PRIVATE_KEY")
		prop              = new(proposer.Proposer)
	)

	s.Nil(prop.InitFromConfig(context.Background(), &proposer.Config{
		ClientConfig: &rpc.ClientConfig{
			L1Endpoint:                  os.Getenv("L1_WS"),
			L2Endpoint:                  os.Getenv("L2_WS"),
			L2EngineEndpoint:            os.Getenv("L2_AUTH"),
			JwtSecret:                   string(jwtSecret),
			TaikoInboxAddress:           common.HexToAddress(os.Getenv("TAIKO_INBOX")),
			TaikoWrapperAddress:         common.HexToAddress(os.Getenv("TAIKO_WRAPPER")),
			ForcedInclusionStoreAddress: common.HexToAddress(os.Getenv("FORCED_INCLUSION_STORE")),
			ProverSetAddress:            common.HexToAddress(os.Getenv("PROVER_SET")),
			TaikoAnchorAddress:          common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
			TaikoTokenAddress:           common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		},
		L1ProposerPrivKey:       l1ProposerPrivKey,
		L2SuggestedFeeRecipient: common.HexToAddress(os.Getenv("L2_SUGGESTED_FEE_RECIPIENT")),
		ProposeInterval:         1024 * time.Hour,
		MaxTxListsPerEpoch:      1,
	}, s.txmgr, s.txmgr))

	s.proposer = prop
	s.proposer.RegisterTxMgrSelectorToBlobServer(s.BlobServer)
}

func (s *ProverTestSuite) TestName() {
	s.Equal("prover", s.p.Name())
}

func (s *ProverTestSuite) TestInitError() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var (
		l1ProverPrivKey = s.KeyFromEnv("L1_PROVER_PRIVATE_KEY")
		p               = new(Prover)
	)

	s.NotNil(InitFromConfig(ctx, p, &Config{
		L1WsEndpoint:          os.Getenv("L1_WS"),
		L2WsEndpoint:          os.Getenv("L2_WS"),
		L2HttpEndpoint:        os.Getenv("L2_HTTP"),
		TaikoInboxAddress:     common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoAnchorAddress:    common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress:     common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		L1ProverPrivKey:       l1ProverPrivKey,
		Dummy:                 true,
		ProveUnassignedBlocks: true,
		RPCTimeout:            10 * time.Minute,
		BackOffRetryInterval:  3 * time.Second,
		BackOffMaxRetries:     12,
	}, s.txmgr, s.txmgr))
}

func (s *ProverTestSuite) TestSubmitProofAggregationOp() {
	s.NotPanics(func() {
		s.p.withRetry(func() error {
			return s.p.submitProofAggregationOp(&proofProducer.BatchProofs{
				ProofResponses: []*proofProducer.ProofResponse{
					{
						BatchID: common.Big1,
						Meta:    &metadata.TaikoDataBlockMetadataPacaya{},
						Proof:   []byte{},
						Opts:    &proofProducer.ProofRequestOptionsPacaya{},
					},
				},
				BatchProof:    []byte{},
				BatchIDs:      []*big.Int{common.Big1},
				ProofType:     proofProducer.ProofTypeOp,
				SgxBatchProof: []byte{},
			})
		})
	})
}

func (s *ProverTestSuite) TestOnBatchesVerified() {
	s.NotPanics(func() {
		s.NotNil(s.p.eventHandlers.batchesVerifiedHandler.HandlePacaya(
			context.Background(),
			&pacayaBindings.TaikoInboxClientBatchesVerified{
				BatchId: testutils.RandomHash().Big().Uint64(),
				Raw: types.Log{
					BlockHash:   testutils.RandomHash(),
					BlockNumber: testutils.RandomHash().Big().Uint64(),
				},
			}))
	})
}

func (s *ProverTestSuite) TestSetApprovalAlreadySetHigher() {
	s.p.cfg.Allowance = common.Big256
	s.Nil(s.p.setApprovalAmount(context.Background(), s.p.cfg.TaikoInboxAddress))

	originalAllowance, err := s.p.rpc.PacayaClients.TaikoToken.Allowance(
		nil,
		s.p.ProverAddress(),
		s.p.cfg.TaikoInboxAddress,
	)
	s.Nil(err)
	s.NotZero(originalAllowance.Uint64())

	s.p.cfg.Allowance = new(big.Int).Sub(originalAllowance, common.Big1)

	s.Nil(s.p.setApprovalAmount(context.Background(), s.p.cfg.TaikoInboxAddress))

	allowance, err := s.p.rpc.PacayaClients.TaikoToken.Allowance(nil, s.p.ProverAddress(), s.p.cfg.TaikoInboxAddress)
	s.Nil(err)

	s.Zero(allowance.Cmp(originalAllowance))
}

func (s *ProverTestSuite) TearDownTest() {
	if s.p.ctx.Err() == nil {
		s.cancel()
	}
	s.p.Close(context.Background())
}

func TestProverTestSuite(t *testing.T) {
	suite.Run(t, new(ProverTestSuite))
}

func (s *ProverTestSuite) initProver(ctx context.Context, key *ecdsa.PrivateKey) {
	decimal, err := s.RPCClient.PacayaClients.TaikoToken.Decimals(nil)
	s.Nil(err)

	p := new(Prover)
	s.Nil(InitFromConfig(ctx, p, &Config{
		L1WsEndpoint:          os.Getenv("L1_WS"),
		L2WsEndpoint:          os.Getenv("L2_WS"),
		L2HttpEndpoint:        os.Getenv("L2_HTTP"),
		TaikoInboxAddress:     common.HexToAddress(os.Getenv("TAIKO_INBOX")),
		TaikoAnchorAddress:    common.HexToAddress(os.Getenv("TAIKO_ANCHOR")),
		TaikoTokenAddress:     common.HexToAddress(os.Getenv("TAIKO_TOKEN")),
		ProverSetAddress:      common.HexToAddress(os.Getenv("PROVER_SET")),
		L1ProverPrivKey:       key,
		Dummy:                 true,
		ProveUnassignedBlocks: true,
		Allowance:             new(big.Int).Exp(big.NewInt(1_000_000_100), new(big.Int).SetUint64(uint64(decimal)), nil),
		RPCTimeout:            3 * time.Second,
		BackOffRetryInterval:  3 * time.Second,
		BackOffMaxRetries:     12,
		SGXProofBufferSize:    1,
		ZKVMProofBufferSize:   1,
		BlockConfirmations:    0,
	}, s.txmgr, s.txmgr))

	s.p = p
}
