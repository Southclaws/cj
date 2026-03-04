package storage

import (
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	ctx, cancel := m.newContext()
	defer cancel()

	_, err = m.chat.InsertOne(ctx, record)
	if err != nil {
		err = errors.Wrap(err, "failed to insert chat log")
	}

	return
}

// GetMessagesForUser returns all messages from the given discord user.
func (m *MongoStorer) GetMessagesForUser(discordUserID string) (messages []ChatLog, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	cursor, err := m.chat.Find(ctx, bson.M{"discorduserid": discordUserID})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &messages)
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
	ctx, cancel := m.newContext()
	defer cancel()

	cursor, err := m.chat.Aggregate(ctx, []bson.M{
		{
			"$group": bson.M{
				"_id": "$discorduserid",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$sort": bson.M{
				"count": -1,
			},
		},
		{
			"$limit": top,
		},
	})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &result)
	// sort.Sort(sort.Reverse(result))
	return
}

func (m *MongoStorer) GetUserMessageCount(discordUserID string) (messageCount int, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	var count int64
	count, err = m.chat.CountDocuments(ctx, bson.M{"discorduserid": discordUserID})
	messageCount = int(count)
	return
}

// GetUserRank returns rank according to most messages sent with the top user having rank 1
func (m *MongoStorer) GetUserRank(discordUserID string) (rank int, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	count, err := m.chat.CountDocuments(ctx, bson.M{"discorduserid": discordUserID})
	if err != nil {
		return
	}

	myPipe, err := m.chat.Aggregate(ctx, []bson.M{
		{
			"$group": bson.M{
				"_id": "$discorduserid",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$match": bson.M{"count": bson.M{"$gt": count}},
		},
		{
			"$count": "rank",
		},
	})
	if err != nil {
		return
	}
	defer myPipe.Close(ctx)

	tempMap := make(map[string]int)
	if myPipe.Next(ctx) {
		err = myPipe.Decode(&tempMap)
		if err != nil {
			return
		}
	} else if myPipe.Err() != nil {
		err = myPipe.Err()
		return
	}
	rank = tempMap["rank"] + 1 // +1 otherwise the top user will have rank 0
	return
}

// GetRandomMessage returns a random message from the database.
func (m *MongoStorer) GetRandomMessage() (log ChatLog, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	cursor, err := m.chat.Aggregate(ctx, []bson.M{
		{"$sample": bson.M{
			"size": 1,
		}},
	})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&log)
		return
	}
	err = cursor.Err()
	return
}

func (m *MongoStorer) GetRandomMessageFromUsers(users []string) (result ChatLog, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	cursor, err := m.chat.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"discorduserid": bson.M{
				"$in": users,
			},
		}},
		{"$sample": bson.M{
			"size": 1,
		}},
	})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		return
	}
	err = cursor.Err()
	return
}

func (m *MongoStorer) GetMessageByID(messageID string) (message ChatLog, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	err = m.chat.FindOne(ctx, bson.M{"discordmessageid": messageID}).Decode(&message)
	return message, err
}
