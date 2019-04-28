package commands

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/types"
)

/*
Verification process:
- > User issues command "verify" to the bot via direct message
- < Bot informs the user of the verification process and queries the user for
    their user profile page URL
- > User provides user profile page URL, with or without the protocol qualifier
    or just the user ID
- < Bot generates a unique token and provides it to the user, asks user to post
    it on their user Bio or Visitor Messages section and reply to the bot with
    either "done" or "cancel", any other string repeats the previous message
- > User posts the unique onto their Bio or Visitor Messages section and writes
    "done" back to the bot
- < Bot checks the Bio and Visitor Messages sections of the user profile page
    and verifies that the unique code provided to the user is present on the
    page then stores the information to the database and purges the Verification
    object from the local cache
*/

func (cm *CommandManager) commandVerify(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	verified, err := cm.Storage.IsUserVerified(message.Author.ID)
	if err != nil {
		return
	}

	if verified {
		err = errors.New("user already verified")
		return
	}

	switch message.Content {
	case "verify":
		err = cm.UserStartsVerification(message)
		context = true

	case "done":
		err = cm.UserConfirmsProfile(message)
		context = false

	case "cancel":
		err = cm.UserCancelsVerification(message)
		context = false

	default:
		err = cm.UserProvidesProfileURL(message)
		context = true
	}

	return
}

// SetVerificationState updates the state of a Verification and ensures it's
// cache entry is updated.
func (cm *CommandManager) SetVerificationState(v *types.Verification, state types.VerificationState) {
	v.VerifyState = state
	cm.Cache.Set(v.DiscordUser.ID, *v, gocache.DefaultExpiration)
}

// UserStartsVerification is called when the user sends the string "verify" to
// the bot.
func (cm *CommandManager) UserStartsVerification(message discordgo.Message) (err error) {
	result, found := cm.Cache.Get(message.Author.ID)
	// At this point, it should not be found because this is the point where
	// a user should be starting their verification and thus there should be
	// no trace of their Verification in the cache.
	if found {
		verification, ok := result.(types.Verification)
		if !ok {
			return errors.New("failed to cast result to verification")
		}
		err = cm.WarnUserVerificationState(message.ChannelID, verification)
		return
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
		`Hi! This process will verify you are the owner of a Burgershot forum account. Please provide your user profile URL or ID.

For example: (Note: These are ***EXAMPLES***, don't just copy paste these...)

• %s
• %s
• %s

Each stage of the verification process will time-out after 5 minutes,
if you take longer than that to respond you will need to start again.

Need some help regarding verification? You can always read the thread that explains how to verify! (%s)`,
		`https://www.burgershot.gg/member.php?action=profile&uid=3`,
		`www.burgershot.gg/member.php?action=profile&uid=3`,
		`3`,
		`https://www.burgershot.gg/showthread.php?tid=480`,
	))

	if err != nil {
		return
	}

	v := types.Verification{DiscordUser: *message.Author, ChannelID: message.ChannelID}
	cm.SetVerificationState(&v, types.VerificationStateAwaitProfileURL)

	return nil
}

// UserProvidesProfileURL is called when the user responds with a profile URL or
// profile ID.
func (cm *CommandManager) UserProvidesProfileURL(message discordgo.Message) (err error) {
	var verification types.Verification

	result, found := cm.Cache.Get(message.Author.ID)
	if !found {
		err = cm.WarnUserNoVerification(message.ChannelID)
		return
	}

	verification, ok := result.(types.Verification)
	if !ok {
		return errors.New("failed to cast result to verification")
	}

	if verification.VerifyState != types.VerificationStateAwaitProfileURL {
		err = cm.WarnUserVerificationState(message.ChannelID, verification)
		return
	}

	var profileURL string
	u, err := url.Parse(message.Content)
	if err == nil {
		if u.Scheme != "https" {
			u.Scheme = "https"
		}
		profileURL = u.String()
	} else {
		var value int
		value, err = strconv.Atoi(message.Content)
		if err != nil {
			err = cm.WarnUserVerificationState(message.ChannelID, verification)
			return
		}
		profileURL = strings.Trim(fmt.Sprintf("https://burgershot.gg/member.php?action=profile&uid=%d", value), " \n")
	}

	verification.ForumUser = profileURL
	if err != nil {
		return
	}

	cm.SetVerificationState(&verification, types.VerificationStateAwaitConfirmation)

	cm.Discord.ChannelMessageSend(message.ChannelID,
		fmt.Sprintf(`Thanks! Now you just need to paste this ID in the "Discord ID" section of your profile: **%s**.
		When you have done this, please reply with the message 'done'.`,
			message.Author.ID))
	return
}

// UserConfirmsProfile is called when the user responds with 'done'
func (cm *CommandManager) UserConfirmsProfile(message discordgo.Message) (err error) {
	var verification types.Verification

	result, found := cm.Cache.Get(message.Author.ID)
	if !found {
		err = cm.WarnUserNoVerification(message.ChannelID)
		return
	}

	verification, ok := result.(types.Verification)
	if !ok {
		return errors.New("failed to cast result to verification")
	}

	if verification.VerifyState != types.VerificationStateAwaitConfirmation {
		err = cm.WarnUserVerificationState(message.ChannelID, verification)
		return
	}

	verification.UserProfile, err = cm.Forum.GetUserProfilePage(verification.ForumUser)
	if err != nil {
		return
	}

	matched, err := cm.CheckUserPageForDiscordID(verification.UserProfile, message.Author.ID)
	if err != nil {
		return
	}

	if !matched {
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"Sorry, your verification failed. Your discord id was not found on your profile page.")
		return
	}

	err = cm.Storage.StoreVerifiedUser(verification)
	if err != nil {
		return
	}

	err = cm.Discord.S.GuildMemberRoleAdd(cm.Config.GuildID, verification.DiscordUser.ID, cm.Config.VerifiedRole)
	if err != nil {
		return
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
		"Congratulations! You have been verified as the owner of the forum account %s. Have a nice day!",
		verification.ForumUser,
	))
	return err
}

// UserCancelsVerification is called when the user responds with 'cancel'
func (cm *CommandManager) UserCancelsVerification(message discordgo.Message) (err error) {
	_, found := cm.Cache.Get(message.Author.ID)
	if !found {
		err = cm.WarnUserNoVerification(message.ChannelID)
		return
	}

	cm.Cache.Delete(message.Author.ID)

	cm.Discord.ChannelMessageSend(message.ChannelID,
		"You have cancelled your verification. You can start again at any time by sending 'verify'.")
	return
}

// WarnUserVerificationState tells a user their current verification state,
// to be used when the user's reply does not match the expected reply according
// to the state of the Verification associated with the user.
//nolint:lll
func (cm *CommandManager) WarnUserVerificationState(channelid string, verification types.Verification) (err error) {
	var stateMessage string
	switch verification.VerifyState {
	case types.VerificationStateNone:
		stateMessage = "Your verification is currently in an invalid state, please try again in 5 minutes!"
	case types.VerificationStateAwaitProfileURL:
		stateMessage = "Your verification is currently awaiting a profile URL or profile ID."
	case types.VerificationStateAwaitConfirmation:
		stateMessage = "Your verification is currently awaiting you to add your discord ID to your profile, once you've done that reply with either `done` or `cancel`"
	}
	cm.Discord.ChannelMessageSend(channelid, stateMessage)
	return
}

// WarnUserNoVerification is simply a message informing the user their
// Verification does not exist and they need to start the process with 'verify'.
func (cm *CommandManager) WarnUserNoVerification(channelid string) (err error) {
	cm.Discord.ChannelMessageSend(channelid, "You need to start your verification by typing 'verify'.")
	return
}

// WarnUserError informs the user of an error and provides them with
// instructions for what to do next.
func (cm *CommandManager) WarnUserError(channelid string, errorString string) (err error) {
	cm.Discord.ChannelMessageSend(channelid, fmt.Sprintf(`An error occurred: "%s"`, errorString))
	return
}

// CheckUserPageForDiscordID checks if a discord ID has been posted by a user.
func (cm *CommandManager) CheckUserPageForDiscordID(page forum.UserProfile, id string) (bool, error) {
	if strings.Contains(page.DiscordID, id) {
		return true, nil
	}
	return false, nil
}
