package main

import (
	"net/http"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/bwmarrin/discordgo"
	"github.com/foize/go.fifo"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

// App stores program state
type App struct {
	config         Config
	mongo          *mgo.Session
	cln            *mgo.Collection
	discordClient  *discordgo.Session
	httpClient     *http.Client
	ready          chan bool
	cache          *cache.Cache
	queue          *fifo.Queue
	locale         Locale
	chatLogger     *ChatLogger
	commandManager *CommandManager
}

// Start starts the app with the specified config and blocks until fatal error
func Start(config Config) {
	app := App{
		config:     Config{},
		httpClient: &http.Client{},
		cache:      cache.New(5*time.Minute, 30*time.Second),
		queue:      fifo.NewQueue(),
	}

	configLocation := os.Getenv("CONFIG_FILE")
	if configLocation == "" {
		configLocation = "config.json"
	}

	logger.Debug("started with debug logging enabled",
		zap.Any("config", app.config))

	app.ConnectDB()
	app.StartChatLogger()
	app.LoadLanguages()
	app.StartCommandManager()
	app.ConnectDiscord()

	done := make(chan bool)
	<-done
}
