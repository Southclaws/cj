package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// User is a recorded and verified SA:MP forum user.
type User struct {
	DiscordUserID    string `gorm:"primary_key;index;not null;unique"`
	ForumUserID      string `gorm:"not null;unique"`
	VerificationCode string `gorm:"not null"`
}

func (app *App) connectDB() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open("sqlite3", filepath.Join(dir, "users.db"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
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
