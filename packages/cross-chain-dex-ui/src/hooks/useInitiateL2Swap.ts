import { useCallback, useState } from 'react';
import {
  Address,
  Hex,
  encodeAbiParameters,
  encodeFunctionData,
  zeroAddress,
} from 'viem';
import { useWalletClient, useSwitchChain, useConfig } from 'wagmi';
import { getWalletClient } from 'wagmi/actions';
import { CrossChainSwapVaultL2ABI, ERC20ABI } from '../lib/contracts';
import {
  L2_VAULT,
  L1_VAULT,
  L2_CHAIN_ID,
  L2_SWAP_GAS_LIMIT,
  L2_USDC_TOKEN,
} from '../lib/constants';
import { l2PublicClient } from '../lib/config';
import { simulateReturnMessage, resolveSimulatedMessage } from '../lib/catalystRpc';
import { SwapDirection } from '../types';
import { useTxStatus } from '../context/TxStatusContext';
import { useTxStatusPolling } from './useTxStatusPolling';

/// Solidity-level ordering from `CrossChainSwapVaultL2.Action`.
/// BRIDGE=0, SWAP_ETH_TO_TOKEN=1, SWAP_TOKEN_TO_ETH=2, ADD_LIQUIDITY=3,
/// REMOVE_LIQUIDITY=4, SWAP_ETH_TO_TOKEN_VIA_L1=5, SWAP_TOKEN_TO_ETH_VIA_L1=6.
const ACTION_SWAP_ETH_TO_TOKEN_VIA_L1 = 5;
const ACTION_SWAP_TOKEN_TO_ETH_VIA_L1 = 6;

const ON_MSG_INVOCATION_SELECTOR_ABI = [
  {
    type: 'function',
    name: 'onMessageInvocation',
    inputs: [{ name: '_data', type: 'bytes' }],
    outputs: [],
    stateMutability: 'payable',
  },
] as const;

type ReturnMessage = {
  id: bigint;
  fee: bigint;
  gasLimit: number;
  from: Address;
  srcChainId: bigint;
  srcOwner: Address;
  destChainId: bigint;
  destOwner: Address;
  to: Address;
  value: bigint;
  data: Hex;
};

interface InitiateParams {
  direction: SwapDirection;
  amountIn: bigint;
  minAmountOut: bigint;
  recipient: Address;
  /// Expected output from the quote hook — used ONLY to seed the placeholder return message's
  /// inner payload. Catalyst's simulation overwrites it with the real L1-computed value.
  expectedAmountOut: bigint;
}

interface UseInitiateL2SwapReturn {
  initiate: (params: InitiateParams) => Promise<boolean>;
  isPending: boolean;
  error: Error | null;
  reset: () => void;
}

/// Initiate an L2→L1→L2 swap via the L1 DEX.
///
/// Pipeline (matches the L1→L2→L1 UserOp overlay):
///   1. [signing]    Wallet approves (USDC→ETH only) + signs the L2 swap tx.
///                   The Catalyst return-message simulation also happens under
///                   this phase since it's a pre-sign RPC and resolves in <1s.
///   2. [sequencing] Tx is in Catalyst's preconf queue, awaiting Raiko proof.
///   3. [proving]    Raiko is generating the ZK proof.
///   4. [proposing]  Proof done; L1 blob tx submitted, awaiting confirmation.
///   5. [complete]   L1 inbox finalized the proposal.
///
/// Post-broadcast phases come from polling Catalyst's `surge_txStatus` by L2
/// tx hash — Catalyst's mempool scan records the tx hash and writes the same
/// status transitions as the UserOp path.
export function useInitiateL2Swap(): UseInitiateL2SwapReturn {
  const { data: walletClient } = useWalletClient();
  const { switchChainAsync } = useSwitchChain();
  const wagmiConfig = useConfig();
  const { setTxStatus } = useTxStatus();
  const { pollStatus } = useTxStatusPolling();
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const reset = useCallback(() => {
    setError(null);
    setIsPending(false);
  }, []);

  const initiate = useCallback(
    async ({
      direction,
      amountIn,
      minAmountOut,
      recipient,
      expectedAmountOut,
    }: InitiateParams): Promise<boolean> => {
      if (!walletClient) {
        setTxStatus({ phase: 'rejected', errorMessage: 'Wallet not connected' });
        return false;
      }
      if (!L2_VAULT || L2_VAULT === zeroAddress) {
        setTxStatus({ phase: 'rejected', errorMessage: 'VITE_L2_VAULT not configured' });
        return false;
      }
      if (!L2_USDC_TOKEN.address) {
        setTxStatus({ phase: 'rejected', errorMessage: 'L2 USDC token address not configured' });
        return false;
      }

      setIsPending(true);
      setError(null);

      try {
        // ------------------------------------------------------------------
        // 1. Make sure the wallet is on L2.
        // ------------------------------------------------------------------
        let activeClient = walletClient;
        if (walletClient.chain?.id !== L2_CHAIN_ID) {
          await switchChainAsync({ chainId: L2_CHAIN_ID });
          activeClient = await getWalletClient(wagmiConfig, { chainId: L2_CHAIN_ID });
        }
        const sender = activeClient.account.address;

        setTxStatus({ phase: 'signing' });

        // ------------------------------------------------------------------
        // 2. For token→ETH direction, approve L2Vault to pull the bUSDC.
        // ------------------------------------------------------------------
        if (direction === 'USDC_TO_ETH') {
          const approveHash = await activeClient.writeContract({
            address: L2_USDC_TOKEN.address!,
            abi: ERC20ABI,
            functionName: 'approve',
            args: [L2_VAULT, amountIn],
            chain: activeClient.chain,
            account: activeClient.account,
          });
          await l2PublicClient.waitForTransactionReceipt({ hash: approveHash });
        }

        // ------------------------------------------------------------------
        // 3. Ask Catalyst to resolve the real L1 return message (pre-sign).
        //    Kept under `signing` since it's <1s and precedes the wallet prompt.
        // ------------------------------------------------------------------
        const action =
          direction === 'ETH_TO_USDC'
            ? ACTION_SWAP_ETH_TO_TOKEN_VIA_L1
            : ACTION_SWAP_TOKEN_TO_ETH_VIA_L1;

        const placeholderInner = encodeAbiParameters(
          [
            { type: 'uint8' },
            { type: 'address' },
            { type: 'uint256' },
          ],
          [action, recipient, expectedAmountOut]
        );
        const placeholderOnMsgData = encodeFunctionData({
          abi: ON_MSG_INVOCATION_SELECTOR_ABI,
          functionName: 'onMessageInvocation',
          args: [placeholderInner],
        });

        const placeholderReturnMsg: ReturnMessage = {
          id: 0n,
          fee: 0n,
          gasLimit: 1_000_000,
          from: zeroAddress,
          srcChainId: 0n,
          srcOwner: L1_VAULT,
          destChainId: BigInt(L2_CHAIN_ID),
          destOwner: L2_VAULT,
          to: L2_VAULT,
          value: 0n,
          data: placeholderOnMsgData,
        };

        const simCalldata = encodeSwapCalldata({
          direction,
          minAmountOut,
          amountIn,
          recipient,
          returnMessage: placeholderReturnMsg,
        });

        const sim = await simulateReturnMessage(
          sender,
          L2_VAULT,
          simCalldata,
          direction === 'ETH_TO_USDC' ? amountIn : undefined
        );
        const realMessage = resolveSimulatedMessage(sim.message);

        // ------------------------------------------------------------------
        // 4. Wallet signs & broadcasts the swap tx.
        // ------------------------------------------------------------------
        // Viem's strict `writeContract` typing collapses `args` to `never` when combined
        // with getWalletClient's generic client type — the runtime call is correct, so
        // we cast the args array to bypass the inference failure.
        const txHash =
          direction === 'ETH_TO_USDC'
            ? await activeClient.writeContract({
                address: L2_VAULT,
                abi: CrossChainSwapVaultL2ABI,
                functionName: 'swapETHForTokenViaL1',
                args: [minAmountOut, recipient, realMessage] as never,
                value: amountIn,
                gas: L2_SWAP_GAS_LIMIT,
                chain: activeClient.chain,
                account: activeClient.account,
              })
            : await activeClient.writeContract({
                address: L2_VAULT,
                abi: CrossChainSwapVaultL2ABI,
                functionName: 'swapTokenForETHViaL1',
                args: [amountIn, minAmountOut, recipient, realMessage] as never,
                gas: L2_SWAP_GAS_LIMIT,
                chain: activeClient.chain,
                account: activeClient.account,
              });

        // ------------------------------------------------------------------
        // 5. Hand off to shared polling: drives sequencing → proving →
        //    proposing → complete from Catalyst's `surge_txStatus` RPC.
        // ------------------------------------------------------------------
        setTxStatus({ phase: 'sequencing', txHash });
        const ok = await pollStatus({ txHash });
        setIsPending(false);
        if (!ok) setError(new Error('Swap failed during proposal pipeline'));
        return ok;
      } catch (err) {
        console.error('L2→L1→L2 swap failed:', err);
        const msg = err instanceof Error ? err.message : 'Swap failed';
        setTxStatus({ phase: 'rejected', errorMessage: truncate(msg) });
        setError(err instanceof Error ? err : new Error(msg));
        setIsPending(false);
        return false;
      }
    },
    [walletClient, switchChainAsync, wagmiConfig, setTxStatus, pollStatus]
  );

  return { initiate, isPending, error, reset };
}

interface EncodeParams {
  direction: SwapDirection;
  amountIn: bigint;
  minAmountOut: bigint;
  recipient: Address;
  returnMessage: ReturnMessage;
}

function encodeSwapCalldata(params: EncodeParams): Hex {
  if (params.direction === 'ETH_TO_USDC') {
    return encodeFunctionData({
      abi: CrossChainSwapVaultL2ABI,
      functionName: 'swapETHForTokenViaL1',
      args: [params.minAmountOut, params.recipient, params.returnMessage],
    });
  }
  return encodeFunctionData({
    abi: CrossChainSwapVaultL2ABI,
    functionName: 'swapTokenForETHViaL1',
    args: [params.amountIn, params.minAmountOut, params.recipient, params.returnMessage],
  });
}

function truncate(s: string): string {
  const firstLine = s.split('\n')[0].trim();
  return firstLine.length > 140 ? `${firstLine.slice(0, 140)}…` : firstLine;
}
