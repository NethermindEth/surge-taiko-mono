import { useCallback, useState } from 'react';
import { Address, encodeFunctionData, zeroAddress, Hex } from 'viem';
import { useWalletClient } from 'wagmi';
import { l2PublicClient } from '../lib/config';
import { L2_BRIDGE, CHAIN_ID } from '../lib/constants';
import { BridgeABI } from '../lib/contracts';
import { useTxStatus } from '../context/TxStatusContext';

interface BridgeOutParams {
  amount: bigint;
  recipient: Address;
}

/// Bridge native L2 funds back to L1 via a direct EOA tx to the L2 bridge.
/// Mirrors the calldata produced by `buildBridgeOutNativeUserOps` but skips
/// the Safe / L1 UserOp builder paths — on the L2 page the user signs from
/// their EOA, so this is just a standard L2 transaction.
export function useBridgeOutEoa() {
  const { data: walletClient } = useWalletClient();
  const { setTxStatus } = useTxStatus();
  const [isPending, setIsPending] = useState(false);

  const initiate = useCallback(
    async ({ amount, recipient }: BridgeOutParams): Promise<boolean> => {
      if (!walletClient) {
        setTxStatus({ phase: 'rejected', errorMessage: 'Wallet not connected' });
        return false;
      }
      setIsPending(true);
      try {
        setTxStatus({ phase: 'signing' });
        const sender = walletClient.account.address;
        const calldata = encodeFunctionData({
          abi: BridgeABI,
          functionName: 'sendMessage',
          args: [{
            id: 0n,
            fee: 0n,
            gasLimit: 1_000_000,
            from: zeroAddress,
            srcChainId: 0n,
            srcOwner: sender,
            destChainId: BigInt(CHAIN_ID),
            destOwner: recipient,
            to: recipient,
            value: amount,
            data: '0x' as Hex,
          }],
        });
        const txHash = await walletClient.sendTransaction({
          to: L2_BRIDGE,
          value: amount,
          data: calldata,
          chain: walletClient.chain,
          account: walletClient.account,
        });
        setTxStatus({ phase: 'sequencing' });
        await l2PublicClient.waitForTransactionReceipt({ hash: txHash });
        setTxStatus({ phase: 'complete', txHash });
        setIsPending(false);
        return true;
      } catch (err) {
        const raw = err instanceof Error ? err.message : 'Bridge withdrawal failed';
        const msg = raw.includes('rejected') || raw.includes('denied')
          ? 'Transaction rejected by user'
          : raw.split(/[.\n]/)[0].trim().slice(0, 160);
        setTxStatus({ phase: 'rejected', errorMessage: msg });
        setIsPending(false);
        return false;
      }
    },
    [walletClient, setTxStatus]
  );

  return { initiate, isPending };
}
