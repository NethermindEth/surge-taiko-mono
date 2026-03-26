import { useState, useCallback, useMemo } from "react";
import { parseUnits, formatUnits, Address } from "viem";
import { TokenInput } from "./TokenInput";
import { useSmartWallet } from "../hooks/useSmartWallet";
import { useTokenBalances } from "../hooks/useTokenBalances";
import { useUserOp } from "../hooks/useUserOp";
import { useSpendingLimit } from "../hooks/useSpendingLimit";
import { ETH_TOKEN, USDC_TOKEN, L1_NATIVE_SYMBOL } from "../lib/constants";

type BridgeToken = typeof L1_NATIVE_SYMBOL | "USDC";

interface BridgeCardProps {
  onSetupWallet: () => void;
}

export function BridgeCard({ onSetupWallet }: BridgeCardProps) {
  const { smartWallet, isConnected } = useSmartWallet();
  const { ethBalance, usdcBalance } = useTokenBalances(smartWallet);
  const { executeBridge, executeBridgeNative, isPending } = useUserOp();
  const { hasExceededL2Limit, wouldExceed, recordSpending, remaining } = useSpendingLimit(smartWallet);

  const [bridgeToken, setBridgeToken] = useState<BridgeToken>("USDC");
  const [inputAmount, setInputAmount] = useState("");
  const [recipient, setRecipient] = useState("");

  const currentToken =
    bridgeToken === L1_NATIVE_SYMBOL ? ETH_TOKEN : USDC_TOKEN;

  const amountIn = useMemo(() => {
    try {
      return inputAmount ? parseUnits(inputAmount, currentToken.decimals) : 0n;
    } catch {
      return 0n;
    }
  }, [inputAmount, currentToken.decimals]);
  const currentBalance =
    bridgeToken === L1_NATIVE_SYMBOL ? ethBalance : usdcBalance;
  const hasInsufficientBalance = amountIn > currentBalance;
  const bridgeAmountUsd = amountIn > 0n ? Number(formatUnits(amountIn, currentToken.decimals)) : 0;
  const exceedsL2Limit = hasExceededL2Limit || (bridgeAmountUsd > 0 && wouldExceed(bridgeAmountUsd));

  const effectiveRecipient = (recipient || smartWallet || "") as Address;

  const handleBridge = useCallback(async () => {
    if (!smartWallet || amountIn === 0n) return;

    let success: boolean;
    if (bridgeToken === L1_NATIVE_SYMBOL) {
      success = await executeBridgeNative({
        amount: amountIn,
        recipient: effectiveRecipient,
        smartWallet,
      });
    } else {
      success = await executeBridge({
        amount: amountIn,
        recipient: effectiveRecipient,
        smartWallet,
      });
    }

    if (success) {
      recordSpending(bridgeAmountUsd);
      setInputAmount("");
    }
  }, [
    smartWallet,
    amountIn,
    bridgeToken,
    bridgeAmountUsd,
    effectiveRecipient,
    executeBridge,
    executeBridgeNative,
    recordSpending,
  ]);

  const getButtonText = () => {
    if (isPending) return "Bridging...";
    if (!isConnected) return "Connect Wallet";
    if (!smartWallet) return "Setup Smart Wallet First";
    if (!amountIn) return "Enter Amount";
    if (hasExceededL2Limit) return "L2 deposit limit reached ($1)";
    if (exceedsL2Limit) return `Exceeds $1 limit ($${remaining.toFixed(2)} left)`;
    if (hasInsufficientBalance) return "Insufficient Balance";
    return `Bridge ${bridgeToken} to L2`;
  };

  const isDisabled =
    isPending ||
    !isConnected ||
    !smartWallet ||
    !amountIn ||
    hasInsufficientBalance ||
    exceedsL2Limit;

  return (
    <div className="flex flex-col md:flex-row items-start gap-4 justify-center w-full relative z-10">
      {/* Left panel — inputs */}
      <div className="w-full md:max-w-md bg-surge-card/80 border border-surge-border/50 rounded-2xl p-4 space-y-3 shadow-xl shadow-black/20 hover-glow transition-all duration-[1000ms] ease-[cubic-bezier(0.16,1,0.3,1)]">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-white">Bridge</h2>
          <span className="text-xs text-gray-400">L1 &rarr; L2</span>
        </div>
        <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg px-3 py-2 text-xs text-yellow-400 flex items-center gap-1.5">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5 shrink-0" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
          </svg>
          We limit swaps and deposits to $1 on this experimental alpha
        </div>

        {/* Token Selector */}
        <div className="flex gap-2">
          {([L1_NATIVE_SYMBOL, "USDC"] as BridgeToken[]).map((t) => (
            <button
              key={t}
              onClick={() => {
                setBridgeToken(t);
                setInputAmount("");
              }}
              className={`flex-1 py-2 rounded-lg text-sm font-medium transition-colors ${
                bridgeToken === t
                  ? "bg-surge-primary text-white"
                  : "bg-surge-dark/50 text-gray-400 hover:text-white border border-surge-border/30"
              }`}
            >
              {t}
            </button>
          ))}
        </div>

        {/* Token Amount */}
        <TokenInput
          token={currentToken}
          amount={inputAmount}
          onAmountChange={setInputAmount}
          balance={currentBalance}
          label="Amount"
        />

        {/* Recipient (optional) */}
        <div className="space-y-1">
          <label className="text-xs text-gray-400">
            Recipient on L2 (optional)
          </label>
          <input
            type="text"
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
            placeholder={
              smartWallet ? `Default: ${smartWallet.slice(0, 10)}...` : "0x..."
            }
            className="w-full bg-surge-dark/50 border border-surge-border/30 rounded-lg px-3 py-2 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-surge-primary/50"
          />
        </div>

        {/* Bridge Button */}
        <button
          onClick={isConnected && !smartWallet ? onSetupWallet : handleBridge}
          disabled={isDisabled}
          className={`w-full py-3 rounded-xl font-semibold text-base transition-all duration-200 ${
            isDisabled
              ? "bg-surge-card/50 text-gray-500 cursor-not-allowed border border-surge-border/30"
              : "bg-gradient-to-r from-surge-primary to-surge-secondary text-white hover:shadow-lg hover:shadow-surge-primary/30 hover:scale-[1.02] active:scale-[0.98]"
          }`}
        >
          {isPending ? (
            <span className="flex items-center justify-center gap-2">
              <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              Bridging...
            </span>
          ) : (
            getButtonText()
          )}
        </button>
      </div>

      {/* Right panel — bridge details (shown when amount is entered) */}
      {amountIn > 0n && (
        <div className="w-full md:max-w-sm bg-surge-card/80 border border-surge-border/50 rounded-2xl p-4 space-y-3 shadow-xl shadow-black/20 animate-panel-in">
          <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-widest">Bridge Details</h3>

          {/* Flow Visualization */}
          <div className="flex items-center justify-center gap-3 py-3">
            <div className="flex items-center gap-2 bg-surge-dark/50 px-3 py-2 rounded-lg">
              <span className="text-xs text-gray-400">L1</span>
              <span className="text-sm text-white font-medium">
                {bridgeToken === L1_NATIVE_SYMBOL ? "Send" : "Lock"}
              </span>
            </div>
            <div className="text-surge-primary">&rarr;</div>
            <div className="flex items-center gap-2 bg-surge-dark/50 px-3 py-2 rounded-lg">
              <span className="text-xs text-gray-400">L2</span>
              <span className="text-sm text-white font-medium">
                {bridgeToken === L1_NATIVE_SYMBOL ? "Receive" : "Mint"}
              </span>
            </div>
          </div>

          {/* Transfer Summary */}
          <div className="bg-surge-dark/30 rounded-lg p-3 space-y-1">
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">You send</span>
              <span className="text-white">
                {formatUnits(amountIn, currentToken.decimals)} {bridgeToken}
              </span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">You receive</span>
              <span className="text-white">
                {formatUnits(amountIn, currentToken.decimals)} {bridgeToken} on L2
              </span>
            </div>
            <div className="flex justify-between text-sm mt-1 pt-1 border-t border-surge-border/30">
              <span className="text-gray-400">Recipient</span>
              <span className="text-white text-xs font-mono">
                {effectiveRecipient
                  ? `${effectiveRecipient.slice(0, 6)}...${effectiveRecipient.slice(-4)}`
                  : "—"}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
