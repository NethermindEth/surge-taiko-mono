import { useState, useCallback, useMemo } from 'react';
import { parseEther, parseUnits, formatEther, formatUnits } from 'viem';
import { TokenInput } from './TokenInput';
import { useSmartWallet } from '../context/SmartWalletContext';
import { useDexReserves } from '../hooks/useDexReserves';
import { useSharedTokenBalances } from '../context/SmartWalletContext';
import { useUserOp } from '../hooks/useUserOp';
import { useSpendingLimit } from '../hooks/useSpendingLimit';
import { useLiquidityPosition } from '../hooks/useLiquidityPosition';
import { ETH_TOKEN, USDC_TOKEN, L1_NATIVE_SYMBOL } from '../lib/constants';
import { DisclaimerModal } from './DisclaimerModal';
import { useDisclaimer } from '../hooks/useDisclaimer';
import { WarningBanner } from './WarningBanner';

type LiquidityTab = 'add' | 'remove';

interface LiquidityCardProps {
  onSetupWallet: () => void;
}

export function LiquidityCard({ onSetupWallet }: LiquidityCardProps) {
  const { smartWallet, isConnected, accountMode } = useSmartWallet();
  const { ethReserve, tokenReserve } = useDexReserves();
  const { ethBalance, usdcBalance } = useSharedTokenBalances();
  const { executeAddLiquidity, executeRemoveLiquidity, isPending } = useUserOp(accountMode);
  const { hasExceededL2Limit, wouldExceed, recordSpending, remaining } = useSpendingLimit(smartWallet);
  const { isDisclaimerOpen, requireDisclaimer, onAccept, onCancel } = useDisclaimer();
  const position = useLiquidityPosition(smartWallet);

  const [tab, setTab] = useState<LiquidityTab>('add');
  const [ethInput, setEthInput] = useState('');
  const [tokenInput, setTokenInput] = useState('');
  const [priceInput, setPriceInput] = useState('1000');

  const hasReserves = ethReserve > 0n && tokenReserve > 0n;
  const price = hasReserves
    ? Number(formatUnits(tokenReserve, USDC_TOKEN.decimals)) / Number(formatEther(ethReserve))
    : Number(priceInput) || 0;

  const handleEthChange = useCallback((value: string) => {
    setEthInput(value);
    if (price > 0 && value) {
      try {
        const ethVal = Number(value);
        if (ethVal > 0) {
          setTokenInput(String(ethVal * price));
        } else {
          setTokenInput('');
        }
      } catch { /* invalid input */ }
    } else if (!value) {
      setTokenInput('');
    }
  }, [price]);

  const handleTokenChange = useCallback((value: string) => {
    setTokenInput(value);
    if (price > 0 && value) {
      try {
        const tokenVal = Number(value);
        if (tokenVal > 0) {
          setEthInput(String(tokenVal / price));
        } else {
          setEthInput('');
        }
      } catch { /* invalid input */ }
    } else if (!value) {
      setEthInput('');
    }
  }, [price]);

  const handlePriceChange = useCallback((value: string) => {
    setPriceInput(value);
    const p = Number(value);
    if (p > 0 && ethInput) {
      const ethVal = Number(ethInput);
      if (ethVal > 0) {
        setTokenInput(String(ethVal * p));
      }
    }
  }, [ethInput]);

  const ethAmount = useMemo(() => {
    try { return ethInput ? parseEther(ethInput) : 0n; } catch { return 0n; }
  }, [ethInput]);

  const tokenAmount = useMemo(() => {
    try { return tokenInput ? parseUnits(tokenInput, USDC_TOKEN.decimals) : 0n; } catch { return 0n; }
  }, [tokenInput]);

  const hasInsufficientETH = ethAmount > ethBalance;
  const hasInsufficientTokens = tokenAmount > usdcBalance;

  const liquidityUsd = (ethAmount > 0n ? Number(formatEther(ethAmount)) : 0)
    + (tokenAmount > 0n ? Number(formatUnits(tokenAmount, USDC_TOKEN.decimals)) : 0);
  const exceedsL2Limit = hasExceededL2Limit || (liquidityUsd > 0 && wouldExceed(liquidityUsd));

  const handleAddLiquidity = useCallback(async () => {
    if (!smartWallet || ethAmount === 0n || tokenAmount === 0n) return;
    const success = await executeAddLiquidity({ ethAmount, tokenAmount, smartWallet });
    if (success) {
      recordSpending(liquidityUsd);
      setEthInput('');
      setTokenInput('');
    }
  }, [smartWallet, ethAmount, tokenAmount, liquidityUsd, executeAddLiquidity, recordSpending]);

  const handleRemoveLiquidity = useCallback(async () => {
    if (!smartWallet) return;
    await executeRemoveLiquidity({ smartWallet });
  }, [smartWallet, executeRemoveLiquidity]);

  const getAddButtonText = () => {
    if (isPending) return 'Adding Liquidity...';
    if (!isConnected) return 'Connect Wallet';
    if (!smartWallet) return 'Setup Smart Wallet First';
    if (!ethAmount || !tokenAmount) return 'Enter Amounts';
    if (hasExceededL2Limit) return 'L2 deposit limit reached ($1)';
    if (exceedsL2Limit) return `Exceeds $1 limit ($${remaining.toFixed(2)} left)`;
    if (hasInsufficientETH) return `Insufficient ${L1_NATIVE_SYMBOL}`;
    if (hasInsufficientTokens) return 'Insufficient USDC Tokens';
    return 'Add Liquidity to L2';
  };

  const isAddDisabled = isPending || !isConnected || !smartWallet || !ethAmount || !tokenAmount || hasInsufficientETH || hasInsufficientTokens || exceedsL2Limit;
  const isRemoveDisabled = isPending || !isConnected || !smartWallet || !position.hasPosition;

  return (
    <div className="flex flex-col md:flex-row items-start gap-4 justify-center w-full relative z-10">
      {/* Left panel — inputs */}
      <div className="w-full md:max-w-md glass-card rounded-2xl p-4 space-y-3 hover-glow">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="inline-block w-1 h-5 rounded-full bg-surge-lavender" />
            <h2 className="text-lg font-semibold text-surge-text">Liquidity</h2>
          </div>
          <span className="text-[11px] font-semibold uppercase tracking-wider px-2.5 py-1 rounded-full border bg-surge-lavender/25 border-surge-lavender/70 text-surge-primary">
            Via Smart Account
          </span>
        </div>
        <WarningBanner />

        {/* Tab Selector */}
        <div className="flex gap-2">
          {(['add', 'remove'] as LiquidityTab[]).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`flex-1 py-2 rounded-lg text-sm font-medium transition-colors ${
                tab === t
                  ? 'bg-surge-primary text-white shadow-sm'
                  : 'bg-surge-card-hover text-surge-muted hover:text-surge-text border border-surge-border'
              }`}
            >
              {t === 'add' ? 'Add Liquidity' : 'Remove Liquidity'}
            </button>
          ))}
        </div>

        {tab === 'add' ? (
          <>
            {/* ETH Input */}
            <TokenInput
              token={ETH_TOKEN}
              amount={ethInput}
              onAmountChange={handleEthChange}
              balance={ethBalance}
              label={`${L1_NATIVE_SYMBOL} Amount`}
            />

            <div className="flex justify-center">
              <div className="text-surge-muted text-lg">+</div>
            </div>

            {/* Token Input */}
            <TokenInput
              token={USDC_TOKEN}
              amount={tokenInput}
              onAmountChange={handleTokenChange}
              balance={usdcBalance}
              label="Token Amount"
            />

            {/* Add Liquidity Button */}
            <button
              onClick={isConnected && !smartWallet ? onSetupWallet : () => requireDisclaimer(handleAddLiquidity)}
              disabled={isAddDisabled}
              className={`w-full py-3 rounded-xl font-semibold text-base transition-all duration-200 ${
                isAddDisabled
                  ? 'bg-surge-card-hover text-surge-muted cursor-not-allowed border border-surge-border'
                  : 'bg-surge-primary text-white hover:bg-surge-secondary hover:shadow-md hover:shadow-surge-secondary/25 active:scale-[0.98]'
              }`}
            >
              {isPending ? (
                <span className="flex items-center justify-center gap-2">
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Adding Liquidity...
                </span>
              ) : (
                getAddButtonText()
              )}
            </button>
          </>
        ) : (
          <>
            {/* Your Position */}
            {position.hasPosition ? (
              <div className="space-y-3">
                <div className="text-xs text-surge-muted font-medium">Your Position</div>

                {/* ETH to receive (disabled input) */}
                <div className="space-y-1">
                  <label className="text-xs text-surge-muted">{L1_NATIVE_SYMBOL} to receive</label>
                  <div className="w-full bg-surge-card-hover border border-surge-border rounded-lg px-3 py-3 text-sm text-surge-text">
                    {formatEther(position.ethAmount)} {L1_NATIVE_SYMBOL}
                  </div>
                </div>

                {/* Token to receive (disabled input) */}
                <div className="space-y-1">
                  <label className="text-xs text-surge-muted">USDC to receive</label>
                  <div className="w-full bg-surge-card-hover border border-surge-border rounded-lg px-3 py-3 text-sm text-surge-text">
                    {formatUnits(position.tokenAmount, USDC_TOKEN.decimals)} USDC
                  </div>
                </div>
              </div>
            ) : (
              <div className="bg-surge-card-hover border border-surge-border rounded-lg p-4 text-center">
                <p className="text-sm text-surge-muted">You have no liquidity position in this pool.</p>
              </div>
            )}

            {/* Remove Liquidity Button */}
            <button
              onClick={isConnected && !smartWallet ? onSetupWallet : () => requireDisclaimer(handleRemoveLiquidity)}
              disabled={isRemoveDisabled}
              className={`w-full py-3 rounded-xl font-semibold text-base transition-all duration-200 ${
                isRemoveDisabled
                  ? 'bg-surge-card-hover text-surge-muted cursor-not-allowed border border-surge-border'
                  : 'bg-surge-peach text-surge-primary hover:bg-surge-amber hover:shadow-md hover:shadow-surge-peach/40 active:scale-[0.98]'
              }`}
            >
              {isPending ? (
                <span className="flex items-center justify-center gap-2">
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Removing Liquidity...
                </span>
              ) : !isConnected ? (
                'Connect Wallet'
              ) : !smartWallet ? (
                'Setup Smart Wallet First'
              ) : !position.hasPosition ? (
                'No Liquidity to Remove'
              ) : (
                'Remove All Liquidity'
              )}
            </button>
          </>
        )}
      </div>

      {/* Right panel — pool info */}
      <div className="w-full md:max-w-sm glass-card rounded-2xl p-4 space-y-3">
        <div className="flex items-center gap-2">
          <span className="inline-block w-1 h-4 rounded-full bg-surge-mint" />
          <h3 className="text-xs font-semibold text-surge-muted uppercase tracking-widest">Pool Info</h3>
        </div>

        {/* Set initial price when pool is empty */}
        {!hasReserves && tab === 'add' && (
          <div className="bg-red-50 border border-red-300 rounded-lg p-3 space-y-2">
            <div className="text-xs text-red-700 font-medium">Pool is empty — set the initial price</div>
            <div className="flex items-center gap-2">
              <span className="text-sm text-surge-muted whitespace-nowrap">1 {L1_NATIVE_SYMBOL} =</span>
              <input
                type="number"
                value={priceInput}
                onChange={(e) => handlePriceChange(e.target.value)}
                className="flex-1 bg-surge-card border border-surge-border rounded-lg px-3 py-1.5 text-surge-text text-sm outline-none focus:border-surge-secondary"
                placeholder="1000"
                min="0"
              />
              <span className="text-sm text-surge-muted">USDC</span>
            </div>
          </div>
        )}

        {/* Current Pool Stats */}
        <div className="bg-surge-card-hover border border-surge-border rounded-lg p-3 space-y-1">
          <div className="text-xs text-surge-muted font-medium mb-2">L2 DEX Reserves</div>
          <div className="flex justify-between text-sm">
            <span className="text-surge-muted">{L1_NATIVE_SYMBOL} Reserve</span>
            <span className="text-surge-text">{formatEther(ethReserve)} {L1_NATIVE_SYMBOL}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-surge-muted">Token Reserve</span>
            <span className="text-surge-text">{formatUnits(tokenReserve, USDC_TOKEN.decimals)} USDC</span>
          </div>
          {price > 0 && (
            <div className="flex justify-between text-sm">
              <span className="text-surge-muted">Price</span>
              <span className="text-surge-text">
                1 {L1_NATIVE_SYMBOL} = {price.toFixed(2)} USDC
              </span>
            </div>
          )}
        </div>

        {/* Your Position Summary (always shown in right panel) */}
        {position.hasPosition && (
          <div className="bg-surge-card-hover border border-surge-border rounded-lg p-3 space-y-1">
            <div className="text-xs text-surge-muted font-medium mb-2">Your Liquidity</div>
            <div className="flex justify-between text-sm">
              <span className="text-surge-muted">{L1_NATIVE_SYMBOL}</span>
              <span className="text-surge-text">{formatEther(position.ethAmount)} {L1_NATIVE_SYMBOL}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-surge-muted">USDC</span>
              <span className="text-surge-text">{formatUnits(position.tokenAmount, USDC_TOKEN.decimals)} USDC</span>
            </div>
          </div>
        )}

        {/* Flow Info (add tab only) */}
        {tab === 'add' && ethAmount > 0n && tokenAmount > 0n && (
          <div className="bg-surge-card-hover border border-surge-border rounded-lg p-3 space-y-1 animate-panel-in">
            <div className="text-xs text-surge-muted font-medium mb-2">Your Deposit</div>
            <div className="flex justify-between text-sm">
              <span className="text-surge-muted">{L1_NATIVE_SYMBOL}</span>
              <span className="text-surge-text">{ethInput} {L1_NATIVE_SYMBOL}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-surge-muted">USDC</span>
              <span className="text-surge-text">{tokenInput} USDC</span>
            </div>
            <div className="flex justify-between text-sm mt-1 pt-1 border-t border-surge-border">
              <span className="text-surge-muted">Operation</span>
              <span className="text-surge-text text-xs">Lock on L1 &rarr; Add to L2 DEX</span>
            </div>
          </div>
        )}
      </div>
      <DisclaimerModal isOpen={isDisclaimerOpen} onAccept={onAccept} onCancel={onCancel} />
    </div>
  );
}
