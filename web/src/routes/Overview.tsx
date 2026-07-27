import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Activity, Clock, Heart, LayoutDashboard, Server, Trophy, Zap } from "lucide-react";

import { ActionRunner, isUserScopedAction } from "../components/ActionRunner";
import { Modal } from "../components/Modal";
import { Card, CardHeader } from "../components/ui/Card";
import { PageHeader } from "../components/ui/PageHeader";
import { PersonTag } from "../components/ui/PersonTag";
import { EmptyState } from "../components/ui/States";
import { StatusDot, statusTone } from "../components/ui/StatusDot";
import {
  fetchActions,
  fetchJobs,
  fetchServices,
  fetchStatus,
  fetchTopMessages,
  fetchTopReactions,
  type ActionInfoRecord,
} from "../lib/api";

function StatusPanel() {
  const { data } = useQuery({ queryKey: ["status"], queryFn: fetchStatus, refetchInterval: 5000 });
  if (!data) return null;
  return (
    <Card>
      <CardHeader title="Status" icon={<Activity className="h-4 w-4" />} />
      <dl className="grid grid-cols-2 gap-x-4 gap-y-2.5 p-5 text-sm">
        <dt className="text-ink-faint">Version</dt>
        <dd className="text-ink">{data.version}</dd>
        <dt className="text-ink-faint">Started</dt>
        <dd className="text-ink">{new Date(data.startedAt).toLocaleString()}</dd>
        <dt className="text-ink-faint">Uptime</dt>
        <dd className="text-ink">{data.uptime}</dd>
      </dl>
    </Card>
  );
}

function ServicesPanel() {
  const { data } = useQuery({ queryKey: ["services"], queryFn: fetchServices, refetchInterval: 10000 });
  if (!data) return null;
  return (
    <Card>
      <CardHeader title="Services" icon={<Server className="h-4 w-4" />} />
      <ul className="flex flex-col gap-2.5 p-5">
        {data.map((service) => (
          <li key={service.name} className="flex items-center justify-between text-sm">
            <div className="flex items-center gap-2">
              <StatusDot tone={statusTone(service.status)} />
              <span className="capitalize text-ink">{service.name}</span>
            </div>
            <span className="text-ink-faint">{service.status.replace("_", " ")}</span>
          </li>
        ))}
      </ul>
    </Card>
  );
}

function JobsPanel() {
  const { data } = useQuery({ queryKey: ["jobs"], queryFn: fetchJobs, refetchInterval: 10000 });
  if (!data) return null;
  return (
    <Card>
      <CardHeader title="Jobs" icon={<Clock className="h-4 w-4" />} />
      <div className="p-5">
        {data.length === 0 ? (
          <EmptyState message="No scheduled jobs have run yet." />
        ) : (
          <ul className="flex flex-col gap-3">
            {data.map((job) => (
              <li key={job.name} className="text-sm">
                <div className="flex items-center justify-between">
                  <span className="text-ink">{job.name}</span>
                  <span className="font-mono text-xs text-ink-faint">{job.schedule}</span>
                </div>
                <div className="text-xs text-ink-faint">
                  {job.lastRun ? `last ran ${new Date(job.lastRun).toLocaleString()}` : "never run"} · {job.runCount} runs
                  {job.lastError && <span className="text-bad"> · {job.lastError}</span>}
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </Card>
  );
}

function TopMessagesPanel() {
  const { data } = useQuery({ queryKey: ["top-messages", 10], queryFn: () => fetchTopMessages(10) });
  if (!data) return null;
  return (
    <Card>
      <CardHeader title="Top messages" icon={<Trophy className="h-4 w-4" />} />
      <div className="p-5">
        {data.length === 0 ? (
          <EmptyState message="No messages recorded yet." />
        ) : (
          <ol className="flex flex-col gap-2.5 text-sm">
            {data.map((entry, index) => (
              <li key={entry.person.id} className="flex items-center justify-between gap-3">
                <div className="flex min-w-0 items-center gap-2">
                  <span className="w-4 shrink-0 text-right text-ink-faint">{index + 1}</span>
                  <PersonTag person={entry.person} size={24} showId={false} />
                </div>
                <span className="shrink-0 rounded-full bg-surface-3 px-2 py-0.5 text-xs text-ink-muted">
                  {entry.messages}
                </span>
              </li>
            ))}
          </ol>
        )}
      </div>
    </Card>
  );
}

function TopReactionsPanel() {
  const { data } = useQuery({ queryKey: ["top-reactions", 10], queryFn: () => fetchTopReactions(10) });
  if (!data) return null;
  return (
    <Card>
      <CardHeader title="Top reactions" icon={<Heart className="h-4 w-4" />} />
      <div className="p-5">
        {data.length === 0 ? (
          <EmptyState message="No reactions recorded yet." />
        ) : (
          <ol className="flex flex-col gap-2.5 text-sm">
            {data.map((entry, index) => (
              <li key={`${entry.person.id}-${entry.reaction}`} className="flex items-center justify-between gap-3">
                <div className="flex min-w-0 items-center gap-2">
                  <span className="w-4 shrink-0 text-right text-ink-faint">{index + 1}</span>
                  <PersonTag person={entry.person} size={24} showId={false} />
                </div>
                <span className="shrink-0 rounded-full bg-surface-3 px-2 py-0.5 text-xs text-ink-muted">
                  {entry.counter} {entry.reaction}
                </span>
              </li>
            ))}
          </ol>
        )}
      </div>
    </Card>
  );
}

function QuickActionsPanel() {
  const { data } = useQuery({ queryKey: ["actions"], queryFn: fetchActions });
  const [active, setActive] = useState<ActionInfoRecord | null>(null);

  const globalActions = (data ?? []).filter((action) => !isUserScopedAction(action.name));

  return (
    <Card className="md:col-span-2">
      <CardHeader title="Server actions" icon={<Zap className="h-4 w-4" />} />
      <div className="p-5">
        {globalActions.length === 0 ? (
          <EmptyState message="No actions available." />
        ) : (
          <div className="grid grid-cols-2 gap-2 lg:grid-cols-3">
            {globalActions.map((action) => (
              <button
                key={action.name}
                type="button"
                onClick={() => setActive(action)}
                className="rounded-lg border border-border bg-surface px-3 py-2.5 text-left font-mono text-xs text-ink-muted transition-colors hover:border-accent/40 hover:bg-accent-soft hover:text-ink"
              >
                {action.name}
              </button>
            ))}
          </div>
        )}
      </div>

      {active && (
        <Modal title="Run action" onClose={() => setActive(null)}>
          <ActionRunner action={active} />
        </Modal>
      )}
    </Card>
  );
}

export function Overview() {
  return (
    <div className="flex flex-col gap-6 p-8">
      <PageHeader icon={LayoutDashboard} title="Dashboard" description="Everything CJ is doing, at a glance." />

      <div className="grid max-w-5xl grid-cols-1 gap-4 md:grid-cols-2">
        <StatusPanel />
        <ServicesPanel />
        <JobsPanel />
        <TopMessagesPanel />
        <TopReactionsPanel />
        <QuickActionsPanel />
      </div>
    </div>
  );
}
