package storage

import (
	"os"
	"testing"
)

var api *MongoStorer

func TestMain(m *testing.M) {
	var err error
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "localhost"
	}

	mongoPort := os.Getenv("MONGO_PORT")
	if mongoPort == "" {
		mongoPort = "27017"
	}

	mongoName := os.Getenv("MONGO_NAME")
	if mongoName == "" {
		mongoName = "local"
	}

	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")

	api, err = New(Config{
		MongoHost: mongoHost,
		MongoPort: mongoPort,
		MongoName: mongoName,
		MongoUser: mongoUser,
		MongoPass: mongoPass,
	})
	if err != nil {
		panic(err)
	}

	api.accounts.DropCollection()
	api.chat.DropCollection()

	ret := m.Run()
	os.Exit(ret)
}
