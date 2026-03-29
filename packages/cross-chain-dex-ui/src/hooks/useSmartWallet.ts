import { useState, useEffect, useCallback } from 'react';
import { Address, decodeEventLog, zeroAddress } from 'viem';
import { useAccount, useWriteContract, useWaitForTransactionReceipt, useReadContract } from 'wagmi';
import toast from 'react-hot-toast';
import { UserOpsSubmitterFactoryABI } from '../lib/contracts';
import { USER_OPS_FACTORY } from '../lib/constants';
import { useUserOp } from './useUserOp';

export function useSmartWallet() {
  const { address: ownerAddress, isConnected } = useAccount();
  const [isInitializing, setIsInitializing] = useState(true);
  const [justCreatedWallet, setJustCreatedWallet] = useState<Address | null>(null);
  const [l2WalletCreated, setL2WalletCreated] = useState(false);

  const { writeContract, data: txHash, isPending: isCreating, reset } = useWriteContract();
  const { data: receipt, isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash: txHash,
  });

  const { executeCreateL2Wallet } = useUserOp();

  // Read smart wallet from factory contract
  const { data: smartWalletFromFactory, isLoading: isLoadingFromFactory, refetch } = useReadContract({
    address: USER_OPS_FACTORY,
    abi: UserOpsSubmitterFactoryABI,
    functionName: 'getSubmitter',
    args: ownerAddress ? [ownerAddress] : undefined,
    query: {
      enabled: !!ownerAddress && isConnected,
    },
  });

  // Determine the smart wallet address (use just-created wallet or factory result)
  const smartWallet = justCreatedWallet
    ? justCreatedWallet
    : (smartWalletFromFactory && smartWalletFromFactory !== zeroAddress
        ? smartWalletFromFactory as Address
        : null);

  // Update initializing state
  useEffect(() => {
    if (!isConnected || !ownerAddress) {
      setIsInitializing(false);
      setJustCreatedWallet(null);
      setL2WalletCreated(false);
      return;
    }
    if (!isLoadingFromFactory) {
      setIsInitializing(false);
    }
  }, [isConnected, ownerAddress, isLoadingFromFactory]);

  // Handle successful wallet creation - parse event logs, then create on L2
  useEffect(() => {
    if (isSuccess && receipt && ownerAddress) {
      // Parse the SubmitterCreated event from logs
      for (const log of receipt.logs) {
        try {
          const decoded = decodeEventLog({
            abi: UserOpsSubmitterFactoryABI,
            data: log.data,
            topics: log.topics,
          });

          if (decoded.eventName === 'SubmitterCreated') {
            const createdAddress = decoded.args.submitter as Address;
            console.log('Smart wallet created on L1:', createdAddress);

            toast.dismiss('create-wallet');
            toast.success(`Smart wallet created: ${createdAddress.slice(0, 8)}...${createdAddress.slice(-6)}`);

            // Set the wallet address immediately
            setJustCreatedWallet(createdAddress);
            refetch();
            reset();

            // Trigger L2 wallet creation via UserOp
            if (!l2WalletCreated) {
              setL2WalletCreated(true);
              toast.loading('Creating L2 wallet...', { id: 'create-l2-wallet' });
              executeCreateL2Wallet({
                owner: ownerAddress,
                smartWallet: createdAddress,
              }).then((success) => {
                toast.dismiss('create-l2-wallet');
                if (success) {
                  toast.success('L2 wallet created at same address');
                } else {
                  toast.error('L2 wallet creation failed — try bridge-out later');
                }
              });
            }

            break;
          }
        } catch {
          // Not a SubmitterCreated event, continue
        }
      }
    }
  }, [isSuccess, receipt, ownerAddress, reset, refetch, executeCreateL2Wallet, l2WalletCreated]);

  const createSmartWallet = useCallback(async () => {
    if (!ownerAddress) {
      throw new Error('Wallet not connected');
    }

    writeContract({
      address: USER_OPS_FACTORY,
      abi: UserOpsSubmitterFactoryABI,
      functionName: 'createSubmitter',
      args: [ownerAddress],
    });
  }, [ownerAddress, writeContract]);

  return {
    smartWallet,
    isLoading: isInitializing || isLoadingFromFactory,
    isCreating: isCreating || isConfirming,
    createSmartWallet,
    ownerAddress,
    isConnected,
    refetch,
  };
}
