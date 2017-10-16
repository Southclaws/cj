package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

// User is a recorded and verified SA:MP forum user.
type User struct {
	DiscordUserID    string `json:"discord_user_id" bson:"discord_user_id"`
	ForumUserID      string `json:"forum_user_id" bson:"forum_user_id"`
	VerificationCode string `json:"verification_code" bson:"verification_code"`
	ForumUserName    string `json:"forum_user_name" bson:"forum_user_name"`
}

// StoreVerifiedUser is for when a user finishes their verification.
func (app App) StoreVerifiedUser(verification Verification) (err error) {
	logger.Debug("storing verified user",
		zap.String("discord_id", verification.discordUser.ID),
		zap.String("forum_user", verification.forumUser),
		zap.String("code", verification.code),
		zap.String("forum_name", verification.userProfile.UserName))

	err = app.accounts.Insert(&User{
		DiscordUserID:    verification.discordUser.ID,
		ForumUserID:      verification.forumUser,
		VerificationCode: verification.code,
		ForumUserName:    verification.userProfile.UserName,
	})

	return
}

// IsUserVerified returns a discord user, a blank string or an error
func (app App) IsUserVerified(discordUserID string) (verified bool, err error) {
	count, err := app.accounts.Find(bson.M{"discord_user_id": discordUserID}).Count()
	if err != nil {
		return
	}
	if count > 0 {
		verified = true
	}
	return
}

// GetDiscordUserForumUser returns a discord user, a blank string or an error
func (app App) GetDiscordUserForumUser(forumUserID string) (discordUserID string, err error) {
	var user User

	err = app.accounts.Find(bson.M{"forum_user_id": forumUserID}).One(&user)
	if err != nil {
		err = errors.Wrap(err, "failed to query user by forum ID")
		return
	}

	discordUserID = user.DiscordUserID
	return
}

// GetForumUserFromDiscordUser returns a link to user's profile, a blank string or an error
func (app App) GetForumUserFromDiscordUser(discordUserID string) (forumUserID string, err error) {
	var user User

	err = app.accounts.Find(bson.M{"discord_user_id": discordUserID}).One(&user)
	if err != nil {
		err = errors.Wrap(err, "failed to query forum ID by discord ID")
		return
	}

	forumUserID = user.ForumUserID
	return
}

// GetForumNameFromDiscordUser returns user's name on SA-MP Forums, a blank string or an error
func (app App) GetForumNameFromDiscordUser(discordUserID string) (forumUserName string, err error) {
	var user User

	err = app.accounts.Find(bson.M{"discord_user_id": discordUserID}).One(&user)
	if err != nil {
		return "", errors.Wrap(err, "failed to query forum name by discord ID")
	}

	forumUserName = user.ForumUserName
	return
}
