import { useState } from "react";
import toast from "react-hot-toast";
import { PageHeader } from "../components/layout/PageHeader";
import { Badge, RoleBadge } from "../components/common/Badge";
import { EmptyState } from "../components/common/EmptyState";
import { Skeleton } from "../components/common/Skeleton";
import { ConfirmDialog } from "../components/common/ConfirmDialog";
import { useDeleteLambda, useLambdas } from "../hooks/useLambdas";
import { useRoleAttributes } from "../hooks/useRoleAttributes";
import { useSyntheticSelectors } from "../hooks/useSyntheticSelectors";
import { findCommonSelector } from "../lib/selectors";
import { decodeLiteral } from "../lib/abiTypes";
import { AdminApiError } from "../lib/apiClient";
import type {
  LambdaRuleView,
  LambdaView,
  RoleAttributesGroup,
  RoleName,
} from "../types/api";

const CONDITION_LABEL: Record<string, string> = {
  eq: "=",
  neq: "≠",
  gt: ">",
  lt: "<",
  gte: "≥",
  lte: "≤",
};

interface RuleContext {
  syntheticByHex: Map<string, string>;
  roleAttrs: Map<RoleName, Map<string, string>>;
  role: RoleName;
}

function syntheticParamType(method: string, argIndex: number): string | null {
  if (method === "eth_getStorageAt") {
    if (argIndex === 0) return "address";
    if (argIndex === 1) return "bytes32";
    return null;
  }
  return argIndex === 0 ? "address" : null;
}

function paramTypeForCalldata(
  selector: string,
  argIndex: number,
  syntheticByHex: Map<string, string>,
): string | null {
  const common = findCommonSelector(selector);
  if (common) return common.params[argIndex]?.type ?? null;
  const method = syntheticByHex.get(selector.toLowerCase());
  if (method) return syntheticParamType(method, argIndex);
  return null;
}

function lhsType(r: LambdaRuleView, ctx: RuleContext): string | null {
  if (r.lhs_kind === "calldata") {
    const argIndex = r.lhs_offset !== null ? (r.lhs_offset - 4) / 32 : 0;
    return paramTypeForCalldata(r.selector, argIndex, ctx.syntheticByHex);
  }
  if (!r.lhs_attribute) return null;
  return ctx.roleAttrs.get(ctx.role)?.get(r.lhs_attribute) ?? null;
}

function selectorSubLabel(
  sel: string,
  syntheticByHex: Map<string, string>,
): string | null {
  const common = findCommonSelector(sel);
  if (common) return common.signature;
  const method = syntheticByHex.get(sel.toLowerCase());
  if (method) return `RPC: ${method}`;
  return null;
}

function lhsLabel(r: LambdaRuleView, ty: string | null): string {
  if (r.lhs_kind === "calldata") {
    const argIndex = r.lhs_offset !== null ? (r.lhs_offset - 4) / 32 : 0;
    return ty ? `arg${argIndex}: ${ty}` : `arg${argIndex}`;
  }
  return ty ? `attr.${r.lhs_attribute}: ${ty}` : `attr.${r.lhs_attribute}`;
}

function rhsLabel(r: LambdaRuleView, ty: string | null): string {
  switch (r.rhs_kind) {
    case "tx_origin":
      return "tx.origin";
    case "msg_sender":
      return "msg.sender";
    case "literal":
      return decodeLiteral(r.rhs_value, ty);
  }
}

function RuleRow({ r, ctx }: { r: LambdaRuleView; ctx: RuleContext }) {
  const sub = selectorSubLabel(r.selector, ctx.syntheticByHex);
  const ty = lhsType(r, ctx);
  return (
    <tr>
      <td className="px-3 py-2 align-top">
        <div className="flex flex-col">
          <span className="font-mono text-[11px] text-surge-text">
            {r.selector}
          </span>
          {sub ? (
            <span className="text-[10px] text-surge-muted">{sub}</span>
          ) : null}
        </div>
      </td>
      <td className="px-3 py-2 align-top font-mono text-[11px] text-surge-text">
        {lhsLabel(r, ty)}
      </td>
      <td className="px-3 py-2 align-top text-center font-mono text-sm text-surge-text">
        {CONDITION_LABEL[r.condition] ?? r.condition}
      </td>
      <td className="px-3 py-2 align-top font-mono text-[11px] text-surge-text">
        {rhsLabel(r, ty)}
      </td>
    </tr>
  );
}

function LambdaCard({
  lambda,
  ctx,
  onDelete,
}: {
  lambda: LambdaView;
  ctx: Omit<RuleContext, "role">;
  onDelete: (l: LambdaView) => void;
}) {
  const rowCtx: RuleContext = { ...ctx, role: lambda.role };
  return (
    <article className="glass-card rounded-2xl p-4">
      <header className="flex flex-wrap items-start justify-between gap-2">
        <div>
          <h3 className="font-mono text-sm font-semibold text-surge-text">
            {lambda.name}
          </h3>
          {lambda.description ? (
            <p className="mt-1 text-sm text-surge-muted">{lambda.description}</p>
          ) : null}
        </div>
        <div className="flex items-center gap-2">
          {lambda.in_use ? (
            <Badge tone="amber" className="text-[10px]">
              in use
            </Badge>
          ) : null}
          <button
            type="button"
            onClick={() => onDelete(lambda)}
            disabled={lambda.in_use}
            title={
              lambda.in_use
                ? "Detach from all rule entries before deleting"
                : undefined
            }
            className="rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent"
          >
            Delete
          </button>
        </div>
      </header>
      <div className="mt-3 overflow-x-auto rounded-lg border border-surge-border">
        <table className="min-w-full divide-y divide-surge-border text-xs">
          <thead className="bg-surge-card-hover/50">
            <tr>
              <th className="px-3 py-2 text-left font-semibold uppercase tracking-wide text-surge-muted">
                Selector
              </th>
              <th className="px-3 py-2 text-left font-semibold uppercase tracking-wide text-surge-muted">
                LHS
              </th>
              <th className="px-3 py-2 text-center font-semibold uppercase tracking-wide text-surge-muted">
                Cond
              </th>
              <th className="px-3 py-2 text-left font-semibold uppercase tracking-wide text-surge-muted">
                RHS
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-surge-border">
            {lambda.rules.length === 0 ? (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-3 text-center text-surge-muted"
                >
                  No rules — this lambda will always evaluate to <em>true</em>.
                </td>
              </tr>
            ) : (
              lambda.rules.map((r) => <RuleRow key={r.id} r={r} ctx={rowCtx} />)
            )}
          </tbody>
        </table>
      </div>
    </article>
  );
}

function buildRoleAttrsMap(
  groups: RoleAttributesGroup[] | undefined,
): Map<RoleName, Map<string, string>> {
  const out = new Map<RoleName, Map<string, string>>();
  for (const g of groups ?? []) {
    const inner = new Map<string, string>();
    for (const a of g.attributes) inner.set(a.name, a.type);
    out.set(g.role, inner);
  }
  return out;
}

export function LambdasPage() {
  const { data, isLoading } = useLambdas();
  const { data: syntheticSelectors } = useSyntheticSelectors();
  const { data: roleAttrGroups } = useRoleAttributes();
  const deleteLambda = useDeleteLambda();
  const [confirmDelete, setConfirmDelete] = useState<LambdaView | null>(null);

  const syntheticByHex = new Map(
    (syntheticSelectors ?? []).map((s) => [s.selector.toLowerCase(), s.method]),
  );
  const roleAttrs = buildRoleAttrsMap(roleAttrGroups);

  const onDelete = async () => {
    if (!confirmDelete) return;
    try {
      await deleteLambda.mutateAsync(confirmDelete.id);
      toast.success(`Lambda "${confirmDelete.name}" deleted`);
      setConfirmDelete(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  if (isLoading) {
    return (
      <div>
        <PageHeader
          title="Lambdas"
          description="Comparison rule sets attached to access rule entries."
        />
        <div className="space-y-3">
          <Skeleton className="h-16 w-full" />
          <Skeleton className="h-16 w-full" />
        </div>
      </div>
    );
  }

  const totalLambdas = (data ?? []).reduce(
    (acc, g) => acc + g.lambdas.length,
    0,
  );

  return (
    <div>
      <PageHeader
        title="Lambdas"
        description="Stored comparison rule sets. Create lambdas while authoring a rule entry; delete unreferenced ones here."
      />

      <div className="space-y-6">
        {data?.map((group) => (
          <section key={group.role}>
            <h2 className="mb-2 flex items-center gap-2 text-sm font-semibold uppercase tracking-wide text-surge-muted">
              <RoleBadge role={group.role} />
              <span>{group.lambdas.length} lambda(s)</span>
            </h2>
            {group.lambdas.length === 0 ? (
              <div className="rounded-2xl border border-dashed border-surge-border bg-surge-card-hover/30 px-4 py-5 text-sm text-surge-muted">
                No lambdas declared for the{" "}
                <span className="font-medium text-surge-text">{group.role}</span>{" "}
                role yet.
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-3 xl:grid-cols-2">
                {group.lambdas.map((l) => (
                  <LambdaCard
                    key={l.id}
                    lambda={l}
                    ctx={{ syntheticByHex, roleAttrs }}
                    onDelete={setConfirmDelete}
                  />
                ))}
              </div>
            )}
          </section>
        ))}
      </div>

      {(!data || totalLambdas === 0) && (
        <EmptyState
          title="No lambdas yet"
          description="Lambdas are created from inside the rule entry editor when you need a comparison-rule predicate on a role."
        />
      )}

      <ConfirmDialog
        open={!!confirmDelete}
        title={`Delete lambda "${confirmDelete?.name ?? ""}"`}
        description="This removes the lambda and all its comparison rules. The action cannot be undone."
        confirmLabel="Delete"
        destructive
        onConfirm={onDelete}
        onCancel={() => setConfirmDelete(null)}
        isLoading={deleteLambda.isPending}
      />
    </div>
  );
}
