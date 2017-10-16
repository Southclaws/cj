package main

import (
	"net/http"
	"os"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	app = App{
		config: Config{
			MongoHost: "localhost",
			MongoPort: "27017",
			MongoName: "cj",
			MongoUser: "root",
			MongoPass: "",
		},
		httpClient: &http.Client{},
	}

	app.ConnectDB()
	app.LoadLanguages()

	os.Exit(m.Run())
}
