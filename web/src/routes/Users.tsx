import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import { IdCard } from "lucide-react";

import { Badge } from "../components/ui/Badge";
import { PersonTag } from "../components/ui/PersonTag";
import { Pagination } from "../components/ui/Pagination";
import { PageHeader } from "../components/ui/PageHeader";
import { SearchInput } from "../components/ui/SearchInput";
import { EmptyState, ErrorState, LoadingState } from "../components/ui/States";
import { Table, Thead, Th, Td, Tr } from "../components/ui/Table";
import { errorMessage } from "../lib/errors";
import { useDebounced } from "../lib/useDebounced";
import { fetchUsers } from "../lib/api";

const PAGE_SIZE = 20;
const snowflakePattern = /^[0-9]{15,20}$/;

export function Users() {
  const navigate = useNavigate();
  const [offset, setOffset] = useState(0);
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebounced(query, 300);

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["data-users", offset, debouncedQuery],
    queryFn: () => fetchUsers(PAGE_SIZE, offset, debouncedQuery),
  });

  function handleEnter() {
    const trimmed = query.trim();
    if (trimmed === "") return;

    if (snowflakePattern.test(trimmed)) {
      navigate(`/users/${trimmed}`);
      return;
    }

    const first = data?.users[0];
    if (first) {
      navigate(`/users/${first.discord_user_id}`);
    }
  }

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader
        icon={IdCard}
        title="Users"
        description="Search by Discord name, ID, forum name, or Burgershot name. Press Enter to jump to the top match, or click any row for the full profile."
      />

      <SearchInput
        placeholder="Search users..."
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          setOffset(0);
        }}
        onKeyDown={(e) => e.key === "Enter" && handleEnter()}
      />

      {isLoading && <LoadingState label="Loading users" />}

      {isError && <ErrorState title="Users are unavailable." message={errorMessage(error)} />}

      {data && data.total === 0 && <EmptyState message="No users match that search." />}

      {data && data.total > 0 && (
        <>
          <div className="max-w-4xl">
            <Table>
              <Thead>
                <Th>User</Th>
                <Th>Forum name</Th>
                <Th>Burgershot name</Th>
                <Th>Verified</Th>
              </Thead>
              <tbody>
                {data.users.map((user) => (
                  <Tr key={user.discord_user_id} clickable onClick={() => navigate(`/users/${user.discord_user_id}`)}>
                    <Td>
                      <PersonTag
                        person={
                          user.person ?? { id: user.discord_user_id, known: false }
                        }
                        fallbackName={user.forum_user_name || user.burger_user_name}
                      />
                    </Td>
                    <Td className="text-ink-muted">{user.forum_user_name || "none"}</Td>
                    <Td className="text-ink-muted">{user.burger_user_name || "none"}</Td>
                    <Td>
                      <Badge tone={user.burgershot_verified ? "good" : "neutral"}>
                        {user.burgershot_verified ? "Yes" : "No"}
                      </Badge>
                    </Td>
                  </Tr>
                ))}
              </tbody>
            </Table>
          </div>

          <Pagination offset={offset} pageSize={PAGE_SIZE} total={data.total} onChange={setOffset} />
        </>
      )}
    </div>
  );
}
