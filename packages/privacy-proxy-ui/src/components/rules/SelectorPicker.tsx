import { useEffect, useState } from "react";
import { useSyntheticSelectors } from "../../hooks/useSyntheticSelectors";
import { isSelector } from "../../lib/format";
import { COMMON_SELECTORS, findCommonSelector } from "../../lib/selectors";

interface SelectorPickerProps {
  value: string;
  onChange: (v: string) => void;
  /** Disable editing (edit mode). */
  disabled?: boolean;
}

const CUSTOM_SENTINEL = "__custom__";

/**
 * Lets the admin pick either a synthetic method name (capability 10) or a
 * raw 4-byte hex. In hex mode, a dropdown of common selectors is offered;
 * the top "Custom…" option reveals a free-form text input. The server
 * normalizes the method name to the synthetic selector on write.
 */
export function SelectorPicker({ value, onChange, disabled }: SelectorPickerProps) {
  const { data: selectors } = useSyntheticSelectors();
  // Determine current input mode based on what's in `value`. Method names
  // start with "eth_"; everything else is treated as raw hex.
  const isMethod = !!value && !value.startsWith("0x");
  const [mode, setMode] = useState<"method" | "hex">(isMethod ? "method" : "hex");

  useEffect(() => {
    if (isMethod && mode !== "method") setMode("method");
  }, [isMethod, mode]);

  const valid = mode === "method" || !value || isSelector(value);

  // In hex mode, surface a dropdown of common selectors. If `value` matches
  // one of them, the dropdown reflects it; otherwise we're in "custom" mode
  // and the free-text input is shown.
  const matchedCommon = findCommonSelector(value);
  const isCustom = mode === "hex" && !matchedCommon;

  return (
    <div>
      <div className="mb-2 inline-flex rounded-lg border border-surge-border bg-surge-card-hover/40 p-0.5 text-xs">
        {(["method", "hex"] as const).map((m) => (
          <button
            key={m}
            type="button"
            disabled={disabled}
            onClick={() => {
              setMode(m);
              // Clear value when switching modes so we don't ship hex tagged as method.
              onChange("");
            }}
            className={`rounded-md px-2.5 py-1 transition ${
              mode === m
                ? "bg-surge-card text-surge-text shadow-sm"
                : "text-surge-muted hover:text-surge-text"
            }`}
          >
            {m === "method" ? "RPC method" : "4-byte hex"}
          </button>
        ))}
      </div>

      {mode === "method" ? (
        <select
          value={value}
          disabled={disabled}
          onChange={(e) => onChange(e.target.value)}
          className="w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 text-sm outline-none focus:border-surge-secondary"
        >
          <option value="">Select an RPC method…</option>
          {selectors?.map((s) => (
            <option key={s.method} value={s.method}>
              {s.method} ({s.selector})
            </option>
          ))}
        </select>
      ) : (
        <div className="space-y-2">
          <select
            value={matchedCommon ? matchedCommon.selector : CUSTOM_SENTINEL}
            disabled={disabled}
            onChange={(e) => {
              const v = e.target.value;
              if (v === CUSTOM_SENTINEL) {
                // Switching to custom — clear so the input renders empty.
                onChange("");
              } else {
                onChange(v);
              }
            }}
            className="w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 text-sm outline-none focus:border-surge-secondary"
          >
            <option value={CUSTOM_SENTINEL}>Custom… (enter a 4-byte selector)</option>
            <option disabled>──────────</option>
            {COMMON_SELECTORS.map((s) => (
              <option key={s.selector} value={s.selector}>
                {s.signature} · {s.selector}
                {s.tag ? ` (${s.tag})` : ""}
              </option>
            ))}
          </select>

          {isCustom ? (
            <input
              type="text"
              value={value}
              disabled={disabled}
              onChange={(e) => onChange(e.target.value)}
              placeholder="0xa9059cbb"
              className="w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 font-mono text-sm outline-none focus:border-surge-secondary"
              autoFocus
            />
          ) : null}
        </div>
      )}
      {!valid ? (
        <p className="mt-1 text-xs text-red-600">
          Selector must be 0x + 8 hex chars.
        </p>
      ) : null}
    </div>
  );
}
