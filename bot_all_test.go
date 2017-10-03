package main

import (
	"net/http"
	"os"
	"testing"
)

var app App

func TestMain(m *testing.M) {
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

	app.config = LoadConfig(configLocation)
	app.ConnectDB()
	app.loadLanguages()

	ret := m.Run()

	if err := os.Remove(dbLocation); err != nil {
		panic(err)
	}

	os.Exit(ret)
}
