type Tone = "good" | "warn" | "bad" | "neutral";

const toneClasses: Record<Tone, string> = {
  good: "bg-good shadow-[0_0_8px_var(--color-good)]",
  warn: "bg-warn shadow-[0_0_8px_var(--color-warn)]",
  bad: "bg-bad shadow-[0_0_8px_var(--color-bad)]",
  neutral: "bg-ink-faint",
};

export function statusTone(status: string): Tone {
  switch (status) {
    case "healthy":
    case "configured":
      return "good";
    case "disabled":
    case "not_configured":
      return "neutral";
    case "unhealthy":
      return "bad";
    default:
      return "warn";
  }
}

interface StatusDotProps {
  tone: Tone;
  pulse?: boolean;
}

export function StatusDot({ tone, pulse = false }: StatusDotProps) {
  return <span className={`inline-block h-2 w-2 rounded-full ${toneClasses[tone]} ${pulse ? "animate-pulse-dot" : ""}`} />;
}
