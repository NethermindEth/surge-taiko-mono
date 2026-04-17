import { useState, useCallback, useMemo, useEffect } from 'react';
import { parseUnits, formatUnits } from 'viem';
import { useAccount } from 'wagmi';
import { TokenInput } from './TokenInput';
import { SwapDetails } from './SwapDetails';
import { SwapPath } from './SwapPath';
import { SwapButton } from './SwapButton';
import { useSmartWallet } from '../context/SmartWalletContext';
import { useDexReserves } from '../hooks/useDexReserves';
import { useSwapQuote } from '../hooks/useSwapQuote';
import { useL1DexQuote } from '../hooks/useL1DexQuote';
import { useSharedTokenBalances } from '../context/SmartWalletContext';
import { useUserOp } from '../hooks/useUserOp';
import { useInitiateL2Swap } from '../hooks/useInitiateL2Swap';
import { SwapDirection, SwapVenue } from '../types';
import { ETH_TOKEN, USDC_TOKEN, DEFAULT_SLIPPAGE } from '../lib/constants';
import { calculateMinOutput } from '../lib/userOp';
import { DisclaimerModal } from './DisclaimerModal';
import { useDisclaimer } from '../hooks/useDisclaimer';
import { WarningBanner } from './WarningBanner';

const MAX_SWAP_AMOUNT = 1; // $1 max per swap (L2_DEX venue only; audience demo guard-rail)

interface SwapCardProps {
  onSetupWallet: () => void;
  onFundWallet?: () => void;
  venue: SwapVenue;
  onVenueChange: (v: SwapVenue) => void;
}

export function SwapCard({ onSetupWallet, onFundWallet: _onFundWallet, venue, onVenueChange: _onVenueChange }: SwapCardProps) {
  const { smartWallet, isConnected, accountMode } = useSmartWallet();
  const { address: eoa } = useAccount();
  const { ethReserve, tokenReserve } = useDexReserves();
  // useSharedTokenBalances reads from the context — already network-aware
  // (L1: smart wallet balances, L2: EOA balances)
  const { ethBalance, usdcBalance } = useSharedTokenBalances();
  const { executeSwap, isPending: isUserOpPending } = useUserOp(accountMode);
  const { initiate: initiateL2Swap, isPending: isL2SwapPending } = useInitiateL2Swap();
  const { isDisclaimerOpen, requireDisclaimer, onAccept, onCancel } = useDisclaimer();

  const [direction, setDirection] = useState<SwapDirection>('ETH_TO_USDC');
  const [inputAmount, setInputAmount] = useState('');

  const isL1Venue = venue === 'L1_DEX';
  const isPending = isL1Venue ? isL2SwapPending : isUserOpPending;

  const inputToken = direction === 'ETH_TO_USDC' ? ETH_TOKEN : USDC_TOKEN;
  const outputToken = direction === 'ETH_TO_USDC' ? USDC_TOKEN : ETH_TOKEN;

  const amountIn = useMemo(() => {
    try {
      return inputAmount ? parseUnits(inputAmount, inputToken.decimals) : 0n;
    } catch {
      return 0n;
    }
  }, [inputAmount, inputToken.decimals]);

  // Quote source is venue-specific. Both hooks are cheap memo/read-only, so we call both
  // but only surface the one for the active venue.
  const l2Quote = useSwapQuote({
    direction,
    amountIn: isL1Venue ? 0n : amountIn,
    ethReserve,
    tokenReserve,
  });
  const l1Quote = useL1DexQuote({
    direction,
    amountIn: isL1Venue ? amountIn : 0n,
  });
  const quote = isL1Venue ? l1Quote : l2Quote;

  // Balances: L2_DEX uses the smart wallet on L1, L1_DEX uses the EOA on L2.
  // Context balances are already network-aware (L1=smart wallet, L2=EOA)
  const inputBalance = direction === 'ETH_TO_USDC' ? ethBalance : usdcBalance;
  const outputBalance = direction === 'ETH_TO_USDC' ? usdcBalance : ethBalance;

  const hasInsufficientBalance = amountIn > inputBalance;
  const exceedsSwapLimit =
    !isL1Venue && amountIn > parseUnits(String(MAX_SWAP_AMOUNT), inputToken.decimals)
      ? `Max ${MAX_SWAP_AMOUNT} ${inputToken.symbol} per swap`
      : undefined;

  // Venue is now controlled by the network selector in App — clear input on change
  useEffect(() => {
    setInputAmount('');
  }, [venue]);

  const handleSwapDirection = useCallback(() => {
    setDirection((prev) => (prev === 'ETH_TO_USDC' ? 'USDC_TO_ETH' : 'ETH_TO_USDC'));
    setInputAmount('');
  }, []);

  const handleSwap = useCallback(async () => {
    if (amountIn === 0n) return;

    if (isL1Venue) {
      if (!eoa) return;
      const minOut = calculateMinOutput(quote.amountOut, DEFAULT_SLIPPAGE);
      const ok = await initiateL2Swap({
        direction,
        amountIn,
        minAmountOut: minOut,
        recipient: eoa,
        expectedAmountOut: quote.amountOut,
      });
      if (ok) setInputAmount('');
      return;
    }

    if (!smartWallet) return;
    const ok = await executeSwap({
      direction,
      amountIn,
      expectedAmountOut: quote.amountOut,
      smartWallet,
    });
    if (ok) setInputAmount('');
  }, [
    isL1Venue,
    eoa,
    smartWallet,
    amountIn,
    direction,
    quote.amountOut,
    executeSwap,
    initiateL2Swap,
  ]);

  // L1 venue bypasses smart-wallet setup — the connected EOA is sufficient.
  const walletGateOk = isL1Venue ? !!eoa : !!smartWallet;
  const needsSetupWallet = isConnected && !isL1Venue && !smartWallet;

  return (
    <div className="flex flex-col md:flex-row items-start gap-4 justify-center w-full relative z-10">
      {/* Left panel — inputs */}
      <div className="w-full md:max-w-md bg-surge-card/80 border border-surge-border/50 rounded-2xl p-4 space-y-3 shadow-xl shadow-black/20 hover-glow transition-all duration-[1000ms] ease-[cubic-bezier(0.16,1,0.3,1)]">
        {/* Header */}
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-white">Swap</h2>
          <span className={`text-xs font-medium ${isL1Venue ? 'text-cyan-400' : 'text-emerald-400'}`}>
            {isL1Venue ? 'Via EOA' : 'Via Smart Account'}
          </span>
        </div>

        <WarningBanner />

        {/* Input Token */}
        <TokenInput
          token={inputToken}
          amount={inputAmount}
          onAmountChange={setInputAmount}
          balance={inputBalance}
          label="From"
        />

        {/* Swap Direction Button */}
        <div className="flex justify-center -my-2 relative z-10">
          <button
            onClick={handleSwapDirection}
            className="p-2 bg-surge-card border border-surge-border rounded-lg hover:bg-surge-dark transition-colors"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5 text-gray-400"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fillRule="evenodd"
                d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                clipRule="evenodd"
              />
              <path
                fillRule="evenodd"
                d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        </div>

        {/* Output Token */}
        <TokenInput
          token={outputToken}
          amount={quote.amountOut > 0n ? formatUnits(quote.amountOut, outputToken.decimals) : ''}
          onAmountChange={() => {}}
          balance={outputBalance}
          label="To"
          disabled
          showMax={false}
        />

        {/* Swap Button */}
        <SwapButton
          onClick={
            needsSetupWallet ? onSetupWallet : () => requireDisclaimer(handleSwap)
          }
          disabled={false}
          isLoading={isPending}
          isConnected={isConnected}
          hasSmartWallet={walletGateOk}
          hasInsufficientBalance={hasInsufficientBalance}
          hasInsufficientLiquidity={quote.insufficientLiquidity}
          hasAmount={amountIn > 0n}
          exceedsSwapLimit={exceedsSwapLimit}
          needsApproval={isL1Venue && direction === 'USDC_TO_ETH'}
        />

      </div>

      {/* Right panel — trade details (always visible; shows "-" placeholders when no input) */}
      <div className="w-full md:max-w-sm bg-surge-card/80 border border-surge-border/50 rounded-2xl p-4 space-y-3 shadow-xl shadow-black/20">
        <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-widest">Trade Details</h3>
        <SwapPath direction={direction} venue={venue} />
        <SwapDetails quote={quote} direction={direction} amountIn={amountIn} />
      </div>
      <DisclaimerModal isOpen={isDisclaimerOpen} onAccept={onAccept} onCancel={onCancel} />
    </div>
  );
}

