import { useState, useEffect } from 'react';
import { WagmiProvider } from 'wagmi';
import { useAccount, useSwitchChain } from 'wagmi';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RainbowKitProvider, darkTheme } from '@rainbow-me/rainbowkit';
import '@rainbow-me/rainbowkit/styles.css';
import { Toaster } from 'react-hot-toast';

import { config, surgeL1Chain, surgeL2Chain } from './lib/config';
// L1_CHAIN_NAME used by Header via props
import { Header } from './components/Header';
import { SwapCard } from './components/SwapCard';
import { BridgeCard } from './components/BridgeCard';
import { LiquidityCard } from './components/LiquidityCard';
import { SmartWalletSetup } from './components/SmartWalletSetup';
import { NetworkSetup } from './components/NetworkSetup';
import { FundWallet } from './components/FundWallet';
import { TxStatusOverlay } from './components/TxStatusOverlay';
import { TxStatusProvider, useTxStatus } from './context/TxStatusContext';
import { useSmartWallet, SmartWalletProvider } from './context/SmartWalletContext';
import { useSharedTokenBalances } from './context/SmartWalletContext';
import { AccountModeSelector } from './components/AccountModeSelector';
import { ActiveTab, SwapVenue } from './types';

const queryClient = new QueryClient();

function AppContent() {
  const { txStatus, setTxStatus } = useTxStatus();
  const {
  smartWallet, isConnected, isLoading, ownerAddress,
  createSmartWallet, isCreating,
  l2WalletExists, createL2Wallet, isCreatingL2Wallet,
  accountMode, has7702Delegation, showModeSelector, selectAccountMode, setShowModeSelector,
} = useSmartWallet();
  const { chainId } = useAccount();
  const { switchChainAsync } = useSwitchChain();
  const { ethBalance, usdcBalance, ethFormatted, usdcFormatted, isLoading: balancesLoading } = useSharedTokenBalances();

  const { selectedNetwork } = useSmartWallet();
  const [activeTab, setActiveTab] = useState<ActiveTab>('swap');
  const [showWalletSetup, setShowWalletSetup] = useState(false);
  const [dismissedWalletSetup, setDismissedWalletSetup] = useState(false);
  const [showNetworkSetup, setShowNetworkSetup] = useState(false);
  const [showFundWallet, setShowFundWallet] = useState(false);
  const [hasShownFundModal, setHasShownFundModal] = useState(false);

  // Derive venue and required chain from selected network
  const venue: SwapVenue = selectedNetwork === 'l2' ? 'L1_DEX' : 'L2_DEX';
  const requiredChainId = selectedNetwork === 'l2' ? surgeL2Chain.id : surgeL1Chain.id;
  const networkSetupTarget = selectedNetwork === 'l2' ? 'l2' as const : 'l1' as const;

  // When L2 is selected, force swap tab (no bridge/liquidity on L2)
  useEffect(() => {
    if (selectedNetwork === 'l2' && activeTab !== 'swap') {
      setActiveTab('swap');
    }
  }, [selectedNetwork, activeTab]);

  // Accept both L1 and L2 as valid networks
  const isWrongNetwork = isConnected && chainId !== surgeL1Chain.id && chainId !== surgeL2Chain.id;

  // Only show network setup if on a completely unknown chain (neither L1 nor L2).
  // If the user is on a valid chain but not the one matching the selected network,
  // silently auto-switch without showing the modal.
  useEffect(() => {
    if (!isConnected || !chainId) return;
    if (isWrongNetwork) {
      // On an unknown chain — try to switch, show modal if that fails
      switchChainAsync({ chainId: requiredChainId }).catch(() => {
        setShowNetworkSetup(true);
      });
    } else if (chainId !== requiredChainId) {
      // On a valid chain (L1 or L2) but not the one for the selected network — silent switch
      switchChainAsync({ chainId: requiredChainId }).catch(() => {});
      setShowNetworkSetup(false);
    } else {
      setShowNetworkSetup(false);
    }
  }, [isConnected, chainId, requiredChainId, isWrongNetwork, switchChainAsync]);

  // Reset dismissed flag when wallet connects/disconnects
  useEffect(() => {
    setDismissedWalletSetup(false);
  }, [isConnected, ownerAddress]);

  // Auto-show wallet setup if connected, on correct network, but no smart wallet
  useEffect(() => {
    if (isConnected && !isWrongNetwork && !smartWallet && !isLoading && !dismissedWalletSetup && accountMode === 'safe' && !showModeSelector) {
      setShowWalletSetup(true);
    } else if (smartWallet && showWalletSetup) {
      setShowWalletSetup(false);
    }
  }, [isConnected, isWrongNetwork, smartWallet, isLoading, showWalletSetup, dismissedWalletSetup, accountMode, showModeSelector]);

  // Auto-show fund wallet modal
  useEffect(() => {
    if (accountMode === 'ambire') return;
    if (!smartWallet || hasShownFundModal || balancesLoading || isLoading || showNetworkSetup || showWalletSetup) return;
    const needsFunding = ethBalance === 0n && usdcBalance === 0n;
    const needsL2 = accountMode === 'safe' && !l2WalletExists;
    if (needsFunding || needsL2) {
      setShowFundWallet(true);
      setHasShownFundModal(true);
    }
  }, [smartWallet, ethBalance, usdcBalance, balancesLoading, hasShownFundModal, isLoading, l2WalletExists, showNetworkSetup, showWalletSetup, accountMode]);

  const availableTabs: ActiveTab[] = selectedNetwork === 'l2'
    ? ['swap']
    : ['swap', 'liquidity', 'bridge'];

  return (
    <div className="h-screen overflow-hidden bg-surge-dark flex flex-col">
      <Header onSetupWallet={() => has7702Delegation ? setShowModeSelector(true) : setShowWalletSetup(true)} />

      <main className="flex-1 min-h-0 relative flex items-center justify-center px-4">
        {/* Tab Navigation */}
        <div className="absolute top-8 left-1/2 -translate-x-1/2 flex gap-1 bg-surge-card/50 rounded-xl p-1 border border-surge-border/30 z-10">
          {availableTabs.map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-5 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${
                activeTab === tab
                  ? 'bg-gradient-to-r from-surge-primary to-surge-secondary text-white shadow-md'
                  : 'text-gray-400 hover:text-white hover:bg-surge-dark/50'
              }`}
            >
              {tab === 'swap' ? 'Swap' : tab === 'bridge' ? 'Bridge' : 'Liquidity'}
            </button>
          ))}
        </div>

        {/* Footer tagline */}
        <div className="absolute bottom-5 left-1/2 -translate-x-1/2 text-center whitespace-nowrap">
          <p className="text-sm text-gray-400">Powered by Surge Protocol</p>
          <p className="text-sm text-gray-500 mt-1">
            {selectedNetwork === 'l2'
              ? 'Synchronous cross-chain settlement'
              : 'L1 swaps through L2 liquidity \u2022 Real time cross chain settlement'}
          </p>
        </div>

        {/* Active Panel */}
        <div className="w-full flex items-center justify-center">
          {activeTab === 'swap' && (
            <SwapCard
              onSetupWallet={() => setShowWalletSetup(true)}
              onFundWallet={() => setShowFundWallet(true)}
              venue={venue}
              onVenueChange={() => {}}
            />
          )}
          {activeTab === 'bridge' && selectedNetwork === 'l1' && (
            <BridgeCard
              onSetupWallet={() => setShowWalletSetup(true)}
              onFundWallet={() => setShowFundWallet(true)}
            />
          )}
          {activeTab === 'liquidity' && selectedNetwork === 'l1' && (
            <LiquidityCard
              onSetupWallet={() => setShowWalletSetup(true)}
            />
          )}
        </div>
      </main>

      <NetworkSetup
        isOpen={showNetworkSetup}
        onClose={() => setShowNetworkSetup(false)}
        targetChain={networkSetupTarget}
      />

      <AccountModeSelector
        isOpen={showModeSelector}
        onSelect={selectAccountMode}
        onClose={() => setShowModeSelector(false)}
      />

      <SmartWalletSetup
        isOpen={showWalletSetup && !isWrongNetwork && !showModeSelector}
        onClose={() => {
          setShowWalletSetup(false);
          setDismissedWalletSetup(true);
        }}
        ownerAddress={ownerAddress}
        isCreating={isCreating}
        createSmartWallet={createSmartWallet}
      />

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

      <TxStatusOverlay
        state={txStatus}
        onClose={() => setTxStatus({ phase: 'idle' })}
      />

      <Toaster
        position="bottom-right"
        toastOptions={{
          style: {
            background: '#0f2847',
            color: '#e2e8f0',
            border: '1px solid #1e4976',
          },
          success: {
            iconTheme: {
              primary: '#10b981',
              secondary: '#fff',
            },
          },
          error: {
            iconTheme: {
              primary: '#ef4444',
              secondary: '#fff',
            },
          },
        }}
      />
    </div>
  );
}

function App() {
  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider theme={darkTheme({ accentColor: '#10b981' })}>
          <TxStatusProvider>
            <SmartWalletProvider>
              <AppContent />
            </SmartWalletProvider>
          </TxStatusProvider>
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}

export default App;
