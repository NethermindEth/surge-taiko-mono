import { ConnectButton } from "@rainbow-me/rainbowkit";
import { Navigate, useLocation } from "react-router-dom";
import { useAccount } from "wagmi";
import { useSignIn } from "../../hooks/useSignIn";
import { useAuth } from "../../context/AuthContext";
import { APP_NAME } from "../../lib/config";

export function LoginGate() {
  const { isConnected, address } = useAccount();
  const { signIn, isLoading, error } = useSignIn();
  const { session } = useAuth();
  const location = useLocation();

  if (session) {
    const from =
      (location.state as { from?: { pathname?: string } } | null)?.from
        ?.pathname ?? "/";
    return <Navigate to={from} replace />;
  }

  return (
    <div className="relative z-10 flex min-h-screen items-center justify-center px-4">
      <div className="glass-card hover-glow w-full max-w-md rounded-2xl p-8">
        <div className="mb-6 flex flex-col items-center">
          <img src="/surge-logo.svg" alt="Surge" className="mb-3 h-10" />
          <h1 className="text-xl font-semibold tracking-tight text-surge-text">
            {APP_NAME}
          </h1>
          <p className="mt-1 text-sm text-surge-muted">
            Admin access by wallet signature
          </p>
        </div>

        <div className="space-y-4">
          <div className="rounded-xl bg-surge-card-hover/60 p-4">
            <p className="text-xs font-medium uppercase tracking-wide text-surge-muted">
              Step 1
            </p>
            <p className="mt-1 text-sm text-surge-text">
              Connect the wallet whose EOA is registered as an admin on this
              proxy.
            </p>
            <div className="mt-3">
              <ConnectButton showBalance={false} chainStatus="icon" />
            </div>
          </div>

          <div
            className={`rounded-xl p-4 transition ${
              isConnected
                ? "bg-surge-card-hover/60"
                : "bg-surge-card-hover/30 opacity-60"
            }`}
          >
            <p className="text-xs font-medium uppercase tracking-wide text-surge-muted">
              Step 2
            </p>
            <p className="mt-1 text-sm text-surge-text">
              Sign the proxy's challenge. One signature unlocks the panel until
              your token expires.
            </p>
            <button
              type="button"
              onClick={signIn}
              disabled={!isConnected || isLoading}
              className="mt-3 inline-flex w-full items-center justify-center rounded-lg bg-surge-primary px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-surge-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {isLoading
                ? "Awaiting signature…"
                : isConnected
                  ? "Sign in"
                  : "Connect a wallet to continue"}
            </button>
            {error ? (
              <p className="mt-2 text-xs text-red-600">{error}</p>
            ) : null}
            {address ? (
              <p className="mt-2 text-xs text-surge-muted">
                Connected as <span className="font-mono">{address}</span>
              </p>
            ) : null}
          </div>
        </div>

        <p className="mt-6 text-center text-xs text-surge-muted">
          You'll only sign once per session. Your bearer token is held in this
          browser and refreshes on its own expiry.
        </p>
      </div>
    </div>
  );
}
