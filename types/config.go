package types

// Config stores configuration variables
//nolint:lll
type Config struct {
	Version               string
	MongoHost             string `split_words:"true" required:"false"`
	MongoPort             string `split_words:"true" required:"false"`
	MongoName             string `split_words:"true" required:"false"`
	MongoUser             string `split_words:"true" required:"false"`
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
	NoInitSync            bool   `split_words:"true" required:"false"` // if set, does not run database role sync on init
	NoDatabase            bool   `split_words:"true" required:"false"` // if set, does not connect to database
	WikiURL               string `split_words:"true" required:"false"` // wiki git repository url to clone
}
