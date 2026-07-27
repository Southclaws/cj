import { z } from "zod";

export class ApiError extends Error {
  code: string;
  requestId: string;
  status: number;

  constructor(status: number, code: string, message: string, requestId: string) {
    super(message);
    this.status = status;
    this.code = code;
    this.requestId = requestId;
  }
}

const errorEnvelopeSchema = z.object({
  error: z.object({
    code: z.string(),
    message: z.string(),
    requestId: z.string(),
  }),
});

export async function apiFetch<T>(path: string, schema: z.ZodType<T>, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: { Accept: "application/json", ...init?.headers },
  });

  const body: unknown = await response.json().catch(() => undefined);

  if (!response.ok) {
    const parsed = errorEnvelopeSchema.safeParse(body);
    if (parsed.success) {
      throw new ApiError(response.status, parsed.data.error.code, parsed.data.error.message, parsed.data.error.requestId);
    }
    throw new ApiError(response.status, "unknown_error", "The dashboard returned an unexpected error.", "");
  }

  return schema.parse(body);
}

export const statusSchema = z.object({
  version: z.string(),
  startedAt: z.string(),
  uptime: z.string(),
});

export type Status = z.infer<typeof statusSchema>;

export function fetchStatus(): Promise<Status> {
  return apiFetch("/api/v1/status", statusSchema);
}

export const configStatusSchema = z.object({
  databaseEnabled: z.boolean(),
  fields: z.array(
    z.object({
      name: z.string(),
      set: z.boolean(),
    }),
  ),
});

export type ConfigStatus = z.infer<typeof configStatusSchema>;

export function fetchConfigStatus(): Promise<ConfigStatus> {
  return apiFetch("/api/v1/config/status", configStatusSchema);
}

export const guildSettingsSchema = z.object({
  searchMessageChannelId: z.string(),
  errorReportChannelId: z.string(),
  leaderboardChannelId: z.string(),
  adsChannelId: z.string(),
  ltfChannelId: z.string(),
  ltfUserIds: z.array(z.string()),
});

export type GuildSettings = z.infer<typeof guildSettingsSchema>;

export function fetchGuildSettings(): Promise<GuildSettings> {
  return apiFetch("/api/v1/data/settings", guildSettingsSchema);
}

export function saveGuildSettings(settings: GuildSettings): Promise<GuildSettings> {
  return apiFetch("/api/v1/data/settings", guildSettingsSchema, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(settings),
  });
}

export const logRecordSchema = z.object({
  timestamp: z.string(),
  level: z.string(),
  component: z.string().optional(),
  message: z.string(),
  fields: z.record(z.string(), z.unknown()).optional(),
});

export type LogRecord = z.infer<typeof logRecordSchema>;

export interface LogFilters {
  level?: string;
  component?: string;
  q?: string;
  correlationId?: string;
}

export function fetchLogs(filters: LogFilters): Promise<LogRecord[]> {
  const params = new URLSearchParams();
  if (filters.level) params.set("level", filters.level);
  if (filters.component) params.set("component", filters.component);
  if (filters.q) params.set("q", filters.q);
  if (filters.correlationId) params.set("correlationId", filters.correlationId);
  const query = params.toString();
  return apiFetch(`/api/v1/logs${query ? `?${query}` : ""}`, z.array(logRecordSchema));
}

export const serviceStatusSchema = z.object({
  name: z.string(),
  status: z.string(),
  detail: z.string().optional(),
});

export type ServiceStatusRecord = z.infer<typeof serviceStatusSchema>;

export function fetchServices(): Promise<ServiceStatusRecord[]> {
  return apiFetch("/api/v1/services", z.array(serviceStatusSchema));
}

export const jobStatusSchema = z.object({
  name: z.string(),
  provider: z.string(),
  schedule: z.string(),
  lastRun: z.string().optional(),
  lastError: z.string().optional(),
  runCount: z.number(),
});

export type JobStatusRecord = z.infer<typeof jobStatusSchema>;

export function fetchJobs(): Promise<JobStatusRecord[]> {
  return apiFetch("/api/v1/jobs", z.array(jobStatusSchema));
}

export const guildSchema = z.object({
  id: z.string(),
  name: z.string(),
  ownerId: z.string(),
  memberCount: z.number(),
  icon: z.string().optional(),
});

export type GuildRecord = z.infer<typeof guildSchema>;

export function fetchGuild(): Promise<GuildRecord> {
  return apiFetch("/api/v1/discord/guild", guildSchema);
}

export const channelSchema = z.object({
  id: z.string(),
  name: z.string(),
  type: z.string(),
  position: z.number(),
});

export const channelCategorySchema = z.object({
  id: z.string(),
  name: z.string(),
  position: z.number(),
  channels: z.array(channelSchema),
});

export type ChannelCategoryRecord = z.infer<typeof channelCategorySchema>;

export function fetchChannels(): Promise<ChannelCategoryRecord[]> {
  return apiFetch("/api/v1/discord/channels", z.array(channelCategorySchema));
}

export const roleSchema = z.object({
  id: z.string(),
  name: z.string(),
  color: z.number(),
  position: z.number(),
  managed: z.boolean(),
  memberCount: z.number(),
  permissions: z.array(z.string()),
  cjCanManage: z.boolean(),
  reasons: z.array(z.string()).optional(),
});

export type RoleRecord = z.infer<typeof roleSchema>;

export function fetchRoles(): Promise<RoleRecord[]> {
  return apiFetch("/api/v1/discord/roles", z.array(roleSchema));
}

export const personRefSchema = z.object({
  id: z.string(),
  username: z.string().optional(),
  nickname: z.string().optional(),
  avatarUrl: z.string().optional(),
  known: z.boolean(),
});

export type PersonRefRecord = z.infer<typeof personRefSchema>;

export const memberSchema = z.object({
  id: z.string(),
  username: z.string(),
  nickname: z.string().optional(),
  avatarUrl: z.string().optional(),
  roles: z.array(z.string()),
  joinedAt: z.string().optional(),
  cjCanActOn: z.boolean(),
  reasons: z.array(z.string()).optional(),
});

export type MemberRecord = z.infer<typeof memberSchema>;

export const membersPageSchema = z.object({
  members: z.array(memberSchema),
  total: z.number(),
  limit: z.number(),
  offset: z.number(),
  cachedAt: z.string(),
});

export type MembersPageRecord = z.infer<typeof membersPageSchema>;

export function fetchMembers(limit: number, offset: number, query?: string): Promise<MembersPageRecord> {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  if (query) params.set("q", query);
  return apiFetch(`/api/v1/discord/members?${params}`, membersPageSchema);
}

export const memberPermissionsSchema = z.object({
  memberId: z.string(),
  permissions: z.array(z.string()),
  isAdministrator: z.boolean(),
});

export type MemberPermissionsRecord = z.infer<typeof memberPermissionsSchema>;

export function fetchMemberPermissions(id: string): Promise<MemberPermissionsRecord> {
  return apiFetch(`/api/v1/discord/members/${id}/permissions`, memberPermissionsSchema);
}

export const auditLogEntrySchema = z.object({
  id: z.string(),
  actionType: z.number(),
  actionName: z.string(),
  targetId: z.string().optional(),
  userId: z.string().optional(),
  user: personRefSchema,
  reason: z.string().optional(),
});

export type AuditLogEntryRecord = z.infer<typeof auditLogEntrySchema>;

export function fetchAuditLog(): Promise<AuditLogEntryRecord[]> {
  return apiFetch("/api/v1/discord/audit-log", z.array(auditLogEntrySchema));
}

export const reactionCounterSchema = z.object({
  Counter: z.number(),
  Reaction: z.string(),
});

export const dataUserSchema = z.object({
  discord_user_id: z.string(),
  forum_user_id: z.string(),
  burger_user_id: z.string(),
  forum_user_name: z.string(),
  burger_user_name: z.string(),
  burgershot_verified: z.boolean(),
  received_reactions: z.array(reactionCounterSchema).nullish(),
  person: personRefSchema.optional(),
});

export type DataUserRecord = z.infer<typeof dataUserSchema>;

export const usersPageSchema = z.object({
  users: z.array(dataUserSchema),
  total: z.number(),
  limit: z.number(),
  offset: z.number(),
});

export type UsersPageRecord = z.infer<typeof usersPageSchema>;

export function fetchUsers(limit: number, offset: number, query?: string): Promise<UsersPageRecord> {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  if (query) params.set("q", query);
  return apiFetch(`/api/v1/data/users?${params}`, usersPageSchema);
}

export const chatLogSchema = z.object({
  timestamp: z.number(),
  discordUserId: z.string(),
  discordChannel: z.string(),
  message: z.string(),
  discordMessageId: z.string(),
  person: personRefSchema,
});

export type ChatLogRecord = z.infer<typeof chatLogSchema>;

export function searchMessages(query: string, userId?: string): Promise<ChatLogRecord[]> {
  const params = new URLSearchParams({ q: query });
  if (userId) params.set("userId", userId);
  return apiFetch(`/api/v1/data/chat?${params}`, z.array(chatLogSchema));
}

export function fetchMessageByID(id: string): Promise<ChatLogRecord> {
  return apiFetch(`/api/v1/data/messages/${id}`, chatLogSchema);
}

export const actionInfoSchema = z.object({
  name: z.string(),
  description: z.string(),
  risk: z.string(),
});

export type ActionInfoRecord = z.infer<typeof actionInfoSchema>;

export function fetchActions(): Promise<ActionInfoRecord[]> {
  return apiFetch("/api/v1/actions", z.array(actionInfoSchema));
}

export const actionOutcomeSchema = z.object({
  summary: z.string(),
  detail: z.record(z.string(), z.unknown()).optional(),
});

export type ActionOutcomeRecord = z.infer<typeof actionOutcomeSchema>;

export function previewAction(name: string, input: Record<string, unknown>): Promise<ActionOutcomeRecord> {
  return apiFetch(`/api/v1/actions/${name}/preview`, actionOutcomeSchema, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
}

export function executeAction(name: string, input: Record<string, unknown>): Promise<ActionOutcomeRecord> {
  return apiFetch(`/api/v1/actions/${name}/execute`, actionOutcomeSchema, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
}

export const actionRunSchema = z.object({
  id: z.string(),
  timestamp: z.string(),
  action: z.string(),
  risk: z.string(),
  mode: z.string(),
  outcome: z.string(),
  summary: z.string().optional(),
  errorMessage: z.string().optional(),
  durationMs: z.number(),
  requestId: z.string().optional(),
});

export type ActionRunRecord = z.infer<typeof actionRunSchema>;

export function fetchActionRuns(): Promise<ActionRunRecord[]> {
  return apiFetch("/api/v1/action-runs", z.array(actionRunSchema));
}

export function fetchCommands(): Promise<string[]> {
  return apiFetch("/api/v1/commands", z.array(z.string()));
}

export const profileMessageSchema = z.object({
  timestamp: z.number(),
  channelId: z.string(),
  channelName: z.string(),
  message: z.string(),
  messageId: z.string(),
});

export type ProfileMessageRecord = z.infer<typeof profileMessageSchema>;

export const profileSchema = z.object({
  discordUserId: z.string(),
  member: memberSchema.optional(),
  permissions: memberPermissionsSchema.optional(),
  user: dataUserSchema.optional(),
  recentMessages: z.array(profileMessageSchema),
  messageCount: z.number(),
  rank: z.number(),
});

export type ProfileRecord = z.infer<typeof profileSchema>;

export function fetchProfile(id: string): Promise<ProfileRecord> {
  return apiFetch(`/api/v1/profile/${id}`, profileSchema);
}

export const topMessagesEntrySchema = z.object({
  person: personRefSchema,
  messages: z.number(),
});

export type TopMessagesEntryRecord = z.infer<typeof topMessagesEntrySchema>;

export function fetchTopMessages(limit: number): Promise<TopMessagesEntryRecord[]> {
  return apiFetch(`/api/v1/leaderboards/top-messages?limit=${limit}`, z.array(topMessagesEntrySchema));
}

export const topReactionsEntrySchema = z.object({
  person: personRefSchema,
  counter: z.number(),
  reaction: z.string(),
});

export type TopReactionsEntryRecord = z.infer<typeof topReactionsEntrySchema>;

export function fetchTopReactions(limit: number): Promise<TopReactionsEntryRecord[]> {
  return apiFetch(`/api/v1/leaderboards/top-reactions?limit=${limit}`, z.array(topReactionsEntrySchema));
}
