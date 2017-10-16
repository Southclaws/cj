package main

import (
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

// App stores program state
type App struct {
	config         Config
	mongo          *mgo.Session
	accounts       *mgo.Collection
	chat           *mgo.Collection
	discordClient  *discordgo.Session
	httpClient     *http.Client
	ready          chan bool
	cache          *cache.Cache
	locale         Locale
	commandManager *CommandManager
}

// Start starts the app with the specified config and blocks until fatal error
func Start(config Config) {
	app := App{
		config:     Config{},
		httpClient: &http.Client{},
		cache:      cache.New(5*time.Minute, 30*time.Second),
	}

	configLocation := os.Getenv("CONFIG_FILE")
	if configLocation == "" {
		configLocation = "config.json"
	}

	logger.Debug("started with debug logging enabled",
		zap.Any("config", app.config))

	app.ConnectDB()
	app.LoadLanguages()
	app.StartCommandManager()
	app.ConnectDiscord()

	done := make(chan bool)
	<-done
}
