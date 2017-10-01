package main

import (
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

	if len(message.Mentions) == 0 {
		return false, false, nil
	}

	for _, user := range message.Mentions {
		if count == 5 {
			break
		}
		count++

		if user.ID == cm.App.config.BotID {
			result += cm.App.locale.GetLangString("en", "CommandWhoisCJ") + " "
			continue
		}

		verified, err = cm.App.IsUserVerified(user.ID)
		if err != nil {
			result += err.Error() + " "
			continue
		}

		if !verified {
			result += cm.App.locale.GetLangString("en", "CommandWhoisNotVerified", user.ID) + " "
		} else {
			username, err = cm.App.GetForumNameFromDiscordUser(user.ID)
			if err != nil {
				return false, false, err
			}

			link, err = cm.App.GetForumUserFromDiscordUser(user.ID)
			if err != nil {
				return false, false, err
			}

			result += cm.App.locale.GetLangString("en", "CommandWhoisProfile", user.ID, username, link) + " "
		}
	}

	_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		return false, false, err
	}

	return true, false, nil
}
