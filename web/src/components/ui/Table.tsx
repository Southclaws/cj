import type { ReactNode } from "react";

export function Table({ children }: { children: ReactNode }) {
  return (
    <div className="overflow-x-auto rounded-xl border border-border bg-surface-2">
      <table className="w-full text-left text-sm">{children}</table>
    </div>
  );
}

export function Thead({ children }: { children: ReactNode }) {
  return (
    <thead>
      <tr className="border-b border-border text-xs font-medium uppercase tracking-wider text-ink-faint">
        {children}
      </tr>
    </thead>
  );
}

export function Th({ children }: { children?: ReactNode }) {
  return <th className="px-4 py-3 font-medium">{children}</th>;
}

export function Td({ children, className = "" }: { children?: ReactNode; className?: string }) {
  return <td className={`px-4 py-3 ${className}`}>{children}</td>;
}

interface TrProps {
  children: ReactNode;
  onClick?: () => void;
  clickable?: boolean;
}

export function Tr({ children, onClick, clickable }: TrProps) {
  return (
    <tr
      onClick={onClick}
      className={`border-b border-border/60 last:border-b-0 hover:bg-surface-3/50 ${clickable ? "cursor-pointer" : ""}`}
    >
      {children}
    </tr>
  );
}
