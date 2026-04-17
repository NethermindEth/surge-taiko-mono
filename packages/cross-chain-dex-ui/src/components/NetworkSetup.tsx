import { useState } from 'react';
import { surgeL1Chain, surgeL2Chain } from '../lib/config';
import { L1_RPC_URL, L2_RPC_URL } from '../lib/constants';
import toast from 'react-hot-toast';

type TargetChain = 'l1' | 'l2';

interface NetworkSetupProps {
  isOpen: boolean;
  onClose: () => void;
  targetChain?: TargetChain;
}

export function NetworkSetup({ isOpen, onClose, targetChain = 'l1' }: NetworkSetupProps) {
  const [isAdding, setIsAdding] = useState(false);
  const [showManual, setShowManual] = useState(false);

  if (!isOpen) return null;

  const chain = targetChain === 'l2' ? surgeL2Chain : surgeL1Chain;
  const rpcUrl = targetChain === 'l2' ? L2_RPC_URL : L1_RPC_URL;

  const networkConfig = {
    networkName: chain.name,
    rpcUrl,
    chainId: chain.id,
    currencySymbol: chain.nativeCurrency.symbol,
  };

  const addNetwork = async () => {
    if (!window.ethereum) {
      toast.error('No wallet detected');
      return;
    }

    setIsAdding(true);
    try {
      await window.ethereum.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId: `0x${chain.id.toString(16)}` }],
      });
      toast.success(`Switched to ${chain.name}!`);
      onClose();
    } catch (switchError: any) {
      if (switchError.code === 4902) {
        try {
          await window.ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [{
              chainId: `0x${chain.id.toString(16)}`,
              chainName: chain.name,
              nativeCurrency: chain.nativeCurrency,
              rpcUrls: [rpcUrl],
            }],
          });
          toast.success('Network added! Please switch to it.');
        } catch (addError: any) {
          if (addError.code === -32602 || addError.message?.includes('HTTPS')) {
            setShowManual(true);
            toast.error('Wallet requires HTTPS. Please add network manually.');
          } else {
            toast.error('Failed to add network');
          }
        }
      } else {
        toast.error('Failed to switch network');
      }
    } finally {
      setIsAdding(false);
    }
  };

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    toast.success(`${label} copied!`);
  };

  return (
    <div className="fixed inset-0 bg-black/75 flex items-center justify-center z-50">
      <div className="bg-surge-card border border-surge-border/50 rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl hover-glow">
        <div className="flex items-center gap-3 mb-4">
          <div className="w-10 h-10 bg-red-500/20 rounded-full flex items-center justify-center">
            <svg className="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <div>
            <h2 className="text-xl font-bold text-white">Wrong Network</h2>
            <p className="text-sm text-gray-400">Connect to {chain.name} to continue</p>
          </div>
        </div>

        {!showManual ? (
          <>
            <p className="text-gray-400 text-sm mb-6">
              Click the button below to add {chain.name} network to your wallet and switch to it.
            </p>

            <button
              onClick={addNetwork}
              disabled={isAdding}
              className="w-full py-3 bg-surge-primary hover:bg-surge-secondary disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors flex items-center justify-center gap-2 mb-4"
            >
              {isAdding ? (
                <>
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Adding Network...
                </>
              ) : (
                `Add ${chain.name} Network`
              )}
            </button>

            <button
              onClick={() => setShowManual(true)}
              className="w-full py-2 text-gray-400 hover:text-white text-sm transition-colors"
            >
              Add manually instead
            </button>
          </>
        ) : (
          <>
            <p className="text-gray-400 text-sm mb-4">
              Your wallet requires HTTPS URLs. Please add the network manually:
            </p>

            <div className="space-y-3 mb-6">
              {[
                { label: 'Network Name', value: networkConfig.networkName },
                { label: 'RPC URL', value: networkConfig.rpcUrl, truncate: true },
                { label: 'Chain ID', value: String(networkConfig.chainId) },
                { label: 'Currency Symbol', value: networkConfig.currencySymbol },
              ].map(({ label, value, truncate }) => (
                <div key={label} className="bg-surge-dark rounded-lg p-3">
                  <div className="flex justify-between items-center">
                    <div className={truncate ? 'flex-1 min-w-0' : ''}>
                      <div className="text-xs text-gray-500 mb-1">{label}</div>
                      <div className={`text-sm text-white font-mono ${truncate ? 'truncate' : ''}`}>{value}</div>
                    </div>
                    <button
                      onClick={() => copyToClipboard(value, label)}
                      className="p-2 hover:bg-surge-border rounded transition-colors ml-2"
                    >
                      <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                      </svg>
                    </button>
                  </div>
                </div>
              ))}
            </div>

            <div className="text-xs text-gray-500 mb-4">
              <strong>Steps:</strong> MetaMask → Settings → Networks → Add Network
            </div>

            <button
              onClick={onClose}
              className="w-full py-3 bg-surge-primary hover:bg-surge-secondary text-white rounded-lg font-medium transition-colors"
            >
              I've Added the Network
            </button>
          </>
        )}
      </div>
    </div>
  );
}
