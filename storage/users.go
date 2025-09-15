package storage

import (
	"regexp"
	// "time" // ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// User is a recorded and verified burgershot forum user.
type User struct {
	DiscordUserID     string            `json:"discord_user_id" bson:"discord_user_id"`
	ForumUserID       string            `json:"forum_user_id" bson:"forum_user_id"`
	BurgerUserID      string            `json:"burger_user_id" bson:"burger_user_id"`
	ForumUserName     string            `json:"forum_user_name" bson:"forum_user_name"`
	BurgerUserName    string            `json:"burger_user_name" bson:"burger_user_name"`
	BurgerVerify      bool              `json:"burgershot_verified" bson:"burgershot_verified"`
	ReceivedReactions []ReactionCounter `json:"received_reactions" bson:"received_reactions,omitempty"`
	/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
	TrackedRoles      []TrackedRole     `json:"tracked_roles" bson:"tracked_roles,omitempty"`
	*/
}

type ReactionCounter struct {
	Counter  int
	Reaction string
}

/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
// TrackedRole tracks the user roles
type TrackedRole struct {
	RoleID   string `json:"role_id" bson:"role_id"`
	RoleName string `json:"role_name" bson:"role_name"`
	AddedAt  int64  `json:"added_at" bson:"added_at"`
}
*/

/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
// AddTrackedRole adds a tracked role to a user
func (m *MongoStorer) AddTrackedRole(discordUserID, roleID, roleName string) error {
	user, err := m.GetUserOrCreate(discordUserID)
	if err != nil {
		return err
	}

	for _, role := range user.TrackedRoles {
		if role.RoleID == roleID {
			return nil
		}
	}

	newRole := TrackedRole{
		RoleID:   roleID,
		RoleName: roleName,
		AddedAt:  time.Now().Unix(),
	}
	user.TrackedRoles = append(user.TrackedRoles, newRole)
	return m.UpdateUser(user)
}
*/

/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
// RemoveTrackedRole removes a tracked role
func (m *MongoStorer) RemoveTrackedRole(discordUserID, roleID string) error {
	user, err := m.GetUserOrCreate(discordUserID)
	if err != nil {
		return err
	}

	for i, role := range user.TrackedRoles {
		if role.RoleID == roleID {
			user.TrackedRoles = append(user.TrackedRoles[:i], user.TrackedRoles[i+1:]...)
			break
		}
	}

	return m.UpdateUser(user)
}
*/

/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
// GetTrackedRoles returns all tracked roles for a user
func (m *MongoStorer) GetTrackedRoles(discordUserID string) ([]TrackedRole, error) {
	user, err := m.GetUserOrCreate(discordUserID)
	if err != nil {
		return nil, err
	}

	return user.TrackedRoles, nil
}
*/

// GetUserOrCreate gets a user or creates one and returns it
func (m *MongoStorer) GetUserOrCreate(discordUserID string) (user User, err error) {
	err = m.accounts.Find(bson.M{"discord_user_id": discordUserID}).One(&user)
	if err != nil {
		user.DiscordUserID = discordUserID
		err = m.accounts.Insert(&User{
			DiscordUserID: discordUserID,
			BurgerVerify:  false,
		})
		return user, err
	}
	return user, nil
}

// UpdateUser aims to update a full document of a user
func (m *MongoStorer) UpdateUser(user User) (err error) {
	err = m.accounts.Update(bson.M{"discord_user_id": user.DiscordUserID},
		bson.M{"$set": user})
	return err
}

// AddEmojiReactionToUser records an emoji reaction to a message of a discordUser.
func (m *MongoStorer) AddEmojiReactionToUser(discordUserID string, emoji string) (err error) {
	user, err := m.GetUserOrCreate(discordUserID)
	if err != nil {
		return err
	}
	var found = false
	for i, v := range user.ReceivedReactions {
		if v.Reaction == emoji {
			found = true
			user.ReceivedReactions[i].Counter++
		}
	}
	if !found {
		entry := ReactionCounter{
			Counter:  1,
			Reaction: emoji,
		}
		user.ReceivedReactions = append(user.ReceivedReactions, entry)
	}
	err = m.UpdateUser(user)
	return err
}

type TopReactionEntry struct {
	UserID   string `bson:"discord_user_id"`
	Counter  int    `bson:"counter"`
	Reaction string `bson:"reaction"`
}

// GetTopReactions gets the top <top> amount of people who received reaction <reaction>
func (m *MongoStorer) GetTopReactions(top int, reaction string) (result []TopReactionEntry, err error) {
	pipeline := []bson.M{
		{
			"$unwind": "$received_reactions",
		},
		{
			"$match": bson.M{
				"received_reactions.reaction": reaction,
			},
		},
		{
			"$project": bson.M{
				"discord_user_id": "$discord_user_id",
				"counter":         "$received_reactions.counter",
				"reaction":        "$received_reactions.reaction",
			},
		},
		{
			"$sort": bson.M{
				"counter": -1,
			},
		},
		{
			"$limit": top,
		},
	}
	// Get the top stats overall when no argument is supplied.
	// Just remove the $match pipeline stage.
	if reaction == "" {
		pipeline = append(pipeline[0:1], pipeline[2:]...)
	}
	m.accounts.Pipe(pipeline).All(&result)
	return
}

// RemoveEmojiReactionFromUser records an emoji reaction to a message of a discordUser.
func (m *MongoStorer) RemoveEmojiReactionFromUser(discordUserID string, emoji string) (err error) {
	user, err := m.GetUserOrCreate(discordUserID)
	if err != nil {
		return err
	}
	for i, v := range user.ReceivedReactions {
		if v.Reaction == emoji {
			user.ReceivedReactions[i].Counter--
		}
	}
	err = m.UpdateUser(user)
	return
}

// UpdateUserUsername updates a person's Burgershot forum name in the database. In case they have their name changed.
func (m *MongoStorer) UpdateUserUsername(discordUserID string, username string) (err error) {
	err = m.accounts.Update(
		bson.D{
			{Name: "discord_user_id", Value: discordUserID},
		},
		bson.D{
			{Name: "$set", Value: bson.D{
				{Name: "burger_user_name", Value: username},
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
			{Name: "discord_user_id", Value: discordUserID},
			{Name: "burgershot_verified", Value: bson.D{
				{Name: "$exists", Value: true},
			}},
		}).Count()
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
			{Name: "discord_user_id", Value: discordUserID},
			{Name: "burgershot_verified", Value: bson.D{
				{Name: "$exists", Value: false},
			}},
		}).Count()
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

	err = m.accounts.Find(bson.M{"burger_user_id": forumUserID}).One(&user)
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

// GetRandomUser returns a random Discord user ID
func (m *MongoStorer) GetRandomUser() (discordUserID string, err error) {
	var user User

	err = m.accounts.Pipe([]bson.M{
		{"$sample": bson.M{
			"size": 1,
		}},
	}).One(&user)
	if err != nil {
		err = errors.Wrap(err, "failed to get random user")
		return
	}

	discordUserID = user.DiscordUserID
	return
}
