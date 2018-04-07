package main

import (
	"log"
	"net/http"
	"os"
	"testing"

	scraper "github.com/cardigann/go-cloudflare-scraper"
)

var app App

func TestMain(m *testing.M) {
	scrpr, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		log.Fatal(err)
	}

	app = App{
		config: Config{
			MongoHost: "localhost",
			MongoPort: "27017",
			MongoName: "cj",
			MongoUser: "root",
			MongoPass: "",
		},
		httpClient: &http.Client{Transport: scrpr},
	}

	app.ConnectDB()

	os.Exit(m.Run())
}
