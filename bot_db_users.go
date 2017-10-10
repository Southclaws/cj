package main

import (
	"go.uber.org/zap"
)

// User is a recorded and verified SA:MP forum user.
type User struct {
	DiscordUserID    string `json:"discord_user_id"`
	ForumUserID      string `json:"forum_user_id"`
	VerificationCode string `json:"verification_code"`
	ForumUserName    string `json:"forum_user_name"`
}

// StoreVerifiedUser is for when a user finishes their verification.
func (app App) StoreVerifiedUser(verification Verification) (err error) {
	logger.Debug("storing verified user",
		zap.String("discord_id", verification.discordUser.ID),
		zap.String("forum_user", verification.forumUser),
		zap.String("code", verification.code),
		zap.String("forum_name", verification.userProfile.UserName))

	// err := app.db.Insert(&User{
	// 	DiscordUserID:    verification.discordUser.ID,
	// 	ForumUserID:      verification.forumUser,
	// 	VerificationCode: verification.code,
	// 	ForumUserName:    verification.userProfile.UserName,
	// }).Error

	return
}

// IsUserVerified returns a discord user, a blank string or an error
func (app App) IsUserVerified(discordUserID string) (verified bool, err error) {
	// var count int
	// result := app.db.Model(&User{}).Where(&User{DiscordUserID: discordUserID}).Count(&count)

	// if result.Error != nil {
	// 	err = errors.Wrap(result.Error, "failed to check user existence by discord ID")
	// } else if count == 0 {
	// 	verified = false
	// }

	return
}

// GetDiscordUserForumUser returns a discord user, a blank string or an error
func (app App) GetDiscordUserForumUser(forumUserID string) (string, error) {
	var user User

	// err := app.db.First(&user, &User{ForumUserID: forumUserID}).Error
	// if err != nil {
	// 	return "", errors.Wrap(err, "failed to query user by forum ID")
	// }

	return user.DiscordUserID, nil
}

// GetForumUserFromDiscordUser returns a link to user's profile, a blank string or an error
func (app App) GetForumUserFromDiscordUser(discordUserID string) (string, error) {
	var user User

	// err := app.db.First(&user, &User{DiscordUserID: discordUserID}).Error
	// if err != nil {
	// 	return "", errors.Wrap(err, "failed to query forum ID by discord ID")
	// }

	return user.ForumUserID, nil
}

// GetForumNameFromDiscordUser returns user's name on SA-MP Forums, a blank string or an error
func (app App) GetForumNameFromDiscordUser(discordUserID string) (string, error) {
	var user User

	// err := app.db.First(&user, &User{DiscordUserID: discordUserID}).Error
	// if err != nil {
	// 	return "", errors.Wrap(err, "failed to query forum name by discord ID")
	// }

	return user.ForumUserName, nil
}
