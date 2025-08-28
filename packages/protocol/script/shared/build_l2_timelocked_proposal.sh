#!/bin/sh

# This script builds a timelocked proposal for sending a cross-chain message to the L2 DelegateOwner contract
set -e

# L1 configuration
export L1_BRIDGE=${L1_BRIDGE:-"0xe736998f06239e3F2c7900Ebd4Ce4D8a89E9584A"}
export L1_TIMELOCK_CONTROLLER=${L1_TIMELOCK_CONTROLLER:-"0x468D938578431337C81F572090e65EB1101Cf21A"}

# L2 configuration
export L2_CHAINID=${L2_CHAINID:-167004}
export L2_DELEGATE_OWNER=${L2_DELEGATE_OWNER:-"0x175fB89509f96097ba2aea0feDfdD2Da01ed58Df"}

# Message configuration
export TARGET_CONTRACT=${TARGET_CONTRACT:-"0xb13cA5Fc7F10598Ed8F774Ee0762e297Ce2e385c"}
export IS_DELEGATE_CALL=${IS_DELEGATE_CALL:-false}
export CALL_DATA=${CALL_DATA:-"0x60fe47b10000000000000000000000000000000000000000000000000000000000000005"}
export GAS_LIMIT=${GAS_LIMIT:-1000000}
export FEE=${VALUE:-100000000000} # 100 gwei
export VALUE=${VALUE:-0}

# Timelock configuration
export TIMELOCK_DELAY=${TIMELOCK_DELAY:-600} # 10 minutes

# Foundry profile
export FOUNDRY_PROFILE=${FOUNDRY_PROFILE:-"shared"}

# Parameterize log level
export LOG_LEVEL=${LOG_LEVEL:-"-vvvv"}

# Parameterize block gas limit
export BLOCK_GAS_LIMIT=${BLOCK_GAS_LIMIT:-200000000}

# Validate required environment variables
if [ "$L1_BRIDGE" = "0x0000000000000000000000000000000000000000" ]; then
    echo "Error: L1_BRIDGE environment variable is required"
    exit 1
fi

if [ "$L1_TIMELOCK_CONTROLLER" = "0x0000000000000000000000000000000000000000" ]; then
    echo "Error: L1_TIMELOCK_CONTROLLER environment variable is required"
    exit 1
fi

if [ "$L2_DELEGATE_OWNER" = "0x0000000000000000000000000000000000000000" ]; then
    echo "Error: L2_DELEGATE_OWNER environment variable is required"
    exit 1
fi

if [ "$TARGET_CONTRACT" = "0x0000000000000000000000000000000000000000" ]; then
    echo "Error: TARGET_CONTRACT environment variable is required"
    exit 1
fi

# Build L2 timelocked proposal
forge script ./script/shared/L2TimelockedProposalBuilder.sol:L2TimelockedProposalBuilder \
    --ffi \
    $LOG_LEVEL \
    --block-gas-limit $BLOCK_GAS_LIMIT
