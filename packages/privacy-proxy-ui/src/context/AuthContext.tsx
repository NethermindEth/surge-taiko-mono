import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { useAccount } from "wagmi";
import {
  clearActiveSession,
  clearSession,
  clearSessionFor,
  getSessionFor,
  readSession,
  setActiveSession,
  writeSession,
  type Session,
} from "../lib/tokenStorage";
import { setUnauthorizedHandler } from "../lib/apiClient";

interface AuthContextValue {
  session: Session | null;
  setSession: (s: Session) => void;
  signOut: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const { address } = useAccount();
  const [session, setSessionState] = useState<Session | null>(() =>
    readSession(),
  );

  // Drop the session automatically when its TTL elapses.
  useEffect(() => {
    const id = setInterval(() => {
      if (session && session.expiresAt <= Math.floor(Date.now() / 1000)) {
        clearSessionFor(session.eoa);
        clearActiveSession();
        setSessionState(null);
      }
    }, 30_000);
    return () => clearInterval(id);
  }, [session]);

  // When the connected wallet address changes, swap in that EOA's
  // cached session (if any) instead of forcing a re-sign-in.
  useEffect(() => {
    if (!address) {
      // Wallet disconnected — drop the active mirror so apiClient stops
      // sending the bearer, but keep the cache so a future re-connect
      // picks the session right back up.
      if (session) {
        clearActiveSession();
        setSessionState(null);
      }
      return;
    }
    if (session && session.eoa.toLowerCase() === address.toLowerCase()) return;

    const cached = getSessionFor(address);
    if (cached) {
      setActiveSession(cached);
      setSessionState(cached);
    } else {
      // No session for the newly-connected EOA — clear the mirror so
      // LoginGate takes over, but do NOT touch other accounts' cached
      // sessions in case the user flips back.
      clearActiveSession();
      setSessionState(null);
    }
  }, [address, session]);

  // Register the 401 handler synchronously during the first render so
  // it's wired before any child can fire a query. apiClient reads the
  // bearer straight from localStorage, so there's no token-getter race.
  const initialized = useRef(false);
  if (!initialized.current) {
    initialized.current = true;
    setUnauthorizedHandler(() => {
      // Don't blow away every cached EOA's session — just the active
      // one. The 401 implies the proxy rejected this token specifically.
      const active = readSession();
      if (active) clearSessionFor(active.eoa);
      else clearActiveSession();
      setSessionState(null);
    });
  }

  const setSession = useCallback((s: Session) => {
    writeSession(s);
    setSessionState(s);
  }, []);

  const signOut = useCallback(() => {
    // Explicit user sign-out: drop EVERYTHING (cache + active mirror).
    clearSession();
    setSessionState(null);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({ session, setSession, signOut }),
    [session, setSession, signOut],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
