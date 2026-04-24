import { useState, useRef, useEffect, useCallback } from 'react';
import { useAccount, useDisconnect } from 'wagmi';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import toast from 'react-hot-toast';
import { useSmartWallet } from '../context/SmartWalletContext';
import { useSharedTokenBalances } from '../context/SmartWalletContext';
import { useUserOp } from '../hooks/useUserOp';
import { ETH_TOKEN, L1_CHAIN_NAME } from '../lib/constants';
import { DisclaimerModal } from './DisclaimerModal';
import { useDisclaimer } from '../hooks/useDisclaimer';

interface HeaderProps {
  onSetupWallet: () => void;
}

export function Header({ onSetupWallet }: HeaderProps) {
  const { smartWallet, isConnected, ownerAddress, accountMode, clearAccountMode, selectedNetwork, setSelectedNetwork, l2WalletExists } = useSmartWallet();
  const { address: eoaAddress } = useAccount();
  const { disconnect } = useDisconnect();
  const { executeWithdraw, isPending } = useUserOp(accountMode);
  const { isDisclaimerOpen, requireDisclaimer, onAccept, onCancel } = useDisclaimer();
  const { ethBalance, usdcBalance, ethFormatted, usdcFormatted } = useSharedTokenBalances();

  const [swDropdownOpen, setSwDropdownOpen] = useState(false);
  const [eoaDropdownOpen, setEoaDropdownOpen] = useState(false);
  const [networkDropdownOpen, setNetworkDropdownOpen] = useState(false);
  const swDropdownRef = useRef<HTMLDivElement>(null);
  const eoaDropdownRef = useRef<HTMLDivElement>(null);
  const networkDropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (swDropdownRef.current && !swDropdownRef.current.contains(e.target as Node)) {
        setSwDropdownOpen(false);
      }
      if (eoaDropdownRef.current && !eoaDropdownRef.current.contains(e.target as Node)) {
        setEoaDropdownOpen(false);
      }
      if (networkDropdownRef.current && !networkDropdownRef.current.contains(e.target as Node)) {
        setNetworkDropdownOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, []);

  const handleWithdraw = useCallback(async () => {
    if (!smartWallet || !ownerAddress) return;
    if (ethBalance === 0n && usdcBalance === 0n) {
      toast.error('No funds to withdraw');
      return;
    }
    setSwDropdownOpen(false);
    const success = await executeWithdraw({ owner: ownerAddress, smartWallet, ethBalance, usdcBalance });
    if (success) {
      toast.success('Withdrawal submitted');
    }
  }, [smartWallet, ownerAddress, ethBalance, usdcBalance, executeWithdraw]);

  return (
    <>
    <header className="w-full px-6 py-3 flex items-center justify-between border-b border-surge-border bg-white/70 backdrop-blur-md relative z-10">
      <div className="flex items-center gap-4">
        <a
          href="https://www.surge.wtf/"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center"
        >
          <img src="/surge-logo.svg" alt="Surge" className="h-10 w-auto" />
        </a>

        {/* Network dropdown */}
        <div className="relative" ref={networkDropdownRef}>
          <button
            onClick={() => { setNetworkDropdownOpen((p) => !p); setSwDropdownOpen(false); setEoaDropdownOpen(false); }}
            className={`flex items-center gap-2 px-3 py-1.5 rounded-lg border text-sm font-medium transition-colors ${
              selectedNetwork === 'l1'
                ? 'bg-surge-mint/25 border-surge-mint/70 hover:bg-surge-mint/35'
                : 'bg-surge-secondary/20 border-surge-secondary/60 hover:bg-surge-secondary/30'
            }`}
          >
            <div className={`w-2 h-2 rounded-full ${selectedNetwork === 'l1' ? 'bg-surge-mint' : 'bg-surge-secondary'}`} />
            <span className="text-surge-primary">{selectedNetwork === 'l1' ? L1_CHAIN_NAME : 'Surge L2'}</span>
            <svg className="w-3.5 h-3.5 text-surge-primary/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          {networkDropdownOpen && (
            <div className="absolute left-0 top-full mt-1 bg-surge-card border border-surge-border rounded-lg shadow-xl shadow-surge-primary/10 overflow-hidden min-w-[180px] z-50">
              <button
                onClick={() => { setSelectedNetwork('l1'); setNetworkDropdownOpen(false); }}
                className={`w-full px-4 py-2.5 text-left text-sm flex items-center gap-2 transition-colors ${
                  selectedNetwork === 'l1' ? 'text-surge-primary bg-surge-mint/20' : 'text-surge-muted hover:bg-surge-card-hover hover:text-surge-text'
                }`}
              >
                <div className="w-2 h-2 rounded-full bg-surge-mint" />
                {L1_CHAIN_NAME}
              </button>
              <button
                onClick={() => { setSelectedNetwork('l2'); setNetworkDropdownOpen(false); }}
                className={`w-full px-4 py-2.5 text-left text-sm flex items-center gap-2 transition-colors border-t border-surge-border ${
                  selectedNetwork === 'l2' ? 'text-surge-primary bg-surge-secondary/15' : 'text-surge-muted hover:bg-surge-card-hover hover:text-surge-text'
                }`}
              >
                <div className="w-2 h-2 rounded-full bg-surge-secondary" />
                Surge L2
              </button>
            </div>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3">
        {/* Smart Wallet: show "Setup on L2" button when on L2 network and L2 wallet doesn't exist */}
        {isConnected && smartWallet && selectedNetwork === 'l2' && !l2WalletExists && (
          <button
            onClick={onSetupWallet}
            className="px-4 py-2 bg-surge-secondary/15 hover:bg-surge-secondary/25 text-surge-primary rounded-lg text-sm font-medium transition-colors border border-surge-secondary/40"
          >
            Setup Smart Wallet on L2
          </button>
        )}

        {/* Smart Wallet display — show when wallet exists on the selected network */}
        {isConnected && smartWallet && (selectedNetwork === 'l1' || l2WalletExists) && (
          <div className="hidden md:flex items-center relative" ref={swDropdownRef}>
            <div className="flex items-center gap-2 px-3 py-2 bg-surge-card rounded-lg border border-surge-border">
              <div className="w-2 h-2 bg-surge-mint rounded-full" />
              <button
                onClick={() => {
                  navigator.clipboard.writeText(smartWallet);
                  toast.success('Smart wallet address copied!');
                }}
                className="text-sm text-surge-text hover:text-surge-primary transition-colors flex items-center gap-1"
                title="Click to copy"
              >
                Smart Wallet: {smartWallet.slice(0, 6)}...{smartWallet.slice(-4)}
                <svg className="w-3 h-3 text-surge-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              </button>
              <span className="text-xs text-surge-muted/60">|</span>
              <span className="text-xs text-surge-muted">
                {parseFloat(ethFormatted).toFixed(4)} {ETH_TOKEN.symbol}
              </span>
              <span className="text-xs text-surge-muted">
                {parseFloat(usdcFormatted).toFixed(2)} USDC
              </span>
              {accountMode === 'safe' && (
                <>
                  <span className="text-xs text-surge-muted/60">|</span>
                  <button
                    onClick={() => { setSwDropdownOpen((p) => !p); setEoaDropdownOpen(false); }}
                    className="text-surge-muted hover:text-surge-primary transition-colors p-0.5"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                  </button>
                </>
              )}
            </div>

            {swDropdownOpen && accountMode === 'safe' && (
              <div className="absolute right-0 top-full mt-1 bg-surge-card border border-surge-border rounded-lg shadow-xl shadow-surge-primary/10 overflow-hidden min-w-[240px] z-50">
                <button
                  onClick={() => requireDisclaimer(handleWithdraw)}
                  disabled={isPending || (ethBalance === 0n && usdcBalance === 0n)}
                  className="w-full px-4 py-3 text-left text-sm text-surge-text hover:bg-surge-card-hover hover:text-surge-primary transition-colors disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg className="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                  {isPending ? 'Withdrawing...' : 'Withdraw all funds to owner'}
                </button>
              </div>
            )}
          </div>
        )}

        {isConnected && !smartWallet && (
          <button
            onClick={onSetupWallet}
            className="px-4 py-2 bg-surge-primary hover:opacity-90 text-white rounded-lg text-sm font-medium transition-opacity shadow-sm"
          >
            Setup Smart Wallet
          </button>
        )}

        {/* EOA Wallet */}
        <ConnectButton.Custom>
          {({ account, chain, openConnectModal, mounted }) => {
            const connected = mounted && account && chain;
            return (
              <div {...(!mounted && { 'aria-hidden': true, style: { opacity: 0, pointerEvents: 'none' as const, userSelect: 'none' as const } })}>
                {connected ? (
                  <div className="relative" ref={eoaDropdownRef}>
                    <button
                      onClick={() => { setEoaDropdownOpen((p) => !p); setSwDropdownOpen(false); }}
                      className="px-4 py-2 bg-surge-card hover:bg-surge-card-hover text-surge-text rounded-lg text-sm font-medium transition-colors border border-surge-border flex items-center gap-2"
                    >
                      <span>{account.displayName}</span>
                      {account.displayBalance && (
                        <span className="text-xs text-surge-muted">({account.displayBalance})</span>
                      )}
                      <svg className="w-3 h-3 text-surge-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                    {eoaDropdownOpen && (
                      <div className="absolute right-0 top-full mt-1 bg-surge-card border border-surge-border rounded-lg shadow-xl shadow-surge-primary/10 overflow-hidden min-w-[180px] z-50">
                        <button
                          onClick={() => {
                            if (eoaAddress) {
                              navigator.clipboard.writeText(eoaAddress);
                              toast.success('EOA address copied!');
                            }
                            setEoaDropdownOpen(false);
                          }}
                          className="w-full px-4 py-3 text-left text-sm text-surge-text hover:bg-surge-card-hover hover:text-surge-primary transition-colors flex items-center gap-2"
                        >
                          <svg className="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                          </svg>
                          Copy Address
                        </button>
                        <button
                          onClick={() => { setEoaDropdownOpen(false); clearAccountMode(); disconnect(); }}
                          className="w-full px-4 py-3 text-left text-sm text-surge-amber hover:bg-surge-peach/20 hover:text-surge-amber transition-colors flex items-center gap-2 border-t border-surge-border"
                        >
                          <svg className="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                          </svg>
                          Disconnect
                        </button>
                      </div>
                    )}
                  </div>
                ) : (
                  <button
                    onClick={openConnectModal}
                    className="px-4 py-2 bg-surge-primary hover:opacity-90 text-white rounded-lg text-sm font-medium transition-opacity shadow-sm"
                  >
                    Connect Wallet
                  </button>
                )}
              </div>
            );
          }}
        </ConnectButton.Custom>
      </div>
    </header>
    <DisclaimerModal isOpen={isDisclaimerOpen} onAccept={onAccept} onCancel={onCancel} />
    </>
  );
}
