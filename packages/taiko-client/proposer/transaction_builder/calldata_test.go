package builder

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
)

func (s *TransactionBuilderTestSuite) TestBuildCalldata() {
	_, err := s.calldataTxBuilder.BuildOntake(context.Background(), [][]byte{{1}, {2}})
	s.Nil(err)
}
