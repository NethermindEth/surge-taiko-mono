#!/bin/sh

# Step 1b: Deploy internal verifiers and configure SurgeVerifier.
# Requires Step 1a (DeployRollupCore) to have been run first.
# Must use the same PRIVATE_KEY as Step 1a.
set -e

# Deployer private key (must match Step 1a)
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# L2 configuration
export L2_CHAIN_ID=${L2_CHAIN_ID:-167004}

# Verifier configuration
export DEPLOY_RISC0_RETH_VERIFIER=${DEPLOY_RISC0_RETH_VERIFIER:-true}
export DEPLOY_SP1_RETH_VERIFIER=${DEPLOY_SP1_RETH_VERIFIER:-true}
export USE_DUMMY_VERIFIER=${USE_DUMMY_VERIFIER:-false}
export DUMMY_VERIFIER_SIGNER=${DUMMY_VERIFIER_SIGNER:-"0x0000000000000000000000000000000000000000"}

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

# Load addresses from Step 1a JSON (if not already set via env vars)
DEPLOY_DIR="$(cd "$(dirname "$0")/../../../.." && pwd)/deployments"
STEP1_JSON="$DEPLOY_DIR/deploy_rollup_core.json"

if [ -z "$SURGE_VERIFIER" ] && [ -f "$STEP1_JSON" ] && command -v jq >/dev/null 2>&1; then
    export SURGE_VERIFIER=$(jq -r '.surge_verifier // empty' "$STEP1_JSON")
fi

if [ -z "$EFFECTIVE_OWNER" ] && [ -f "$STEP1_JSON" ] && command -v jq >/dev/null 2>&1; then
    export EFFECTIVE_OWNER=$(jq -r '.effective_owner // empty' "$STEP1_JSON")
fi

echo "====================================="
echo "Step 1b: Deploy Verifiers"
echo "====================================="
echo ""
echo "  SURGE_VERIFIER:             $SURGE_VERIFIER"
echo "  EFFECTIVE_OWNER:            $EFFECTIVE_OWNER"
echo "  USE_DUMMY_VERIFIER:         $USE_DUMMY_VERIFIER"
echo "  DEPLOY_RISC0_RETH_VERIFIER: $DEPLOY_RISC0_RETH_VERIFIER"
echo "  DEPLOY_SP1_RETH_VERIFIER:   $DEPLOY_SP1_RETH_VERIFIER"
echo ""

if [ "$BROADCAST" = "true" ]; then
    echo "Running in BROADCAST mode - transactions will be executed"
else
    echo "Running in SIMULATION mode - set BROADCAST=true to execute transactions"
fi
echo ""

forge script ./script/layer1/surge/deploy-modular/DeployVerifiers.s.sol:DeployVerifiers \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $VERIFY_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT
