import type { ReactNode } from "react";

type Tone = "neutral" | "accent" | "good" | "warn" | "bad";

const toneClasses: Record<Tone, string> = {
  neutral: "border-border-strong bg-surface-3 text-ink-muted",
  accent: "border-accent/30 bg-accent-soft text-accent-ink",
  good: "border-good/30 bg-good-soft text-good",
  warn: "border-warn/30 bg-warn-soft text-warn",
  bad: "border-bad/30 bg-bad-soft text-bad",
};

interface BadgeProps {
  children: ReactNode;
  tone?: Tone;
  icon?: ReactNode;
}

export function Badge({ children, tone = "neutral", icon }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium ${toneClasses[tone]}`}
    >
      {icon}
      {children}
    </span>
  );
}

const riskTones: Record<string, Tone> = {
  "read-only": "neutral",
  safe: "good",
  "state-changing": "warn",
  dangerous: "bad",
};

export function riskTone(risk: string): Tone {
  return riskTones[risk] ?? "neutral";
}
