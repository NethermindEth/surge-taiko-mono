#!/bin/bash
# run_anvil_tdx_testnet.sh
#
# This script:
# 1. Starts Anvil with preloaded Automata contracts
# 2. Deploys AzureTdxVerifier
# 3. Sets up TDX collaterals (TCB info, QE identity, PCS certificates)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PROTOCOL_DIR="$REPO_ROOT/packages/protocol"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Configuration
ANVIL_PORT=${ANVIL_PORT:-8545}
ANVIL_CHAIN_ID=${ANVIL_CHAIN_ID:-31337}
PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}
DEPLOYER_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

# Automata contract addresses (from network_params.yaml)
export TDX_PCS_DAO_ADDRESS="0x928826C6D0d1986bD0465697984fa3722ADE16E1"
export TDX_FMSPC_TCB_DAO_ADDRESS="0xA8C0F6F6Deb3dA48Be03A99C112737000a5a3088"
export TDX_ENCLAVE_IDENTITY_DAO_ADDRESS="0x5d1122a0d55b5095C0f03FBEa106A2e9722cb13F"
export TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS="0x870D17e2C12aF1C47dD1f0e4aFd36e28c830D558"
export TDX_AUTOMATA_DCAP_ATTESTATION_ADDRESS="0x06303d57212EF0AA0d712694F3f4410EB7120f4E"

# Collateral paths (relative to protocol package)
TDX_ASSETS_DIR="$PROTOCOL_DIR/test/layer1/automata-attestation/assets/0923"

cleanup() {
    log_info "Cleaning up..."
    if [ -n "$ANVIL_PID" ]; then
        kill $ANVIL_PID 2>/dev/null || true
    fi
}

trap cleanup EXIT

# =============================================================================
# Step 1: Generate initial_state.json if not exists
# =============================================================================
generate_initial_state() {
    log_info "Generating initial_state.json..."

    cd "$REPO_ROOT"

    python3 << 'PYTHON_SCRIPT'
import yaml
import json
import sys

with open('network_params.yaml', 'r') as f:
    data = yaml.safe_load(f)

preloaded = data.get('network_params', {}).get('additional_preloaded_contracts', {})

if not preloaded:
    print("ERROR: No additional_preloaded_contracts found", file=sys.stderr)
    sys.exit(1)

accounts = {}
for addr, contract in preloaded.items():
    accounts[addr.lower()] = {
        "nonce": contract.get("nonce", 1),
        "balance": contract.get("balance", "0x0"),
        "code": contract.get("code", "0x"),
        "storage": contract.get("storage", {})
    }

# Add funded accounts
funded_accounts = [
    "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
    "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
    "0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc",
    "0x90f79bf6eb2c4f870365e785982e1f101e93b906",
    "0x15d34aaf54267db7d7c367839aaf71a00a2c6a65",
    "0x9965507d1a55bcc2695c58ba16fb37d819b0a4dc",
    "0x976ea74026e726554db657fa54763abd0c3a0aa9",
    "0x14dc79964da2c08b23698b3d3cc7ca32193d9955",
    "0x23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f",
    "0xa0ee7a142d267c1f36714e4a8f75612f20a79720",
]
for acc in funded_accounts:
    if acc not in accounts:
        accounts[acc] = {
            "nonce": 0,
            "balance": "0x21e19e0c9bab2400000",
            "code": "0x",
            "storage": {}
        }

state = {
    "block": {
        "number": "0x0",
        "beneficiary": "0x0000000000000000000000000000000000000000",
        "timestamp": "0x0",
        "gas_limit": 30000000,
        "basefee": 1000000000,
        "difficulty": "0x0",
        "prevrandao": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "blob_excess_gas_and_price": {"excess_blob_gas": 0, "blob_gasprice": 1}
    },
    "accounts": accounts,
    "best_block_number": 0,
    "blocks": [],
    "transactions": []
}

with open('initial_state.json', 'w') as f:
    json.dump(state, f, indent=2)

print(f"Created initial_state.json with {len(accounts)} accounts")
PYTHON_SCRIPT
}

# =============================================================================
# Step 2: Start Anvil
# =============================================================================
start_anvil() {
    log_info "Starting Anvil on port $ANVIL_PORT..."

    cd "$REPO_ROOT"

    # Check if initial_state.json exists
    if [ ! -f "initial_state.json" ]; then
        generate_initial_state
    fi

    anvil \
        --load-state initial_state.json \
        --port $ANVIL_PORT \
        --chain-id $ANVIL_CHAIN_ID \
        --block-time 1 \
        --gas-limit 200000000 \
        --silent &

    ANVIL_PID=$!

    # Wait for anvil to be ready
    log_info "Waiting for Anvil to be ready..."
    for i in {1..30}; do
        if cast chain-id --rpc-url http://localhost:$ANVIL_PORT 2>/dev/null; then
            log_info "Anvil is ready!"
            break
        fi
        sleep 1
    done

    # Verify contracts are loaded
    log_info "Verifying preloaded contracts..."
    PCS_CODE=$(cast code $TDX_PCS_DAO_ADDRESS --rpc-url http://localhost:$ANVIL_PORT 2>/dev/null || echo "0x")
    if [ "$PCS_CODE" = "0x" ]; then
        log_error "PcsDao contract not loaded!"
        exit 1
    fi
    log_info "Automata contracts loaded successfully"
}

# =============================================================================
# Step 3: Deploy AzureTdxVerifier
# =============================================================================
deploy_azure_tdx_verifier() {
    log_info "Deploying AzureTdxVerifier..."

    cd "$PROTOCOL_DIR"

    # Set environment for deployment
    export PRIVATE_KEY
    export FORK_URL="http://localhost:$ANVIL_PORT"
    export FOUNDRY_PROFILE="layer1"
    export BROADCAST="true"
    export VERIFY="false"
    export LOG_LEVEL="-vvv"
    export BLOCK_GAS_LIMIT="200000000"

    # L2 configuration (dummy values for local testing)
    export L2_NETWORK="devnet"
    export L2_CHAINID="167004"
    export L2_GENESIS_HASH="0xee1950562d42f0da28bd4550d88886bc90894c77c9c9eaefef775d4c8223f259"

    # Liveness configuration
    export MAX_VERIFICATION_DELAY="100"
    export MIN_VERIFICATION_STREAK="10"
    export LIVENESS_BOND_BASE="1000000000000000000"
    export COOLDOWN_WINDOW="86400"

    # Owner configuration
    export USE_TIMELOCKED_OWNER="false"
    export OWNER_MULTISIG="$DEPLOYER_ADDRESS"
    export OWNER_MULTISIG_SIGNERS="$DEPLOYER_ADDRESS"
    export TIMELOCK_PERIOD="0"
    export DAO="$DEPLOYER_ADDRESS"

    # Preconf configuration
    export USE_PRECONF="false"

    # Forced inclusion configuration
    export INCLUSION_WINDOW="24"
    export INCLUSION_FEE_IN_GWEI="100"

    # Only deploy Azure TDX verifier
    export DEPLOY_SGX_RETH_VERIFIER="false"
    export DEPLOY_SGX_GETH_VERIFIER="false"
    export DEPLOY_TDX_VERIFIER="false"
    export DEPLOY_AZURE_TDX_VERIFIER="true"
    export DEPLOY_RISC0_RETH_VERIFIER="false"
    export DEPLOY_SP1_RETH_VERIFIER="false"

    # Run deployment
    forge script ./script/layer1/surge/DeploySurgeL1.s.sol:DeploySurgeL1 \
        --fork-url $FORK_URL \
        --broadcast \
        --ffi \
        $LOG_LEVEL \
        --private-key $PRIVATE_KEY \
        --block-gas-limit $BLOCK_GAS_LIMIT

    # Get the deployed verifier address from the deployment output
    if [ -f "deployments/azure_tdx_nethermind_verifier.json" ]; then
        AZURE_TDX_VERIFIER_ADDRESS=$(jq -r '.address' deployments/azure_tdx_nethermind_verifier.json)
        log_info "AzureTdxVerifier deployed at: $AZURE_TDX_VERIFIER_ADDRESS"
    else
        # Try to find it in the broadcast logs
        AZURE_TDX_VERIFIER_ADDRESS=$(grep -r "azure_tdx_nethermind_verifier" broadcast/ 2>/dev/null | grep -oE '0x[a-fA-F0-9]{40}' | head -1 || echo "")
        if [ -z "$AZURE_TDX_VERIFIER_ADDRESS" ]; then
            log_warn "Could not automatically detect AzureTdxVerifier address"
            log_warn "Check deployments/ directory for the address"
        else
            log_info "AzureTdxVerifier deployed at: $AZURE_TDX_VERIFIER_ADDRESS"
        fi
    fi

    export AZURE_TDX_VERIFIER_ADDRESS
}

# =============================================================================
# Step 4: Setup TDX Collaterals
# =============================================================================
setup_tdx_collaterals() {
    log_info "Setting up TDX collaterals..."

    cd "$PROTOCOL_DIR"

    # Check if we have the verifier address
    if [ -z "$AZURE_TDX_VERIFIER_ADDRESS" ]; then
        log_error "AZURE_TDX_VERIFIER_ADDRESS not set. Please set it manually."
        log_info "You can find it in the deployment output or deployments/ directory"
        read -p "Enter AzureTdxVerifier address: " AZURE_TDX_VERIFIER_ADDRESS
    fi

    # Set environment
    export PRIVATE_KEY
    export FORK_URL="http://localhost:$ANVIL_PORT"
    export FOUNDRY_PROFILE="layer1"
    export BROADCAST="true"
    export LOG_LEVEL="-vvv"
    export BLOCK_GAS_LIMIT="200000000"

    # The new owner after setup (deployer for local testing)
    export NEW_OWNER="$DEPLOYER_ADDRESS"

    # Collateral paths (relative to packages/protocol)
    export AZURE_TDX_TCB_INFO_PATH="/test/layer1/automata-attestation/assets/0923/tdx_tcb_90c06f000000.json"
    export AZURE_TDX_QE_IDENTITY_PATH="/test/layer1/automata-attestation/assets/0923/tdx_identity.json"
    export AZURE_TDX_PCS_CERT_PATH="/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex"
    export AZURE_TDX_ROOT_PCS_CERT_PATH="/test/layer1/automata-attestation/assets/0923/tdx_root_pcs_cert.hex"

    log_info "Running setup script with collaterals..."
    log_info "  TCB Info: $AZURE_TDX_TCB_INFO_PATH"
    log_info "  QE Identity: $AZURE_TDX_QE_IDENTITY_PATH"
    log_info "  PCS Cert: $AZURE_TDX_PCS_CERT_PATH"
    log_info "  Root PCS Cert: $AZURE_TDX_ROOT_PCS_CERT_PATH"

    # Run the setup script
    forge script script/layer1/surge/SetupAzureTDXVerifier.s.sol \
        --fork-url $FORK_URL \
        --broadcast \
        --ffi \
        $LOG_LEVEL \
        --private-key $PRIVATE_KEY \
        --block-gas-limit $BLOCK_GAS_LIMIT

    log_info "TDX collaterals setup complete!"
}

# =============================================================================
# Main
# =============================================================================
main() {
    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║     Anvil TDX Testnet Setup                                  ║"
    echo "║     - Preloaded Automata contracts                           ║"
    echo "║     - AzureTdxVerifier deployment                            ║"
    echo "║     - TDX collaterals configuration                          ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""

    # Parse arguments
    SKIP_DEPLOY=false
    SKIP_SETUP=false
    KEEP_RUNNING=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-deploy)
                SKIP_DEPLOY=true
                shift
                ;;
            --skip-setup)
                SKIP_SETUP=true
                shift
                ;;
            --keep-running)
                KEEP_RUNNING=true
                shift
                ;;
            --verifier-address)
                AZURE_TDX_VERIFIER_ADDRESS="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --skip-deploy        Skip AzureTdxVerifier deployment"
                echo "  --skip-setup         Skip TDX collaterals setup"
                echo "  --keep-running       Keep Anvil running after setup"
                echo "  --verifier-address   Specify existing AzureTdxVerifier address"
                echo "  --help               Show this help message"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done

    # Step 1: Start Anvil
    start_anvil

    # Step 2: Deploy AzureTdxVerifier
    if [ "$SKIP_DEPLOY" = false ]; then
        deploy_azure_tdx_verifier
    else
        log_info "Skipping AzureTdxVerifier deployment"
        if [ -z "$AZURE_TDX_VERIFIER_ADDRESS" ]; then
            log_error "Please provide --verifier-address when using --skip-deploy"
            exit 1
        fi
    fi

    # Step 3: Setup TDX Collaterals
    if [ "$SKIP_SETUP" = false ]; then
        setup_tdx_collaterals
    else
        log_info "Skipping TDX collaterals setup"
    fi

    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║     Setup Complete!                                          ║"
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║  Anvil RPC:          http://localhost:$ANVIL_PORT                    ║"
    echo "║  Chain ID:           $ANVIL_CHAIN_ID                                 ║"
    if [ -n "$AZURE_TDX_VERIFIER_ADDRESS" ]; then
    echo "║  AzureTdxVerifier:   $AZURE_TDX_VERIFIER_ADDRESS  ║"
    fi
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║  Automata Contracts:                                         ║"
    echo "║    PcsDao:           $TDX_PCS_DAO_ADDRESS  ║"
    echo "║    FmspcTcbDao:      $TDX_FMSPC_TCB_DAO_ADDRESS  ║"
    echo "║    EnclaveIdentityDao: $TDX_ENCLAVE_IDENTITY_DAO_ADDRESS  ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""

    if [ "$KEEP_RUNNING" = true ]; then
        log_info "Anvil is running. Press Ctrl+C to stop."
        wait $ANVIL_PID
    else
        log_info "Setup complete. Anvil PID: $ANVIL_PID"
        log_info "Run 'kill $ANVIL_PID' to stop Anvil"
        # Disown the process so it keeps running
        disown $ANVIL_PID 2>/dev/null || true
        trap - EXIT
    fi
}

main "$@"
