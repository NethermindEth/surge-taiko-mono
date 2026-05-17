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
import {
  clearSession,
  readSession,
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
  const [session, setSessionState] = useState<Session | null>(() =>
    readSession(),
  );

  // Re-check expiry every 30s so the topbar countdown stays honest and the
  // session is dropped automatically when the TTL elapses.
  useEffect(() => {
    const id = setInterval(() => {
      if (session && session.expiresAt <= Math.floor(Date.now() / 1000)) {
        clearSession();
        setSessionState(null);
      }
    }, 30_000);
    return () => clearInterval(id);
  }, [session]);

  // Register the 401 handler synchronously during the first render so it's
  // wired before any child component can fire a query. apiClient reads the
  // bearer straight from localStorage, so there's no token-getter race.
  const initialized = useRef(false);
  if (!initialized.current) {
    initialized.current = true;
    setUnauthorizedHandler(() => {
      clearSession();
      setSessionState(null);
    });
  }

  const setSession = useCallback((s: Session) => {
    writeSession(s);
    setSessionState(s);
  }, []);

  const signOut = useCallback(() => {
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
