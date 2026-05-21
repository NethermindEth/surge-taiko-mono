#!/bin/bash
# setup_anvil_with_collaterals.sh
#
# This script:
# 1. Starts Anvil with preloaded Automata contracts
# 2. Fetches TDX collaterals from Intel's Trusted Services API
# 3. Adds the collaterals to the Automata DAOs
# 4. Deploys and initializes AzureTdxVerifier

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PROTOCOL_DIR="$REPO_ROOT/packages/protocol"
ASSETS_DIR="$REPO_ROOT/azure-tdx-assets"
TMP_DIR="$REPO_ROOT/.tmp"

# Create local temp directory
mkdir -p "$TMP_DIR"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Configuration
ANVIL_PORT=${ANVIL_PORT:-8545}
PRIVATE_KEY=${PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}
DEPLOYER_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
RPC_URL="http://localhost:$ANVIL_PORT"

# Intel Trusted Services API URLs
# FMSPC can be overridden - default is for Azure TDX
FMSPC=${FMSPC:-"90c06f000000"}
AZURE_TDX_TCB_LINK=${AZURE_TDX_TCB_LINK:-"https://api.trustedservices.intel.com/tdx/certification/v4/tcb?fmspc=$FMSPC"}
AZURE_TDX_QE_IDENTITY_LINK=${AZURE_TDX_QE_IDENTITY_LINK:-"https://api.trustedservices.intel.com/tdx/certification/v4/qe/identity"}

# Automata contract addresses (from network_params.yaml)
PCS_DAO="0x928826C6D0d1986bD0465697984fa3722ADE16E1"
FMSPC_TCB_DAO="0xA8C0F6F6Deb3dA48Be03A99C112737000a5a3088"
ENCLAVE_IDENTITY_DAO="0x5d1122a0d55b5095C0f03FBEa106A2e9722cb13F"
ENCLAVE_IDENTITY_HELPER="0x870D17e2C12aF1C47dD1f0e4aFd36e28c830D558"
AUTOMATA_DCAP_ATTESTATION="0x06303d57212EF0AA0d712694F3f4410EB7120f4E"

# AzureTdxVerifier deployment config
# These can be overridden via environment or flags
L2_CHAIN_ID=${L2_CHAIN_ID:-167004}
# Mock addresses - use anvil account #1 as default verifier (can be overridden)
TAIKO_INBOX=${TAIKO_INBOX:-"0x0000000000000000000000000000000000000001"}
PROOF_VERIFIER=${PROOF_VERIFIER:-"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}
VERIFIER_OWNER=${VERIFIER_OWNER:-"$DEPLOYER_ADDRESS"}

# Use static files as fallback
USE_STATIC_ASSETS=${USE_STATIC_ASSETS:-false}
STATIC_TCB_INFO_PATH="$PROTOCOL_DIR/test/layer1/automata-attestation/assets/0923/tdx_tcb_90c06f000000.json"
STATIC_QE_IDENTITY_PATH="$PROTOCOL_DIR/test/layer1/automata-attestation/assets/0923/tdx_identity.json"
STATIC_PCS_CERT_PATH="$PROTOCOL_DIR/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex"
STATIC_ROOT_PCS_CERT_PATH="$PROTOCOL_DIR/test/layer1/automata-attestation/assets/0923/tdx_root_pcs_cert.hex"

cleanup() {
    if [ -n "$ANVIL_PID" ] && [ "$KEEP_RUNNING" != "true" ]; then
        log_info "Stopping Anvil..."
        kill $ANVIL_PID 2>/dev/null || true
    fi
}

trap cleanup EXIT

# =============================================================================
# URL decode helper
# =============================================================================
url_decode() {
    local url_encoded="${1//+/ }"
    printf '%b' "${url_encoded//%/\\x}"
}

# =============================================================================
# Fetch TDX Collaterals from Intel API
# =============================================================================
fetch_collaterals() {
    log_info "Fetching TDX collaterals from Intel Trusted Services..."

    mkdir -p "$ASSETS_DIR"

    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║ Fetching TDX collaterals...                                  ║"
    echo "║ TCB Link: $AZURE_TDX_TCB_LINK"
    echo "║ QE Identity Link: $AZURE_TDX_QE_IDENTITY_LINK"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""

    # -------------------------------------------------------------------------
    # Fetch TCB Info
    # -------------------------------------------------------------------------
    log_info "Downloading TCB info..."
    TCB_RESPONSE=$(curl -s -D - -X GET "${AZURE_TDX_TCB_LINK}")

    # Save TCB JSON body
    echo "$TCB_RESPONSE" | sed '1,/^\r$/d' > "$ASSETS_DIR/tcb_full.json"

    # Minify TCB JSON
    if ! jq -c . "$ASSETS_DIR/tcb_full.json" > "$ASSETS_DIR/tcb.json" 2>/dev/null; then
        log_error "Failed to parse TCB JSON response"
        log_warn "Falling back to static assets..."
        USE_STATIC_ASSETS=true
        return
    fi
    log_info "TCB info saved to $ASSETS_DIR/tcb.json"

    # -------------------------------------------------------------------------
    # Extract TCB certificate chain from headers
    # -------------------------------------------------------------------------
    log_info "Extracting TCB certificate chain..."
    TCB_CERT_CHAIN=$(echo "$TCB_RESPONSE" | grep -i "^Tcb-Info-Issuer-Chain:" | cut -d' ' -f2- | tr -d '\r\n')

    if [ -z "$TCB_CERT_CHAIN" ]; then
        log_warn "Could not find Tcb-Info-Issuer-Chain header"
        log_warn "Falling back to static assets for certificates..."
        USE_STATIC_ASSETS=true
        return
    fi

    TCB_CERT_CHAIN_DECODED=$(url_decode "$TCB_CERT_CHAIN")

    # Extract certificates from chain
    echo "$TCB_CERT_CHAIN_DECODED" | awk '
        /-----BEGIN CERTIFICATE-----/ {
            cert_num++
            in_cert=1
        }
        in_cert {
            cert[cert_num] = cert[cert_num] $0 "\n"
        }
        /-----END CERTIFICATE-----/ {
            in_cert=0
        }
        END {
            for (i=1; i<=cert_num; i++) {
                print cert[i] > "'"$ASSETS_DIR"'/temp_tcb_cert_" i ".pem"
            }
            print "Found " cert_num " TCB certificates"
        }
    '

    # Process TCB Signing Certificate (first in chain)
    if [ -f "$ASSETS_DIR/temp_tcb_cert_1.pem" ]; then
        log_info "Processing TCB Signing Certificate..."
        mv "$ASSETS_DIR/temp_tcb_cert_1.pem" "$ASSETS_DIR/tdx_tcb_signing_cert.pem"
        openssl x509 -in "$ASSETS_DIR/tdx_tcb_signing_cert.pem" -outform DER -out "$ASSETS_DIR/tdx_tcb_signing_cert.der"
        echo -n "0x" > "$ASSETS_DIR/tdx_pcs_cert.hex"
        xxd -p -c 1000000 "$ASSETS_DIR/tdx_tcb_signing_cert.der" | tr -d '\n' >> "$ASSETS_DIR/tdx_pcs_cert.hex"
    fi

    # Process TCB Root CA Certificate (second in chain)
    if [ -f "$ASSETS_DIR/temp_tcb_cert_2.pem" ]; then
        log_info "Processing TCB Root CA Certificate..."
        mv "$ASSETS_DIR/temp_tcb_cert_2.pem" "$ASSETS_DIR/tdx_tcb_root_cert.pem"
        openssl x509 -in "$ASSETS_DIR/tdx_tcb_root_cert.pem" -outform DER -out "$ASSETS_DIR/tdx_tcb_root_cert.der"
        echo -n "0x" > "$ASSETS_DIR/tdx_root_pcs_cert.hex"
        xxd -p -c 1000000 "$ASSETS_DIR/tdx_tcb_root_cert.der" | tr -d '\n' >> "$ASSETS_DIR/tdx_root_pcs_cert.hex"
    fi

    # Cleanup temp files
    rm -f "$ASSETS_DIR/temp_tcb_cert_"*.pem

    # -------------------------------------------------------------------------
    # Fetch QE Identity
    # -------------------------------------------------------------------------
    log_info "Downloading QE identity..."
    if ! curl -s "${AZURE_TDX_QE_IDENTITY_LINK}" -o "$ASSETS_DIR/qe_identity_full.json"; then
        log_error "Failed to download QE identity"
        USE_STATIC_ASSETS=true
        return
    fi

    if ! jq -c . "$ASSETS_DIR/qe_identity_full.json" > "$ASSETS_DIR/qe_identity.json" 2>/dev/null; then
        log_error "Failed to parse QE identity JSON"
        USE_STATIC_ASSETS=true
        return
    fi
    log_info "QE identity saved to $ASSETS_DIR/qe_identity.json"

    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║ TDX collaterals fetched successfully!                        ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""

    # Copy fetched assets to protocol test directory (for foundry fs_permissions)
    PROTOCOL_ASSETS_DIR="$PROTOCOL_DIR/test/azure-tdx-assets"
    mkdir -p "$PROTOCOL_ASSETS_DIR"
    cp "$ASSETS_DIR/tcb.json" "$PROTOCOL_ASSETS_DIR/"
    cp "$ASSETS_DIR/qe_identity.json" "$PROTOCOL_ASSETS_DIR/"
    cp "$ASSETS_DIR/tdx_pcs_cert.hex" "$PROTOCOL_ASSETS_DIR/"
    cp "$ASSETS_DIR/tdx_root_pcs_cert.hex" "$PROTOCOL_ASSETS_DIR/"

    # Set paths to fetched assets (inside protocol/test for foundry access)
    TCB_INFO_PATH="$PROTOCOL_ASSETS_DIR/tcb.json"
    QE_IDENTITY_PATH="$PROTOCOL_ASSETS_DIR/qe_identity.json"
    PCS_CERT_PATH="$PROTOCOL_ASSETS_DIR/tdx_pcs_cert.hex"
    ROOT_PCS_CERT_PATH="$PROTOCOL_ASSETS_DIR/tdx_root_pcs_cert.hex"
}

# =============================================================================
# Use static assets (fallback)
# =============================================================================
use_static_assets() {
    log_info "Using static test assets..."
    TCB_INFO_PATH="$STATIC_TCB_INFO_PATH"
    QE_IDENTITY_PATH="$STATIC_QE_IDENTITY_PATH"
    PCS_CERT_PATH="$STATIC_PCS_CERT_PATH"
    ROOT_PCS_CERT_PATH="$STATIC_ROOT_PCS_CERT_PATH"

    # Verify files exist
    for f in "$TCB_INFO_PATH" "$QE_IDENTITY_PATH" "$PCS_CERT_PATH" "$ROOT_PCS_CERT_PATH"; do
        if [ ! -f "$f" ]; then
            log_error "Static asset not found: $f"
            exit 1
        fi
    done
}

# =============================================================================
# Generate initial_state.json
# =============================================================================
generate_state() {
    log_info "Generating initial_state.json..."
    cd "$REPO_ROOT"

    python3 << 'EOF'
import yaml
import json

with open('network_params.yaml', 'r') as f:
    data = yaml.safe_load(f)

preloaded = data.get('network_params', {}).get('additional_preloaded_contracts', {})
if not preloaded:
    raise Exception("No additional_preloaded_contracts found")

accounts = {}
for addr, contract in preloaded.items():
    accounts[addr.lower()] = {
        "nonce": contract.get("nonce", 1),
        "balance": contract.get("balance", "0x0"),
        "code": contract.get("code", "0x"),
        "storage": contract.get("storage", {})
    }

# Add funded accounts
for acc in ["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
            "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
            "0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc"]:
    if acc not in accounts:
        accounts[acc] = {"nonce": 0, "balance": "0x21e19e0c9bab2400000", "code": "0x", "storage": {}}

state = {
    "block": {
        "number": "0x0", "beneficiary": "0x0000000000000000000000000000000000000000",
        "timestamp": "0x0", "gas_limit": 30000000, "basefee": 1000000000, "difficulty": "0x0",
        "prevrandao": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "blob_excess_gas_and_price": {"excess_blob_gas": 0, "blob_gasprice": 1}
    },
    "accounts": accounts, "best_block_number": 0, "blocks": [], "transactions": []
}

with open('initial_state.json', 'w') as f:
    json.dump(state, f, indent=2)
print(f"Created initial_state.json with {len(accounts)} accounts")
EOF
}

# =============================================================================
# Start Anvil
# =============================================================================
start_anvil() {
    log_info "Starting Anvil..."
    cd "$REPO_ROOT"

    [ ! -f "initial_state.json" ] && generate_state

    anvil --load-state initial_state.json --port $ANVIL_PORT --block-time 1 --gas-limit 200000000 --silent &
    ANVIL_PID=$!

    # Wait for ready
    for i in {1..30}; do
        cast chain-id --rpc-url $RPC_URL 2>/dev/null && break
        sleep 1
    done

    # Verify
    if [ "$(cast code $PCS_DAO --rpc-url $RPC_URL 2>/dev/null)" = "0x" ]; then
        log_error "Automata contracts not loaded!"
        exit 1
    fi
    log_info "Anvil ready with Automata contracts"
}

# =============================================================================
# Add collaterals using cast and forge
# =============================================================================
add_collaterals() {
    log_info "Adding TDX collaterals..."

    cd "$PROTOCOL_DIR"

    # Verify collateral files exist
    for f in "$TCB_INFO_PATH" "$QE_IDENTITY_PATH" "$PCS_CERT_PATH" "$ROOT_PCS_CERT_PATH"; do
        if [ ! -f "$f" ]; then
            log_error "Collateral file not found: $f"
            exit 1
        fi
    done

    log_info "Using collaterals:"
    log_info "  TCB Info: $TCB_INFO_PATH"
    log_info "  QE Identity: $QE_IDENTITY_PATH"
    log_info "  PCS Cert: $PCS_CERT_PATH"
    log_info "  Root PCS Cert: $ROOT_PCS_CERT_PATH"

    # Read certificate hex files
    ROOT_CERT_HEX=$(cat "$ROOT_PCS_CERT_PATH")
    PCS_CERT_HEX=$(cat "$PCS_CERT_PATH")

    # Add Root PCS Certificate (CA.ROOT = 0)
    log_info "Adding Root PCS Certificate..."
    cast send $PCS_DAO "upsertPcsCertificates(uint8,bytes)(bytes32)" 0 "$ROOT_CERT_HEX" \
        --rpc-url $RPC_URL \
        --private-key $PRIVATE_KEY \
        --gas-limit 5000000 2>/dev/null || log_warn "Root cert may already exist"

    # Add TCB Signing Certificate (CA.SIGNING = 3)
    log_info "Adding TCB Signing Certificate..."
    cast send $PCS_DAO "upsertPcsCertificates(uint8,bytes)(bytes32)" 3 "$PCS_CERT_HEX" \
        --rpc-url $RPC_URL \
        --private-key $PRIVATE_KEY \
        --gas-limit 5000000 2>/dev/null || log_warn "Signing cert may already exist"

    # For TCB info and QE identity, we need to use forge script since they require JSON parsing
    log_info "Adding TCB Info and QE Identity via forge script..."

    # Create a minimal setup script
    cat > "$TMP_DIR/AddCollaterals.s.sol" << 'SOLIDITY'
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";
import "solady/src/utils/JSONParserLib.sol";
import "solady/src/utils/LibString.sol";

interface IFmspcTcbDao {
    function upsertFmspcTcb(TcbInfoJsonObj memory tcbInfoJson) external;
}

interface IEnclaveIdentityDao {
    function upsertEnclaveIdentity(uint256 id, uint256 isvsvn, EnclaveIdentityJsonObj memory identityJson) external;
}

interface IEnclaveIdentityHelper {
    function parseIdentityString(string memory identityStr) external pure returns (IdentityObj memory identity, string memory success);
}

struct TcbInfoJsonObj { string tcbInfoStr; bytes signature; }
struct EnclaveIdentityJsonObj { string identityStr; bytes signature; }
enum EnclaveId { QE, QVE, TD_QE }
struct IdentityObj {
    EnclaveId id; uint32 version; uint64 issueDateTimestamp; uint64 nextUpdateTimestamp;
    uint32 tcbEvaluationDataNumber; bytes4 miscselect; bytes4 miscselectMask;
    bytes16 attributes; bytes16 attributesMask; bytes32 mrsigner; uint16 isvprodid; Tcb[] tcb;
}
enum EnclaveIdTcbStatus { SGX_ENCLAVE_REPORT_ISVSVN_NOT_SUPPORTED, OK, SGX_ENCLAVE_REPORT_ISVSVN_REVOKED, SGX_ENCLAVE_REPORT_ISVSVN_OUT_OF_DATE }
struct Tcb { uint16 isvsvn; uint256 dateTimestamp; EnclaveIdTcbStatus status; }

contract AddCollaterals is Script {
    using JSONParserLib for JSONParserLib.Item;
    using LibString for string;

    function run() external {
        uint256 pk = vm.envUint("PRIVATE_KEY");
        address fmspcTcbDao = vm.envAddress("FMSPC_TCB_DAO");
        address enclaveIdentityDao = vm.envAddress("ENCLAVE_IDENTITY_DAO");
        address enclaveIdentityHelper = vm.envAddress("ENCLAVE_IDENTITY_HELPER");
        string memory tcbPath = vm.envString("TCB_INFO_PATH");
        string memory qePath = vm.envString("QE_IDENTITY_PATH");

        vm.startBroadcast(pk);

        // Add TCB Info
        string memory tcbJson = vm.readFile(tcbPath);
        TcbInfoJsonObj memory tcbObj = parseTcbInfoJson(tcbJson);
        try IFmspcTcbDao(fmspcTcbDao).upsertFmspcTcb(tcbObj) {
            console2.log("TCB Info added");
        } catch {
            console2.log("TCB Info may already exist");
        }

        // Add QE Identity
        string memory qeJson = vm.readFile(qePath);
        EnclaveIdentityJsonObj memory qeObj = parseEnclaveIdentityJson(qeJson);
        (IdentityObj memory identity,) = IEnclaveIdentityHelper(enclaveIdentityHelper).parseIdentityString(qeObj.identityStr);
        try IEnclaveIdentityDao(enclaveIdentityDao).upsertEnclaveIdentity(uint256(identity.id), 4, qeObj) {
            console2.log("QE Identity added");
        } catch {
            console2.log("QE Identity may already exist");
        }

        vm.stopBroadcast();
    }

    function parseTcbInfoJson(string memory jsonStr) internal pure returns (TcbInfoJsonObj memory result) {
        JSONParserLib.Item memory root = JSONParserLib.parse(jsonStr);
        JSONParserLib.Item[] memory children = root.children();
        for (uint256 i = 0; i < root.size(); i++) {
            string memory key = children[i].key();
            if (LibString.eq(key, "\"tcbInfo\"")) result.tcbInfoStr = children[i].value();
            else if (LibString.eq(key, "\"signature\"")) result.signature = vm.parseBytes(JSONParserLib.decodeString(children[i].value()));
        }
    }

    function parseEnclaveIdentityJson(string memory jsonStr) internal pure returns (EnclaveIdentityJsonObj memory result) {
        JSONParserLib.Item memory root = JSONParserLib.parse(jsonStr);
        JSONParserLib.Item[] memory children = root.children();
        for (uint256 i = 0; i < root.size(); i++) {
            string memory key = children[i].key();
            if (LibString.eq(key, "\"enclaveIdentity\"")) result.identityStr = children[i].value();
            else if (LibString.eq(key, "\"signature\"")) result.signature = vm.parseBytes(JSONParserLib.decodeString(children[i].value()));
        }
    }
}
SOLIDITY

    # Run the script
    PRIVATE_KEY=$PRIVATE_KEY \
    FMSPC_TCB_DAO=$FMSPC_TCB_DAO \
    ENCLAVE_IDENTITY_DAO=$ENCLAVE_IDENTITY_DAO \
    ENCLAVE_IDENTITY_HELPER=$ENCLAVE_IDENTITY_HELPER \
    TCB_INFO_PATH="$TCB_INFO_PATH" \
    QE_IDENTITY_PATH="$QE_IDENTITY_PATH" \
    FOUNDRY_PROFILE=layer1 \
    forge script "$TMP_DIR/AddCollaterals.s.sol:AddCollaterals" \
        --fork-url $RPC_URL \
        --broadcast \
        --ffi \
        -vvv \
        --block-gas-limit 200000000

    log_info "Collaterals added successfully!"
}

# =============================================================================
# Deploy AzureTdxVerifier
# =============================================================================
deploy_azure_tdx_verifier() {
    log_info "Deploying AzureTdxVerifier..."

    cd "$PROTOCOL_DIR"

    # Create deployment script
    cat > "$TMP_DIR/DeployAzureTdxVerifier.s.sol" << 'SOLIDITY'
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";
import "src/layer1/verifiers/AzureTdxVerifier.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployAzureTdxVerifier is Script {
    function run() external {
        uint256 pk = vm.envUint("PRIVATE_KEY");
        uint64 l2ChainId = uint64(vm.envUint("L2_CHAIN_ID"));
        address taikoInbox = vm.envAddress("TAIKO_INBOX");
        address proofVerifier = vm.envAddress("PROOF_VERIFIER");
        address automataDcapAttestation = vm.envAddress("AUTOMATA_DCAP_ATTESTATION");
        address owner = vm.envAddress("VERIFIER_OWNER");

        vm.startBroadcast(pk);

        // Deploy implementation
        AzureTdxVerifier impl = new AzureTdxVerifier(
            l2ChainId,
            taikoInbox,
            proofVerifier,
            automataDcapAttestation
        );
        console2.log("AzureTdxVerifier implementation deployed at:", address(impl));

        // Deploy proxy
        bytes memory initData = abi.encodeCall(AzureTdxVerifier.init, owner);
        ERC1967Proxy proxy = new ERC1967Proxy(address(impl), initData);
        console2.log("AzureTdxVerifier proxy deployed at:", address(proxy));

        // Verify initialization
        AzureTdxVerifier verifier = AzureTdxVerifier(address(proxy));
        console2.log("Verifier owner:", verifier.owner());
        console2.log("Taiko chain ID:", verifier.taikoChainId());
        console2.log("Taiko inbox:", verifier.taikoInbox());
        console2.log("Proof verifier:", verifier.taikoProofVerifier());
        console2.log("Automata DCAP:", verifier.automataDcapAttestation());

        vm.stopBroadcast();

        // Output the proxy address for the script to capture
        console2.log("AZURE_TDX_VERIFIER_ADDRESS:", address(proxy));
    }
}
SOLIDITY

    # Run deployment
    PRIVATE_KEY=$PRIVATE_KEY \
    L2_CHAIN_ID=$L2_CHAIN_ID \
    TAIKO_INBOX=$TAIKO_INBOX \
    PROOF_VERIFIER=$PROOF_VERIFIER \
    AUTOMATA_DCAP_ATTESTATION=$AUTOMATA_DCAP_ATTESTATION \
    VERIFIER_OWNER=$VERIFIER_OWNER \
    FOUNDRY_PROFILE=layer1 \
    forge script "$TMP_DIR/DeployAzureTdxVerifier.s.sol:DeployAzureTdxVerifier" \
        --fork-url $RPC_URL \
        --broadcast \
        -vvv \
        --block-gas-limit 200000000 2>&1 | tee "$TMP_DIR/deploy_output.log"

    # Extract the deployed address
    AZURE_TDX_VERIFIER_ADDRESS=$(grep "AZURE_TDX_VERIFIER_ADDRESS:" "$TMP_DIR/deploy_output.log" | tail -1 | awk '{print $2}')

    if [ -z "$AZURE_TDX_VERIFIER_ADDRESS" ]; then
        log_error "Failed to extract AzureTdxVerifier address from deployment output"
        return 1
    fi

    log_info "AzureTdxVerifier deployed at: $AZURE_TDX_VERIFIER_ADDRESS"

    # Save to file for reference
    echo "$AZURE_TDX_VERIFIER_ADDRESS" > "$REPO_ROOT/azure-tdx-assets/verifier_address.txt"
}

# =============================================================================
# Main
# =============================================================================
main() {
    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║  Anvil with Automata + AzureTdxVerifier                      ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""

    KEEP_RUNNING=false
    SKIP_COLLATERALS=false
    SKIP_DEPLOY=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --keep-running) KEEP_RUNNING=true; shift ;;
            --skip-collaterals) SKIP_COLLATERALS=true; shift ;;
            --skip-deploy) SKIP_DEPLOY=true; shift ;;
            --use-static) USE_STATIC_ASSETS=true; shift ;;
            --fmspc) FMSPC="$2"; AZURE_TDX_TCB_LINK="https://api.trustedservices.intel.com/tdx/certification/v4/tcb?fmspc=$FMSPC"; shift 2 ;;
            --taiko-inbox) TAIKO_INBOX="$2"; shift 2 ;;
            --proof-verifier) PROOF_VERIFIER="$2"; shift 2 ;;
            --l2-chain-id) L2_CHAIN_ID="$2"; shift 2 ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --keep-running       Keep Anvil running after setup"
                echo "  --skip-collaterals   Only start Anvil, skip adding collaterals"
                echo "  --skip-deploy        Skip AzureTdxVerifier deployment"
                echo "  --use-static         Use static test assets instead of fetching from Intel API"
                echo "  --fmspc <FMSPC>      Specify FMSPC for TCB info (default: 90c06f000000)"
                echo "  --taiko-inbox <ADDR> TaikoInbox address (default: mock 0x...0001)"
                echo "  --proof-verifier <ADDR> ProofVerifier address (default: anvil account #1)"
                echo "  --l2-chain-id <ID>   L2 chain ID (default: 167004)"
                echo ""
                echo "Environment variables:"
                echo "  ANVIL_PORT                  Anvil port (default: 8545)"
                echo "  AZURE_TDX_TCB_LINK          Custom TCB info URL"
                echo "  AZURE_TDX_QE_IDENTITY_LINK  Custom QE identity URL"
                echo "  TAIKO_INBOX                 TaikoInbox contract address"
                echo "  PROOF_VERIFIER              ProofVerifier contract address"
                echo "  L2_CHAIN_ID                 L2 chain ID"
                exit 0 ;;
            *) log_error "Unknown option: $1"; exit 1 ;;
        esac
    done

    # Fetch or use static collaterals
    if [ "$SKIP_COLLATERALS" = false ]; then
        if [ "$USE_STATIC_ASSETS" = true ]; then
            use_static_assets
        else
            fetch_collaterals
            # If fetch failed, use static
            if [ "$USE_STATIC_ASSETS" = true ]; then
                use_static_assets
            fi
        fi
    fi

    start_anvil

    if [ "$SKIP_COLLATERALS" = false ]; then
        add_collaterals
    fi

    # Deploy AzureTdxVerifier
    if [ "$SKIP_DEPLOY" = false ]; then
        deploy_azure_tdx_verifier
    fi

    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║  Setup Complete!                                             ║"
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║  RPC URL: http://localhost:$ANVIL_PORT                               ║"
    echo "║  FMSPC: $FMSPC                                         ║"
    echo "║  L2 Chain ID: $L2_CHAIN_ID                                       ║"
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║  Automata Contracts:                                         ║"
    echo "║    PcsDao:             $PCS_DAO  ║"
    echo "║    FmspcTcbDao:        $FMSPC_TCB_DAO  ║"
    echo "║    EnclaveIdentityDao: $ENCLAVE_IDENTITY_DAO  ║"
    if [ -n "$AZURE_TDX_VERIFIER_ADDRESS" ]; then
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║  AzureTdxVerifier:     $AZURE_TDX_VERIFIER_ADDRESS  ║"
    echo "║    TaikoInbox:         $TAIKO_INBOX  ║"
    echo "║    ProofVerifier:      $PROOF_VERIFIER  ║"
    fi
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""

    if [ "$KEEP_RUNNING" = true ]; then
        log_info "Anvil running. Press Ctrl+C to stop."
        trap - EXIT
        wait $ANVIL_PID
    else
        log_info "Anvil PID: $ANVIL_PID"
        disown $ANVIL_PID 2>/dev/null || true
        trap - EXIT
    fi
}

main "$@"
