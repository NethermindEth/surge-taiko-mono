import { useState, useCallback, useRef, useEffect } from 'react';
import { Address, Hex } from 'viem';
import { useWalletClient, useSwitchChain } from 'wagmi';
import { SwapDirection } from '../types';
import {
  buildSwapUserOps,
  buildBridgeUserOps,
  buildBridgeNativeUserOps,
  buildAddLiquidityUserOps,
  buildRemoveLiquidityUserOps,
  userOpsToSafeTx,
  sendUserOpToBuilder,
  calculateMinOutput,
  queryUserOpStatus,
} from '../lib/userOp';
import { getSafeNonce, buildSafeTxTypedData, buildExecTransactionCalldata } from '../lib/safeOp';
import { UserOp } from '../types';
import { CHAIN_ID, L2_CHAIN_ID, DEFAULT_SLIPPAGE } from '../lib/constants';
import { l1PublicClient, l2PublicClient } from '../lib/config';
import { useTxStatus } from '../context/TxStatusContext';

interface UseUserOpReturn {
  executeSwap: (params: ExecuteSwapParams) => Promise<boolean>;
  executeBridge: (params: ExecuteBridgeParams) => Promise<boolean>;
  executeBridgeNative: (params: ExecuteBridgeNativeParams) => Promise<boolean>;
  executeAddLiquidity: (params: ExecuteAddLiquidityParams) => Promise<boolean>;
  executeRemoveLiquidity: (params: { smartWallet: Address }) => Promise<boolean>;
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
  const { switchChainAsync } = useSwitchChain();
  const { setTxStatus } = useTxStatus();
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const pollIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const txHashRef = useRef<string | undefined>(undefined);

  useEffect(() => {
    return () => {
      if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
    };
  }, []);

  const pollStatus = useCallback((userOpId: number): Promise<boolean> => {
    return new Promise((resolve) => {
      setTxStatus({ phase: 'sequencing' });

      pollIntervalRef.current = setInterval(async () => {
        const status = await queryUserOpStatus(userOpId);
        if (!status) return;

        if (status.status === 'Pending') {
          setTxStatus({ phase: 'sequencing' });
        } else if (status.status === 'Processing') {
          txHashRef.current = status.tx_hash;
          setTxStatus({ phase: 'proposing' });
        } else if (status.status === 'ProvingBlock') {
          setTxStatus({ phase: 'proving' });
        } else if (status.status === 'Executed') {
          if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
          pollIntervalRef.current = null;
          setTxStatus({ phase: 'complete', txHash: txHashRef.current });
          setIsPending(false);
          resolve(true);
        } else if (status.status === 'Rejected') {
          if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
          pollIntervalRef.current = null;
          setTxStatus({ phase: 'rejected', errorMessage: status.reason });
          setError(new Error(status.reason));
          setIsPending(false);
          resolve(false);
        }
      }, 1000);
    });
  }, [setTxStatus]);

  const executeSwap = useCallback(
    async ({
      direction,
      amountIn,
      expectedAmountOut,
      smartWallet,
      slippage = DEFAULT_SLIPPAGE,
    }: ExecuteSwapParams): Promise<boolean> => {
      if (!walletClient) {
        setTxStatus({ phase: 'rejected', errorMessage: 'Wallet not connected' });
        return false;
      }

      setIsPending(true);
      setError(null);
      txHashRef.current = undefined;

      try {
        const minAmountOut = calculateMinOutput(expectedAmountOut, slippage);
        const ops = buildSwapUserOps(direction, amountIn, minAmountOut, smartWallet);

        setTxStatus({ phase: 'signing' });

        // Fetch nonce from the Safe on L1
        const nonce = await getSafeNonce(l1PublicClient, smartWallet);

        // Convert ops to a single SafeTxParams
        const safeTx = userOpsToSafeTx(ops);

        // Build Safe EIP-712 typed data
        const typedData = buildSafeTxTypedData(smartWallet, CHAIN_ID, nonce, safeTx);

        const signature = await walletClient.signTypedData(typedData);
        console.log('Signature:', signature);

        // Encode execTransaction calldata
        const calldata = buildExecTransactionCalldata(safeTx, signature as Hex);

        const result = await sendUserOpToBuilder(smartWallet, calldata);

        if (result.success && result.userOpId !== undefined) {
          return await pollStatus(result.userOpId);
        } else if (result.success) {
          setTxStatus({ phase: 'complete' });
          setIsPending(false);
          return true;
        } else {
          setTxStatus({ phase: 'rejected', errorMessage: result.error || 'Failed to submit swap' });
          setError(new Error(result.error || 'Failed to submit swap'));
          setIsPending(false);
          return false;
        }
      } catch (err) {
        console.error('Swap failed:', err);
        const msg = err instanceof Error ? err.message : 'Swap failed';
        setTxStatus({ phase: 'rejected', errorMessage: msg });
        setError(err instanceof Error ? err : new Error(msg));
        setIsPending(false);
        return false;
      }
    },
    [walletClient, pollStatus, setTxStatus]
  );

  const executeGenericOps = useCallback(
    async (ops: UserOp[], smartWallet: Address, chainId?: number): Promise<boolean> => {
      if (!walletClient) {
        setTxStatus({ phase: 'rejected', errorMessage: 'Wallet not connected' });
        return false;
      }

      setIsPending(true);
      setError(null);
      txHashRef.current = undefined;

      try {
        setTxStatus({ phase: 'signing' });

        // Determine which chain this Safe lives on
        const targetChainId = chainId ?? CHAIN_ID;
        const publicClient = targetChainId === L2_CHAIN_ID ? l2PublicClient : l1PublicClient;

        // Switch chain if needed (e.g. bridge-out: signing on L2)
        if (chainId !== undefined && chainId !== walletClient.chain?.id) {
          await switchChainAsync({ chainId });
        }

        // Fetch nonce from the Safe on the correct chain
        const nonce = await getSafeNonce(publicClient, smartWallet);

        // Convert ops to a single SafeTxParams
        const safeTx = userOpsToSafeTx(ops);

        // Build Safe EIP-712 typed data
        const typedData = buildSafeTxTypedData(smartWallet, targetChainId, nonce, safeTx);

        const signature = await walletClient.signTypedData(typedData);

        // Encode execTransaction calldata
        const calldata = buildExecTransactionCalldata(safeTx, signature as Hex);

        const result = await sendUserOpToBuilder(smartWallet, calldata, chainId);

        if (result.success && result.userOpId !== undefined) {
          return await pollStatus(result.userOpId);
        } else if (result.success) {
          setTxStatus({ phase: 'complete' });
          setIsPending(false);
          return true;
        } else {
          setTxStatus({ phase: 'rejected', errorMessage: result.error || 'Failed to submit' });
          setError(new Error(result.error || 'Failed to submit'));
          setIsPending(false);
          return false;
        }
      } catch (err) {
        console.error('Operation failed:', err);
        const msg = err instanceof Error ? err.message : 'Operation failed';
        setTxStatus({ phase: 'rejected', errorMessage: msg });
        setError(err instanceof Error ? err : new Error(msg));
        setIsPending(false);
        return false;
      }
    },
    [walletClient, switchChainAsync, pollStatus, setTxStatus]
  );

  const executeBridge = useCallback(
    async ({ amount, recipient, smartWallet }: ExecuteBridgeParams): Promise<boolean> => {
      const ops = buildBridgeUserOps(amount, recipient);
      return executeGenericOps(ops, smartWallet);
    },
    [executeGenericOps]
  );

  const executeBridgeNative = useCallback(
    async ({ amount, recipient, smartWallet }: ExecuteBridgeNativeParams): Promise<boolean> => {
      const ops = buildBridgeNativeUserOps(amount, recipient, smartWallet);
      return executeGenericOps(ops, smartWallet);
    },
    [executeGenericOps]
  );

  const executeAddLiquidity = useCallback(
    async ({ ethAmount, tokenAmount, smartWallet }: ExecuteAddLiquidityParams): Promise<boolean> => {
      const ops = buildAddLiquidityUserOps(ethAmount, tokenAmount);
      return executeGenericOps(ops, smartWallet);
    },
    [executeGenericOps]
  );

  const executeRemoveLiquidity = useCallback(
    async ({ smartWallet }: { smartWallet: Address }): Promise<boolean> => {
      const ops = buildRemoveLiquidityUserOps();
      return executeGenericOps(ops, smartWallet);
    },
    [executeGenericOps]
  );

  return {
    executeSwap,
    executeBridge,
    executeBridgeNative,
    executeAddLiquidity,
    executeRemoveLiquidity,
    isPending,
    error,
  };
}
