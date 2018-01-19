package main

import (
	"os"
	"strconv"

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
	MongoHost             string
	MongoPort             string
	MongoName             string
	MongoUser             string
	MongoPass             string
	DiscordToken          string // discord API token
	AdministrativeChannel string // administrative channel where someone can speak as bot
	PrimaryChannel        string // main channel the bot hangs out in
	Heartbeat             int    // Heartbeat time in minutes, a heartbeat is when the bot chimes in to the server, sometimes with a random message
	BotID                 string // the bot's client ID
	GuildID               string // the discord channel ID
	VerifiedRole          string // ID of the role for verified members
	DebugUser             string // when set, only accept commands from this user
	Admin                 string // user who has control over the bot
	LanguageData          string // `lang` data directory location, defaults to `./lang`
	Language              string // The code of the language used in bot
	NoInitSync            bool   // if set, does not run database role sync on init
}

func main() {
	Start(Config{
		MongoHost:             configStrFromEnv("MONGO_HOST"),
		MongoPort:             configStrFromEnv("MONGO_PORT"),
		MongoName:             configStrFromEnv("MONGO_NAME"),
		MongoUser:             configStrFromEnv("MONGO_USER"),
		MongoPass:             os.Getenv("MONGO_PASS"),
		DiscordToken:          configStrFromEnv("DISCORD_TOKEN"),
		AdministrativeChannel: configStrFromEnv("ADMINISTRATIVE_CHANNEL"),
		PrimaryChannel:        configStrFromEnv("PRIMARY_CHANNEL"),
		Heartbeat:             configIntFromEnv("HEARTBEAT"),
		BotID:                 configStrFromEnv("BOT_ID"),
		GuildID:               configStrFromEnv("GUILD_ID"),
		VerifiedRole:          configStrFromEnv("VERIFIED_ROLE"),
		DebugUser:             os.Getenv("DEBUG_USER"),
		Admin:                 configStrFromEnv("ADMIN"),
		LanguageData:          configStrFromEnv("LANGUAGE_DATA"),
		Language:              configStrFromEnv("LANGUAGE"),
		NoInitSync:            os.Getenv("NO_INIT_SYNC") != "",
	})
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
