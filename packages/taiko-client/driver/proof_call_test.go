package driver

import (
	"context"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/internal/testutils"
)

// executionWitness mirrors the NMC Witness type (camelCase JSON, hex-encoded byte arrays).
// It is returned nested under the "witness" field of a proof_call result.
type executionWitness struct {
	Codes   []string `json:"codes"`
	State   []string `json:"state"`
	Keys    []string `json:"keys"`
	Headers []string `json:"headers"`
}

// proofCallResult is the proof_call response envelope (CallResultWithProof on the NMC side): the EVM
// return data, an optional in-VM error, and the execution witness captured during the call. An in-VM
// failure (revert, out-of-gas, ...) does NOT fail the RPC — it is reported in Error while the witness
// is still captured, so a verifier can re-prove the failure.
type proofCallResult struct {
	Result  string           `json:"result"`
	Error   *proofCallError  `json:"error,omitempty"`
	Witness executionWitness `json:"witness"`
}

// proofCallError is the in-payload error object proof_call returns on an in-VM failure.
type proofCallError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// callRequest is the JSON-RPC call object (a partial transaction) passed to proof_call.
type callRequest struct {
	To   string `json:"to"`
	Data string `json:"data,omitempty"`
	Gas  string `json:"gas,omitempty"`
}

// TestProofCallBasic calls proof_call on L2 NMC with a simple call and verifies the
// witness structure is returned.
func (s *DriverTestSuite) TestProofCallBasic() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("proof_call only available on NMC")
	}

	// Advance L2 past genesis so we have a valid block to query.
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	l2Head, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.Greater(l2Head.Number.Uint64(), uint64(0))

	// Call proof_call with a simple read of the Taiko anchor contract.
	// Any valid address works — the witness captures all state accessed during the call.
	taikoAnchor := common.HexToAddress("0x1670010000000000000000000000000000010001")
	blockHex := "0x" + l2Head.Number.Text(16)

	var result proofCallResult
	err = s.RPCClient.L2.CallContext(
		context.Background(),
		&result,
		"proof_call",
		callRequest{
			To:  taikoAnchor.Hex(),
			Gas: "0x100000",
		},
		blockHex,
	)
	s.Nil(err, "proof_call should succeed")

	// The witness should have state nodes — the EVM reads the account's nonce/balance/codeHash
	// and any storage accessed during the call.
	witness := result.Witness
	s.NotEmpty(witness.State, "witness should contain state trie nodes")
	s.T().Logf("Witness: %d state nodes, %d codes, %d keys, %d headers",
		len(witness.State), len(witness.Codes), len(witness.Keys), len(witness.Headers))
}

// TestProofCallWithData calls proof_call with actual calldata (a function selector)
// and verifies the witness captures code.
func (s *DriverTestSuite) TestProofCallWithData() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("proof_call only available on NMC")
	}

	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	l2Head, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Call with the paused() selector (0x5c975abb) on the Taiko anchor contract.
	// Even if it reverts, proof_call still succeeds at the RPC layer and captures the witness.
	taikoAnchor := common.HexToAddress("0x1670010000000000000000000000000000010001")
	blockHex := "0x" + l2Head.Number.Text(16)

	var result proofCallResult
	err = s.RPCClient.L2.CallContext(
		context.Background(),
		&result,
		"proof_call",
		callRequest{
			To:   taikoAnchor.Hex(),
			Data: "0x5c975abb",
			Gas:  "0x100000",
		},
		blockHex,
	)
	s.Nil(err, "proof_call should succeed even if the call reverts internally")
	witness := result.Witness
	s.NotEmpty(witness.State, "witness should contain state trie nodes")

	// When calling a contract with code, the witness should capture the bytecode.
	s.NotEmpty(witness.Codes, "witness should contain contract bytecodes")
	s.T().Logf("Witness with calldata: %d state nodes, %d codes, %d keys, %d headers",
		len(witness.State), len(witness.Codes), len(witness.Keys), len(witness.Headers))
}

// TestProofCallGenesisBlock verifies that requesting a witness for the genesis
// block returns an error (genesis has no parent header to walk back from).
func (s *DriverTestSuite) TestProofCallGenesisBlock() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("proof_call only available on NMC")
	}

	taikoAnchor := common.HexToAddress("0x1670010000000000000000000000000000010001")

	var result proofCallResult
	err := s.RPCClient.L2.CallContext(
		context.Background(),
		&result,
		"proof_call",
		callRequest{To: taikoAnchor.Hex(), Gas: "0x100000"},
		"0x0", // genesis block
	)
	s.NotNil(err, "Should fail for genesis block")
	s.T().Logf("Genesis block error (expected): %v", err)
}

// TestProofCallNonExistentBlock verifies that requesting a witness for a block
// that doesn't exist returns an error.
func (s *DriverTestSuite) TestProofCallNonExistentBlock() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("proof_call only available on NMC")
	}

	taikoAnchor := common.HexToAddress("0x1670010000000000000000000000000000010001")

	var result proofCallResult
	err := s.RPCClient.L2.CallContext(
		context.Background(),
		&result,
		"proof_call",
		callRequest{To: taikoAnchor.Hex(), Gas: "0x100000"},
		"0xFFFFFF", // block 16777215, far in the future
	)
	s.NotNil(err, "Should fail for non-existent block")
	s.T().Logf("Non-existent block error (expected): %v", err)
}

// TestProofCallLatestBlock verifies that the "latest" block parameter works.
func (s *DriverTestSuite) TestProofCallLatestBlock() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("proof_call only available on NMC")
	}

	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	taikoAnchor := common.HexToAddress("0x1670010000000000000000000000000000010001")

	var result proofCallResult
	err := s.RPCClient.L2.CallContext(
		context.Background(),
		&result,
		"proof_call",
		callRequest{To: taikoAnchor.Hex(), Gas: "0x100000"},
		"latest",
	)
	s.Nil(err, "proof_call should work with 'latest'")
	s.NotEmpty(result.Witness.State, "witness should contain state trie nodes")
}

// TestProofCallToEOA verifies that calling an EOA (no code) still returns a valid
// witness — just with account state, no codes.
func (s *DriverTestSuite) TestProofCallToEOA() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("proof_call only available on NMC")
	}

	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	l2Head, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Call to a random EOA — no code, should still produce a witness with account state.
	randomEOA := common.HexToAddress("0x9999999999999999999999999999999999999999")
	blockHex := "0x" + l2Head.Number.Text(16)

	var result proofCallResult
	err = s.RPCClient.L2.CallContext(
		context.Background(),
		&result,
		"proof_call",
		callRequest{To: randomEOA.Hex(), Gas: "0x100000"},
		blockHex,
	)
	s.Nil(err, "proof_call to EOA should succeed")
	// Witness should have state (sender + receiver account nodes) but likely no codes.
	witness := result.Witness
	s.NotEmpty(witness.State, "witness should contain account state nodes even for EOA")
	s.T().Logf("EOA witness: %d state nodes, %d codes", len(witness.State), len(witness.Codes))
}

// TestL1STATICCALLEndToEndWithWitness is the full pipeline test:
// 1. Deploy L1 view contract
// 2. Send L2 tx calling L1STATICCALL precompile
// 3. Propose block, sync, verify receipt
// 4. Generate execution witness (via proof_call) for a call on that L2 block
// 5. Verify the witness captures state
func (s *DriverTestSuite) TestL1STATICCALLEndToEndWithWitness() {
	if os.Getenv("L2_NODE") != testutils.L2NodeNMC {
		s.T().Skip("L1STATICCALL + proof_call only supported on NMC")
	}

	expectedValue := common.BigToHash(big.NewInt(0x42424242))

	// Step 1: Deploy minimal view contract on L1.
	s.setupL1ViewContract(expectedValue)

	l1Head, err := s.RPCClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Step 2: Advance L2 chain so l1Origin covers our L1 block.
	s.ProposeAndInsertValidBlock(s.p, s.d.ChainSyncer().EventSyncer())

	// Step 3: Send L2 tx calling L1STATICCALL.
	calldata := buildL1STATICCALLCalldata(l1TestContractAddr, l1Head.Number, nil)
	_, err = testutils.SendDynamicFeeTx(
		s.RPCClient.L2,
		s.TestAddrPrivKey,
		&l1STATICCALLPrecompileAddr,
		common.Big0,
		calldata,
	)
	s.Nil(err)

	// Step 4: Propose and sync the block containing the L1STATICCALL tx.
	s.ProposeValidBlock(s.p)
	s.Nil(s.d.ChainSyncer().EventSyncer().ProcessL1Blocks(context.Background()))
	s.Nil(s.RPCClient.WaitTillL2ExecutionEngineSynced(context.Background()))

	// Step 5: Verify the tx receipt.
	l2Head, err := s.RPCClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.Greater(l2Head.Number.Uint64(), uint64(1))

	txCount, err := s.RPCClient.L2.TransactionCount(context.Background(), l2Head.Hash())
	s.Nil(err)
	s.GreaterOrEqual(txCount, uint(2))

	// Step 6: Generate execution witness for a call on the block that processed the L1STATICCALL tx.
	// This simulates what raiko would do to generate proving data for the block.
	taikoAnchor := common.HexToAddress("0x1670010000000000000000000000000000010001")
	blockHex := "0x" + l2Head.Number.Text(16)

	var result proofCallResult
	err = s.RPCClient.L2.CallContext(
		context.Background(),
		&result,
		"proof_call",
		callRequest{To: taikoAnchor.Hex(), Gas: "0x100000"},
		blockHex,
	)
	s.Nil(err, "proof_call should succeed on block with L1STATICCALL tx")
	witness := result.Witness
	s.NotEmpty(witness.State, "witness should contain state trie nodes")
	s.T().Logf("E2E witness: %d state nodes, %d codes, %d keys, %d headers",
		len(witness.State), len(witness.Codes), len(witness.Keys), len(witness.Headers))
}
