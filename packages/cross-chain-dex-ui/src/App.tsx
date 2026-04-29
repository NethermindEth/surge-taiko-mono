import { useState, useEffect } from 'react';
import { WagmiProvider } from 'wagmi';
import { useAccount, useSwitchChain } from 'wagmi';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RainbowKitProvider, lightTheme } from '@rainbow-me/rainbowkit';
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
  //
  // While a tx is in flight (e.g. a bridge withdrawal triggered from the L1 page that
  // needs to sign on L2), suppress the auto-switch so we don't fight the in-flight
  // hook-driven chain switch. Once the tx settles ('complete' / 'rejected' / 'idle'),
  // the effect re-fires and pulls the wallet back to the selected page network.
  const isTxInFlight = txStatus.phase !== 'idle' && txStatus.phase !== 'complete' && txStatus.phase !== 'rejected';
  useEffect(() => {
    if (!isConnected || !chainId) return;
    if (isTxInFlight) return;
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
  }, [isConnected, chainId, requiredChainId, isWrongNetwork, switchChainAsync, isTxInFlight]);

  // Reset dismissed flag when wallet connects/disconnects
  useEffect(() => {
    setDismissedWalletSetup(false);
  }, [isConnected, ownerAddress]);

  // Auto-show wallet setup if connected on L1 with no smart wallet.
  // On L2, swaps use the connected EOA — smart wallet is an L1-only feature.
  useEffect(() => {
    if (isConnected && selectedNetwork === 'l1' && !isWrongNetwork && !smartWallet && !isLoading && !dismissedWalletSetup && accountMode === 'safe' && !showModeSelector) {
      setShowWalletSetup(true);
    } else if (smartWallet && showWalletSetup) {
      setShowWalletSetup(false);
    }
  }, [isConnected, isWrongNetwork, smartWallet, isLoading, showWalletSetup, dismissedWalletSetup, accountMode, showModeSelector, selectedNetwork]);

  // Auto-show fund wallet modal — L1 only (smart wallet is L1-only feature).
  useEffect(() => {
    if (selectedNetwork !== 'l1') return;
    if (accountMode === 'ambire') return;
    if (!smartWallet || hasShownFundModal || balancesLoading || isLoading || showNetworkSetup || showWalletSetup) return;
    const needsFunding = ethBalance === 0n || usdcBalance === 0n;
    const needsL2 = accountMode === 'safe' && !l2WalletExists;
    if (needsFunding || needsL2) {
      setShowFundWallet(true);
      setHasShownFundModal(true);
    }
  }, [smartWallet, ethBalance, usdcBalance, balancesLoading, hasShownFundModal, isLoading, l2WalletExists, showNetworkSetup, showWalletSetup, accountMode, selectedNetwork]);

  const availableTabs: ActiveTab[] = selectedNetwork === 'l2'
    ? ['swap']
    : ['swap', 'liquidity', 'bridge'];

  return (
    <div className="h-screen overflow-hidden bg-surge-dark flex flex-col">
      <Header onSetupWallet={() => has7702Delegation ? setShowModeSelector(true) : setShowWalletSetup(true)} />

      <main className="flex-1 min-h-0 relative flex items-center justify-center px-4">
        {/* Tab Navigation */}
        <div className="absolute top-8 left-1/2 -translate-x-1/2 flex gap-1 bg-white/70 backdrop-blur-md rounded-xl p-1 border border-surge-border z-10 shadow-sm">
          {availableTabs.map((tab) => {
            const active = activeTab === tab;
            const pastelHover =
              tab === 'swap' ? 'hover:bg-surge-mint/30'
              : tab === 'liquidity' ? 'hover:bg-surge-lavender/30'
              : 'hover:bg-surge-secondary/20';
            return (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`px-5 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${
                  active
                    ? 'bg-surge-primary text-white shadow-md'
                    : `text-surge-muted hover:text-surge-primary ${pastelHover}`
                }`}
              >
                {tab === 'swap' ? 'Swap' : tab === 'bridge' ? 'Bridge' : 'Liquidity'}
              </button>
            );
          })}
        </div>

        {/* Footer — Powered by Nethermind */}
        <div className="absolute bottom-5 left-1/2 -translate-x-1/2 text-center">
          <a
            href="https://www.nethermind.io/"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-block opacity-80 hover:opacity-100 transition-opacity"
          >
            <img src="/powered-by-nethermind.svg" alt="Powered by Nethermind" className="h-5 w-auto" />
          </a>
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
        isOpen={showModeSelector && selectedNetwork === 'l1'}
        onSelect={selectAccountMode}
        onClose={() => setShowModeSelector(false)}
      />

      <SmartWalletSetup
        isOpen={showWalletSetup && !isWrongNetwork && !showModeSelector && selectedNetwork === 'l1'}
        onClose={() => {
          setShowWalletSetup(false);
          setDismissedWalletSetup(true);
        }}
        ownerAddress={ownerAddress}
        isCreating={isCreating}
        createSmartWallet={createSmartWallet}
      />

      {smartWallet && selectedNetwork === 'l1' && (
        <FundWallet
          isOpen={showFundWallet}
          onClose={() => setShowFundWallet(false)}
          smartWallet={smartWallet}
          ethBalance={ethFormatted}
          usdcBalance={usdcFormatted}
          targetChainId={requiredChainId}
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
            background: '#ffffff',
            color: '#172342',
            border: '1px solid #e5e7eb',
            boxShadow: '0 8px 24px rgba(23, 35, 66, 0.08)',
          },
          success: {
            iconTheme: {
              primary: '#8ce8ab',
              secondary: '#172342',
            },
          },
          error: {
            iconTheme: {
              primary: '#fabeab',
              secondary: '#172342',
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
        <RainbowKitProvider theme={lightTheme({ accentColor: '#172342', accentColorForeground: '#ffffff', borderRadius: 'medium' })}>
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
