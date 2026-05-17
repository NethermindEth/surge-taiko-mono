import type { ReactNode } from "react";

export interface Column<T> {
  key: string;
  header: ReactNode;
  render: (row: T) => ReactNode;
  width?: string;
  align?: "left" | "right";
}

interface DataTableProps<T> {
  columns: Column<T>[];
  rows: T[];
  rowKey: (row: T) => string;
  onRowClick?: (row: T) => void;
  emptyState?: ReactNode;
  loading?: boolean;
  loadingNode?: ReactNode;
}

export function DataTable<T>({
  columns,
  rows,
  rowKey,
  onRowClick,
  emptyState,
  loading,
  loadingNode,
}: DataTableProps<T>) {
  if (loading) {
    return <div className="rounded-2xl border border-surge-border bg-surge-card">{loadingNode}</div>;
  }

  if (rows.length === 0 && emptyState) {
    return (
      <div className="rounded-2xl border border-surge-border bg-surge-card">
        {emptyState}
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-surge-border bg-surge-card">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-surge-border text-sm">
          <thead className="bg-surge-card-hover/60">
            <tr>
              {columns.map((col) => (
                <th
                  key={col.key}
                  scope="col"
                  style={col.width ? { width: col.width } : undefined}
                  className={`px-4 py-3 text-xs font-semibold uppercase tracking-wide text-surge-muted ${
                    col.align === "right" ? "text-right" : "text-left"
                  }`}
                >
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-surge-border bg-surge-card">
            {rows.map((row) => (
              <tr
                key={rowKey(row)}
                onClick={onRowClick ? () => onRowClick(row) : undefined}
                className={`transition ${
                  onRowClick
                    ? "cursor-pointer hover:bg-surge-card-hover"
                    : ""
                }`}
              >
                {columns.map((col) => (
                  <td
                    key={col.key}
                    className={`px-4 py-3 align-middle text-surge-text ${
                      col.align === "right" ? "text-right" : ""
                    }`}
                  >
                    {col.render(row)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
