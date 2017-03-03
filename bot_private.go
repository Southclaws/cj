package main

import (
	"log"

	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandlePrivateMessage processes a private message sent directly to the bot
// usually for direct commands such as account verification.
func (app *App) HandlePrivateMessage(message discordgo.Message) error {
	debug("[private:HandlePrivateMessage] message '%s'", message.Content)

	var err error

	// Convert incoming commands to lowercase.
	message.Content = strings.ToLower(message.Content)

	if err != nil {
		log.Printf("HandlePrivateMessage caused an error: %s", err)
		err = app.WarnUserError(message.ChannelID, err.Error())
		if err != nil {
			log.Printf("WarnUserError caused an error: %s", err.Error())
		}
	}

	return nil
}
