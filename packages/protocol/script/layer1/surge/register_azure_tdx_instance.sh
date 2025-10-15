#!/bin/sh

# This script registers TDX instances for the Surge protocol
set -e

# Deployer private key
export PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}

# Network configuration
export FORK_URL=${FORK_URL:-"http://localhost:8545"}

# TDX configuration
export TDX_PCS_DAO_ADDRESS=${TDX_PCS_DAO_ADDRESS:-"0x0000000000000000000000000000000000000000"}
export TDX_FMSPC_TCB_DAO_ADDRESS=${TDX_FMSPC_TCB_DAO_ADDRESS:-"0x0000000000000000000000000000000000000000"}
export TDX_ENCLAVE_IDENTITY_DAO_ADDRESS=${TDX_ENCLAVE_IDENTITY_DAO_ADDRESS:-"0x0000000000000000000000000000000000000000"}
export TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS=${TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS:-"0x0000000000000000000000000000000000000000"}
### export AUTOMATA_DCAP_ATTESTATION=${AUTOMATA_DCAP_ATTESTATION:-"0x0000000000000000000000000000000000000000"}
export AZURE_TDX_VERIFIER=${AZURE_TDX_VERIFIER:-"0x0000000000000000000000000000000000000000"}
### export PEM_CERT_CHAIN_LIB=${PEM_CERT_CHAIN_LIB:-"0x0000000000000000000000000000000000000000"}

# Trusted parameters configuration
export AZURE_TDX_TRUSTED_PARAMS_BYTES=${AZURE_TDX_TRUSTED_PARAMS_BYTES:-"0x"}

# Attestation configuration
export AZURE_TDX_QE_IDENTITY_PATH=${AZURE_TDX_QE_IDENTITY_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_identity.json"}
export AZURE_TDX_TCB_INFO_PATH=${AZURE_TDX_TCB_INFO_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_tcb_90c06f000000.json"}
export AZURE_TDX_PCS_CERT_PATH=${AZURE_TDX_PCS_CERT_PATH:-"/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex"}
export AZURE_TDX_QUOTE_BYTES=${AZURE_TDX_QUOTE_BYTES:-"0x"}

# Foundry configuration
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

forge script ./script/layer1/surge/RegisterAzureTDXInstance.s.sol:RegisterAzureTDXInstance \
    --fork-url $FORK_URL \
    $BROADCAST_ARG \
    --ffi \
    $LOG_LEVEL \
    --private-key $PRIVATE_KEY \
    --block-gas-limit $BLOCK_GAS_LIMIT 