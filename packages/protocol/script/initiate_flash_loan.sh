#!/bin/sh

# =============================================================================
# Flash Loan Initiation Script (L2 side)
# =============================================================================
# Initiates a synchronous L2→L1→L2 flash loan by:
#   1. Calling Catalyst's surge_simulateReturnMessage to get the exact
#      L1→L2 return Message the L1 callback will produce.
#   2. Building FlashLoanExecutorL2.execute(amount, beneficiary, returnMessage)
#      calldata with the simulated Message.
#   3. Submitting a normal L2 transaction to the L2 mempool via cast send.
#
# Catalyst picks up the tx from the mempool during block building, detects
# the outbound Bridge.sendMessage, and injects the return signal into the
# anchor's fast signals automatically.
#
# Usage:
#   CATALYST_RPC=http://localhost:4545 \
#   PRIVATE_KEY=0x... \
#   L2_RPC=http://localhost:8547 \
#   AMOUNT=1000000000 \            # 1000 tokens (6 decimals)
#   BENEFICIARY=0x...  \           # receives the fee profit
#   ./script/initiate_flash_loan.sh
# =============================================================================

set -e

export CATALYST_RPC=${CATALYST_RPC:-"http://localhost:4545"}
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
L1_CHAIN_ID=${L1_CHAIN_ID:-$(cast chain-id --rpc-url "${L1_RPC:-http://localhost:32003}" 2>/dev/null || echo 1)}

echo "============================================="
echo " Initiate Flash Loan (L2 → L1 → L2)"
echo "============================================="
echo "  Catalyst RPC:  $CATALYST_RPC"
echo "  L2 RPC:        $L2_RPC"
echo "  L2 Executor:   $L2_EXECUTOR"
echo "  L1 Callback:   $L1_CALLBACK"
echo "  Sender:        $SENDER"
echo "  Beneficiary:   $BENEFICIARY"
echo "  Amount:        $AMOUNT (raw units)"
echo "  L2 Chain ID:   $L2_CHAIN_ID"
echo "  L1 Chain ID:   $L1_CHAIN_ID"
echo ""

# ---------------------------------------------------------------
# Step 1: Build placeholder execute() calldata for simulation
# ---------------------------------------------------------------
# We need calldata to send to surge_simulateReturnMessage. The
# placeholder returnMessage is a zeroed-out Message struct — the
# simulation traces the tx to find the outbound Bridge.sendMessage
# (which doesn't depend on returnMessage), then simulates the L1
# callback to get the real return.

FEE=$(python3 -c "print($AMOUNT * 100 // 10000)")
TOTAL=$(python3 -c "print($AMOUNT + $FEE)")

INNER_PAYLOAD=$(cast abi-encode "f(uint256,address)" "$TOTAL" "$L2_EXECUTOR")
RETURN_MSG_DATA=$(cast calldata "onMessageInvocation(bytes)" "$INNER_PAYLOAD")

# Placeholder Message for simulation (zeroed fields — will be replaced)
PLACEHOLDER_MSG="(0,0,1000000,0x0000000000000000000000000000000000000000,$L1_CHAIN_ID,$L1_CALLBACK,$L2_CHAIN_ID,$L2_EXECUTOR,$L2_EXECUTOR,0,$RETURN_MSG_DATA)"

SIM_CALLDATA=$(cast calldata \
    "execute(uint256,address,(uint64,uint64,uint32,address,uint64,address,uint64,address,address,uint256,bytes))" \
    "$AMOUNT" \
    "$BENEFICIARY" \
    "$PLACEHOLDER_MSG")

echo "Step 1: Simulating return message via Catalyst..."

# ---------------------------------------------------------------
# Step 2: Call surge_simulateReturnMessage
# ---------------------------------------------------------------
SIM_REQUEST=$(python3 -c "
import json
print(json.dumps({
    'jsonrpc': '2.0',
    'method': 'surge_simulateReturnMessage',
    'params': [{
        'from': '$L2_EXECUTOR',
        'to': '$L2_EXECUTOR',
        'data': '$SIM_CALLDATA'
    }],
    'id': 1
}))")

SIM_RESPONSE=$(curl -s -X POST "$CATALYST_RPC" \
    -H "Content-Type: application/json" \
    --data "$SIM_REQUEST")

echo "Simulation response:"
echo "$SIM_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$SIM_RESPONSE"
echo ""

# Extract the simulated Message fields
MSG_JSON=$(echo "$SIM_RESPONSE" | python3 -c "
import sys, json
d = json.load(sys.stdin)
if 'error' in d:
    print('ERROR: ' + json.dumps(d['error']), file=sys.stderr)
    sys.exit(1)
msg = d['result']['message']
print(json.dumps(msg))
") || { echo "Simulation failed"; exit 1; }

# Parse fields
MSG_ID=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
MSG_FEE=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['fee'])")
MSG_GAS_LIMIT=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['gasLimit'])")
MSG_FROM=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['from'])")
MSG_SRC_CHAIN_ID=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['srcChainId'])")
MSG_SRC_OWNER=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['srcOwner'])")
MSG_DEST_CHAIN_ID=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['destChainId'])")
MSG_DEST_OWNER=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['destOwner'])")
MSG_TO=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['to'])")
MSG_VALUE=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['value'])")
MSG_DATA=$(echo "$MSG_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['data'])")

echo "Step 2: Got simulated return Message:"
echo "  id=$MSG_ID, gasLimit=$MSG_GAS_LIMIT, srcChainId=$MSG_SRC_CHAIN_ID"
echo "  from=$MSG_FROM, to=$MSG_TO"
echo ""

# ---------------------------------------------------------------
# Step 3: Build the real execute() calldata with simulated Message
# ---------------------------------------------------------------
RETURN_MSG="($MSG_ID,$MSG_FEE,$MSG_GAS_LIMIT,$MSG_FROM,$MSG_SRC_CHAIN_ID,$MSG_SRC_OWNER,$MSG_DEST_CHAIN_ID,$MSG_DEST_OWNER,$MSG_TO,$MSG_VALUE,$MSG_DATA)"

REAL_CALLDATA=$(cast calldata \
    "execute(uint256,address,(uint64,uint64,uint32,address,uint64,address,uint64,address,address,uint256,bytes))" \
    "$AMOUNT" \
    "$BENEFICIARY" \
    "$RETURN_MSG")

echo "Step 3: Submitting L2 transaction to mempool..."

# ---------------------------------------------------------------
# Step 4: Submit as a normal L2 transaction
# ---------------------------------------------------------------
# cast send with raw calldata: use the sig+args form with empty sig
TX_HASH=$(cast send "$L2_EXECUTOR" \
    --private-key "$PRIVATE_KEY" \
    --rpc-url "$L2_RPC" \
    -- "$REAL_CALLDATA" \
    2>&1 | grep "transactionHash" | awk '{print $2}') || true

# Fallback: use cast publish with a signed raw tx
if [ -z "$TX_HASH" ]; then
    TX_RESULT=$(cast send "$L2_EXECUTOR" \
        --private-key "$PRIVATE_KEY" \
        --rpc-url "$L2_RPC" \
        "execute(uint256,address,(uint64,uint64,uint32,address,uint64,address,uint64,address,address,uint256,bytes))" \
        "$AMOUNT" \
        "$BENEFICIARY" \
        "$RETURN_MSG" \
        2>&1)
    TX_HASH=$(echo "$TX_RESULT" | grep "transactionHash" | awk '{print $2}')
fi

echo ""
echo "============================================="
echo " Flash Loan Submitted"
echo "============================================="
echo "  L2 TX Hash: $TX_HASH"
echo ""
echo "  Catalyst will pick this up from the mempool,"
echo "  detect the outbound signal, inject the return"
echo "  signal into the anchor, and include the tx."
echo ""
echo "  Check status:"
echo "    cast receipt $TX_HASH --rpc-url $L2_RPC"
echo ""
echo "  Or poll Catalyst's surge_txStatus:"
echo "    curl -s -X POST $CATALYST_RPC \\"
echo "      -H 'Content-Type: application/json' \\"
echo "      --data '{\"jsonrpc\":\"2.0\",\"method\":\"surge_txStatus\",\"params\":[{\"txHash\":\"$TX_HASH\"}],\"id\":1}'"
