import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";

import type { ComboboxOption } from "../components/ui/Combobox";
import { useDebounced } from "./useDebounced";
import { fetchChannels, fetchMembers, fetchRoles, fetchUsers } from "./api";

const CHANNEL_ROLE_STALE_TIME = 5 * 60 * 1000;
const MAX_RESULTS = 8;

export function useMemberAndUserSearch(query: string) {
  const debounced = useDebounced(query, 300);
  const enabled = debounced.length > 0;

  const members = useQuery({
    queryKey: ["entity-search-members", debounced],
    queryFn: () => fetchMembers(MAX_RESULTS, 0, debounced),
    enabled,
  });
  const users = useQuery({
    queryKey: ["entity-search-cj-users", debounced],
    queryFn: () => fetchUsers(MAX_RESULTS, 0, debounced),
    enabled,
  });

  const options = useMemo(() => {
    const seen = new Set<string>();
    const out: ComboboxOption[] = [];

    for (const member of members.data?.members ?? []) {
      seen.add(member.id);
      out.push({
        id: member.id,
        label: member.nickname ? `${member.nickname} (${member.username})` : member.username,
        sublabel: member.id,
      });
    }
    for (const user of users.data?.users ?? []) {
      if (seen.has(user.discord_user_id)) continue;
      seen.add(user.discord_user_id);
      out.push({
        id: user.discord_user_id,
        label: user.forum_user_name || user.burger_user_name || user.discord_user_id,
        sublabel: user.discord_user_id,
      });
    }
    return out;
  }, [members.data, users.data]);

  return { options, isLoading: members.isLoading || users.isLoading };
}

export function useChannelOptions(query: string) {
  const { data, isLoading } = useQuery({
    queryKey: ["entity-search-channels"],
    queryFn: fetchChannels,
    staleTime: CHANNEL_ROLE_STALE_TIME,
  });

  const all = useMemo(() => {
    const out: ComboboxOption[] = [];
    for (const category of data ?? []) {
      for (const channel of category.channels) {
        out.push({ id: channel.id, label: `#${channel.name}`, sublabel: channel.id });
      }
    }
    return out;
  }, [data]);

  const options = useMemo(() => {
    if (!query) return [];
    const q = query.toLowerCase();
    return all.filter((option) => option.label.toLowerCase().includes(q) || option.id.includes(q)).slice(0, MAX_RESULTS);
  }, [all, query]);

  return { options, all, isLoading };
}

export function useRoleOptions(query: string) {
  const { data, isLoading } = useQuery({
    queryKey: ["entity-search-roles"],
    queryFn: fetchRoles,
    staleTime: CHANNEL_ROLE_STALE_TIME,
  });

  const all = useMemo(() => {
    return (data ?? []).map((role): ComboboxOption => ({ id: role.id, label: role.name, sublabel: role.id }));
  }, [data]);

  const options = useMemo(() => {
    if (!query) return [];
    const q = query.toLowerCase();
    return all.filter((option) => option.label.toLowerCase().includes(q) || option.id.includes(q)).slice(0, MAX_RESULTS);
  }, [all, query]);

  return { options, all, isLoading };
}
