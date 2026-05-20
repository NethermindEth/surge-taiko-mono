import { Link } from "react-router-dom";
import { PageHeader } from "../components/layout/PageHeader";
import { useListMembers } from "../hooks/members/useMembers";
import { useListBindings, useListRules } from "../hooks/rules/useRules";
import { useRoles } from "../hooks/useRoles";

interface StatCardProps {
  label: string;
  value: string | number;
  to: string;
  hint?: string;
}

function StatCard({ label, value, to, hint }: StatCardProps) {
  return (
    <Link
      to={to}
      className="glass-card hover-glow group flex flex-col rounded-2xl p-5 transition"
    >
      <span className="text-xs font-medium uppercase tracking-wide text-surge-muted">
        {label}
      </span>
      <span className="mt-2 text-3xl font-semibold tracking-tight text-surge-text">
        {value}
      </span>
      {hint ? (
        <span className="mt-1 text-xs uppercase tracking-wide text-surge-muted">
          {hint}
        </span>
      ) : null}
      <span className="mt-3 text-xs font-medium text-surge-secondary opacity-0 transition group-hover:opacity-100">
        Open →
      </span>
    </Link>
  );
}

export function DashboardPage() {
  const members = useListMembers();
  const rules = useListRules();
  const bindings = useListBindings();
  const roles = useRoles();

  const memberCount = members.data?.length ?? 0;
  const adminCount =
    members.data?.filter((m) => m.role === "admin").length ?? 0;
  const userCount = members.data?.filter((m) => m.role === "user").length ?? 0;
  const ruleCount = rules.data?.length ?? 0;
  const bindingCount = bindings.data?.length ?? 0;
  const entryCount =
    rules.data?.reduce((acc, r) => acc + r.entries.length, 0) ?? 0;

  return (
    <div>
      <PageHeader
        title="Overview"
        description="At-a-glance status of every artifact the proxy enforces."
      />
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard
          label="Members"
          value={members.isLoading ? "…" : memberCount}
          hint={`${adminCount} admin · ${userCount} user`}
          to="/members"
        />
        <StatCard
          label="Access rules"
          value={rules.isLoading ? "…" : ruleCount}
          hint={`${entryCount} role entries`}
          to="/rules"
        />
        <StatCard
          label="Contract bindings"
          value={bindings.isLoading ? "…" : bindingCount}
          hint="(contract, selector) → rule"
          to="/rules"
        />
        <StatCard
          label="Declared roles"
          value={roles.isLoading ? "…" : roles.data?.length ?? 0}
          to="/roles"
        />
      </div>
    </div>
  );
}
