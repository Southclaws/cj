package api

import (
	"net/http"
	"time"

	"github.com/Southclaws/cj/admin/actions"
	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/logs"
	"github.com/Southclaws/cj/admin/middleware"
	"github.com/Southclaws/cj/admin/readmodel"
)

type Config struct {
	AllowedHosts   []string
	RequestTimeout time.Duration
	ActionTimeout  time.Duration
	MaxBodyBytes   int64
	Status         *readmodel.StatusProvider
	ConfigStatus   *readmodel.ConfigStatusProvider
	Settings       *readmodel.SettingsProvider
	Services       *readmodel.ServicesProvider
	Jobs           *readmodel.JobsProvider
	Discord        *readmodel.DiscordProvider
	Users          *readmodel.UsersProvider
	Chat           *readmodel.ChatProvider
	ActionRuns     *readmodel.ActionRunsProvider
	Actions        *actions.Runner
	Commands       *readmodel.CommandsProvider
	Profile        *readmodel.ProfileProvider
	Leaderboards   *readmodel.LeaderboardsProvider
	LogBuffer      *logs.Buffer
	LogStream      *logs.Stream
}

func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/status", handleStatus(cfg.Status))
	mux.HandleFunc("GET /api/v1/config/status", handleConfigStatus(cfg.ConfigStatus))
	mux.HandleFunc("GET /api/v1/data/settings", handleGetSettings(cfg.Settings))
	mux.HandleFunc("PUT /api/v1/data/settings", handlePutSettings(cfg.Settings))
	mux.HandleFunc("GET /api/v1/logs", handleLogsList(cfg.LogBuffer))
	mux.HandleFunc("GET /api/v1/services", handleServices(cfg.Services))
	mux.HandleFunc("GET /api/v1/jobs", handleJobs(cfg.Jobs))
	mux.HandleFunc("GET /api/v1/discord/guild", handleDiscordGuild(cfg.Discord))
	mux.HandleFunc("GET /api/v1/discord/channels", handleDiscordChannels(cfg.Discord))
	mux.HandleFunc("GET /api/v1/discord/roles", handleDiscordRoles(cfg.Discord))
	mux.HandleFunc("GET /api/v1/discord/members", handleDiscordMembers(cfg.Discord))
	mux.HandleFunc("GET /api/v1/discord/members/{id}", handleDiscordMember(cfg.Discord))
	mux.HandleFunc("GET /api/v1/discord/members/{id}/permissions", handleDiscordMemberPermissions(cfg.Discord))
	mux.HandleFunc("GET /api/v1/discord/audit-log", handleDiscordAuditLog(cfg.Discord))
	mux.HandleFunc("GET /api/v1/data/users", handleDataUsers(cfg.Users))
	mux.HandleFunc("GET /api/v1/data/users/{id}", handleDataUser(cfg.Users))
	mux.HandleFunc("GET /api/v1/data/chat", handleDataChat(cfg.Chat))
	mux.HandleFunc("GET /api/v1/data/messages/{id}", handleDataMessage(cfg.Chat))
	mux.HandleFunc("GET /api/v1/actions", handleActionsList(cfg.Actions))
	mux.HandleFunc("GET /api/v1/actions/{name}", handleActionGet(cfg.Actions))
	mux.HandleFunc("POST /api/v1/actions/{name}/preview", handleActionPreview(cfg.Actions))
	mux.HandleFunc("GET /api/v1/action-runs", handleActionRunsList(cfg.ActionRuns))
	mux.HandleFunc("GET /api/v1/action-runs/{id}", handleActionRunGet(cfg.ActionRuns))
	mux.HandleFunc("GET /api/v1/commands", handleCommandsList(cfg.Commands))
	mux.HandleFunc("GET /api/v1/profile/{id}", handleProfile(cfg.Profile))
	mux.HandleFunc("GET /api/v1/leaderboards/top-messages", handleLeaderboardsTopMessages(cfg.Leaderboards))
	mux.HandleFunc("GET /api/v1/leaderboards/top-reactions", handleLeaderboardsTopReactions(cfg.Leaderboards))
	mux.HandleFunc("/api/", handleAPINotFound)
	mux.Handle("/", staticHandler())

	timed := middleware.Timeout(cfg.RequestTimeout)(mux)

	top := http.NewServeMux()
	top.HandleFunc("GET /api/v1/logs/stream", handleLogsStream(cfg.LogBuffer, cfg.LogStream))
	top.Handle("POST /api/v1/actions/{name}/execute",
		middleware.Timeout(cfg.ActionTimeout)(handleActionExecute(cfg.Actions)))
	top.Handle("/", timed)

	var handler http.Handler = top
	handler = middleware.MaxBytes(cfg.MaxBodyBytes)(handler)
	handler = middleware.NetworkGuard(cfg.AllowedHosts)(handler)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Recover(handler)
	return handler
}

func handleStatus(status *readmodel.StatusProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apierrors.WriteJSON(w, http.StatusOK, status.Status())
	}
}

func handleConfigStatus(status *readmodel.ConfigStatusProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apierrors.WriteJSON(w, http.StatusOK, status.Status())
	}
}
