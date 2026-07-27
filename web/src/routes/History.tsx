import { useQuery } from "@tanstack/react-query";
import { History as HistoryIcon } from "lucide-react";

import { Badge } from "../components/ui/Badge";
import { PageHeader } from "../components/ui/PageHeader";
import { EmptyState, ErrorState, LoadingState } from "../components/ui/States";
import { Table, Thead, Th, Td, Tr } from "../components/ui/Table";
import { errorMessage } from "../lib/errors";
import { fetchActionRuns } from "../lib/api";

export function History() {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["action-runs"],
    queryFn: fetchActionRuns,
    refetchInterval: 10000,
  });

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader
        icon={HistoryIcon}
        title="History"
        description="Every action preview and run, newest first. This is a read-only, append-only record: there is no edit or delete."
      />

      {isLoading && <LoadingState label="Loading history" />}

      {isError && <ErrorState title="Action history is unavailable." message={errorMessage(error)} />}

      {data && data.length === 0 && <EmptyState message="No actions have been previewed or run yet." />}

      {data && data.length > 0 && (
        <div className="max-w-4xl">
          <Table>
            <Thead>
              <Th>Time</Th>
              <Th>Action</Th>
              <Th>Mode</Th>
              <Th>Risk</Th>
              <Th>Outcome</Th>
              <Th>Duration</Th>
              <Th>Summary</Th>
            </Thead>
            <tbody>
              {data.map((run) => (
                <Tr key={run.id}>
                  <Td className="whitespace-nowrap align-top text-ink-faint">{new Date(run.timestamp).toLocaleString()}</Td>
                  <Td className="align-top font-mono text-xs text-ink">{run.action}</Td>
                  <Td className="align-top text-ink-muted">{run.mode}</Td>
                  <Td className="align-top text-ink-muted">{run.risk}</Td>
                  <Td className="align-top">
                    <Badge tone={run.outcome === "success" ? "good" : "bad"}>{run.outcome}</Badge>
                  </Td>
                  <Td className="align-top text-ink-muted">{run.durationMs}ms</Td>
                  <Td className="align-top text-ink-muted">{run.errorMessage ?? run.summary ?? ""}</Td>
                </Tr>
              ))}
            </tbody>
          </Table>
        </div>
      )}
    </div>
  );
}
