import { Search } from "lucide-react";
import type { InputHTMLAttributes } from "react";

export function SearchInput({ className = "", ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <div className={`relative max-w-sm ${className}`}>
      <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-ink-faint" />
      <input
        type="text"
        className="w-full rounded-lg border border-border bg-surface-2 py-1.5 pl-9 pr-3 text-sm text-ink placeholder:text-ink-faint focus:border-accent/50 focus:outline-none focus:ring-2 focus:ring-accent/20"
        {...props}
      />
    </div>
  );
}
