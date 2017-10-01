package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
	"github.com/foize/go.fifo"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
)

var logger *zap.Logger

func initLogger(debug bool) {
	config := zap.NewDevelopmentConfig()
	config.DisableCaller = true

	if debug {
		dyn := zap.NewAtomicLevel()
		dyn.SetLevel(zap.DebugLevel)
		config.Level = dyn
	}

	_logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	logger = _logger.With(
		zap.String("@version", os.Getenv("GIT_HASH")),
		zap.Namespace("@fields"),
	)
}

// Config stores configuration variables
type Config struct {
	DiscordToken          string `json:"discord_token"`          // discord API token
	AdministrativeChannel string `json:"administrative_channel"` // administrative channel where someone can speak as bot
	PrimaryChannel        string `json:"primary_channel"`        // main channel the bot hangs out in
	Heartbeat             int    `json:"heartbeat"`              // Heartbeat time in minutes, a heartbeat is when the bot chimes in to the server, sometimes with a random message
	BotID                 string `json:"bot_id"`                 // the bot's client ID
	GuildID               string `json:"guild_id"`               // the discord channel ID
	VerifiedRole          string `json:"verified_role"`          // ID of the role for verified members
	NormalRole            string `json:"normal_role"`            // role assigned to all users automatically
	DebugLogs             bool   `json:"debug_logs"`             // debug logging
	DebugUser             string `json:"debug_user"`             // when set, only accept commands from this user
	Admin                 string `json:"admin"`                  // user who has control over the bot
	LogFlushAt            int    `json:"log_flush_at"`           // size chat log can reach before being flushed to db
	LogFlushInterval      int    `json:"log_flush_interval"`     // interval between automatic chat log flushes
}

// App stores program state
type App struct {
	config         Config
	discordClient  *discordgo.Session
	httpClient     *http.Client
	ready          chan bool
	cache          *gocache.Cache
	queue          *fifo.Queue
	locale         Locale
	chatLogger     *ChatLogger
	commandManager *CommandManager
	db             *gorm.DB
	done           chan bool
}

func main() {
	var err error
	app := App{
		config:     Config{},
		httpClient: &http.Client{},
		cache:      gocache.New(5*time.Minute, 30*time.Second),
		queue:      fifo.NewQueue(),
	}

	configLocation := os.Getenv("CONFIG_FILE")
	if configLocation == "" {
		configLocation = "config.json"
	}

	dbLocation := os.Getenv("DB_FILE")
	if dbLocation == "" {
		dbLocation = "users.db"
	}

	app.LoadConfig(configLocation)
	initLogger(app.config.DebugLogs)
	logger.Debug("started with debug logging enabled",
		zap.Any("config", app.config))

	app.ConnectDB(dbLocation)
	app.StartChatLogger()
	app.loadLanguages()
	app.StartCommandManager()

	var count int
	app.db.Model(&User{}).Count(&count)

	err = app.connect()
	if err != nil {
		panic(err)
	}

	app.done = make(chan bool)
	<-app.done

	err1 := app.discordClient.Close()
	err2 := app.db.Close()
	logger.Fatal("shutting down",
		zap.Error(err1),
		zap.Error(err2))
}

// LoadConfig loads the specified config JSON file and returns the contents as
// a pointer to a Config object.
func (app *App) LoadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	app.config = Config{}
	err = json.NewDecoder(file).Decode(&app.config)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
}
