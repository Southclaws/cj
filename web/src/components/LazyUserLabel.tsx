import { useQuery } from "@tanstack/react-query";

import { fetchUsers } from "../lib/api";

export function LazyUserLabel({ id }: { id: string }) {
  const { data } = useQuery({
    queryKey: ["lazy-user-label", id],
    queryFn: () => fetchUsers(1, 0, id),
    staleTime: 60000,
  });

  const match = data?.users.find((user) => user.discord_user_id === id);
  const label = match ? match.forum_user_name || match.burger_user_name || id : id;

  return <span className={match ? "" : "font-mono"}>{label}</span>;
}
