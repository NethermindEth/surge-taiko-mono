import { useState, useCallback, useRef, useEffect } from 'react';
import { Address, Hex } from 'viem';
import { useWalletClient } from 'wagmi';
import { SwapDirection } from '../types';
import {
  buildSwapUserOps,
  buildBridgeUserOps,
  buildBridgeNativeUserOps,
  buildBridgeOutNativeUserOps,
  buildCreateL2WalletUserOps,
  buildAddLiquidityUserOps,
  buildExecuteBatchTypedData,
  sendUserOpToBuilder,
  calculateMinOutput,
  queryUserOpStatus,
} from '../lib/userOp';
import { UserOp } from '../types';
import { DEFAULT_SLIPPAGE, L2_CHAIN_ID } from '../lib/constants';
import { useTxStatus } from '../context/TxStatusContext';

interface ExecuteBridgeOutNativeParams {
  amount: bigint;
  recipient: Address;
  smartWallet: Address;
}

interface ExecuteCreateL2WalletParams {
  owner: Address;
  smartWallet: Address;
}

interface UseUserOpReturn {
  executeSwap: (params: ExecuteSwapParams) => Promise<boolean>;
  executeBridge: (params: ExecuteBridgeParams) => Promise<boolean>;
  executeBridgeNative: (params: ExecuteBridgeNativeParams) => Promise<boolean>;
  executeBridgeOutNative: (params: ExecuteBridgeOutNativeParams) => Promise<boolean>;
  executeCreateL2Wallet: (params: ExecuteCreateL2WalletParams) => Promise<boolean>;
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

        const typedData = buildExecuteBatchTypedData(smartWallet, ops);
        const signature = await walletClient.signTypedData(typedData);
        console.log('Signature:', signature);

        const result = await sendUserOpToBuilder(smartWallet, ops, signature as Hex);

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

        const typedData = buildExecuteBatchTypedData(smartWallet, ops, chainId);
        const signature = await walletClient.signTypedData(typedData);

        const result = await sendUserOpToBuilder(smartWallet, ops, signature as Hex, chainId);

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
    [walletClient, pollStatus, setTxStatus]
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

  const executeBridgeOutNative = useCallback(
    async ({ amount, recipient, smartWallet }: ExecuteBridgeOutNativeParams): Promise<boolean> => {
      const ops = buildBridgeOutNativeUserOps(amount, recipient, smartWallet);
      // Sign with L2 chain ID — catalyst auto-detects the target chain from the signature
      return executeGenericOps(ops, smartWallet, L2_CHAIN_ID);
    },
    [executeGenericOps]
  );

  const executeCreateL2Wallet = useCallback(
    async ({ owner, smartWallet }: ExecuteCreateL2WalletParams): Promise<boolean> => {
      const ops = buildCreateL2WalletUserOps(owner, smartWallet);
      // L1 UserOp — bridge.sendMessage routes createSubmitter to L2 via processMessage
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

  return {
    executeSwap,
    executeBridge,
    executeBridgeNative,
    executeBridgeOutNative,
    executeCreateL2Wallet,
    executeAddLiquidity,
    isPending,
    error,
  };
}
