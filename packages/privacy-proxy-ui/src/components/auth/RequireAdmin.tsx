import { useEffect } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAccount } from "wagmi";
import { useAuth } from "../../context/AuthContext";

/**
 * Route guard. Redirects to /login when there's no live session, and signs
 * the user out if their connected wallet differs from the EOA that issued
 * the current token.
 */
export function RequireAdmin({ children }: { children: React.ReactNode }) {
  const { session, signOut } = useAuth();
  const { address } = useAccount();
  const location = useLocation();

  // If the user switched accounts in the wallet, drop the session — the
  // token belongs to the previous EOA.
  useEffect(() => {
    if (!session) return;
    if (address && address.toLowerCase() !== session.eoa.toLowerCase()) {
      signOut();
    }
  }, [address, session, signOut]);

  if (!session) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return <>{children}</>;
}
