package bot

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Southclaws/cj/admin"
	"github.com/Southclaws/cj/admin/logs"
	"github.com/Southclaws/cj/bot/admod"
	"github.com/Southclaws/cj/bot/commands"
	"github.com/Southclaws/cj/bot/heartbeat"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// App stores program state
type App struct {
	config         *types.Config
	discordClient  *discord.Session
	storage        storage.Storer
	forum          *forum.ForumClient
	ready          chan error
	extensions     []Extension
	channels       map[string]*discordgo.Channel
	admin          *admin.Server
	heartbeat      *heartbeat.Heartbeat
	commandManager *commands.CommandManager
}

// Extension represents an extension to the bot that receives a pointer to the
// storage backend.
type Extension interface {
	Init(*types.Config, *discord.Session, storage.Storer, *forum.ForumClient) error
	OnMessage(discordgo.Message) error
}

// Start starts the app with the specified config and blocks until fatal error
func Start(config *types.Config) {
	app := App{
		config:   config,
		ready:    make(chan error),
		channels: make(map[string]*discordgo.Channel),
	}

	var err error

	if !config.NoDatabase {
		app.storage, err = storage.New(storage.Config{
			MongoHost: config.MongoHost,
			MongoPort: config.MongoPort,
			MongoName: config.MongoName,
			MongoUser: config.MongoUser,
			MongoPass: config.MongoPass,
		})
		if err != nil {
			zap.L().Fatal("failed to connect to database", zap.Error(err))
		}
	} else {
		app.storage = &storage.Memory{}
	}

	app.forum, err = forum.NewForumClient()
	if err != nil {
		zap.L().Fatal("failed to initialise forum client", zap.Error(err))
	}

	err = app.ConnectDiscord()
	if err != nil {
		zap.L().Fatal("failed to connect to discord", zap.Error(err))
	}

	app.heartbeat = &heartbeat.Heartbeat{}
	app.commandManager = &commands.CommandManager{}
	app.extensions = []Extension{
		app.commandManager,
		app.heartbeat,
		&admod.Watcher{},
	}

	for _, ex := range app.extensions {
		zap.L().Debug("initialising extension")
		err = ex.Init(config, app.discordClient, app.storage, app.forum)
		if err != nil {
			zap.L().Fatal("failed to initialise extension", zap.Error(err))
		}
	}

	zap.L().Debug("started with debug logging enabled",
		zap.Int("extensions", len(app.extensions)),
		zap.Any("config", config))

	app.startDashboard(config)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals

	if app.admin != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := app.admin.Shutdown(ctx); err != nil {
			zap.L().Error("failed to shut down dashboard", zap.Error(err))
		}
		cancel()
	}

	if closer, ok := app.storage.(interface{ Close() error }); ok {
		err = closer.Close()
		if err != nil {
			zap.L().Error("failed to close storage", zap.Error(err))
		}
	}
}

func (app *App) startDashboard(config *types.Config) {
	if !config.DashboardEnabled {
		return
	}

	logBuffer := logs.NewBuffer(config.DashboardLogBufferSize, config.DashboardLogRetention)
	logStream := logs.NewStream()
	logRedactor := logs.NewRedactor(config.DiscordToken, config.MongoPass)
	logCore := logs.NewCore(zap.NewAtomicLevelAt(zap.DebugLevel), logBuffer, logStream, logRedactor)

	zap.ReplaceGlobals(zap.L().WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, logCore)
	})))

	server, err := admin.New(config, config.Version, admin.Dependencies{
		Storage:        app.storage,
		Discord:        app.discordClient,
		Heartbeat:      app.heartbeat,
		CommandManager: app.commandManager,
		LogBuffer:      logBuffer,
		LogStream:      logStream,
	})
	if errors.Is(err, admin.ErrDisabled) {
		return
	}
	if err != nil {
		zap.L().Error("dashboard configuration invalid, dashboard disabled", zap.Error(err))
		return
	}

	if err := server.Start(context.Background()); err != nil {
		zap.L().Error("failed to start dashboard", zap.Error(err))
		return
	}

	app.admin = server
}
