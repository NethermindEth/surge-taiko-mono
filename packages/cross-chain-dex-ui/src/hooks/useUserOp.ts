import { useState, useCallback, useRef, useEffect, createElement } from 'react';
import { Address, Hex } from 'viem';
import { useWalletClient } from 'wagmi';
import toast from 'react-hot-toast';
import { SwapDirection } from '../types';
import {
  buildSwapUserOps,
  buildBridgeUserOps,
  buildBridgeNativeUserOps,
  buildAddLiquidityUserOps,
  computeUserOpsDigest,
  sendUserOpToBuilder,
  calculateMinOutput,
  queryUserOpStatus,
} from '../lib/userOp';
import { UserOp } from '../types';
import { DEFAULT_SLIPPAGE } from '../lib/constants';

// Inline SVG spinner components for toast icons
function Spinner({ color = '#60a5fa' }: { color?: string }) {
  return createElement('div', {
    style: {
      width: 20, height: 20,
      border: `2px solid ${color}33`,
      borderTop: `2px solid ${color}`,
      borderRadius: '50%',
      animation: 'spin 0.8s linear infinite',
    },
  });
}

function LightningSpinner() {
  return createElement('div', {
    style: { position: 'relative' as const, width: 20, height: 20 },
  },
    createElement('div', {
      style: {
        position: 'absolute' as const, inset: 0,
        border: '2px solid #a78bfa33',
        borderTop: '2px solid #a78bfa',
        borderRadius: '50%',
        animation: 'spin 0.8s linear infinite',
      },
    }),
    createElement('svg', {
      viewBox: '0 0 24 24', fill: '#a78bfa',
      style: { position: 'absolute' as const, inset: 3, width: 14, height: 14 },
    },
      createElement('path', { d: 'M13 2L3 14h9l-1 8 10-12h-9l1-8z' })
    ),
  );
}

function PenSpinner() {
  return createElement('div', {
    style: { position: 'relative' as const, width: 20, height: 20 },
  },
    createElement('div', {
      style: {
        position: 'absolute' as const, inset: 0,
        border: '2px solid #fbbf2433',
        borderTop: '2px solid #fbbf24',
        borderRadius: '50%',
        animation: 'spin 0.8s linear infinite',
      },
    }),
    createElement('svg', {
      viewBox: '0 0 24 24', fill: 'none', stroke: '#fbbf24', strokeWidth: 2,
      strokeLinecap: 'round' as const, strokeLinejoin: 'round' as const,
      style: { position: 'absolute' as const, inset: 3, width: 14, height: 14 },
    },
      createElement('path', { d: 'M17 3a2.83 2.83 0 114 4L7.5 20.5 2 22l1.5-5.5L17 3z' })
    ),
  );
}

interface UseUserOpReturn {
  executeSwap: (params: ExecuteSwapParams) => Promise<boolean>;
  executeBridge: (params: ExecuteBridgeParams) => Promise<boolean>;
  executeBridgeNative: (params: ExecuteBridgeNativeParams) => Promise<boolean>;
  executeAddLiquidity: (params: ExecuteAddLiquidityParams) => Promise<boolean>;
  isPending: boolean;
  error: Error | null;
}

interface ExecuteBridgeNativeParams {
  amount: bigint;
  recipient: Address;
  smartWallet: Address;
}

interface ExecuteSwapParams {
  direction: SwapDirection;
  amountIn: bigint;
  expectedAmountOut: bigint;
  smartWallet: Address;
  slippage?: number;
}

interface ExecuteBridgeParams {
  amount: bigint;
  recipient: Address;
  smartWallet: Address;
}

interface ExecuteAddLiquidityParams {
  ethAmount: bigint;
  tokenAmount: bigint;
  smartWallet: Address;
}

export function useUserOp(): UseUserOpReturn {
  const { data: walletClient } = useWalletClient();
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const pollIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Clean up polling on unmount
  useEffect(() => {
    return () => {
      if (pollIntervalRef.current) {
        clearInterval(pollIntervalRef.current);
      }
    };
  }, []);

  const pollStatus = useCallback((userOpId: number, toastId: string = 'swap'): Promise<boolean> => {
    return new Promise((resolve) => {
      toast.loading('Sending to builder', { id: toastId, icon: createElement(Spinner) });

      pollIntervalRef.current = setInterval(async () => {
        const status = await queryUserOpStatus(userOpId);
        if (!status) return;

        if (status.status === 'Pending') {
          toast.loading('Pending', { id: toastId, icon: createElement(Spinner) });
        } else if (status.status === 'Processing') {
          toast.loading(`Proposing (tx: ${status.tx_hash.slice(0, 10)})`, { id: toastId, icon: createElement(Spinner, { color: '#34d399' }) });
        } else if (status.status === 'ProvingBlock') {
          toast.loading('Generating ZK Proof', { id: toastId, icon: createElement(LightningSpinner) });
        } else if (status.status === 'Executed') {
          if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
          pollIntervalRef.current = null;
          toast.success('Execution complete', { id: toastId });
          setIsPending(false);
          resolve(true);
        } else if (status.status === 'Rejected') {
          if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
          pollIntervalRef.current = null;
          toast.error(`Rejected: ${status.reason}`, { id: toastId });
          setError(new Error(status.reason));
          setIsPending(false);
          resolve(false);
        }
      }, 1000);
    });
  }, []);

  const executeSwap = useCallback(
    async ({
      direction,
      amountIn,
      expectedAmountOut,
      smartWallet,
      slippage = DEFAULT_SLIPPAGE,
    }: ExecuteSwapParams): Promise<boolean> => {
      if (!walletClient) {
        toast.error('Wallet not connected');
        return false;
      }

      setIsPending(true);
      setError(null);

      try {
        // Calculate minimum output with slippage
        const minAmountOut = calculateMinOutput(expectedAmountOut, slippage);

        // Build UserOp(s)
        const ops = buildSwapUserOps(direction, amountIn, minAmountOut, smartWallet);

        toast.loading('Signing transaction', { id: 'swap', icon: createElement(PenSpinner) });

        // Compute digest
        const digest = computeUserOpsDigest(ops);
        console.log('UserOps:', ops);
        console.log('Digest to sign:', digest);
        console.log('Signer address:', walletClient.account.address);

        // Sign the digest using signMessage (standard personal_sign)
        // The contract uses MessageHashUtils.toEthSignedMessageHash to add the Ethereum prefix
        // before recovering the signer address
        const signature = await walletClient.signMessage({
          message: { raw: digest as `0x${string}` },
        });
        console.log('Signature:', signature);

        // Send to builder RPC
        const result = await sendUserOpToBuilder(smartWallet, ops, signature as Hex);

        if (result.success && result.userOpId !== undefined) {
          // Poll for status
          return await pollStatus(result.userOpId, 'swap');
        } else if (result.success) {
          // No ID returned, can't poll - just show success
          toast.success('Swap submitted successfully!', { id: 'swap' });
          setIsPending(false);
          return true;
        } else {
          toast.error(result.error || 'Failed to submit swap', { id: 'swap' });
          setError(new Error(result.error || 'Failed to submit swap'));
          setIsPending(false);
          return false;
        }
      } catch (err) {
        console.error('Swap failed:', err);
        const errorMessage = err instanceof Error ? err.message : 'Swap failed';
        toast.error(errorMessage, { id: 'swap' });
        setError(err instanceof Error ? err : new Error(errorMessage));
        setIsPending(false);
        return false;
      }
    },
    [walletClient, pollStatus]
  );

  const executeGenericOps = useCallback(
    async (ops: UserOp[], smartWallet: Address, toastId: string, successMsg: string): Promise<boolean> => {
      if (!walletClient) {
        toast.error('Wallet not connected');
        return false;
      }

      setIsPending(true);
      setError(null);

      try {
        toast.loading('Signing transaction', { id: toastId, icon: createElement(PenSpinner) });

        const digest = computeUserOpsDigest(ops);
        const signature = await walletClient.signMessage({
          message: { raw: digest as `0x${string}` },
        });

        const result = await sendUserOpToBuilder(smartWallet, ops, signature as Hex);

        if (result.success && result.userOpId !== undefined) {
          return await pollStatus(result.userOpId, toastId);
        } else if (result.success) {
          toast.success(successMsg, { id: toastId });
          setIsPending(false);
          return true;
        } else {
          toast.error(result.error || 'Failed to submit', { id: toastId });
          setError(new Error(result.error || 'Failed to submit'));
          setIsPending(false);
          return false;
        }
      } catch (err) {
        console.error('Operation failed:', err);
        const errorMessage = err instanceof Error ? err.message : 'Operation failed';
        toast.error(errorMessage, { id: toastId });
        setError(err instanceof Error ? err : new Error(errorMessage));
        setIsPending(false);
        return false;
      }
    },
    [walletClient, pollStatus]
  );

  const executeBridge = useCallback(
    async ({ amount, recipient, smartWallet }: ExecuteBridgeParams): Promise<boolean> => {
      const ops = buildBridgeUserOps(amount, recipient);
      return executeGenericOps(ops, smartWallet, 'bridge', 'Bridge submitted successfully!');
    },
    [executeGenericOps]
  );

  const executeBridgeNative = useCallback(
    async ({ amount, recipient, smartWallet }: ExecuteBridgeNativeParams): Promise<boolean> => {
      const ops = buildBridgeNativeUserOps(amount, recipient, smartWallet);
      return executeGenericOps(ops, smartWallet, 'bridge', 'Native bridge submitted!');
    },
    [executeGenericOps]
  );

  const executeAddLiquidity = useCallback(
    async ({ ethAmount, tokenAmount, smartWallet }: ExecuteAddLiquidityParams): Promise<boolean> => {
      const ops = buildAddLiquidityUserOps(ethAmount, tokenAmount);
      return executeGenericOps(ops, smartWallet, 'liquidity', 'Liquidity addition submitted!');
    },
    [executeGenericOps]
  );

  return {
    executeSwap,
    executeBridge,
    executeBridgeNative,
    executeAddLiquidity,
    isPending,
    error,
  };
}
