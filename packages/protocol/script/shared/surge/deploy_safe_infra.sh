#!/bin/sh
set -e

export PRIVATE_KEY=${PRIVATE_KEY:?"PRIVATE_KEY required"}
export L1_RPC=${L1_RPC:?"L1_RPC required"}
export L2_RPC=${L2_RPC:?"L2_RPC required"}
export FOUNDRY_PROFILE=shared

echo "=============================="
echo "Deploying Safe Infrastructure"
echo "=============================="

echo "=== L1 ==="
forge script ./script/shared/surge/DeploySafeInfra.s.sol:DeploySafeInfra \
    --fork-url $L1_RPC \
    --broadcast \
    --evm-version paris \
    --private-key $PRIVATE_KEY \
    -vvvv

echo "=== L2 ==="
forge script ./script/shared/surge/DeploySafeInfra.s.sol:DeploySafeInfra \
    --fork-url $L2_RPC \
    --broadcast \
    --evm-version paris \
    --private-key $PRIVATE_KEY \
    -vvvv

echo "Done! Verify addresses match on both chains."
