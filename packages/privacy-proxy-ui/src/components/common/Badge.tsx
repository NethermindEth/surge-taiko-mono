import type { ReactNode } from "react";

type Tone =
  | "neutral"
  | "primary"
  | "secondary"
  | "mint"
  | "amber"
  | "peach"
  | "lavender"
  | "aqua";

const TONE_STYLES: Record<Tone, string> = {
  neutral: "bg-surge-card-hover text-surge-text border-surge-border",
  primary: "bg-surge-primary/10 text-surge-primary border-surge-primary/20",
  secondary:
    "bg-surge-secondary/15 text-surge-primary border-surge-secondary/30",
  mint: "bg-surge-mint/25 text-emerald-800 border-surge-mint/40",
  amber: "bg-surge-amber/20 text-amber-900 border-surge-amber/40",
  peach: "bg-surge-peach/30 text-orange-900 border-surge-peach/50",
  lavender:
    "bg-surge-lavender/20 text-surge-primary border-surge-lavender/40",
  aqua: "bg-surge-aqua/30 text-teal-900 border-surge-aqua/50",
};

interface BadgeProps {
  children: ReactNode;
  tone?: Tone;
  className?: string;
}

export function Badge({ children, tone = "neutral", className = "" }: BadgeProps) {
  return (
    <span className={`pill ${TONE_STYLES[tone]} ${className}`}>{children}</span>
  );
}

export function RoleBadge({ role }: { role: string }) {
  const tone: Tone =
    role === "admin" ? "primary" : role === "user" ? "secondary" : "lavender";
  return <Badge tone={tone}>{role}</Badge>;
}

export function ModeBadge({ mode }: { mode: "allow" | "deny" }) {
  return (
    <Badge tone={mode === "allow" ? "mint" : "amber"}>{mode}</Badge>
  );
}
