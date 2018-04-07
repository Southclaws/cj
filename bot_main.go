package main

import (
	"os"
	"strconv"

	// loads environment variables from .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	var config zap.Config
	debug := os.Getenv("DEBUG")

	if os.Getenv("TESTING") != "" {
		config = zap.NewDevelopmentConfig()
		config.DisableCaller = true
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.MessageKey = "@message"
		config.EncoderConfig.TimeKey = "@timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		if debug != "0" && debug != "" {
			dyn := zap.NewAtomicLevel()
			dyn.SetLevel(zap.DebugLevel)
			config.Level = dyn
		}
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
	MongoHost             string `split_words:"true" required:"true"`
	MongoPort             string `split_words:"true" required:"true"`
	MongoName             string `split_words:"true" required:"true"`
	MongoUser             string `split_words:"true" required:"true"`
	MongoPass             string `split_words:"true" required:"false"`
	DiscordToken          string `split_words:"true" required:"true"`  // discord API token
	AdministrativeChannel string `split_words:"true" required:"true"`  // administrative channel where someone can speak as bot
	PrimaryChannel        string `split_words:"true" required:"true"`  // main channel the bot hangs out in
	Heartbeat             int    `split_words:"true" required:"true"`  // Heartbeat time in minutes, a heartbeat is when the bot chimes in to the server, sometimes with a random message
	BotID                 string `split_words:"true" required:"true"`  // the bot's client ID
	GuildID               string `split_words:"true" required:"true"`  // the discord channel ID
	VerifiedRole          string `split_words:"true" required:"true"`  // ID of the role for verified members
	DebugUser             string `split_words:"true" required:"false"` // when set, only accept commands from this user
	Admin                 string `split_words:"true" required:"true"`  // user who has control over the bot
	LanguageData          string `split_words:"true" required:"true"`  // `lang` data directory location, defaults to `./lang`
	Language              string `split_words:"true" required:"true"`  // The code of the language used in bot
	NoInitSync            bool   `split_words:"true" required:"false"` // if set, does not run database role sync on init
}

func main() {
	config := Config{}
	err := envconfig.Process("CJ", &config)
	if err != nil {
		logger.Fatal("failed to load configuration",
			zap.Error(err))
	}

	Start(config)
}

func configStrFromEnv(name string) (value string) {
	value = os.Getenv(name)
	if value == "" {
		logger.Fatal("environment variable not set",
			zap.String("name", name))
	}
	return
}

func configIntFromEnv(name string) (value int) {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		logger.Fatal("environment variable not set",
			zap.String("name", name))
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logger.Fatal("failed to convert environment variable to int",
			zap.Error(err),
			zap.String("name", name))
	}
	return
}
