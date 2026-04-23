import { SwapDirection, SwapVenue } from '../types';
import { ETH_TOKEN, USDC_TOKEN, L1_DEX_NAME } from '../lib/constants';

interface SwapPathProps {
  direction: SwapDirection;
  venue: SwapVenue;
}

export function SwapPath({ direction, venue }: SwapPathProps) {
  const inputToken = direction === 'ETH_TO_USDC' ? ETH_TOKEN : USDC_TOKEN;
  const outputToken = direction === 'ETH_TO_USDC' ? USDC_TOKEN : ETH_TOKEN;

  const isL1Dex = venue === 'L1_DEX';

  // Chain labels: where the user starts, where the DEX runs, where they end up.
  const leftChain = isL1Dex ? 'L2' : 'L1';
  const middleChain = isL1Dex ? 'L1' : 'L2';
  const rightChain = isL1Dex ? 'L2' : 'L1';

  // Colours: L1→L2→L1 uses emerald→cyan (green-leaning), L2→L1→L2 uses cyan→amber
  // so the audience can read direction at a glance.
  const leftColor = isL1Dex ? 'text-cyan-400' : 'text-emerald-500';
  const middleColor = isL1Dex ? 'text-amber-400' : 'text-cyan-500';
  const rightColor = isL1Dex ? 'text-cyan-400' : 'text-emerald-500';

  const firstGradient = isL1Dex
    ? 'from-cyan-500/50 to-amber-500/50'
    : 'from-emerald-500/50 to-cyan-500/50';
  const firstArrowColor = isL1Dex ? 'border-l-amber-500/70' : 'border-l-cyan-500/70';
  const secondGradient = isL1Dex
    ? 'from-amber-500/50 to-cyan-500/50'
    : 'from-cyan-500/50 to-emerald-500/50';
  const secondArrowColor = isL1Dex ? 'border-l-cyan-500/70' : 'border-l-emerald-500/70';

  const dexBg = isL1Dex
    ? 'bg-gradient-to-br from-amber-500/20 to-cyan-500/20 border border-amber-500/30'
    : 'bg-gradient-to-br from-emerald-500/20 to-cyan-500/20 border border-emerald-500/30';
  const dexIconColor = isL1Dex ? 'text-amber-400' : 'text-emerald-400';
  const dexLabel = isL1Dex ? L1_DEX_NAME : 'DEX';

  return (
    <div className="bg-surge-dark/50 rounded-xl p-4 border border-surge-border/30">
      <div className="text-xs text-gray-500 mb-3 text-center">Swap Route</div>

      {/* Icons Row — all vertically centered */}
      <div className="flex items-center">
        {/* Input Token */}
        <div className="flex flex-col items-center flex-shrink-0 w-10">
          <div className="w-10 h-10 rounded-full bg-surge-card border border-surge-border/50 flex items-center justify-center">
            <img src={inputToken.logo} alt={inputToken.symbol} className="w-6 h-6" />
          </div>
          <span className="text-[10px] text-gray-400 mt-1">{inputToken.symbol}</span>
        </div>

        {/* First bridge connector */}
        <div className="flex-1 flex items-center mx-1 relative self-start mt-5">
          <div className={`flex-1 h-[2px] bg-gradient-to-r ${firstGradient}`} />
          <div className={`w-0 h-0 border-t-[5px] border-t-transparent border-b-[5px] border-b-transparent border-l-[7px] ${firstArrowColor}`} />
          <span className="absolute -bottom-4 left-1/2 -translate-x-1/2 text-[9px] text-gray-500">bridge</span>
        </div>

        {/* DEX Icon */}
        <div className="flex flex-col items-center flex-shrink-0 w-10">
          <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${dexBg}`}>
            <svg className={`w-5 h-5 ${dexIconColor}`} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M7 16V4M7 4L3 8M7 4L11 8M17 8V20M17 20L21 16M17 20L13 16" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </div>
          <span className="text-[10px] text-gray-400 mt-1 whitespace-nowrap">{dexLabel}</span>
        </div>

        {/* Second bridge connector */}
        <div className="flex-1 flex items-center mx-1 relative self-start mt-5">
          <div className={`flex-1 h-[2px] bg-gradient-to-r ${secondGradient}`} />
          <div className={`w-0 h-0 border-t-[5px] border-t-transparent border-b-[5px] border-b-transparent border-l-[7px] ${secondArrowColor}`} />
          <span className="absolute -bottom-4 left-1/2 -translate-x-1/2 text-[9px] text-gray-500">bridge</span>
        </div>

        {/* Output Token */}
        <div className="flex flex-col items-center flex-shrink-0 w-10">
          <div className="w-10 h-10 rounded-full bg-surge-card border border-surge-border/50 flex items-center justify-center">
            <img src={outputToken.logo} alt={outputToken.symbol} className="w-6 h-6" />
          </div>
          <span className="text-[10px] text-gray-400 mt-1">{outputToken.symbol}</span>
        </div>
      </div>

      {/* Chain labels row */}
      <div className="flex items-start mt-4">
        <div className="w-10 flex justify-center flex-shrink-0">
          <span className={`text-[11px] font-bold ${leftColor}`}>{leftChain}</span>
        </div>
        <div className="flex-1" />
        <div className="w-10 flex justify-center flex-shrink-0">
          <span className={`text-[11px] font-bold ${middleColor}`}>{middleChain}</span>
        </div>
        <div className="flex-1" />
        <div className="w-10 flex justify-center flex-shrink-0">
          <span className={`text-[11px] font-bold ${rightColor}`}>{rightChain}</span>
        </div>
      </div>
    </div>
  );
}
