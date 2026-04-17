import { Address } from 'viem';

export interface Token {
  symbol: string;
  name: string;
  decimals: number;
  address: Address | null; // null for ETH
  logo: string;
}

export interface UserOp {
  target: Address;
  value: bigint;
  data: `0x${string}`;
}

export interface SwapParams {
  tokenIn: Token;
  tokenOut: Token;
  amountIn: bigint;
  minAmountOut: bigint;
  recipient: Address;
}

export interface SwapQuote {
  amountOut: bigint;
  priceImpact: number;
  fee: bigint;
  rate: number;
  insufficientLiquidity: boolean;
}

export interface DexReserves {
  ethReserve: bigint;
  tokenReserve: bigint;
}

export type SwapDirection = 'ETH_TO_USDC' | 'USDC_TO_ETH';

/// Which DEX actually runs the swap. Drives which chain the user transacts on,
/// which wallet flow is used (UserOp vs direct EOA), and how proceeds settle.
///   L2_DEX — L1→L2→L1 swap via L2 SimpleDEX (existing, UserOp-signed on L1).
///   L1_DEX — L2→L1→L2 swap via L1 DEX / Uniswap V2 (direct EOA tx on L2).
export type SwapVenue = 'L2_DEX' | 'L1_DEX';

export type BridgeDirection = 'L1_TO_L2' | 'L2_TO_L1';

export type ActiveTab = 'swap' | 'bridge' | 'liquidity';

export type TxOverlayPhase =
  | 'idle'
  // L1→L2→L1 UserOp lifecycle
  | 'signing'
  | 'sequencing'
  | 'proposing'
  | 'proving'
  // L2→L1→L2 direct-tx lifecycle
  | 'simulating'
  | 'submitting'
  | 'included'
  // Terminal
  | 'complete'
  | 'rejected';

export interface TxOverlayState {
  phase: TxOverlayPhase;
  txHash?: string;
  errorMessage?: string;
}

export type AccountMode = 'safe' | 'ambire';
