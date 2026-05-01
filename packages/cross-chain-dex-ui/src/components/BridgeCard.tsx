import { useState, useCallback, useMemo, useEffect } from "react";
import { parseUnits, formatUnits, Address, isAddress } from "viem";
import { useAccount } from "wagmi";
import { TokenInput } from "./TokenInput";
import { useSmartWallet } from "../context/SmartWalletContext";
import { useSharedTokenBalances } from "../context/SmartWalletContext";
import { useTokenBalances } from "../hooks/useTokenBalances";
import { useUserOp } from "../hooks/useUserOp";
import { useBridgeOutEoa } from "../hooks/useBridgeOutEoa";
import { useSpendingLimit } from "../hooks/useSpendingLimit";
import { ETH_TOKEN, L1_NATIVE_SYMBOL, L1_CHAIN_NAME } from "../lib/constants";
import { l1PublicClient, l2PublicClient } from "../lib/config";
import { DisclaimerModal } from "./DisclaimerModal";
import { useDisclaimer } from "../hooks/useDisclaimer";
import { WarningBanner } from "./WarningBanner";

interface BridgeCardProps {
  network: 'l1' | 'l2';
  onSetupWallet: () => void;
  onFundWallet?: () => void;
}

export function BridgeCard({ network, onSetupWallet, onFundWallet }: BridgeCardProps) {
  const isL1 = network === 'l1';
  const { smartWallet, isConnected, l2WalletExists, accountMode, setSelectedNetwork } = useSmartWallet();
  const { address: eoaAddress } = useAccount();

  // L1 deposit uses smart wallet's L1 balance; L2 withdraw uses EOA's L2 balance.
  const smartWalletBalances = useSharedTokenBalances();
  const eoaL2Balances = useTokenBalances(eoaAddress ?? null, 'l2');
  const { ethBalance } = isL1 ? smartWalletBalances : eoaL2Balances;

  const { executeBridgeNative, isPending: isUserOpPending } = useUserOp(accountMode);
  const { initiate: initiateBridgeOut, isPending: isBridgeOutPending } = useBridgeOutEoa();
  const isPending = isL1 ? isUserOpPending : isBridgeOutPending;

  const { hasExceededL2Limit, wouldExceed, recordSpending, remaining } = useSpendingLimit(smartWallet);
  const { isDisclaimerOpen, requireDisclaimer, onAccept, onCancel } = useDisclaimer();

  // Bridge currently only supports the native xDAI in both directions.
  // USDC bridge-out (L2 -> L1) is not implemented at the contract layer.
  const [inputAmount, setInputAmount] = useState("");
  const [recipient, setRecipient] = useState("");

  const amountIn = useMemo(() => {
    try {
      return inputAmount ? parseUnits(inputAmount, ETH_TOKEN.decimals) : 0n;
    } catch {
      return 0n;
    }
  }, [inputAmount]);

  const hasInsufficientBalance = amountIn > ethBalance;
  const bridgeAmountUsd = amountIn > 0n ? Number(formatUnits(amountIn, ETH_TOKEN.decimals)) : 0;
  const exceedsL2Limit = isL1 && (hasExceededL2Limit || (bridgeAmountUsd > 0 && wouldExceed(bridgeAmountUsd)));

  // Default recipient: the connected EOA on the destination chain.
  const defaultRecipient = eoaAddress || "";
  const effectiveRecipient = (recipient || defaultRecipient) as Address;

  // If the user types a custom recipient, check whether it's a contract on the
  // destination chain — bridging to a non-receiving contract can strand funds.
  const [recipientIsContract, setRecipientIsContract] = useState(false);
  useEffect(() => {
    setRecipientIsContract(false);
    const trimmed = recipient.trim();
    if (!trimmed || !isAddress(trimmed)) return;
    let cancelled = false;
    const destClient = isL1 ? l2PublicClient : l1PublicClient;
    destClient
      .getCode({ address: trimmed as Address })
      .then((code) => {
        if (!cancelled) setRecipientIsContract(!!code && code !== '0x');
      })
      .catch(() => { /* ignore network blips */ });
    return () => { cancelled = true; };
  }, [recipient, isL1]);

  const handleSubmit = useCallback(async () => {
    if (amountIn === 0n) return;

    let success = false;
    if (isL1) {
      if (!smartWallet) return;
      success = await executeBridgeNative({ amount: amountIn, recipient: effectiveRecipient, smartWallet });
      if (success) recordSpending(bridgeAmountUsd);
    } else {
      if (!eoaAddress) return;
      success = await initiateBridgeOut({ amount: amountIn, recipient: effectiveRecipient });
    }

    if (success) setInputAmount("");
  }, [
    isL1,
    smartWallet,
    eoaAddress,
    amountIn,
    effectiveRecipient,
    bridgeAmountUsd,
    executeBridgeNative,
    initiateBridgeOut,
    recordSpending,
  ]);

  const getButtonText = () => {
    if (isPending) return "Bridging...";
    if (!isConnected) return "Connect Wallet";
    if (isL1 && !smartWallet) return "Setup Smart Wallet First";
    if (isL1 && !l2WalletExists) return "Create L2 wallet first";
    if (!amountIn) return "Enter Amount";
    if (isL1 && hasExceededL2Limit) return "L2 deposit limit reached ($1)";
    if (isL1 && exceedsL2Limit) return `Exceeds $1 limit ($${remaining.toFixed(2)} left)`;
    if (hasInsufficientBalance) return "Insufficient Balance";
    return isL1 ? `Bridge ${L1_NATIVE_SYMBOL} to L2` : `Withdraw ${L1_NATIVE_SYMBOL} to L1`;
  };

  const needsL2WalletSetup = isL1 && !!smartWallet && !l2WalletExists && !!onFundWallet;

  const isDisabled =
    isPending ||
    !isConnected ||
    (isL1 && !smartWallet) ||
    (isL1 && !l2WalletExists) ||
    !amountIn ||
    hasInsufficientBalance ||
    exceedsL2Limit;

  const venueBadge = isL1
    ? { label: 'Via Smart Account', cls: 'bg-surge-secondary/15 border-surge-secondary/50 text-surge-primary' }
    : { label: 'Via EOA', cls: 'bg-surge-secondary/20 border-surge-secondary/60 text-surge-primary' };

  return (
    <div className="flex flex-col md:flex-row items-start gap-4 justify-center w-full relative z-10">
      {/* Left panel — inputs */}
      <div className="w-full md:max-w-md glass-card rounded-2xl p-4 space-y-3 hover-glow transition-all duration-[1000ms] ease-[cubic-bezier(0.16,1,0.3,1)]">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="inline-block w-1 h-5 rounded-full bg-surge-secondary" />
            <h2 className="text-lg font-semibold text-surge-text">
              {isL1 ? 'Deposit' : 'Withdraw'}
            </h2>
          </div>
          <span className={`text-[11px] font-semibold uppercase tracking-wider px-2.5 py-1 rounded-full border ${venueBadge.cls}`}>
            {venueBadge.label}
          </span>
        </div>
        <WarningBanner />

        {/* Cross-network notice */}
        <button
          type="button"
          onClick={() => setSelectedNetwork(isL1 ? 'l2' : 'l1')}
          className="w-full text-left bg-surge-card-hover border border-surge-border rounded-lg px-3 py-2 text-xs text-surge-muted hover:bg-surge-secondary/10 hover:border-surge-secondary/40 hover:text-surge-primary transition-colors"
        >
          {isL1
            ? <>Need to withdraw L2 → L1? <span className="font-semibold text-surge-primary">Switch to Surge L2 →</span></>
            : <>Need to deposit L1 → L2? <span className="font-semibold text-surge-primary">Switch to {L1_CHAIN_NAME} →</span></>}
        </button>

        {/* Token Amount */}
        <TokenInput
          token={ETH_TOKEN}
          amount={inputAmount}
          onAmountChange={setInputAmount}
          balance={ethBalance}
          label="Amount"
        />

        {/* Recipient (optional) */}
        <div className="space-y-1">
          <label className="text-xs text-surge-muted">
            {isL1 ? "Recipient on L2 (optional)" : "Recipient on L1 (optional)"}
          </label>
          <input
            type="text"
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
            placeholder={
              defaultRecipient ? `Default: ${defaultRecipient.slice(0, 10)}...` : "0x..."
            }
            className="w-full bg-surge-card-hover border border-surge-border rounded-lg px-3 py-2 text-sm text-surge-text placeholder-surge-muted/60 focus:outline-none focus:border-surge-secondary"
          />
          {recipientIsContract && (
            <div className="bg-red-50 border border-red-300 rounded-lg px-3 py-2 text-xs text-red-700 flex items-start gap-1.5 mt-1">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5 shrink-0 mt-0.5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
              <span>This recipient is a contract on {isL1 ? "L2" : "L1"}. Funds may be stranded if it can't accept native transfers — double-check before sending.</span>
            </div>
          )}
        </div>

        {/* Bridge Button */}
        <button
          onClick={
            isL1 && isConnected && !smartWallet
              ? onSetupWallet
              : needsL2WalletSetup
                ? onFundWallet
                : () => requireDisclaimer(handleSubmit)
          }
          disabled={!needsL2WalletSetup && isDisabled}
          className={`w-full py-3 rounded-xl font-semibold text-base transition-all duration-200 ${
            isDisabled
              ? "bg-surge-card-hover text-surge-muted cursor-not-allowed border border-surge-border"
              : "bg-surge-primary text-white hover:bg-surge-secondary hover:shadow-md hover:shadow-surge-secondary/25 active:scale-[0.98]"
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
        <div className="w-full md:max-w-sm glass-card rounded-2xl p-4 space-y-3 animate-panel-in">
          <div className="flex items-center gap-2">
            <span className="inline-block w-1 h-4 rounded-full bg-surge-secondary" />
            <h3 className="text-xs font-semibold text-surge-muted uppercase tracking-widest">Bridge Details</h3>
          </div>

          {/* Flow Visualization */}
          <div className="flex items-center justify-center gap-3 py-3">
            <div className="flex items-center gap-2 bg-surge-card-hover border border-surge-border px-3 py-2 rounded-lg">
              <span className="text-xs text-surge-muted">{isL1 ? "L1" : "L2"}</span>
              <span className="text-sm text-surge-text font-medium">Send</span>
            </div>
            <div className="text-surge-secondary">&rarr;</div>
            <div className="flex items-center gap-2 bg-surge-card-hover border border-surge-border px-3 py-2 rounded-lg">
              <span className="text-xs text-surge-muted">{isL1 ? "L2" : "L1"}</span>
              <span className="text-sm text-surge-text font-medium">Receive</span>
            </div>
          </div>

          {/* Transfer Summary */}
          <div className="bg-surge-card-hover border border-surge-border rounded-lg p-3 space-y-1">
            <div className="flex justify-between text-sm">
              <span className="text-surge-muted">You send</span>
              <span className="text-surge-text">
                {formatUnits(amountIn, ETH_TOKEN.decimals)} {L1_NATIVE_SYMBOL}
              </span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-surge-muted">You receive</span>
              <span className="text-surge-text">
                {formatUnits(amountIn, ETH_TOKEN.decimals)} {L1_NATIVE_SYMBOL} on {isL1 ? "L2" : "L1"}
              </span>
            </div>
            <div className="flex justify-between text-sm mt-1 pt-1 border-t border-surge-border">
              <span className="text-surge-muted">Recipient</span>
              <span className="text-surge-text text-xs font-mono">
                {effectiveRecipient
                  ? `${effectiveRecipient.slice(0, 6)}...${effectiveRecipient.slice(-4)}`
                  : "—"}
              </span>
            </div>
          </div>
        </div>
      )}
      <DisclaimerModal isOpen={isDisclaimerOpen} onAccept={onAccept} onCancel={onCancel} />
    </div>
  );
}
