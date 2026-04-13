#!/bin/sh

# =============================================================================
# L2 -> L1 -> L2 Flash Loan Full Deployment Script
# =============================================================================
# Deploys the complete flash loan system on L1 and L2, links the two sides,
# and outputs a summary of deployed addresses.
#
# Usage:
#   PRIVATE_KEY=0x... ./script/deploy_flash_loan.sh
# =============================================================================

set -e

export PRIVATE_KEY=${PRIVATE_KEY:-""}

if [ -z "$PRIVATE_KEY" ]; then
    echo "ERROR: PRIVATE_KEY is required"
    echo "Usage: PRIVATE_KEY=0x... ./script/deploy_flash_loan.sh"
    exit 1
fi

export L1_RPC=${L1_RPC:-""}
export L2_RPC=${L2_RPC:-""}
export L1_BRIDGE=${L1_BRIDGE:-"0xc1e59A201cE4CD58590FC3Ab45081921cF186550"}
export L2_BRIDGE=${L2_BRIDGE:-"0x7633740000000000000000000000000000000001"}
export TOKEN_DECIMALS=${TOKEN_DECIMALS:-"6"}
export INITIAL_POOL_LIQUIDITY=${INITIAL_POOL_LIQUIDITY:-"1000000000000"} # 1M tokens with 6 decimals
export LOG_LEVEL=${LOG_LEVEL:-"-vvvv"}
export DRY_RUN=${DRY_RUN:-"false"}

DEPLOY_DIR="deployments"
L1_DEPLOY_JSON="$DEPLOY_DIR/flash-loan-l1.json"
L2_DEPLOY_JSON="$DEPLOY_DIR/flash-loan-l2.json"

echo "============================================="
echo " Flash Loan L2->L1->L2 Full Deployment"
echo "============================================="
echo ""

L1_CHAIN_ID=$(cast chain-id --rpc-url "$L1_RPC")
L2_CHAIN_ID=$(cast chain-id --rpc-url "$L2_RPC")
SENDER=$(cast wallet address --private-key "$PRIVATE_KEY")

echo "  Deployer:      $SENDER"
echo "  L1 RPC:        $L1_RPC"
echo "  L2 RPC:        $L2_RPC"
echo "  L1 Chain ID:   $L1_CHAIN_ID"
echo "  L2 Chain ID:   $L2_CHAIN_ID"
echo "  L1 Bridge:     $L1_BRIDGE"
echo "  L2 Bridge:     $L2_BRIDGE"
echo "  Token Decimals:$TOKEN_DECIMALS"
echo "  Pool Liquidity:$INITIAL_POOL_LIQUIDITY"
echo ""

if [ "$DRY_RUN" = "true" ]; then
    echo "  *** DRY RUN — no transactions will be broadcast ***"
    BROADCAST_ARG=""
else
    BROADCAST_ARG="--broadcast"
fi
echo ""

# Step 1: L1
echo "============================================="
echo " Step 1/3: Deploying L1 contracts"
echo "============================================="
export L2_CHAIN_ID
export FOUNDRY_PROFILE="layer1"

forge script ./script/layer1/surge/flash-loan/DeployFlashLoanL1.s.sol:DeployFlashLoanL1 \
    --fork-url "$L1_RPC" \
    $BROADCAST_ARG \
    $LOG_LEVEL \
    --private-key "$PRIVATE_KEY"

if [ "$DRY_RUN" = "true" ]; then
    echo "DRY RUN: Skipping remaining steps"
    exit 0
fi

L1_TOKEN=$(python3 -c "import json; print(json.load(open('$L1_DEPLOY_JSON'))['FlashLoanToken'])")
L1_CALLBACK=$(python3 -c "import json; print(json.load(open('$L1_DEPLOY_JSON'))['FlashLoanCallbackL1'])")

echo ""
echo "  L1 Token:      $L1_TOKEN"
echo "  L1 Callback:   $L1_CALLBACK"
echo ""

# Step 2: L2
echo "============================================="
echo " Step 2/3: Deploying L2 contracts"
echo "============================================="
export L1_CHAIN_ID
export FOUNDRY_PROFILE="layer2"

forge script ./script/layer2/surge/flash-loan/DeployFlashLoanL2.s.sol:DeployFlashLoanL2 \
    --fork-url "$L2_RPC" \
    $BROADCAST_ARG \
    $LOG_LEVEL \
    --private-key "$PRIVATE_KEY"

L2_TOKEN=$(python3 -c "import json; print(json.load(open('$L2_DEPLOY_JSON'))['FlashLoanToken'])")
L2_POOL=$(python3 -c "import json; print(json.load(open('$L2_DEPLOY_JSON'))['FlashLoanPool'])")
L2_EXECUTOR=$(python3 -c "import json; print(json.load(open('$L2_DEPLOY_JSON'))['FlashLoanExecutorL2'])")

echo ""
echo "  L2 Token:      $L2_TOKEN"
echo "  L2 Pool:       $L2_POOL"
echo "  L2 Executor:   $L2_EXECUTOR"
echo ""

# Step 3: Link
echo "============================================="
echo " Step 3/3: Linking L1 and L2"
echo "============================================="
echo "Setting L2 executor on L1 callback..."
cast send "$L1_CALLBACK" "setL2Executor(address)" "$L2_EXECUTOR" \
    --private-key "$PRIVATE_KEY" --rpc-url "$L1_RPC" > /dev/null

echo "Setting L1 callback on L2 executor..."
cast send "$L2_EXECUTOR" "setL1Callback(address)" "$L1_CALLBACK" \
    --private-key "$PRIVATE_KEY" --rpc-url "$L2_RPC" > /dev/null

echo "  Done."
echo ""

echo "============================================="
echo " Deployment Complete"
echo "============================================="
echo ""
echo " L1 Contracts:"
echo "   FlashLoanToken (L1):     $L1_TOKEN"
echo "   FlashLoanCallbackL1:     $L1_CALLBACK"
echo ""
echo " L2 Contracts:"
echo "   FlashLoanToken (L2):     $L2_TOKEN"
echo "   FlashLoanPool:           $L2_POOL"
echo "   FlashLoanExecutorL2:     $L2_EXECUTOR"
echo ""
echo "============================================="
