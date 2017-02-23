package main

import (
	"log"
	"net/http"
	"os"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	log.Print("initialising...")
	var err error
	app = App{
		config:     Config{},
		httpClient: &http.Client{},
	}
	configLocation := os.Getenv("CONFIG_FILE")
	if configLocation == "" {
		configLocation = "config_example.json"
	}

	err = app.LoadConfig(configLocation)
	if err != nil {
		log.Fatal(err)
	}
	app.ConnectDB()

	log.Print("initialised.")

	os.Exit(m.Run())
}
