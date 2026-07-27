import { useQuery } from "@tanstack/react-query";
import { Wrench } from "lucide-react";

import { Card } from "../components/ui/Card";
import { PageHeader } from "../components/ui/PageHeader";
import { ErrorState, LoadingState } from "../components/ui/States";
import { StatusDot } from "../components/ui/StatusDot";
import { errorMessage } from "../lib/errors";
import { fetchConfigStatus } from "../lib/api";

function fieldLabel(name: string): string {
  return name
    .split("_")
    .map((part) => part[0].toUpperCase() + part.slice(1))
    .join(" ");
}

export function Configuration() {
  const { data, error, isLoading, isError } = useQuery({
    queryKey: ["config-status"],
    queryFn: fetchConfigStatus,
  });

  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader
        icon={Wrench}
        title="Configuration"
        description="Shows whether required configuration is present. Values themselves, including secrets, are never sent to the browser."
      />

      {isLoading && <LoadingState label="Loading configuration status" />}

      {isError && <ErrorState title="Configuration status is unavailable." message={errorMessage(error)} />}

      {data && (
        <Card className="max-w-md p-6 text-sm">
          <div className="mb-4 flex items-center justify-between border-b border-border pb-3">
            <span className="text-ink-muted">Database</span>
            <span className="flex items-center gap-2">
              <StatusDot tone={data.databaseEnabled ? "good" : "bad"} />
              <span className={data.databaseEnabled ? "text-ink" : "text-bad"}>
                {data.databaseEnabled ? "Enabled" : "Disabled"}
              </span>
            </span>
          </div>
          <ul className="flex flex-col gap-2.5">
            {data.fields.map((field) => (
              <li key={field.name} className="flex items-center justify-between">
                <span className="text-ink-muted">{fieldLabel(field.name)}</span>
                <span className="flex items-center gap-2">
                  <StatusDot tone={field.set ? "good" : "bad"} />
                  <span className={field.set ? "text-ink" : "text-bad"}>{field.set ? "Set" : "Missing"}</span>
                </span>
              </li>
            ))}
          </ul>
        </Card>
      )}
    </div>
  );
}
