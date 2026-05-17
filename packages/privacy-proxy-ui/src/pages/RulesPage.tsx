import { useState } from "react";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { PageHeader } from "../components/layout/PageHeader";
import { DataTable, type Column } from "../components/common/DataTable";
import { AddressDisplay } from "../components/common/AddressDisplay";
import { ModeBadge, RoleBadge, Badge } from "../components/common/Badge";
import { EmptyState } from "../components/common/EmptyState";
import { TableSkeleton } from "../components/common/Skeleton";
import { ConfirmDialog } from "../components/common/ConfirmDialog";
import { RuleDrawer } from "../components/rules/RuleDrawer";
import {
  useDeleteRule,
  useListRules,
} from "../hooks/rules/useRules";
import { useSyntheticSelectors } from "../hooks/useSyntheticSelectors";
import { isAddress } from "../lib/format";
import { findCommonSelector } from "../lib/selectors";
import { AdminApiError } from "../lib/apiClient";
import type { RuleView } from "../types/api";

function SelectorCell({ selector }: { selector: string }) {
  const { data: selectors } = useSyntheticSelectors();
  const method = selectors?.find(
    (s) => s.selector.toLowerCase() === selector.toLowerCase(),
  );
  const common = !method ? findCommonSelector(selector) : undefined;
  // Sub-label below the hex: "RPC: eth_getBalance" for gated RPC endpoints
  // (lower-case styling to read as a label, preserving the method's camelCase),
  // or the human signature for well-known contract selectors.
  const subLabel =
    method ? `RPC: ${method.method}` : common ? common.signature : null;
  return (
    <div className="flex flex-col">
      <span className="font-mono text-xs text-surge-text">{selector}</span>
      {subLabel ? (
        <span className="text-[11px] font-normal normal-case text-surge-muted">
          {subLabel}
        </span>
      ) : null}
    </div>
  );
}

export function RulesPage() {
  const navigate = useNavigate();
  const [contractFilter, setContractFilter] = useState<string>("");
  const [modeFilter, setModeFilter] = useState<"" | "allow" | "deny">("");
  const list = useListRules(
    contractFilter && isAddress(contractFilter) ? contractFilter : undefined,
  );
  const del = useDeleteRule();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingId, setEditingId] = useState<number | undefined>(undefined);
  const [confirmDelete, setConfirmDelete] = useState<RuleView | null>(null);

  const rows = (list.data ?? []).filter((r) =>
    modeFilter ? r.mode === modeFilter : true,
  );

  const onCreate = () => {
    setEditingId(undefined);
    setDrawerOpen(true);
  };
  const onEdit = (rule: RuleView) => {
    setEditingId(rule.id);
    setDrawerOpen(true);
  };

  const onConfirmDelete = async () => {
    if (!confirmDelete) return;
    try {
      await del.mutateAsync(confirmDelete.id);
      toast.success(`Rule #${confirmDelete.id} deleted`);
      setConfirmDelete(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const columns: Column<RuleView>[] = [
    {
      key: "id",
      header: "#",
      width: "60px",
      render: (r) => <span className="text-xs text-surge-muted">{r.id}</span>,
    },
    {
      key: "contract",
      header: "Contract",
      render: (r) => <AddressDisplay value={r.contract_address} />,
    },
    {
      key: "selector",
      header: "Selector",
      render: (r) => <SelectorCell selector={r.function_selector} />,
    },
    {
      key: "mode",
      header: "Mode",
      render: (r) => <ModeBadge mode={r.mode} />,
    },
    {
      key: "entries",
      header: "Entries",
      render: (r) =>
        r.entries.length === 0 ? (
          <span className="text-xs text-surge-muted">—</span>
        ) : (
          <div className="flex flex-wrap gap-1">
            {r.entries.map((e) => (
              <span key={e.id} className="inline-flex items-center gap-1">
                <RoleBadge role={e.role} />
                {e.lambda_name ? (
                  <Badge tone="aqua" className="font-mono text-[10px]">
                    {e.lambda_name}
                  </Badge>
                ) : null}
              </span>
            ))}
          </div>
        ),
    },
    {
      key: "actions",
      header: "",
      align: "right",
      render: (r) => (
        <div className="flex items-center justify-end gap-1">
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              onEdit(r);
            }}
            className="rounded-md px-2 py-1 text-xs font-medium text-surge-secondary hover:bg-surge-card-hover"
          >
            Edit
          </button>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setConfirmDelete(r);
            }}
            className="rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50"
          >
            Delete
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Access rules"
        description="Per-(contract, selector) allow/deny rules with role entries."
        actions={
          <button
            type="button"
            onClick={onCreate}
            className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90"
          >
            + Create rule
          </button>
        }
      />

      <div className="mb-4 flex flex-wrap items-center gap-2">
        <input
          type="text"
          value={contractFilter}
          onChange={(e) => setContractFilter(e.target.value)}
          placeholder="Filter by contract address (0x...)"
          className="w-80 rounded-lg border border-surge-border bg-surge-card px-3 py-2 font-mono text-sm outline-none focus:border-surge-secondary"
        />
        <select
          value={modeFilter}
          onChange={(e) =>
            setModeFilter(e.target.value as "" | "allow" | "deny")
          }
          className="rounded-lg border border-surge-border bg-surge-card px-3 py-2 text-sm outline-none focus:border-surge-secondary"
        >
          <option value="">All modes</option>
          <option value="allow">allow</option>
          <option value="deny">deny</option>
        </select>
        <button
          type="button"
          onClick={() => list.refetch()}
          className="rounded-lg border border-surge-border bg-surge-card px-3 py-2 text-sm text-surge-muted hover:bg-surge-card-hover"
        >
          Refresh
        </button>
      </div>

      <DataTable
        columns={columns}
        rows={rows}
        rowKey={(r) => String(r.id)}
        onRowClick={(r) => navigate(`/rules/${r.id}`)}
        loading={list.isLoading}
        loadingNode={<TableSkeleton />}
        emptyState={
          <EmptyState
            title="No rules yet"
            description="Contracts/selectors with no rule are freely callable. Create your first rule to gate one."
            action={
              <button
                type="button"
                onClick={onCreate}
                className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white"
              >
                + Create rule
              </button>
            }
          />
        }
      />

      <RuleDrawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        editingId={editingId}
      />

      <ConfirmDialog
        open={!!confirmDelete}
        title={`Delete rule #${confirmDelete?.id ?? ""}`}
        description="Removes the rule and all of its entries. The (contract, selector) pair becomes freely callable again."
        confirmString={String(confirmDelete?.id ?? "")}
        destructive
        confirmLabel="Delete"
        onConfirm={onConfirmDelete}
        onCancel={() => setConfirmDelete(null)}
        isLoading={del.isPending}
      />
    </div>
  );
}
