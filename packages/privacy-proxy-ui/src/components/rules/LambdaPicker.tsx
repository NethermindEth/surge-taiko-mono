import { useMemo, useState } from "react";
import { useLambdas } from "../../hooks/useLambdas";
import { useSyntheticSelectors } from "../../hooks/useSyntheticSelectors";
import type { LambdaView, RoleName } from "../../types/api";
import { Badge } from "../common/Badge";
import { LambdaBuilder } from "./LambdaBuilder";

interface LambdaPickerProps {
  role: RoleName;
  /**
   * The access rule's selector value as held in the parent form. May be a
   * 4-byte hex (e.g. "0xa9059cbb") or a synthetic method name (e.g.
   * "eth_getBalance"). The picker normalizes for downstream use.
   */
  selector: string;
  value: number | null;
  onChange: (id: number | null) => void;
}

export function LambdaPicker({ role, selector, value, onChange }: LambdaPickerProps) {
  const lambdasQ = useLambdas();
  const syntheticQ = useSyntheticSelectors();
  const [creating, setCreating] = useState(false);

  const hexSelector = useMemo(() => {
    if (!selector) return "";
    if (selector.startsWith("0x")) return selector.toLowerCase();
    const found = (syntheticQ.data ?? []).find((s) => s.method === selector);
    return found ? found.selector.toLowerCase() : "";
  }, [selector, syntheticQ.data]);

  const lambdasForRole = useMemo(() => {
    return lambdasQ.data?.find((g) => g.role === role)?.lambdas ?? [];
  }, [lambdasQ.data, role]);

  const selected: LambdaView | undefined = useMemo(() => {
    if (value === null) return undefined;
    return lambdasForRole.find((l) => l.id === value);
  }, [lambdasForRole, value]);

  const matchesSelector = (l: LambdaView): boolean =>
    !hexSelector ||
    l.rules.some((r) => r.selector.toLowerCase() === hexSelector);

  if (creating) {
    if (!hexSelector) {
      return (
        <div className="rounded-xl border border-dashed border-amber-400 bg-amber-50/40 p-3 text-xs text-amber-900">
          Pick a selector first — the lambda's rules are scoped to it.
          <div className="mt-2">
            <button
              type="button"
              onClick={() => setCreating(false)}
              className="rounded-md border border-amber-400 px-2 py-1 text-xs font-medium hover:bg-amber-100"
            >
              Cancel
            </button>
          </div>
        </div>
      );
    }
    return (
      <LambdaBuilder
        role={role}
        selector={hexSelector}
        onCancel={() => setCreating(false)}
        onCreated={(l) => {
          onChange(l.id);
          setCreating(false);
        }}
      />
    );
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2">
        <select
          value={value ?? ""}
          onChange={(e) => onChange(e.target.value === "" ? null : Number(e.target.value))}
          className="flex-1 rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
        >
          <option value="">(no lambda)</option>
          {lambdasForRole.map((l) => (
            <option key={l.id} value={l.id}>
              {l.name}
              {!matchesSelector(l) ? " — no rules for this selector" : ""}
            </option>
          ))}
        </select>
        <button
          type="button"
          onClick={() => setCreating(true)}
          className="rounded-md border border-surge-secondary bg-surge-secondary/10 px-2 py-1.5 text-xs font-semibold text-surge-primary hover:bg-surge-secondary/20"
        >
          + New lambda
        </button>
      </div>

      {selected && !matchesSelector(selected) ? (
        <p className="text-[11px] text-amber-700">
          Heads up: lambda{" "}
          <span className="font-mono">{selected.name}</span> has no rules for
          selector{" "}
          <span className="font-mono">{hexSelector || selector}</span>. It will
          evaluate to <em>true</em> on every call until you add one.
        </p>
      ) : null}

      {selected ? (
        <div className="flex flex-wrap items-center gap-2 rounded-md border border-surge-border bg-surge-card-hover/30 px-2 py-1.5">
          <Badge tone="aqua" className="font-mono text-[10px]">
            {selected.name}
          </Badge>
          <span className="text-[11px] text-surge-muted">
            {selected.rules.length} rule(s)
          </span>
          {selected.description ? (
            <span className="text-[11px] text-surge-muted">
              · {selected.description}
            </span>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}
