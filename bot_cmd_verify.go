package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandVerify(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	verified, err := cm.App.IsUserVerified(message.Author.ID)
	if err != nil {
		return false, false, err
	}

	if verified {
		_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, "You are already verified! If you are experiencing problems with the bot or verification, please contact an admin.")
		return true, false, err
	}

	switch message.Content {
	case "verify":
		err = cm.App.UserStartsVerification(message)
	case "done":
		err = cm.App.UserConfirmsProfile(message)
		return true, false, err
	case "cancel":
		err = cm.App.UserCancelsVerification(message)
	default:
		err = cm.App.UserProvidesProfileURL(message)
	}

	return true, true, err
}
