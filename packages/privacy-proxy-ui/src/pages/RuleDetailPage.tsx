import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import toast from "react-hot-toast";
import { PageHeader } from "../components/layout/PageHeader";
import { Badge, ModeBadge, RoleBadge } from "../components/common/Badge";
import { Skeleton } from "../components/common/Skeleton";
import { ConfirmDialog } from "../components/common/ConfirmDialog";
import { RuleDrawer } from "../components/rules/RuleDrawer";
import { LambdaPicker } from "../components/rules/LambdaPicker";
import {
  useAddEntry,
  useDeleteEntry,
  useGetRule,
  useUpdateEntry,
} from "../hooks/rules/useRules";
import { useRoles } from "../hooks/useRoles";
import { useSyntheticSelectors } from "../hooks/useSyntheticSelectors";
import { AdminApiError } from "../lib/apiClient";
import type { EntryView, RoleName } from "../types/api";

export function RuleDetailPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);
  const detail = useGetRule(Number.isNaN(id) ? undefined : id);
  const roles = useRoles();
  const selectors = useSyntheticSelectors();
  const addEntry = useAddEntry();
  const updateEntry = useUpdateEntry();
  const deleteEntry = useDeleteEntry();

  const [editOpen, setEditOpen] = useState(false);
  const [newEntryRole, setNewEntryRole] = useState<RoleName>("user");
  const [newEntryLambda, setNewEntryLambda] = useState<number | null>(null);
  const [editingLambda, setEditingLambda] = useState<{
    entry: EntryView;
    lambdaId: number | null;
  } | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<EntryView | null>(null);

  if (Number.isNaN(id)) {
    return <p className="text-sm text-red-600">Invalid rule id.</p>;
  }

  if (detail.isLoading || !detail.data) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-7 w-48" />
        <Skeleton className="h-20 w-full" />
      </div>
    );
  }

  const rule = detail.data;
  const selectorMethod = selectors.data?.find(
    (s) => s.selector.toLowerCase() === rule.selector.toLowerCase(),
  );

  const onAddEntry = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await addEntry.mutateAsync({
        ruleId: rule.id,
        body: {
          role: newEntryRole,
          lambda_id: newEntryLambda,
        },
      });
      toast.success("Entry added");
      setNewEntryRole("user");
      setNewEntryLambda(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const onSaveLambda = async () => {
    if (!editingLambda) return;
    try {
      await updateEntry.mutateAsync({
        ruleId: rule.id,
        entryId: editingLambda.entry.id,
        body: {
          lambda_id: editingLambda.lambdaId,
        },
      });
      toast.success("Entry updated");
      setEditingLambda(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const onDeleteEntry = async () => {
    if (!confirmDelete) return;
    try {
      await deleteEntry.mutateAsync({
        ruleId: rule.id,
        entryId: confirmDelete.id,
      });
      toast.success("Entry removed");
      setConfirmDelete(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const usedRoles = new Set(rule.entries.map((e) => e.role));
  const availableRoles =
    roles.data?.filter((r) => r.name !== "admin" && !usedRoles.has(r.name)) ?? [];

  return (
    <div>
      <PageHeader
        title={`Rule #${rule.id}`}
        description="Edit role entries individually or replace the whole rule."
        actions={
          <div className="flex gap-2">
            <Link
              to="/rules"
              className="rounded-lg border border-surge-border bg-surge-card px-3 py-2 text-sm text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
            >
              ← All rules
            </Link>
            <button
              type="button"
              onClick={() => setEditOpen(true)}
              className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90"
            >
              Replace rule
            </button>
          </div>
        }
      />

      <section className="glass-card mb-6 rounded-2xl p-5">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-surge-muted">
              Name
            </p>
            <p className="mt-1 text-sm font-semibold text-surge-text">
              {rule.name}
            </p>
            {rule.description ? (
              <p className="mt-1 text-xs text-surge-muted">{rule.description}</p>
            ) : null}
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-surge-muted">
              Selector
            </p>
            <p className="mt-1 font-mono text-sm text-surge-text">
              {rule.selector}
            </p>
            {selectorMethod ? (
              <p className="text-[10px] uppercase tracking-wide text-surge-muted">
                {selectorMethod.method}
              </p>
            ) : null}
            <p className="mt-1 text-[11px] text-surge-muted">
              Used by {rule.binding_count} contract(s)
            </p>
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-surge-muted">
              Mode
            </p>
            <div className="mt-1">
              <ModeBadge mode={rule.mode} />
            </div>
          </div>
        </div>
      </section>

      <section>
        <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-surge-muted">
          Role entries ({rule.entries.length})
        </h2>

        <div className="overflow-hidden rounded-2xl border border-surge-border bg-surge-card">
          {rule.entries.length === 0 ? (
            <p className="px-6 py-10 text-center text-sm text-surge-muted">
              No entries on this rule yet.
            </p>
          ) : (
            <table className="min-w-full divide-y divide-surge-border text-sm">
              <thead className="bg-surge-card-hover/60">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-surge-muted">
                    Role
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-surge-muted">
                    Lambda
                  </th>
                  <th className="px-4 py-3 text-right" />
                </tr>
              </thead>
              <tbody className="divide-y divide-surge-border">
                {rule.entries.map((e) => {
                  const isEditing = editingLambda?.entry.id === e.id;
                  return (
                    <tr key={e.id}>
                      <td className="px-4 py-3 align-top">
                        <RoleBadge role={e.role} />
                      </td>
                      <td className="px-4 py-3 align-top">
                        {isEditing ? (
                          <LambdaPicker
                            role={e.role}
                            selector={rule.selector}
                            value={editingLambda!.lambdaId}
                            onChange={(id) =>
                              setEditingLambda({ entry: e, lambdaId: id })
                            }
                          />
                        ) : e.lambda_name ? (
                          <Badge tone="aqua" className="font-mono text-[10px]">
                            {e.lambda_name}
                          </Badge>
                        ) : (
                          <span className="text-xs text-surge-muted">
                            (no lambda)
                          </span>
                        )}
                      </td>
                      <td className="px-4 py-3 text-right align-top">
                        {isEditing ? (
                          <div className="flex justify-end gap-1">
                            <button
                              type="button"
                              onClick={() => setEditingLambda(null)}
                              className="rounded-md px-2 py-1 text-xs font-medium text-surge-muted hover:bg-surge-card-hover"
                            >
                              Cancel
                            </button>
                            <button
                              type="button"
                              onClick={onSaveLambda}
                              disabled={updateEntry.isPending}
                              className="rounded-md bg-surge-primary px-2 py-1 text-xs font-semibold text-white hover:bg-surge-primary/90"
                            >
                              {updateEntry.isPending ? "Saving…" : "Save"}
                            </button>
                          </div>
                        ) : (
                          <div className="flex justify-end gap-1">
                            <button
                              type="button"
                              onClick={() =>
                                setEditingLambda({
                                  entry: e,
                                  lambdaId: e.lambda_id ?? null,
                                })
                              }
                              className="rounded-md px-2 py-1 text-xs font-medium text-surge-secondary hover:bg-surge-card-hover"
                            >
                              Edit lambda
                            </button>
                            <button
                              type="button"
                              onClick={() => setConfirmDelete(e)}
                              className="rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50"
                            >
                              Remove
                            </button>
                          </div>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>

        <form
          onSubmit={onAddEntry}
          className="mt-4 space-y-3 rounded-xl border border-dashed border-surge-border bg-surge-card-hover/30 p-3"
        >
          <div className="grid grid-cols-[minmax(0,1fr)_auto] items-end gap-2">
            <div>
              <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
                Role
              </label>
              <select
                value={newEntryRole}
                onChange={(e) => {
                  setNewEntryRole(e.target.value as RoleName);
                  setNewEntryLambda(null);
                }}
                className="mt-1 w-full rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
              >
                {availableRoles.length === 0 ? (
                  <option value="" disabled>
                    All roles already attached
                  </option>
                ) : (
                  availableRoles.map((r) => (
                    <option key={r.id} value={r.name}>
                      {r.name}
                    </option>
                  ))
                )}
              </select>
            </div>
            <button
              type="submit"
              disabled={availableRoles.length === 0 || addEntry.isPending}
              className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {addEntry.isPending ? "Adding…" : "Add entry"}
            </button>
          </div>
          <div>
            <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
              Lambda
            </label>
            <div className="mt-1">
              <LambdaPicker
                role={newEntryRole}
                selector={rule.selector}
                value={newEntryLambda}
                onChange={(id) => setNewEntryLambda(id)}
              />
            </div>
          </div>
        </form>
      </section>

      <RuleDrawer
        open={editOpen}
        onClose={() => setEditOpen(false)}
        editingId={rule.id}
      />

      <ConfirmDialog
        open={!!confirmDelete}
        title="Remove role entry"
        description={`Removes the ${confirmDelete?.role ?? ""} entry from this rule.`}
        confirmLabel="Remove"
        destructive
        onConfirm={onDeleteEntry}
        onCancel={() => setConfirmDelete(null)}
        isLoading={deleteEntry.isPending}
      />
    </div>
  );
}
