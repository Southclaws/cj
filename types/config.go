package types

import (
	"time"
)

// Config stores configuration variables
//
//nolint:lll
type Config struct {
	Version                string
	MongoHost              string        `split_words:"true" required:"false"`
	MongoPort              string        `split_words:"true" required:"false"`
	MongoName              string        `split_words:"true" required:"false"`
	MongoUser              string        `split_words:"true" required:"false"`
	MongoPass              string        `split_words:"true" required:"false"`
	DiscordToken           string        `split_words:"true" required:"true"`                               // discord API token
	ReadmeChannel          string        `split_words:"true" required:"true"`                               // the default readme channel if not set in DB
	AdsChannel             string        `split_words:"true" required:"false" default:"763069609712025641"` // the default server ads channel if not set in DB
	DefaultChannel         string        `split_words:"true" required:"true"`                               // the default command channel if not set in DB
	DefaultRole            string        `split_words:"true" required:"true"`                               // the default command role if not set in DB
	DefaultCooldown        time.Duration `split_words:"true" required:"true"`                               // the default command cooldown if not set in DB
	Heartbeat              int           `split_words:"true" required:"true"`                               // Heartbeat time in minutes, a heartbeat is when the bot chimes in to the server, sometimes with a random message
	BotID                  string        `split_words:"true" required:"true"`                               // the bot's client ID
	GuildID                string        `split_words:"true" required:"true"`                               // the discord channel ID
	DebugUser              string        `split_words:"true" required:"false"`                              // when set, only accept commands from this user
	NoInitSync             bool          `split_words:"true" required:"false"`                              // if set, does not run database role sync on init
	NoDatabase             bool          `split_words:"true" required:"false"`                              // if set, does not connect to database
	ReadmeGithubOwner      string        `split_words:"true" required:"true" default:"openmultiplayer"`     // the readme repository owner
	ReadmeGithubRepository string        `split_words:"true" required:"true" default:"discord-rules"`       // the readme repository name
	ReadmeFileName         string        `split_words:"true" required:"true" default:"README.md"`           // the file name from readme in repository

	AdministrativeChannel string `split_words:"true" required:"false"` // DEPRECATED
	PrimaryChannel        string `split_words:"true" required:"false"` // DEPRECATED
	Admin                 string `split_words:"true" required:"false"` // DEPRECATED
}
