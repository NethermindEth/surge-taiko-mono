#!/bin/sh

# This script sets up the Surge protocol on L2
set -e

# Deployer private key
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# Network configuration
export FORK_URL=http://localhost:8545

# L1 configuration
export L1_CHAINID=1
export L1_BRIDGE=0x0000000000000000000000000000000000000000
export L1_SIGNAL_SERVICE=0x0000000000000000000000000000000000000000
export L1_ERC20_VAULT=0x0000000000000000000000000000000000000000
export L1_ERC721_VAULT=0x0000000000000000000000000000000000000000
export L1_ERC1155_VAULT=0x0000000000000000000000000000000000000000

# Owner and executor configuration
export OWNER_MULTISIG=0x60997970C51812dc3A010C7d01b50e0d17dc79C8
export OWNER_MULTISIG_SIGNERS="0x1237810000000000000000000000000000000002,0x1237810000000000000000000000000000000003,0x1237810000000000000000000000000000000004"
export TIMELOCK_PERIOD=86400

# Deploy Surge protocol
export FOUNDRY_PROFILE="layer2"

# Broadcast transactions
export BROADCAST=false

# Parameterize broadcasting
export BROADCAST_ARG=""
if [ "$BROADCAST" = "true" ]; then
    BROADCAST_ARG="--broadcast"
fi

forge script ./script/layer1/surge/SetupSurgeL2.s.sol:SetupSurgeL2 \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    --ffi \
    -vvvv \
    --private-key $PRIVATE_KEY \
    --block-gas-limit 200000000 