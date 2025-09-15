package storage

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo"
	"github.com/patrickmn/go-cache"

	"github.com/Southclaws/cj/types"
)

// Storer describes a type that is capable of persisting data
type Storer interface {
	RecordChatLog(discordUserID, discordChannel, message, messageID string) (err error)
	GetMessagesForUser(discordUserID string) (messages []ChatLog, err error)
	GetTopMessages(top int) (result TopMessages, err error)
	GetUserMessageCount(discordUserID string) (messageCount int, err error)
	GetUserRank(discordUserID string) (rank int, err error)
	GetRandomMessage() (result ChatLog, err error)
	GetRandomMessageFromUsers([]string) (result ChatLog, err error)
	GetRandomUser() (result string, err error)
	GetMessageByID(messageID string) (message ChatLog, err error)

	GetUserOrCreate(discordUserID string) (result User, err error)
	UpdateUserUsername(discordUserID, username string) (err error)
	UpdateUser(user User) (err error)
	RemoveUser(id string) (err error)
	IsUserVerified(discordUserID string) (verified bool, err error)
	IsUserLegacyVerified(discordUserID string) (verified bool, err error)
	GetDiscordUserForumUser(forumUserID string) (discordUserID string, err error)
	GetForumUserFromDiscordUser(discordUserID string) (legacyUserID, burgerUserID string, err error)
	GetForumNameFromDiscordUser(discordUserID string) (legacyUserName, burgerUserName string, err error)
	GetDiscordUserFromForumName(forumName string) (legacyUserID, burgerUserID string, err error)
	AddEmojiReactionToUser(discordUserID, emoji string) (err error)
	RemoveEmojiReactionFromUser(discordUserID, emoji string) (err error)
	GetTopReactions(top int, reaction string) (result []TopReactionEntry, err error)

	SetCommandSettings(command string, settings types.CommandSettings) (err error)
	GetCommandSettings(command string) (settings types.CommandSettings, found bool, err error)

	GetReadmeMessage() (message string, err error)
	FetchReadmeMessage(githubOwner string, githubRepoistory string, fileName string) (message string, err error)
	UpdateReadmeMessage(session *discordgo.Session, original *discordgo.Message, upstream string) (err error)

	/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
	AddTrackedRole(discordUserID, roleID, roleName string) error
	RemoveTrackedRole(discordUserID, roleID string) error
	GetTrackedRoles(discordUserID string) ([]TrackedRole, error)
	*/
}

// MongoStorer exposes a storage MongoStorer for the bot
type MongoStorer struct {
	mongo    *mgo.Session
	accounts *mgo.Collection
	chat     *mgo.Collection
	settings *mgo.Collection
	cache    *cache.Cache
}

// Config represents database connection info
type Config struct {
	MongoHost string `split_words:"true" required:"true"`
	MongoPort string `split_words:"true" required:"true"`
	MongoName string `split_words:"true" required:"true"`
	MongoUser string `split_words:"true" required:"true"`
	MongoPass string `split_words:"true" required:"false"`
}

// New constructs a new storage API and connects to the database
func New(config Config) (m *MongoStorer, err error) {
	m = new(MongoStorer)
	m.mongo, err = mgo.Dial(fmt.Sprintf("%s:%s", config.MongoHost, config.MongoPort))
	if err != nil {
		return
	}

	if config.MongoPass != "" {
		err = m.mongo.Login(&mgo.Credential{
			Username: config.MongoUser,
			Password: config.MongoPass,
		})
		if err != nil {
			return
		}
	}

	m.accounts = m.mongo.DB(config.MongoName).C("accounts")
	m.chat = m.mongo.DB(config.MongoName).C("chat")
	m.settings = m.mongo.DB(config.MongoName).C("settings")
	m.cache = cache.New(time.Hour*24, time.Hour*12)

	return
}
