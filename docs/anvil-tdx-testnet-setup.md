# Running Anvil Local Testnet with Automata and AzureTdxVerifier

This guide explains how to set up a local Anvil testnet with the Automata DCAP attestation infrastructure and AzureTdxVerifier contracts preloaded, including the necessary TDX collaterals.

## Overview

The setup involves:

1. Creating an `initial_state.json` with preloaded Automata contracts
2. Fetching TDX collaterals (TCB info, QE identity, PCS certificates)
3. Starting Anvil with the preloaded state
4. Deploying AzureTdxVerifier
5. Setting up the verifier with collaterals

## Prerequisites

- [Foundry](https://book.getfoundry.sh/getting-started/installation) installed
- `curl`, `jq`, `openssl`, `xxd` available
- Access to Intel's Trusted Services API

## Step 1: Create Initial State JSON

The Automata contracts need to be preloaded into Anvil. Extract the `additional_preloaded_contracts` from `network_params.yaml` and convert to Anvil's state format.

### Automata Contract Addresses

| Contract                | Address                                      |
| ----------------------- | -------------------------------------------- |
| AutomataDcapAttestation | `0x06303d57212EF0AA0d712694F3f4410EB7120f4E` |
| PcsDao                  | `0x928826C6D0d1986bD0465697984fa3722ADE16E1` |
| FmspcTcbDao             | `0xA8C0F6F6Deb3dA48Be03A99C112737000a5a3088` |
| EnclaveIdentityDao      | `0x5d1122a0d55b5095C0f03FBEa106A2e9722cb13F` |
| EnclaveIdentityHelper   | `0x870D17e2C12aF1C47dD1f0e4aFd36e28c830D558` |
| P256Verifier            | `0x2ff6f69c51b27799df6b797bfd2dbfb6f35aa495` |
| X509Helper              | `0xde437f4bf3b738f98260ae1bb882c38c92144e51` |
| FmspcTcbHelper          | `0xcce17f2fe091a2124bf541d4da0453a68a171f0a` |
| DcapV4Router            | `0x4c5f86b5251a5a62336747379575778881467640` |
| SigVerify               | `0x9b85d55246461989bb3ab491cf6047bbf33c2eb1` |
| CREATE2 Deployer        | `0x4e59b44847b379578588920ca78fbf26c0b4956c` |

### Generate initial_state.json

Create a script to extract contracts from `network_params.yaml`:

```bash
#!/bin/bash
# generate_initial_state.sh

# Extract additional_preloaded_contracts from network_params.yaml and convert to anvil state format
python3 << 'EOF'
import yaml
import json

with open('network_params.yaml', 'r') as f:
    content = f.read()

# Parse YAML
data = yaml.safe_load(content)

# Get the preloaded contracts (under network_params)
preloaded = data.get('network_params', {}).get('additional_preloaded_contracts', {})

# Convert to anvil state format
accounts = {}
for addr, contract in preloaded.items():
    accounts[addr.lower()] = {
        "nonce": contract.get("nonce", 1),
        "balance": contract.get("balance", "0x0"),
        "code": contract.get("code", "0x"),
        "storage": contract.get("storage", {})
    }

# Add default funded accounts (anvil default accounts)
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
            "balance": "0x21e19e0c9bab2400000",  # 10000 ETH
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
        "blob_excess_gas_and_price": {
            "excess_blob_gas": 0,
            "blob_gasprice": 1
        }
    },
    "accounts": accounts,
    "best_block_number": 0,
    "blocks": [],
    "transactions": []
}

with open('initial_state.json', 'w') as f:
    json.dump(state, f, indent=2)

print(f"Created initial_state.json with {len(accounts)} accounts")
EOF
```

Alternatively, use the pre-built state from `network_params.yaml` using yq/jq:

```bash
# Quick extraction using yq and jq
yq -o=json '.network_params.additional_preloaded_contracts' network_params.yaml > contracts.json

# Create the full state file
jq '{
  block: {
    number: "0x0",
    beneficiary: "0x0000000000000000000000000000000000000000",
    timestamp: "0x0",
    gas_limit: 30000000,
    basefee: 1000000000,
    difficulty: "0x0",
    prevrandao: "0x0000000000000000000000000000000000000000000000000000000000000000",
    blob_excess_gas_and_price: { excess_blob_gas: 0, blob_gasprice: 1 }
  },
  accounts: (to_entries | map({key: .key | ascii_downcase, value: .value}) | from_entries),
  best_block_number: 0,
  blocks: [],
  transactions: []
}' contracts.json > initial_state.json
```

A pre-generated `initial_state.json` is available in the repository root.

## Step 2: Fetch TDX Collaterals

TDX attestation requires collaterals from Intel's Trusted Services API. The `scripts/setup_anvil_with_collaterals.sh` script handles this automatically, but you can also fetch them manually:

### Automatic Fetching (Recommended)

The setup script will automatically fetch collaterals from Intel's API:

```bash
# Fetches collaterals for default FMSPC (90c06f000000)
./scripts/setup_anvil_with_collaterals.sh --keep-running

# Specify a different FMSPC
./scripts/setup_anvil_with_collaterals.sh --fmspc 00606a000000 --keep-running

# Use static test assets instead of fetching
./scripts/setup_anvil_with_collaterals.sh --use-static --keep-running
```

### Manual Fetching

Use the following script to fetch collaterals manually:

```bash
#!/bin/bash
# fetch_tdx_collaterals.sh

set -e

# Configuration - update these URLs based on your TDX platform
# For Azure TDX, use the Azure attestation service URLs
AZURE_TDX_TCB_LINK="${AZURE_TDX_TCB_LINK:-https://api.trustedservices.intel.com/tdx/certification/v4/tcb?fmspc=90c06f000000}"
AZURE_TDX_QE_IDENTITY_LINK="${AZURE_TDX_QE_IDENTITY_LINK:-https://api.trustedservices.intel.com/tdx/certification/v4/qe/identity}"

OUTPUT_DIR="${OUTPUT_DIR:-./azure-tdx-assets}"

echo "Creating output directory: $OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

url_decode() {
    local url_encoded="${1//+/ }"
    printf '%b' "${url_encoded//%/\\x}"
}

echo "=========================================="
echo "Downloading TCB info from $AZURE_TDX_TCB_LINK"
echo "=========================================="

TCB_RESPONSE=$(curl -s -D - -X GET "${AZURE_TDX_TCB_LINK}")

# Save full TCB response
echo "$TCB_RESPONSE" | sed '1,/^\r$/d' > "$OUTPUT_DIR/tcb_full.json"

# Minify TCB JSON
jq -c . "$OUTPUT_DIR/tcb_full.json" > "$OUTPUT_DIR/tcb.json"
echo "TCB info saved to $OUTPUT_DIR/tcb.json"

echo "=========================================="
echo "Extracting TCB cert chain from headers"
echo "=========================================="

TCB_CERT_CHAIN=$(echo "$TCB_RESPONSE" | grep -i "^Tcb-Info-Issuer-Chain:" | cut -d' ' -f2- | tr -d '\r\n')

if [ -z "$TCB_CERT_CHAIN" ]; then
    echo "Warning: Could not find Tcb-Info-Issuer-Chain header"
    echo "Using static certificates from test assets..."
else
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
                print cert[i] > "'"$OUTPUT_DIR"'/temp_tcb_cert_" i ".pem"
            }
            print "Found " cert_num " TCB certificates"
        }
    '

    # Process signing certificate (first in chain)
    if [ -f "$OUTPUT_DIR/temp_tcb_cert_1.pem" ]; then
        echo "Processing TCB Signing Certificate..."
        mv "$OUTPUT_DIR/temp_tcb_cert_1.pem" "$OUTPUT_DIR/tdx_tcb_signing_cert.pem"
        openssl x509 -in "$OUTPUT_DIR/tdx_tcb_signing_cert.pem" -outform DER -out "$OUTPUT_DIR/tdx_tcb_signing_cert.der"
        echo -n "0x" > "$OUTPUT_DIR/tdx_tcb_signing_cert.hex"
        xxd -p -c 1000000 "$OUTPUT_DIR/tdx_tcb_signing_cert.der" | tr -d '\n' >> "$OUTPUT_DIR/tdx_tcb_signing_cert.hex"
    fi

    # Process root certificate (second in chain)
    if [ -f "$OUTPUT_DIR/temp_tcb_cert_2.pem" ]; then
        echo "Processing TCB Root CA Certificate..."
        mv "$OUTPUT_DIR/temp_tcb_cert_2.pem" "$OUTPUT_DIR/tdx_tcb_root_cert.pem"
        openssl x509 -in "$OUTPUT_DIR/tdx_tcb_root_cert.pem" -outform DER -out "$OUTPUT_DIR/tdx_tcb_root_cert.der"
        echo -n "0x" > "$OUTPUT_DIR/tdx_tcb_root_cert.hex"
        xxd -p -c 1000000 "$OUTPUT_DIR/tdx_tcb_root_cert.der" | tr -d '\n' >> "$OUTPUT_DIR/tdx_tcb_root_cert.hex"
    fi

    # Cleanup temp files
    rm -f "$OUTPUT_DIR/temp_tcb_cert_"*.pem
fi

echo "=========================================="
echo "Downloading QE identity"
echo "=========================================="

curl -s "${AZURE_TDX_QE_IDENTITY_LINK}" -o "$OUTPUT_DIR/qe_identity_full.json"
jq -c . "$OUTPUT_DIR/qe_identity_full.json" > "$OUTPUT_DIR/qe_identity.json"
echo "QE identity saved to $OUTPUT_DIR/qe_identity.json"

echo "=========================================="
echo "TDX collaterals fetched successfully!"
echo "=========================================="
echo "Files created:"
ls -la "$OUTPUT_DIR"
```

### Using Static Test Assets

For local testing, you can use the pre-existing test assets:

```bash
# Copy test assets
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_tcb_90c06f000000.json ./azure-tdx-assets/tcb.json
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_identity.json ./azure-tdx-assets/qe_identity.json
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex ./azure-tdx-assets/tdx_pcs_cert.hex
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_root_pcs_cert.hex ./azure-tdx-assets/tdx_root_pcs_cert.hex
```

## Step 3: Start Anvil with Preloaded State

```bash
# Start anvil with the initial state
anvil --load-state initial_state.json \
      --block-time 1 \
      --chain-id 31337 \
      --gas-limit 200000000
```

Verify the contracts are loaded:

```bash
# Check if Automata contracts have code
cast code 0x928826C6D0d1986bD0465697984fa3722ADE16E1  # PcsDao
cast code 0xA8C0F6F6Deb3dA48Be03A99C112737000a5a3088  # FmspcTcbDao
cast code 0x5d1122a0d55b5095C0f03FBEa106A2e9722cb13F  # EnclaveIdentityDao
```

## Step 4: Deploy AzureTdxVerifier

Navigate to the protocol package and deploy:

```bash
cd packages/protocol

# Set environment variables
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
export FORK_URL=http://localhost:8545
export FOUNDRY_PROFILE=layer1

# Deploy the verifier (as part of Surge L1 deployment or standalone)
export DEPLOY_AZURE_TDX_VERIFIER=true
export DEPLOY_SGX_RETH_VERIFIER=false
export DEPLOY_SGX_GETH_VERIFIER=false
export DEPLOY_TDX_VERIFIER=false
export DEPLOY_RISC0_RETH_VERIFIER=false
export DEPLOY_SP1_RETH_VERIFIER=false

# Set Automata contract addresses
export TDX_PCS_DAO_ADDRESS=0x928826C6D0d1986bD0465697984fa3722ADE16E1
export TDX_FMSPC_TCB_DAO_ADDRESS=0xA8C0F6F6Deb3dA48Be03A99C112737000a5a3088
export TDX_ENCLAVE_IDENTITY_DAO_ADDRESS=0x5d1122a0d55b5095C0f03FBEa106A2e9722cb13F
export TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS=0x870D17e2C12aF1C47dD1f0e4aFd36e28c830D558
export TDX_AUTOMATA_DCAP_ATTESTATION_ADDRESS=0x06303d57212EF0AA0d712694F3f4410EB7120f4E

# Run deployment
./script/layer1/surge/deploy_surge_l1.sh
```

Or deploy just the AzureTdxVerifier:

```bash
# Get the deployed verifier address from the deployment output
# Then setup the verifier with collaterals
```

## Step 5: Setup AzureTdxVerifier with Collaterals

```bash
cd packages/protocol

# Set environment variables
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
export FORK_URL=http://localhost:8545
export FOUNDRY_PROFILE=layer1
export BROADCAST=true

# Set verifier address (from deployment output)
export AZURE_TDX_VERIFIER_ADDRESS=<deployed-verifier-address>

# Set Automata contract addresses
export TDX_PCS_DAO_ADDRESS=0x928826C6D0d1986bD0465697984fa3722ADE16E1
export TDX_FMSPC_TCB_DAO_ADDRESS=0xA8C0F6F6Deb3dA48Be03A99C112737000a5a3088
export TDX_ENCLAVE_IDENTITY_DAO_ADDRESS=0x5d1122a0d55b5095C0f03FBEa106A2e9722cb13F
export TDX_ENCLAVE_IDENTITY_HELPER_ADDRESS=0x870D17e2C12aF1C47dD1f0e4aFd36e28c830D558

# Set collateral paths (relative to packages/protocol)
export AZURE_TDX_TCB_INFO_PATH=/azure-tdx-assets/tcb.json
export AZURE_TDX_QE_IDENTITY_PATH=/azure-tdx-assets/qe_identity.json
export AZURE_TDX_PCS_CERT_PATH=/azure-tdx-assets/tdx_pcs_cert.hex
export AZURE_TDX_ROOT_PCS_CERT_PATH=/azure-tdx-assets/tdx_root_pcs_cert.hex

# Set the new owner (for ownership transfer after setup)
export NEW_OWNER=<timelock-or-owner-address>

# Run the setup script
./script/layer1/surge/setup_azure_tdx_verifier.sh
```

## Step 6: Verify Setup

```bash
# Check that collaterals are configured
# The setup script will output logs like:
# ** TDX_ROOT_PCS_CERTIFICATES configured
# ** TDX_PCS_CERTIFICATES configured
# ** TDX_QE_IDENTITY configured
# ** TDX_TCB_INFO configured
# ** AzureTdxVerifier ownership transferred to: <owner>
```

## Complete Example Script

Here's a complete script that runs all steps:

```bash
#!/bin/bash
# run_anvil_tdx_testnet.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$SCRIPT_DIR"

echo "=========================================="
echo "Step 1: Generate initial_state.json"
echo "=========================================="

cd "$REPO_ROOT"

# Generate state file from network_params.yaml
python3 << 'PYTHON_SCRIPT'
import yaml
import json
import sys

try:
    with open('network_params.yaml', 'r') as f:
        data = yaml.safe_load(f)
except Exception as e:
    print(f"Error reading network_params.yaml: {e}", file=sys.stderr)
    sys.exit(1)

preloaded = data.get('network_params', {}).get('additional_preloaded_contracts', {})

accounts = {}
for addr, contract in preloaded.items():
    accounts[addr.lower()] = {
        "nonce": contract.get("nonce", 1),
        "balance": contract.get("balance", "0x0"),
        "code": contract.get("code", "0x"),
        "storage": contract.get("storage", {})
    }

# Add funded accounts
for acc in ["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"]:
    if acc not in accounts:
        accounts[acc] = {"nonce": 0, "balance": "0x21e19e0c9bab2400000", "code": "0x", "storage": {}}

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

echo "=========================================="
echo "Step 2: Copy TDX test assets"
echo "=========================================="

mkdir -p azure-tdx-assets
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_tcb_90c06f000000.json azure-tdx-assets/tcb.json
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_identity.json azure-tdx-assets/qe_identity.json
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_pcs_cert.hex azure-tdx-assets/tdx_pcs_cert.hex
cp packages/protocol/test/layer1/automata-attestation/assets/0923/tdx_root_pcs_cert.hex azure-tdx-assets/tdx_root_pcs_cert.hex

echo "TDX assets copied to azure-tdx-assets/"

echo "=========================================="
echo "Step 3: Start Anvil"
echo "=========================================="

# Start anvil in background
anvil --load-state initial_state.json \
      --block-time 1 \
      --chain-id 31337 \
      --gas-limit 200000000 &

ANVIL_PID=$!
echo "Anvil started with PID: $ANVIL_PID"

# Wait for anvil to be ready
sleep 3

# Verify contracts are loaded
echo "Verifying preloaded contracts..."
PCS_CODE=$(cast code 0x928826C6D0d1986bD0465697984fa3722ADE16E1 2>/dev/null || echo "0x")
if [ "$PCS_CODE" != "0x" ]; then
    echo "PcsDao contract loaded successfully"
else
    echo "Warning: PcsDao contract not loaded"
fi

echo "=========================================="
echo "Anvil is running on http://localhost:8545"
echo "Press Ctrl+C to stop"
echo "=========================================="

# Wait for anvil process
wait $ANVIL_PID
```

## Troubleshooting

### Contract code not loaded

- Ensure `initial_state.json` has lowercase addresses
- Verify the `code` field is a valid hex string starting with `0x`

### Collateral setup fails

- Check that collateral files exist at the specified paths
- Verify JSON files are valid (use `jq . file.json` to validate)
- Ensure the certificate hex files have the `0x` prefix

### Anvil fails to start with state

- The `--load-state` and `--init` flags are mutually exclusive
- Ensure the state JSON matches Anvil's expected format

## References

- [Automata Network On-Chain PCCS](https://github.com/automata-network/automata-on-chain-pccs)
- [Automata DCAP Attestation](https://github.com/automata-network/automata-dcap-attestation)
- [Intel SGX DCAP](https://www.intel.com/content/www/us/en/developer/articles/technical/quote-verification-attestation-with-intel-sgx-dcap.html)
- [Foundry Anvil Reference](https://book.getfoundry.sh/anvil/)
