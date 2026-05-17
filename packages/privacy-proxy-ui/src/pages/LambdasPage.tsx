import { PageHeader } from "../components/layout/PageHeader";
import { Badge, RoleBadge } from "../components/common/Badge";
import { EmptyState } from "../components/common/EmptyState";
import { Skeleton } from "../components/common/Skeleton";
import { useLambdas } from "../hooks/useLambdas";
import { useSyntheticSelectors } from "../hooks/useSyntheticSelectors";
import { findCommonSelector } from "../lib/selectors";

export function LambdasPage() {
  const { data, isLoading } = useLambdas();
  const { data: syntheticSelectors } = useSyntheticSelectors();
  const syntheticByHex = new Map(
    (syntheticSelectors ?? []).map((s) => [s.selector.toLowerCase(), s.method]),
  );

  /**
   * Sub-label rendered beneath the raw hex selector chip. Prefers (in
   * order): the common-selector signature, then the synthetic RPC method
   * with an "RPC:" prefix. Returns null when the selector is unknown to
   * both tables — the chip then shows just the hex.
   */
  const subLabelFor = (sel: string): string | null => {
    const common = findCommonSelector(sel);
    if (common) return common.signature;
    const method = syntheticByHex.get(sel.toLowerCase());
    if (method) return `RPC: ${method}`;
    return null;
  };

  if (isLoading) {
    return (
      <div>
        <PageHeader
          title="Lambdas"
          description="In-build predicates available to rule entries."
          tag={<Badge tone="amber">Read-only</Badge>}
        />
        <div className="space-y-3">
          <Skeleton className="h-16 w-full" />
          <Skeleton className="h-16 w-full" />
        </div>
      </div>
    );
  }

  return (
    <div>
      <PageHeader
        title="Lambdas"
        description="In-build predicates available to rule entries. Each lambda runs against its role's typed attribute struct."
        tag={<Badge tone="amber">Read-only</Badge>}
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
                role. Entries on this role cannot carry a lambda.
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                {group.lambdas.map((l) => (
                  <article
                    key={`${group.role}-${l.name}`}
                    className="glass-card hover-glow rounded-2xl p-4"
                  >
                    <header className="flex flex-wrap items-start justify-between gap-2">
                      <h3 className="font-mono text-sm font-semibold text-surge-text">
                        {l.name}
                      </h3>
                      {l.expected_selectors.length === 0 ? (
                        <Badge tone="neutral" className="text-[10px]">
                          any selector
                        </Badge>
                      ) : (
                        <div className="flex flex-wrap justify-end gap-2">
                          {l.expected_selectors.map((sel) => {
                            const sub = subLabelFor(sel);
                            return (
                              <div
                                key={sel}
                                className="flex flex-col items-end"
                              >
                                <Badge
                                  tone="aqua"
                                  className="font-mono text-[10px]"
                                >
                                  {sel}
                                </Badge>
                                {sub ? (
                                  <span className="mt-0.5 text-[10px] normal-case text-surge-muted">
                                    {sub}
                                  </span>
                                ) : null}
                              </div>
                            );
                          })}
                        </div>
                      )}
                    </header>
                    <p className="mt-2 text-sm text-surge-muted">
                      {l.description}
                    </p>
                  </article>
                ))}
              </div>
            )}
          </section>
        ))}
      </div>

      {(!data || data.length === 0) && (
        <EmptyState
          title="No lambdas in this build"
          description="Lambdas ship with the binary. See docs/system-design.md to add new ones."
        />
      )}

      <div className="mt-6 rounded-xl border border-surge-border bg-surge-card-hover/30 p-4 text-sm text-surge-muted">
        Want a new lambda? Open{" "}
        <a
          className="text-surge-secondary hover:underline"
          href="https://github.com/NethermindEth/surge-taiko-mono/blob/main/packages/privacy-proxy/docs/system-design.md#5-lambdas"
          target="_blank"
          rel="noreferrer"
        >
          system-design.md §5
        </a>{" "}
        for the procedure (add a sibling module under the target role's
        directory and register a <code className="font-mono">LambdaSpec</code>).
      </div>
    </div>
  );
}
