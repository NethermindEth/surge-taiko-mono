#!/bin/sh

# This script sets up the Cross-Chain DEX by linking L1Vault into L2Vault.
# Run this on L2 after both L1 and L2 contracts are deployed.
set -e

# Foundry keystore account for the signer (must be the L2Vault admin)
export ACCOUNT=${ACCOUNT:-"surge_gnosis_deployer"}
export PASSWORD_FILE=${PASSWORD_FILE:-"/tmp/.keystore-pw"}

# Network configuration
export L2_RPC=${L2_RPC:-"https://rpc.realtime.surge.wtf"}

# Vault addresses (must be set after deployment)
export L1_VAULT=${L1_VAULT:-""}
export L2_VAULT=${L2_VAULT:-""}

if [ -z "$L1_VAULT" ] || [ -z "$L2_VAULT" ]; then
    echo "ERROR: L1_VAULT and L2_VAULT must be set"
    echo "Usage: L1_VAULT=0x... L2_VAULT=0x... ./setup_cross_chain_dex_l2.sh"
    exit 1
fi

# Broadcast transactions
export BROADCAST=${BROADCAST:-false}

export BROADCAST_ARG=""
if [ "$BROADCAST" = "true" ]; then
    BROADCAST_ARG="--broadcast"
fi

export LOG_LEVEL=${LOG_LEVEL:-"-vvvv"}
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"layer2"}

echo "=====================================";
echo "Setting up Cross-Chain DEX (L2)";
echo "=====================================";
echo "L2 RPC: $L2_RPC"
echo "L1 Vault: $L1_VAULT"
echo "L2 Vault: $L2_VAULT"
echo ""

if [ "$BROADCAST" = "true" ]; then
    echo "Running in BROADCAST mode - transactions will be executed"
else
    echo "Running in SIMULATION mode - set BROADCAST=true to execute transactions"
fi
echo ""

forge script ./script/layer2/surge/cross-chain-dex/SetupCrossChainDexL2.s.sol:SetupCrossChainDexL2 \
    --fork-url $L2_RPC \
    $BROADCAST_ARG \
    $LOG_LEVEL \
    --account $ACCOUNT \
    --password-file $PASSWORD_FILE
