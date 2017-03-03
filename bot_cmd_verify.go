package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandVerify(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, error) {
	debug("[cmd_verify:commandVerify] args: '%s'", args)

	verified, err := cm.App.IsUserVerified(message.Author.ID)

	if verified {
		_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, "You are already verified! If you are experiencing problems with the bot or verification, please contact an admin.")
		return true, err
	}
	switch args {
	case "verify":
		err = cm.App.UserStartsVerification(message)
	case "done":
		err = cm.App.UserConfirmsProfile(message)
	case "cancel":
		err = cm.App.UserCancelsVerification(message)
	default:
		err = cm.App.UserProvidesProfileURL(message)
	}

	return false, err
}
