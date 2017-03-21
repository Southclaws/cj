package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func commandWhois(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	var verified bool
	var err error

	if len(message.Mentions) == 0 {
		return false, false, nil
	}

	for _, user := range message.Mentions {
		if user.ID == cm.App.config.BotID {
			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString("en", "CommandWhoisCJ"))
			if err != nil {
				log.Print(err)
			}
			continue
		}

		verified, err = cm.App.IsUserVerified(user.ID)
		if err != nil {
			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, err.Error())
			if err != nil {
				log.Print(err)
			}
			continue
		}

		if !verified {
			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString("en", "UserNotVerified", user.ID))
			if err != nil {
				log.Print(err)
			}
		} else {
			username, err := cm.App.GetForumNameFromDiscordUser(user.ID)
			if err != nil {
				log.Print(err)
			}

			link, err := cm.App.GetForumUserFromDiscordUser(user.ID)
			if err != nil {
				log.Print(err)
			}

			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString("en", "CommandWhoisProfile", user.ID, username, link))
			if err != nil {
				log.Print(err)
			}
		}
	}

	return true, false, err
}
