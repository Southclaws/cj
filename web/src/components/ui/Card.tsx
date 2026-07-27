import type { ReactNode } from "react";

interface CardProps {
  children: ReactNode;
  className?: string;
  interactive?: boolean;
}

export function Card({ children, className = "", interactive = false }: CardProps) {
  return (
    <div
      className={[
        "rounded-xl border border-border bg-surface-2 shadow-[0_1px_0_0_rgba(255,255,255,0.02)_inset]",
        interactive ? "transition-colors hover:border-border-strong hover:bg-surface-3" : "",
        className,
      ].join(" ")}
    >
      {children}
    </div>
  );
}

interface CardHeaderProps {
  title: string;
  description?: string;
  icon?: ReactNode;
  action?: ReactNode;
}

export function CardHeader({ title, description, icon, action }: CardHeaderProps) {
  return (
    <div className="flex items-start justify-between gap-3 border-b border-border px-5 py-4">
      <div className="flex items-start gap-2.5">
        {icon && <span className="mt-0.5 text-ink-faint">{icon}</span>}
        <div>
          <h2 className="text-xs font-semibold uppercase tracking-wider text-ink-muted">{title}</h2>
          {description && <p className="mt-1 text-sm text-ink-faint">{description}</p>}
        </div>
      </div>
      {action}
    </div>
  );
}
