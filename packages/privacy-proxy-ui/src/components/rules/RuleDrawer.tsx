import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { Drawer } from "../common/Drawer";
import { SelectorPicker } from "./SelectorPicker";
import { EntriesEditor } from "./EntriesEditor";
import {
  useCreateRule,
  useGetRule,
  useReplaceRule,
} from "../../hooks/rules/useRules";
import { AdminApiError } from "../../lib/apiClient";
import type { EntryInput, RuleView } from "../../types/api";

interface RuleDrawerProps {
  open: boolean;
  onClose: () => void;
  /** When provided, drawer is in edit mode. */
  editingId?: number;
  /** Pre-filled selector when launched from a contract-binding flow. */
  presetSelector?: string;
  /** Fired after a successful create — caller may chain into binding flow. */
  onCreated?: (rule: RuleView) => void;
}

interface FormState {
  name: string;
  description: string;
  selector: string;
  mode: "allow" | "deny";
  entries: EntryInput[];
}

const EMPTY: FormState = {
  name: "",
  description: "",
  selector: "",
  mode: "allow",
  entries: [],
};

export function RuleDrawer({
  open,
  onClose,
  editingId,
  presetSelector,
  onCreated,
}: RuleDrawerProps) {
  const isEdit = editingId !== undefined;
  const detail = useGetRule(isEdit ? editingId : undefined);
  const create = useCreateRule();
  const replace = useReplaceRule();
  const [form, setForm] = useState<FormState>(EMPTY);

  useEffect(() => {
    if (!open) return;
    if (isEdit && detail.data) {
      setForm({
        name: detail.data.name,
        description: detail.data.description ?? "",
        selector: detail.data.selector,
        mode: detail.data.mode,
        entries: detail.data.entries.map((e) => ({
          role: e.role,
          lambda_id: e.lambda_id,
        })),
      });
    } else if (!isEdit) {
      setForm({
        ...EMPTY,
        selector: presetSelector ?? "",
      });
    }
  }, [open, isEdit, detail.data, presetSelector]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) {
      toast.error("Give the rule a name.");
      return;
    }
    if (!form.selector) {
      toast.error("Pick a selector or RPC method.");
      return;
    }
    if (form.entries.length === 0) {
      toast.error("Add at least one role entry before saving.");
      return;
    }

    try {
      if (isEdit) {
        await replace.mutateAsync({
          id: editingId!,
          body: {
            name: form.name.trim(),
            description: form.description.trim() || null,
            mode: form.mode,
            entries: form.entries,
          },
        });
        toast.success("Rule updated");
        onClose();
      } else {
        const created = await create.mutateAsync({
          name: form.name.trim(),
          description: form.description.trim() || null,
          selector: form.selector,
          mode: form.mode,
          entries: form.entries,
        });
        toast.success(`Rule "${created.name}" created`);
        onCreated?.(created);
        onClose();
      }
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const submitting = create.isPending || replace.isPending;
  const hasMinimumEntries = form.entries.length > 0;

  return (
    <Drawer
      open={open}
      onClose={onClose}
      title={isEdit ? `Edit rule #${editingId}` : "Create rule"}
      widthClass="max-w-2xl"
      footer={
        <div className="flex items-center justify-end gap-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg px-3 py-2 text-sm font-medium text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
          >
            Cancel
          </button>
          <button
            type="submit"
            form="rule-form"
            disabled={submitting || !hasMinimumEntries}
            title={!hasMinimumEntries ? "Add at least one role entry" : undefined}
            className="rounded-lg bg-surge-primary px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {submitting ? "Saving…" : isEdit ? "Save changes" : "Create rule"}
          </button>
        </div>
      }
    >
      <form id="rule-form" onSubmit={onSubmit} className="space-y-5">
        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Name
          </label>
          <input
            type="text"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            placeholder="e.g. erc20-self-balance-only"
            className="mt-1 w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 text-sm outline-none focus:border-surge-secondary"
          />
        </div>

        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Description (optional)
          </label>
          <textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            placeholder="What does this rule do?"
            rows={2}
            className="mt-1 w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 text-sm outline-none focus:border-surge-secondary"
          />
        </div>

        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Function selector
          </label>
          <div className="mt-1">
            <SelectorPicker
              value={form.selector}
              onChange={(v) => setForm({ ...form, selector: v })}
              disabled={isEdit}
            />
          </div>
          {isEdit ? (
            <p className="mt-1 text-[11px] text-surge-muted">
              Selector is immutable after creation. Delete and recreate the
              rule to change it.
            </p>
          ) : null}
        </div>

        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Mode
          </label>
          <div className="mt-1 grid grid-cols-2 gap-2">
            {(["allow", "deny"] as const).map((m) => (
              <button
                type="button"
                key={m}
                onClick={() => setForm({ ...form, mode: m })}
                className={`rounded-lg border px-3 py-2 text-sm font-medium transition ${
                  form.mode === m
                    ? m === "allow"
                      ? "border-emerald-600 bg-emerald-600 text-white"
                      : "border-amber-600 bg-amber-600 text-white"
                    : "border-surge-border bg-surge-card text-surge-text hover:bg-surge-card-hover"
                }`}
              >
                {m}
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Role entries
          </label>
          <div className="mt-3">
            <EntriesEditor
              entries={form.entries}
              onChange={(entries) => setForm({ ...form, entries })}
              selector={form.selector}
            />
          </div>
        </div>
      </form>
    </Drawer>
  );
}
