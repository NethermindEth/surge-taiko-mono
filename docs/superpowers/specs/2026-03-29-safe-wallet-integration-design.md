# Safe Wallet Integration Design

## Goal

Replace the custom `UserOpsSubmitter` smart wallet with Gnosis Safe (1-of-1) across the Surge cross-chain DEX stack. This gives us ERC-1271, battle-tested security, wallet compatibility (Ambire, MetaMask), and a proper smart account standard.

## Context

### Why

The custom `UserOpsSubmitter` lacks:
- ERC-1271 (`isValidSignature`) — Ambire's post-sign validation fails on L2 because its deployless validator uses Cancun opcodes (MCOPY) unsupported on our L2
- Nonce management — current contract explicitly says "No nonce checks. Do not use in production."
- Ecosystem compatibility — wallets, block explorers, and tools don't recognize it

Safe gives us all of this out of the box.

### Current Architecture

- `UserOpsSubmitterFactory` deploys wallets via CREATE2 (same address on L1 + L2)
- `UserOpsSubmitter.executeBatch(ops[], signature)` verifies EIP-712 signature and forwards calls
- Catalyst receives UserOps via `surge_sendUserOp` RPC, decodes `executeBatch` calldata
- UI builds `ExecuteBatch` EIP-712 typed data, signs with wallet extension

### What Changes

| Component | Current | New |
|-----------|---------|-----|
| Smart wallet | UserOpsSubmitter | Safe 1-of-1 proxy |
| Factory | UserOpsSubmitterFactory (CREATE2) | SafeProxyFactory (CREATE2) |
| Execution | `executeBatch(ops[], sig)` | `execTransaction(to, value, data, ...)` |
| Batching | Built into executeBatch | MultiSend contract |
| Signature | Custom EIP-712 `ExecuteBatch` type | Safe's EIP-712 `SafeTx` type |
| ERC-1271 | Not supported | Built-in |
| Nonces | None | Built-in (prevents replay) |

## Architecture

### Deployment (both chains)

Deploy from the **same EOA at matched nonces** on L1 (Gnosis) and L2 (Surge) to get identical addresses:

1. **Safe Singleton** (`SafeL2.sol`) — the implementation contract
2. **SafeProxyFactory** — deploys proxy wallets via CREATE2
3. **MultiSend** — batches multiple calls into one `execTransaction`
4. **CompatibilityFallbackHandler** — provides ERC-1271 + other interfaces

All compiled with `--evm-version paris` (no PUSH0/MCOPY) since L2 doesn't support Cancun.

### Wallet Creation

**L1 creation:**
- UI calls `SafeProxyFactory.createProxyWithNonce(singleton, initializer, saltNonce)`
- `initializer` = `Safe.setup([ownerEOA], 1, ...)` — 1-of-1 with the user's EOA as owner
- `saltNonce` = deterministic (e.g. `uint256(keccak256(ownerEOA))`)
- CREATE2 produces a deterministic address

**L2 creation:**
- Same flow, triggered via bridge relay (existing `CrossChainRelay` pattern)
- L1 smart wallet calls `bridge.sendMessage` → relay → `SafeProxyFactory.createProxyWithNonce` on L2
- Same singleton address + same initializer + same saltNonce = same proxy address on L2

### Transaction Execution

**Single call (e.g. bridge deposit):**
```
Safe.execTransaction(
  to: bridgeAddress,
  value: amount,
  data: sendMessage(...),
  operation: CALL,
  safeTxGas: 0,
  baseGas: 0,
  gasPrice: 0,
  gasToken: address(0),
  refundReceiver: address(0),
  signatures: ownerSignature
)
```

**Batched calls (e.g. approve + swap):**
```
Safe.execTransaction(
  to: multiSendAddress,
  value: 0,
  data: MultiSend.multiSend(encodedTransactions),
  operation: DELEGATECALL,
  ...
  signatures: ownerSignature
)
```

### EIP-712 Signature

Safe uses its own EIP-712 domain and `SafeTx` type:

```
Domain: {
  chainId,
  verifyingContract: safeAddress
}

SafeTx: {
  to, value, data, operation,
  safeTxGas, baseGas, gasPrice,
  gasToken, refundReceiver, nonce
}
```

The `nonce` comes from `Safe.nonce()` — auto-incrementing, prevents replay.

### Catalyst Changes

The `surge_sendUserOp` RPC receives:
```json
{
  "submitter": "0x...",     // Safe address
  "calldata": "0x...",      // execTransaction ABI-encoded calldata
  "chainId": 100            // target chain (0 or L1 = default L1)
}
```

Catalyst changes:
- Decode `execTransaction` instead of `executeBatch` in the L1 simulation path
- For L2 direct UserOps: construct L2 tx calling `execTransaction` on the L2 Safe
- The rest (signal slots, bridge relay, L1 multicall) stays the same

### UI Changes

- `useSmartWallet` — call `SafeProxyFactory.createProxyWithNonce` instead of `createSubmitter`
- `useUserOp` / `userOp.ts`:
  - Build `SafeTx` typed data instead of `ExecuteBatch`
  - Use `MultiSend.multiSend` encoding for batched ops
  - Fetch nonce from `Safe.nonce()` before signing
  - `sendUserOpToBuilder` sends `execTransaction` calldata
- `BridgeCard` — no changes (ops are abstracted behind the hooks)
- L2 wallet creation — bridge relay calls `SafeProxyFactory.createProxyWithNonce` on L2

### Cross-Chain Relay

Same pattern as current: L1 Safe calls `bridge.sendMessage` targeting the `CrossChainRelay` on L2, which forwards to `SafeProxyFactory.createProxyWithNonce`. The relay is already deployed and working.

## Contract Deployment Plan

All compiled with `--evm-version paris`.

| Contract | Source | Purpose |
|----------|--------|---------|
| SafeL2 | `@safe-global/safe-contracts/SafeL2.sol` | Singleton implementation |
| SafeProxyFactory | `@safe-global/safe-contracts/proxies/SafeProxyFactory.sol` | CREATE2 proxy deployer |
| MultiSend | `@safe-global/safe-contracts/libraries/MultiSend.sol` | Batch execution |
| CompatibilityFallbackHandler | `@safe-global/safe-contracts/handler/CompatibilityFallbackHandler.sol` | ERC-1271 + token callbacks |

Deploy from same EOA at matched nonces on both chains.

## Branch Strategy

New branch off `surge-alethia-real-time-driver` (or `feat/bridge-out-create2`). Does NOT modify the current UserOpsSubmitter code — clean replacement. The existing branch continues to work with UserOpsSubmitter.

## Success Criteria

1. Safe proxy deployed at same address on L1 and L2
2. Bridge-in (L1→L2 deposit) works through Safe
3. Bridge-out (L2→L1 withdrawal) works — Ambire signs on L2 without validation errors
4. Swap works through Safe
5. L2 wallet creation via bridge relay works
6. Nonce management prevents replay
