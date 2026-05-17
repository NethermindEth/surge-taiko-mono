import { NavLink } from "react-router-dom";
import type { ReactNode } from "react";

interface NavItem {
  to: string;
  label: string;
  icon: ReactNode;
  end?: boolean;
}

const ICONS = {
  dashboard: (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <rect x="3" y="3" width="7" height="9" rx="1" />
      <rect x="14" y="3" width="7" height="5" rx="1" />
      <rect x="14" y="12" width="7" height="9" rx="1" />
      <rect x="3" y="16" width="7" height="5" rx="1" />
    </svg>
  ),
  members: (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <circle cx="9" cy="8" r="3" />
      <path d="M2.5 20a6.5 6.5 0 0 1 13 0" />
      <circle cx="17" cy="9" r="2.5" />
      <path d="M15 20a5 5 0 0 1 7 0" />
    </svg>
  ),
  rules: (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <path d="M4 4h12l4 4v12a0 0 0 0 1 0 0H4z" />
      <path d="M8 9h6M8 13h8M8 17h5" />
    </svg>
  ),
  lambdas: (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <path d="M5 4h6l8 16h-6L9 12l-4 8H3z" />
    </svg>
  ),
  selectors: (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <path d="M7 7h10v10H7z" />
      <path d="M3 9v6M21 9v6M9 3h6M9 21h6" />
    </svg>
  ),
  roles: (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <path d="M12 2l8 4v6c0 5-3.5 8.5-8 10-4.5-1.5-8-5-8-10V6z" />
    </svg>
  ),
};

const NAV: NavItem[] = [
  { to: "/", label: "Dashboard", icon: ICONS.dashboard, end: true },
  { to: "/members", label: "Members", icon: ICONS.members },
  { to: "/rules", label: "Access rules", icon: ICONS.rules },
  { to: "/lambdas", label: "Lambdas", icon: ICONS.lambdas },
  { to: "/selectors", label: "Gated RPC endpoints", icon: ICONS.selectors },
  { to: "/roles", label: "Roles", icon: ICONS.roles },
];

export function Sidebar() {
  return (
    <aside className="hidden w-60 shrink-0 border-r border-surge-border bg-surge-card/70 backdrop-blur-md md:flex md:flex-col">
      <div className="flex flex-col items-start gap-1.5 px-5 py-5">
        <img src="/surge-logo.svg" alt="Surge" className="h-7" />
        <span className="text-sm font-semibold tracking-tight text-surge-text">
          Access Control Panel
        </span>
      </div>
      <nav className="flex-1 space-y-0.5 px-3">
        {NAV.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end}
            className={({ isActive }) =>
              `flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition ${
                isActive
                  ? "bg-surge-primary text-white shadow-sm"
                  : "text-surge-muted hover:bg-surge-card-hover hover:text-surge-text"
              }`
            }
          >
            <span>{item.icon}</span>
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>
      <div className="px-5 py-4">
        <img
          src="/powered-by-nethermind.svg"
          alt="Powered by Nethermind"
          className="h-6 opacity-60"
        />
      </div>
    </aside>
  );
}
