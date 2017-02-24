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

	app = App{
		config:     Config{},
		httpClient: &http.Client{},
	}

	configLocation := os.Getenv("CONFIG_FILE")
	if configLocation == "" {
		configLocation = "config_example.json"
	}

	dbLocation := os.Getenv("DB_FILE")
	if dbLocation == "" {
		dbLocation = "users_test.db"
	}

	app.LoadConfig(configLocation)
	app.ConnectDB(dbLocation)

	log.Print("initialised.")

	ret := m.Run()

	app.db.Close()
	log.Printf("removing test database %s", dbLocation)
	if err := os.Remove(dbLocation); err != nil {
		panic(err)
	}

	os.Exit(ret)
}
