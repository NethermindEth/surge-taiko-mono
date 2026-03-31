import { useState, useEffect, useCallback, useRef } from 'react';
import { type Address, type Hex, keccak256, encodePacked, encodeAbiParameters, getContractAddress, decodeEventLog, concat, toHex } from 'viem';
import { useAccount, useWriteContract, useWaitForTransactionReceipt } from 'wagmi';
import toast from 'react-hot-toast';
import { SafeProxyFactoryABI, SafeProxyFactoryFullABI } from '../lib/contracts';
import { SAFE_PROXY_FACTORY, SAFE_SINGLETON, SAFE_FALLBACK_HANDLER } from '../lib/constants';
import { buildSafeSetupCalldata } from '../lib/safeOp';
import { l1PublicClient, l2PublicClient } from '../lib/config';
import { useUserOp } from './useUserOp';
import { AccountMode } from '../types';
import { detect7702Delegation, isAmbireAccount } from '../lib/ambireOp';

const STORAGE_KEY = 'surge_safe_address_';

let cachedProxyCreationCode: Hex | undefined;

/**
 * Predict the CREATE2 address of a Safe proxy.
 * SafeProxyFactory salt = keccak256(keccak256(initializer) ++ uint256(saltNonce))
 * deploymentData = proxyCreationCode ++ abi.encode(singleton)
 */
async function predictSafeAddress(owner: Address): Promise<Address> {
  const initializer = buildSafeSetupCalldata(owner, SAFE_FALLBACK_HANDLER);
  const saltNonce = BigInt(keccak256(encodePacked(['address'], [owner])));

  if (!cachedProxyCreationCode) {
    cachedProxyCreationCode = await l1PublicClient.readContract({
      address: SAFE_PROXY_FACTORY,
      abi: SafeProxyFactoryFullABI,
      functionName: 'proxyCreationCode',
    });
  }

  const salt = keccak256(
    concat([keccak256(initializer), toHex(saltNonce, { size: 32 })])
  );

  const deploymentData = concat([
    cachedProxyCreationCode,
    encodeAbiParameters([{ type: 'address' }], [SAFE_SINGLETON]),
  ]);

  return getContractAddress({
    from: SAFE_PROXY_FACTORY,
    salt,
    bytecode: deploymentData,
    opcode: 'CREATE2',
  });
}

function getSavedSafe(owner: string): Address | null {
  try {
    const saved = localStorage.getItem(STORAGE_KEY + owner.toLowerCase());
    return saved as Address | null;
  } catch {
    return null;
  }
}

function saveSafe(owner: string, safe: Address): void {
  try {
    localStorage.setItem(STORAGE_KEY + owner.toLowerCase(), safe);
  } catch {}
}

const MODE_STORAGE_KEY = 'surge_account_mode_';

function getSavedMode(owner: string): AccountMode | null {
  try {
    const saved = localStorage.getItem(MODE_STORAGE_KEY + owner.toLowerCase());
    return saved === 'ambire' || saved === 'safe' ? saved : null;
  } catch {
    return null;
  }
}

function saveMode(owner: string, mode: AccountMode): void {
  try {
    localStorage.setItem(MODE_STORAGE_KEY + owner.toLowerCase(), mode);
  } catch {}
}

export function useSmartWallet() {
  const { address: ownerAddress, isConnected } = useAccount();
  const [isInitializing, setIsInitializing] = useState(true);
  const [smartWallet, setSmartWallet] = useState<Address | null>(null);
  const [l2WalletExists, setL2WalletExists] = useState(false);
  const [isCreatingL2Wallet, setIsCreatingL2Wallet] = useState(false);
  const [accountMode, setAccountMode] = useState<AccountMode>('safe');
  const [has7702Delegation, setHas7702Delegation] = useState(false);
  const [showModeSelector, setShowModeSelector] = useState(false);

  const { writeContract, data: txHash, isPending: isCreating, reset } = useWriteContract();
  const { data: receipt, isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash: txHash,
  });

  const { executeCreateL2Wallet } = useUserOp(accountMode);
  const l2PollRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Cleanup L2 poll on unmount
  useEffect(() => {
    return () => {
      if (l2PollRef.current) clearInterval(l2PollRef.current);
    };
  }, []);

  // On connect: check localStorage or predict the deterministic Safe address and verify on-chain.
  useEffect(() => {
    if (!isConnected || !ownerAddress) {
      setSmartWallet(null);
      setIsInitializing(false);
      return;
    }

    let cancelled = false;
    setIsInitializing(true);

    const detectWallet = async () => {
      // Skip Safe detection in ambire mode
      if (accountMode === 'ambire') {
        setIsInitializing(false);
        return;
      }

      // First try localStorage
      const saved = getSavedSafe(ownerAddress);
      if (saved) {
        try {
          const code = await l1PublicClient.getCode({ address: saved });
          if (cancelled) return;
          if (code && code !== '0x') {
            setSmartWallet(saved);
            setIsInitializing(false);
            return;
          }
          // Stale entry — remove it
          localStorage.removeItem(STORAGE_KEY + ownerAddress.toLowerCase());
        } catch (err) {
          if (cancelled) return;
          console.warn('Failed to verify saved Safe address:', err);
        }
      }

      // If not in localStorage (or no code), predict the CREATE2 address and check
      try {
        const predicted = await predictSafeAddress(ownerAddress);
        if (cancelled) return;
        const code = await l1PublicClient.getCode({ address: predicted });
        if (cancelled) return;
        if (code && code !== '0x') {
          setSmartWallet(predicted);
          saveSafe(ownerAddress, predicted);
        }
      } catch (err) {
        if (cancelled) return;
        console.warn('Failed to predict Safe address:', err);
      }

      if (!cancelled) setIsInitializing(false);
    };

    detectWallet();
    return () => { cancelled = true; };
  }, [isConnected, ownerAddress, accountMode]);

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
        if (savedMode === 'ambire') {
          setHas7702Delegation(true);
          setSmartWallet(ownerAddress);
          setL2WalletExists(true);
          setIsInitializing(false);
        }
        return;
      }

      // No saved preference — check for 7702 delegation
      const delegationTarget = await detect7702Delegation(l1PublicClient, ownerAddress);
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
      setAccountMode('safe');
    };

    detectMode();
    return () => { cancelled = true; };
  }, [isConnected, ownerAddress]);

  // After a successful creation tx, parse the ProxyCreation event to get the proxy address.
  useEffect(() => {
    if (!isSuccess || !receipt || !ownerAddress) return;

    for (const log of receipt.logs) {
      try {
        const decoded = decodeEventLog({
          abi: SafeProxyFactoryABI,
          data: log.data,
          topics: log.topics,
        });

        if (decoded.eventName === 'ProxyCreation') {
          const proxyAddress = (decoded.args as { proxy: Address }).proxy;
          console.log('Safe created:', proxyAddress);
          toast.dismiss('create-wallet');
          toast.success(
            `Safe wallet created: ${proxyAddress.slice(0, 8)}...${proxyAddress.slice(-6)}`,
          );
          setSmartWallet(proxyAddress);
          saveSafe(ownerAddress, proxyAddress);
          reset();
          break;
        }
      } catch {
        // Not a ProxyCreation log — skip.
      }
    }
  }, [isSuccess, receipt, ownerAddress, reset]);

  // After L1 Safe is known, check if the same address has code on L2
  useEffect(() => {
    if (!smartWallet) {
      setL2WalletExists(false);
      return;
    }

    let cancelled = false;
    l2PublicClient
      .getCode({ address: smartWallet })
      .then((code) => {
        if (!cancelled) setL2WalletExists(!!(code && code !== '0x'));
      })
      .catch(() => { if (!cancelled) setL2WalletExists(false); });
    return () => { cancelled = true; };
  }, [smartWallet]);

  const createSmartWallet = useCallback(async () => {
    if (!ownerAddress) throw new Error('Wallet not connected');

    const initializer = buildSafeSetupCalldata(ownerAddress, SAFE_FALLBACK_HANDLER);
    const saltNonce = BigInt(keccak256(encodePacked(['address'], [ownerAddress])));

    writeContract({
      address: SAFE_PROXY_FACTORY,
      abi: SafeProxyFactoryABI,
      functionName: 'createProxyWithNonce',
      args: [SAFE_SINGLETON, initializer, saltNonce],
    });
  }, [ownerAddress, writeContract]);

  const createL2Wallet = useCallback(async (): Promise<void> => {
    if (!ownerAddress || !smartWallet) {
      toast.error('Smart wallet not ready');
      return;
    }

    setIsCreatingL2Wallet(true);
    try {
      const success = await executeCreateL2Wallet({ owner: ownerAddress, smartWallet });
      if (success) {
        toast.success('L2 Safe creation submitted via bridge');
        // Poll L2 for wallet code every 5s until found (max 60 attempts = 5 min)
        let attempts = 0;
        if (l2PollRef.current) clearInterval(l2PollRef.current);
        l2PollRef.current = setInterval(async () => {
          attempts++;
          try {
            const code = await l2PublicClient.getCode({ address: smartWallet });
            if (code && code !== '0x') {
              setL2WalletExists(true);
              if (l2PollRef.current) clearInterval(l2PollRef.current);
              l2PollRef.current = null;
            }
          } catch {}
          if (attempts >= 60 && l2PollRef.current) {
            clearInterval(l2PollRef.current);
            l2PollRef.current = null;
          }
        }, 5000);
      } else {
        toast.error('Failed to create L2 Safe');
      }
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to create L2 Safe';
      toast.error(msg);
    } finally {
      setIsCreatingL2Wallet(false);
    }
  }, [ownerAddress, smartWallet, executeCreateL2Wallet]);

  const selectAccountMode = useCallback((mode: AccountMode) => {
    if (!ownerAddress) return;
    saveMode(ownerAddress, mode);
    setAccountMode(mode);
    setShowModeSelector(false);

    if (mode === 'ambire') {
      setSmartWallet(ownerAddress);
      setL2WalletExists(true);
      setIsInitializing(false);
    }
  }, [ownerAddress]);

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
}
