#!/bin/sh

# This script sets up the TDX verifier after deployment
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Required verifier configuration
# TDX Automata contract addresses
export TDX_PCS_DAO_ADDRESS=${TDX_PCS_DAO_ADDRESS:-"0x45CF7485A0D394130153a3630EA0729999511C2e"}
export TDX_FMSPC_TCB_DAO_ADDRESS=${TDX_FMSPC_TCB_DAO_ADDRESS:-"0x63eF330eAaadA189861144FCbc9176dae41A5BAf"}
export TDX_ENCLAVE_IDENTITY_DAO_ADDRESS=${TDX_ENCLAVE_IDENTITY_DAO_ADDRESS:-"0xc3ea5Ff40263E16cD2f4413152A77e7A6b10B0C9"}
export TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS=${TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS:-"0x635A8A01e84cDcE1475FCeB7D57FEcadD3d1a0A0"}

# TDX collateral paths
export TDX_QE_IDENTITY_PATH=${TDX_QE_IDENTITY_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_identity.json"}
export TDX_TCB_INFO_PATH=${TDX_TCB_INFO_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_tcb_00806F050000.json"}
export TDX_PCS_CERT_PATH=${TDX_PCS_CERT_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex"}

export TDX_VERIFIER_ADDRESS=${TDX_VERIFIER_ADDRESS:-""}

# This should be the L1 timelock controller
export NEW_OWNER=${NEW_OWNER:-""}

# TDX configuration (optional)
export TDX_TRUSTED_PARAMS_BYTES=${TDX_TRUSTED_PARAMS_BYTES:-""}
export TDX_QUOTE_BYTES=${TDX_QUOTE_BYTES:-""}

# TDX collateral paths (optional)
export TDX_PCS_CERT_PATH=${TDX_PCS_CERT_PATH:-""}
export TDX_QE_IDENTITY_PATH=${TDX_QE_IDENTITY_PATH:-""}
export TDX_TCB_INFO_PATH=${TDX_TCB_INFO_PATH:-""}

# Foundry profile
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"layer1"}

# Broadcast transactions
export BROADCAST=${BROADCAST:-true}

# Verify smart contracts
export VERIFY=${VERIFY:-false}

# Required environment variable validation
if [ -z "$TDX_VERIFIER_ADDRESS" ]; then
    echo "Error: TDX_VERIFIER_ADDRESS not set"
    exit 1
fi

if [ -z "$TDX_PCS_DAO_ADDRESS" ]; then
    echo "Error: TDX_PCS_DAO_ADDRESS not set"
    exit 1
fi

if [ -z "$TDX_FMSPC_TCB_DAO_ADDRESS" ]; then
    echo "Error: TDX_FMSPC_TCB_DAO_ADDRESS not set"
    exit 1
fi

if [ -z "$TDX_ENCLAVE_IDENTITY_DAO_ADDRESS" ]; then
    echo "Error: TDX_ENCLAVE_IDENTITY_DAO_ADDRESS not set"
    exit 1
fi

if [ -z "$TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS" ]; then
    echo "Error: TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS not set"
    exit 1
fi

if [ -z "$NEW_OWNER" ]; then
    echo "Error: NEW_OWNER not set (should be timelock controller address)"
    exit 1
fi

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

echo "Setting up TDX verifier..."
echo "TDX Verifier address: $TDX_VERIFIER_ADDRESS"
echo "TDX PCS DAO address: $TDX_PCS_DAO_ADDRESS"
echo "TDX FMSPC TCB DAO address: $TDX_FMSPC_TCB_DAO_ADDRESS"
echo "TDX Enclave Identity DAO address: $TDX_ENCLAVE_IDENTITY_DAO_ADDRESS"
echo "TDX Enclave Identity Helper address: $TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS"
echo "New owner: $NEW_OWNER"

if [ -n "$TDX_TRUSTED_PARAMS_BYTES" ]; then
    echo "TDX Trusted Params provided"
fi

if [ -n "$TDX_PCS_CERT_PATH" ]; then
    echo "TDX PCS Certificate path: $TDX_PCS_CERT_PATH"
fi

if [ -n "$TDX_QE_IDENTITY_PATH" ]; then
    echo "TDX QE Identity path: $TDX_QE_IDENTITY_PATH"
fi

if [ -n "$TDX_TCB_INFO_PATH" ]; then
    echo "TDX TCB Info path: $TDX_TCB_INFO_PATH"
fi

if [ -n "$TDX_QUOTE_BYTES" ]; then
    echo "TDX Quote bytes provided for instance registration"
fi

# Run the setup script
forge script script/layer1/surge/SetupTDXVerifier.s.sol \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $VERIFY_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT

echo "TDX verifier setup completed successfully!"