import { Address } from 'viem';
import { useAccount } from 'wagmi';

// TODO(Task 5): Rewrite to use Gnosis Safe (1-of-1) instead of UserOpsSubmitter.
// The UserOpsSubmitterFactoryABI / USER_OPS_FACTORY references have been removed.
// For now this hook returns smartWallet: null so the app compiles.

export function useSmartWallet() {
  const { address: ownerAddress, isConnected } = useAccount();

  return {
    smartWallet: null as Address | null,
    isLoading: false,
    isCreating: false,
    createSmartWallet: async () => {
      throw new Error('Safe wallet creation not yet implemented (Task 5)');
    },
    ownerAddress,
    isConnected,
    refetch: async () => {},
  };
}
