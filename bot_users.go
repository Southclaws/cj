package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// User is a recorded and verified SA:MP forum user.
type User struct {
	DiscordUserID    string `gorm:"primary_key;index;not null;unique"`
	ForumUserID      string `gorm:"not null;unique"`
	VerificationCode string `gorm:"not null"`
}

// ConnectDB connects the app to the database
func (app *App) ConnectDB(dbpath string) {
	db, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	db.AutoMigrate(&User{})

	app.db = db
	log.Print("Connected to database")
}

// StoreVerifiedUser is for when a user finishes their verification.
func (app App) StoreVerifiedUser(verification Verification) error {
	err := app.db.Create(&User{
		DiscordUserID:    verification.discordUser.ID,
		ForumUserID:      verification.forumUser,
		VerificationCode: verification.code,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

// IsUserVerified returns a discord user, a blank string or an error
func (app App) IsUserVerified(discordUserID string) (bool, error) {
	var count int
	result := app.db.Model(&User{}).Where(&User{DiscordUserID: discordUserID}).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// GetDiscordUserForumUser returns a discord user, a blank string or an error
func (app App) GetDiscordUserForumUser(forumUserID string) (string, error) {
	var user User

	err := app.db.First(&user, &User{ForumUserID: forumUserID}).Error
	if err != nil {
		return "", err
	}

	return user.DiscordUserID, nil
}

// GetForumUserFromDiscordUser returns a discord user, a blank string or an error
func (app App) GetForumUserFromDiscordUser(discordUserID string) (string, error) {
	var user User

	err := app.db.First(&user, &User{DiscordUserID: discordUserID}).Error
	if err != nil {
		return "", err
	}

	return user.ForumUserID, nil
}
