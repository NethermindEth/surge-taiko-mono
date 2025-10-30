#!/bin/sh

# This script sets up the TDX verifier after deployment
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# Required verifier configuration
# TDX Automata contract addresses
export TDX_PCS_DAO_ADDRESS=${TDX_PCS_DAO_ADDRESS:-"0x928826C6D0d1986bD0465697984fa3722ADE16E1"}
export TDX_FMSPC_TCB_DAO_ADDRESS=${TDX_FMSPC_TCB_DAO_ADDRESS:-"0xA8C0F6F6Deb3dA48Be03A99C112737000a5a3088"}
export TDX_ENCLAVE_IDENTITY_DAO_ADDRESS=${TDX_ENCLAVE_IDENTITY_DAO_ADDRESS:-"0x5d1122a0d55b5095C0f03FBEa106A2e9722cb13F"}
export TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS=${TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS:-"0x870D17e2C12aF1C47dD1f0e4aFd36e28c830D558"}

# TDX collateral paths
export TDX_QE_IDENTITY_PATH=${TDX_QE_IDENTITY_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_identity.json"}
export TDX_TCB_INFO_PATH=${TDX_TCB_INFO_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_tcb_90c06f000000.json"}
export TDX_PCS_CERT_PATH=${TDX_PCS_CERT_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex"}
export TDX_ROOT_PCS_CERT_PATH=${TDX_ROOT_PCS_CERT_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_root_pcs_cert.hex"}

export AZURE_TDX_VERIFIER_ADDRESS=${AZURE_TDX_VERIFIER_ADDRESS:-""}

# This should be the L1 timelock controller
export NEW_OWNER=${NEW_OWNER:-""}

# TDX configuration (optional)
export AZURE_TDX_TRUSTED_PARAMS_BYTES=${AZURE_TDX_TRUSTED_PARAMS_BYTES:-""}
export AZURE_TDX_QUOTE_BYTES=${AZURE_TDX_QUOTE_BYTES:-""}

# TDX collateral paths (optional)
export AZURE_TDX_PCS_CERT_PATH=${AZURE_TDX_PCS_CERT_PATH:-""}
export AZURE_TDX_QE_IDENTITY_PATH=${AZURE_TDX_QE_IDENTITY_PATH:-""}
export AZURE_TDX_TCB_INFO_PATH=${AZURE_TDX_TCB_INFO_PATH:-""}

# Foundry profile
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"layer1"}

# Broadcast transactions
export BROADCAST=${BROADCAST:-true}

# Verify smart contracts
export VERIFY=${VERIFY:-false}

# Required environment variable validation
if [ -z "$AZURE_TDX_VERIFIER_ADDRESS" ]; then
    echo "Error: AZURE_TDX_VERIFIER_ADDRESS not set"
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
echo "Azure TDX Verifier address: $AZURE_TDX_VERIFIER_ADDRESS"
echo "TDX PCS DAO address: $TDX_PCS_DAO_ADDRESS"
echo "TDX FMSPC TCB DAO address: $TDX_FMSPC_TCB_DAO_ADDRESS"
echo "TDX Enclave Identity DAO address: $TDX_ENCLAVE_IDENTITY_DAO_ADDRESS"
echo "TDX Enclave Identity Helper address: $TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS"
echo "New owner: $NEW_OWNER"

if [ -n "$AZURE_TDX_TRUSTED_PARAMS_BYTES" ]; then
    echo "Azure TDX Trusted Params provided"
fi

if [ -n "$AZURE_TDX_PCS_CERT_PATH" ]; then
    echo "TDX PCS Certificate path: $AZURE_TDX_PCS_CERT_PATH"
fi

if [ -n "$AZURE_TDX_QE_IDENTITY_PATH" ]; then
    echo "TDX QE Identity path: $AZURE_TDX_QE_IDENTITY_PATH"
fi

if [ -n "$AZURE_TDX_TCB_INFO_PATH" ]; then
    echo "TDX TCB Info path: $AZURE_TDX_TCB_INFO_PATH"
fi

if [ -n "$AZURE_TDX_QUOTE_BYTES" ]; then
    echo "Azure TDX Quote bytes provided for instance registration"
fi

# Run the setup script
forge script script/layer1/surge/SetupAzureTDXVerifier.s.sol \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    $VERIFY_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT

echo "Azure TDX verifier setup completed successfully!"