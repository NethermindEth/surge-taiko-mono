import { useState, useEffect, useCallback } from 'react';
import { Address, decodeEventLog, zeroAddress } from 'viem';
import { useAccount, useWriteContract, useWaitForTransactionReceipt, useReadContract } from 'wagmi';
import toast from 'react-hot-toast';
import { UserOpsSubmitterFactoryABI } from '../lib/contracts';
import { USER_OPS_FACTORY } from '../lib/constants';
import { l2PublicClient } from '../lib/config';
import { useUserOp } from './useUserOp';

export function useSmartWallet() {
  const { address: ownerAddress, isConnected } = useAccount();
  const [isInitializing, setIsInitializing] = useState(true);
  const [justCreatedWallet, setJustCreatedWallet] = useState<Address | null>(null);
  const [l2WalletExists, setL2WalletExists] = useState<boolean | null>(null);
  const [isCreatingL2, setIsCreatingL2] = useState(false);

  const { writeContract, data: txHash, isPending: isCreating, reset } = useWriteContract();
  const { data: receipt, isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash: txHash,
  });

  const { executeCreateL2Wallet } = useUserOp();

  // Read smart wallet from L1 factory
  const { data: smartWalletFromFactory, isLoading: isLoadingFromFactory, refetch } = useReadContract({
    address: USER_OPS_FACTORY,
    abi: UserOpsSubmitterFactoryABI,
    functionName: 'getSubmitter',
    args: ownerAddress ? [ownerAddress] : undefined,
    query: {
      enabled: !!ownerAddress && isConnected,
    },
  });

  // Determine the smart wallet address
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
      setL2WalletExists(null);
      return;
    }
    if (!isLoadingFromFactory) {
      setIsInitializing(false);
    }
  }, [isConnected, ownerAddress, isLoadingFromFactory]);

  // Check if L2 wallet exists whenever L1 wallet is known
  useEffect(() => {
    if (!smartWallet || !ownerAddress) {
      setL2WalletExists(null);
      return;
    }

    l2PublicClient.readContract({
      address: USER_OPS_FACTORY,
      abi: UserOpsSubmitterFactoryABI,
      functionName: 'getSubmitter',
      args: [ownerAddress],
    }).then((l2Submitter) => {
      setL2WalletExists(l2Submitter !== zeroAddress);
    }).catch(() => {
      setL2WalletExists(null);
    });
  }, [smartWallet, ownerAddress]);

  // Handle successful L1 wallet creation - parse event logs
  useEffect(() => {
    if (isSuccess && receipt && ownerAddress) {
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

            setJustCreatedWallet(createdAddress);
            refetch();
            reset();
            break;
          }
        } catch {
          // Not a SubmitterCreated event, continue
        }
      }
    }
  }, [isSuccess, receipt, ownerAddress, reset, refetch]);

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

  // Explicitly create L2 wallet — called by the UI after funding
  const createL2Wallet = useCallback(async () => {
    if (!smartWallet || !ownerAddress) return false;
    if (l2WalletExists) return true;

    setIsCreatingL2(true);
    toast.loading('Creating L2 wallet...', { id: 'create-l2-wallet' });

    const success = await executeCreateL2Wallet({
      owner: ownerAddress,
      smartWallet,
    });

    toast.dismiss('create-l2-wallet');
    setIsCreatingL2(false);

    if (success) {
      toast.success('L2 wallet created at same address');
      setL2WalletExists(true);
    } else {
      toast.error('L2 wallet creation failed');
    }

    return success;
  }, [smartWallet, ownerAddress, l2WalletExists, executeCreateL2Wallet]);

  return {
    smartWallet,
    l2WalletExists,
    isLoading: isInitializing || isLoadingFromFactory,
    isCreating: isCreating || isConfirming,
    isCreatingL2,
    createSmartWallet,
    createL2Wallet,
    ownerAddress,
    isConnected,
    refetch,
  };
}
