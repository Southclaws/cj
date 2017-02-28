package main

import (
	"strings"

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
		case "/userinfo":
			// Check if we have some parameter for the command, if not show usage message.
			if len(command) == 1 {
				app.discordClient.ChannelMessageSend(message.ChannelID, app.locale.GetLangString("en", "CommandUserInfoUsage"))
			} else {
				for _, user := range message.Mentions {
					verified, _ := app.IsUserVerified(user.ID)

					if verified == false {
						app.discordClient.ChannelMessageSend(message.ChannelID, app.locale.GetLangString("en", "UserNotVerified", user.ID))
					} else {
						link, _ := app.GetForumUserFromDiscordUser(user.ID)
						profile, _ := app.GetUserProfilePage(link)

						app.discordClient.ChannelMessageSend(message.ChannelID, app.locale.GetLangString("en", "CommandUserInfoProfile", profile.UserName, profile.JoinDate, profile.TotalPosts, profile.Reputation))
					}
				}
			}
		case "/whois":
			// Check if we have some parameter for the command, if not show usage message.
			if len(command) == 1 {
				app.discordClient.ChannelMessageSend(message.ChannelID, app.locale.GetLangString("en", "CommandWhoisUsage"))
			} else {
				for _, user := range message.Mentions {
					verified, _ := app.IsUserVerified(user.ID)

					if verified == false {
						app.discordClient.ChannelMessageSend(message.ChannelID, app.locale.GetLangString("en", "UserNotVerified", user.ID))
					} else {
						username, _ := app.GetForumNameFromDiscordUser(user.ID)
						link, _ := app.GetForumUserFromDiscordUser(user.ID)

						app.discordClient.ChannelMessageSend(message.ChannelID, app.locale.GetLangString("en", "CommandWhoisProfile", user.ID, username, link))
					}
				}
			}
		}
	}
	return nil
}
