import { createContext, useContext, ReactNode } from 'react';
import { useSmartWalletInternal, SmartWalletState } from '../hooks/useSmartWallet';

const SmartWalletContext = createContext<SmartWalletState | null>(null);

export function SmartWalletProvider({ children }: { children: ReactNode }) {
  const value = useSmartWalletInternal();
  return (
    <SmartWalletContext.Provider value={value}>
      {children}
    </SmartWalletContext.Provider>
  );
}

export function useSmartWallet(): SmartWalletState {
  const ctx = useContext(SmartWalletContext);
  if (!ctx) throw new Error('useSmartWallet must be used within SmartWalletProvider');
  return ctx;
}
