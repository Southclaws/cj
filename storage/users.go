package storage

import (
	"regexp"

	"github.com/Southclaws/cj/types"
	"github.com/pkg/errors"
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
func (api *API) StoreVerifiedUser(verification types.Verification) (err error) {
	err = api.accounts.Insert(&User{
		DiscordUserID:    verification.DiscordUser.ID,
		ForumUserID:      verification.ForumUser,
		VerificationCode: verification.Code,
		ForumUserName:    verification.UserProfile.UserName,
	})

	return
}

// RemoveUser removes a user by their Discord ID
func (api *API) RemoveUser(id string) (err error) {
	return api.accounts.Remove(bson.M{"discord_user_id": id})
}

// IsUserVerified returns a discord user, a blank string or an error
func (api *API) IsUserVerified(discordUserID string) (verified bool, err error) {
	count, err := api.accounts.Find(bson.M{"discord_user_id": discordUserID}).Count()
	if err != nil {
		return
	}
	if count > 0 {
		verified = true
	}
	return
}

// GetDiscordUserForumUser returns a discord user, a blank string or an error
func (api *API) GetDiscordUserForumUser(forumUserID string) (discordUserID string, err error) {
	var user User

	err = api.accounts.Find(bson.M{"forum_user_id": forumUserID}).One(&user)
	if err != nil {
		err = errors.Wrap(err, "failed to query user by forum ID")
		return
	}

	discordUserID = user.DiscordUserID
	return
}

// GetForumUserFromDiscordUser returns a link to user's profile, a blank string or an error
func (api *API) GetForumUserFromDiscordUser(discordUserID string) (forumUserID string, err error) {
	var user User

	err = api.accounts.Find(bson.M{"discord_user_id": discordUserID}).One(&user)
	if err != nil {
		err = errors.Wrap(err, "failed to query forum ID by discord ID")
		return
	}

	forumUserID = user.ForumUserID
	return
}

// GetForumNameFromDiscordUser returns user's name on SA-MP Forums, a blank string or an error
func (api *API) GetForumNameFromDiscordUser(discordUserID string) (forumUserName string, err error) {
	var user User

	err = api.accounts.Find(bson.M{"discord_user_id": discordUserID}).One(&user)
	if err != nil {
		return "", errors.Wrap(err, "failed to query forum name by discord ID")
	}

	forumUserName = user.ForumUserName
	return
}

// GetDiscordUserFromForumName returns user's name on SA-MP Forums, a blank string or an error
func (api *API) GetDiscordUserFromForumName(forumName string) (discordUserID string, err error) {
	var user User

	regex := bson.M{"$regex": bson.RegEx{Pattern: "^" + regexp.QuoteMeta(forumName) + "$", Options: "i"}}

	err = api.accounts.Find(bson.M{"forum_user_name": regex}).One(&user)
	if err != nil {
		return "", errors.Wrap(err, "failed to query user by forum name")
	}

	discordUserID = user.DiscordUserID
	return
}
