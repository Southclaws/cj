package storage

import (
	"context"
	"encoding/base64"

	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v28/github"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SetCommandSettings upsets command settings
func (m *MongoStorer) SetCommandSettings(command string, settings types.CommandSettings) (err error) {
	settings.Command = command
	ctx, cancel := m.newContext()
	defer cancel()

	_, err = m.settings.UpdateOne(
		ctx,
		bson.M{"command": command},
		bson.M{"$set": settings},
		options.UpdateOne().SetUpsert(true),
	)
	return
}

// GetCommandSettings returns command settings and uses a cache
func (m *MongoStorer) GetCommandSettings(command string) (settings types.CommandSettings, found bool, err error) {
	if c, ok := m.cache.Get(command); ok {
		if settings, ok = c.(types.CommandSettings); ok {
			return
		}
	}
	ctx, cancel := m.newContext()
	defer cancel()

	err = m.settings.FindOne(ctx, bson.M{"command": command}).Decode(&settings)
	switch err {
	case mongo.ErrNoDocuments:
		err = nil
	case nil:
		found = true
	}
	return
}

// Readme holds the readme message id from the mongo database
type Readme struct {
	MessageID string `json:"readme_message_id" bson:"readme_message_id"`
}

// GetReadmeMessage gets the readme message id from the database
func (m *MongoStorer) GetReadmeMessage() (message string, err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	var readme Readme
	err = m.settings.FindOne(ctx, bson.M{"readme_message_id": bson.M{"$exists": true}}).Decode(&readme)
	if err != nil {
		return
	}

	message = readme.MessageID

	return
}

// FetchReadmeMessage fetches already sent message to upstream
func (m *MongoStorer) FetchReadmeMessage(githubOwner string, githubRepoistory string, fileName string) (message string, err error) {
	client := github.NewClient(nil)
	file, _, _, err := client.Repositories.GetContents(context.Background(), githubOwner, githubRepoistory, fileName, nil)
	if err != nil {
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(*file.Content)
	if err != nil {
		return
	}

	message = string(decoded)

	return
}

// UpdateReadmeMessage updates the message in both channel and database
func (m *MongoStorer) UpdateReadmeMessage(session *discordgo.Session, original *discordgo.Message, upstream string) (err error) {
	_, err = session.ChannelMessageEdit(original.ChannelID, original.ID, upstream)
	if err != nil {
		session.ChannelMessageSend("948604467887083550", err.Error())
	}

	ctx, cancel := m.newContext()
	defer cancel()

	var readme Readme
	err = m.settings.FindOne(ctx, bson.M{"readme_message_id": bson.M{"$exists": true}}).Decode(&readme)
	if err == mongo.ErrNoDocuments {
		_, err = m.settings.InsertOne(ctx, bson.D{{Key: "readme_message_id", Value: original.ID}})
		if err != nil {
			session.ChannelMessageSend("948604467887083550", err.Error())
			return
		}
	} else if err != nil {
		session.ChannelMessageSend("948604467887083550", err.Error())
		return
	} else {
		_, err = m.settings.UpdateOne(
			ctx,
			bson.M{"readme_message_id": bson.M{"$exists": true}},
			bson.M{"$set": bson.M{"readme_message_id": original.ID}},
		)
		if err != nil {
			session.ChannelMessageSend("948604467887083550", err.Error())
			return
		}
	}
	return
}
