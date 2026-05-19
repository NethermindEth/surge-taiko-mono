import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";

/**
 * Route guard. Redirects to /login when there's no live session for the
 * currently-connected EOA. Account-switch handling lives in
 * {@link AuthProvider}, which restores a cached session if one exists
 * for the new EOA — so this guard just observes the resulting state.
 */
export function RequireAdmin({ children }: { children: React.ReactNode }) {
  const { session } = useAuth();
  const location = useLocation();

  if (!session) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return <>{children}</>;
}
