package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/foize/go.fifo"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
	dbg "github.com/tj/go-debug"
)

var debug = dbg.Debug("main")

// Config stores configuration variables
type Config struct {
	LogFile          string `json:"log_file"`           // log file
	DiscordToken     string `json:"discord_token"`      // discord API token
	PrimaryChannel   string `json:"primary_channel"`    // main channel the bot hangs out in
	Heartbeat        int    `json:"heartbeat"`          // Heartbeat time in minutes, a heartbeat is when the bot chimes in to the server, sometimes with a random message
	BotID            string `json:"bot_id"`             // the bot's client ID
	Debug            bool   `json:"debug"`              // debug mode
	DebugUser        string `json:"debug_user"`         // when set, only accept commands from this user
	Admin            string `json:"admin"`              // user who has control over the bot
	LogFlushAt       int    `json:"log_flush_at"`       // size chat log can reach before being flushed to db
	LogFlushInterval int    `json:"log_flush_interval"` // interval between automatic chat log flushes
}

// App stores program state
type App struct {
	config        Config
	discordClient *discordgo.Session
	httpClient    *http.Client
	ready         chan bool
	cache         *gocache.Cache
	queue         *fifo.Queue
	locale        Locale
	chatLogger    *ChatLogger
	db            *gorm.DB
	done          chan bool
}

func main() {
	log.Print("Initialising SA:MP Forum Discord bot by Southclaws")

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

	f, err := os.Open(app.config.LogFile)
	if err != nil {
		f, err = os.Create(app.config.LogFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("Config:\n")
	log.Printf("- LogFile: %s\n", app.config.LogFile)
	log.Printf("- DiscordToken: (%d chars)\n", len(app.config.DiscordToken))
	log.Printf("- PrimaryChannel: %v\n", app.config.PrimaryChannel)
	log.Printf("- Heartbeat: %v\n", app.config.Heartbeat)
	log.Printf("- BotID: %v\n", app.config.BotID)
	log.Printf("- Debug: %v\n", app.config.Debug)
	log.Printf("- DebugUser: %v\n", app.config.DebugUser)
	log.Printf("- Admin: %v\n", app.config.Admin)
	log.Printf("- LogFlushAt: %v\n", app.config.LogFlushAt)
	log.Printf("- LogFlushInterval: %v\n", app.config.LogFlushInterval)
	log.Printf("~\n")

	app.ConnectDB(dbLocation)
	app.StartChatLogger()
	app.loadLanguages()

	var count int
	app.db.Model(&User{}).Count(&count)
	log.Printf("Verified users: %d", count)

	if app.config.Debug {
		dbg.Enable("main")
		debug("Debug mode enabled")
	}

	err = app.connect()
	if err != nil {
		log.Fatal(err)
	}

	app.done = make(chan bool)
	<-app.done

	err1 := app.discordClient.Close()
	err2 := app.db.Close()
	log.Printf("Closed database, shutting down %v %v", err1, err2)
}

// LoadConfig loads the specified config JSON file and returns the contents as
// a pointer to a Config object.
func (app *App) LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	app.config = Config{}
	err = json.NewDecoder(file).Decode(&app.config)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
