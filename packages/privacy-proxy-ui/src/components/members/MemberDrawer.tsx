import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { Drawer } from "../common/Drawer";
import {
  useGetMember,
  useUpsertMember,
} from "../../hooks/members/useMembers";
import { isAddress, normalizeAddress } from "../../lib/format";
import { AdminApiError } from "../../lib/apiClient";
import type { RoleName, UpsertMemberRequest } from "../../types/api";

interface MemberDrawerProps {
  open: boolean;
  onClose: () => void;
  /** When provided, drawer is in edit mode. */
  editingEoa?: string;
}

interface FormState {
  eoa: string;
  role: RoleName;
  kyc: boolean;
  blacklisted: boolean;
}

const EMPTY: FormState = {
  eoa: "",
  role: "user",
  kyc: false,
  blacklisted: false,
};

export function MemberDrawer({ open, onClose, editingEoa }: MemberDrawerProps) {
  const isEdit = !!editingEoa;
  const detail = useGetMember(isEdit ? editingEoa : undefined);
  const upsert = useUpsertMember();
  const [form, setForm] = useState<FormState>(EMPTY);
  const [eoaError, setEoaError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    if (isEdit && detail.data) {
      setForm({
        eoa: detail.data.eoa_address,
        role: detail.data.role,
        kyc: detail.data.attributes?.kyc ?? false,
        blacklisted: detail.data.attributes?.blacklisted ?? false,
      });
    } else if (!isEdit) {
      setForm(EMPTY);
    }
    setEoaError(null);
  }, [open, isEdit, detail.data]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isEdit) {
      if (!isAddress(form.eoa)) {
        setEoaError("Enter a 20-byte hex address (0x + 40 chars)");
        return;
      }
    }

    const body: UpsertMemberRequest =
      form.role === "admin"
        ? { role: "admin" }
        : {
            role: "user",
            attributes: { kyc: form.kyc, blacklisted: form.blacklisted },
          };

    try {
      await upsert.mutateAsync({
        eoa: normalizeAddress(form.eoa),
        body,
      });
      toast.success(isEdit ? "Member updated" : "Member created");
      onClose();
    } catch (err) {
      const msg =
        err instanceof AdminApiError ? err.message : (err as Error).message;
      toast.error(msg);
    }
  };

  return (
    <Drawer
      open={open}
      onClose={onClose}
      title={isEdit ? "Edit member" : "Add member"}
      subtitle={
        isEdit
          ? "Update this EOA's role and typed attributes."
          : "Register a new EOA on this proxy."
      }
      footer={
        <div className="flex items-center justify-end gap-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg px-3 py-2 text-sm font-medium text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
          >
            Cancel
          </button>
          <button
            type="submit"
            form="member-form"
            disabled={upsert.isPending}
            className="rounded-lg bg-surge-primary px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-surge-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {upsert.isPending ? "Saving…" : isEdit ? "Save changes" : "Create"}
          </button>
        </div>
      }
    >
      <form id="member-form" onSubmit={onSubmit} className="space-y-5">
        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            EOA address
          </label>
          <input
            type="text"
            value={form.eoa}
            onChange={(e) => setForm({ ...form, eoa: e.target.value })}
            placeholder="0x..."
            disabled={isEdit}
            className="mt-1 w-full rounded-lg border border-surge-border bg-surge-card-hover/40 px-3 py-2 font-mono text-sm outline-none focus:border-surge-secondary disabled:opacity-70"
          />
          {eoaError ? (
            <p className="mt-1 text-xs text-red-600">{eoaError}</p>
          ) : null}
        </div>

        <div>
          <label className="block text-xs font-medium uppercase tracking-wide text-surge-muted">
            Role
          </label>
          <div className="mt-1 grid grid-cols-2 gap-2">
            {(["admin", "user"] as const).map((r) => (
              <button
                type="button"
                key={r}
                onClick={() => setForm({ ...form, role: r })}
                className={`rounded-lg border px-3 py-2 text-sm font-medium transition ${
                  form.role === r
                    ? "border-surge-primary bg-surge-primary text-white"
                    : "border-surge-border bg-surge-card text-surge-text hover:bg-surge-card-hover"
                }`}
              >
                {r}
              </button>
            ))}
          </div>
        </div>

        {form.role === "user" ? (
          <div className="rounded-xl border border-surge-border bg-surge-card-hover/30 p-4">
            <p className="text-xs font-medium uppercase tracking-wide text-surge-muted">
              User attributes
            </p>
            <p className="mt-1 text-xs text-surge-muted">
              Admin-managed flags. Omitted fields on update preserve the
              stored value.
            </p>
            <div className="mt-3 space-y-2">
              <label className="flex cursor-pointer items-center justify-between rounded-lg bg-surge-card px-3 py-2 text-sm">
                <span>
                  <span className="font-medium text-surge-text">KYC</span>
                  <span className="ml-2 text-xs text-surge-muted">
                    Identity verified
                  </span>
                </span>
                <input
                  type="checkbox"
                  checked={form.kyc}
                  onChange={(e) =>
                    setForm({ ...form, kyc: e.target.checked })
                  }
                  className="h-4 w-4 accent-surge-primary"
                />
              </label>
              <label className="flex cursor-pointer items-center justify-between rounded-lg bg-surge-card px-3 py-2 text-sm">
                <span>
                  <span className="font-medium text-surge-text">Blacklisted</span>
                  <span className="ml-2 text-xs text-surge-muted">
                    Hard-blocked from access
                  </span>
                </span>
                <input
                  type="checkbox"
                  checked={form.blacklisted}
                  onChange={(e) =>
                    setForm({ ...form, blacklisted: e.target.checked })
                  }
                  className="h-4 w-4 accent-surge-primary"
                />
              </label>
            </div>
          </div>
        ) : (
          <div className="rounded-xl border border-surge-border bg-surge-card-hover/30 p-4 text-sm text-surge-muted">
            Admin role is identity-only — no attribute fields. Switching this
            member to <span className="font-medium text-surge-text">admin</span>{" "}
            removes any stored user attributes.
          </div>
        )}
      </form>
    </Drawer>
  );
}
