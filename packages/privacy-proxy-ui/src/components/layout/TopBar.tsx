import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useAuth } from "../../context/AuthContext";
import { AddressDisplay } from "../common/AddressDisplay";
import { durationUntil } from "../../lib/format";
import { request, AdminApiError } from "../../lib/apiClient";

export function TopBar() {
  const { session, signOut } = useAuth();
  const [, force] = useState(0);

  // 1Hz tick so the expiry pill counts down without re-rendering the world.
  useEffect(() => {
    const id = setInterval(() => force((x) => x + 1), 1000);
    return () => clearInterval(id);
  }, []);

  if (!session) return null;

  const handleSignOut = async () => {
    try {
      await request<{ revoked: number }>(
        `/admin/members/${session.eoa}/tokens`,
        { method: "DELETE" },
      );
    } catch (err) {
      if (err instanceof AdminApiError && err.status !== 401) {
        toast.error(`Best-effort revoke failed: ${err.message}`);
      }
    } finally {
      signOut();
      toast.success("Signed out");
    }
  };

  return (
    <header className="flex items-center justify-between gap-4 border-b border-surge-border bg-surge-card/70 px-6 py-3 backdrop-blur-md">
      <div className="flex flex-col">
        <span className="text-xs font-medium uppercase tracking-wide text-surge-muted">
          Signed in
        </span>
        <AddressDisplay value={session.eoa} className="text-sm" />
      </div>
      <div className="flex items-center gap-3">
        <span
          className="pill border-surge-border bg-surge-card-hover text-surge-muted"
          title="Bearer token expiry"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
            <circle cx="12" cy="12" r="9" />
            <path d="M12 7v5l3 2" />
          </svg>
          Auth token expires in {durationUntil(session.expiresAt)}
        </span>
        <button
          type="button"
          onClick={handleSignOut}
          className="rounded-lg border border-surge-border bg-surge-card px-3 py-1.5 text-sm font-medium text-surge-text hover:bg-surge-card-hover"
        >
          Sign out
        </button>
      </div>
    </header>
  );
}
