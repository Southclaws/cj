package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

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

func (cm *CommandManager) commandReVerify(
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

	verified, err = cm.Storage.IsUserLegacyVerified(message.Author.ID)
	if err != nil {
		return
	}

	if !verified {
		cm.Discord.ChannelMessageSend(message.ChannelID, "You are not legacy verified, please use 'verify' instead.")
		return
	}

	switch message.Content {
	case "reverify":
		err = cm.LegacyUserStartsVerification(message)
		context = true

	case "done":
		err = cm.LegacyUserConfirmsProfile(message)
		context = false

	case "cancel":
		err = cm.LegacyUserCancelsReVerification(message)
		context = false

	default:
		err = cm.UserProvidesProfileURL(message)
		context = true
	}

	return
}

// LegacyUserStartsVerification is called when the user sends the string "reverify" to
// the bot.
func (cm *CommandManager) LegacyUserStartsVerification(message discordgo.Message) (err error) {
	result, found := cm.Cache.Get(message.Author.ID)
	// At this point, it should not be found because this is the point where
	// a user should be starting their verification and thus there should be
	// no trace of their Verification in the cache.
	if found {
		verification, ok := result.(types.Verification)
		if !ok {
			return errors.New("failed to cast result to re-verification")
		}
		err = cm.WarnLegacyUserReVerificationState(message.ChannelID, verification)
		return
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
		`Hi! This process will verify you are the owner of a burgershot forum account. Please provide your user profile URL or ID.

For example: (Note: These are ***EXAMPLES***, don't just copy paste these...)

• %s
• %s
• %s

Each stage of the verification process will time-out after 5 minutes,
if you take longer than that to respond you will need to start again.`,
		`https://burgershot.gg/member.php?action=profile&uid=3`,
		`burgershot.gg/member.php?action=profile&uid=3`,
		`3`,
	))

	if err != nil {
		return
	}

	v := types.Verification{DiscordUser: *message.Author, ChannelID: message.ChannelID}
	cm.SetVerificationState(&v, types.VerificationStateAwaitProfileURL)

	return nil
}

// LegacyUserConfirmsProfile is called when the user responds with 'done'
func (cm *CommandManager) LegacyUserConfirmsProfile(message discordgo.Message) (err error) {
	var verification types.Verification

	result, found := cm.Cache.Get(message.Author.ID)
	if !found {
		err = cm.WarnUserNoVerification(message.ChannelID)
		return
	}

	verification, ok := result.(types.Verification)
	if !ok {
		return errors.New("failed to cast result to re-verification")
	}

	if verification.VerifyState != types.VerificationStateAwaitConfirmation {
		err = cm.WarnUserVerificationState(message.ChannelID, verification)
		return
	}

	verification.UserProfile, err = cm.Forum.GetUserProfilePage(verification.ForumUser)
	if err != nil {
		return
	}

	verified, err := cm.CheckUserPageForDiscordID(verification.UserProfile, message.Author.ID)
	if err != nil {
		return
	}

	if !verified {
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"Sorry, your re-verification failed. Your Discord ID was not found on your profile page.")
		return
	}

	err = cm.Storage.SetLegacyUserToVerified(verification)
	if err != nil {
		return
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
		"Congratulations! You have been re-verified as the owner of the forum account %s. Have a nice day!",
		verification.ForumUser,
	))
	return err
}

// LegacyUserCancelsReVerification is called when the user responds with 'cancel'
func (cm *CommandManager) LegacyUserCancelsReVerification(message discordgo.Message) (err error) {
	_, found := cm.Cache.Get(message.Author.ID)
	if !found {
		err = cm.WarnUserNoVerification(message.ChannelID)
		return
	}

	cm.Cache.Delete(message.Author.ID)

	cm.Discord.ChannelMessageSend(message.ChannelID,
		"You have cancelled your re-verification. You can start again at any time by sending 'verify'.")
	return
}

// WarnLegacyUserReVerificationState tells a user their current re-verification state,
// to be used when the user's reply does not match the expected reply according
// to the state of the re-verification associated with the user.
func (cm *CommandManager) WarnLegacyUserReVerificationState(channelid string, verification types.Verification) (err error) {
	var stateMessage string
	switch verification.VerifyState {
	case types.VerificationStateNone:
		stateMessage = "Your re-verification is currently in an invalid state, please try again in 5 minutes!"
	case types.VerificationStateAwaitProfileURL:
		stateMessage = "Your re-verification is currently awaiting a profile URL or profile ID."
	case types.VerificationStateAwaitConfirmation:
		stateMessage = "Your re-verification is currently awaiting you to add your discord ID to your profile, once you've done that reply with either 'done' or 'cancel'"
	}
	cm.Discord.ChannelMessageSend(channelid, stateMessage)
	return
}

// WarnLegacyUserNoVerification is simply a message informing the user their
// Re-verification does not exist and they need to start the process with 'reverify'.
func (cm *CommandManager) WarnLegacyUserNoVerification(channelid string) (err error) {
	cm.Discord.ChannelMessageSend(channelid, "You need to start your re-verification by typing 're-verify'.")
	return
}
