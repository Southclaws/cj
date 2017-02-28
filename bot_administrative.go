package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandleAdministrativeMessage processes any public message from the channel that is
// designated as the administrative channel for the bot.
// (see api.config.NewsChannel)
func (app App) HandleAdministrativeMessage(message discordgo.Message) error {
	debug("[channel:HandleAdministrativeMessage] %v", message)

	if message.Content[:1] == "/" {
		command := strings.Split(strings.ToLower(message.Content), " ")

		switch command[0] {
		case "/say":
			// Check if we have some parameter for the command, if not show usage message.
			if len(command) == 1 {
				app.discordClient.ChannelMessageSend(message.ChannelID, app.locale.GetLangString("en", "CommandSayUsage"))
			} else {
				// Remove the command from message (+ space) and send it to app.config.PrimaryChannel.
				app.discordClient.ChannelMessageSend(app.config.PrimaryChannel, message.Content[len(command[0]):len(message.Content)])
			}
		}
	}

	return nil
}
