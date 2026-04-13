#!/bin/sh

# =============================================================================
# Flash Loan Initiation Script (L2 side)
# =============================================================================
# Submits an L2 UserOp to Catalyst's surge_sendUserOp RPC that initiates a
# synchronous L2→L1→L2 flash loan on L2.
#
# The UserOp calldata targets FlashLoanExecutorL2.execute(amount, beneficiary,
# returnMessage). The `returnMessage` field is populated with a placeholder
# struct; Catalyst's builder is responsible for simulating the L1 callback,
# recomputing the actual return Message (with correct bridge id/srcChainId),
# and substituting it into the UserOp calldata when assembling the L2 block.
#
# Usage:
#   CATALYST_RPC=http://localhost:9545 \
#   PRIVATE_KEY=0x... \
#   AMOUNT=1000000000 \            # 1000 tokens (6 decimals)
#   BENEFICIARY=0x...  \           # receives the fee profit
#   ./script/initiate_flash_loan.sh
# =============================================================================

set -e

export CATALYST_RPC=${CATALYST_RPC:-"http://localhost:9545"}
export PRIVATE_KEY=${PRIVATE_KEY:-""}
export L2_RPC=${L2_RPC:-""}

if [ -z "$PRIVATE_KEY" ]; then
    echo "ERROR: PRIVATE_KEY is required"
    exit 1
fi

if [ -z "$L2_RPC" ]; then
    echo "ERROR: L2_RPC is required"
    exit 1
fi

# Default amount: 1000 tokens at 6 decimals
export AMOUNT=${AMOUNT:-"1000000000"}

# Read deployment JSON for executor address
DEPLOY_JSON="deployments/flash-loan-l2.json"
L1_DEPLOY_JSON="deployments/flash-loan-l1.json"

L2_EXECUTOR=${L2_EXECUTOR:-$(python3 -c "import json; print(json.load(open('$DEPLOY_JSON'))['FlashLoanExecutorL2'])")}
L1_CALLBACK=${L1_CALLBACK:-$(python3 -c "import json; print(json.load(open('$L1_DEPLOY_JSON'))['FlashLoanCallbackL1'])")}

SENDER=$(cast wallet address --private-key "$PRIVATE_KEY")
BENEFICIARY=${BENEFICIARY:-"$SENDER"}

L2_CHAIN_ID=$(cast chain-id --rpc-url "$L2_RPC")

echo "============================================="
echo " Initiate Flash Loan (L2 -> L1 -> L2)"
echo "============================================="
echo "  Catalyst RPC:  $CATALYST_RPC"
echo "  L2 Executor:   $L2_EXECUTOR"
echo "  L1 Callback:   $L1_CALLBACK"
echo "  Sender:        $SENDER"
echo "  Beneficiary:   $BENEFICIARY"
echo "  Amount:        $AMOUNT (raw units)"
echo "  L2 Chain ID:   $L2_CHAIN_ID"
echo ""

# ---------------------------------------------------------------
# Build execute() calldata
# ---------------------------------------------------------------
#
# Signature:
#   function execute(
#       uint256 amount,
#       address beneficiary,
#       IBridge.Message calldata returnMessage
#   ) external
#
# IBridge.Message struct:
#   (uint64 id, uint64 fee, uint32 gasLimit, address from, uint64 srcChainId,
#    address srcOwner, uint64 destChainId, address destOwner, address to,
#    uint256 value, bytes data)
#
# The returnMessage fields are placeholders. Catalyst's builder will:
#   1. Simulate the L1 callback to determine the actual Message (in particular
#      the auto-assigned `id` from L1 Bridge.sendMessage).
#   2. Rebuild the L2 block with the accurate Message substituted into the
#      UserOp calldata.
#   3. Include the slot derived from hashMessage(Message) in both
#      `requiredReturnSignals` of the ProposeInputV2 and the anchor's fast
#      signals.
# ---------------------------------------------------------------

L1_CHAIN_ID=${L1_CHAIN_ID:-$(cast chain-id --rpc-url "${L1_RPC:-http://localhost:8545}" 2>/dev/null || echo 1)}

# Placeholder Message — the builder will overwrite during simulation.
# `to` must be the L2 executor (where the bridge will deliver the return).
# `from` will be filled by L1 bridge to be the L1 callback.
# `data` is abi.encodeWithSignature("onMessageInvocation(bytes)", abi.encode(total, beneficiary))
# with total = amount + fee (fee = amount * 100 / 10000 = 1%).
FEE=$(python3 -c "print($AMOUNT * 100 // 10000)")
TOTAL=$(python3 -c "print($AMOUNT + $FEE)")

# abi.encode(uint256 total, address beneficiary)
INNER_PAYLOAD=$(cast abi-encode "f(uint256,address)" "$TOTAL" "$L2_EXECUTOR")
# abi.encodeWithSignature("onMessageInvocation(bytes)", INNER_PAYLOAD)
RETURN_MSG_DATA=$(cast calldata "onMessageInvocation(bytes)" "$INNER_PAYLOAD")

# Build the struct fields for IBridge.Message — placeholder id=0 is filled by builder
RETURN_MSG="(0,0,1000000,0x0000000000000000000000000000000000000000,$L1_CHAIN_ID,$L1_CALLBACK,$L2_CHAIN_ID,$L2_EXECUTOR,$L2_EXECUTOR,0,$RETURN_MSG_DATA)"

CALLDATA=$(cast calldata \
    "execute(uint256,address,(uint64,uint64,uint32,address,uint64,address,uint64,address,address,uint256,bytes))" \
    "$AMOUNT" \
    "$BENEFICIARY" \
    "$RETURN_MSG")

echo "execute() calldata (length: ${#CALLDATA}):"
echo "$CALLDATA" | head -c 200
echo "..."
echo ""

# ---------------------------------------------------------------
# Submit via Catalyst's surge_sendUserOp RPC
# ---------------------------------------------------------------

REQUEST=$(python3 -c "
import json
print(json.dumps({
    'jsonrpc': '2.0',
    'method': 'surge_sendUserOp',
    'params': [{
        'submitter': '$L2_EXECUTOR',
        'calldata': '$CALLDATA',
        'chainId': $L2_CHAIN_ID
    }],
    'id': 1
}))")

echo "Sending to $CATALYST_RPC ..."
RESPONSE=$(curl -s -X POST "$CATALYST_RPC" \
    -H "Content-Type: application/json" \
    --data "$REQUEST")

echo ""
echo "Response:"
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
echo ""

USER_OP_ID=$(echo "$RESPONSE" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('result',''))" 2>/dev/null || echo "")

if [ -n "$USER_OP_ID" ]; then
    echo "UserOp submitted. id = $USER_OP_ID"
    echo ""
    echo "Poll status with:"
    echo "  curl -s -X POST $CATALYST_RPC \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    --data '{\"jsonrpc\":\"2.0\",\"method\":\"surge_userOpStatus\",\"params\":[$USER_OP_ID],\"id\":1}'"
else
    echo "Submission failed. See response above."
    exit 1
fi
