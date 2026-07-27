import { useQuery } from "@tanstack/react-query";
import { ClipboardList } from "lucide-react";

import { CopyableId } from "../components/ui/CopyableId";
import { PageHeader } from "../components/ui/PageHeader";
import { PersonTag } from "../components/ui/PersonTag";
import { EmptyState, ErrorState, LoadingState } from "../components/ui/States";
import { Table, Thead, Th, Td, Tr } from "../components/ui/Table";
import { errorMessage } from "../lib/errors";
import { fetchAuditLog } from "../lib/api";

export function AuditLog() {
  const { data, isLoading, isError, error } = useQuery({ queryKey: ["discord-audit-log"], queryFn: fetchAuditLog });

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader
        icon={ClipboardList}
        title="Audit Log"
        description="The guild's Discord audit log. Per-member permissions live on that member's profile now."
      />

      {isLoading && <LoadingState label="Loading audit log" />}

      {isError && (
        <ErrorState
          title="Audit log is unavailable."
          message={`CJ may not have the View Audit Log permission. ${errorMessage(error)}`}
        />
      )}

      {data && data.length === 0 && <EmptyState message="No audit log entries yet." />}

      {data && data.length > 0 && (
        <div className="max-w-2xl">
          <Table>
            <Thead>
              <Th>Action</Th>
              <Th>Target</Th>
              <Th>By</Th>
              <Th>Reason</Th>
            </Thead>
            <tbody>
              {data.map((entry) => (
                <Tr key={entry.id}>
                  <Td className="text-ink">{entry.actionName}</Td>
                  <Td className="text-xs">{entry.targetId ? <CopyableId value={entry.targetId} /> : "none"}</Td>
                  <Td>{entry.userId ? <PersonTag person={entry.user} size={24} showId={false} /> : "unknown"}</Td>
                  <Td className="text-ink-muted">{entry.reason ?? ""}</Td>
                </Tr>
              ))}
            </tbody>
          </Table>
        </div>
      )}
    </div>
  );
}
