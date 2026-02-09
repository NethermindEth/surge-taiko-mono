#!/bin/bash

# Generate go contract bindings.
# ref: https://geth.ethereum.org/docs/dapp/native-bindings
#
# Usage:
#   TAIKO_GETH_DIR=/path/to/taiko-geth ./gen_bindings.sh [shasta|surge]
#
# Examples:
#   TAIKO_GETH_DIR=~/code/taiko-geth ./gen_bindings.sh shasta
#   TAIKO_GETH_DIR=~/code/taiko-geth ./gen_bindings.sh surge

set -eou pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"

# Default to shasta if no fork specified
FORK=${1:-shasta}

# Validate fork name
if [[ "$FORK" != "shasta" && "$FORK" != "surge" ]]; then
  echo "Error: Invalid fork name '$FORK'. Must be 'shasta' or 'surge'."
  exit 1
fi

echo ""
echo "=========================================="
echo "Generating Go contract bindings for: ${FORK}"
echo "=========================================="
echo ""
echo "TAIKO_GETH_DIR: ${TAIKO_GETH_DIR}"
echo ""

# Build abigen from taiko-geth
cd ${TAIKO_GETH_DIR} &&
  make all &&
  cd -

# Compile protocol contracts
cd ../protocol &&
  pnpm clean &&
  pnpm compile &&
  cd -

ABIGEN_BIN=$TAIKO_GETH_DIR/build/bin/abigen

# Create bindings directory if it doesn't exist
mkdir -p $DIR/../bindings/${FORK}

echo ""
echo "Start generating Go contract bindings for ${FORK}..."
echo ""

# Generate bindings based on fork
if [[ "$FORK" == "shasta" ]]; then
  # Shasta bindings (from taiko-alethia-protocol-v3.0.0)
  cat ../protocol/out/layer1/MainnetInbox.sol/MainnetInbox.json |
    jq .abi |
    ${ABIGEN_BIN} --abi - --type ShastaInboxClient --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_shasta_inbox.go

  cat ../protocol/out/layer2/Anchor.sol/Anchor.json |
    jq .abi |
    ${ABIGEN_BIN} --abi - --type ShastaAnchor --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_shasta_anchor.go

  cat ../protocol/out/layer1/ComposeVerifier.sol/ComposeVerifier.json |
    jq .abi |
    ${ABIGEN_BIN} --abi - --type ComposeVerifier --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_compose_verifier.go

elif [[ "$FORK" == "surge" ]]; then
  # Surge bindings (from surge-alethia-protocol-v3.0.0)
  # These protocols may differ from taiko-alethia-protocol-v3.0.0

  # Core Inbox contract
  cat ../protocol/out/layer1/MainnetInbox.sol/MainnetInbox.json |
    jq .abi |
    ${ABIGEN_BIN} --abi - --type SurgeInboxClient --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_surge_inbox.go
  echo "  Generated: gen_surge_inbox.go"

  # L2 Anchor contract
  cat ../protocol/out/layer2/Anchor.sol/Anchor.json |
    jq .abi |
    ${ABIGEN_BIN} --abi - --type SurgeAnchor --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_surge_anchor.go
  echo "  Generated: gen_surge_anchor.go"

  # ComposeVerifier (if exists)
  if [[ -f "../protocol/out/layer1/ComposeVerifier.sol/ComposeVerifier.json" ]]; then
    cat ../protocol/out/layer1/ComposeVerifier.sol/ComposeVerifier.json |
      jq .abi |
      ${ABIGEN_BIN} --abi - --type ComposeVerifier --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_compose_verifier.go
    echo "  Generated: gen_compose_verifier.go"
  fi

  # Surge-specific contracts
  if [[ -f "../protocol/out/layer1/SurgeVerifier.sol/SurgeVerifier.json" ]]; then
    cat ../protocol/out/layer1/SurgeVerifier.sol/SurgeVerifier.json |
      jq .abi |
      ${ABIGEN_BIN} --abi - --type SurgeVerifier --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_surge_verifier.go
    echo "  Generated: gen_surge_verifier.go"
  fi

  if [[ -f "../protocol/out/layer1/LibProofBitmap.sol/LibProofBitmap.json" ]]; then
    cat ../protocol/out/layer1/LibProofBitmap.sol/LibProofBitmap.json |
      jq .abi |
      ${ABIGEN_BIN} --abi - --type LibProofBitmap --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_lib_proof_bitmap.go
    echo "  Generated: gen_lib_proof_bitmap.go"
  fi

  # BondManager interface
  if [[ -f "../protocol/out/layer1/IBondManager.sol/IBondManager.json" ]]; then
    cat ../protocol/out/layer1/IBondManager.sol/IBondManager.json |
      jq .abi |
      ${ABIGEN_BIN} --abi - --type BondManager --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_bond_manager.go
    echo "  Generated: gen_bond_manager.go"
  fi

  # CCIP StateStore (Surge-specific)
  if [[ -f "../protocol/out/layer1/CCIPStateStore.sol/CCIPStateStore.json" ]]; then
    cat ../protocol/out/layer1/CCIPStateStore.sol/CCIPStateStore.json |
      jq .abi |
      ${ABIGEN_BIN} --abi - --type CCIPStateStore --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_ccip_state_store.go
    echo "  Generated: gen_ccip_state_store.go"
  fi

  # SurgeTimelockController (Surge-specific)
  if [[ -f "../protocol/out/layer1/SurgeTimelockController.sol/SurgeTimelockController.json" ]]; then
    cat ../protocol/out/layer1/SurgeTimelockController.sol/SurgeTimelockController.json |
      jq .abi |
      ${ABIGEN_BIN} --abi - --type SurgeTimelockController --pkg ${FORK} --out $DIR/../bindings/${FORK}/gen_surge_timelock_controller.go
    echo "  Generated: gen_surge_timelock_controller.go"
  fi
fi

# Record the git commit for traceability
git -C ../../ log --format="%H" -n 1 >./bindings/${FORK}/.githead

echo ""
echo "🍻 Go contract bindings for ${FORK} generated!"
echo "   Bindings location: bindings/${FORK}/"
echo "   Git HEAD recorded: $(cat ./bindings/${FORK}/.githead)"
