import { Fragment, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Link } from "react-router";

import { MemberUserPicker } from "./MemberUserPicker";
import { MultiPicker } from "./MultiPicker";
import { Avatar } from "./ui/Avatar";
import { Badge, riskTone } from "./ui/Badge";
import { Button } from "./ui/Button";
import { errorMessage } from "../lib/errors";
import { useRoleOptions } from "../lib/entitySearch";
import {
  type ActionInfoRecord,
  type ActionOutcomeRecord,
  executeAction,
  fetchCommands,
  fetchJobs,
  fetchServices,
  previewAction,
} from "../lib/api";

const fieldClass =
  "rounded-lg border border-border bg-base px-3 py-1.5 text-sm text-ink placeholder:text-ink-faint focus:border-accent/50 focus:outline-none focus:ring-2 focus:ring-accent/20";

function humanizeKey(key: string): string {
  const spaced = key.replace(/([A-Z])/g, " $1").toLowerCase();
  return spaced.charAt(0).toUpperCase() + spaced.slice(1);
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isLeaderboardEntry(value: unknown): value is Record<string, unknown> & { userId: string } {
  return isRecord(value) && typeof value.userId === "string";
}

function isWikiHit(value: unknown): value is Record<string, unknown> & { pageName: string; url: string } {
  return isRecord(value) && typeof value.pageName === "string" && typeof value.url === "string";
}

function LeaderboardList({ entries }: { entries: (Record<string, unknown> & { userId: string })[] }) {
  return (
    <ol className="flex flex-col gap-2">
      {entries.map((entry, index) => {
        const name = (entry.username as string) || entry.userId;
        return (
          <li key={entry.userId} className="flex items-center justify-between gap-3">
            <div className="flex min-w-0 items-center gap-2">
              <span className="w-4 shrink-0 text-right text-ink-faint">{index + 1}</span>
              <Avatar id={entry.userId} name={name} url={entry.avatarUrl as string | undefined} size={22} />
              <Link to={`/users/${entry.userId}`} className="truncate text-ink hover:text-accent">
                {name}
              </Link>
            </div>
            <span className="shrink-0 text-xs text-ink-muted">
              {typeof entry.messages === "number" ? `${entry.messages} messages` : `${entry.counter} ${entry.reaction}`}
            </span>
          </li>
        );
      })}
    </ol>
  );
}

function WikiHitList({ hits }: { hits: (Record<string, unknown> & { pageName: string; url: string })[] }) {
  return (
    <ul className="flex flex-col gap-1.5">
      {hits.map((hit) => (
        <li key={hit.url}>
          <a href={hit.url} className="text-accent hover:text-accent-hover">
            {hit.pageName}
          </a>
        </li>
      ))}
    </ul>
  );
}

function DetailValue({ value }: { value: unknown }) {
  if (value === null || value === undefined || value === "") {
    return <span className="text-ink-faint">none</span>;
  }
  if (typeof value === "boolean") return <span>{value ? "Yes" : "No"}</span>;

  if (Array.isArray(value)) {
    if (value.length === 0) return <span className="text-ink-faint">none</span>;
    if (value.every(isLeaderboardEntry)) return <LeaderboardList entries={value} />;
    if (value.every(isWikiHit)) return <WikiHitList hits={value} />;
    if (value.every((item) => !isRecord(item) && !Array.isArray(item))) {
      return <span>{value.join(", ")}</span>;
    }
    return (
      <ul className="flex flex-col gap-1">
        {value.map((item, index) => (
          <li key={index}>
            <DetailValue value={item} />
          </li>
        ))}
      </ul>
    );
  }

  if (isRecord(value)) return <DetailFields detail={value} />;

  return <span>{String(value)}</span>;
}

function DetailFields({ detail }: { detail: Record<string, unknown> }) {
  const fields = Object.entries(detail).filter(([key]) => key !== "kind");
  return (
    <dl className="grid grid-cols-[auto_1fr] items-baseline gap-x-3 gap-y-1.5">
      {fields.map(([key, value]) => (
        <Fragment key={key}>
          <dt className="text-xs text-ink-faint">{humanizeKey(key)}</dt>
          <dd className="text-ink">
            <DetailValue value={value} />
          </dd>
        </Fragment>
      ))}
    </dl>
  );
}

function Outcome({ outcome }: { outcome: ActionOutcomeRecord }) {
  const detail = outcome.detail;
  const hasDetail = detail && Object.keys(detail).length > 0;
  const isLeaderboard = hasDetail && detail.kind === "leaderboard" && Array.isArray(detail.entries);

  return (
    <div className="rounded-lg border border-good/25 bg-good-soft p-3 text-sm text-ink">
      <p>{outcome.summary}</p>
      {hasDetail && (
        <div className="mt-2.5 border-t border-good/20 pt-2.5 text-sm">
          {isLeaderboard ? (
            <LeaderboardList entries={detail.entries as (Record<string, unknown> & { userId: string })[]} />
          ) : (
            <DetailFields detail={detail} />
          )}
        </div>
      )}
    </div>
  );
}

interface ActionRunnerProps {
  action: ActionInfoRecord;
  presetUserId?: string;
}

export function ActionRunner({ action, presetUserId }: ActionRunnerProps) {
  const needsService = action.name === "services.test";
  const needsJob = action.name === "jobs.run";
  const needsCommand = action.name === "commands.get-settings" || action.name === "commands.set-settings";
  const needsSettingsJson = action.name === "commands.set-settings";
  const needsWikiTerm = action.name === "wiki.search";
  const needsUserID = action.name === "leaderboards.user-rank";
  const needsLimit = action.name === "leaderboards.top-messages" || action.name === "leaderboards.top-reactions";
  const needsReaction = action.name === "leaderboards.top-reactions";

  const { data: services } = useQuery({ queryKey: ["services"], queryFn: fetchServices, enabled: needsService });
  const { data: jobs } = useQuery({ queryKey: ["jobs"], queryFn: fetchJobs, enabled: needsJob });
  const { data: commands } = useQuery({ queryKey: ["commands"], queryFn: fetchCommands, enabled: needsCommand });
  const { all: allRoles } = useRoleOptions("");

  const [service, setService] = useState("");
  const [job, setJob] = useState("");
  const [command, setCommand] = useState("");
  const [roleIds, setRoleIds] = useState<string[]>([]);
  const [cooldownSeconds, setCooldownSeconds] = useState("");
  const [miscJson, setMiscJson] = useState("{}");
  const [wikiTerm, setWikiTerm] = useState("");
  const [userID, setUserID] = useState(presetUserId ?? "");
  const [limit, setLimit] = useState("");
  const [reaction, setReaction] = useState("");
  const [outcome, setOutcome] = useState<ActionOutcomeRecord | null>(null);

  let input: Record<string, unknown> = {};
  let ready = true;
  let miscJsonError: string | null = null;

  if (needsService) {
    input = { service };
    ready = service !== "";
  } else if (needsJob) {
    input = { name: job };
    ready = job !== "";
  } else if (needsWikiTerm) {
    input = { term: wikiTerm };
    ready = wikiTerm.trim().length >= 3;
  } else if (needsUserID) {
    input = { userId: userID };
    ready = userID.trim() !== "";
  } else if (needsCommand) {
    ready = command !== "";
    if (needsSettingsJson) {
      let misc: unknown = {};
      try {
        misc = miscJson.trim() === "" ? {} : JSON.parse(miscJson);
      } catch {
        ready = false;
        miscJsonError = "Misc must be valid JSON.";
      }
      const cooldownNumber = cooldownSeconds.trim() === "" ? 0 : Number(cooldownSeconds);
      input = {
        command,
        settings: { roles: roleIds, cooldown: cooldownNumber * 1_000_000_000, misc },
      };
    } else {
      input = { command };
    }
  } else if (needsLimit) {
    input = { limit: limit === "" ? undefined : Number(limit), reaction: needsReaction ? reaction || undefined : undefined };
  }

  const previewMutation = useMutation({
    mutationFn: () => previewAction(action.name, input),
    onSuccess: setOutcome,
  });

  const executeMutation = useMutation({
    mutationFn: () => executeAction(action.name, input),
    onSuccess: setOutcome,
  });

  const busy = previewMutation.isPending || executeMutation.isPending;
  const activeError = previewMutation.error ?? executeMutation.error;

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-start justify-between gap-3">
        <div>
          <h2 className="font-mono text-sm font-medium text-ink">{action.name}</h2>
          <p className="mt-1 text-sm text-ink-faint">{action.description}</p>
        </div>
        {(action.risk === "state-changing" || action.risk === "dangerous") && (
          <Badge tone={riskTone(action.risk)}>{action.risk}</Badge>
        )}
      </div>

      {needsService && (
        <select value={service} onChange={(e) => setService(e.target.value)} className={fieldClass}>
          <option value="">Select a service...</option>
          {services?.map((s) => (
            <option key={s.name} value={s.name}>
              {s.name}
            </option>
          ))}
        </select>
      )}

      {needsJob && (
        <select value={job} onChange={(e) => setJob(e.target.value)} className={fieldClass}>
          <option value="">Select a job...</option>
          {jobs?.map((j) => (
            <option key={j.name} value={j.name}>
              {j.name}
            </option>
          ))}
        </select>
      )}

      {needsCommand && (
        <select value={command} onChange={(e) => setCommand(e.target.value)} className={fieldClass}>
          <option value="">Select a command...</option>
          {commands?.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
      )}

      {needsSettingsJson && (
        <div className="flex flex-col gap-3">
          <div className="flex flex-col gap-1">
            <span className="text-xs text-ink-faint">Roles allowed to use this command</span>
            <MultiPicker
              ids={roleIds}
              onChange={setRoleIds}
              useOptions={useRoleOptions}
              placeholder="Search roles..."
              renderChipLabel={(id) => allRoles.find((r) => r.id === id)?.label ?? id}
            />
          </div>

          <label className="flex flex-col gap-1">
            <span className="text-xs text-ink-faint">Cooldown, in seconds (0 for none)</span>
            <input
              type="number"
              min="0"
              value={cooldownSeconds}
              onChange={(e) => setCooldownSeconds(e.target.value)}
              placeholder="0"
              className={`${fieldClass} w-40`}
            />
          </label>

          <label className="flex flex-col gap-1">
            <span className="text-xs text-ink-faint">Misc (advanced, raw JSON)</span>
            <textarea
              value={miscJson}
              onChange={(e) => setMiscJson(e.target.value)}
              rows={3}
              className={`${fieldClass} font-mono`}
            />
            {miscJsonError && <p className="text-xs text-bad">{miscJsonError}</p>}
          </label>

          <p className="text-xs text-ink-faint">
            This replaces the command's full settings. Leaving a field at its default clears it, it doesn't leave the
            existing value untouched: check "commands.get-settings" first if you only want to change one part.
          </p>
        </div>
      )}

      {needsWikiTerm && (
        <input
          type="text"
          placeholder="Search term (3+ characters)..."
          value={wikiTerm}
          onChange={(e) => setWikiTerm(e.target.value)}
          className={fieldClass}
        />
      )}

      {needsUserID &&
        (presetUserId ? (
          <p className="text-sm text-ink-muted">
            For user <span className="font-mono text-ink">{presetUserId}</span>
          </p>
        ) : (
          <MemberUserPicker value={userID} onChange={setUserID} placeholder="Search by name or ID..." />
        ))}

      {needsLimit && (
        <div className="flex gap-3">
          <input
            type="number"
            placeholder="Limit (default 10)"
            value={limit}
            onChange={(e) => setLimit(e.target.value)}
            className={`${fieldClass} w-40`}
          />
          {needsReaction && (
            <input
              type="text"
              placeholder="Reaction (optional)"
              value={reaction}
              onChange={(e) => setReaction(e.target.value)}
              className={fieldClass}
            />
          )}
        </div>
      )}

      <div className="flex items-center gap-3">
        <Button variant="secondary" disabled={busy || !ready} onClick={() => previewMutation.mutate()}>
          Preview
        </Button>
        <Button variant="primary" disabled={busy || !ready} onClick={() => executeMutation.mutate()}>
          Run
        </Button>
      </div>

      {activeError && <p className="text-sm text-bad">{errorMessage(activeError)}</p>}
      {outcome && <Outcome outcome={outcome} />}
    </div>
  );
}

export function isUserScopedAction(name: string): boolean {
  return name === "leaderboards.user-rank";
}
