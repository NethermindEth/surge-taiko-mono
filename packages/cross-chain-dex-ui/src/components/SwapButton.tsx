interface SwapButtonProps {
  onClick: () => void;
  disabled: boolean;
  isLoading: boolean;
  isConnected: boolean;
  hasSmartWallet: boolean;
  hasInsufficientBalance: boolean;
  hasInsufficientLiquidity: boolean;
  hasAmount: boolean;
  exceedsSwapLimit?: string;
  needsApproval?: boolean;
}

export function SwapButton({
  onClick,
  disabled,
  isLoading,
  isConnected,
  hasSmartWallet,
  hasInsufficientBalance,
  hasInsufficientLiquidity,
  hasAmount,
  exceedsSwapLimit,
  needsApproval,
}: SwapButtonProps) {
  const getButtonText = () => {
    if (isLoading) return needsApproval ? 'Approving and Swapping...' : 'Swapping...';
    if (!isConnected) return 'Connect Wallet';
    if (!hasSmartWallet) return 'Setup Smart Wallet First';
    if (!hasAmount) return 'Enter Amount';
    if (exceedsSwapLimit) return exceedsSwapLimit;
    if (hasInsufficientLiquidity) return 'Insufficient Liquidity';
    if (hasInsufficientBalance) return 'Insufficient Balance';
    return needsApproval ? 'Approve and Swap' : 'Swap';
  };

  const isDisabled = disabled || isLoading || !isConnected || !hasSmartWallet || !hasAmount || hasInsufficientBalance || hasInsufficientLiquidity || !!exceedsSwapLimit;

  return (
    <button
      onClick={onClick}
      disabled={isDisabled}
      className={`w-full py-3 rounded-xl font-semibold text-base transition-all duration-200 ${
        isDisabled
          ? 'bg-surge-card-hover text-surge-muted cursor-not-allowed border border-surge-border'
          : 'bg-surge-primary text-white hover:bg-surge-secondary hover:shadow-md hover:shadow-surge-secondary/25 active:scale-[0.98]'
      }`}
    >
      {isLoading ? (
        <span className="flex items-center justify-center gap-2">
          <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
          {needsApproval ? 'Approving and Swapping...' : 'Swapping...'}
        </span>
      ) : (
        getButtonText()
      )}
    </button>
  );
}
