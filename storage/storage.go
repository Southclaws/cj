package storage

import (
	"fmt"

	"github.com/globalsign/mgo"

	"github.com/Southclaws/cj/types"
)

// Storer describes a type that is capable of persisting data
type Storer interface {
	RecordChatLog(discordUserID string, discordChannel string, message string) (err error)
	GetMessagesForUser(discordUserID string) (messages []ChatLog, err error)
	GetTopMessages(top int) (result TopMessages, err error)
	StoreVerifiedUser(verification types.Verification) (err error)
	SetLegacyUserToVerified(verification types.Verification) (err error)
	RemoveUser(id string) (err error)
	IsUserVerified(discordUserID string) (verified bool, err error)
	IsUserLegacyVerified(discordUserID string) (verified bool, err error)
	GetDiscordUserForumUser(forumUserID string) (discordUserID string, err error)
	GetForumUserFromDiscordUser(discordUserID string) (forumUserID string, err error)
	GetForumNameFromDiscordUser(discordUserID string) (forumUserName string, err error)
	GetDiscordUserFromForumName(forumName string) (discordUserID string, err error)
}

// MongoStorer exposes a storage MongoStorer for the bot
type MongoStorer struct {
	mongo    *mgo.Session
	accounts *mgo.Collection
	chat     *mgo.Collection
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
			Source:   config.MongoName,
			Username: config.MongoUser,
			Password: config.MongoPass,
		})
		if err != nil {
			return
		}
	}

	m.accounts = m.mongo.DB(config.MongoName).C("accounts")
	m.chat = m.mongo.DB(config.MongoName).C("chat")

	return
}
