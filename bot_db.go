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
	db.LogMode(app.config.DebugLogs)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&ChatLog{})

	app.db = db
	log.Printf("Connected to database '%s'", dbpath)
}
