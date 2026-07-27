import { useQuery } from "@tanstack/react-query";
import { ShieldCheck } from "lucide-react";

import { Badge } from "../components/ui/Badge";
import { CopyableId } from "../components/ui/CopyableId";
import { PageHeader } from "../components/ui/PageHeader";
import { ErrorState, LoadingState } from "../components/ui/States";
import { Table, Thead, Th, Td, Tr } from "../components/ui/Table";
import { errorMessage } from "../lib/errors";
import { fetchRoles } from "../lib/api";

function roleColor(color: number): string {
  if (color === 0) {
    return "#99a1af";
  }
  return `#${color.toString(16).padStart(6, "0")}`;
}

export function Roles() {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["discord-roles"],
    queryFn: fetchRoles,
  });

  const sorted = data ? [...data].sort((a, b) => b.position - a.position) : [];

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader icon={ShieldCheck} title="Roles" description="Every role in the guild, ranked by hierarchy." />

      {isLoading && <LoadingState label="Loading roles" />}

      {isError && <ErrorState title="Roles are unavailable." message={errorMessage(error)} />}

      {data && (
        <div className="max-w-4xl">
          <Table>
            <Thead>
              <Th>Role</Th>
              <Th>Members</Th>
              <Th>Permissions</Th>
              <Th>CJ can manage</Th>
            </Thead>
            <tbody>
              {sorted.map((role) => (
                <Tr key={role.id}>
                  <Td>
                    <span className="flex items-center gap-2 text-ink">
                      <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: roleColor(role.color) }} />
                      {role.name}
                      {role.managed && <span className="text-xs text-ink-faint">(managed)</span>}
                    </span>
                    <CopyableId value={role.id} className="text-xs" />
                  </Td>
                  <Td className="text-ink-muted">{role.memberCount}</Td>
                  <Td className="text-ink-muted">
                    {role.permissions.length === 0 ? <span className="text-ink-faint">none</span> : role.permissions.join(", ")}
                  </Td>
                  <Td>
                    <Badge tone={role.cjCanManage ? "good" : "neutral"}>{role.cjCanManage ? "Yes" : "No"}</Badge>
                    {!role.cjCanManage && role.reasons && role.reasons.length > 0 && (
                      <div className="mt-1 text-xs text-ink-faint">{role.reasons.join("; ")}</div>
                    )}
                  </Td>
                </Tr>
              ))}
            </tbody>
          </Table>
        </div>
      )}
    </div>
  );
}
