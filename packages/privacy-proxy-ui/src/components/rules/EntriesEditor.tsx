import { useMemo } from "react";
import { useLambdas } from "../../hooks/useLambdas";
import { useRoles } from "../../hooks/useRoles";
import type { EntryInput, RoleName } from "../../types/api";

interface EntriesEditorProps {
  entries: EntryInput[];
  onChange: (entries: EntryInput[]) => void;
}

/**
 * Roles that can carry a rule entry. Admin is excluded — admins aren't
 * gated by rule entries (they bypass user-role allow/deny matrices), and
 * the proxy rejects admin entries with a lambda anyway.
 */
const ENTRY_ROLE_EXCLUDE = new Set<string>(["admin"]);

export function EntriesEditor({ entries, onChange }: EntriesEditorProps) {
  const roles = useRoles();
  const lambdas = useLambdas();

  const assignableRoles = useMemo(
    () => (roles.data ?? []).filter((r) => !ENTRY_ROLE_EXCLUDE.has(r.name)),
    [roles.data],
  );

  const lambdasByRole = useMemo(() => {
    const map = new Map<string, string[]>();
    lambdas.data?.forEach((g) => {
      map.set(
        g.role,
        g.lambdas.map((l) => l.name),
      );
    });
    return map;
  }, [lambdas.data]);

  const addRow = () => {
    const usedRoles = new Set(entries.map((e) => e.role));
    const nextRole = assignableRoles.find((r) => !usedRoles.has(r.name))?.name;
    if (!nextRole) return; // all assignable roles already attached
    onChange([...entries, { role: nextRole, lambda_name: null }]);
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
      {entries.map((entry, i) => {
        const available = lambdasByRole.get(entry.role) ?? [];
        const noLambdas = available.length === 0;
        return (
          <div
            key={i}
            className="grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)_auto] items-center gap-2 rounded-xl border border-surge-border bg-surge-card-hover/30 p-3"
          >
            <div>
              <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                Role
              </label>
              <select
                value={entry.role}
                onChange={(e) =>
                  updateRow(i, {
                    role: e.target.value as RoleName,
                    lambda_name: null,
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
            <div>
              <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                Lambda
              </label>
              {noLambdas ? (
                <p className="mt-1 rounded-md border border-dashed border-surge-border bg-surge-card px-2 py-1.5 text-xs text-surge-muted">
                  No lambdas available for this role.
                </p>
              ) : (
                <select
                  value={entry.lambda_name ?? ""}
                  onChange={(e) =>
                    updateRow(i, {
                      lambda_name: e.target.value || null,
                    })
                  }
                  className="mt-1 w-full rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
                >
                  <option value="">(no lambda)</option>
                  {available.map((name) => (
                    <option key={name} value={name}>
                      {name}
                    </option>
                  ))}
                </select>
              )}
            </div>
            <button
              type="button"
              onClick={() => removeRow(i)}
              aria-label="Remove entry"
              className="self-end rounded-md p-1.5 text-surge-muted hover:bg-surge-card-hover hover:text-red-600"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
                <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M6 6l1 14a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-14" />
              </svg>
            </button>
          </div>
        );
      })}
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
