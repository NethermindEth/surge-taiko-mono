import { TxOverlayPhase, TxOverlayState } from '../types';
import { EXPLORER_URL } from '../lib/constants';

// ─── Step definitions ────────────────────────────────────────────────────────

type ActivePhase = Exclude<TxOverlayPhase, 'idle' | 'rejected'>;

interface StepDef {
  phase: ActivePhase;
  label: string;
  color: string;
}

const STEPS: StepDef[] = [
  { phase: 'signing',    label: 'Signing',             color: '#fbbf24' },
  { phase: 'sequencing', label: 'Sequencing',           color: '#60a5fa' },
  { phase: 'proving',    label: 'Generating ZK Proof',  color: '#a78bfa' },
  { phase: 'proposing',  label: 'Submitting Block',      color: '#34d399' },
  { phase: 'complete',   label: 'Execution Complete',   color: '#10b981' },
];

const PHASE_TO_IDX: Partial<Record<TxOverlayPhase, number>> = {
  signing: 0, sequencing: 1, proving: 2, proposing: 3, complete: 4,
};

const ITEM_H = 88;   // px per slot
const WINDOW_H = ITEM_H * 3; // 3 slots visible

// ─── Icon sub-components ─────────────────────────────────────────────────────

function SpinnerRing({
  color,
  spinning,
  children,
}: {
  color: string;
  spinning: boolean;
  children: React.ReactNode;
}) {
  return (
    <div style={{ position: 'relative', width: 56, height: 56, flexShrink: 0 }}>
      {/* Rotating ring */}
      <div
        style={{
          position: 'absolute',
          inset: 0,
          borderRadius: '50%',
          border: `2.5px solid ${color}28`,
          borderTop: `2.5px solid ${color}`,
          animation: spinning ? 'spin 0.8s linear infinite' : 'none',
        }}
      />
      {/* Centered icon */}
      <div
        style={{
          position: 'absolute',
          inset: 11,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        {children}
      </div>
    </div>
  );
}

function SigningIcon({ color, spinning }: { color: string; spinning: boolean }) {
  return (
    <SpinnerRing color={color} spinning={spinning}>
      <svg viewBox="0 0 24 24" fill="none" stroke={color} strokeWidth="2"
        strokeLinecap="round" strokeLinejoin="round" style={{ width: 22, height: 22 }}>
        <path d="M12 20h9" />
        <path d="M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4 9.5-9.5z" />
      </svg>
    </SpinnerRing>
  );
}

function SequencingIcon({ color, spinning }: { color: string; spinning: boolean }) {
  return (
    <SpinnerRing color={color} spinning={spinning}>
      <svg viewBox="0 0 24 24" fill="none" stroke={color} strokeWidth="2"
        strokeLinecap="round" strokeLinejoin="round" style={{ width: 22, height: 22 }}>
        <polygon points="12 2 2 7 12 12 22 7 12 2" />
        <polyline points="2 17 12 22 22 17" />
        <polyline points="2 12 12 17 22 12" />
      </svg>
    </SpinnerRing>
  );
}

function ProposingIcon({ color, spinning }: { color: string; spinning: boolean }) {
  return (
    <SpinnerRing color={color} spinning={spinning}>
      <svg viewBox="0 0 24 24" fill="none" stroke={color} strokeWidth="2"
        strokeLinecap="round" strokeLinejoin="round" style={{ width: 22, height: 22 }}>
        <path d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
        <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
        <line x1="12" y1="22.08" x2="12" y2="12" />
      </svg>
    </SpinnerRing>
  );
}

function ProvingIcon({ color, spinning }: { color: string; spinning: boolean }) {
  return (
    <SpinnerRing color={color} spinning={spinning}>
      <svg viewBox="0 0 24 24" fill={color} style={{ width: 22, height: 22 }}>
        <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
      </svg>
    </SpinnerRing>
  );
}

function CompleteCircle({ color }: { color: string }) {
  return (
    <div
      style={{
        width: 56,
        height: 56,
        borderRadius: '50%',
        border: `2.5px solid ${color}`,
        background: `${color}18`,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        flexShrink: 0,
      }}
    >
      <svg viewBox="0 0 24 24" fill="none" stroke={color} strokeWidth="2.5"
        strokeLinecap="round" strokeLinejoin="round" style={{ width: 26, height: 26 }}>
        <polyline points="20 6 9 17 4 12" />
      </svg>
    </div>
  );
}

// ─── Step row rendered inside the wheel ──────────────────────────────────────

const ICON_MAP: Record<ActivePhase, React.ComponentType<{ color: string; spinning: boolean }>> = {
  signing:    SigningIcon,
  sequencing: SequencingIcon,
  proposing:  ProposingIcon,
  proving:    ProvingIcon,
  complete:   () => null, // replaced by CompleteCircle
};

function StepRow({ step, offset, isActive }: { step: StepDef; offset: number; isActive: boolean }) {
  const { color, label, phase } = step;
  const absOffset = Math.abs(offset);
  const opacity   = absOffset === 0 ? 1 : absOffset === 1 ? 0.28 : 0;
  const scale     = absOffset === 0 ? 1 : 0.86;
  const spinning  = offset === 0 && isActive;

  const Icon = ICON_MAP[phase];

  return (
    <div
      style={{
        height: ITEM_H,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 20,
        opacity,
        transform: `scale(${scale})`,
        transition: 'opacity 0.45s ease, transform 0.45s ease',
      }}
    >
      {phase === 'complete' ? (
        <CompleteCircle color={color} />
      ) : (
        <Icon color={color} spinning={spinning} />
      )}

      <span
        style={{
          fontSize: absOffset === 0 ? 21 : 17,
          fontWeight: absOffset === 0 ? 600 : 400,
          color: absOffset === 0 ? '#f1f5f9' : '#64748b',
          whiteSpace: 'nowrap',
          letterSpacing: '-0.01em',
          transition: 'font-size 0.45s ease, color 0.45s ease, font-weight 0.45s ease',
          fontFamily: 'inherit',
        }}
      >
        {label}
      </span>
    </div>
  );
}

// ─── Main overlay ─────────────────────────────────────────────────────────────

interface TxStatusOverlayProps {
  state: TxOverlayState;
  onClose: () => void;
}

export function TxStatusOverlay({ state, onClose }: TxStatusOverlayProps) {
  if (state.phase === 'idle') return null;

  const currentIdx = PHASE_TO_IDX[state.phase] ?? 0;
  const translateY = ITEM_H * (1 - currentIdx);
  const isInProgress = !['complete', 'rejected'].includes(state.phase);

  // ── Rejected view ────────────────────────────────────────────────────────
  if (state.phase === 'rejected') {
    return (
      <Backdrop>
        <Card>
          <div style={{ padding: '40px 32px 32px', textAlign: 'center' }}>
            {/* Red X */}
            <div style={{ display: 'flex', justifyContent: 'center', marginBottom: 20 }}>
              <div style={{
                width: 72, height: 72, borderRadius: '50%',
                border: '2.5px solid #ef4444',
                background: '#ef444414',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
              }}>
                <svg viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2.5"
                  strokeLinecap="round" style={{ width: 34, height: 34 }}>
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </div>
            </div>
            <p style={{ color: '#f1f5f9', fontWeight: 600, fontSize: 18, marginBottom: 8, fontFamily: 'inherit' }}>
              Transaction Failed
            </p>
            {state.errorMessage && (
              <p style={{ color: '#94a3b8', fontSize: 14, marginBottom: 24, fontFamily: 'inherit', lineHeight: 1.5 }}>
                {state.errorMessage}
              </p>
            )}
            <button
              onClick={onClose}
              className="w-full py-3 rounded-xl text-sm font-medium bg-surge-card border border-surge-border/50 text-white hover:bg-surge-border/50 transition-colors"
            >
              Close
            </button>
          </div>
        </Card>
      </Backdrop>
    );
  }

  // ── Slot-wheel view (in-progress + complete) ─────────────────────────────
  return (
    <Backdrop>
      <Card>
        <div style={{ padding: '36px 32px 32px' }}>
          {/* Header label */}
          <p style={{
            textAlign: 'center',
            fontSize: 11,
            fontWeight: 500,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
            color: '#475569',
            marginBottom: 28,
            fontFamily: 'inherit',
          }}>
            {isInProgress ? 'Transaction in Progress' : 'Transaction Complete'}
          </p>

          {/* ── Slot wheel ── */}
          <div
            style={{
              height: WINDOW_H,
              overflow: 'hidden',
              position: 'relative',
              WebkitMaskImage:
                'linear-gradient(to bottom, transparent 0%, black 26%, black 74%, transparent 100%)',
              maskImage:
                'linear-gradient(to bottom, transparent 0%, black 26%, black 74%, transparent 100%)',
            }}
          >
            <div
              style={{
                transform: `translateY(${translateY}px)`,
                transition: 'transform 0.55s cubic-bezier(0.4, 0, 0.2, 1)',
              }}
            >
              {STEPS.map((step, i) => (
                <StepRow
                  key={step.phase}
                  step={step}
                  offset={i - currentIdx}
                  isActive={isInProgress}
                />
              ))}
            </div>
          </div>

          {/* ── Complete buttons ── */}
          {state.phase === 'complete' && (
            <div className="animate-fade-up" style={{ display: 'flex', gap: 10, marginTop: 28 }}>
              {state.txHash && (
                <a
                  href={`${EXPLORER_URL}/tx/${state.txHash}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex-1 py-3 rounded-xl text-sm font-medium text-center bg-surge-card border border-surge-border/50 text-white hover:bg-surge-border/40 transition-colors"
                >
                  View on Explorer
                </a>
              )}
              <button
                onClick={onClose}
                className="flex-1 py-3 rounded-xl text-sm font-medium bg-gradient-to-r from-surge-primary to-surge-secondary text-white hover:shadow-lg hover:shadow-surge-primary/30 hover:scale-[1.02] active:scale-[0.98] transition-all"
              >
                Close
              </button>
            </div>
          )}
        </div>
      </Card>
    </Backdrop>
  );
}

// ─── Layout helpers ───────────────────────────────────────────────────────────

function Backdrop({ children }: { children: React.ReactNode }) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      style={{ backgroundColor: 'rgba(3, 5, 8, 0.92)' }}
    >
      {children}
    </div>
  );
}

function Card({ children }: { children: React.ReactNode }) {
  return (
    <div
      className="w-full mx-4"
      style={{
        maxWidth: 400,
        background: 'linear-gradient(160deg, #0d1929 0%, #070c16 55%, #090e1c 100%)',
        border: '1px solid rgba(255, 255, 255, 0.07)',
        borderRadius: 24,
        boxShadow:
          'inset 0 1px 0 rgba(255, 255, 255, 0.05), 0 32px 64px rgba(0, 0, 0, 0.8), 0 0 60px rgba(16, 185, 129, 0.06)',
      }}
    >
      {children}
    </div>
  );
}
