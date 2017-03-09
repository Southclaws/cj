package main

import (
	"github.com/bwmarrin/discordgo"
)

func commandWhois(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	var verified bool
	var err error
	var count = 0

	for _, user := range message.Mentions {
		if count == 5 {
			break
		}
		count++

		verified, err = cm.App.IsUserVerified(user.ID)
		if err != nil {
			cm.App.discordClient.ChannelMessageSend(message.ChannelID, err.Error())
			continue
		}

		if verified == false {
			cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString("en", "CommandWhoisNotVerified", user.ID))
		} else {
			username, _ := cm.App.GetForumNameFromDiscordUser(user.ID)
			link, _ := cm.App.GetForumUserFromDiscordUser(user.ID)

			cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString("en", "CommandWhoisProfile", user.ID, username, link))
		}
	}

	return true, false, nil
}
