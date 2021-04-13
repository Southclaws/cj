package storage

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// ChatLog represents a single logged chat message from Discord
type ChatLog struct {
	Timestamp        int64
	DiscordUserID    string
	DiscordChannel   string
	Message          string
	DiscordMessageID string
}

// RecordChatLog records a chat message from a user in a channel
func (m *MongoStorer) RecordChatLog(discordUserID, discordChannel, message, messageID string) (err error) {
	record := ChatLog{
		time.Now().Unix(),
		discordUserID,
		discordChannel,
		message,
		messageID,
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

func (m *MongoStorer) GetUserMessageCount(discordUserID string) (messageCount int, err error) {
	messageCount, err = m.chat.Find(bson.M{"discorduserid": discordUserID}).Count()
	return
}

// GetUserRank returns rank according to most messages sent with the top user having rank 1
func (m *MongoStorer) GetUserRank(discordUserID string) (rank int, err error) {
	messageCount, err := m.chat.Find(bson.M{"discorduserid": discordUserID}).Count()

	myPipe := m.chat.Pipe([]bson.M{
		bson.M{
			"$group": bson.M{
				"_id": "$discorduserid",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$match": bson.M{"count": bson.M{"$gt": messageCount}},
		},
		bson.M{
			"$count": "rank",
		},
	})
	tempMap := make(map[string]int)
	err = myPipe.One(tempMap)
	rank = tempMap["rank"] + 1 // +1 otherwise the top user will have rank 0
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

func (m *MongoStorer) GetRandomMessageFromUsers(users []string) (result ChatLog, err error) {
	err = m.chat.Pipe([]bson.M{
		{"$match": bson.M{
			"DiscordUserID": bson.M{
				"$in": users,
			},
		}},
		{"$sample": bson.M{
			"size": 1,
		}},
	}).One(&result)
	return
}

func (m *MongoStorer) GetMessageByID(messageID string) (message ChatLog, err error) {
	err = m.chat.Find(bson.M{"discordmessageid": messageID}).One(&message)
	return message, err
}
