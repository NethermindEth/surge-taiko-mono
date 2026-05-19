import { useMemo } from "react";
import { useRoles } from "../../hooks/useRoles";
import type { EntryInput, RoleName } from "../../types/api";
import { LambdaPicker } from "./LambdaPicker";

interface EntriesEditorProps {
  entries: EntryInput[];
  onChange: (entries: EntryInput[]) => void;
  selector: string;
}

const ENTRY_ROLE_EXCLUDE = new Set<string>(["admin"]);

export function EntriesEditor({ entries, onChange, selector }: EntriesEditorProps) {
  const roles = useRoles();

  const assignableRoles = useMemo(
    () => (roles.data ?? []).filter((r) => !ENTRY_ROLE_EXCLUDE.has(r.name)),
    [roles.data],
  );

  const addRow = () => {
    const usedRoles = new Set(entries.map((e) => e.role));
    const nextRole = assignableRoles.find((r) => !usedRoles.has(r.name))?.name;
    if (!nextRole) return;
    onChange([...entries, { role: nextRole, lambda_id: null }]);
  };

  const updateRow = (index: number, patch: Partial<EntryInput>) => {
    onChange(entries.map((e, i) => (i === index ? { ...e, ...patch } : e)));
  };

  const removeRow = (index: number) => {
    onChange(entries.filter((_, i) => i !== index));
  };

  const usedRoles = new Set(entries.map((e) => e.role));
  const canAddMore = assignableRoles.some((r) => !usedRoles.has(r.name));

  return (
    <div className="space-y-3">
      {entries.map((entry, i) => (
        <div
          key={i}
          className="space-y-2 rounded-xl border border-surge-border bg-surge-card-hover/30 p-3"
        >
          <div className="flex items-start gap-2">
            <div className="flex-1">
              <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                Role
              </label>
              <select
                value={entry.role}
                onChange={(e) =>
                  updateRow(i, {
                    role: e.target.value as RoleName,
                    lambda_id: null,
                  })
                }
                className="mt-1 w-full rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
              >
                {assignableRoles.map((r) => (
                  <option key={r.id} value={r.name}>
                    {r.name}
                  </option>
                ))}
              </select>
            </div>
            <button
              type="button"
              onClick={() => removeRow(i)}
              aria-label="Remove entry"
              className="mt-5 rounded-md p-1.5 text-surge-muted hover:bg-surge-card-hover hover:text-red-600"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
                <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M6 6l1 14a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-14" />
              </svg>
            </button>
          </div>
          <div>
            <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
              Lambda
            </label>
            <div className="mt-1">
              <LambdaPicker
                role={entry.role}
                selector={selector}
                value={entry.lambda_id ?? null}
                onChange={(id) => updateRow(i, { lambda_id: id })}
              />
            </div>
          </div>
        </div>
      ))}
      <button
        type="button"
        onClick={addRow}
        disabled={!canAddMore}
        className="rounded-lg border border-dashed border-surge-border bg-surge-card-hover/30 px-3 py-2 text-sm font-medium text-surge-muted hover:border-surge-secondary hover:text-surge-text disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:border-surge-border disabled:hover:text-surge-muted"
      >
        + Add role entry
      </button>
    </div>
  );
}
