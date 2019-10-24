package storage

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/go-github/v28/github"
)

// Readme to be written
type Readme struct {
	MessageID string `json:"readme_message_id" bson:"readme_message_id"`
}

// GetReadmeMessage gets the readme message id from the database
func (m *MongoStorer) GetReadmeMessage() (message string, err error) {
	var readme Readme
	err = m.chat.Find(bson.M{"readme_message_id": bson.M{"$exists": true}}).One(&readme)
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

	message = *(gist.Files[gistFile].Content)

	return
}

// UpdateReadmeMessage updates the message in both channel and database
func (m *MongoStorer) UpdateReadmeMessage(session *discordgo.Session, original *discordgo.Message, upstream string) (err error) {
	if original.Content != upstream {
		session.ChannelMessageEdit(original.ChannelID, original.ID, upstream)
	} else {
		var readme Readme
		err = m.chat.Find(bson.M{"readme_message_id": bson.M{"$exists": true}}).One(&readme)
		if err == mgo.ErrNotFound {
			err = m.chat.Insert(bson.D{{Name: "readme_message_id", Value: original.ID}})
			if err != nil {
				return
			}
		} else if err != nil {
			return
		} else {
			err = m.chat.Update(bson.M{}, bson.M{"$set": bson.M{"readme_message_id": original.ID}})
			if err != nil {
				return
			}
		}
	}
	return
}
