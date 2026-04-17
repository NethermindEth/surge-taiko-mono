import { useState, useEffect, useCallback } from 'react';
import { zeroAddress } from 'viem';
import { UniswapV2RouterABI } from '../lib/contracts';
import { L1_ROUTER, L1_DEX_WETH, USDC_TOKEN } from '../lib/constants';
import { l1PublicClient } from '../lib/config';
import { usePageVisible } from './usePageVisible';
import { SwapDirection, SwapQuote } from '../types';

interface UseL1DexQuoteParams {
  direction: SwapDirection;
  amountIn: bigint;
}

/// Quote source for swaps routed through the L1 DEX (L2→L1→L2 venue).
///
/// Calls `IUniswapV2Router02.getAmountsOut(amountIn, [WETH, USDC])` on the configured
/// L1 router. This works identically against our `SimpleDEXL1` (test mode) and a live
/// Uniswap V2 router — they share the V2 ABI.
export function useL1DexQuote({ direction, amountIn }: UseL1DexQuoteParams): SwapQuote {
  const pageVisible = usePageVisible();
  const [amountOut, setAmountOut] = useState<bigint>(0n);

  const fetchQuote = useCallback(async () => {
    if (
      !L1_ROUTER ||
      L1_ROUTER === zeroAddress ||
      !L1_DEX_WETH ||
      L1_DEX_WETH === zeroAddress ||
      !USDC_TOKEN.address ||
      amountIn === 0n
    ) {
      setAmountOut(0n);
      return;
    }

    const path =
      direction === 'ETH_TO_USDC'
        ? ([L1_DEX_WETH, USDC_TOKEN.address] as const)
        : ([USDC_TOKEN.address, L1_DEX_WETH] as const);

    try {
      const amounts = await l1PublicClient.readContract({
        address: L1_ROUTER,
        abi: UniswapV2RouterABI,
        functionName: 'getAmountsOut',
        args: [amountIn, path as readonly `0x${string}`[]],
      });
      setAmountOut(amounts[amounts.length - 1]);
    } catch (err) {
      console.warn('L1 DEX quote failed:', err);
      setAmountOut(0n);
    }
  }, [direction, amountIn]);

  useEffect(() => {
    if (!pageVisible) return;
    fetchQuote();
    const interval = setInterval(fetchQuote, 10_000);
    return () => clearInterval(interval);
  }, [fetchQuote, pageVisible]);

  const fee = (amountIn * 3n) / 1000n; // Uniswap V2 constant
  const insufficientLiquidity = amountIn > 0n && amountOut === 0n;

  // Rate approximation for display only — normalizes decimal mismatch so the UI field
  // stays readable. Price impact is omitted (would need a spot-price query).
  const inputDecimals = direction === 'ETH_TO_USDC' ? 18 : USDC_TOKEN.decimals;
  const outputDecimals = direction === 'ETH_TO_USDC' ? USDC_TOKEN.decimals : 18;
  const rate =
    amountIn > 0n && amountOut > 0n
      ? (Number(amountOut) / Number(amountIn)) *
        10 ** (inputDecimals - outputDecimals)
      : 0;

  return {
    amountOut,
    priceImpact: 0,
    fee,
    rate,
    insufficientLiquidity,
  };
}
