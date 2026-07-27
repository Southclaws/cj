import type { ReactNode } from "react";
import { ArrowDown, ArrowUp, ArrowUpDown } from "lucide-react";

import type { SortDir } from "../../lib/useSort";

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
      <tr className="sticky top-0 z-10 border-b border-border bg-surface-2 text-xs font-medium uppercase tracking-wider text-ink-faint">
        {children}
      </tr>
    </thead>
  );
}

interface ThProps {
  children?: ReactNode;
  onSort?: () => void;
  sortDir?: SortDir | null;
}

export function Th({ children, onSort, sortDir }: ThProps) {
  if (!onSort) {
    return <th className="px-4 py-3 font-medium">{children}</th>;
  }

  const Icon = sortDir === "asc" ? ArrowUp : sortDir === "desc" ? ArrowDown : ArrowUpDown;

  return (
    <th className="px-4 py-3 font-medium">
      <button
        type="button"
        onClick={onSort}
        className={`flex items-center gap-1 transition-colors hover:text-ink ${sortDir ? "text-ink" : ""}`}
      >
        {children}
        <Icon className={`h-3 w-3 shrink-0 ${sortDir ? "" : "opacity-40"}`} />
      </button>
    </th>
  );
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
