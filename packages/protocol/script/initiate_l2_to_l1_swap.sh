#!/usr/bin/env bash

# =============================================================================
# Initiate L2 → L1 → L2 Swap (ETH → bUSDC) — devnet smoke test
# =============================================================================
# Drives one end-to-end synchronous L2→L1→L2 swap on the bidirectional cross-chain
# DEX. Mirrors what the UI will do for the "Via L1 DEX" swap venue:
#
#   1. Read deployment addresses (cross-chain-dex-l{1,2}.json).
#   2. Build the L2 tx calldata with a PLACEHOLDER return message.
#   3. Call Catalyst's `surge_simulateReturnMessage` with that calldata to get
#      the REAL return message the bridge will emit after the L1 swap.
#   4. Re-encode the calldata with the real message embedded.
#   5. Broadcast the L2 tx with an explicit --gas-limit (estimation would fail,
#      the fast signal isn't present until the block is built).
#
# Required env:
#   PRIVATE_KEY     - L2 EOA that initiates the swap
#   L1_RPC          - L1 RPC URL
#   L2_RPC          - L2 RPC URL (the one Catalyst drives)
#   CATALYST_RPC    - Catalyst JSON-RPC endpoint serving `surge_simulateReturnMessage`
#
# Optional env:
#   AMOUNT_ETH      - wei to swap from L2 ETH (default: 0.1 ETH)
#   MIN_TOKEN_OUT   - minimum bUSDC accepted (default: 0 — no slippage protection)
#   RECIPIENT       - L2 bUSDC recipient (default: sender)
# =============================================================================

set -euo pipefail

: "${PRIVATE_KEY:?PRIVATE_KEY is required}"
: "${L1_RPC:?L1_RPC is required}"
: "${L2_RPC:?L2_RPC is required}"
: "${CATALYST_RPC:?CATALYST_RPC is required}"

AMOUNT_ETH=${AMOUNT_ETH:-"100000000000000000"}  # 0.1 ETH
MIN_TOKEN_OUT=${MIN_TOKEN_OUT:-"0"}

SENDER=$(cast wallet address --private-key "$PRIVATE_KEY")
RECIPIENT=${RECIPIENT:-$SENDER}

L1_DEPLOY_JSON="deployments/cross-chain-dex-l1.json"
L2_DEPLOY_JSON="deployments/cross-chain-dex-l2.json"

if [ ! -f "$L1_DEPLOY_JSON" ] || [ ! -f "$L2_DEPLOY_JSON" ]; then
    echo "ERROR: deployment files missing. Run ./script/deploy_cross_chain_dex.sh first."
    exit 1
fi

L1_VAULT=$(python3 -c "import json; print(json.load(open('$L1_DEPLOY_JSON'))['CrossChainSwapVaultL1'])")
L2_VAULT=$(python3 -c "import json; print(json.load(open('$L2_DEPLOY_JSON'))['CrossChainSwapVaultL2'])")
L1_CHAIN_ID=$(cast chain-id --rpc-url "$L1_RPC")
L2_CHAIN_ID=$(cast chain-id --rpc-url "$L2_RPC")

echo "============================================="
echo " L2→L1→L2 swap (ETH → bUSDC) smoke test"
echo "============================================="
echo "  Sender:       $SENDER"
echo "  Recipient:    $RECIPIENT"
echo "  L1 Vault:     $L1_VAULT"
echo "  L2 Vault:     $L2_VAULT"
echo "  L1 Chain ID:  $L1_CHAIN_ID"
echo "  L2 Chain ID:  $L2_CHAIN_ID"
echo "  Amount:       $AMOUNT_ETH wei"
echo "  Min out:      $MIN_TOKEN_OUT"
echo ""

# ---------------------------------------------------------------
# Action enum values (must match CrossChainSwapVaultL2.Action)
#   0 BRIDGE
#   1 SWAP_ETH_TO_TOKEN
#   2 SWAP_TOKEN_TO_ETH
#   3 ADD_LIQUIDITY
#   4 REMOVE_LIQUIDITY
#   5 SWAP_ETH_TO_TOKEN_VIA_L1   <-- this flow
#   6 SWAP_TOKEN_TO_ETH_VIA_L1
# ---------------------------------------------------------------
ACTION=5

# ---------------------------------------------------------------
# 1. Build placeholder return message
#    The real L1Vault return payload is abi.encode(uint8 action, address recipient, uint256 tokenOut)
#    wrapped in onMessageInvocation(bytes). For the placeholder we plug in tokenOut=0 and let
#    Catalyst's simulation overwrite the whole message with the real values.
# ---------------------------------------------------------------

PLACEHOLDER_INNER=$(cast abi-encode "f(uint8,address,uint256)" "$ACTION" "$RECIPIENT" "0")
PLACEHOLDER_ONMSG=$(cast calldata "onMessageInvocation(bytes)" "$PLACEHOLDER_INNER")

MSG_TUPLE_TYPE="(uint64,uint64,uint32,address,uint64,address,uint64,address,address,uint256,bytes)"
PLACEHOLDER_MSG="(0,0,1000000,0x0000000000000000000000000000000000000000,0,$L1_VAULT,$L2_CHAIN_ID,$L2_VAULT,$L2_VAULT,0,$PLACEHOLDER_ONMSG)"

# ---------------------------------------------------------------
# 2. Build L2 tx calldata with placeholder
# ---------------------------------------------------------------

SIM_CALLDATA=$(cast calldata \
    "swapETHForTokenViaL1(uint256,address,$MSG_TUPLE_TYPE)" \
    "$MIN_TOKEN_OUT" "$RECIPIENT" "$PLACEHOLDER_MSG")

# ---------------------------------------------------------------
# 3. Ask Catalyst to simulate the return message
# ---------------------------------------------------------------

echo "Calling surge_simulateReturnMessage..."
SIM_RESPONSE=$(curl -s -X POST "$CATALYST_RPC" \
    -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"method\":\"surge_simulateReturnMessage\",\"params\":[{\"from\":\"$SENDER\",\"to\":\"$L2_VAULT\",\"data\":\"$SIM_CALLDATA\",\"value\":\"$AMOUNT_ETH\"}],\"id\":1}" \
    | python3 -c "import json,sys; r=json.load(sys.stdin); print(json.dumps(r.get('result')) if 'result' in r else f'ERROR:{r.get(\"error\",{}).get(\"message\",\"unknown\")}')")

if [ -z "$SIM_RESPONSE" ] || echo "$SIM_RESPONSE" | grep -q "^ERROR:"; then
    echo "ERROR: surge_simulateReturnMessage failed: $SIM_RESPONSE"
    exit 1
fi

echo "  ✓ Catalyst returned a simulated return message"
echo ""

# ---------------------------------------------------------------
# 4. Extract real message fields and rebuild calldata
# ---------------------------------------------------------------

REAL_MSG_TUPLE=$(python3 - <<PY
import json, sys
resp = json.loads('''$SIM_RESPONSE''')
m = resp["message"]
# All numerics come back as decimal strings from the RPC handler.
print("({},{},{},{},{},{},{},{},{},{},{})".format(
    m["id"], m["fee"], m["gasLimit"], m["from"], m["srcChainId"],
    m["srcOwner"], m["destChainId"], m["destOwner"], m["to"],
    m["value"], m["data"],
))
PY
)

REAL_CALLDATA=$(cast calldata \
    "swapETHForTokenViaL1(uint256,address,$MSG_TUPLE_TYPE)" \
    "$MIN_TOKEN_OUT" "$RECIPIENT" "$REAL_MSG_TUPLE")

# ---------------------------------------------------------------
# 5. Broadcast on L2
# ---------------------------------------------------------------

echo "Broadcasting L2 tx with --gas-limit 3000000..."
TX_HASH=$(cast send "$L2_VAULT" \
    --value "$AMOUNT_ETH" \
    --gas-limit 3000000 \
    --rpc-url "$L2_RPC" \
    --private-key "$PRIVATE_KEY" \
    "$REAL_CALLDATA" \
    --json | python3 -c "import json,sys; print(json.load(sys.stdin)['transactionHash'])")

echo ""
echo "  L2 tx hash:  $TX_HASH"
echo ""

# ---------------------------------------------------------------
# 6. Wait for receipt
# ---------------------------------------------------------------

RECEIPT=$(cast receipt "$TX_HASH" --rpc-url "$L2_RPC" --json)
STATUS=$(echo "$RECEIPT" | python3 -c "import json,sys; print(json.load(sys.stdin)['status'])")

if [ "$STATUS" = "0x1" ]; then
    echo "  ✓ L2 tx succeeded"
    echo ""
    echo "Next: observe the L2 receipt logs — expect L1DexSwapInitiatedETHForToken"
    echo "      and L1DexSwapCompletedETHForToken from CrossChainSwapVaultL2."
else
    echo "  ✗ L2 tx FAILED (status=$STATUS)"
    echo "  Full receipt:"
    echo "$RECEIPT" | python3 -m json.tool
    exit 1
fi
