import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import { Users } from "lucide-react";

import { Badge } from "../components/ui/Badge";
import { PersonTag } from "../components/ui/PersonTag";
import { Pagination } from "../components/ui/Pagination";
import { PageHeader } from "../components/ui/PageHeader";
import { SearchInput } from "../components/ui/SearchInput";
import { EmptyState, ErrorState, LoadingState } from "../components/ui/States";
import { Table, Thead, Th, Td, Tr } from "../components/ui/Table";
import { errorMessage } from "../lib/errors";
import { fetchMembers, fetchRoles } from "../lib/api";

const PAGE_SIZE = 50;

export function Members() {
  const navigate = useNavigate();
  const [offset, setOffset] = useState(0);
  const [query, setQuery] = useState("");

  const members = useQuery({
    queryKey: ["discord-members", offset, query],
    queryFn: () => fetchMembers(PAGE_SIZE, offset, query),
  });
  const roles = useQuery({ queryKey: ["discord-roles"], queryFn: fetchRoles });

  const roleNames = new Map((roles.data ?? []).map((role) => [role.id, role.name]));

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader
        icon={Users}
        title="Members"
        description={
          members.data
            ? `Search by username, nickname, or Discord ID. Member list as of ${new Date(members.data.cachedAt).toLocaleString()}.`
            : "Search by username, nickname, or Discord ID."
        }
      />

      <SearchInput
        placeholder="Search members..."
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          setOffset(0);
        }}
      />

      {members.isLoading && <LoadingState label="Loading members" />}

      {members.isError && <ErrorState title="Members are unavailable." message={errorMessage(members.error)} />}

      {members.data && members.data.total === 0 && <EmptyState message="No members match that search." />}

      {members.data && members.data.total > 0 && (
        <>
          <div className="max-w-4xl">
            <Table>
              <Thead>
                <Th>Member</Th>
                <Th>Roles</Th>
                <Th>Joined</Th>
                <Th>CJ can act on</Th>
              </Thead>
              <tbody>
                {members.data.members.map((member) => (
                  <Tr key={member.id} clickable onClick={() => navigate(`/users/${member.id}`)}>
                    <Td>
                      <PersonTag
                        person={{
                          id: member.id,
                          username: member.username,
                          nickname: member.nickname,
                          avatarUrl: member.avatarUrl,
                          known: true,
                        }}
                      />
                    </Td>
                    <Td className="text-ink-muted">
                      {member.roles.length === 0 ? "none" : member.roles.map((id) => roleNames.get(id) ?? id).join(", ")}
                    </Td>
                    <Td className="text-ink-muted">
                      {member.joinedAt ? new Date(member.joinedAt).toLocaleDateString() : "unknown"}
                    </Td>
                    <Td>
                      <Badge tone={member.cjCanActOn ? "good" : "neutral"}>{member.cjCanActOn ? "Yes" : "No"}</Badge>
                      {!member.cjCanActOn && member.reasons && member.reasons.length > 0 && (
                        <div className="mt-1 text-xs text-ink-faint">{member.reasons.join("; ")}</div>
                      )}
                    </Td>
                  </Tr>
                ))}
              </tbody>
            </Table>
          </div>

          <Pagination offset={offset} pageSize={PAGE_SIZE} total={members.data.total} onChange={setOffset} />
        </>
      )}
    </div>
  );
}
