package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/bot/commands"
	"github.com/Southclaws/cj/bot/stats"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// App stores program state
type App struct {
	config        *types.Config
	discordClient *discord.Session
	storage       storage.Storer
	forum         *forum.ForumClient
	ready         chan error
	extensions    []Extension
	channels      map[string]*discordgo.Channel
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

	app.storage, err = storage.New(storage.Config{
		MongoHost: config.MongoHost,
		MongoPort: config.MongoPort,
		MongoName: config.MongoName,
		MongoUser: config.MongoUser,
		MongoPass: config.MongoPass,
	})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	app.forum, err = forum.NewForumClient()
	if err != nil {
		logger.Fatal("failed to initialise forum client", zap.Error(err))
	}

	err = app.ConnectDiscord()
	if err != nil {
		logger.Fatal("failed to connect to discord", zap.Error(err))
	}

	app.extensions = []Extension{
		&commands.CommandManager{},
		&stats.Aggregator{},
	}

	for _, ex := range app.extensions {
		err = ex.Init(config, app.discordClient, app.storage, app.forum)
		if err != nil {
			logger.Fatal("failed to initialise extension", zap.Error(err))
		}
	}

	app.forum.NewPostAlert("3", func() {
		//nolint:errcheck
		app.discordClient.S.ChannelMessageSend(
			config.PrimaryChannel,
			"New Kalcor Post: http://forum.sa-mp.com/search.php?do=finduser&u=3",
		)
	})

	_, err = app.discordClient.S.ChannelMessageSend(
		config.PrimaryChannel,
		fmt.Sprintf("Hey, what's cracking now? CJ initialised with version %s", config.Version))
	if err != nil {
		logger.Fatal("failed to send initialisation message", zap.Error(err))
	}

	logger.Debug("started with debug logging enabled",
		zap.Int("extensions", len(app.extensions)),
		zap.Any("config", config))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGKILL)
	<-signals
}
