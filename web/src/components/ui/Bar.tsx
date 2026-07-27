interface RelativeBarProps {
  value: number;
  max: number;
  color?: string;
}

export function RelativeBar({ value, max, color }: RelativeBarProps) {
  const pct = max > 0 ? Math.min(Math.max((value / max) * 100, value > 0 ? 3 : 0), 100) : 0;
  return (
    <div className="h-1.5 w-full overflow-hidden rounded-full bg-surface-3">
      <div
        className="animate-bar-grow h-full rounded-full bg-accent"
        style={{ width: `${pct}%`, backgroundColor: color }}
      />
    </div>
  );
}
