package main

import "github.com/bwmarrin/discordgo"

// HandlePrivateMessage processes a private message sent directly to the bot
// usually for direct commands such as account verification.
func (app *App) HandlePrivateMessage(message discordgo.Message) error {
	debug("[private:HandlePrivateMessage] message '%s'", message.Content)

	if app.config.DebugUser != "" {
		if message.Author.ID != app.config.DebugUser {
			debug("[private:HandlePrivateMessage] app.config.Debug true, user ID does not match app.config.DebugUser")
			return nil
		}
	}

	var err error

	verified, err := app.IsUserVerified(message.Author.ID)

	if verified {
		if message.Content == "kill" {
			if message.Author.ID == app.config.Admin {
				debug("[private:HandlePrivateMessage] kill signal received from '%s'", message.Author.ID)
				app.done <- true
			}
		}

		// todo: build command system for verified users and Discord admins.

	} else {
		if message.Content == "verify" {
			app.UserStartsVerification(message)
		} else if message.Content == "done" {
			app.UserConfirmsProfile(message)
		} else if message.Content == "cancel" {
			app.UserCancelsVerification(message)
		} else {
			app.UserProvidesProfileURL(message)
		}
	}

	if err != nil {
		app.WarnUserError(message.ChannelID, err.Error())
	}

	return nil
}
