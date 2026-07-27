package admin

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/Southclaws/cj/admin/actions"
	"github.com/Southclaws/cj/admin/api"
	"github.com/Southclaws/cj/admin/logs"
	"github.com/Southclaws/cj/admin/readmodel"
	"github.com/Southclaws/cj/bot/commands"
	"github.com/Southclaws/cj/bot/heartbeat"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

var ErrDisabled = errors.New("dashboard is disabled")

const (
	defaultRequestTimeout = 10 * time.Second
	defaultActionTimeout  = 120 * time.Second
	defaultMaxBodyBytes   = 1 << 20
	defaultLogBufferSize  = 5000
	defaultLogRetention   = time.Hour
)

type dashboardConfig struct {
	ListenAddr   string
	AllowedHosts string
}

type Dependencies struct {
	Storage        storage.Storer
	Discord        *discord.Session
	Heartbeat      *heartbeat.Heartbeat
	CommandManager *commands.CommandManager
	LogBuffer      *logs.Buffer
	LogStream      *logs.Stream
}

type Server struct {
	httpServer *http.Server
	listener   net.Listener
	addr       string
	addrSource string
}

func New(cfg *types.Config, version string, deps Dependencies) (*Server, error) {
	if !cfg.DashboardEnabled {
		return nil, ErrDisabled
	}

	dc := dashboardConfig{
		ListenAddr:   cfg.DashboardListenAddr,
		AllowedHosts: cfg.DashboardAllowedHosts,
	}

	addr, source := resolveListenAddr(dc)

	allowedHosts, err := resolveAllowedHosts(dc, addr)
	if err != nil {
		return nil, err
	}

	logBuffer := deps.LogBuffer
	if logBuffer == nil {
		logBuffer = logs.NewBuffer(defaultLogBufferSize, defaultLogRetention)
	}
	logStream := deps.LogStream
	if logStream == nil {
		logStream = logs.NewStream()
	}

	servicesProvider := readmodel.NewServicesProvider(cfg, deps.Storage, deps.Discord)
	discordProvider := readmodel.NewDiscordProvider(deps.Discord)
	leaderboardsProvider := readmodel.NewLeaderboardsProvider(deps.Storage, discordProvider)
	usersProvider := readmodel.NewUsersProvider(deps.Storage, discordProvider)
	profileProvider := readmodel.NewProfileProvider(deps.Storage, discordProvider, usersProvider, leaderboardsProvider)

	registry := actions.NewRegistry()
	registry.Register(actions.NewTestServiceAction(servicesProvider))
	registry.Register(actions.NewWikiSearchAction())
	registry.Register(actions.NewTopMessagesAction(leaderboardsProvider))
	registry.Register(actions.NewTopReactionsAction(leaderboardsProvider))
	registry.Register(actions.NewUserRankAction(leaderboardsProvider))
	if deps.CommandManager != nil {
		registry.Register(actions.NewRefreshCommandsAction(deps.CommandManager))
		registry.Register(actions.NewGetCommandSettingsAction(deps.CommandManager))
		registry.Register(actions.NewSetCommandSettingsAction(deps.CommandManager))
	}
	if deps.Heartbeat != nil {
		registry.Register(actions.NewRefreshReadmeAction(deps.Heartbeat))
		registry.Register(actions.NewRunJobAction(deps.Heartbeat))
	}
	if deps.Discord != nil {
		registry.Register(actions.NewRefreshMembersAction(deps.Discord))
	}
	runner := actions.NewRunner(registry, deps.Storage)

	handler := api.NewRouter(api.Config{
		AllowedHosts:   allowedHosts,
		RequestTimeout: defaultRequestTimeout,
		ActionTimeout:  defaultActionTimeout,
		MaxBodyBytes:   defaultMaxBodyBytes,
		Status:         readmodel.NewStatusProvider(version, time.Now()),
		ConfigStatus:   readmodel.NewConfigStatusProvider(cfg),
		Settings:       readmodel.NewSettingsProvider(deps.Storage),
		Services:       servicesProvider,
		Jobs:           readmodel.NewJobsProvider(deps.Heartbeat),
		Discord:        discordProvider,
		Users:          usersProvider,
		Chat:           readmodel.NewChatProvider(deps.Storage, discordProvider),
		ActionRuns:     readmodel.NewActionRunsProvider(deps.Storage),
		Actions:        runner,
		Commands:       readmodel.NewCommandsProvider(deps.CommandManager),
		Profile:        profileProvider,
		Leaderboards:   leaderboardsProvider,
		LogBuffer:      logBuffer,
		LogStream:      logStream,
	})

	return &Server{
		httpServer: &http.Server{Handler: handler},
		addr:       addr,
		addrSource: source,
	}, nil
}

func (s *Server) Start(context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln

	zap.L().Info("dashboard listening",
		zap.String("addr", ln.Addr().String()),
		zap.String("source", s.addrSource))

	go func() {
		if err := s.httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("dashboard server stopped unexpectedly", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.addr
}
