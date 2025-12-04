#!/bin/sh

# This script accepts ownership of multiple contracts using the AcceptOwnership script.
set -e

# Deployer private key (must be the pending owner)
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Contract addresses to accept ownership (comma-separated)
export CONTRACT_ADDRESSES=${CONTRACT_ADDRESSES:-""}

# Broadcast transactions
export BROADCAST=${BROADCAST:-false}

# Parameterize broadcasting
export BROADCAST_ARG=""
if [ "$BROADCAST" = "true" ]; then
    BROADCAST_ARG="--broadcast"
fi

# Parameterize log level
export LOG_LEVEL=${LOG_LEVEL:-"-vvvv"}

echo "Contract addresses to accept ownership:"
echo "$CONTRACT_ADDRESSES"
echo ""

if [ "$BROADCAST" = "true" ]; then
    echo "Running in BROADCAST mode - transactions will be executed"
else
    echo "Running in SIMULATION mode - set BROADCAST=true to execute transactions"
fi
echo ""

forge script ./script/layer1/surge/AcceptOwnership.s.sol:AcceptOwnership \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY
