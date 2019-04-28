package storage

import (
	"regexp"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

// User is a recorded and verified burgershot forum user.
type User struct {
	DiscordUserID  string `json:"discord_user_id" bson:"discord_user_id"`
	ForumUserID    string `json:"forum_user_id" bson:"forum_user_id"`
	BurgerUserID   string `json:"burger_user_id" bson:"burger_user_id"`
	ForumUserName  string `json:"forum_user_name" bson:"forum_user_name"`
	BurgerUserName string `json:"burger_user_name" bson:"burger_user_name"`
	BurgerVerify   bool   `json:"burgershot_verified" bson:"burgershot_verified"`
}

// StoreVerifiedUser is for when a user finishes their verification.
func (m *MongoStorer) StoreVerifiedUser(verification types.Verification) (err error) {
	err = m.accounts.Insert(&User{
		DiscordUserID: verification.DiscordUser.ID,
		ForumUserID:   verification.ForumUser,
		ForumUserName: verification.UserProfile.UserName,
		BurgerVerify:  true,
	})

	return
}

// SetLegacyUserToVerified is for re-verification when a legacy verified user gets set to verified.
func (m *MongoStorer) SetLegacyUserToVerified(verification types.Verification) (err error) {
	err = m.accounts.Update(
		bson.D{
			{"discord_user_id", verification.DiscordUser.ID},
		},
		bson.D{
			{"$set", bson.D{
				{"burger_user_id", verification.ForumUser},
				{"burger_user_name", verification.UserProfile.UserName},
				{"burgershot_verified", true},
			}},
		})

	return
}

// RemoveUser removes a user by their Discord ID
func (m *MongoStorer) RemoveUser(id string) (err error) {
	return m.accounts.Remove(bson.M{"discord_user_id": id})
}

// IsUserVerified returns a discord user, a blank string or an error
// Difference between IsUserLegacyVerified: this specifically checks if the user verified on burgershot.
func (m *MongoStorer) IsUserVerified(discordUserID string) (verified bool, err error) {
	count, err := m.accounts.Find(
		bson.D{
			{"discord_user_id", discordUserID},
			{"burgershot_verified", bson.D{
				{"$exists", true},
			}}}).Count()
	if err != nil {
		return
	}
	if count > 0 {
		verified = true
	}
	return
}

// IsUserLegacyVerified returns a discord user, a blank string or an error
// Difference between IsUserVerified: this specifically checks if the user verified on SA:MP.
func (m *MongoStorer) IsUserLegacyVerified(discordUserID string) (verified bool, err error) {
	count, err := m.accounts.Find(
		bson.D{
			{"discord_user_id", discordUserID},
			{"burgershot_verified", bson.D{
				{"$exists", false},
			}}}).Count()
	if err != nil {

		return
	}
	if count > 0 {
		verified = true
	}
	return
}

// GetDiscordUserForumUser returns a discord user, a blank string or an error
func (m *MongoStorer) GetDiscordUserForumUser(forumUserID string) (discordUserID string, err error) {
	var user User

	err = m.accounts.Find(bson.M{"forum_user_id": forumUserID}).One(&user)
	if err != nil {
		err = errors.Wrap(err, "failed to query user by forum ID")
		return
	}

	discordUserID = user.DiscordUserID
	return
}

// GetForumUserFromDiscordUser returns a link to user's profile, a blank string or an error
func (m *MongoStorer) GetForumUserFromDiscordUser(discordUserID string) (legacyUserID string, burgerUserID string, err error) {
	var user User

	err = m.accounts.Find(bson.M{"discord_user_id": discordUserID}).One(&user)
	if err != nil {
		err = errors.Wrap(err, "failed to query forum ID by discord ID")
		return
	}

	legacyUserID = user.ForumUserID
	burgerUserID = user.BurgerUserID
	return
}

// GetForumNameFromDiscordUser returns user's name on SA-MP Forums & Burgershot, a blank string or an error
func (m *MongoStorer) GetForumNameFromDiscordUser(discordUserID string) (legacyUserName string, burgerUserName string, err error) {
	var user User

	err = m.accounts.Find(bson.M{"discord_user_id": discordUserID}).One(&user)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to query forum name by discord ID")
	}

	legacyUserName = user.ForumUserName
	burgerUserName = user.BurgerUserName
	return
}

// GetDiscordUserFromForumName returns user's name on SA-MP Forums, a blank string or an error
func (m *MongoStorer) GetDiscordUserFromForumName(forumName string) (legacyUserID string, burgerUserID string, err error) {
	var legacyUser User
	var burgerUser User

	regex := bson.M{"$regex": bson.RegEx{Pattern: "^" + regexp.QuoteMeta(forumName) + "$", Options: "i"}}

	erro := m.accounts.Find(bson.M{"forum_user_name": regex}).One(&legacyUser)
	if erro != nil {
		legacyUserID = ""
	} else {
		legacyUserID = legacyUser.DiscordUserID
	}

	erro = m.accounts.Find(bson.M{"burger_user_name": regex}).One(&burgerUser)
	if erro != nil {
		burgerUserID = ""
	} else {
		burgerUserID = burgerUser.DiscordUserID
	}

	if len(burgerUserID) == 0 && len(legacyUserID) == 0 {
		err = errors.New("user not found")
	}

	return
}
