package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

// ConnectDB connects the app to the database
func (app *App) ConnectDB(dbpath string) {
	db, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(app.config.Debug)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&ChatLog{})

	app.db = db
	log.Print("Connected to database")
}
