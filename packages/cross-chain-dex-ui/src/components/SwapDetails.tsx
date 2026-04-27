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
    <div className="bg-surge-card-hover border border-surge-border rounded-xl p-4 space-y-3">
      <div className="flex justify-between items-center text-sm">
        <span className="text-surge-muted">Rate</span>
        <span className="text-surge-text">
          {hasInput ? `1 ${inputSymbol} = ${quote.rate.toFixed(6)} ${outputSymbol}` : '-'}
        </span>
      </div>

      <div className="flex justify-between items-center text-sm">
        <span className="text-surge-muted">Fee ({FEE_PERCENT}%)</span>
        <span className="text-surge-text">
          {hasInput
            ? `${Number(formatUnits(quote.fee, inputToken.decimals)).toFixed(6)} ${inputSymbol}`
            : '-'}
        </span>
      </div>

      <div className="flex justify-between items-center text-sm">
        <span className="text-surge-muted">Price Impact</span>
        <span
          className={
            hasInput
              ? quote.priceImpact > 5
                ? 'text-surge-amber'
                : quote.priceImpact > 1
                ? 'text-surge-peach'
                : 'text-surge-primary'
              : 'text-surge-text'
          }
        >
          {hasInput ? `${quote.priceImpact.toFixed(2)}%` : '-'}
        </span>
      </div>

      <div className="flex justify-between items-center text-sm">
        <span className="text-surge-muted">Expected Output</span>
        <span className="text-surge-text font-medium">
          {hasInput
            ? `${Number(formatUnits(quote.amountOut, outputToken.decimals)).toFixed(6)} ${outputSymbol}`
            : '-'}
        </span>
      </div>
    </div>
  );
}
