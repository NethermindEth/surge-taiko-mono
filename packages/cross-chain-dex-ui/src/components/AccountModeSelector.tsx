import { AccountMode } from '../types';

interface AccountModeSelectorProps {
  isOpen: boolean;
  onSelect: (mode: AccountMode) => void;
  onClose: () => void;
}

export function AccountModeSelector({ isOpen, onSelect, onClose }: AccountModeSelectorProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-surge-primary/70 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-surge-card border border-surge-border rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl shadow-surge-primary/20 hover-glow">
        <h2 className="text-xl font-bold text-surge-text mb-2">Choose Account Type</h2>
        <p className="text-surge-muted text-sm mb-6">
          Your wallet supports Ambire Smart Account (EIP-7702). Choose how you'd like to interact with Surge.
        </p>

        <div className="space-y-3">
          <button
            onClick={() => onSelect('safe')}
            className="w-full text-left p-4 bg-surge-card-hover rounded-xl border border-surge-border hover:border-surge-secondary transition-colors group"
          >
            <div className="flex items-center gap-3 mb-1">
              <div className="w-8 h-8 bg-surge-secondary/20 rounded-lg flex items-center justify-center">
                <svg className="w-4 h-4 text-surge-secondary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <span className="text-surge-text font-medium group-hover:text-surge-primary transition-colors">Safe Wallet</span>
            </div>
            <p className="text-xs text-surge-muted ml-11">
              Creates a dedicated Safe. Works with any wallet.
            </p>
          </button>

          <button
            onClick={() => onSelect('ambire')}
            className="w-full text-left p-4 bg-surge-card-hover rounded-xl border border-surge-border hover:border-surge-lavender transition-colors group"
          >
            <div className="flex items-center gap-3 mb-1">
              <div className="w-8 h-8 bg-surge-lavender/25 rounded-lg flex items-center justify-center">
                <svg className="w-4 h-4 text-surge-lavender" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <span className="text-surge-text font-medium group-hover:text-surge-lavender transition-colors">Ambire Account</span>
            </div>
            <p className="text-xs text-surge-muted ml-11">
              Uses your existing 7702 smart account. No extra wallet needed.
            </p>
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
