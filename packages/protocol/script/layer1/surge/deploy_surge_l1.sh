#!/bin/sh

# This script deploys the Surge protocol on L1
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Contract owner configuration
export CONTRACT_OWNER=${CONTRACT_OWNER:-"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}

# L2 configuration
export L2_CHAIN_ID=${L2_CHAIN_ID:-167004}
export L2_GENESIS_HASH=${L2_GENESIS_HASH:-"0xee1950562d42f0da28bd4550d88886bc90894c77c9c9eaefef775d4c8223f259"}

# Verifier deployment flags (only RISC0 and SP1)
export DEPLOY_RISC0_RETH_VERIFIER=${DEPLOY_RISC0_RETH_VERIFIER:-true}
export DEPLOY_SP1_RETH_VERIFIER=${DEPLOY_SP1_RETH_VERIFIER:-true}

# Use dummy verifier for testing (default: false for production, set to true for devnet testing)
export USE_DUMMY_VERIFIER=${USE_DUMMY_VERIFIER:-false}

# Protocol mode: true for Taiko (InboxOptimized2, CodecOptimized, AnyTwoVerifier)
#                false for Surge (SurgeInbox, SurgeCodecSimple, SurgeVerifier)
export USE_TAIKO=${USE_TAIKO:-false}

# Inbox configuration
# ---------------------------------------------------------------
# Bond token address (Use 0-address for ETH)
export BOND_TOKEN=${BOND_TOKEN:-"0x0000000000000000000000000000000000000000"}

# Proving window in seconds (default: 2 hours)
export PROVING_WINDOW=${PROVING_WINDOW:-7200}

# Extended proving window in seconds (default: 4 hours)
export EXTENDED_PROVING_WINDOW=${EXTENDED_PROVING_WINDOW:-14400}

# Maximum number of proposals that can be finalized in one transaction
export MAX_FINALIZATION_COUNT=${MAX_FINALIZATION_COUNT:-16}

# Finalization grace period in seconds (default: 768 seconds = 2 epochs)
export FINALIZATION_GRACE_PERIOD=${FINALIZATION_GRACE_PERIOD:-768}

# Ring buffer size for storing proposal hashes
export RING_BUFFER_SIZE=${RING_BUFFER_SIZE:-100}

# Percentage of basefee paid to coinbase (0-100, default: 75 for devnet)
export BASEFEE_SHARING_PCTG=${BASEFEE_SHARING_PCTG:-75}

# Minimum number of forced inclusions to process if due
export MIN_FORCED_INCLUSION_COUNT=${MIN_FORCED_INCLUSION_COUNT:-1}

# Delay for forced inclusions in seconds (default: 0 for devnet)
export FORCED_INCLUSION_DELAY=${FORCED_INCLUSION_DELAY:-0}

# Base fee for forced inclusions in Gwei (default: 10,000,000 = 0.01 ETH)
export FORCED_INCLUSION_FEE_IN_GWEI=${FORCED_INCLUSION_FEE_IN_GWEI:-10000000}

# Queue size at which the fee doubles
export FORCED_INCLUSION_FEE_DOUBLE_THRESHOLD=${FORCED_INCLUSION_FEE_DOUBLE_THRESHOLD:-50}

# Minimum delay between checkpoints in seconds (default: 384 seconds = 1 epoch)
export MIN_CHECKPOINT_DELAY=${MIN_CHECKPOINT_DELAY:-384}

# Multiplier for permissionless inclusion window
export PERMISSIONLESS_INCLUSION_MULTIPLIER=${PERMISSIONLESS_INCLUSION_MULTIPLIER:-5}

# Composite key version for proof verification
export COMPOSITE_KEY_VERSION=${COMPOSITE_KEY_VERSION:-1}

# Finality gadget configuration
# ---------------------------------------------------------------
# Optimistic fallback delay in seconds (default: 7 days)
# Delay before a single proof with no conflicts can finalize a transition
export OPTIMISTIC_FALLBACK_DELAY=${OPTIMISTIC_FALLBACK_DELAY:-604800}

# Minimum number of distinct proofs required for a transition to finalize immediately
export FINALISING_PROOF_COUNT=${FINALISING_PROOF_COUNT:-2}

# Deploy Surge protocol
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
export LOG_LEVEL=${LOG_LEVEL:-"-vvvv"}

# Parameterize block gas limit
export BLOCK_GAS_LIMIT=${BLOCK_GAS_LIMIT:-200000000}

forge script ./script/layer1/surge/DeploySurgeL1.s.sol:DeploySurgeL1 \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $VERIFY_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT
