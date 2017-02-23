package main

import "github.com/bwmarrin/discordgo"

// HandlePrivateMessage processes a private message sent directly to the bot
// usually for direct commands such as account verification.
func (app App) HandlePrivateMessage(message discordgo.Message) error {
	debug("[private:HandlePrivateMessage] message '%s'", message.Content)

	if app.config.Debug && app.config.DebugUser != "" {
		if message.Author.ID != app.config.DebugUser {
			return nil
		}
	}

	var err error

	if message.Content == "verify" {
		app.UserStartsVerification(message)
	} else if message.Content == "done" {
		app.UserConfirmsProfile(message)
	} else if message.Content == "cancel" {
		app.UserCancelsVerification(message)
	} else {
		app.UserProvidesProfileURL(message)
	}

	if err != nil {
		app.WarnUserError(message.ChannelID, err.Error())
	}

	return nil
}
