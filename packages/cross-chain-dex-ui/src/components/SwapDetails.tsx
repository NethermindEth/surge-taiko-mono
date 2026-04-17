import { SwapQuote, SwapDirection } from '../types';
import { formatUnits } from 'viem';
import { FEE_PERCENT, ETH_TOKEN, USDC_TOKEN } from '../lib/constants';

interface SwapDetailsProps {
  quote: SwapQuote;
  direction: SwapDirection;
  amountIn: bigint;
}

export function SwapDetails({ quote, direction, amountIn }: SwapDetailsProps) {
  const inputToken = direction === 'ETH_TO_USDC' ? ETH_TOKEN : USDC_TOKEN;
  const outputToken = direction === 'ETH_TO_USDC' ? USDC_TOKEN : ETH_TOKEN;
  const inputSymbol = inputToken.symbol;
  const outputSymbol = outputToken.symbol;

  const hasInput = amountIn > 0n;

  return (
    <div className="bg-surge-dark rounded-xl p-4 space-y-3">
      <div className="flex justify-between items-center text-sm">
        <span className="text-gray-400">Rate</span>
        <span className="text-white">
          {hasInput ? `1 ${inputSymbol} = ${quote.rate.toFixed(6)} ${outputSymbol}` : '-'}
        </span>
      </div>

      <div className="flex justify-between items-center text-sm">
        <span className="text-gray-400">Fee ({FEE_PERCENT}%)</span>
        <span className="text-white">
          {hasInput
            ? `${Number(formatUnits(quote.fee, inputToken.decimals)).toFixed(6)} ${inputSymbol}`
            : '-'}
        </span>
      </div>

      <div className="flex justify-between items-center text-sm">
        <span className="text-gray-400">Price Impact</span>
        <span
          className={
            hasInput
              ? quote.priceImpact > 5
                ? 'text-red-500'
                : quote.priceImpact > 1
                ? 'text-yellow-500'
                : 'text-green-500'
              : 'text-white'
          }
        >
          {hasInput ? `${quote.priceImpact.toFixed(2)}%` : '-'}
        </span>
      </div>

      <div className="flex justify-between items-center text-sm">
        <span className="text-gray-400">Expected Output</span>
        <span className="text-white font-medium">
          {hasInput
            ? `${Number(formatUnits(quote.amountOut, outputToken.decimals)).toFixed(6)} ${outputSymbol}`
            : '-'}
        </span>
      </div>
    </div>
  );
}
