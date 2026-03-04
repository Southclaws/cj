package storage

import (
	"context"
	"os"
	"testing"
	"time"
)

var api *MongoStorer

func TestMain(m *testing.M) {
	var err error
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "127.0.0.1"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = api.accounts.Drop(ctx)
	_ = api.chat.Drop(ctx)

	ret := m.Run()
	_ = api.Close()
	os.Exit(ret)
}
