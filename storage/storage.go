package storage

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

// API exposes a storage API for the bot
type API struct {
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
func New(config Config) (api *API, err error) {
	api = new(API)
	api.mongo, err = mgo.Dial(fmt.Sprintf("%s:%s", config.MongoHost, config.MongoPort))
	if err != nil {
		return
	}

	if config.MongoPass != "" {
		err = api.mongo.Login(&mgo.Credential{
			Source:   config.MongoName,
			Username: config.MongoUser,
			Password: config.MongoPass,
		})
		if err != nil {
			return
		}
	}

	api.accounts = api.mongo.DB(config.MongoName).C("accounts")
	api.chat = api.mongo.DB(config.MongoName).C("chat")

	return
}
