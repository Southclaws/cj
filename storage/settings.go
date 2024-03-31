package storage

import (
	"context"

	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
)

// SetCommandSettings upsets command settings
func (m *MongoStorer) SetCommandSettings(command string, settings types.CommandSettings) (err error) {
	settings.Command = command
	_, err = m.settings.Upsert(bson.M{"command": command}, settings)
	return
}

// GetCommandSettings returns command settings and uses a cache
func (m *MongoStorer) GetCommandSettings(command string) (settings types.CommandSettings, found bool, err error) {
	if c, ok := m.cache.Get(command); ok {
		if settings, ok = c.(types.CommandSettings); ok {
			return
		}
	}
	err = m.settings.Find(bson.M{"command": command}).One(&settings)
	switch err {
	case mgo.ErrNotFound:
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
	var readme Readme
	err = m.settings.Find(bson.M{"readme_message_id": bson.M{"$exists": true}}).One(&readme)
	if err != nil {
		return
	}

	message = readme.MessageID

	return
}

// FetchReadmeMessage fetches already sent message to upstream
func (m *MongoStorer) FetchReadmeMessage(gistID string, gistFile github.GistFilename) (message string, err error) {
	client := github.NewClient(nil)
	gist, _, err := client.Gists.Get(context.Background(), gistID)
	if err != nil {
		return
	}

	file, ok := gist.Files[gistFile]
	if !ok {
		err = errors.New("gist file not found (\"" + string(gistFile) + "\")")
		return
	}
	message = *(file.Content)

	return
}

// UpdateReadmeMessage updates the message in both channel and database
func (m *MongoStorer) UpdateReadmeMessage(session *discordgo.Session, original *discordgo.Message, upstream string) (err error) {
	if original.Content != upstream {
		session.ChannelMessageEdit(original.ChannelID, original.ID, upstream)
	} else {
		var readme Readme
		err = m.settings.Find(bson.M{"readme_message_id": bson.M{"$exists": true}}).One(&readme)
		if err == mgo.ErrNotFound {
			err = m.settings.Insert(bson.D{{Name: "readme_message_id", Value: original.ID}})
			if err != nil {
				return
			}
		} else if err != nil {
			return
		} else {
			err = m.settings.Update(bson.M{}, bson.M{"$set": bson.M{"readme_message_id": original.ID}})
			if err != nil {
				return
			}
		}
	}
	return
}
