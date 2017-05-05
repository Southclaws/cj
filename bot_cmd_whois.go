package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func commandWhois(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	var (
		verified bool
		err      error
		count    = 0
		username string
		link     string
		result   string
	)

	for _, user := range message.Mentions {
		if count == 5 {
			break
		}
		count++

		verified, err = cm.App.IsUserVerified(user.ID)
		if err != nil {
			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, err.Error())
			if err != nil {
				log.Print(err)
			}
			continue
		}

		if !verified {
			result += cm.App.locale.GetLangString("en", "CommandWhoisNotVerified", user.ID) + " "
		} else {
			username, err = cm.App.GetForumNameFromDiscordUser(user.ID)
			if err != nil {
				log.Print(err)
			}

			link, err = cm.App.GetForumUserFromDiscordUser(user.ID)
			if err != nil {
				log.Print(err)
			}

			result += cm.App.locale.GetLangString("en", "CommandWhoisProfile", user.ID, username, link) + " "
		}
	}

	_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		log.Print(err)
	}

	return true, false, nil
}
