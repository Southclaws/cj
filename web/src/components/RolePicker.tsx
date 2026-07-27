import { useState } from "react";

import { Combobox } from "./ui/Combobox";
import { useRoleOptions } from "../lib/entitySearch";

interface RolePickerProps {
  value: string;
  onChange: (id: string) => void;
  placeholder?: string;
}

export function RolePicker({ value, onChange, placeholder = "Search roles..." }: RolePickerProps) {
  const [query, setQuery] = useState("");
  const { options, all, isLoading } = useRoleOptions(query);

  const selected = value ? (all.find((option) => option.id === value) ?? { id: value, label: value }) : null;

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
