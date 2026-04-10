package driver

import (
	"context"
	"math/big"
	"os"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/testutils"
)

var l1STATICCALLPrecompileAddr = common.HexToAddress("0x0000000000000000000000000000000000010002")

// l1TestContractAddr is the address where we deploy a minimal test contract on L1.
var l1TestContractAddr = common.HexToAddress("0x4444444444444444444444444444444444444444")

// minimalViewContractCode is the bytecode of a contract that:
//   - On any call, returns the 32-byte value stored at storage slot 0.
//
// Solidity equivalent: function fallback() { return storage[0]; }
// Assembly: PUSH0 SLOAD PUSH0 MSTORE PUSH1 0x20 PUSH0 RETURN
// Hex:      5F545F5260205FF3
var minimalViewContractCode = common.FromHex("0x5F545F5260205FF3")

// dynamicSlotViewContractCode is the bytecode of a contract that:
//   - Reads the first 32 bytes of calldata as a storage slot key
//   - Returns the 32-byte value stored at that slot
//
// Assembly: PUSH0 CALLDATALOAD SLOAD PUSH0 MSTORE PUSH1 0x20 PUSH0 RETURN
// Hex:      5F35545F5260205FF3
var dynamicSlotViewContractCode = common.FromHex("0x5F35545F5260205FF3")

// l1DynamicContractAddr is the address for the dynamic-slot view contract.
var l1DynamicContractAddr = common.HexToAddress("0x6666666666666666666666666666666666666666")

// buildL1STATICCALLCalldata constructs the variable-length input for the L1STATICCALL precompile.
// Format: targetAddress(20) || l1BlockNumber(32) || calldata(variable)
func buildL1STATICCALLCalldata(target common.Address, blockNum *big.Int, fnCalldata []byte) []byte {
	data := make([]byte, 0, 20+32+len(fnCalldata))
	data = append(data, target.Bytes()...)
	data = append(data, common.BigToHash(blockNum).Bytes()...)
	data = append(data, fnCalldata...)
	return data
}

// setupL1ViewContract deploys a minimal view contract on L1 at l1TestContractAddr
// that returns the value from storage slot 0. Uses Anvil's anvil_setCode and anvil_setStorageAt.
func (s *DriverTestSuite) setupL1ViewContract(storageValue common.Hash) {
	// Deploy bytecode at the test address.
	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setCode", l1TestContractAddr, common.Bytes2Hex(minimalViewContractCode),
	))
	// Set storage slot 0 to the desired value.
	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setStorageAt",
		l1TestContractAddr,
		common.Hash{}, // slot 0
		storageValue,
	))
	s.L1Mine()
}

// TestL1STATICCALLViewFunction verifies that eth_call to the L1STATICCALL precompile on L2
// returns the correct result from an L1 view function call.
func (s *DriverTestSuite) TestL1STATICCALLViewFunction() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	expectedValue := common.BigToHash(big.NewInt(0xdeadbeef))

	// Deploy a minimal view contract on L1 that returns storage slot 0.
	s.setupL1ViewContract(expectedValue)

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Advance the L2 chain so it has an L1 origin >= our L1 block.
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// Call L1STATICCALL precompile on L2 with empty function calldata
	// (triggers the fallback which returns slot 0).
	calldata := buildL1STATICCALLCalldata(l1TestContractAddr, l1Head.Number, nil)

	result, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata,
		Gas:  200_000,
	}, nil)
	s.Nil(err)
	s.Equal(expectedValue, common.BytesToHash(result),
		"L1STATICCALL should return the value from the L1 contract's storage slot 0")
}

// TestL1STATICCALLTransaction sends an actual L2 transaction to the L1STATICCALL precompile,
// proposes a block, and verifies the tx succeeded on-chain.
func (s *DriverTestSuite) TestL1STATICCALLTransaction() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	expectedValue := common.BigToHash(big.NewInt(0xcafebabe))
	s.setupL1ViewContract(expectedValue)

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Send an actual L2 tx calling the L1STATICCALL precompile.
	calldata := buildL1STATICCALLCalldata(l1TestContractAddr, l1Head.Number, nil)
	_, err = testutils.SendDynamicFeeTx(
		s.RPCClient.L2,
		s.TestAddrPrivKey,
		&l1STATICCALLPrecompileAddr,
		common.Big0,
		calldata,
	)
	s.Nil(err)

	// Propose and sync.
	s.ProposeValidBlock(s.p)
	s.Nil(backoff.Retry(func() error {
		return s.d.ChainSyncer().EventSyncer().ProcessL1Blocks(context.Background())
	}, backoff.NewExponentialBackOff()))
	s.Nil(s.RPCClient.WaitTillL2ExecutionEngineSynced(context.Background()))

	// Find the L1STATICCALL tx in the latest block.
	l2Head, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.Greater(l2Head.Number.Uint64(), uint64(1))

	txCount, err := s.RPCClient.L2.TransactionCount(context.Background(), l2Head.Hash())
	s.Nil(err)
	s.GreaterOrEqual(txCount, uint(2), "block should have anchor tx + L1STATICCALL tx")

	var userTx *types.Transaction
	for idx := uint(0); idx < txCount; idx++ {
		tx, err := s.RPCClient.L2.TransactionInBlock(context.Background(), l2Head.Hash(), idx)
		s.Nil(err)
		if tx.To() != nil && *tx.To() == l1STATICCALLPrecompileAddr {
			userTx = tx
			break
		}
	}
	s.NotNil(userTx, "L1STATICCALL tx not found in block")

	receipt, err := s.RPCClient.L2.TransactionReceipt(context.Background(), userTx.Hash())
	s.Nil(err)
	s.Equal(types.ReceiptStatusSuccessful, receipt.Status)
	s.T().Logf("L1STATICCALL tx gas used: %d", receipt.GasUsed)
}

// TestL1STATICCALLInvalidBlockRejected verifies that requesting a non-existent (future)
// L1 block via the precompile fails.
func (s *DriverTestSuite) TestL1STATICCALLInvalidBlockRejected() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Request a block far in the future.
	futureBlock := new(big.Int).Add(l1Head.Number, big.NewInt(1000))
	calldata := buildL1STATICCALLCalldata(l1TestContractAddr, futureBlock, nil)

	_, err = s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata,
		Gas:  200_000,
	}, nil)
	s.NotNil(err, "Expected error when querying future L1 block")
}

// TestL1STATICCALLRevertPropagation verifies that an L1 call to a non-existent contract
// (no code) returns empty bytes through the precompile.
func (s *DriverTestSuite) TestL1STATICCALLRevertPropagation() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	// Call a non-existent contract (no code) — eth_call returns empty, which
	// the precompile should handle. The precompile returns whatever eth_call
	// returns, so an empty-code address returns empty bytes (not a revert).
	emptyAddr := common.HexToAddress("0x5555555555555555555555555555555555555555")

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// Call with a function selector that doesn't exist on the empty address.
	fnSelector := common.FromHex("0xdeadbeef")
	calldata := buildL1STATICCALLCalldata(emptyAddr, l1Head.Number, fnSelector)

	result, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata,
		Gas:  200_000,
	}, nil)
	// Calling a non-existent contract via eth_call returns empty data (not an error).
	// The precompile should return this empty result successfully.
	s.Nil(err)
	s.Empty(result, "Call to empty address should return empty bytes")
}

// TestL1STATICCALLShortInput verifies that the precompile rejects input shorter than
// the minimum 52 bytes (20 address + 32 block number).
func (s *DriverTestSuite) TestL1STATICCALLShortInput() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// Only 20 bytes (just the address, missing block number).
	shortCalldata := common.FromHex("0x4444444444444444444444444444444444444444")

	_, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: shortCalldata,
		Gas:  200_000,
	}, nil)
	s.NotNil(err, "Expected error for input shorter than 52 bytes")
	s.T().Logf("Short input error (expected): %v", err)
}

// TestL1STATICCALLWithCalldataPassthrough verifies that the precompile correctly passes
// variable-length calldata to the L1 contract. Uses a dynamic-slot contract that reads
// a storage slot key from calldata[0:32] and returns the value.
func (s *DriverTestSuite) TestL1STATICCALLWithCalldataPassthrough() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	// Deploy the dynamic-slot view contract on L1.
	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setCode",
		l1DynamicContractAddr, common.Bytes2Hex(dynamicSlotViewContractCode),
	))

	// Set two different storage slots with distinct values.
	slot1 := common.BigToHash(big.NewInt(1))
	value1 := common.BigToHash(big.NewInt(0xAAAA))
	slot2 := common.BigToHash(big.NewInt(2))
	value2 := common.BigToHash(big.NewInt(0xBBBB))

	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setStorageAt",
		l1DynamicContractAddr, slot1, value1,
	))
	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setStorageAt",
		l1DynamicContractAddr, slot2, value2,
	))
	s.L1Mine()

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// Call with calldata = slot1 key → should return value1.
	calldata1 := buildL1STATICCALLCalldata(l1DynamicContractAddr, l1Head.Number, slot1.Bytes())
	result1, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata1,
		Gas:  200_000,
	}, nil)
	s.Nil(err)
	s.Equal(value1, common.BytesToHash(result1),
		"Calldata with slot 1 should return value1")

	// Call with calldata = slot2 key → should return value2.
	calldata2 := buildL1STATICCALLCalldata(l1DynamicContractAddr, l1Head.Number, slot2.Bytes())
	result2, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata2,
		Gas:  200_000,
	}, nil)
	s.Nil(err)
	s.Equal(value2, common.BytesToHash(result2),
		"Calldata with slot 2 should return value2")

	s.T().Logf("Calldata passthrough: slot1→0x%x, slot2→0x%x", result1, result2)
}

// l1ExpensiveContractAddr is the address for the expensive (multi-SLOAD) contract.
var l1ExpensiveContractAddr = common.HexToAddress("0x7777777777777777777777777777777777777777")

// expensiveViewContractCode does 10 cold SLOADs (slots 0-9), returns slot 9 value.
// Each cold SLOAD costs 2100 gas on L1 → total ~21,000+ gas.
// Assembly: (PUSH1 n SLOAD POP) × 9, then PUSH1 9 SLOAD PUSH0 MSTORE PUSH1 0x20 PUSH0 RETURN
var expensiveViewContractCode = common.FromHex(
	"0x6000545060015450600254506003545060045450600554506006545060075450600854506009545F5260205FF3",
)

// TestL1STATICCALLGasIncludesL1Cost verifies that a transaction calling L1STATICCALL
// is charged more than just the static overhead, proving L1 consumed gas is included.
func (s *DriverTestSuite) TestL1STATICCALLGasIncludesL1Cost() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	expectedValue := common.BigToHash(big.NewInt(0xfeedface))
	s.setupL1ViewContract(expectedValue)

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	calldata := buildL1STATICCALLCalldata(l1TestContractAddr, l1Head.Number, nil)
	_, err = testutils.SendDynamicFeeTx(
		s.RPCClient.L2,
		s.TestAddrPrivKey,
		&l1STATICCALLPrecompileAddr,
		common.Big0,
		calldata,
	)
	s.Nil(err)

	s.ProposeValidBlock(s.p)
	s.Nil(backoff.Retry(func() error {
		return s.d.ChainSyncer().EventSyncer().ProcessL1Blocks(context.Background())
	}, backoff.NewExponentialBackOff()))
	s.Nil(s.RPCClient.WaitTillL2ExecutionEngineSynced(context.Background()))

	// Use TransactionInBlock rather than tx.Hash() because the proposer re-encodes
	// txs via blob tx lists, which can change the hash.
	l2Head, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	txCount, err := s.RPCClient.L2.TransactionCount(context.Background(), l2Head.Hash())
	s.Nil(err)
	s.GreaterOrEqual(txCount, uint(2), "block should have anchor tx + L1STATICCALL tx")

	var userTx *types.Transaction
	for idx := uint(0); idx < txCount; idx++ {
		t, err := s.RPCClient.L2.TransactionInBlock(context.Background(), l2Head.Hash(), idx)
		s.Nil(err)
		if t.To() != nil && *t.To() == l1STATICCALLPrecompileAddr {
			userTx = t
			break
		}
	}
	s.NotNil(userTx, "L1STATICCALL tx not found in block")

	receipt, err := s.RPCClient.L2.TransactionReceipt(context.Background(), userTx.Hash())
	s.Nil(err)
	s.Equal(types.ReceiptStatusSuccessful, receipt.Status)

	// Static overhead alone is ~33,000 (intrinsic 21000 + base 2000 + per-call 10000).
	// With dynamic L1 gas charging, the minimal view contract's SLOAD (~2100 gas on L1)
	// should push gasUsed above this baseline.
	s.T().Logf("L1STATICCALL gasUsed=%d (should include L1 consumed gas)", receipt.GasUsed)
	s.Greater(receipt.GasUsed, uint64(33_000),
		"gasUsed should exceed static overhead, proving L1 gas is charged")
}

// TestL1STATICCALLExpensiveContractHigherGas verifies that calling a more expensive L1 contract
// (10 SLOADs) results in higher gas consumption than a cheap one (1 SLOAD).
func (s *DriverTestSuite) TestL1STATICCALLExpensiveContractHigherGas() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	// Deploy cheap contract (1 SLOAD).
	storageVal := common.BigToHash(big.NewInt(42))
	s.setupL1ViewContract(storageVal)

	// Deploy expensive contract (10 SLOADs).
	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setCode",
		l1ExpensiveContractAddr, common.Bytes2Hex(expensiveViewContractCode),
	))
	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setStorageAt",
		l1ExpensiveContractAddr, common.BigToHash(big.NewInt(9)), storageVal,
	))
	s.L1Mine()

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// eth_call with cheap contract.
	cheapCalldata := buildL1STATICCALLCalldata(l1TestContractAddr, l1Head.Number, nil)
	cheapGas, err := s.RPCClient.L2.EstimateGas(context.Background(), ethereum.CallMsg{
		From: common.Address{},
		To:   &l1STATICCALLPrecompileAddr,
		Data: cheapCalldata,
		Gas:  500_000,
	})
	s.Nil(err)

	// eth_call with expensive contract.
	expensiveCalldata := buildL1STATICCALLCalldata(l1ExpensiveContractAddr, l1Head.Number, nil)
	expensiveGas, err := s.RPCClient.L2.EstimateGas(context.Background(), ethereum.CallMsg{
		From: common.Address{},
		To:   &l1STATICCALLPrecompileAddr,
		Data: expensiveCalldata,
		Gas:  500_000,
	})
	s.Nil(err)

	s.T().Logf("Gas estimates: cheap=%d, expensive=%d", cheapGas, expensiveGas)
	s.Greater(expensiveGas, cheapGas,
		"Expensive L1 contract (10 SLOADs) should cost more gas than cheap one (1 SLOAD)")
}

// TestL1STATICCALLLowGasLimitFails verifies that calling an expensive L1 contract with
// insufficient gas fails — the remaining L2 gas bounds the L1 call's gas limit.
func (s *DriverTestSuite) TestL1STATICCALLLowGasLimitFails() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	// Deploy expensive contract (10 SLOADs, ~21,000 L1 gas).
	s.Nil(s.RPCClient.L1.CallContext(
		context.Background(), nil, "anvil_setCode",
		l1ExpensiveContractAddr, common.Bytes2Hex(expensiveViewContractCode),
	))
	s.L1Mine()

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	calldata := buildL1STATICCALLCalldata(l1ExpensiveContractAddr, l1Head.Number, nil)

	// With only 13,000 gas total: after base (2000) + static overhead (10000),
	// only ~1,000 gas remains for the L1 call. The 10-SLOAD contract needs ~21,000.
	// The L1 call should OOG, causing precompile failure.
	_, err = s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata,
		Gas:  13_000,
	}, nil)
	s.NotNil(err, "Expensive L1 call with insufficient gas should fail")
	s.T().Logf("Low gas error (expected): %v", err)

	// Same call with plenty of gas should succeed.
	result, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata,
		Gas:  500_000,
	}, nil)
	s.Nil(err, "Same L1 call with sufficient gas should succeed")
	s.NotEmpty(result)
	s.T().Logf("High gas success: returned %d bytes", len(result))
}

// TestL1STATICCALLZeroBlockNumber verifies behavior when block number 0 is requested.
func (s *DriverTestSuite) TestL1STATICCALLZeroBlockNumber() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL only supported on NMC")
	}

	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// Block 0 (genesis) — should work since L1 genesis exists and has valid state.
	calldata := buildL1STATICCALLCalldata(
		l1TestContractAddr,
		big.NewInt(0),
		nil,
	)

	// This may succeed (returns empty/zero since no contract at genesis) or fail
	// depending on whether block 0 is within the lookback window of l1Origin.
	// Either way, it should not panic.
	_, err := s.RPCClient.L2.CallContract(context.Background(), ethereum.CallMsg{
		To:   &l1STATICCALLPrecompileAddr,
		Data: calldata,
		Gas:  200_000,
	}, nil)
	s.T().Logf("Block 0 result: err=%v", err)
}
