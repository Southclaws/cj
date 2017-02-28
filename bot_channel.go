package main

import (
	"strings"

	"fmt"

	"github.com/bwmarrin/discordgo"
)

// HandleChannelMessage processes any public message from the channel that is
// designated as the primary channel for the bot.
// (see api.config.PrimaryChannel)
func (app App) HandleChannelMessage(message discordgo.Message) error {
	debug("[channel:HandleChannelMessage] %v", message)

	if message.Content[:1] == "/" {
		command := strings.Split(strings.ToLower(message.Content), " ")

		switch command[0] {
		case "/whois":
			// Check if we have some parameter for the command, if not show usage message.
			if len(command) == 1 {
				app.discordClient.ChannelMessageSend(message.ChannelID, "Command: **/whois [user(s)]**\nDescription: Get username on SA-MP Forums.\nExample: */whois @Southclaws#1657*")
			} else {
				for _, user := range message.Mentions {
					verified, _ := app.IsUserVerified(user.ID)

					if verified == false {
						app.discordClient.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> is not verified.", user.ID))
					} else {
						username, _ := app.GetForumNameFromDiscordUser(user.ID)
						app.discordClient.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> is **%s** on SA-MP forums.", user.ID, username))
					}
				}
			}
		}
	}
	return nil
}
