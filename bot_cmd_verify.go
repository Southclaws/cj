package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func commandVerify(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	verified, err := cm.App.IsUserVerified(message.Author.ID)
	if err != nil {
		return false, false, err
	}

	if verified {
		err = errors.New("user are already verified")
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
		return true, false, err
	default:
		err = cm.App.UserProvidesProfileURL(message)
	}

	return true, true, err
}
