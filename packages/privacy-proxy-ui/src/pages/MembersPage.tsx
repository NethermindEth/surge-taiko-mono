import { useMemo, useState } from "react";
import toast from "react-hot-toast";
import { PageHeader } from "../components/layout/PageHeader";
import { DataTable, type Column } from "../components/common/DataTable";
import { AddressDisplay } from "../components/common/AddressDisplay";
import { RoleBadge } from "../components/common/Badge";
import { EmptyState } from "../components/common/EmptyState";
import { TableSkeleton } from "../components/common/Skeleton";
import { ConfirmDialog } from "../components/common/ConfirmDialog";
import { MemberDrawer } from "../components/members/MemberDrawer";
import {
  useDeleteMember,
  useListMembers,
  useRevokeMemberTokens,
} from "../hooks/members/useMembers";
import { useRoles } from "../hooks/useRoles";
import { timeAgo } from "../lib/format";
import { useAuth } from "../context/AuthContext";
import { AdminApiError } from "../lib/apiClient";
import type { MemberView } from "../types/api";

export function MembersPage() {
  const { session } = useAuth();
  const roles = useRoles();
  const [roleFilter, setRoleFilter] = useState<string>("");
  const [search, setSearch] = useState<string>("");
  const list = useListMembers(roleFilter || undefined);
  const del = useDeleteMember();
  const revoke = useRevokeMemberTokens();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingEoa, setEditingEoa] = useState<string | undefined>(undefined);
  const [confirmDelete, setConfirmDelete] = useState<MemberView | null>(null);
  const [confirmRevoke, setConfirmRevoke] = useState<MemberView | null>(null);

  const filtered = useMemo(() => {
    const rows = list.data ?? [];
    if (!search.trim()) return rows;
    const q = search.toLowerCase();
    return rows.filter((m) => m.eoa_address.toLowerCase().includes(q));
  }, [list.data, search]);

  const onCreate = () => {
    setEditingEoa(undefined);
    setDrawerOpen(true);
  };
  const onEdit = (row: MemberView) => {
    setEditingEoa(row.eoa_address);
    setDrawerOpen(true);
  };

  const onConfirmDelete = async () => {
    if (!confirmDelete) return;
    try {
      await del.mutateAsync(confirmDelete.eoa_address);
      toast.success("Member deleted");
      setConfirmDelete(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const onConfirmRevoke = async () => {
    if (!confirmRevoke) return;
    try {
      const res = await revoke.mutateAsync(confirmRevoke.eoa_address);
      toast.success(`Revoked ${res.revoked} session(s)`);
      setConfirmRevoke(null);
    } catch (err) {
      toast.error(
        err instanceof AdminApiError ? err.message : (err as Error).message,
      );
    }
  };

  const columns: Column<MemberView>[] = [
    {
      key: "eoa",
      header: "EOA",
      render: (m) => <AddressDisplay value={m.eoa_address} />,
    },
    {
      key: "role",
      header: "Role",
      render: (m) => <RoleBadge role={m.role} />,
    },
    {
      key: "attributes",
      header: "Attributes",
      render: (m) =>
        m.attributes ? (
          <span className="text-xs text-surge-muted">
            <span
              className={`font-medium ${
                m.attributes.kyc ? "text-emerald-700" : "text-surge-muted"
              }`}
            >
              KYC {m.attributes.kyc ? "✓" : "✗"}
            </span>
            <span className="mx-2 text-surge-border">·</span>
            <span
              className={`font-medium ${
                m.attributes.blacklisted ? "text-red-600" : "text-surge-muted"
              }`}
            >
              Blacklist {m.attributes.blacklisted ? "✓" : "✗"}
            </span>
          </span>
        ) : (
          <span className="text-xs text-surge-muted">—</span>
        ),
    },
    {
      key: "created_at",
      header: "Created",
      render: (m) => (
        <span className="text-xs text-surge-muted">{timeAgo(m.created_at)}</span>
      ),
    },
    {
      key: "actions",
      header: "",
      align: "right",
      render: (m) => {
        const isSelf =
          session?.eoa.toLowerCase() === m.eoa_address.toLowerCase();
        return (
          <div className="flex items-center justify-end gap-1">
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                onEdit(m);
              }}
              className="rounded-md px-2 py-1 text-xs font-medium text-surge-secondary hover:bg-surge-card-hover"
            >
              Edit
            </button>
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                setConfirmRevoke(m);
              }}
              className="rounded-md px-2 py-1 text-xs font-medium text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
            >
              Revoke tokens
            </button>
            <button
              type="button"
              disabled={isSelf}
              title={isSelf ? "Cannot delete yourself" : undefined}
              onClick={(e) => {
                e.stopPropagation();
                if (!isSelf) setConfirmDelete(m);
              }}
              className="rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent"
            >
              Delete
            </button>
          </div>
        );
      },
    },
  ];

  return (
    <div>
      <PageHeader
        title="Members"
        description="Authenticated accounts on this proxy: admins and users."
        actions={
          <button
            type="button"
            onClick={onCreate}
            className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90"
          >
            + Add member
          </button>
        }
      />

      <div className="mb-4 flex flex-wrap items-center gap-2">
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by EOA…"
          className="w-64 rounded-lg border border-surge-border bg-surge-card px-3 py-2 text-sm outline-none focus:border-surge-secondary"
        />
        <select
          value={roleFilter}
          onChange={(e) => setRoleFilter(e.target.value)}
          className="rounded-lg border border-surge-border bg-surge-card px-3 py-2 text-sm outline-none focus:border-surge-secondary"
        >
          <option value="">All roles</option>
          {roles.data?.map((r) => (
            <option key={r.id} value={r.name}>
              {r.name}
            </option>
          ))}
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
        rows={filtered}
        rowKey={(m) => m.eoa_address}
        onRowClick={onEdit}
        loading={list.isLoading}
        loadingNode={<TableSkeleton />}
        emptyState={
          <EmptyState
            title="No members yet"
            description="Sign in with another wallet, or click 'Add member' to register one manually."
            action={
              <button
                type="button"
                onClick={onCreate}
                className="rounded-lg bg-surge-primary px-3 py-2 text-sm font-semibold text-white"
              >
                + Add member
              </button>
            }
          />
        }
      />

      <MemberDrawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        editingEoa={editingEoa}
      />

      <ConfirmDialog
        open={!!confirmDelete}
        title="Delete member"
        description={`This removes the member's row and revokes every active token. The member can sign in again as 'user' afterwards.`}
        confirmString={confirmDelete?.eoa_address}
        destructive
        confirmLabel="Delete"
        onConfirm={onConfirmDelete}
        onCancel={() => setConfirmDelete(null)}
        isLoading={del.isPending}
      />

      <ConfirmDialog
        open={!!confirmRevoke}
        title="Revoke all sessions"
        description="Drops every active bearer token for this EOA. The member will need to re-sign-in. Their role and attributes are preserved."
        confirmString={confirmRevoke?.eoa_address}
        confirmLabel="Revoke"
        onConfirm={onConfirmRevoke}
        onCancel={() => setConfirmRevoke(null)}
        isLoading={revoke.isPending}
      />
    </div>
  );
}
