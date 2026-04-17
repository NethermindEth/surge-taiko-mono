import { createContext, useContext, useState, ReactNode, useMemo, useEffect } from 'react';
import { useAccount } from 'wagmi';
import { useSmartWalletInternal, SmartWalletState } from '../hooks/useSmartWallet';
import { useTokenBalances as useTokenBalancesHook } from '../hooks/useTokenBalances';

type SelectedNetwork = 'l1' | 'l2';

const SELECTED_NETWORK_STORAGE_KEY = 'surge-dex:selected-network';

function loadSelectedNetwork(): SelectedNetwork {
  if (typeof window === 'undefined') return 'l1';
  const raw = window.localStorage.getItem(SELECTED_NETWORK_STORAGE_KEY);
  return raw === 'l2' ? 'l2' : 'l1';
}

interface SmartWalletContextValue extends SmartWalletState {
  tokenBalances: ReturnType<typeof useTokenBalancesHook>;
  selectedNetwork: SelectedNetwork;
  setSelectedNetwork: (n: SelectedNetwork) => void;
}

const SmartWalletContext = createContext<SmartWalletContextValue | null>(null);

export function SmartWalletProvider({ children }: { children: ReactNode }) {
  const wallet = useSmartWalletInternal();
  const { address: eoaAddress } = useAccount();
  const [selectedNetwork, setSelectedNetwork] = useState<SelectedNetwork>(loadSelectedNetwork);

  useEffect(() => {
    window.localStorage.setItem(SELECTED_NETWORK_STORAGE_KEY, selectedNetwork);
  }, [selectedNetwork]);
  // L1: show Smart Wallet balances. L2: show EOA balances (L2 swaps use EOA directly).
  const balanceAddress = selectedNetwork === 'l2' ? (eoaAddress ?? null) : wallet.smartWallet;
  const tokenBalances = useTokenBalancesHook(balanceAddress, selectedNetwork);
  const value = useMemo(
    () => ({ ...wallet, tokenBalances, selectedNetwork, setSelectedNetwork }),
    [wallet, tokenBalances, selectedNetwork]
  );
  return (
    <SmartWalletContext.Provider value={value}>
      {children}
    </SmartWalletContext.Provider>
  );
}

export function useSmartWallet(): SmartWalletContextValue {
  const ctx = useContext(SmartWalletContext);
  if (!ctx) throw new Error('useSmartWallet must be used within SmartWalletProvider');
  return ctx;
}

export function useSharedTokenBalances() {
  const { tokenBalances } = useSmartWallet();
  return tokenBalances;
}
