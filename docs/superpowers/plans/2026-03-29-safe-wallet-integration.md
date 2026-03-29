# Safe Wallet Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the custom UserOpsSubmitter smart wallet with Gnosis Safe (1-of-1) across both chains and the cross-chain DEX UI.

**Architecture:** Deploy Safe infrastructure (SafeL2, SafeProxyFactory, MultiSend, FallbackHandler) from the same deployer at matched nonces on L1 and L2 to get identical addresses. Replace the UI's wallet creation, signing, and transaction building to use Safe's `execTransaction` and `MultiSend`. Catalyst needs no changes — it treats calldata as opaque bytes forwarded to the submitter address.

**Tech Stack:** Solidity (Safe v1.5.0 contracts, paris EVM), TypeScript/React (viem), Rust (Catalyst — no changes needed)

**Branch:** New branch `feat/safe-wallet` off `surge-alethia-real-time-driver`

---

## File Map

| File | Action | Responsibility |
|------|--------|---------------|
| `packages/protocol/script/shared/surge/DeploySafeInfra.s.sol` | Create | Forge script to deploy all Safe contracts |
| `packages/protocol/script/shared/surge/deploy_safe_infra.sh` | Create | Shell wrapper for deployment on both chains |
| `packages/cross-chain-dex-ui/src/lib/contracts.ts` | Modify | Replace UserOpsSubmitter ABIs with Safe ABIs |
| `packages/cross-chain-dex-ui/src/lib/constants.ts` | Modify | Replace factory address, add Safe addresses |
| `packages/cross-chain-dex-ui/src/lib/safeOp.ts` | Create | Safe tx building, EIP-712, MultiSend encoding |
| `packages/cross-chain-dex-ui/src/lib/userOp.ts` | Modify | Update `sendUserOpToBuilder` to use Safe calldata |
| `packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts` | Modify | Use SafeProxyFactory for wallet creation |
| `packages/cross-chain-dex-ui/src/hooks/useUserOp.ts` | Modify | Build Safe txs instead of executeBatch |

---

### Task 1: Create new branch and install Safe contracts

**Files:**
- Modify: `packages/protocol/foundry.toml` (add remapping)
- Install: `lib/safe-smart-account`

- [ ] **Step 1: Create branch**

```bash
git checkout surge-alethia-real-time-driver
git checkout -b feat/safe-wallet
```

- [ ] **Step 2: Install Safe contracts**

```bash
cd packages/protocol
forge install safe-global/safe-smart-account
```

- [ ] **Step 3: Add remapping to foundry.toml**

Add to `[profile.layer1]` remappings in `packages/protocol/foundry.toml`:

```toml
"@safe/=lib/safe-smart-account/contracts/"
```

- [ ] **Step 4: Verify compilation**

```bash
FOUNDRY_PROFILE=layer1 forge build lib/safe-smart-account/contracts/SafeL2.sol --evm-version paris
```

Expected: Compiler run successful

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "chore: install safe-smart-account v1.5.0 contracts"
```

---

### Task 2: Write and run deployment script

**Files:**
- Create: `packages/protocol/script/shared/surge/DeploySafeInfra.s.sol`
- Create: `packages/protocol/script/shared/surge/deploy_safe_infra.sh`

- [ ] **Step 1: Create the Forge deployment script**

```solidity
// packages/protocol/script/shared/surge/DeploySafeInfra.s.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "forge-std/src/Script.sol";
import "forge-std/src/console2.sol";

import { SafeL2 } from "@safe/SafeL2.sol";
import { SafeProxyFactory } from "@safe/proxies/SafeProxyFactory.sol";
import { MultiSend } from "@safe/libraries/MultiSend.sol";
import { MultiSendCallOnly } from "@safe/libraries/MultiSendCallOnly.sol";
import { CompatibilityFallbackHandler } from "@safe/handler/CompatibilityFallbackHandler.sol";

contract DeploySafeInfra is Script {
    function run() external {
        vm.startBroadcast();

        SafeL2 singleton = new SafeL2();
        console2.log("SafeL2 singleton:", address(singleton));

        SafeProxyFactory factory = new SafeProxyFactory();
        console2.log("SafeProxyFactory:", address(factory));

        MultiSend multiSend = new MultiSend();
        console2.log("MultiSend:", address(multiSend));

        MultiSendCallOnly multiSendCallOnly = new MultiSendCallOnly();
        console2.log("MultiSendCallOnly:", address(multiSendCallOnly));

        CompatibilityFallbackHandler fallbackHandler = new CompatibilityFallbackHandler();
        console2.log("FallbackHandler:", address(fallbackHandler));

        vm.stopBroadcast();
    }
}
```

- [ ] **Step 2: Create the shell wrapper**

```bash
#!/bin/sh
# packages/protocol/script/shared/surge/deploy_safe_infra.sh
set -e

export PRIVATE_KEY=${PRIVATE_KEY:?"PRIVATE_KEY required"}
export L1_RPC=${L1_RPC:?"L1_RPC required"}
export L2_RPC=${L2_RPC:?"L2_RPC required"}
export FOUNDRY_PROFILE=layer1

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
```

- [ ] **Step 3: Generate a fresh deployer EOA**

```bash
cast wallet new
```

Save the address and private key. Fund it on both L1 and L2 (needs ~0.01 ETH/xDAI each).

- [ ] **Step 4: Deploy on both chains**

```bash
cd packages/protocol
PRIVATE_KEY=0x<key> L1_RPC=https://rpc.gnosis.gateway.fm L2_RPC=https://rpc.realtime.surge.wtf \
  bash script/shared/surge/deploy_safe_infra.sh
```

Verify: both chains show the same 5 addresses.

- [ ] **Step 5: Record deployed addresses**

Save the addresses in a `deployments/safe-infra.json`:

```json
{
  "safeL2Singleton": "0x...",
  "safeProxyFactory": "0x...",
  "multiSend": "0x...",
  "multiSendCallOnly": "0x...",
  "fallbackHandler": "0x..."
}
```

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: deploy Safe infrastructure on L1 and L2"
```

---

### Task 3: Add Safe ABIs and constants to UI

**Files:**
- Modify: `packages/cross-chain-dex-ui/src/lib/contracts.ts`
- Modify: `packages/cross-chain-dex-ui/src/lib/constants.ts`

- [ ] **Step 1: Add Safe ABIs to contracts.ts**

Replace `UserOpsSubmitterFactoryABI` and `UserOpsSubmitterABI` with:

```typescript
export const SafeProxyFactoryABI = [
  {
    type: 'function',
    name: 'createProxyWithNonce',
    inputs: [
      { name: '_singleton', type: 'address' },
      { name: 'initializer', type: 'bytes' },
      { name: 'saltNonce', type: 'uint256' },
    ],
    outputs: [{ name: 'proxy', type: 'address' }],
    stateMutability: 'nonpayable',
  },
] as const;

export const SafeABI = [
  {
    type: 'function',
    name: 'setup',
    inputs: [
      { name: '_owners', type: 'address[]' },
      { name: '_threshold', type: 'uint256' },
      { name: 'to', type: 'address' },
      { name: 'data', type: 'bytes' },
      { name: 'fallbackHandler', type: 'address' },
      { name: 'paymentToken', type: 'address' },
      { name: 'payment', type: 'uint256' },
      { name: 'paymentReceiver', type: 'address' },
    ],
    outputs: [],
    stateMutability: 'nonpayable',
  },
  {
    type: 'function',
    name: 'execTransaction',
    inputs: [
      { name: 'to', type: 'address' },
      { name: 'value', type: 'uint256' },
      { name: 'data', type: 'bytes' },
      { name: 'operation', type: 'uint8' },
      { name: 'safeTxGas', type: 'uint256' },
      { name: 'baseGas', type: 'uint256' },
      { name: 'gasPrice', type: 'uint256' },
      { name: 'gasToken', type: 'address' },
      { name: 'refundReceiver', type: 'address' },
      { name: 'signatures', type: 'bytes' },
    ],
    outputs: [{ name: 'success', type: 'bool' }],
    stateMutability: 'payable',
  },
  {
    type: 'function',
    name: 'nonce',
    inputs: [],
    outputs: [{ name: '', type: 'uint256' }],
    stateMutability: 'view',
  },
  {
    type: 'function',
    name: 'getOwners',
    inputs: [],
    outputs: [{ name: '', type: 'address[]' }],
    stateMutability: 'view',
  },
  {
    type: 'function',
    name: 'getTransactionHash',
    inputs: [
      { name: 'to', type: 'address' },
      { name: 'value', type: 'uint256' },
      { name: 'data', type: 'bytes' },
      { name: 'operation', type: 'uint8' },
      { name: 'safeTxGas', type: 'uint256' },
      { name: 'baseGas', type: 'uint256' },
      { name: 'gasPrice', type: 'uint256' },
      { name: 'gasToken', type: 'address' },
      { name: 'refundReceiver', type: 'address' },
      { name: '_nonce', type: 'uint256' },
    ],
    outputs: [{ name: '', type: 'bytes32' }],
    stateMutability: 'view',
  },
  {
    type: 'event',
    name: 'SafeSetup',
    inputs: [
      { name: 'initiator', type: 'address', indexed: true },
      { name: 'owners', type: 'address[]', indexed: false },
      { name: 'threshold', type: 'uint256', indexed: false },
      { name: 'initializer', type: 'address', indexed: false },
      { name: 'fallbackHandler', type: 'address', indexed: false },
    ],
  },
] as const;

export const MultiSendABI = [
  {
    type: 'function',
    name: 'multiSend',
    inputs: [{ name: 'transactions', type: 'bytes' }],
    outputs: [],
    stateMutability: 'payable',
  },
] as const;
```

- [ ] **Step 2: Update constants.ts with Safe addresses**

Replace `USER_OPS_FACTORY` with:

```typescript
export const SAFE_PROXY_FACTORY = import.meta.env.VITE_SAFE_PROXY_FACTORY as `0x${string}`;
export const SAFE_SINGLETON = import.meta.env.VITE_SAFE_SINGLETON as `0x${string}`;
export const SAFE_MULTISEND = import.meta.env.VITE_SAFE_MULTISEND as `0x${string}`;
export const SAFE_FALLBACK_HANDLER = import.meta.env.VITE_SAFE_FALLBACK_HANDLER as `0x${string}`;
```

Update `.env` and `.env.example` with the deployed addresses from Task 2.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "feat(ui): add Safe ABIs and contract constants"
```

---

### Task 4: Create safeOp.ts — Safe transaction builder

**Files:**
- Create: `packages/cross-chain-dex-ui/src/lib/safeOp.ts`

- [ ] **Step 1: Create safeOp.ts**

This file handles Safe-specific EIP-712 signing, MultiSend encoding, and `execTransaction` calldata building.

```typescript
// packages/cross-chain-dex-ui/src/lib/safeOp.ts
import { encodeFunctionData, encodePacked, Address, Hex, concat, pad, toHex, PublicClient } from 'viem';
import { SafeABI, MultiSendABI } from './contracts';
import { SAFE_MULTISEND } from './constants';

// Safe EIP-712 types (no name/version in domain — Safe uses chainId + verifyingContract only)
export const SafeTxTypes = {
  SafeTx: [
    { name: 'to', type: 'address' },
    { name: 'value', type: 'uint256' },
    { name: 'data', type: 'bytes' },
    { name: 'operation', type: 'uint8' },
    { name: 'safeTxGas', type: 'uint256' },
    { name: 'baseGas', type: 'uint256' },
    { name: 'gasPrice', type: 'uint256' },
    { name: 'gasToken', type: 'address' },
    { name: 'refundReceiver', type: 'address' },
    { name: 'nonce', type: 'uint256' },
  ],
} as const;

export function getSafeDomain(safeAddress: Address, chainId: number) {
  return {
    chainId,
    verifyingContract: safeAddress,
  };
}

export interface SafeTxParams {
  to: Address;
  value: bigint;
  data: Hex;
  operation: 0 | 1; // 0 = CALL, 1 = DELEGATECALL
}

/**
 * Build EIP-712 typed data for a Safe transaction.
 */
export function buildSafeTxTypedData(
  safeAddress: Address,
  chainId: number,
  nonce: bigint,
  tx: SafeTxParams
) {
  return {
    domain: getSafeDomain(safeAddress, chainId),
    types: SafeTxTypes,
    primaryType: 'SafeTx' as const,
    message: {
      to: tx.to,
      value: tx.value,
      data: tx.data,
      operation: tx.operation,
      safeTxGas: 0n,
      baseGas: 0n,
      gasPrice: 0n,
      gasToken: '0x0000000000000000000000000000000000000000' as Address,
      refundReceiver: '0x0000000000000000000000000000000000000000' as Address,
      nonce,
    },
  };
}

/**
 * Encode execTransaction calldata for a single Safe tx.
 */
export function buildExecTransactionCalldata(
  tx: SafeTxParams,
  signature: Hex
): Hex {
  return encodeFunctionData({
    abi: SafeABI,
    functionName: 'execTransaction',
    args: [
      tx.to,
      tx.value,
      tx.data,
      tx.operation,
      0n, // safeTxGas
      0n, // baseGas
      0n, // gasPrice
      '0x0000000000000000000000000000000000000000', // gasToken
      '0x0000000000000000000000000000000000000000', // refundReceiver
      signature,
    ],
  });
}

/**
 * Encode multiple calls into MultiSend format.
 * Each tx is packed as: operation(1) + to(20) + value(32) + dataLength(32) + data(variable)
 */
export function encodeMultiSend(txs: SafeTxParams[]): Hex {
  const encoded = txs.map((tx) =>
    encodePacked(
      ['uint8', 'address', 'uint256', 'uint256', 'bytes'],
      [tx.operation, tx.to, tx.value, BigInt(tx.data.length / 2 - 1), tx.data]
    )
  );
  return concat(encoded);
}

/**
 * Build a Safe tx that calls MultiSend with batched operations.
 * Returns the SafeTxParams to sign (operation = DELEGATECALL to MultiSend).
 */
export function buildMultiSendSafeTx(txs: SafeTxParams[]): SafeTxParams {
  const multiSendData = encodeFunctionData({
    abi: MultiSendABI,
    functionName: 'multiSend',
    args: [encodeMultiSend(txs)],
  });

  return {
    to: SAFE_MULTISEND,
    value: 0n,
    data: multiSendData,
    operation: 1, // DELEGATECALL
  };
}

/**
 * Read the current nonce from a Safe contract.
 */
export async function getSafeNonce(
  client: PublicClient,
  safeAddress: Address
): Promise<bigint> {
  return client.readContract({
    address: safeAddress,
    abi: SafeABI,
    functionName: 'nonce',
  });
}

/**
 * Build the initializer calldata for Safe.setup().
 */
export function buildSafeSetupCalldata(
  owner: Address,
  fallbackHandler: Address
): Hex {
  return encodeFunctionData({
    abi: SafeABI,
    functionName: 'setup',
    args: [
      [owner],                                                    // owners
      1n,                                                         // threshold
      '0x0000000000000000000000000000000000000000' as Address,    // to (no delegate call)
      '0x' as Hex,                                                // data
      fallbackHandler,                                            // fallbackHandler
      '0x0000000000000000000000000000000000000000' as Address,    // paymentToken
      0n,                                                         // payment
      '0x0000000000000000000000000000000000000000' as Address,    // paymentReceiver
    ],
  });
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat(ui): create safeOp.ts — Safe tx builder and MultiSend encoder"
```

---

### Task 5: Update useSmartWallet to use Safe Proxy Factory

**Files:**
- Modify: `packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts`

- [ ] **Step 1: Rewrite useSmartWallet**

Replace the entire hook to use `SafeProxyFactory.createProxyWithNonce` instead of `UserOpsSubmitterFactory.createSubmitter`.

Key changes:
- Call `createProxyWithNonce(singleton, setupCalldata, saltNonce)` where `saltNonce = uint256(keccak256(owner))`
- Read existing Safe by computing the expected CREATE2 address and checking if code exists
- For L2 creation: use bridge relay to call `createProxyWithNonce` on L2

The Safe proxy address is deterministic: `CREATE2(factory, salt, proxyCreationCode + singleton)` where `salt = keccak256(keccak256(initializer) + saltNonce)`.

To check if a Safe exists: `getCode(expectedAddress).length > 0` (simpler than a factory registry).

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat(ui): useSmartWallet uses Safe Proxy Factory"
```

---

### Task 6: Update useUserOp to build Safe transactions

**Files:**
- Modify: `packages/cross-chain-dex-ui/src/hooks/useUserOp.ts`
- Modify: `packages/cross-chain-dex-ui/src/lib/userOp.ts`

- [ ] **Step 1: Update executeGenericOps**

Replace `buildExecuteBatchTypedData` with `buildSafeTxTypedData`. For single ops, use `execTransaction` directly. For batched ops (approve + swap), use `buildMultiSendSafeTx` to wrap them.

Key changes:
- Fetch nonce from the Safe before signing: `getSafeNonce(client, safeAddress)`
- Build `SafeTx` typed data with the nonce
- Sign with `walletClient.signTypedData`
- Build `execTransaction` calldata with the signature
- Send to builder via `sendUserOpToBuilder(safeAddress, calldata, chainId)`

- [ ] **Step 2: Update sendUserOpToBuilder**

Change the function signature — instead of encoding `executeBatch(ops, signature)` calldata internally, accept pre-built calldata:

```typescript
export async function sendUserOpToBuilder(
  submitter: Address,
  calldata: Hex,
  chainId?: number
): Promise<{ success: boolean; result?: unknown; error?: string; userOpId?: number }>
```

The caller (useUserOp) now builds the full `execTransaction` calldata and passes it directly.

- [ ] **Step 3: Update op builders**

Update `buildSwapUserOps`, `buildBridgeUserOps`, `buildBridgeNativeUserOps`, `buildBridgeOutNativeUserOps`, `buildCreateL2WalletUserOps` to return `SafeTxParams[]` instead of `UserOp[]`. The structure is the same (`to`, `value`, `data`) — just add `operation: 0` (CALL) to each.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat(ui): useUserOp builds Safe execTransaction instead of executeBatch"
```

---

### Task 7: Update L2 wallet creation to use Safe

**Files:**
- Modify: `packages/cross-chain-dex-ui/src/lib/safeOp.ts` (add L2 creation builder)

- [ ] **Step 1: Add buildCreateL2SafeUserOps**

The L2 Safe creation goes through the bridge relay (same pattern as before):
- L1 Safe calls `bridge.sendMessage` targeting the relay on L2
- Relay calls `SafeProxyFactory.createProxyWithNonce(singleton, initializer, saltNonce)` on L2
- Same params = same address

```typescript
export function buildCreateL2SafeOps(
  owner: Address,
  safeProxyFactory: Address,
  safeSingleton: Address,
  fallbackHandler: Address,
  l1Bridge: Address,
  l2Relay: Address,
  l2ChainId: number,
  sender: Address
): SafeTxParams[] {
  const initializer = buildSafeSetupCalldata(owner, fallbackHandler);
  const saltNonce = BigInt(keccak256(encodePacked(['address'], [owner])));

  const createProxyCalldata = encodeFunctionData({
    abi: SafeProxyFactoryABI,
    functionName: 'createProxyWithNonce',
    args: [safeSingleton, initializer, saltNonce],
  });

  // Encode for relay: abi.encode(target, calldata)
  const relayPayload = encodeAbiParameters(
    [{ type: 'address' }, { type: 'bytes' }],
    [safeProxyFactory, createProxyCalldata],
  );

  const onMessageInvocationData = encodeFunctionData({
    abi: [{ type: 'function', name: 'onMessageInvocation', inputs: [{ name: '_data', type: 'bytes' }], outputs: [], stateMutability: 'payable' }],
    functionName: 'onMessageInvocation',
    args: [relayPayload],
  });

  // Bridge sendMessage targeting relay on L2
  return [{
    to: l1Bridge,
    value: 0n,
    data: encodeFunctionData({
      abi: BridgeABI,
      functionName: 'sendMessage',
      args: [{
        id: 0n, fee: 0n, gasLimit: 2_000_000,
        from: '0x0000000000000000000000000000000000000000' as Address,
        srcChainId: 0n, srcOwner: sender,
        destChainId: BigInt(l2ChainId), destOwner: sender,
        to: l2Relay, value: 0n, data: onMessageInvocationData,
      }],
    }),
    operation: 0,
  }];
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat(ui): L2 Safe creation via bridge relay"
```

---

### Task 8: End-to-end testing

- [ ] **Step 1: Deploy Safe infra on both chains** (Task 2)
- [ ] **Step 2: Update .env with deployed addresses**
- [ ] **Step 3: Start catalyst + UI**
- [ ] **Step 4: Test wallet creation on L1** — connect wallet, create Safe
- [ ] **Step 5: Test funding** — send xDAI to Safe on L1
- [ ] **Step 6: Test L2 wallet creation** — bridge relay creates Safe on L2 at same address
- [ ] **Step 7: Test bridge-in** — deposit xDAI to L2 Safe
- [ ] **Step 8: Test bridge-out** — withdraw from L2 Safe (switch to L2, sign, Ambire validates)
- [ ] **Step 9: Test swap** — swap xDAI for USDC via Safe

- [ ] **Step 10: Commit and push**

```bash
git add -A
git commit -m "feat: Safe wallet integration complete"
git push origin feat/safe-wallet
```
