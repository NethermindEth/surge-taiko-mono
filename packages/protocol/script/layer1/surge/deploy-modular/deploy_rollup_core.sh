#!/bin/sh

# Step 1a: Deploy core rollup infrastructure (EmptyImpl, Inbox proxy, SurgeVerifier, optional Timelock)
# This is the first step in the modular Surge L1 deployment.
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Contract owner configuration
export CONTRACT_OWNER=${CONTRACT_OWNER:-"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}

# L2 configuration
export L2_CHAIN_ID=${L2_CHAIN_ID:-167004}

# SurgeVerifier configuration
export NUM_PROOFS_THRESHOLD=${NUM_PROOFS_THRESHOLD:-2}

# Timelock configuration (optional)
export USE_TIMELOCK=${USE_TIMELOCK:-false}
export TIMELOCK_MIN_DELAY=${TIMELOCK_MIN_DELAY:-86400}
export TIMELOCK_MIN_FINALIZATION_STREAK=${TIMELOCK_MIN_FINALIZATION_STREAK:-604800}
export TIMELOCK_PROPOSERS=${TIMELOCK_PROPOSERS:-""}
export TIMELOCK_EXECUTORS=${TIMELOCK_EXECUTORS:-""}

# Forge configuration
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"layer1"}
export VERIFY=${VERIFY:-false}
export BROADCAST=${BROADCAST:-false}
export LOG_LEVEL=${LOG_LEVEL:-"-vvv"}
export BLOCK_GAS_LIMIT=${BLOCK_GAS_LIMIT:-200000000}

BROADCAST_ARG=""
if [ "$BROADCAST" = "true" ]; then
    BROADCAST_ARG="--broadcast"
fi

VERIFY_ARG=""
if [ "$VERIFY" = "true" ]; then
    VERIFY_ARG="--verify"
fi

echo "====================================="
echo "Step 1a: Deploy Rollup Core"
echo "====================================="
echo ""
echo "  CONTRACT_OWNER:       $CONTRACT_OWNER"
echo "  L2_CHAIN_ID:          $L2_CHAIN_ID"
echo "  NUM_PROOFS_THRESHOLD: $NUM_PROOFS_THRESHOLD"
echo "  USE_TIMELOCK:         $USE_TIMELOCK"
echo ""

if [ "$BROADCAST" = "true" ]; then
    echo "Running in BROADCAST mode - transactions will be executed"
else
    echo "Running in SIMULATION mode - set BROADCAST=true to execute transactions"
fi
echo ""

forge script ./script/layer1/surge/deploy-modular/DeployRollupCore.s.sol:DeployRollupCore \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $VERIFY_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT
