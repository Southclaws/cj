package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

// User is a recorded and verified SA:MP forum user.
type User struct {
	gorm.Model
	DiscordUserID    string `gorm:"primary_key;index;not null;unique"`
	ForumUserID      string `gorm:"not null;unique"`
	VerificationCode string `gorm:"not null"`
	ForumUserName    string `gorm:"not null"`
}

// StoreVerifiedUser is for when a user finishes their verification.
func (app App) StoreVerifiedUser(verification Verification) error {
	debug("[users:StoreVerifiedUser] storing '%s' '%s' '%s' '%s'", verification.discordUser.ID, verification.forumUser, verification.code, verification.userProfile.UserName)
	err := app.db.Create(&User{
		DiscordUserID:    verification.discordUser.ID,
		ForumUserID:      verification.forumUser,
		VerificationCode: verification.code,
		ForumUserName:    verification.userProfile.UserName,
	}).Error

	return err
}

// IsUserVerified returns a discord user, a blank string or an error
func (app App) IsUserVerified(discordUserID string) (bool, error) {
	var count int
	result := app.db.Model(&User{}).Where(&User{DiscordUserID: discordUserID}).Count(&count)

	if result.Error != nil {
		return false, errors.Wrap(result.Error, "failed to check user existence by discord ID")
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
		return "", errors.Wrap(err, "failed to query user by forum ID")
	}

	return user.DiscordUserID, nil
}

// GetForumUserFromDiscordUser returns a link to user's profile, a blank string or an error
func (app App) GetForumUserFromDiscordUser(discordUserID string) (string, error) {
	var user User

	err := app.db.First(&user, &User{DiscordUserID: discordUserID}).Error
	if err != nil {
		return "", errors.Wrap(err, "failed to query forum ID by discord ID")
	}

	return user.ForumUserID, nil
}

// GetForumNameFromDiscordUser returns user's name on SA-MP Forums, a blank string or an error
func (app App) GetForumNameFromDiscordUser(discordUserID string) (string, error) {
	var user User

	err := app.db.First(&user, &User{DiscordUserID: discordUserID}).Error
	if err != nil {
		return "", errors.Wrap(err, "failed to query forum name by discord ID")
	}

	return user.ForumUserName, nil
}
