import { useState } from 'react';
import { Address, parseEther, parseUnits } from 'viem';
import { useWalletClient, useSwitchChain } from 'wagmi';
import toast from 'react-hot-toast';
import { surgeL1Chain, surgeL2Chain } from '../lib/config';
import { L1_NATIVE_SYMBOL, USDC_TOKEN, L2_USDC_TOKEN } from '../lib/constants';
import { ERC20ABI } from '../lib/contracts';
import { WarningBannerWrapped } from './WarningBanner';

interface FundWalletProps {
  isOpen: boolean;
  onClose: () => void;
  smartWallet: Address;
  ethBalance: string;
  usdcBalance: string;
  /// Chain the transfers must land on. Driven by the app's selected network; passed
  /// in so the fund CTAs never race against a pending wagmi chain switch.
  targetChainId: number;
  l2WalletExists?: boolean;
  onCreateL2Wallet?: () => Promise<void>;
  isCreatingL2Wallet?: boolean;
}

export function FundWallet({
  isOpen,
  onClose,
  smartWallet,
  ethBalance,
  usdcBalance,
  targetChainId,
  l2WalletExists = false,
  onCreateL2Wallet,
  isCreatingL2Wallet = false,
}: FundWalletProps) {
  const { data: walletClient } = useWalletClient();
  const { switchChainAsync } = useSwitchChain();
  const [isFunding, setIsFunding] = useState(false);
  const [isFundingUsdc, setIsFundingUsdc] = useState(false);

  if (!isOpen) return null;

  const hasFunds = parseFloat(ethBalance) > 0 || parseFloat(usdcBalance) > 0;

  // Target chain drives which USDC token we touch. The CTA label is always "USDC"
  // regardless of the underlying contract's symbol.
  const isOnL2 = targetChainId === surgeL2Chain.id;
  const activeUsdc = isOnL2 ? L2_USDC_TOKEN : USDC_TOKEN;
  const targetChain = isOnL2 ? surgeL2Chain : surgeL1Chain;
  const hasEth = parseFloat(ethBalance) > 0;
  const hasUsdc = parseFloat(usdcBalance) > 0;

  /// Ensure the wallet is on the target chain before submitting. Without this, a
  /// pending wagmi chain-switch can silently divert the transfer to the opposite
  /// chain (e.g. sending L1 USDC when the user is viewing the L2 page).
  const ensureTargetChain = async (): Promise<boolean> => {
    if (walletClient?.chain?.id === targetChainId) return true;
    try {
      await switchChainAsync({ chainId: targetChainId });
      return true;
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Chain switch failed';
      if (!msg.includes('rejected')) toast.error(msg);
      return false;
    }
  };

  const fundWallet = async () => {
    if (!walletClient) return;
    setIsFunding(true);
    try {
      if (!(await ensureTargetChain())) return;
      const hash = await walletClient.sendTransaction({
        chain: targetChain,
        to: smartWallet,
        value: parseEther('0.01'),
      });
      toast.success(`Sent 0.01 ${L1_NATIVE_SYMBOL} (tx: ${hash.slice(0, 10)}...)`);
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Transfer failed';
      if (!msg.includes('rejected')) toast.error(msg);
    } finally {
      setIsFunding(false);
    }
  };

  const fundWalletUsdc = async () => {
    if (!walletClient || !activeUsdc.address) return;
    setIsFundingUsdc(true);
    try {
      if (!(await ensureTargetChain())) return;
      const hash = await walletClient.writeContract({
        chain: targetChain,
        address: activeUsdc.address,
        abi: ERC20ABI,
        functionName: 'transfer',
        args: [smartWallet, parseUnits('1', activeUsdc.decimals)],
      });
      toast.success(`Sent 1 USDC (tx: ${hash.slice(0, 10)}...)`);
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Transfer failed';
      if (!msg.includes('rejected')) toast.error(msg);
    } finally {
      setIsFundingUsdc(false);
    }
  };

  const copyAddress = () => {
    navigator.clipboard.writeText(smartWallet);
    toast.success('Address copied!');
  };

  return (
    <div className="fixed inset-0 bg-black/75 flex items-center justify-center z-50">
      <div className="bg-surge-card border border-surge-border/50 rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl hover-glow relative">
        <button
          onClick={onClose}
          aria-label="Close"
          className="absolute top-4 right-4 p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-surge-border/50 transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <div className="flex items-center gap-3 mb-4 pr-8">
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

        <WarningBannerWrapped />

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

        <div className="flex gap-3 mb-4">
          <button
            onClick={fundWallet}
            disabled={isFunding || !walletClient || hasEth}
            className="flex-1 py-3 bg-surge-primary/80 hover:bg-surge-primary disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
          >
            {isFunding ? 'Sending...' : `Send 0.01 ${L1_NATIVE_SYMBOL}`}
          </button>
          <button
            onClick={fundWalletUsdc}
            disabled={isFundingUsdc || !walletClient || !activeUsdc.address || hasUsdc}
            className="flex-1 py-3 bg-surge-primary/80 hover:bg-surge-primary disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
          >
            {isFundingUsdc ? 'Sending...' : 'Send 1 USDC'}
          </button>
        </div>

        <div className="text-xs text-gray-500 mb-4">
          <strong>Note:</strong> Send funds on {targetChain.name} (Chain ID: {targetChain.id})
        </div>

        {/* L2 Safe status / creation — only show after wallet is funded */}
        {hasFunds && !l2WalletExists && onCreateL2Wallet && (
          <div className="mb-4">
            <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg px-3 py-2 text-xs text-yellow-400 mb-3">
              Your Safe wallet does not yet exist on L2. Create it via the bridge to enable L2 DEX operations.
            </div>
            <button
              onClick={onCreateL2Wallet}
              disabled={isCreatingL2Wallet}
              className="w-full py-3 bg-yellow-600 hover:bg-yellow-500 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
            >
              {isCreatingL2Wallet ? 'Creating L2 Wallet...' : 'Create L2 Wallet'}
            </button>
          </div>
        )}

        {l2WalletExists && (
          <div className="bg-green-500/10 border border-green-500/30 rounded-lg px-3 py-2 text-xs text-green-400 mb-4">
            L2 Safe wallet is active at the same address.
          </div>
        )}

        {hasFunds && (
          <button
            onClick={onClose}
            className={`w-full py-3 rounded-lg font-medium transition-colors ${
              l2WalletExists || !onCreateL2Wallet
                ? 'bg-surge-primary hover:bg-surge-secondary text-white'
                : 'bg-surge-card border border-surge-border/50 text-gray-400 hover:text-white hover:border-surge-border'
            }`}
          >
            {l2WalletExists || !onCreateL2Wallet
              ? 'Done'
              : 'Skip L2 wallet setup for now'}
          </button>
        )}
      </div>
    </div>
  );
}
