import { PageHeader } from "../components/layout/PageHeader";
import { Badge, RoleBadge } from "../components/common/Badge";
import { Skeleton } from "../components/common/Skeleton";
import { useRoles } from "../hooks/useRoles";

const ROLE_BLURBS: Record<string, string> = {
  admin: "Identity-only. Holds the keys to this admin surface.",
  user: "Carries the typed attributes kyc + blacklisted. Self-registers on first sign-in.",
};

export function RolesPage() {
  const { data, isLoading } = useRoles();

  return (
    <div>
      <PageHeader
        title="Roles"
        description="The role set this build recognizes. Roles are code-declared — adding one is a code change plus a migration."
        tag={<Badge tone="amber">Read-only</Badge>}
      />

      {isLoading ? (
        <div className="space-y-3">
          <Skeleton className="h-16 w-full" />
          <Skeleton className="h-16 w-full" />
        </div>
      ) : (
        <div className="space-y-3">
          {data?.map((r) => (
            <article
              key={r.id}
              className="glass-card hover-glow flex items-start gap-4 rounded-2xl p-4"
            >
              <RoleBadge role={r.name} />
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-surge-text">{r.name}</p>
                <p className="text-sm text-surge-muted">
                  {ROLE_BLURBS[r.name] ?? "Custom role declared by this build."}
                </p>
              </div>
              <span className="text-xs text-surge-muted">#{r.id}</span>
            </article>
          ))}
        </div>
      )}

      <div className="mt-6 rounded-xl border border-surge-border bg-surge-card-hover/30 p-4 text-sm text-surge-muted">
        Want a new role? Open{" "}
        <a
          className="text-surge-secondary hover:underline"
          href="https://github.com/NethermindEth/surge-taiko-mono/blob/main/packages/privacy-proxy/docs/system-design.md#6-roles"
          target="_blank"
          rel="noreferrer"
        >
          system-design.md §6
        </a>{" "}
        for the step-by-step procedure (8 touchpoints across the proxy crate).
      </div>
    </div>
  );
}
