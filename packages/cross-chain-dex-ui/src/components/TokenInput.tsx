import { Token } from '../types';
import { formatUnits } from 'viem';

interface TokenInputProps {
  token: Token;
  amount: string;
  onAmountChange: (value: string) => void;
  balance: bigint;
  label: string;
  disabled?: boolean;
  showMax?: boolean;
}

export function TokenInput({
  token,
  amount,
  onAmountChange,
  balance,
  label,
  disabled = false,
  showMax = true,
}: TokenInputProps) {
  const handleMaxClick = () => {
    onAmountChange(formatUnits(balance, token.decimals));
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    // Allow empty string, numbers, and one decimal point
    if (value === '' || /^\d*\.?\d*$/.test(value)) {
      onAmountChange(value);
    }
  };

  return (
    <div className="bg-surge-card-hover rounded-xl p-3 border border-surge-border">
      <div className="flex justify-between items-center mb-2">
        <span className="text-sm text-surge-muted">{label}</span>
        <div className="flex items-center gap-2">
          <span className="text-sm text-surge-muted">
            Balance: {Number(formatUnits(balance, token.decimals)).toFixed(4)}
          </span>
          {showMax && !disabled && balance > 0n && (
            <button
              onClick={handleMaxClick}
              className="text-xs font-semibold tracking-wider text-surge-primary px-2.5 py-0.5 bg-surge-lavender/40 border border-surge-lavender/70 rounded-full hover:bg-surge-lavender/60 transition-colors"
            >
              MAX
            </button>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3">
        {/* Token selector - fixed width */}
        <div className="flex items-center gap-2 bg-surge-card px-3 py-2 rounded-lg shrink-0 border border-surge-border">
          <img src={token.logo} alt={token.symbol} className="w-6 h-6" />
          <span className="text-surge-text font-medium">{token.symbol}</span>
        </div>

        {/* Input - takes remaining space */}
        <input
          type="text"
          value={amount}
          onChange={handleChange}
          disabled={disabled}
          placeholder="0.0"
          className="flex-1 min-w-0 bg-transparent text-2xl text-surge-text text-right outline-none placeholder-surge-muted/50 disabled:opacity-50"
        />
      </div>
    </div>
  );
}
