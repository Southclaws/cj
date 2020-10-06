package types

import (
	"time"

	"github.com/google/go-github/v28/github"
)

// Config stores configuration variables
//nolint:lll
type Config struct {
	Version            string
	MongoHost          string              `split_words:"true" required:"false"`
	MongoPort          string              `split_words:"true" required:"false"`
	MongoName          string              `split_words:"true" required:"false"`
	MongoUser          string              `split_words:"true" required:"false"`
	MongoPass          string              `split_words:"true" required:"false"`
	DiscordToken       string              `split_words:"true" required:"true"`                                                    // discord API token
	ReadmeChannel      string              `split_words:"true" required:"true"`                                                    // the default readme channel if not set in DB
	AdsChannel         string              `split_words:"true" required:"true" default:"763069609712025641"`                       // the default server ads channel if not set in DB
	DefaultChannel     string              `split_words:"true" required:"true"`                                                    // the default command channel if not set in DB
	DefaultRole        string              `split_words:"true" required:"true"`                                                    // the default command role if not set in DB
	DefaultCooldown    time.Duration       `split_words:"true" required:"true"`                                                    // the default command cooldown if not set in DB
	Heartbeat          int                 `split_words:"true" required:"true"`                                                    // Heartbeat time in minutes, a heartbeat is when the bot chimes in to the server, sometimes with a random message
	BotID              string              `split_words:"true" required:"true"`                                                    // the bot's client ID
	GuildID            string              `split_words:"true" required:"true"`                                                    // the discord channel ID
	VerifiedRole       string              `split_words:"true" required:"true"`                                                    // ID of the role for verified members
	DebugUser          string              `split_words:"true" required:"false"`                                                   // when set, only accept commands from this user
	NoInitSync         bool                `split_words:"true" required:"false"`                                                   // if set, does not run database role sync on init
	NoDatabase         bool                `split_words:"true" required:"false"`                                                   // if set, does not connect to database
	WikiURL            string              `split_words:"true" required:"false" default:"https://github.com/openmultiplayer/wiki"` // wiki git url
	ReadmeGistID       string              `split_words:"true" required:"true" default:"37f3548342e51365a7291132502b71b8"`         // the gist id for the readme
	ReadmeGistFileName github.GistFilename `split_words:"true" required:"true" default:"samp-discord-rules.md"`                    // the file name which contains the readme in the gist

	AdministrativeChannel string `split_words:"true" required:"false"` // DEPRECATED
	PrimaryChannel        string `split_words:"true" required:"false"` // DEPRECATED
	Admin                 string `split_words:"true" required:"false"` // DEPRECATED
}
