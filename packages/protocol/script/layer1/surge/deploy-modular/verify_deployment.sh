#!/bin/sh

# Step 1e: Read-only verification of all registrations and ownership state.
# Requires Steps 1a-1d to have been run first.
# This script does NOT broadcast any transactions.
set -e

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# L2 configuration
export L2_CHAIN_ID=${L2_CHAIN_ID:-167004}

# Verifier configuration (affects ownership verification)
export USE_DUMMY_VERIFIER=${USE_DUMMY_VERIFIER:-false}

# Forge configuration
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"layer1"}
export LOG_LEVEL=${LOG_LEVEL:-"-vvv"}

# Load addresses from step JSON files (if not already set via env vars)
DEPLOY_DIR="$(cd "$(dirname "$0")/../../../.." && pwd)/deployments"

# From Step 1a
STEP1_JSON="$DEPLOY_DIR/deploy_rollup_core.json"
if [ -f "$STEP1_JSON" ] && command -v jq >/dev/null 2>&1; then
    [ -z "$SURGE_INBOX" ] && export SURGE_INBOX=$(jq -r '.surge_inbox // empty' "$STEP1_JSON")
    [ -z "$SURGE_VERIFIER" ] && export SURGE_VERIFIER=$(jq -r '.surge_verifier // empty' "$STEP1_JSON")
    [ -z "$EFFECTIVE_OWNER" ] && export EFFECTIVE_OWNER=$(jq -r '.effective_owner // empty' "$STEP1_JSON")
fi

# From Step 1b (may not exist if verifiers not deployed)
STEP2_JSON="$DEPLOY_DIR/deploy_verifiers.json"
if [ -f "$STEP2_JSON" ] && command -v jq >/dev/null 2>&1; then
    [ -z "$RISC0_VERIFIER" ] && export RISC0_VERIFIER=$(jq -r '.risc0_verifier // empty' "$STEP2_JSON")
    [ -z "$SP1_VERIFIER" ] && export SP1_VERIFIER=$(jq -r '.sp1_verifier // empty' "$STEP2_JSON")
fi

# Default to zero address if verifiers not deployed
export RISC0_VERIFIER=${RISC0_VERIFIER:-"0x0000000000000000000000000000000000000000"}
export SP1_VERIFIER=${SP1_VERIFIER:-"0x0000000000000000000000000000000000000000"}

# From Step 1c
STEP3_JSON="$DEPLOY_DIR/deploy_shared_contracts.json"
if [ -f "$STEP3_JSON" ] && command -v jq >/dev/null 2>&1; then
    [ -z "$SHARED_RESOLVER" ] && export SHARED_RESOLVER=$(jq -r '.shared_resolver // empty' "$STEP3_JSON")
    [ -z "$SIGNAL_SERVICE" ] && export SIGNAL_SERVICE=$(jq -r '.signal_service // empty' "$STEP3_JSON")
    [ -z "$BRIDGE" ] && export BRIDGE=$(jq -r '.bridge // empty' "$STEP3_JSON")
    [ -z "$ERC20_VAULT" ] && export ERC20_VAULT=$(jq -r '.erc20_vault // empty' "$STEP3_JSON")
    [ -z "$ERC721_VAULT" ] && export ERC721_VAULT=$(jq -r '.erc721_vault // empty' "$STEP3_JSON")
    [ -z "$ERC1155_VAULT" ] && export ERC1155_VAULT=$(jq -r '.erc1155_vault // empty' "$STEP3_JSON")
fi

echo "====================================="
echo "Step 1e: Verify Deployment"
echo "====================================="
echo ""
echo "  SURGE_INBOX:      $SURGE_INBOX"
echo "  SURGE_VERIFIER:   $SURGE_VERIFIER"
echo "  EFFECTIVE_OWNER:  $EFFECTIVE_OWNER"
echo "  SHARED_RESOLVER:  $SHARED_RESOLVER"
echo "  SIGNAL_SERVICE:   $SIGNAL_SERVICE"
echo "  BRIDGE:           $BRIDGE"
echo "  ERC20_VAULT:      $ERC20_VAULT"
echo "  ERC721_VAULT:     $ERC721_VAULT"
echo "  ERC1155_VAULT:    $ERC1155_VAULT"
echo "  RISC0_VERIFIER:   $RISC0_VERIFIER"
echo "  SP1_VERIFIER:     $SP1_VERIFIER"
echo ""

forge script ./script/layer1/surge/deploy-modular/VerifyDeployment.s.sol:VerifyDeployment \
    --fork-url $FORK_URL \
    --ffi \
    $LOG_LEVEL
