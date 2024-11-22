#!/bin/sh

# This script is only used by `pnpm test:deploy:l1`.
set -e

export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
L2_CHAINID=123781 \
L2_GENESIS_HASH=0xee1950562d42f0da28bd4550d88886bc90894c77c9c9eaefef775d4c8223f259 \
OWNER_MULTISIG=0x1237810000000000000000000000000000000001 \
OWNER_MULTISIG_SIGNERS="0x1237810000000000000000000000000000000002,0x1237810000000000000000000000000000000003,0x1237810000000000000000000000000000000004" \
FOUNDRY_PROFILE="layer1" \
forge script ./script/layer1/DeploySurgeOnL1.s.sol:DeploySurgeOnL1 \
    --fork-url http://localhost:8545 \
    --broadcast \
    --ffi \
    -vvvv \
    --private-key $PRIVATE_KEY \
    --block-gas-limit 200000000
