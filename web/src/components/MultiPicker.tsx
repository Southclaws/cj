import { useState, type ReactNode } from "react";
import { X } from "lucide-react";

import { Combobox, type ComboboxOption } from "./ui/Combobox";

interface MultiPickerProps {
  ids: string[];
  onChange: (ids: string[]) => void;
  useOptions: (query: string) => { options: ComboboxOption[]; isLoading: boolean };
  renderChipLabel?: (id: string) => ReactNode;
  placeholder?: string;
}

export function MultiPicker({ ids, onChange, useOptions, renderChipLabel, placeholder }: MultiPickerProps) {
  const [query, setQuery] = useState("");
  const { options, isLoading } = useOptions(query);
  const available = options.filter((option) => !ids.includes(option.id));

  function add(id: string) {
    if (!ids.includes(id)) onChange([...ids, id]);
    setQuery("");
  }
  function remove(id: string) {
    onChange(ids.filter((existing) => existing !== id));
  }

  return (
    <div className="flex flex-col gap-2">
      {ids.length > 0 && (
        <ul className="flex flex-wrap gap-1.5">
          {ids.map((id) => (
            <li
              key={id}
              className="flex items-center gap-1.5 rounded-full border border-border-strong bg-surface-3 px-2.5 py-1 text-xs text-ink"
            >
              <span className={renderChipLabel ? "" : "font-mono"}>{renderChipLabel ? renderChipLabel(id) : id}</span>
              <button type="button" onClick={() => remove(id)} className="text-ink-faint transition-colors hover:text-ink">
                <X className="h-3 w-3" />
              </button>
            </li>
          ))}
        </ul>
      )}
      <Combobox
        selected={null}
        onSelect={(option) => option && add(option.id)}
        query={query}
        onQueryChange={setQuery}
        options={available}
        placeholder={placeholder}
        isLoading={isLoading}
      />
    </div>
  );
}
