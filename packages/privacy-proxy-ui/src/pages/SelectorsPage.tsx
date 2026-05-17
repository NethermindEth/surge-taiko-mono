import { PageHeader } from "../components/layout/PageHeader";
import { DataTable, type Column } from "../components/common/DataTable";
import { Badge } from "../components/common/Badge";
import { TableSkeleton } from "../components/common/Skeleton";
import { useSyntheticSelectors } from "../hooks/useSyntheticSelectors";
import type { SyntheticSelector } from "../types/api";

const columns: Column<SyntheticSelector>[] = [
  {
    key: "method",
    header: "RPC method",
    render: (s) => (
      <span className="font-mono text-sm text-surge-text">{s.method}</span>
    ),
  },
  {
    key: "selector",
    header: "Synthetic selector",
    render: (s) => (
      <span className="font-mono text-sm text-surge-muted">{s.selector}</span>
    ),
  },
];

export function SelectorsPage() {
  const { data, isLoading } = useSyntheticSelectors();

  return (
    <div>
      <PageHeader
        title="Gated RPC endpoints"
        description="These are reserved function selectors that can be used to gate access to entire RPC endpoints, instead of a specific contract function."
        tag={<Badge tone="amber">Read-only</Badge>}
      />
      <DataTable
        columns={columns}
        rows={data ?? []}
        rowKey={(s) => s.selector}
        loading={isLoading}
        loadingNode={<TableSkeleton rows={3} />}
      />
    </div>
  );
}
