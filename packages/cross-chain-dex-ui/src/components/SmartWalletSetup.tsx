import toast from 'react-hot-toast';
import { Address } from 'viem';

interface SmartWalletSetupProps {
  isOpen: boolean;
  onClose: () => void;
  ownerAddress?: Address;
  isCreating: boolean;
  createSmartWallet: () => Promise<void>;
}

export function SmartWalletSetup({ isOpen, onClose, ownerAddress, isCreating, createSmartWallet }: SmartWalletSetupProps) {
  if (!isOpen) return null;

  const handleCreate = async () => {
    try {
      await createSmartWallet();
      toast.loading('Creating smart wallet...', { id: 'create-wallet' });
    } catch (error) {
      toast.error('Failed to create smart wallet');
    }
  };

  return (
    <div className="fixed inset-0 bg-surge-primary/70 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-surge-card border border-surge-border rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl shadow-surge-primary/20 hover-glow">
        <h2 className="text-xl font-bold text-surge-text mb-2">Setup Surge Smart Wallet</h2>
        <p className="text-surge-muted text-sm mb-6">
          A Safe smart wallet is required to execute cross-chain swaps. Your connected wallet will sign transactions that the Safe executes.
        </p>

        <div>
          <p className="text-sm text-surge-muted mb-4">
            This will deploy a new Safe wallet with your connected wallet as the owner.
          </p>
          <div className="bg-surge-card-hover border border-surge-border rounded-lg p-3 mb-4">
            <div className="text-xs text-surge-muted mb-1">Owner (your EOA)</div>
            <div className="text-sm text-surge-text font-mono">
              {ownerAddress}
            </div>
          </div>
          <button
            onClick={handleCreate}
            disabled={isCreating}
            className="w-full py-3 bg-surge-primary hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-opacity flex items-center justify-center gap-2 shadow-sm"
          >
            {isCreating ? (
              <>
                <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                Creating Smart Wallet...
              </>
            ) : (
              'Create Smart Wallet'
            )}
          </button>
        </div>

        <button
          onClick={onClose}
          className="w-full mt-4 py-2 text-surge-muted hover:text-surge-primary text-sm transition-colors"
        >
          Cancel
        </button>
      </div>
    </div>
  );
}
