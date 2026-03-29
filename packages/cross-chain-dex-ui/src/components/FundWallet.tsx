import { Address } from 'viem';
import toast from 'react-hot-toast';
import { surgeL1Chain } from '../lib/config';
import { L1_NATIVE_SYMBOL } from '../lib/constants';

interface FundWalletProps {
  isOpen: boolean;
  onClose: () => void;
  onCreateL2Wallet?: () => Promise<boolean>;
  l2WalletExists?: boolean | null;
  smartWallet: Address;
  ethBalance: string;
  usdcBalance: string;
}

export function FundWallet({ isOpen, onClose, onCreateL2Wallet, l2WalletExists, smartWallet, ethBalance, usdcBalance }: FundWalletProps) {
  if (!isOpen) return null;

  const copyAddress = () => {
    navigator.clipboard.writeText(smartWallet);
    toast.success('Address copied!');
  };

  return (
    <div className="fixed inset-0 bg-black/75 flex items-center justify-center z-50">
      <div className="bg-surge-card border border-surge-border/50 rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl hover-glow">
        <div className="flex items-center gap-3 mb-4">
          <div className="w-10 h-10 bg-yellow-500/20 rounded-full flex items-center justify-center">
            <svg className="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h2 className="text-xl font-bold text-white">Fund Your Smart Wallet</h2>
            <p className="text-sm text-gray-400">Add {L1_NATIVE_SYMBOL} or USDC to start swapping</p>
          </div>
        </div>

        <div className="bg-red-500/10 border border-red-500/30 rounded-lg px-3 py-2 text-xs text-red-400 flex items-start gap-1.5 mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5 shrink-0 mt-0.5" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
          <span>Experimental Alpha - transaction limit of US $1. <a href="https://surge.wtf/alpha-disclaimer" target="_blank" rel="noopener noreferrer" className="underline hover:text-red-300">See disclaimer</a></span>
        </div>

        <p className="text-gray-400 text-sm mb-6">
          Your smart wallet needs funds to execute swaps. Send {L1_NATIVE_SYMBOL} or USDC to the address below.
        </p>

        {/* Smart Wallet Address */}
        <div className="bg-surge-dark rounded-lg p-4 mb-6">
          <div className="text-xs text-gray-500 mb-2">Smart Wallet Address</div>
          <div className="flex items-center gap-2">
            <code className="text-sm text-white font-mono flex-1 break-all">
              {smartWallet}
            </code>
            <button
              onClick={copyAddress}
              className="p-2 hover:bg-surge-border rounded transition-colors shrink-0"
            >
              <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>

        {/* Current Balances */}
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="bg-surge-dark rounded-lg p-3">
            <div className="text-xs text-gray-500 mb-1">{L1_NATIVE_SYMBOL} Balance</div>
            <div className="text-lg font-semibold text-white">
              {parseFloat(ethBalance).toFixed(4)} {L1_NATIVE_SYMBOL}
            </div>
          </div>
          <div className="bg-surge-dark rounded-lg p-3">
            <div className="text-xs text-gray-500 mb-1">USDC Balance</div>
            <div className="text-lg font-semibold text-white">
              {parseFloat(usdcBalance).toFixed(2)} USDC
            </div>
          </div>
        </div>

        <div className="text-xs text-gray-500 mb-4">
          <strong>Note:</strong> Send funds on {surgeL1Chain.name} (Chain ID: {surgeL1Chain.id})
        </div>

        <button
          onClick={async () => {
            if (l2WalletExists === false && onCreateL2Wallet) {
              await onCreateL2Wallet();
            }
            onClose();
          }}
          className="w-full py-3 bg-surge-primary hover:bg-surge-secondary text-white rounded-lg font-medium transition-colors"
        >
          {l2WalletExists === false ? "I've Funded My Wallet — Create L2 Wallet" : "I've Funded My Wallet"}
        </button>
      </div>
    </div>
  );
}
