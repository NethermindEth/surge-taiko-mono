#!/bin/sh

# Step 1d: Deploy SurgeInbox implementation, upgrade proxy, and initialize.
# Requires Step 1a (DeployRollupCore) and Step 1c (DeploySharedContracts) to have been run first.
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Bond configuration
export BOND_TOKEN=${BOND_TOKEN:-"0x0000000000000000000000000000000000000000"}
export MIN_BOND=${MIN_BOND:-0}
export LIVENESS_BOND=${LIVENESS_BOND:-128000000000}
export WITHDRAWAL_DELAY=${WITHDRAWAL_DELAY:-3600}

# Inbox configuration
export PROVING_WINDOW=${PROVING_WINDOW:-7200}
export MAX_PROOF_SUBMISSION_DELAY=${MAX_PROOF_SUBMISSION_DELAY:-14400}
export RING_BUFFER_SIZE=${RING_BUFFER_SIZE:-16000}
export BASEFEE_SHARING_PCTG=${BASEFEE_SHARING_PCTG:-75}
export MIN_FORCED_INCLUSION_COUNT=${MIN_FORCED_INCLUSION_COUNT:-1}
export FORCED_INCLUSION_DELAY=${FORCED_INCLUSION_DELAY:-0}
export FORCED_INCLUSION_FEE_IN_GWEI=${FORCED_INCLUSION_FEE_IN_GWEI:-10000000}
export FORCED_INCLUSION_FEE_DOUBLE_THRESHOLD=${FORCED_INCLUSION_FEE_DOUBLE_THRESHOLD:-50}
export MIN_CHECKPOINT_DELAY=${MIN_CHECKPOINT_DELAY:-384}
export PERMISSIONLESS_INCLUSION_MULTIPLIER=${PERMISSIONLESS_INCLUSION_MULTIPLIER:-5}

# Finalization streak configuration
export MAX_FINALIZATION_DELAY_BEFORE_STREAK_RESET=${MAX_FINALIZATION_DELAY_BEFORE_STREAK_RESET:-3600}

# Rollback configuration
export MAX_FINALIZATION_DELAY_BEFORE_ROLLBACK=${MAX_FINALIZATION_DELAY_BEFORE_ROLLBACK:-604800}

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

if [ -z "$SURGE_INBOX" ] && [ -f "$STEP1_JSON" ] && command -v jq >/dev/null 2>&1; then
    export SURGE_INBOX=$(jq -r '.surge_inbox // empty' "$STEP1_JSON")
fi

if [ -z "$SURGE_VERIFIER" ] && [ -f "$STEP1_JSON" ] && command -v jq >/dev/null 2>&1; then
    export SURGE_VERIFIER=$(jq -r '.surge_verifier // empty' "$STEP1_JSON")
fi

if [ -z "$EFFECTIVE_OWNER" ] && [ -f "$STEP1_JSON" ] && command -v jq >/dev/null 2>&1; then
    export EFFECTIVE_OWNER=$(jq -r '.effective_owner // empty' "$STEP1_JSON")
fi

# Load addresses from Step 1c JSON
STEP3_JSON="$DEPLOY_DIR/deploy_shared_contracts.json"

if [ -z "$SIGNAL_SERVICE" ] && [ -f "$STEP3_JSON" ] && command -v jq >/dev/null 2>&1; then
    export SIGNAL_SERVICE=$(jq -r '.signal_service // empty' "$STEP3_JSON")
fi

echo "====================================="
echo "Step 1d: Setup Inbox"
echo "====================================="
echo ""
echo "  SURGE_INBOX:      $SURGE_INBOX"
echo "  SURGE_VERIFIER:   $SURGE_VERIFIER"
echo "  SIGNAL_SERVICE:   $SIGNAL_SERVICE"
echo "  EFFECTIVE_OWNER:  $EFFECTIVE_OWNER"
echo ""

if [ "$BROADCAST" = "true" ]; then
    echo "Running in BROADCAST mode - transactions will be executed"
else
    echo "Running in SIMULATION mode - set BROADCAST=true to execute transactions"
fi
echo ""

forge script ./script/layer1/surge/deploy-modular/SetupInbox.s.sol:SetupInbox \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $VERIFY_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT
