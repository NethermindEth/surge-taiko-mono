package transaction

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/testutils"
	producer "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_producer"
)

func (s *TransactionTestSuite) TestBuildTxs() {
	header, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.NotNil(header)

	builder := s.builder.BuildProveBatchesPacaya(&producer.BatchProofs{
		ProofResponses: []*producer.ProofResponse{{
			BatchID:    common.Big1,
			Meta:       metadata.NewTaikoDataBlockMetadataPacaya(&pacayaBindings.TaikoInboxClientBatchProposed{}),
			Proof1:     testutils.RandomBytes(100),
			ProofType1: producer.ProofTypeZKR0,
			Proof2:     testutils.RandomBytes(100),
			ProofType2: producer.ProofTypeZKSP1,
			Opts:       &producer.ProofRequestOptionsPacaya{Headers: []*types.Header{header}},
		}},
	})
	_, err = builder(&bind.TransactOpts{})
	s.Nil(err)
}
