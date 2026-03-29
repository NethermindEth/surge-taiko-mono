import { useState, useCallback, useMemo } from 'react';
import { parseUnits, formatUnits } from 'viem';
import { TokenInput } from './TokenInput';
import { SwapDetails } from './SwapDetails';
import { SwapPath } from './SwapPath';
import { SwapButton } from './SwapButton';
import { useSmartWallet } from '../hooks/useSmartWallet';
import { useDexReserves } from '../hooks/useDexReserves';
import { useSwapQuote } from '../hooks/useSwapQuote';
import { useTokenBalances } from '../hooks/useTokenBalances';
import { useUserOp } from '../hooks/useUserOp';
import { SwapDirection } from '../types';
import { ETH_TOKEN, USDC_TOKEN } from '../lib/constants';
import { DisclaimerModal } from './DisclaimerModal';
import { useDisclaimer } from '../hooks/useDisclaimer';

const MAX_SWAP_AMOUNT = 1; // $1 max per swap

interface SwapCardProps {
  onSetupWallet: () => void;
  onFundWallet?: () => void;
}

export function SwapCard({ onSetupWallet, onFundWallet }: SwapCardProps) {
  const { smartWallet, isConnected } = useSmartWallet();
  const { ethReserve, tokenReserve, isLoading: reservesLoading } = useDexReserves();
  const { ethBalance, usdcBalance } = useTokenBalances(smartWallet);
  const { executeSwap, isPending } = useUserOp();
  const { isDisclaimerOpen, requireDisclaimer, onAccept, onCancel } = useDisclaimer();

  const [direction, setDirection] = useState<SwapDirection>('ETH_TO_USDC');
  const [inputAmount, setInputAmount] = useState('');

  const inputToken = direction === 'ETH_TO_USDC' ? ETH_TOKEN : USDC_TOKEN;
  const outputToken = direction === 'ETH_TO_USDC' ? USDC_TOKEN : ETH_TOKEN;

  const amountIn = useMemo(() => {
    try {
      return inputAmount ? parseUnits(inputAmount, inputToken.decimals) : 0n;
    } catch {
      return 0n;
    }
  }, [inputAmount, inputToken.decimals]);

  const quote = useSwapQuote({
    direction,
    amountIn,
    ethReserve,
    tokenReserve,
  });
  const inputBalance = direction === 'ETH_TO_USDC' ? ethBalance : usdcBalance;
  const outputBalance = direction === 'ETH_TO_USDC' ? usdcBalance : ethBalance;

  const hasInsufficientBalance = amountIn > inputBalance;
  const exceedsSwapLimit = amountIn > parseUnits(String(MAX_SWAP_AMOUNT), inputToken.decimals)
    ? `Max ${MAX_SWAP_AMOUNT} ${inputToken.symbol} per swap`
    : undefined;

  const handleSwapDirection = useCallback(() => {
    setDirection((prev) => (prev === 'ETH_TO_USDC' ? 'USDC_TO_ETH' : 'ETH_TO_USDC'));
    setInputAmount('');
  }, []);

  const handleSwap = useCallback(async () => {
    if (!smartWallet || amountIn === 0n) return;

    const success = await executeSwap({
      direction,
      amountIn,
      expectedAmountOut: quote.amountOut,
      smartWallet,
    });

    if (success) {
      setInputAmount('');
    }
  }, [smartWallet, amountIn, direction, quote.amountOut, executeSwap]);

  return (
    <div className="flex flex-col md:flex-row items-start gap-4 justify-center w-full relative z-10">
      {/* Left panel — inputs */}
      <div className="w-full md:max-w-md bg-surge-card/80 border border-surge-border/50 rounded-2xl p-4 space-y-3 shadow-xl shadow-black/20 hover-glow transition-all duration-[1000ms] ease-[cubic-bezier(0.16,1,0.3,1)]">
        {/* Header */}
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-white">Swap</h2>
          {reservesLoading && (
            <span className="text-xs text-gray-400">Loading reserves...</span>
          )}
        </div>
        <div className="bg-red-500/10 border border-red-500/30 rounded-lg px-3 py-2 text-xs text-red-400 flex items-center gap-1.5">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5 shrink-0" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
          Experimental Alpha - transaction limit of US $1. <a href="https://surge.wtf/alpha-disclaimer" target="_blank" rel="noopener noreferrer" className="underline hover:text-red-300">See disclaimer</a>
        </div>

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
          onClick={isConnected && !smartWallet ? onSetupWallet : () => requireDisclaimer(handleSwap)}
          disabled={false}
          isLoading={isPending}
          isConnected={isConnected}
          hasSmartWallet={!!smartWallet}
          hasInsufficientBalance={hasInsufficientBalance}
          hasInsufficientLiquidity={quote.insufficientLiquidity}
          hasAmount={amountIn > 0n}
          exceedsSwapLimit={exceedsSwapLimit}
        />

      </div>

      {/* Right panel — trade details (shown when amount is entered) */}
      {amountIn > 0n && (
        <div className="w-full md:max-w-sm bg-surge-card/80 border border-surge-border/50 rounded-2xl p-4 space-y-3 shadow-xl shadow-black/20 animate-panel-in">
          <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-widest">Trade Details</h3>
          <SwapDetails quote={quote} direction={direction} amountIn={amountIn} />
          <SwapPath direction={direction} show={true} />
        </div>
      )}
      <DisclaimerModal isOpen={isDisclaimerOpen} onAccept={onAccept} onCancel={onCancel} />
    </div>
  );
}
