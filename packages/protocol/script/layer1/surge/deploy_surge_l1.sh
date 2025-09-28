#!/bin/sh

# This script deploys the Surge protocol on L1
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Owner and executor configuration
export USE_TIMELOCKED_OWNER=${USE_TIMELOCKED_OWNER:-true}
# (Remaining are only required if USE_TIMELOCKED_OWNER is true)
export OWNER_MULTISIG=${OWNER_MULTISIG:-"0x60997970C51812dc3A010C7d01b50e0d17dc79C8"}
export OWNER_MULTISIG_SIGNERS=${OWNER_MULTISIG_SIGNERS:-"0x1237810000000000000000000000000000000002,0x1237810000000000000000000000000000000003,0x1237810000000000000000000000000000000004"}
export TIMELOCK_PERIOD=${TIMELOCK_PERIOD:-86400}

# DAO configuration
export DAO=${DAO:-"0x1237810000000000000000000000000000000001"}

# L2 configuration
export L2_NETWORK=${L2_NETWORK:-"devnet"}
export L2_CHAINID=${L2_CHAINID:-167004}
export L2_GENESIS_HASH=${L2_GENESIS_HASH:-"0xee1950562d42f0da28bd4550d88886bc90894c77c9c9eaefef775d4c8223f259"}

# Liveness configuration
export MAX_VERIFICATION_DELAY=${MAX_VERIFICATION_DELAY:-100}
export MIN_VERIFICATION_STREAK=${MIN_VERIFICATION_STREAK:-10}
export LIVENESS_BOND_BASE=${LIVENESS_BOND_BASE:-1000000000000000000}
export COOLDOWN_WINDOW=${COOLDOWN_WINDOW:-604800}

# Preconf configuration
export USE_PRECONF=${USE_PRECONF:-false}
# (Remaining are only required if USE_PRECONF is true)
export FALLBACK_PRECONF=${FALLBACK_PRECONF:-"0x0000000000000000000000000000000000000000"}

# Forced inclusion configuration
export INCLUSION_WINDOW=${INCLUSION_WINDOW:-24}
export INCLUSION_FEE_IN_GWEI=${INCLUSION_FEE_IN_GWEI:-100}

# Verifier deployment flags
export DEPLOY_SGX_RETH_VERIFIER=${DEPLOY_SGX_RETH_VERIFIER:-true}
export DEPLOY_SGX_GETH_VERIFIER=${DEPLOY_SGX_GETH_VERIFIER:-false}
export DEPLOY_AZURE_TDX_VERIFIER=${DEPLOY_AZURE_TDX_VERIFIER:-true}
export DEPLOY_RISC0_RETH_VERIFIER=${DEPLOY_RISC0_RETH_VERIFIER:-true}
export DEPLOY_SP1_RETH_VERIFIER=${DEPLOY_SP1_RETH_VERIFIER:-true}

# Note: Verifier configuration is now handled in separate setup scripts

# TDX Automata contract addresses
export TDX_PCS_DAO_ADDRESS=${TDX_PCS_DAO_ADDRESS:-"0x45CF7485A0D394130153a3630EA0729999511C2e"}
export TDX_FMSPC_TCB_DAO_ADDRESS=${TDX_FMSPC_TCB_DAO_ADDRESS:-"0x63eF330eAaadA189861144FCbc9176dae41A5BAf"}
export TDX_ENCLAVE_IDENTITY_DAO_ADDRESS=${TDX_ENCLAVE_IDENTITY_DAO_ADDRESS:-"0xc3ea5Ff40263E16cD2f4413152A77e7A6b10B0C9"}
export TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS=${TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS:-"0x635A8A01e84cDcE1475FCeB7D57FEcadD3d1a0A0"}

# TDX collateral paths
export TDX_QE_IDENTITY_PATH=${TDX_QE_IDENTITY_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_identity.json"}
export TDX_TCB_INFO_PATH=${TDX_TCB_INFO_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_tcb_00806F050000.json"}
export TDX_PCS_CERT_PATH=${TDX_PCS_CERT_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex"}

# Deploy Surge protocol
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"layer1"}

# Broadcast transactions
export BROADCAST=${BROADCAST:-false}

# Parameterize broadcasting
export BROADCAST_ARG=""
if [ "$BROADCAST" = "true" ]; then
    BROADCAST_ARG="--broadcast"
fi

# Parameterize log level
export LOG_LEVEL=${LOG_LEVEL:-"-vvvv"}

# Parameterize block gas limit
export BLOCK_GAS_LIMIT=${BLOCK_GAS_LIMIT:-200000000}

forge script ./script/layer1/surge/DeploySurgeL1.s.sol:DeploySurgeL1 \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT 