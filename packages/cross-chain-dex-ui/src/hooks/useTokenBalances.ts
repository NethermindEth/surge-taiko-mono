import { useState, useEffect, useCallback } from 'react';
import { Address, formatEther, formatUnits } from 'viem';
import { ERC20ABI } from '../lib/contracts';
import { USDC_TOKEN, L2_USDC_TOKEN } from '../lib/constants';
import { l1PublicClient, l2PublicClient } from '../lib/config';
import { usePageVisible } from './usePageVisible';

type SelectedNetwork = 'l1' | 'l2';

interface TokenBalances {
  ethBalance: bigint;
  usdcBalance: bigint;
  ethFormatted: string;
  usdcFormatted: string;
  isLoading: boolean;
  error: Error | null;
  refetch: () => void;
}

export function useTokenBalances(
  smartWallet: Address | null,
  selectedNetwork: SelectedNetwork = 'l1'
): TokenBalances {
  const pageVisible = usePageVisible();
  const [ethBalance, setEthBalance] = useState<bigint>(0n);
  const [usdcBalance, setUsdcBalance] = useState<bigint>(0n);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const client = selectedNetwork === 'l2' ? l2PublicClient : l1PublicClient;
  const tokenConfig = selectedNetwork === 'l2' ? L2_USDC_TOKEN : USDC_TOKEN;

  const fetchBalances = useCallback(async () => {
    if (!smartWallet) {
      setEthBalance(0n);
      setUsdcBalance(0n);
      return;
    }

    try {
      setError(null);

      const ethBal = await client.getBalance({ address: smartWallet });
      setEthBalance(ethBal);

      if (tokenConfig.address && tokenConfig.address !== '0x0000000000000000000000000000000000000000') {
        try {
          const usdcBal = await client.readContract({
            address: tokenConfig.address,
            abi: ERC20ABI,
            functionName: 'balanceOf',
            args: [smartWallet],
          });
          setUsdcBalance(usdcBal);
        } catch {
          // Token contract may not exist on this chain
          setUsdcBalance(0n);
        }
      }
    } catch (err) {
      console.error(`Failed to fetch ${selectedNetwork} balances:`, err);
      setError(err instanceof Error ? err : new Error('Failed to fetch balances'));
    } finally {
      setIsLoading(false);
    }
  }, [smartWallet, client, selectedNetwork, tokenConfig.address]);

  // Reset balances on network switch
  useEffect(() => {
    setEthBalance(0n);
    setUsdcBalance(0n);
    setIsLoading(true);
  }, [selectedNetwork]);

  useEffect(() => {
    if (!pageVisible) return;
    fetchBalances();

    const interval = setInterval(fetchBalances, 5000);
    return () => clearInterval(interval);
  }, [fetchBalances, pageVisible]);

  return {
    ethBalance,
    usdcBalance,
    ethFormatted: formatEther(ethBalance),
    usdcFormatted: formatUnits(usdcBalance, USDC_TOKEN.decimals),
    isLoading,
    error,
    refetch: fetchBalances,
  };
}
