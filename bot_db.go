package main

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// ConnectDB connects the app to the database
func (app *App) ConnectDB(dbpath string) {
	db, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		logger.Fatal("failed to open sqlite database", zap.Error(err))
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&ChatLog{})

	app.db = db
}
