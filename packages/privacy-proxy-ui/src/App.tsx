import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RainbowKitProvider, lightTheme } from "@rainbow-me/rainbowkit";
import { WagmiProvider } from "wagmi";
import {
  BrowserRouter,
  Navigate,
  Route,
  Routes,
} from "react-router-dom";
import { Toaster } from "react-hot-toast";
import { wagmiConfig } from "./lib/config";
import { AuthProvider } from "./context/AuthContext";
import { LoginGate } from "./components/auth/LoginGate";
import { RequireAdmin } from "./components/auth/RequireAdmin";
import { AppShell } from "./components/layout/AppShell";
import { DashboardPage } from "./pages/DashboardPage";
import { MembersPage } from "./pages/MembersPage";
import { RulesPage } from "./pages/RulesPage";
import { RuleDetailPage } from "./pages/RuleDetailPage";
import { RolesPage } from "./pages/RolesPage";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        // Don't retry on auth failures — useless.
        if (
          error instanceof Error &&
          error.name === "AdminApiError" &&
          // @ts-expect-error optional status field on subclass
          (error.status === 401 || error.status === 403)
        ) {
          return false;
        }
        return failureCount < 1;
      },
      staleTime: 5_000,
    },
  },
});

export default function App() {
  return (
    <WagmiProvider config={wagmiConfig}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider
          theme={lightTheme({
            accentColor: "#172342",
            accentColorForeground: "#ffffff",
            borderRadius: "medium",
          })}
        >
          <AuthProvider>
            <BrowserRouter>
              <Routes>
                <Route path="/login" element={<LoginGate />} />
                <Route
                  element={
                    <RequireAdmin>
                      <AppShell />
                    </RequireAdmin>
                  }
                >
                  <Route index element={<DashboardPage />} />
                  <Route path="members" element={<MembersPage />} />
                  <Route path="rules" element={<RulesPage />} />
                  <Route path="rules/:id" element={<RuleDetailPage />} />
                  <Route path="roles" element={<RolesPage />} />
                </Route>
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </BrowserRouter>
            <Toaster
              position="top-right"
              toastOptions={{
                style: {
                  borderRadius: "0.75rem",
                  background: "#ffffff",
                  color: "#172342",
                  border: "1px solid rgba(23, 35, 66, 0.08)",
                  boxShadow:
                    "0 8px 24px rgba(23, 35, 66, 0.08), 0 1px 2px rgba(23, 35, 66, 0.04)",
                  fontSize: "0.875rem",
                },
              }}
            />
          </AuthProvider>
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}
