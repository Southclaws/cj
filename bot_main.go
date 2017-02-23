package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	gocache "github.com/patrickmn/go-cache"
	dbg "github.com/tj/go-debug"
)

var debug = dbg.Debug("main")

// Config stores configuration variables
type Config struct {
	DiscordToken   string `json:"discord_token"`   // discord API token
	PrimaryChannel string `json:"primary_channel"` // main channel the bot hangs out in
	Heartbeat      int    `json:"heartbeat"`       // Heartbeat time in minutes, a heartbeat is when the bot chimes in to the server, sometimes with a random message
	BotID          string `json:"bot_id"`          // the bot's client ID
	Debug          bool   `json:"debug"`           // debug mode
	DebugUser      string `json:"debug_user"`      // when set, only accept commands from this user
}

// App stores program state
type App struct {
	config        Config
	discordClient *discordgo.Session
	httpClient    *http.Client
	ready         chan bool
	cache         *gocache.Cache
	locale        Locale
}

func main() {
	log.Print("Initialising SA:MP Forum Discord bot by Southclaws")

	var err error
	app := App{
		config:     Config{},
		httpClient: &http.Client{},
		cache:      gocache.New(5*time.Minute, 30*time.Second),
	}
	configLocation := os.Getenv("CONFIG_FILE")
	if configLocation == "" {
		configLocation = "config.json"
	}

	err = loadConfig(configLocation, &app.config)
	if err != nil {
		log.Fatal(err)
	}

	app.loadLanguages()

	log.Printf("Config:\n")
	log.Printf("- DiscordToken: (%d chars)\n", len(app.config.DiscordToken))
	log.Printf("- PrimaryChannel: %s\n", app.config.PrimaryChannel)
	log.Printf("- Debug: %v\n", app.config.Debug)
	log.Printf("~\n")

	log.Print(app.locale.GetLangString("en", "WarnUserError", "test"))

	if app.config.Debug {
		dbg.Enable("main")
		debug("Debug mode enabled")
	}

	err = app.connect()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	<-done
}

func loadConfig(filename string, cfg *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	json.NewDecoder(file).Decode(&cfg)

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
