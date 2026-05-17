import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import toast from "react-hot-toast";
import { PageHeader } from "../components/layout/PageHeader";
import { AddressDisplay } from "../components/common/AddressDisplay";
import { Badge, ModeBadge, RoleBadge } from "../components/common/Badge";
import { Skeleton } from "../components/common/Skeleton";
import { ConfirmDialog } from "../components/common/ConfirmDialog";
import { RuleDrawer } from "../components/rules/RuleDrawer";
import {
  useAddEntry,
  useDeleteEntry,
  useGetRule,
  useUpdateEntry,
} from "../hooks/rules/useRules";
import { useLambdas } from "../hooks/useLambdas";
import { useRoles } from "../hooks/useRoles";
import { useSyntheticSelectors } from "../hooks/useSyntheticSelectors";
import { AdminApiError } from "../lib/apiClient";
import type { EntryView, RoleName } from "../types/api";

export function RuleDetailPage() {
  const params = useParams<{ id: string }>();
  const id = Number(params.id);
  const detail = useGetRule(Number.isNaN(id) ? undefined : id);
  const lambdas = useLambdas();
  const roles = useRoles();
  const selectors = useSyntheticSelectors();
  const addEntry = useAddEntry();
  const updateEntry = useUpdateEntry();
  const deleteEntry = useDeleteEntry();

  const [editOpen, setEditOpen] = useState(false);
  const [newEntryRole, setNewEntryRole] = useState<RoleName>("user");
  const [newEntryLambda, setNewEntryLambda] = useState<string>("");
  const [editingLambda, setEditingLambda] = useState<{
    entry: EntryView;
    lambda: string;
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
    (s) => s.selector.toLowerCase() === rule.function_selector.toLowerCase(),
  );

  const lambdasForRole = (role: string) =>
    lambdas.data?.find((g) => g.role === role)?.lambdas ?? [];

  const onAddEntry = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await addEntry.mutateAsync({
        ruleId: rule.id,
        body: {
          role: newEntryRole,
          lambda_name: newEntryLambda || null,
        },
      });
      toast.success("Entry added");
      setNewEntryRole("user");
      setNewEntryLambda("");
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
          lambda_name: editingLambda.lambda || null,
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

  // Roles not yet on the rule — used to populate the add-entry role
  // dropdown. Admin is excluded for the same reason as in EntriesEditor:
  // admins aren't gated by rule entries.
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
              Contract
            </p>
            <div className="mt-1">
              <AddressDisplay value={rule.contract_address} />
            </div>
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-surge-muted">
              Selector
            </p>
            <p className="mt-1 font-mono text-sm text-surge-text">
              {rule.function_selector}
            </p>
            {selectorMethod ? (
              <p className="text-[10px] uppercase tracking-wide text-surge-muted">
                {selectorMethod.method}
              </p>
            ) : null}
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
                  const options = lambdasForRole(e.role);
                  return (
                    <tr key={e.id}>
                      <td className="px-4 py-3 align-middle">
                        <RoleBadge role={e.role} />
                      </td>
                      <td className="px-4 py-3 align-middle">
                        {isEditing ? (
                          options.length === 0 ? (
                            <span className="text-xs text-surge-muted">
                              No lambdas available for {e.role}.
                            </span>
                          ) : (
                            <select
                              value={editingLambda!.lambda}
                              onChange={(ev) =>
                                setEditingLambda({
                                  entry: e,
                                  lambda: ev.target.value,
                                })
                              }
                              className="rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
                            >
                              <option value="">(no lambda)</option>
                              {options.map((l) => (
                                <option key={l.name} value={l.name}>
                                  {l.name}
                                </option>
                              ))}
                            </select>
                          )
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
                      <td className="px-4 py-3 text-right align-middle">
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
                                  lambda: e.lambda_name ?? "",
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
          className="mt-4 grid grid-cols-[minmax(0,1fr)_minmax(0,2fr)_auto] items-end gap-2 rounded-xl border border-dashed border-surge-border bg-surge-card-hover/30 p-3"
        >
          <div>
            <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
              Role
            </label>
            <select
              value={newEntryRole}
              onChange={(e) => {
                setNewEntryRole(e.target.value as RoleName);
                setNewEntryLambda("");
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
          <div>
            <label className="block text-[10px] font-semibold uppercase tracking-wide text-surge-muted">
              Lambda
            </label>
            <select
              value={newEntryLambda}
              onChange={(e) => setNewEntryLambda(e.target.value)}
              className="mt-1 w-full rounded-md border border-surge-border bg-surge-card px-2 py-1.5 text-sm outline-none focus:border-surge-secondary"
              disabled={lambdasForRole(newEntryRole).length === 0}
            >
              <option value="">(no lambda)</option>
              {lambdasForRole(newEntryRole).map((l) => (
                <option key={l.name} value={l.name}>
                  {l.name}
                </option>
              ))}
            </select>
          </div>
          <button
            type="submit"
            disabled={availableRoles.length === 0 || addEntry.isPending}
            className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {addEntry.isPending ? "Adding…" : "Add entry"}
          </button>
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
