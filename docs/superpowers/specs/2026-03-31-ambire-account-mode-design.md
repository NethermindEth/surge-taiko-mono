# Ambire Account Mode — Design Spec

## Overview

Add an optional AmbireAccount (EIP-7702) wallet mode to the Surge DEX UI alongside the existing Safe wallet mode. When a user connects an EOA that has a 7702 delegation to AmbireAccount, the UI offers a choice between Safe mode (current behavior) and Ambire Account mode. The user's choice is persisted per address.

## Motivation

Ambire wallet now supports Safe accounts, but its `signTypedData` implementation wraps all signatures in a SafeMessage EIP-712 envelope for Safe accounts. This makes the resulting signature incompatible with Safe's `execTransaction`, which expects a raw ECDSA signature over the SafeTx hash. Rather than fighting this limitation, we support AmbireAccount natively — users connect with their EOA (which has 7702 delegation to AmbireAccount) and sign via `personal_sign`, producing signatures compatible with AmbireAccount's `execute()`.

## Two Wallet Modes

| Aspect                 | Safe mode (default, unchanged)               | Ambire Account mode (new)                          |
| ---------------------- | -------------------------------------------- | -------------------------------------------------- |
| L1 smart wallet        | Safe proxy (created by UI)                   | AmbireAccount via 7702 delegation (pre-existing)   |
| L2 actions             | Safe on L2 (bridged creation)                | Raw EOA (no 7702 on L2)                            |
| Signing method         | `signTypedData` (SafeTx EIP-712)             | `personal_sign` (raw execute hash, EthSign Mode 1) |
| Catalyst calldata      | `Safe.execTransaction(...)`                  | `AmbireAccount.execute(txns, signature)`           |
| Wallet creation needed | Yes — deploy Safe on L1, bridge-create on L2 | None                                               |
| Balance reads          | From Safe address                            | From EOA address                                   |

## Detection & Mode Selection

### 7702 Delegation Detection

1. On wallet connect, call `getCode(eoaAddress)` on L1.
2. Check if the returned bytecode starts with `0xef0100` (EIP-7702 delegation designator).
3. If so, extract the 20-byte delegation target address from bytes 3..23.
4. Verify the target is AmbireAccount (check against known deployment address, or `getCode(target)` and check for `execute` selector).
5. If delegation target is NOT AmbireAccount, default to Safe mode without offering choice.

### Mode Selection UX

- If 7702 delegation to AmbireAccount is detected: show `AccountModeSelector` modal with two options:
  - **Safe Wallet** (default) — "Creates a dedicated Safe. Works with any wallet."
  - **Ambire Account** — "Uses your existing 7702 smart account. Native batching, no extra wallet needed."
- If no delegation detected: skip modal, proceed with Safe flow.
- Persist choice in `localStorage` keyed by EOA address (e.g., `surge_account_mode_<address>`).
- Modal only shown on first connect per address (or if stored preference is cleared).

## Signing & Calldata Flow

### L1 Operations (Ambire Account mode)

Applies to: swap, bridge, add/remove liquidity, bridge native.

1. Build `Transaction[]` array — same inner operations as Safe mode (approve, swap, bridge calls etc.), structured as AmbireAccount `Transaction{to, value, data}`.
2. Read current `nonce` from AmbireAccount contract on L1 (`publicClient.readContract({functionName: 'nonce'})`).
3. Compute execute hash: `keccak256(abi.encode(eoaAddress, chainId, nonce, txns))`.
4. Sign via `walletClient.signMessage(executeHash)` — this is `personal_sign`, producing ECDSA over `keccak256("\x19Ethereum Signed Message:\n32" + executeHash)`.
5. Append EthSign mode byte `0x01` to the 65-byte signature.
6. Encode `execute(Transaction[],bytes)` calldata with `(txns, signature)`.
7. Send to Catalyst: `surge_sendUserOp({submitter: eoaAddress, calldata})`.

### L2 Operations (Ambire Account mode)

No 7702 delegation on L2, so the EOA acts directly:

1. Build the transaction (e.g., bridge `sendMessage`).
2. Send as raw EOA call data to Catalyst: `surge_sendUserOp({submitter: eoaAddress, calldata})`.

### AmbireAccount Signature Verification (on-chain)

AmbireAccount's `execute()` computes:

```
hash = keccak256(abi.encode(address(this), block.chainid, currentNonce, txns))
```

Passes to `SignatureValidator.recoverAddrImpl(hash, signature, true)`.

With mode byte `0x01` (EthSign), the validator does:

```
hash = keccak256("\x19Ethereum Signed Message:\n32", hash)
signer = ecrecover(hash, v, r, s)
```

This matches `personal_sign` behavior — the EOA's signature is verified correctly, and `privileges[signer]` is checked.

## UI Component Changes

### `useSmartWallet` hook

- New state: `accountMode: "safe" | "ambire"`.
- Reads persisted mode from localStorage on connect.
- **Safe mode**: existing behavior unchanged.
- **Ambire mode**:
  - `smartWallet` = connected EOA address (the AmbireAccount IS the EOA).
  - `l2WalletExists` = always `true` (L2 uses raw EOA).
  - Skips all Safe creation/detection logic (no CREATE2 prediction, no factory calls).

### `SmartWalletSetup` modal

Not shown in Ambire mode (no wallet to create).

### `FundWallet` modal

Still shown (user still needs funds in the EOA), but skips the "Create L2 Wallet" step in Ambire mode.

### New `AccountModeSelector` modal

- Shown once when 7702 delegation to AmbireAccount is detected.
- Two options: Safe Wallet / Ambire Account.
- Persists choice to localStorage.
- Only appears on first connect per address.

### `useUserOp` hook

`executeGenericOps` branches by mode:

- **Safe mode**: existing flow — `getSafeNonce` → `buildSafeTxTypedData` → `signTypedData` → `buildExecTransactionCalldata` → `sendUserOpToBuilder`.
- **Ambire mode**: new flow — `getAmbireNonce` → compute execute hash → `signMessage` → append mode byte → encode `execute` calldata → `sendUserOpToBuilder`.

### Swap/Bridge/Liquidity cards

No changes. They call the same hook functions. Mode branching is internal to hooks.

## Catalyst Changes

Minimal change — execution routing by function selector.

When Catalyst receives `surge_sendUserOp({submitter, calldata})`:

1. Read first 4 bytes of `calldata` (function selector).
2. If `execTransaction` selector (`0x6a761202`) → existing Safe execution path.
3. If `execute(Transaction[],bytes)` selector (`0x6171d1c9`) → new AmbireAccount execution path:
   - Decode `(Transaction[], bytes)` from calldata.
   - Call `submitter.execute(txns, signature)` on-chain.
   - Same proving/proposing lifecycle as today.

No new RPC methods. No schema changes. Same `surge_userOpStatus` polling works for both modes.

## What's NOT Changing

- Builder RPC API (`surge_sendUserOp` / `surge_userOpStatus`)
- DEX reserves, swap quotes (pure read operations)
- L1 contracts (Vault, Bridge, USDC)
- L2 contracts (SimpleDEX, Bridge)
- RainbowKit / wagmi config
- TxStatusOverlay / status polling
- Connection flow (same wagmi `useAccount`)

## AmbireAccount Contract Reference

Contract: `AmbireAccount.sol` from [AmbireTech/wallet](https://github.com/AmbireTech/wallet)

Key interface:

```solidity
struct Transaction {
    address to;
    uint value;
    bytes data;
}

function execute(Transaction[] calldata txns, bytes calldata signature) public;
function nonce() public view returns (uint);
function privileges(address) public view returns (bytes32);
```

Signature format: `{r}{s}{v}{mode_byte}` where mode byte `0x01` = EthSign.

SignatureValidator modes:

- `0x00` (EIP712): `ecrecover(hash, v, r, s)` — direct
- `0x01` (EthSign): `ecrecover(keccak256("\x19Ethereum Signed Message:\n32" + hash), v, r, s)`
- `0x02` (SmartWallet): ERC-1271 `isValidSignature` call
