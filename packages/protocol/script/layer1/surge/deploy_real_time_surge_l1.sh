#!/bin/sh

# This script deploys the Surge protocol with real-time proving inbox on L1
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Contract owner configuration
export CONTRACT_OWNER=${CONTRACT_OWNER:-"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}

# L2 configuration
export L2_CHAIN_ID=${L2_CHAIN_ID:-167004}

# Verifier deployment flags (only RISC0 and SP1)
export DEPLOY_RISC0_RETH_VERIFIER=${DEPLOY_RISC0_RETH_VERIFIER:-true}
export DEPLOY_SP1_RETH_VERIFIER=${DEPLOY_SP1_RETH_VERIFIER:-true}

# Use dummy verifier for testing (default: false for production, set to true for devnet testing)
export USE_DUMMY_VERIFIER=${USE_DUMMY_VERIFIER:-false}

# Signer address for ProofVerifierDummy (required if USE_DUMMY_VERIFIER=true)
export DUMMY_VERIFIER_SIGNER=${DUMMY_VERIFIER_SIGNER:-"0x0000000000000000000000000000000000000000"}

# SurgeVerifier configuration
# ---------------------------------------------------------------
# Minimum number of distinct proofs required for a transition to finalize
export NUM_PROOFS_THRESHOLD=${NUM_PROOFS_THRESHOLD:-2}

# Inbox configuration
# ---------------------------------------------------------------
# Percentage of basefee paid to coinbase (0-100, default: 75 for devnet)
export BASEFEE_SHARING_PCTG=${BASEFEE_SHARING_PCTG:-75}

# Genesis configuration
# ---------------------------------------------------------------
# The genesis proposal hash to activate the inbox with
export GENESIS_PROPOSAL_HASH=${GENESIS_PROPOSAL_HASH:-"0x0000000000000000000000000000000000000000000000000000000000000001"}

# Deploy configuration
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"layer1"}

# Verify smart contracts
export VERIFY=${VERIFY:-false}

# Broadcast transactions
export BROADCAST=${BROADCAST:-false}

# Parameterize broadcasting
export BROADCAST_ARG=""
if [ "$BROADCAST" = "true" ]; then
    BROADCAST_ARG="--broadcast"
fi

# Parameterize verification
export VERIFY_ARG=""
if [ "$VERIFY" = "true" ]; then
    VERIFY_ARG="--verify"
fi

# Parameterize log level
export LOG_LEVEL=${LOG_LEVEL:-"-vvv"}

# Parameterize block gas limit
export BLOCK_GAS_LIMIT=${BLOCK_GAS_LIMIT:-200000000}

forge script ./script/layer1/surge/DeployRealTimeSurgeL1.s.sol:DeployRealTimeSurgeL1 \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $VERIFY_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT
