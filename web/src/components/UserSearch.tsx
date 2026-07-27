import { useState } from "react";
import { useNavigate } from "react-router";

import { Combobox } from "./ui/Combobox";
import { useMemberAndUserSearch } from "../lib/entitySearch";

const snowflakePattern = /^[0-9]{15,20}$/;

export function UserSearch() {
  const navigate = useNavigate();
  const [query, setQuery] = useState("");
  const { options } = useMemberAndUserSearch(query);

  function goToProfile(id: string) {
    setQuery("");
    navigate(`/users/${id}`);
  }

  function handleEnter() {
    const trimmed = query.trim();
    if (trimmed === "") return;

    if (snowflakePattern.test(trimmed)) {
      goToProfile(trimmed);
      return;
    }

    const first = options[0];
    if (first) {
      goToProfile(first.id);
    }
  }

  return (
    <div className="relative px-3" onKeyDown={(e) => e.key === "Enter" && handleEnter()}>
      <Combobox
        selected={null}
        onSelect={(option) => option && goToProfile(option.id)}
        query={query}
        onQueryChange={setQuery}
        options={options}
        placeholder="Search users..."
        inputId="global-user-search"
      />
      {query === "" && (
        <kbd className="pointer-events-none absolute right-5 top-1/2 -translate-y-1/2 rounded border border-border-strong bg-surface px-1.5 py-0.5 font-mono text-[10px] text-ink-faint">
          /
        </kbd>
      )}
    </div>
  );
}
