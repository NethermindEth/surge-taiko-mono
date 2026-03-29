import { useState, useEffect, useCallback } from 'react';
import { Address, formatEther, formatUnits } from 'viem';
import { ERC20ABI } from '../lib/contracts';
import { USDC_TOKEN } from '../lib/constants';
import { l2PublicClient } from '../lib/config';

interface L2TokenBalances {
  ethBalance: bigint;
  usdcBalance: bigint;
  ethFormatted: string;
  usdcFormatted: string;
  isLoading: boolean;
  refetch: () => void;
}

export function useL2TokenBalances(account: Address | null): L2TokenBalances {
  const [ethBalance, setEthBalance] = useState<bigint>(0n);
  const [usdcBalance, setUsdcBalance] = useState<bigint>(0n);
  const [isLoading, setIsLoading] = useState(true);

  const fetchBalances = useCallback(async () => {
    if (!account) {
      setEthBalance(0n);
      setUsdcBalance(0n);
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);

      const ethBal = await l2PublicClient.getBalance({ address: account });
      setEthBalance(ethBal);

      if (USDC_TOKEN.address && USDC_TOKEN.address !== '0x0000000000000000000000000000000000000000') {
        try {
          const usdcBal = await l2PublicClient.readContract({
            address: USDC_TOKEN.address,
            abi: ERC20ABI,
            functionName: 'balanceOf',
            args: [account],
          });
          setUsdcBalance(usdcBal);
        } catch {
          // Token may not exist on L2
          setUsdcBalance(0n);
        }
      }
    } catch (err) {
      console.error('Failed to fetch L2 balances:', err);
    } finally {
      setIsLoading(false);
    }
  }, [account]);

  useEffect(() => {
    fetchBalances();
    const interval = setInterval(fetchBalances, 5000);
    return () => clearInterval(interval);
  }, [fetchBalances]);

  return {
    ethBalance,
    usdcBalance,
    ethFormatted: formatEther(ethBalance),
    usdcFormatted: formatUnits(usdcBalance, USDC_TOKEN.decimals),
    isLoading,
    refetch: fetchBalances,
  };
}
