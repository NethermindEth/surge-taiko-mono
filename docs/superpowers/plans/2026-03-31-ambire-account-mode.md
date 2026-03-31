# Ambire Account Mode Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add optional AmbireAccount (EIP-7702) wallet mode to the Surge DEX UI, auto-detected when an EOA has 7702 delegation, allowing users to choose between Safe and AmbireAccount execution paths.

**Architecture:** The `useSmartWallet` hook gains an `accountMode` state that controls all downstream behavior. A new `ambireOp.ts` library handles AmbireAccount-specific hash computation, signing, and calldata encoding. The existing `useUserOp` hook branches on `accountMode` to use either the Safe or Ambire execution path. A new `AccountModeSelector` modal appears when 7702 delegation is detected.

**Tech Stack:** React 18, viem, wagmi, TypeScript

**Spec:** `docs/superpowers/specs/2026-03-31-ambire-account-mode-design.md`

---

### Task 1: Add AccountMode type and AmbireAccount ABI

**Files:**

- Modify: `packages/cross-chain-dex-ui/src/types/index.ts`
- Modify: `packages/cross-chain-dex-ui/src/lib/contracts.ts`

- [ ] **Step 1: Add AccountMode type**

In `packages/cross-chain-dex-ui/src/types/index.ts`, add at the end:

```typescript
export type AccountMode = "safe" | "ambire";
```

- [ ] **Step 2: Add AmbireAccount ABI**

In `packages/cross-chain-dex-ui/src/lib/contracts.ts`, add at the end:

```typescript
export const AmbireAccountABI = [
  {
    type: "function",
    name: "execute",
    inputs: [
      {
        name: "txns",
        type: "tuple[]",
        components: [
          { name: "to", type: "address" },
          { name: "value", type: "uint256" },
          { name: "data", type: "bytes" },
        ],
      },
      { name: "signature", type: "bytes" },
    ],
    outputs: [],
    stateMutability: "nonpayable",
  },
  {
    type: "function",
    name: "nonce",
    inputs: [],
    outputs: [{ name: "", type: "uint256" }],
    stateMutability: "view",
  },
  {
    type: "function",
    name: "privileges",
    inputs: [{ name: "", type: "address" }],
    outputs: [{ name: "", type: "bytes32" }],
    stateMutability: "view",
  },
] as const;
```

- [ ] **Step 3: Commit**

```bash
git add packages/cross-chain-dex-ui/src/types/index.ts packages/cross-chain-dex-ui/src/lib/contracts.ts
git commit -m "feat(dex-ui): add AccountMode type and AmbireAccount ABI"
```

---

### Task 2: Create ambireOp.ts library

**Files:**

- Create: `packages/cross-chain-dex-ui/src/lib/ambireOp.ts`

This file handles all AmbireAccount-specific logic: 7702 detection, nonce reading, execute hash computation, signature formatting, and calldata encoding.

- [ ] **Step 1: Create ambireOp.ts**

Create `packages/cross-chain-dex-ui/src/lib/ambireOp.ts`:

```typescript
import {
  type Address,
  type Hex,
  type PublicClient,
  encodeAbiParameters,
  keccak256,
  concat,
  toHex,
  encodeFunctionData,
} from "viem";
import { AmbireAccountABI } from "./contracts";
import { UserOp } from "../types";

// EIP-7702 delegation designator prefix
const DELEGATION_PREFIX = "0xef0100" as const;

/**
 * AmbireAccount Transaction struct matching the on-chain struct.
 */
export interface AmbireTransaction {
  to: Address;
  value: bigint;
  data: Hex;
}

/**
 * Check if an address has an EIP-7702 delegation.
 * Returns the delegation target address if found, null otherwise.
 */
export async function detect7702Delegation(
  client: PublicClient,
  address: Address,
): Promise<Address | null> {
  try {
    const code = await client.getCode({ address });
    if (!code || code === "0x") return null;
    // EIP-7702 delegation designator: 0xef0100 + 20-byte address
    if (
      code.toLowerCase().startsWith(DELEGATION_PREFIX) &&
      code.length === 46
    ) {
      // Extract 20-byte address (bytes 3..23, hex chars 8..48 after 0x prefix)
      return `0x${code.slice(8)}` as Address;
    }
    return null;
  } catch {
    return null;
  }
}

/**
 * Check if a delegation target is an AmbireAccount by looking for the
 * execute function selector in its bytecode.
 */
export async function isAmbireAccount(
  client: PublicClient,
  delegationTarget: Address,
): Promise<boolean> {
  try {
    const code = await client.getCode({ address: delegationTarget });
    if (!code || code === "0x") return false;
    // Check for execute((address,uint256,bytes)[],bytes) selector: 0x6171d1c9
    return code.toLowerCase().includes("6171d1c9");
  } catch {
    return false;
  }
}

/**
 * Read the current nonce from an AmbireAccount.
 */
export async function getAmbireNonce(
  client: PublicClient,
  account: Address,
): Promise<bigint> {
  return (await client.readContract({
    address: account,
    abi: AmbireAccountABI,
    functionName: "nonce",
  })) as bigint;
}

/**
 * Compute the execute hash that AmbireAccount.execute() computes on-chain:
 * keccak256(abi.encode(address(this), block.chainid, currentNonce, txns))
 */
export function computeExecuteHash(
  account: Address,
  chainId: number,
  nonce: bigint,
  txns: AmbireTransaction[],
): Hex {
  // Encode the Transaction[] as a tuple array
  const txnTuples = txns.map((t) => ({
    to: t.to,
    value: t.value,
    data: t.data,
  }));

  return keccak256(
    encodeAbiParameters(
      [
        { type: "address" },
        { type: "uint256" },
        { type: "uint256" },
        {
          type: "tuple[]",
          components: [
            { name: "to", type: "address" },
            { name: "value", type: "uint256" },
            { name: "data", type: "bytes" },
          ],
        },
      ],
      [account, BigInt(chainId), nonce, txnTuples],
    ),
  );
}

/**
 * Append EthSign mode byte (0x01) to a signature for SignatureValidatorV2.
 * Input: 65-byte ECDSA signature (r + s + v) as hex.
 * Output: 66-byte signature with mode byte appended.
 */
export function appendEthSignMode(signature: Hex): Hex {
  return concat([signature, toHex(1, { size: 1 })]) as Hex;
}

/**
 * Convert UserOp[] to AmbireTransaction[].
 */
export function userOpsToAmbireTransactions(
  ops: UserOp[],
): AmbireTransaction[] {
  return ops.map((op) => ({
    to: op.target,
    value: op.value,
    data: op.data,
  }));
}

/**
 * Encode the full AmbireAccount.execute(txns, signature) calldata.
 */
export function buildAmbireExecuteCalldata(
  txns: AmbireTransaction[],
  signature: Hex,
): Hex {
  return encodeFunctionData({
    abi: AmbireAccountABI,
    functionName: "execute",
    args: [txns, signature],
  });
}
```

- [ ] **Step 2: Commit**

```bash
git add packages/cross-chain-dex-ui/src/lib/ambireOp.ts
git commit -m "feat(dex-ui): add ambireOp library for AmbireAccount hash, signing, and calldata"
```

---

### Task 3: Create AccountModeSelector modal component

**Files:**

- Create: `packages/cross-chain-dex-ui/src/components/AccountModeSelector.tsx`

- [ ] **Step 1: Create the component**

Create `packages/cross-chain-dex-ui/src/components/AccountModeSelector.tsx`:

```typescript
import { AccountMode } from '../types';

interface AccountModeSelectorProps {
  isOpen: boolean;
  onSelect: (mode: AccountMode) => void;
  onClose: () => void;
}

export function AccountModeSelector({ isOpen, onSelect, onClose }: AccountModeSelectorProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/75 flex items-center justify-center z-50">
      <div className="bg-surge-card border border-surge-border/50 rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl hover-glow">
        <h2 className="text-xl font-bold text-white mb-2">Choose Account Type</h2>
        <p className="text-gray-400 text-sm mb-6">
          Your wallet supports Ambire Smart Account (EIP-7702). Choose how you'd like to interact with Surge.
        </p>

        <div className="space-y-3">
          <button
            onClick={() => onSelect('safe')}
            className="w-full text-left p-4 bg-surge-dark rounded-xl border border-surge-border/30 hover:border-surge-primary/50 transition-colors group"
          >
            <div className="flex items-center gap-3 mb-1">
              <div className="w-8 h-8 bg-blue-500/20 rounded-lg flex items-center justify-center">
                <svg className="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <span className="text-white font-medium group-hover:text-surge-primary transition-colors">Safe Wallet</span>
              <span className="ml-auto text-xs text-gray-500 bg-surge-card px-2 py-0.5 rounded">Default</span>
            </div>
            <p className="text-xs text-gray-500 ml-11">
              Creates a dedicated Safe. Works with any wallet.
            </p>
          </button>

          <button
            onClick={() => onSelect('ambire')}
            className="w-full text-left p-4 bg-surge-dark rounded-xl border border-surge-border/30 hover:border-surge-secondary/50 transition-colors group"
          >
            <div className="flex items-center gap-3 mb-1">
              <div className="w-8 h-8 bg-purple-500/20 rounded-lg flex items-center justify-center">
                <svg className="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <span className="text-white font-medium group-hover:text-surge-secondary transition-colors">Ambire Account</span>
            </div>
            <p className="text-xs text-gray-500 ml-11">
              Uses your existing 7702 smart account. No extra wallet needed.
            </p>
          </button>
        </div>

        <button
          onClick={onClose}
          className="w-full mt-4 py-2 text-gray-400 hover:text-white text-sm transition-colors"
        >
          Cancel
        </button>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add packages/cross-chain-dex-ui/src/components/AccountModeSelector.tsx
git commit -m "feat(dex-ui): add AccountModeSelector modal for 7702 wallet detection"
```

---

### Task 4: Update useSmartWallet hook to support account modes

**Files:**

- Modify: `packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts`

This is the core change. The hook gains `accountMode` state, 7702 detection, and mode-aware behavior.

- [ ] **Step 1: Add imports and localStorage helpers for account mode**

In `packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts`, add to the existing imports:

```typescript
import { AccountMode } from "../types";
import { detect7702Delegation, isAmbireAccount } from "../lib/ambireOp";
```

Add new constants and helpers after the existing `STORAGE_KEY`:

```typescript
const MODE_STORAGE_KEY = "surge_account_mode_";

function getSavedMode(owner: string): AccountMode | null {
  try {
    const saved = localStorage.getItem(MODE_STORAGE_KEY + owner.toLowerCase());
    return saved === "ambire" || saved === "safe" ? saved : null;
  } catch {
    return null;
  }
}

function saveMode(owner: string, mode: AccountMode): void {
  try {
    localStorage.setItem(MODE_STORAGE_KEY + owner.toLowerCase(), mode);
  } catch {}
}
```

- [ ] **Step 2: Add new state variables to the hook**

Inside `useSmartWallet()`, add after existing state declarations:

```typescript
const [accountMode, setAccountMode] = useState<AccountMode>("safe");
const [has7702Delegation, setHas7702Delegation] = useState(false);
const [showModeSelector, setShowModeSelector] = useState(false);
```

- [ ] **Step 3: Add 7702 detection effect**

Add a new `useEffect` after the existing connect/detect effect. This runs when the wallet connects to check for 7702 delegation:

```typescript
// Detect 7702 delegation and determine account mode
useEffect(() => {
  if (!isConnected || !ownerAddress) {
    setHas7702Delegation(false);
    setShowModeSelector(false);
    return;
  }

  let cancelled = false;

  const detectMode = async () => {
    // Check saved preference first
    const savedMode = getSavedMode(ownerAddress);
    if (savedMode) {
      setAccountMode(savedMode);
      if (savedMode === "ambire") {
        setHas7702Delegation(true);
        // In ambire mode, the smart wallet IS the EOA
        setSmartWallet(ownerAddress);
        setL2WalletExists(true);
        setIsInitializing(false);
      }
      return;
    }

    // No saved preference — check for 7702 delegation
    const delegationTarget = await detect7702Delegation(
      l1PublicClient,
      ownerAddress,
    );
    if (cancelled) return;

    if (delegationTarget) {
      const isAmbire = await isAmbireAccount(l1PublicClient, delegationTarget);
      if (cancelled) return;

      if (isAmbire) {
        setHas7702Delegation(true);
        setShowModeSelector(true);
        return;
      }
    }

    // No 7702 or not AmbireAccount — default to safe mode
    setHas7702Delegation(false);
    setAccountMode("safe");
  };

  detectMode();
  return () => {
    cancelled = true;
  };
}, [isConnected, ownerAddress]);
```

- [ ] **Step 4: Add mode selection handler**

Add a callback after the existing `createL2Wallet` callback:

```typescript
const selectAccountMode = useCallback(
  (mode: AccountMode) => {
    if (!ownerAddress) return;
    saveMode(ownerAddress, mode);
    setAccountMode(mode);
    setShowModeSelector(false);

    if (mode === "ambire") {
      // In ambire mode, the smart wallet IS the EOA
      setSmartWallet(ownerAddress);
      setL2WalletExists(true);
      setIsInitializing(false);
    }
    // If 'safe', the existing Safe detection effect will handle it
  },
  [ownerAddress],
);
```

- [ ] **Step 5: Update the existing Safe detection effect to skip in ambire mode**

Wrap the body of the existing `detectWallet` effect (the one starting with `if (!isConnected || !ownerAddress)`) so it skips when in ambire mode. Add this check right after the `setIsInitializing(true)` line inside `detectWallet`:

```typescript
// Skip Safe detection in ambire mode
if (accountMode === "ambire") {
  setIsInitializing(false);
  return;
}
```

- [ ] **Step 6: Update the return object**

Update the return statement to include the new state:

```typescript
return {
  smartWallet,
  isLoading: isInitializing,
  isCreating: isCreating || isConfirming,
  createSmartWallet,
  ownerAddress,
  isConnected,
  refetch: () => {},
  l2WalletExists,
  createL2Wallet,
  isCreatingL2Wallet,
  accountMode,
  has7702Delegation,
  showModeSelector,
  selectAccountMode,
  setShowModeSelector,
};
```

- [ ] **Step 7: Commit**

```bash
git add packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts
git commit -m "feat(dex-ui): add 7702 detection and account mode selection to useSmartWallet"
```

---

### Task 5: Update useUserOp hook to branch on account mode

**Files:**

- Modify: `packages/cross-chain-dex-ui/src/hooks/useUserOp.ts`

- [ ] **Step 1: Add imports**

Add to the imports in `useUserOp.ts`:

```typescript
import { AccountMode } from "../types";
import {
  getAmbireNonce,
  computeExecuteHash,
  appendEthSignMode,
  userOpsToAmbireTransactions,
  buildAmbireExecuteCalldata,
} from "../lib/ambireOp";
```

- [ ] **Step 2: Add accountMode parameter**

Update the `useUserOp` function signature to accept `accountMode`:

```typescript
export function useUserOp(accountMode: AccountMode = 'safe'): UseUserOpReturn {
```

- [ ] **Step 3: Add Ambire execution path in executeGenericOps**

In the `executeGenericOps` callback, add the Ambire branch. Replace the existing try block body (after `setTxStatus({ phase: 'signing' })`) with mode-aware logic.

**Important:** In Ambire mode on L2 (`chainId === L2_CHAIN_ID`), there is no 7702 delegation, so L2 operations (bridge-out) are sent as direct EOA transactions instead of going through `AmbireAccount.execute()`.

The full updated `executeGenericOps` becomes:

```typescript
const executeGenericOps = useCallback(
  async (
    ops: UserOp[],
    smartWallet: Address,
    chainId?: number,
  ): Promise<boolean> => {
    if (!walletClient) {
      setTxStatus({ phase: "rejected", errorMessage: "Wallet not connected" });
      return false;
    }

    setIsPending(true);
    setError(null);
    txHashRef.current = undefined;

    try {
      setTxStatus({ phase: "signing" });

      const targetChainId = chainId ?? CHAIN_ID;
      const publicClient =
        targetChainId === L2_CHAIN_ID ? l2PublicClient : l1PublicClient;

      if (chainId !== undefined && chainId !== walletClient.chain?.id) {
        await switchChainAsync({ chainId });
      }

      let calldata: Hex;

      if (accountMode === "ambire" && targetChainId === L2_CHAIN_ID) {
        // Ambire mode on L2: no 7702 delegation, send as direct EOA transaction
        // L2 operations are single ops (e.g., bridge sendMessage)
        const op = ops[0];
        await walletClient.sendTransaction({
          to: op.target,
          value: op.value,
          data: op.data,
          chain: walletClient.chain,
          account: walletClient.account,
        });
        setTxStatus({ phase: "complete" });
        setIsPending(false);
        return true;
      } else if (accountMode === "ambire") {
        // AmbireAccount path on L1: personal_sign + execute()
        const txns = userOpsToAmbireTransactions(ops);
        const nonce = await getAmbireNonce(publicClient, smartWallet);
        const executeHash = computeExecuteHash(
          smartWallet,
          targetChainId,
          nonce,
          txns,
        );

        // personal_sign: wallet signs keccak256("\x19Ethereum Signed Message:\n32" + hash)
        const rawSignature = await walletClient.signMessage({
          message: { raw: executeHash },
        });

        const signature = appendEthSignMode(rawSignature as Hex);
        calldata = buildAmbireExecuteCalldata(txns, signature);
      } else {
        // Safe path: signTypedData + execTransaction()
        const nonce = await getSafeNonce(publicClient, smartWallet);
        const safeTx = userOpsToSafeTx(ops);
        const typedData = buildSafeTxTypedData(
          smartWallet,
          targetChainId,
          nonce,
          safeTx,
        );
        const signature = await walletClient.signTypedData(typedData);
        calldata = buildExecTransactionCalldata(safeTx, signature as Hex);
      }

      const result = await sendUserOpToBuilder(smartWallet, calldata, chainId);

      if (result.success && result.userOpId !== undefined) {
        return await pollStatus(result.userOpId);
      } else if (result.success) {
        setTxStatus({ phase: "complete" });
        setIsPending(false);
        return true;
      } else {
        setTxStatus({
          phase: "rejected",
          errorMessage: result.error || "Failed to submit",
        });
        setError(new Error(result.error || "Failed to submit"));
        setIsPending(false);
        return false;
      }
    } catch (err) {
      console.error("Operation failed:", err);
      const msg = err instanceof Error ? err.message : "Operation failed";
      setTxStatus({ phase: "rejected", errorMessage: msg });
      setError(err instanceof Error ? err : new Error(msg));
      setIsPending(false);
      return false;
    }
  },
  [walletClient, switchChainAsync, pollStatus, setTxStatus, accountMode],
);
```

- [ ] **Step 4: Update executeWithdraw for ambire mode**

The `executeWithdraw` function also builds Safe calldata directly. Update it to branch on mode:

```typescript
const executeWithdraw = useCallback(
  async ({
    owner,
    smartWallet,
    ethBalance,
    usdcBalance,
  }: {
    owner: Address;
    smartWallet: Address;
    ethBalance: bigint;
    usdcBalance: bigint;
  }): Promise<boolean> => {
    if (!walletClient) return false;

    const ops = buildWithdrawUserOps(owner, ethBalance, usdcBalance);
    if (ops.length === 0) return false;

    setIsPending(true);
    setError(null);

    try {
      let calldata: Hex;

      if (accountMode === "ambire") {
        const txns = userOpsToAmbireTransactions(ops);
        const nonce = await getAmbireNonce(l1PublicClient, smartWallet);
        const executeHash = computeExecuteHash(
          smartWallet,
          CHAIN_ID,
          nonce,
          txns,
        );
        const rawSignature = await walletClient.signMessage({
          message: { raw: executeHash },
        });
        const signature = appendEthSignMode(rawSignature as Hex);
        calldata = buildAmbireExecuteCalldata(txns, signature);
      } else {
        const nonce = await getSafeNonce(l1PublicClient, smartWallet);
        const safeTx = userOpsToSafeTx(ops);
        const typedData = buildSafeTxTypedData(
          smartWallet,
          CHAIN_ID,
          nonce,
          safeTx,
        );
        const signature = await walletClient.signTypedData(typedData);
        calldata = buildExecTransactionCalldata(safeTx, signature as Hex);
      }

      // In ambire mode, withdraw also goes through builder (no direct tx send)
      // In safe mode, withdraw is sent as a direct transaction
      if (accountMode === "ambire") {
        const result = await sendUserOpToBuilder(smartWallet, calldata);
        setIsPending(false);
        return result.success;
      } else {
        await walletClient.sendTransaction({
          to: smartWallet,
          data: calldata,
          chain: walletClient.chain,
          account: walletClient.account,
        });
        setIsPending(false);
        return true;
      }
    } catch (err) {
      console.error("Withdraw failed:", err);
      setError(err instanceof Error ? err : new Error("Withdraw failed"));
      setIsPending(false);
      return false;
    }
  },
  [walletClient, accountMode],
);
```

- [ ] **Step 5: Commit**

```bash
git add packages/cross-chain-dex-ui/src/hooks/useUserOp.ts
git commit -m "feat(dex-ui): add AmbireAccount execution path to useUserOp"
```

---

### Task 6: Wire up App.tsx with mode selector and pass accountMode through

**Files:**

- Modify: `packages/cross-chain-dex-ui/src/App.tsx`
- Modify: `packages/cross-chain-dex-ui/src/components/FundWallet.tsx`
- Modify: `packages/cross-chain-dex-ui/src/components/SmartWalletSetup.tsx`

- [ ] **Step 1: Update App.tsx imports**

Add the new import at the top of `App.tsx`:

```typescript
import { AccountModeSelector } from "./components/AccountModeSelector";
```

- [ ] **Step 2: Destructure new state from useSmartWallet**

Update the destructuring in `AppContent`:

```typescript
const {
  smartWallet,
  isConnected,
  isLoading,
  ownerAddress,
  createSmartWallet,
  isCreating,
  l2WalletExists,
  createL2Wallet,
  isCreatingL2Wallet,
  accountMode,
  showModeSelector,
  selectAccountMode,
  setShowModeSelector,
} = useSmartWallet();
```

- [ ] **Step 3: Update SmartWalletSetup visibility**

The SmartWalletSetup modal should only show in safe mode. Update the auto-show effect:

```typescript
// Auto-show wallet setup if connected, on correct network, but no smart wallet (safe mode only)
useEffect(() => {
  if (
    isConnected &&
    !isWrongNetwork &&
    !smartWallet &&
    !isLoading &&
    !dismissedWalletSetup &&
    accountMode === "safe"
  ) {
    setShowWalletSetup(true);
  } else if (smartWallet && showWalletSetup) {
    setShowWalletSetup(false);
  }
}, [
  isConnected,
  isWrongNetwork,
  smartWallet,
  isLoading,
  showWalletSetup,
  dismissedWalletSetup,
  accountMode,
]);
```

- [ ] **Step 4: Update FundWallet visibility and props**

In the auto-show effect for FundWallet, skip the L2 wallet check in ambire mode:

```typescript
useEffect(() => {
  if (
    !smartWallet ||
    hasShownFundModal ||
    balancesLoading ||
    isLoading ||
    showNetworkSetup ||
    showWalletSetup
  )
    return;
  const needsFunding = ethBalance === 0n && usdcBalance === 0n;
  const needsL2 = accountMode === "safe" && !l2WalletExists;
  if (needsFunding || needsL2) {
    setShowFundWallet(true);
    setHasShownFundModal(true);
  }
}, [
  smartWallet,
  ethBalance,
  usdcBalance,
  balancesLoading,
  hasShownFundModal,
  isLoading,
  l2WalletExists,
  showNetworkSetup,
  showWalletSetup,
  accountMode,
]);
```

Update FundWallet rendering to hide L2 creation in ambire mode:

```typescript
{smartWallet && (
  <FundWallet
    isOpen={showFundWallet}
    onClose={() => setShowFundWallet(false)}
    smartWallet={smartWallet}
    ethBalance={ethFormatted}
    usdcBalance={usdcFormatted}
    l2WalletExists={accountMode === 'ambire' ? true : l2WalletExists}
    onCreateL2Wallet={accountMode === 'ambire' ? undefined : createL2Wallet}
    isCreatingL2Wallet={isCreatingL2Wallet}
  />
)}
```

- [ ] **Step 5: Add AccountModeSelector to the render tree**

Add after the `SmartWalletSetup` component in the JSX:

```typescript
<AccountModeSelector
  isOpen={showModeSelector}
  onSelect={selectAccountMode}
  onClose={() => setShowModeSelector(false)}
/>
```

- [ ] **Step 6: Commit**

```bash
git add packages/cross-chain-dex-ui/src/App.tsx packages/cross-chain-dex-ui/src/components/FundWallet.tsx
git commit -m "feat(dex-ui): wire AccountModeSelector into App and adapt modals for account mode"
```

---

### Task 7: Pass accountMode to useUserOp consumers

**Files:**

- Modify: `packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts` (the useUserOp call inside it)
- Modify: `packages/cross-chain-dex-ui/src/components/SwapCard.tsx`
- Modify: `packages/cross-chain-dex-ui/src/components/BridgeCard.tsx`
- Modify: `packages/cross-chain-dex-ui/src/components/LiquidityCard.tsx`

Each card component independently calls `useSmartWallet()` and `useUserOp()`. Since `useSmartWallet` now returns `accountMode` (read from localStorage), each card can destructure it and pass it to `useUserOp`.

- [ ] **Step 1: Update useSmartWallet's internal useUserOp call**

In `packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts`, update:

```typescript
const { executeCreateL2Wallet } = useUserOp(accountMode);
```

- [ ] **Step 2: Update SwapCard.tsx**

In `packages/cross-chain-dex-ui/src/components/SwapCard.tsx`, update the `useSmartWallet` destructuring and `useUserOp` call:

```typescript
// Change this line:
const { smartWallet, isConnected } = useSmartWallet();
// To:
const { smartWallet, isConnected, accountMode } = useSmartWallet();

// Change this line:
const { executeSwap, isPending } = useUserOp();
// To:
const { executeSwap, isPending } = useUserOp(accountMode);
```

- [ ] **Step 3: Update BridgeCard.tsx**

In `packages/cross-chain-dex-ui/src/components/BridgeCard.tsx`, apply the same pattern:

```typescript
// Change this line:
const { smartWallet, isConnected, l2WalletExists } = useSmartWallet();
// To:
const { smartWallet, isConnected, l2WalletExists, accountMode } =
  useSmartWallet();

// Change this line:
const {
  executeBridge,
  executeBridgeNative,
  executeBridgeOutNative,
  isPending,
} = useUserOp();
// To:
const {
  executeBridge,
  executeBridgeNative,
  executeBridgeOutNative,
  isPending,
} = useUserOp(accountMode);
```

- [ ] **Step 4: Update LiquidityCard.tsx**

In `packages/cross-chain-dex-ui/src/components/LiquidityCard.tsx`, apply the same pattern:

```typescript
// Change this line:
const { smartWallet, isConnected } = useSmartWallet();
// To:
const { smartWallet, isConnected, accountMode } = useSmartWallet();

// Change this line:
const { executeAddLiquidity, executeRemoveLiquidity, isPending } = useUserOp();
// To:
const { executeAddLiquidity, executeRemoveLiquidity, isPending } =
  useUserOp(accountMode);
```

- [ ] **Step 5: Commit**

```bash
git add packages/cross-chain-dex-ui/src/hooks/useSmartWallet.ts packages/cross-chain-dex-ui/src/components/SwapCard.tsx packages/cross-chain-dex-ui/src/components/BridgeCard.tsx packages/cross-chain-dex-ui/src/components/LiquidityCard.tsx
git commit -m "feat(dex-ui): pass accountMode to all useUserOp consumers"
```

---

### Task 8: Manual integration test

**Files:** None (testing only)

- [ ] **Step 1: Start the dev server**

```bash
cd packages/cross-chain-dex-ui && pnpm dev
```

- [ ] **Step 2: Test Safe mode (default behavior)**

1. Connect with a regular EOA (MetaMask or similar)
2. Verify NO `AccountModeSelector` modal appears
3. Verify SmartWalletSetup modal appears as before
4. Verify swap/bridge/liquidity flows work unchanged

- [ ] **Step 3: Test 7702 detection**

1. Connect with an EOA that has 7702 delegation to AmbireAccount
2. Verify `AccountModeSelector` modal appears
3. Select "Safe Wallet" — verify normal Safe flow proceeds
4. Disconnect, reconnect — verify saved preference is used (no modal)

- [ ] **Step 4: Test Ambire Account mode**

1. Connect with a 7702-delegated EOA
2. Select "Ambire Account"
3. Verify NO SmartWalletSetup modal appears
4. Verify FundWallet modal shows but without L2 wallet creation step
5. Verify balances are read from the EOA address
6. Attempt a swap — verify `personal_sign` prompt appears (not `signTypedData`)
7. Verify the signed transaction is submitted to Catalyst

- [ ] **Step 5: Test mode persistence**

1. Select "Ambire Account" mode
2. Refresh the page
3. Verify the mode is restored from localStorage (no modal)
4. Clear localStorage, reconnect — verify modal reappears
