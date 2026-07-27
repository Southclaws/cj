import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router";
import { BarChart3, MessagesSquare, Terminal, User } from "lucide-react";
import type { LucideIcon } from "lucide-react";

import { ActionRunner, isUserScopedAction } from "../components/ActionRunner";
import { Avatar } from "../components/ui/Avatar";
import { Badge } from "../components/ui/Badge";
import { Card, CardHeader } from "../components/ui/Card";
import { CopyableId } from "../components/ui/CopyableId";
import { EmptyState, ErrorState, LoadingState } from "../components/ui/States";
import { errorMessage } from "../lib/errors";
import { fetchActions, fetchProfile, type ProfileRecord } from "../lib/api";

type Tab = "overview" | "messages" | "stats" | "commands";

const tabs: { id: Tab; label: string; icon: LucideIcon }[] = [
  { id: "overview", label: "Overview", icon: User },
  { id: "messages", label: "Messages", icon: MessagesSquare },
  { id: "stats", label: "Stats", icon: BarChart3 },
  { id: "commands", label: "Commands", icon: Terminal },
];

function tabButtonClass(active: boolean): string {
  return [
    "flex items-center gap-2 border-b-2 px-1 py-2.5 text-sm font-medium transition-colors",
    active ? "border-accent text-ink" : "border-transparent text-ink-faint hover:text-ink-muted",
  ].join(" ");
}

export function UserProfile() {
  const { id = "" } = useParams();
  const [tab, setTab] = useState<Tab>("overview");

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["profile", id],
    queryFn: () => fetchProfile(id),
  });

  const displayName = data?.member?.username ?? data?.user?.forum_user_name ?? "Unknown user";

  return (
    <div className="flex flex-col gap-6 p-8">
      <div className="flex items-center gap-4">
        <Avatar id={id} name={displayName} url={data?.member?.avatarUrl} size={56} />
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-lg font-semibold tracking-tight text-ink">{displayName}</h1>
            {data?.permissions?.isAdministrator && <Badge tone="accent">Administrator</Badge>}
            {data?.user?.burgershot_verified && <Badge tone="good">Verified</Badge>}
            {data?.member === undefined && !isLoading && <Badge tone="neutral">Not in guild</Badge>}
          </div>
          <CopyableId value={id} className="mt-0.5 text-sm" />
        </div>
      </div>

      {isLoading && <LoadingState label="Loading profile" />}
      {isError && <ErrorState title="This profile is unavailable." message={errorMessage(error)} />}

      {data && (
        <>
          <div className="flex gap-5 border-b border-border">
            {tabs.map((t) => {
              const Icon = t.icon;
              return (
                <button key={t.id} type="button" onClick={() => setTab(t.id)} className={tabButtonClass(tab === t.id)}>
                  <Icon className="h-4 w-4" strokeWidth={1.75} />
                  {t.label}
                </button>
              );
            })}
          </div>

          {tab === "overview" && <OverviewTab profile={data} />}
          {tab === "messages" && <MessagesTab profile={data} />}
          {tab === "stats" && <StatsTab profile={data} />}
          {tab === "commands" && <CommandsTab userId={id} />}
        </>
      )}
    </div>
  );
}

function OverviewTab({ profile }: { profile: ProfileRecord }) {
  return (
    <div className="grid max-w-3xl grid-cols-1 gap-4 md:grid-cols-2">
      <Card>
        <CardHeader title="Discord" />
        <div className="p-5 text-sm">
          {profile.member ? (
            <dl className="grid grid-cols-2 gap-y-2.5">
              <dt className="text-ink-faint">Username</dt>
              <dd className="text-ink">{profile.member.username}</dd>
              {profile.member.nickname && (
                <>
                  <dt className="text-ink-faint">Nickname</dt>
                  <dd className="text-ink">{profile.member.nickname}</dd>
                </>
              )}
              <dt className="text-ink-faint">Joined</dt>
              <dd className="text-ink">
                {profile.member.joinedAt ? new Date(profile.member.joinedAt).toLocaleDateString() : "unknown"}
              </dd>
              <dt className="text-ink-faint">Roles</dt>
              <dd className="text-ink">
                {profile.member.roles.length === 0 ? (
                  "none"
                ) : (
                  <ul className="flex flex-col gap-0.5">
                    {profile.member.roles.map((roleId) => (
                      <li key={roleId}>
                        <CopyableId value={roleId} className="text-xs" />
                      </li>
                    ))}
                  </ul>
                )}
              </dd>
              <dt className="text-ink-faint">CJ can act on</dt>
              <dd>
                <Badge tone={profile.member.cjCanActOn ? "good" : "neutral"}>
                  {profile.member.cjCanActOn ? "Yes" : "No"}
                </Badge>
                {!profile.member.cjCanActOn && profile.member.reasons && profile.member.reasons.length > 0 && (
                  <div className="mt-1 text-xs text-ink-faint">{profile.member.reasons.join("; ")}</div>
                )}
              </dd>
              {profile.permissions && (
                <>
                  <dt className="text-ink-faint">Administrator</dt>
                  <dd>
                    <Badge tone={profile.permissions.isAdministrator ? "accent" : "neutral"}>
                      {profile.permissions.isAdministrator ? "Yes" : "No"}
                    </Badge>
                  </dd>
                  <dt className="text-ink-faint">Permissions</dt>
                  <dd className="text-ink-muted">
                    {profile.permissions.permissions.length === 0 ? "none" : profile.permissions.permissions.join(", ")}
                  </dd>
                </>
              )}
            </dl>
          ) : (
            <p className="text-ink-faint">Not currently a member of this guild.</p>
          )}
        </div>
      </Card>

      <Card>
        <CardHeader title="CJ data" />
        <div className="p-5 text-sm">
          {profile.user ? (
            <dl className="grid grid-cols-2 gap-y-2.5">
              <dt className="text-ink-faint">Forum name</dt>
              <dd className="text-ink">{profile.user.forum_user_name || "none"}</dd>
              <dt className="text-ink-faint">Forum ID</dt>
              <dd className="text-xs">
                {profile.user.forum_user_id ? <CopyableId value={profile.user.forum_user_id} /> : "none"}
              </dd>
              <dt className="text-ink-faint">Burgershot name</dt>
              <dd className="text-ink">{profile.user.burger_user_name || "none"}</dd>
              <dt className="text-ink-faint">Burgershot ID</dt>
              <dd className="text-xs">
                {profile.user.burger_user_id ? <CopyableId value={profile.user.burger_user_id} /> : "none"}
              </dd>
              <dt className="text-ink-faint">Verified</dt>
              <dd>
                <Badge tone={profile.user.burgershot_verified ? "good" : "neutral"}>
                  {profile.user.burgershot_verified ? "Yes" : "No"}
                </Badge>
              </dd>
            </dl>
          ) : (
            <p className="text-ink-faint">No CJ record for this user.</p>
          )}
        </div>
      </Card>
    </div>
  );
}

function MessagesTab({ profile }: { profile: ProfileRecord }) {
  return (
    <Card className="max-w-3xl">
      <CardHeader
        title="Recent messages"
        description={`Last ${profile.recentMessages.length} across all channels.`}
        action={
          <Link
            to={`/data/messages?userId=${profile.discordUserId}`}
            className="text-sm text-accent hover:text-accent-hover"
          >
            Search all their messages
          </Link>
        }
      />

      <div className="p-5">
        {profile.recentMessages.length === 0 ? (
          <EmptyState message="No messages recorded for this user." />
        ) : (
          <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full text-left text-sm">
              <thead>
                <tr className="border-b border-border text-ink-faint">
                  <th className="px-4 py-2 font-medium">Time</th>
                  <th className="px-4 py-2 font-medium">Channel</th>
                  <th className="px-4 py-2 font-medium">Message</th>
                </tr>
              </thead>
              <tbody>
                {profile.recentMessages.map((message) => (
                  <tr key={message.messageId} className="border-b border-border last:border-b-0">
                    <td className="px-4 py-2 whitespace-nowrap text-ink-faint">
                      {new Date(message.timestamp * 1000).toLocaleString()}
                    </td>
                    <td className="px-4 py-2 text-ink-muted">
                      #{message.channelName || message.channelId}
                      <div>
                        <CopyableId value={message.channelId} className="text-xs" />
                      </div>
                    </td>
                    <td className="px-4 py-2 text-ink">{message.message}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </Card>
  );
}

function StatsTab({ profile }: { profile: ProfileRecord }) {
  return (
    <div className="grid max-w-3xl grid-cols-1 gap-4 md:grid-cols-2">
      <Card>
        <CardHeader title="Leaderboard" />
        <dl className="grid grid-cols-2 gap-y-2.5 p-5 text-sm">
          <dt className="text-ink-faint">Messages sent</dt>
          <dd className="text-ink">{profile.messageCount}</dd>
          <dt className="text-ink-faint">Rank</dt>
          <dd className="text-ink">{profile.rank > 0 ? `#${profile.rank}` : "unranked"}</dd>
        </dl>
      </Card>

      <Card>
        <CardHeader title="Reactions received" />
        <div className="p-5">
          {profile.user?.received_reactions && profile.user.received_reactions.length > 0 ? (
            <ul className="flex flex-col gap-1.5 text-sm">
              {profile.user.received_reactions.map((reaction) => (
                <li key={reaction.Reaction} className="flex items-center justify-between">
                  <span className="text-ink">{reaction.Reaction}</span>
                  <span className="rounded-full bg-surface-3 px-2 py-0.5 text-xs text-ink-muted">{reaction.Counter}</span>
                </li>
              ))}
            </ul>
          ) : (
            <EmptyState message="No reactions on record." />
          )}
        </div>
      </Card>
    </div>
  );
}

function CommandsTab({ userId }: { userId: string }) {
  const { data } = useQuery({ queryKey: ["actions"], queryFn: fetchActions });
  const userActions = (data ?? []).filter((action) => isUserScopedAction(action.name));

  return (
    <div className="flex max-w-lg flex-col gap-4">
      {userActions.length === 0 ? (
        <EmptyState message="No commands take a specific user as input yet." />
      ) : (
        userActions.map((action) => (
          <Card key={action.name} className="p-5">
            <ActionRunner action={action} presetUserId={userId} />
          </Card>
        ))
      )}
    </div>
  );
}
