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
import { BindingDrawer } from "../components/rules/BindingDrawer";
import {
  useDeleteBinding,
  useDeleteRule,
  useListBindings,
  useListRules,
} from "../hooks/rules/useRules";
import { useSyntheticSelectors } from "../hooks/useSyntheticSelectors";
import { findCommonSelector } from "../lib/selectors";
import { isAddress } from "../lib/format";
import { AdminApiError } from "../lib/apiClient";
import type { BindingView, RuleView } from "../types/api";

function SelectorCell({ selector }: { selector: string }) {
  const { data: selectors } = useSyntheticSelectors();
  const method = selectors?.find(
    (s) => s.selector.toLowerCase() === selector.toLowerCase(),
  );
  const common = !method ? findCommonSelector(selector) : undefined;
  const subLabel =
    method ? `RPC: ${method.method}` : common ? common.signature : null;
  return (
    <div className="flex flex-col">
      <span className="font-mono text-xs text-surge-text">{selector}</span>
      {subLabel ? (
        <span className="text-[11px] normal-case text-surge-muted">
          {subLabel}
        </span>
      ) : null}
    </div>
  );
}

export function RulesPage() {
  const navigate = useNavigate();
  const [contractFilter, setContractFilter] = useState<string>("");

  const bindings = useListBindings({
    contract:
      contractFilter && isAddress(contractFilter)
        ? contractFilter.toLowerCase()
        : undefined,
  });
  const rules = useListRules();
  const deleteBinding = useDeleteBinding();
  const deleteRule = useDeleteRule();

  const [bindingDrawerOpen, setBindingDrawerOpen] = useState(false);
  const [ruleDrawerOpen, setRuleDrawerOpen] = useState(false);
  const [editingRuleId, setEditingRuleId] = useState<number | undefined>(
    undefined,
  );
  const [confirmDeleteBinding, setConfirmDeleteBinding] =
    useState<BindingView | null>(null);
  const [confirmDeleteRule, setConfirmDeleteRule] = useState<RuleView | null>(
    null,
  );

  const onApplyRule = () => setBindingDrawerOpen(true);
  const onCreateRule = () => {
    setEditingRuleId(undefined);
    setRuleDrawerOpen(true);
  };
  const onEditRule = (rule: RuleView) => {
    setEditingRuleId(rule.id);
    setRuleDrawerOpen(true);
  };

  const onConfirmDeleteBinding = async () => {
    if (!confirmDeleteBinding) return;
    try {
      await deleteBinding.mutateAsync(confirmDeleteBinding.id);
      toast.success("Rule unbound from contract");
      setConfirmDeleteBinding(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const onConfirmDeleteRule = async () => {
    if (!confirmDeleteRule) return;
    try {
      await deleteRule.mutateAsync(confirmDeleteRule.id);
      toast.success(`Rule "${confirmDeleteRule.name}" deleted`);
      setConfirmDeleteRule(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const bindingColumns: Column<BindingView>[] = [
    {
      key: "contract",
      header: "Contract",
      render: (b) => <AddressDisplay value={b.contract_address} />,
    },
    {
      key: "selector",
      header: "Selector",
      render: (b) => <SelectorCell selector={b.selector} />,
    },
    {
      key: "rule",
      header: "Rule",
      render: (b) => (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            navigate(`/rules/${b.rule_id}`);
          }}
          className="text-sm font-semibold text-surge-secondary hover:underline"
        >
          {b.rule_name}
        </button>
      ),
    },
    {
      key: "mode",
      header: "Mode",
      render: (b) => <ModeBadge mode={b.mode} />,
    },
    {
      key: "actions",
      header: "",
      align: "right",
      render: (b) => (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            setConfirmDeleteBinding(b);
          }}
          className="rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50"
        >
          Unbind
        </button>
      ),
    },
  ];

  const ruleColumns: Column<RuleView>[] = [
    {
      key: "name",
      header: "Name",
      render: (r) => (
        <div className="flex flex-col">
          <span className="font-semibold text-surge-text">{r.name}</span>
          {r.description ? (
            <span className="text-[11px] text-surge-muted">{r.description}</span>
          ) : null}
        </div>
      ),
    },
    {
      key: "selector",
      header: "Selector",
      render: (r) => <SelectorCell selector={r.selector} />,
    },
    {
      key: "mode",
      header: "Mode",
      render: (r) => <ModeBadge mode={r.mode} />,
    },
    {
      key: "entries",
      header: "Role entries",
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
      key: "usage",
      header: "Bindings",
      render: (r) => (
        <span className="text-xs text-surge-muted">
          {r.binding_count} contract{r.binding_count === 1 ? "" : "s"}
        </span>
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
              onEditRule(r);
            }}
            className="rounded-md px-2 py-1 text-xs font-medium text-surge-secondary hover:bg-surge-card-hover"
          >
            Edit
          </button>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setConfirmDeleteRule(r);
            }}
            disabled={r.binding_count > 0}
            title={
              r.binding_count > 0
                ? "Unbind from every contract before deleting"
                : undefined
            }
            className="rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent"
          >
            Delete
          </button>
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-10">
      <section>
        <PageHeader
          title="Access rules"
          description="Apply reusable rules to a (contract, selector) pair. Each contract can have many bindings, but only one rule per selector."
          actions={
            <button
              type="button"
              onClick={onApplyRule}
              className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90"
            >
              + Apply rule to contract
            </button>
          }
        />

        <div className="mb-4 flex flex-wrap items-center gap-2">
          <input
            type="text"
            value={contractFilter}
            onChange={(e) => setContractFilter(e.target.value)}
            placeholder="Filter bindings by contract (0x...)"
            className="w-80 rounded-lg border border-surge-border bg-surge-card px-3 py-2 font-mono text-sm outline-none focus:border-surge-secondary"
          />
          <button
            type="button"
            onClick={() => bindings.refetch()}
            className="rounded-lg border border-surge-border bg-surge-card px-3 py-2 text-sm text-surge-muted hover:bg-surge-card-hover"
          >
            Refresh
          </button>
        </div>

        <DataTable
          columns={bindingColumns}
          rows={bindings.data ?? []}
          rowKey={(b) => String(b.id)}
          loading={bindings.isLoading}
          loadingNode={<TableSkeleton />}
          emptyState={
            <EmptyState
              title="No contract bindings yet"
              description="A (contract, selector) pair without a binding is freely callable. Apply a rule to gate one."
              action={
                <button
                  type="button"
                  onClick={onApplyRule}
                  className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white"
                >
                  + Apply rule to contract
                </button>
              }
            />
          }
        />
      </section>

      <section>
        <header className="mb-4 flex flex-wrap items-end justify-between gap-3">
          <div>
            <h2 className="text-lg font-semibold text-surge-text">
              Rule library
            </h2>
            <p className="text-sm text-surge-muted">
              Reusable templates. Each rule is bound to its selector at
              creation and can be applied to many contracts.
            </p>
          </div>
          <button
            type="button"
            onClick={onCreateRule}
            className="rounded-lg border border-surge-secondary bg-surge-secondary/10 px-3 py-2 text-sm font-semibold text-surge-primary hover:bg-surge-secondary/20"
          >
            + Create rule
          </button>
        </header>

        <DataTable
          columns={ruleColumns}
          rows={rules.data ?? []}
          rowKey={(r) => String(r.id)}
          onRowClick={(r) => navigate(`/rules/${r.id}`)}
          loading={rules.isLoading}
          loadingNode={<TableSkeleton />}
          emptyState={
            <EmptyState
              title="No rules defined"
              description="Define a reusable allow/deny rule with role entries, then bind it to a contract selector above."
              action={
                <button
                  type="button"
                  onClick={onCreateRule}
                  className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white"
                >
                  + Create rule
                </button>
              }
            />
          }
        />
      </section>

      <BindingDrawer
        open={bindingDrawerOpen}
        onClose={() => setBindingDrawerOpen(false)}
        presetContract={
          contractFilter && isAddress(contractFilter) ? contractFilter : undefined
        }
      />

      <RuleDrawer
        open={ruleDrawerOpen}
        onClose={() => setRuleDrawerOpen(false)}
        editingId={editingRuleId}
      />

      <ConfirmDialog
        open={!!confirmDeleteBinding}
        title="Unbind rule from contract"
        description={`Removes the (contract, selector) → rule "${confirmDeleteBinding?.rule_name ?? ""}" mapping. The selector becomes freely callable again.`}
        destructive
        confirmLabel="Unbind"
        onConfirm={onConfirmDeleteBinding}
        onCancel={() => setConfirmDeleteBinding(null)}
        isLoading={deleteBinding.isPending}
      />

      <ConfirmDialog
        open={!!confirmDeleteRule}
        title={`Delete rule "${confirmDeleteRule?.name ?? ""}"`}
        description="Removes the rule and all its entries. Bindings must be removed first."
        destructive
        confirmLabel="Delete"
        onConfirm={onConfirmDeleteRule}
        onCancel={() => setConfirmDeleteRule(null)}
        isLoading={deleteRule.isPending}
      />
    </div>
  );
}
