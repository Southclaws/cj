import { useEffect, useMemo, useRef, useState, type KeyboardEvent as ReactKeyboardEvent } from "react";
import { useNavigate } from "react-router";
import { ArrowRight, CornerDownLeft, Search, User as UserIcon } from "lucide-react";

import { useMemberAndUserSearch } from "../lib/entitySearch";
import { navItems } from "../lib/navigation";

interface PaletteEntry {
  key: string;
  label: string;
  sublabel?: string;
  icon: typeof ArrowRight;
  go: () => void;
}

export function CommandPalette() {
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [activeIndex, setActiveIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  const { options } = useMemberAndUserSearch(open ? query : "");

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setOpen((prev) => !prev);
        return;
      }
      if (e.key === "Escape" && open) {
        setOpen(false);
      }
    }
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [open]);

  useEffect(() => {
    if (open) {
      setQuery("");
      setActiveIndex(0);
      requestAnimationFrame(() => inputRef.current?.focus());
    }
  }, [open]);

  useEffect(() => {
    setActiveIndex(0);
  }, [query]);

  const entries = useMemo((): PaletteEntry[] => {
    const q = query.trim().toLowerCase();
    const pages: PaletteEntry[] = navItems
      .filter((item) => q === "" || item.label.toLowerCase().includes(q))
      .map((item) => ({
        key: `page-${item.to}`,
        label: item.label,
        sublabel: "Page",
        icon: item.icon,
        go: () => navigate(item.to),
      }));

    const people: PaletteEntry[] =
      q === ""
        ? []
        : options.map((option) => ({
            key: `user-${option.id}`,
            label: option.label,
            sublabel: option.sublabel ?? "Profile",
            icon: UserIcon,
            go: () => navigate(`/users/${option.id}`),
          }));

    return [...pages, ...people];
  }, [query, options, navigate]);

  function activate(entry: PaletteEntry) {
    entry.go();
    setOpen(false);
  }

  function handleInputKeyDown(e: ReactKeyboardEvent<HTMLInputElement>) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIndex((i) => Math.min(i + 1, entries.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIndex((i) => Math.max(i - 1, 0));
    } else if (e.key === "Enter") {
      e.preventDefault();
      const entry = entries[activeIndex];
      if (entry) activate(entry);
    }
  }

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-[70] flex items-start justify-center bg-black/70 p-4 pt-[12vh] backdrop-blur-sm"
      onClick={() => setOpen(false)}
    >
      <div
        className="animate-palette-in w-full max-w-lg overflow-hidden rounded-xl border border-border-strong bg-surface-2 shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center gap-2.5 border-b border-border px-4 py-3">
          <Search className="h-4 w-4 shrink-0 text-ink-faint" />
          <input
            ref={inputRef}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleInputKeyDown}
            placeholder="Jump to a page or search a user..."
            className="w-full bg-transparent text-sm text-ink placeholder:text-ink-faint focus:outline-none"
          />
          <kbd className="shrink-0 rounded border border-border-strong px-1.5 py-0.5 font-mono text-[10px] text-ink-faint">
            esc
          </kbd>
        </div>

        <div className="max-h-80 overflow-y-auto p-1.5">
          {entries.length === 0 && <div className="px-3 py-6 text-center text-sm text-ink-faint">No matches.</div>}
          {entries.map((entry, index) => {
            const Icon = entry.icon;
            const active = index === activeIndex;
            return (
              <button
                key={entry.key}
                type="button"
                onMouseEnter={() => setActiveIndex(index)}
                onClick={() => activate(entry)}
                className={`flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-left text-sm transition-colors ${
                  active ? "bg-accent-soft text-ink" : "text-ink-muted"
                }`}
              >
                <Icon className="h-4 w-4 shrink-0 text-ink-faint" strokeWidth={1.75} />
                <span className="flex-1 truncate">{entry.label}</span>
                {entry.sublabel && <span className="shrink-0 font-mono text-xs text-ink-faint">{entry.sublabel}</span>}
                {active && <CornerDownLeft className="h-3.5 w-3.5 shrink-0 text-ink-faint" />}
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
}
