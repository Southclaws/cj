package storage

import (
	"os"
	"testing"
)

var api *API

func TestMain(m *testing.M) {
	var err error
	api, err = New(Config{
		MongoHost: "localhost",
		MongoPort: "27017",
		MongoName: "local",
		MongoUser: "",
		MongoPass: "",
	})
	if err != nil {
		panic(err)
	}

	api.accounts.DropCollection()
	api.chat.DropCollection()

	ret := m.Run()
	os.Exit(ret)
}
