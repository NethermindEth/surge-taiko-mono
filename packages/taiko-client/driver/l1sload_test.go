package driver

import (
	"context"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/testutils"
)

const l2NodeNMC = "l2_nmc"

var l1SLOADPrecompileAddr = common.HexToAddress("0x0000000000000000000000000000000000010001")

// buildL1SLOADCalldata constructs the 84-byte input for the L1SLOAD precompile.
// Format: targetAddress(20) || storageKey(32) || l1BlockNumber(32)
func buildL1SLOADCalldata(addr common.Address, slot common.Hash, blockNum *big.Int) []byte {
	calldata := make([]byte, 0, 84)
	calldata = append(calldata, addr.Bytes()...)
	calldata = append(calldata, slot.Bytes()...)
	calldata = append(calldata, common.BigToHash(blockNum).Bytes()...)
	return calldata
}

// TestL1SLOADStaticCall verifies that eth_call to the L1SLOAD precompile on L2
// returns the correct storage value from L1 (Anvil).
func (s *DriverTestSuite) TestL1SLOADStaticCall() {
	if os.Getenv("L2_NODE") != l2NodeNMC {
		s.T().Skip("L1SLOAD only supported on NMC")
	}

	targetAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	slot := common.BigToHash(big.NewInt(1))
	value := common.BigToHash(big.NewInt(0xdeadbeef))

	// Plant a known storage value on L1 and mine the block.
	s.SetL1Storage(targetAddr, slot, value)

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Advance the L2 chain so it has an L1 origin >= our storage block.
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// eth_call the L1SLOAD precompile on L2.
	calldata := buildL1SLOADCalldata(targetAddr, slot, l1Head.Number)

	result, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1SLOADPrecompileAddr,
		Data: calldata,
		Gas:  100_000,
	}, nil)
	s.Nil(err)
	s.Equal(value, common.BytesToHash(result))
}

// TestL1SLOADTransaction sends an actual L2 transaction to the L1SLOAD precompile,
// proposes a block, and verifies the tx succeeded on-chain.
func (s *DriverTestSuite) TestL1SLOADTransaction() {
	if os.Getenv("L2_NODE") != l2NodeNMC {
		s.T().Skip("L1SLOAD only supported on NMC")
	}

	targetAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	slot := common.BigToHash(big.NewInt(2))
	value := common.BigToHash(big.NewInt(0xcafebabe))

	s.SetL1Storage(targetAddr, slot, value)

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Send an actual L2 tx calling the precompile.
	calldata := buildL1SLOADCalldata(targetAddr, slot, l1Head.Number)
	_, err = testutils.SendDynamicFeeTx(
		s.RPCClient.L2,
		s.TestAddrPrivKey,
		&l1SLOADPrecompileAddr,
		common.Big0,
		calldata,
	)
	s.Nil(err)

	// Propose and insert a block containing the L1SLOAD tx.
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// Retrieve the L1SLOAD tx from the proposed block.
	// We use TransactionInBlock rather than the client-side hash because the
	// proposer re-encodes txs via blob tx lists, which can change the hash.
	l2Head, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.Greater(l2Head.Number.Uint64(), uint64(1))

	txCount, err := s.RPCClient.L2.TransactionCount(context.Background(), l2Head.Hash())
	s.Nil(err)
	s.GreaterOrEqual(txCount, uint(2), "block should have anchor tx + L1SLOAD tx")

	// Index 0 is the anchor tx; our L1SLOAD tx is at index 1.
	userTx, err := s.RPCClient.L2.TransactionInBlock(context.Background(), l2Head.Hash(), 1)
	s.Nil(err)

	receipt, err := s.RPCClient.L2.TransactionReceipt(context.Background(), userTx.Hash())
	s.Nil(err)
	s.Equal(types.ReceiptStatusSuccessful, receipt.Status)
	s.T().Logf("L1SLOAD tx gas used: %d", receipt.GasUsed)
}

// TestL1SLOADInvalidBlockRejected verifies that requesting a non-existent (future)
// L1 block via the precompile fails.
func (s *DriverTestSuite) TestL1SLOADInvalidBlockRejected() {
	if os.Getenv("L2_NODE") != l2NodeNMC {
		s.T().Skip("L1SLOAD only supported on NMC")
	}

	targetAddr := common.HexToAddress("0x3333333333333333333333333333333333333333")
	slot := common.BigToHash(big.NewInt(3))
	value := common.BigToHash(big.NewInt(0xbaadf00d))

	s.SetL1Storage(targetAddr, slot, value)

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Request a block far in the future that doesn't exist on L1.
	futureBlock := new(big.Int).Add(l1Head.Number, big.NewInt(1000))
	calldata := buildL1SLOADCalldata(targetAddr, slot, futureBlock)

	_, err = s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1SLOADPrecompileAddr,
		Data: calldata,
		Gas:  100_000,
	}, nil)
	s.NotNil(err, "Expected error when querying future L1 block")
}
