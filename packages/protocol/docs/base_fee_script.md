# Base Fee Parameter Calculator Documentation

## Overview

This script calculates optimal parameters for the modified EIP-1559 base fee mechanism.

## Complete Script

```python
import math

# Predefined. Do not change.
ADJUSTMENT_QUOTIENT = 8

# In seconds
BLOCK_TIME = 1

BLOCK_GAS_LIMIT = 200_000_000

MIN_BASE_FEE_GWEI = 0.1

# Block gas target per second
GAS_ISSUANCE_PER_SECOND = BLOCK_GAS_LIMIT // BLOCK_TIME // 2
print("gasIssuancePerSecond: ", GAS_ISSUANCE_PER_SECOND)

target = GAS_ISSUANCE_PER_SECOND
aq = ADJUSTMENT_QUOTIENT
base_fee_wei = MIN_BASE_FEE_GWEI * pow(10, 9)

# Required to set a minimum base fee for every block
MIN_GAS_EXCESS = target * aq * math.log(base_fee_wei * target * aq)

# Reduce to cleaner number
MIN_GAS_EXCESS = (MIN_GAS_EXCESS // 1_000_000) * 1_000_000
eventual_base_fee = math.exp(MIN_GAS_EXCESS / (target * aq)) / (target * aq) / pow(10, 9)
print("minGasExcess: {} (Base fee: {} Gwei)".format(MIN_GAS_EXCESS, eventual_base_fee))
```

## Required Inputs

- `ADJUSTMENT_QUOTIENT`: Controls how quickly the base fee adjusts to changes in demand
- `BLOCK_TIME`: Target time between blocks (in seconds)
- `BLOCK_GAS_LIMIT`: Maximum gas limit per block
- `MIN_BASE_FEE_GWEI`: Desired minimum base fee in Gwei

## Outputs

1. `gasIssuancePerSecond`:

   - Similar to EIP-1559's block gas target, but measured per second
   - Calculated as a portion of the block gas limit divided by block time
   - Represents the target network gas consumption rate

2. `minGasExcess`:
   - Used to enforce a minimum base fee for L2 blocks
   - Calculated using the gas issuance target and desired minimum base fee
   - When network usage is at or below this threshold, the base fee will stabilize around the minimum value

## Usage in Smart Contracts

The calculated values can be used in the base fee configuration of L2 contracts:

```solidity
baseFeeConfig: LibSharedData.BaseFeeConfig({
    adjustmentQuotient: ADJUSTMENT_QUOTIENT,
    sharingPctg: 0,
    gasIssuancePerSecond: GAS_ISSUANCE_PER_SECOND,
    minGasExcess: MIN_GAS_EXCESS,
    maxGasIssuancePerBlock: maxGasIssuancePerBlock  // Set as needed
})
```
