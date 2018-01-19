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
		_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString(cm.App.config.Language, "CommandVerifyAlreadyVerified"),)
		return true, false, err
	}

	switch message.Content {
	case cm.App.locale.GetLangString(cm.App.config.Language, "CommandVerifyAlreadyVerify")
		err = cm.App.UserStartsVerification(message)
	case cm.App.locale.GetLangString(cm.App.config.Language, "CommandVerifyAlreadyDone")
		err = cm.App.UserConfirmsProfile(message)
		return true, false, err
	case cm.App.locale.GetLangString(cm.App.config.Language, "CommandVerifyAlreadyCancel")
		err = cm.App.UserCancelsVerification(message)
		return true, false, err
	default:
		err = cm.App.UserProvidesProfileURL(message)
	}

	return true, true, err
}
