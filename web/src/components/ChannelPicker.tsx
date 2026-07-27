import { useState } from "react";

import { Combobox } from "./ui/Combobox";
import { useChannelOptions } from "../lib/entitySearch";

interface ChannelPickerProps {
  value: string;
  onChange: (id: string) => void;
  placeholder?: string;
}

export function ChannelPicker({ value, onChange, placeholder = "Search channels..." }: ChannelPickerProps) {
  const [query, setQuery] = useState("");
  const { options, all, isLoading } = useChannelOptions(query);

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
