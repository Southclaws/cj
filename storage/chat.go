package storage

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// ChatLog represents a single logged chat message from Discord
type ChatLog struct {
	Timestamp      int64
	DiscordUserID  string
	DiscordChannel string
	Message        string
}

// RecordChatLog records a chat message from a user in a channel
func (m *MongoStorer) RecordChatLog(discordUserID string, discordChannel string, message string) (err error) {
	record := ChatLog{
		time.Now().Unix(),
		discordUserID,
		discordChannel,
		message,
	}

	err = m.chat.Insert(record)
	if err != nil {
		err = errors.Wrap(err, "failed to insert chat log")
	}

	return
}

// GetMessagesForUser returns all messages from the given discord user.
func (m *MongoStorer) GetMessagesForUser(discordUserID string) (messages []ChatLog, err error) {
	err = m.chat.Find(bson.M{"discorduserid": discordUserID}).All(&messages)
	return
}

// TopMessages is a list of users with the most messages
type TopMessages []TopMessagesEntry

// TopMessagesEntry is a user and their message count
type TopMessagesEntry struct {
	User     string `bson:"_id"`
	Messages int    `bson:"count"`
}

func (s TopMessages) Len() int           { return len(s) }
func (s TopMessages) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TopMessages) Less(i, j int) bool { return s[i].Messages < s[j].Messages }

// GetTopMessages returns n users with the most messages
func (m *MongoStorer) GetTopMessages(top int) (result TopMessages, err error) {
	err = m.chat.Pipe([]bson.M{
		bson.M{
			"$group": bson.M{
				"_id": "$discorduserid",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"count": -1,
			},
		},
		bson.M{
			"$limit": top,
		},
	}).All(&result)
	// sort.Sort(sort.Reverse(result))
	return
}

// GetRandomMessage returns a random message from the database.
func (m *MongoStorer) GetRandomMessage() (log ChatLog, err error) {
	err = m.chat.Pipe([]bson.M{
		{"$sample": bson.M{
			"size": 1,
		}},
	}).One(&log)
	return
}
