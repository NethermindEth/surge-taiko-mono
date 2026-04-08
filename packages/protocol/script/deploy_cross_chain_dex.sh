#!/bin/sh

# =============================================================================
# Cross-Chain DEX Full Deployment Script
# =============================================================================
# Deploys the complete Cross-Chain DEX system on L1 and L2, links the vaults,
# and outputs a summary of deployed addresses.
#
# Usage:
#   PRIVATE_KEY=0x... ./script/deploy_cross_chain_dex.sh
#
# Environment variables (all have sensible defaults except PRIVATE_KEY):
#   PRIVATE_KEY      - Deployer private key (required)
#   L1_RPC           - L1 RPC URL
#   L2_RPC           - L2 RPC URL
#   L1_BRIDGE        - L1 Bridge contract address
#   L2_BRIDGE        - L2 Bridge contract address
#   SWAP_TOKEN       - Existing L1 ERC20 address (if empty, deploys a new SwapToken)
#   TOKEN_DECIMALS   - Token decimals (default: 18; set to 6 for real USDC)
#   INITIAL_TOKEN_SUPPLY - Only used when deploying new token (default: 1M tokens)
#   LOG_LEVEL        - Forge verbosity (default: -vvvv)
#   DRY_RUN          - Set to "true" to simulate without broadcasting
# =============================================================================

set -e

# ---------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------

export PRIVATE_KEY=${PRIVATE_KEY:-""}

if [ -z "$PRIVATE_KEY" ]; then
    echo "ERROR: PRIVATE_KEY is required"
    echo "Usage: PRIVATE_KEY=0x... ./script/deploy_cross_chain_dex.sh"
    exit 1
fi

export L1_RPC=${L1_RPC:-""}
export L2_RPC=${L2_RPC:-""}

export L1_BRIDGE=${L1_BRIDGE:-"0xc1e59A201cE4CD58590FC3Ab45081921cF186550"}
export L2_BRIDGE=${L2_BRIDGE:-"0x7633740000000000000000000000000000000001"}

export SWAP_TOKEN=${SWAP_TOKEN:-""}
export TOKEN_DECIMALS=${TOKEN_DECIMALS:-"6"} # (Default: 6 for USDC)

# Only specify if SWAP_TOKEN is empty
export INITIAL_TOKEN_SUPPLY=${INITIAL_TOKEN_SUPPLY:-"1000000000000000000000000"}

export LOG_LEVEL=${LOG_LEVEL:-"-vvvv"}
export DRY_RUN=${DRY_RUN:-"false"}

DEPLOY_DIR="deployments"
L1_DEPLOY_JSON="$DEPLOY_DIR/cross-chain-dex-l1.json"
L2_DEPLOY_JSON="$DEPLOY_DIR/cross-chain-dex-l2.json"

# ---------------------------------------------------------------
# Derived values
# ---------------------------------------------------------------

echo "============================================="
echo " Cross-Chain DEX Full Deployment"
echo "============================================="
echo ""

echo "Resolving configuration..."
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
if [ -n "$SWAP_TOKEN" ]; then
    echo "  Swap Token:    $SWAP_TOKEN (existing)"
else
    echo "  Swap Token:    (new token will be deployed)"
    echo "  Token Supply:  $INITIAL_TOKEN_SUPPLY"
fi
echo "  Token Decimals: $TOKEN_DECIMALS"
echo ""

if [ "$DRY_RUN" = "true" ]; then
    echo "  *** DRY RUN — no transactions will be broadcast ***"
    BROADCAST_ARG=""
else
    BROADCAST_ARG="--broadcast"
fi
echo ""

# ---------------------------------------------------------------
# Step 1: Deploy L1 contracts
# ---------------------------------------------------------------

echo "============================================="
echo " Step 1/4: Deploying L1 contracts"
echo "============================================="
echo ""

export L2_CHAIN_ID
export FOUNDRY_PROFILE="layer1"

forge script ./script/layer1/surge/cross-chain-dex/DeployCrossChainDexL1.s.sol:DeployCrossChainDexL1 \
    --fork-url "$L1_RPC" \
    $BROADCAST_ARG \
    $LOG_LEVEL \
    --private-key "$PRIVATE_KEY"

if [ "$DRY_RUN" = "true" ]; then
    echo ""
    echo "DRY RUN: Skipping remaining steps (L2 deploy, linking, verification)"
    exit 0
fi

# Parse L1 deployment output
L1_VAULT=$(python3 -c "import json; print(json.load(open('$L1_DEPLOY_JSON'))['CrossChainSwapVaultL1'])")
L1_TOKEN=$(python3 -c "import json; print(json.load(open('$L1_DEPLOY_JSON'))['SwapToken'])")

echo ""
echo "  L1 Vault:  $L1_VAULT"
echo "  L1 Token:  $L1_TOKEN"
echo ""

# ---------------------------------------------------------------
# Step 2: Deploy L2 contracts
# ---------------------------------------------------------------

echo "============================================="
echo " Step 2/4: Deploying L2 contracts"
echo "============================================="
echo ""

export L1_CHAIN_ID
export FOUNDRY_PROFILE="layer2"

forge script ./script/layer2/surge/cross-chain-dex/DeployCrossChainDexL2.s.sol:DeployCrossChainDexL2 \
    --fork-url "$L2_RPC" \
    $BROADCAST_ARG \
    $LOG_LEVEL \
    --private-key "$PRIVATE_KEY"

# Parse L2 deployment output
L2_VAULT=$(python3 -c "import json; print(json.load(open('$L2_DEPLOY_JSON'))['CrossChainSwapVaultL2'])")
L2_DEX=$(python3 -c "import json; print(json.load(open('$L2_DEPLOY_JSON'))['SimpleDEX'])")
L2_TOKEN=$(python3 -c "import json; print(json.load(open('$L2_DEPLOY_JSON'))['SwapTokenL2'])")

echo ""
echo "  L2 Vault:  $L2_VAULT"
echo "  L2 DEX:    $L2_DEX"
echo "  L2 Token:  $L2_TOKEN"
echo ""

# ---------------------------------------------------------------
# Step 3: Link vaults (L1 <-> L2)
# ---------------------------------------------------------------

echo "============================================="
echo " Step 3/4: Linking vaults"
echo "============================================="
echo ""

echo "Setting L2 vault on L1 vault..."
cast send "$L1_VAULT" "setL2Vault(address)" "$L2_VAULT" \
    --private-key "$PRIVATE_KEY" \
    --rpc-url "$L1_RPC" > /dev/null

echo "Setting L1 vault on L2 vault..."
cast send "$L2_VAULT" "setL1Vault(address)" "$L1_VAULT" \
    --private-key "$PRIVATE_KEY" \
    --rpc-url "$L2_RPC" > /dev/null

echo "  Done."
echo ""

# ---------------------------------------------------------------
# Step 4: Verify
# ---------------------------------------------------------------

echo "============================================="
echo " Step 4/4: Verifying deployment"
echo "============================================="
echo ""

L1_VAULT_L2=$(cast call "$L1_VAULT" "l2Vault()(address)" --rpc-url "$L1_RPC")
L2_VAULT_L1=$(cast call "$L2_VAULT" "l1Vault()(address)" --rpc-url "$L2_RPC")
L2_MINTER=$(cast call "$L2_TOKEN" "minter()(address)" --rpc-url "$L2_RPC")
L2_LP=$(cast call "$L2_DEX" "liquidityProvider()(address)" --rpc-url "$L2_RPC")

ERRORS=0

if [ "$L1_VAULT_L2" != "$L2_VAULT" ]; then
    echo "  ERROR: L1 vault l2Vault mismatch: $L1_VAULT_L2 != $L2_VAULT"
    ERRORS=$((ERRORS + 1))
else
    echo "  OK: L1 Vault -> L2 Vault linked"
fi

if [ "$L2_VAULT_L1" != "$L1_VAULT" ]; then
    echo "  ERROR: L2 vault l1Vault mismatch: $L2_VAULT_L1 != $L1_VAULT"
    ERRORS=$((ERRORS + 1))
else
    echo "  OK: L2 Vault -> L1 Vault linked"
fi

if [ "$L2_MINTER" != "$L2_VAULT" ]; then
    echo "  ERROR: SwapTokenL2 minter mismatch: $L2_MINTER != $L2_VAULT"
    ERRORS=$((ERRORS + 1))
else
    echo "  OK: SwapTokenL2 minter = L2 Vault"
fi

if [ "$L2_LP" != "$L2_VAULT" ]; then
    echo "  ERROR: SimpleDEX liquidityProvider mismatch: $L2_LP != $L2_VAULT"
    ERRORS=$((ERRORS + 1))
else
    echo "  OK: SimpleDEX liquidityProvider = L2 Vault"
fi

echo ""

if [ "$ERRORS" -gt 0 ]; then
    echo "  DEPLOYMENT VERIFICATION FAILED ($ERRORS errors)"
    exit 1
fi

# ---------------------------------------------------------------
# Summary
# ---------------------------------------------------------------

echo "============================================="
echo " Deployment Complete"
echo "============================================="
echo ""
echo " L1 Contracts:"
echo "   SwapToken:              $L1_TOKEN"
echo "   CrossChainSwapVaultL1:  $L1_VAULT"
echo ""
echo " L2 Contracts:"
echo "   SwapTokenL2:            $L2_TOKEN"
echo "   SimpleDEX:              $L2_DEX"
echo "   CrossChainSwapVaultL2:  $L2_VAULT"
echo ""
echo " Configuration:"
echo "   L1 Bridge:              $L1_BRIDGE"
echo "   L2 Bridge:              $L2_BRIDGE"
echo "   Token Decimals:         $TOKEN_DECIMALS"
echo ""
echo "============================================="
