import { useState } from "react";

import { Combobox } from "./ui/Combobox";
import { useMemberAndUserSearch } from "../lib/entitySearch";

interface MemberUserPickerProps {
  value: string;
  onChange: (id: string) => void;
  placeholder?: string;
}

export function MemberUserPicker({ value, onChange, placeholder = "Search by name or ID..." }: MemberUserPickerProps) {
  const [query, setQuery] = useState("");
  const { options, isLoading } = useMemberAndUserSearch(query);
  const { options: resolvedOptions } = useMemberAndUserSearch(value);

  const selectedFromResults = options.find((option) => option.id === value) ?? resolvedOptions.find((option) => option.id === value);
  const selected = value ? (selectedFromResults ?? { id: value, label: value }) : null;

  return (
    <Combobox
      selected={selected}
      onSelect={(option) => {
        onChange(option?.id ?? "");
        setQuery("");
      }}
      query={query}
      onQueryChange={setQuery}
      options={options}
      placeholder={placeholder}
      isLoading={isLoading}
    />
  );
}
