import { X } from "lucide-react";

import { SearchInput } from "./SearchInput";

export interface ComboboxOption {
  id: string;
  label: string;
  sublabel?: string;
}

interface ComboboxProps {
  selected: ComboboxOption | null;
  onSelect: (option: ComboboxOption | null) => void;
  query: string;
  onQueryChange: (query: string) => void;
  options: ComboboxOption[];
  placeholder?: string;
  isLoading?: boolean;
  inputId?: string;
}

export function Combobox({
  selected,
  onSelect,
  query,
  onQueryChange,
  options,
  placeholder,
  isLoading,
  inputId,
}: ComboboxProps) {
  if (selected) {
    return (
      <div className="flex items-center gap-2 rounded-lg border border-accent/30 bg-accent-soft px-3 py-1.5 text-sm text-ink">
        <span className="truncate">{selected.label}</span>
        {selected.sublabel && <span className="shrink-0 font-mono text-xs text-ink-faint">{selected.sublabel}</span>}
        <button
          type="button"
          onClick={() => onSelect(null)}
          className="ml-auto shrink-0 text-ink-faint transition-colors hover:text-ink"
        >
          <X className="h-3.5 w-3.5" />
        </button>
      </div>
    );
  }

  return (
    <div className="relative">
      <SearchInput
        id={inputId}
        placeholder={placeholder}
        value={query}
        onChange={(e) => onQueryChange(e.target.value)}
        className="max-w-none"
      />

      {query.length > 0 && options.length > 0 && (
        <ul className="absolute z-20 mt-1 max-h-60 w-full overflow-y-auto rounded-lg border border-border-strong bg-surface-2 text-sm shadow-xl">
          {options.map((option) => (
            <li key={option.id}>
              <button
                type="button"
                onClick={() => onSelect(option)}
                className="block w-full truncate px-3 py-1.5 text-left text-ink-muted hover:bg-surface-3 hover:text-ink"
              >
                {option.label}
                {option.sublabel && <span className="ml-1.5 font-mono text-xs text-ink-faint">({option.sublabel})</span>}
              </button>
            </li>
          ))}
        </ul>
      )}

      {query.length > 0 && options.length === 0 && (
        <div className="absolute z-20 mt-1 w-full rounded-lg border border-border-strong bg-surface-2 px-3 py-1.5 text-sm text-ink-faint shadow-xl">
          {isLoading ? "Searching..." : "No matches."}
        </div>
      )}
    </div>
  );
}
