package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

// User is a recorded and verified SA:MP forum user.
type User struct {
	gorm.Model
	DiscordUserID    string `gorm:"primary_key;index"`
	ForumUserID      string
	VerificationCode string
}

func (app App) connectDB() {
	db, err := gorm.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})

	app.db = db
	log.Print("Connected to database")
}

// StoreVerifiedUser is for when a user finishes their verification.
func (app App) StoreVerifiedUser(verification Verification) {
	app.db.Create(&User{
		DiscordUserID:    verification.discordUser.ID,
		ForumUserID:      verification.forumUser,
		VerificationCode: verification.code,
	})
}

// GetDiscordUserForumAccount returns a discord user, a blank string or an error
func (app App) GetDiscordUserForumAccount(forumUserID string) (string, error) {
	var user User

	err := app.db.First(user, &User{ForumUserID: forumUserID}).Error
	if err != nil {
		return "", err
	}

	return user.DiscordUserID, nil
}
