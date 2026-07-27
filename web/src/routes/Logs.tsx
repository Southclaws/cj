import { useEffect, useRef, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Pause, Play, ScrollText } from "lucide-react";

import { Button } from "../components/ui/Button";
import { PageHeader } from "../components/ui/PageHeader";
import { StatusDot } from "../components/ui/StatusDot";
import { fetchLogs, logRecordSchema, type LogFilters, type LogRecord } from "../lib/api";

const MAX_DISPLAYED = 500;

function matchesFilters(record: LogRecord, filters: LogFilters): boolean {
  if (filters.level && record.level.toLowerCase() !== filters.level.toLowerCase()) {
    return false;
  }
  if (filters.component && (record.component ?? "").toLowerCase() !== filters.component.toLowerCase()) {
    return false;
  }
  if (filters.q && !record.message.toLowerCase().includes(filters.q.toLowerCase())) {
    return false;
  }
  if (filters.correlationId) {
    const id = record.fields?.correlation_id;
    if (id !== filters.correlationId) {
      return false;
    }
  }
  return true;
}

function levelColor(level: string): string {
  switch (level.toLowerCase()) {
    case "error":
    case "fatal":
      return "text-bad";
    case "warn":
      return "text-warn";
    case "debug":
      return "text-ink-faint";
    default:
      return "text-ink-muted";
  }
}

const selectClass = "rounded-lg border border-border bg-surface-2 px-2.5 py-1.5 text-sm text-ink-muted";
const inputClass =
  "rounded-lg border border-border bg-surface-2 px-2.5 py-1.5 text-sm text-ink-muted placeholder:text-ink-faint focus:border-accent/50 focus:outline-none";

export function Logs() {
  const [level, setLevel] = useState("");
  const [component, setComponent] = useState("");
  const [q, setQ] = useState("");
  const [correlationId, setCorrelationId] = useState("");
  const [paused, setPaused] = useState(false);
  const [records, setRecords] = useState<LogRecord[]>([]);

  const filters: LogFilters = { level, component, q, correlationId };
  const filtersRef = useRef(filters);
  filtersRef.current = filters;
  const pausedRef = useRef(paused);
  pausedRef.current = paused;

  const { data: backlog } = useQuery({
    queryKey: ["logs", filters],
    queryFn: () => fetchLogs(filters),
  });

  useEffect(() => {
    if (backlog) {
      setRecords([...backlog].reverse());
    }
  }, [backlog]);

  useEffect(() => {
    const source = new EventSource("/api/v1/logs/stream");
    source.addEventListener("log", (event) => {
      if (pausedRef.current) {
        return;
      }
      let parsed: LogRecord;
      try {
        parsed = logRecordSchema.parse(JSON.parse((event as MessageEvent<string>).data));
      } catch {
        return;
      }
      if (!matchesFilters(parsed, filtersRef.current)) {
        return;
      }
      setRecords((prev) => [...prev, parsed].slice(-MAX_DISPLAYED));
    });
    return () => source.close();
  }, []);

  return (
    <div className="flex h-full flex-col gap-4 p-8">
      <PageHeader
        icon={ScrollText}
        title="Logs"
        description="Live structured logs, streamed from the process."
        action={
          <Button variant="secondary" onClick={() => setPaused((prev) => !prev)}>
            <span className="flex items-center gap-1.5">
              {paused ? <Play className="h-3.5 w-3.5" /> : <Pause className="h-3.5 w-3.5" />}
              {paused ? "Resume" : "Pause"}
            </span>
          </Button>
        }
      />

      <div className="flex items-center gap-2 text-xs text-ink-faint">
        <StatusDot tone={paused ? "neutral" : "good"} pulse={!paused} />
        {paused ? "Paused" : "Live"}
      </div>

      <div className="flex flex-wrap gap-3">
        <select value={level} onChange={(e) => setLevel(e.target.value)} className={selectClass}>
          <option value="">All levels</option>
          <option value="debug">Debug</option>
          <option value="info">Info</option>
          <option value="warn">Warn</option>
          <option value="error">Error</option>
        </select>
        <input
          type="text"
          placeholder="Component"
          value={component}
          onChange={(e) => setComponent(e.target.value)}
          className={inputClass}
        />
        <input
          type="text"
          placeholder="Search text"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          className={inputClass}
        />
        <input
          type="text"
          placeholder="Correlation ID"
          value={correlationId}
          onChange={(e) => setCorrelationId(e.target.value)}
          className={inputClass}
        />
      </div>

      <div className="flex-1 overflow-y-auto rounded-xl border border-border bg-base p-3 font-mono text-xs">
        {records.length === 0 && <p className="text-ink-faint">No log records yet.</p>}
        {records.map((record, index) => (
          <div key={index} className="flex gap-2 border-b border-border/60 py-1 last:border-b-0">
            <span className="shrink-0 text-ink-faint">{new Date(record.timestamp).toLocaleTimeString()}</span>
            <span className={`shrink-0 uppercase ${levelColor(record.level)}`}>{record.level}</span>
            {record.component && <span className="shrink-0 text-ink-faint">[{record.component}]</span>}
            <span className="text-ink-muted">{record.message}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
