package commands

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
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
	case "done":
		err = cm.UserConfirmsProfile(message)
		return
	case "cancel":
		err = cm.UserCancelsVerification(message)
		return
	default:
		err = cm.UserProvidesProfileURL(message)
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

	_, err = cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
		`Hi! This process will verify you are the owner of a SA:MP forum account. Please provide your user profile URL or ID.

Examples:

- %s
- %s
- %s

Each stage of the verification process will time-out after 5 minutes,
if you take longer than that to respond you will need to start again.`,
		`http://forum.sa-mp.com/member.php?u=50199`,
		`forum.sa-mp.com/member.php?u=50199`,
		`50199`,
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

	matched, err := regexp.MatchString(`^(http:\/\/)?forum\.sa-mp\.com\/member\.php\?u=[0-9]*$`, message.Content)
	if err != nil {
		err = cm.WarnUserVerificationState(message.ChannelID, verification)
		return
	}

	var profileURL string

	// If it didn't match, check if it's just a user ID
	if matched {
		if strings.HasPrefix(message.Content, "http://") {
			profileURL = message.Content
		} else {
			profileURL = "http://" + message.Content
		}
	} else {
		var value int
		value, err = strconv.Atoi(message.Content)
		if err != nil {
			err = cm.WarnUserVerificationState(message.ChannelID, verification)
			return
		}
		profileURL = fmt.Sprintf("http://forum.sa-mp.com/member.php?u=%d", value)
	}

	verification.ForumUser = profileURL
	verification.Code, err = GenerateRandomString(8)
	if err != nil {
		return
	}

	cm.SetVerificationState(&verification, types.VerificationStateAwaitConfirmation)

	//nolint:lll
	_, err = cm.Discord.ChannelMessageSend(message.ChannelID,
		fmt.Sprintf(`Thanks! Now you just need to paste this 8-digit verification code into your Bio section then reply with 'done'.

**%s**

You can edit this section here: http://forum.sa-mp.com/profile.php?do=editprofile in the 'Additional Information' section at the bottom.

You must ensure your "About Me" section is visible to *Everyone*, you can change this setting here: http://forum.sa-mp.com/profile.php?do=privacy`,
			verification.Code))
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

	verified, err := cm.CheckUserPageForCode(verification.UserProfile, verification.Code)
	if err != nil {
		return
	}

	if !verified {
		_, err = cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"Sorry, your verification failed. The code was not found on your profile page.")
		return
	}

	err = cm.Storage.StoreVerifiedUser(verification)
	if err != nil {
		return
	}

	err = cm.Discord.GuildMemberRoleAdd(cm.Config.GuildID, verification.DiscordUser.ID, cm.Config.VerifiedRole)
	if err != nil {
		return
	}

	_, err = cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
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

	_, err = cm.Discord.ChannelMessageSend(message.ChannelID,
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
		stateMessage = "Your verification is currently awaiting you to post the verification code on your Profile Bio, once you've done that reply with either 'done' or 'cancel'"
	}
	_, err = cm.Discord.ChannelMessageSend(channelid, stateMessage)
	return
}

// WarnUserNoVerification is simply a message informing the user their
// Verification does not exist and they need to start the process with 'verify'.
func (cm *CommandManager) WarnUserNoVerification(channelid string) (err error) {
	_, err = cm.Discord.ChannelMessageSend(channelid, "You need to start your verification by typing 'verify'.")
	return
}

// WarnUserError informs the user of an error and provides them with
// instructions for what to do next.
func (cm *CommandManager) WarnUserError(channelid string, errorString string) (err error) {
	_, err = cm.Discord.ChannelMessageSend(channelid, fmt.Sprintf(`An error occurred: "%s"`, errorString))
	return
}

// CheckUserPageForCode checks if a verification code has been posted by a user.
func (cm *CommandManager) CheckUserPageForCode(page forum.UserProfile, code string) (bool, error) {
	if strings.Contains(page.BioText, code) {
		return true, nil
	}
	return false, nil
}

/*
Author: Matt Silverlock
Date: 2014-05-24
Accessed: 2017-02-22
https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand
*/

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b) //nolint:gas
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
