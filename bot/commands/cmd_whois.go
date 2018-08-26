package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandWhois(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	var (
		verified bool
		count    = 0
		username string
		link     string
		result   string
	)

	if len(message.Mentions) == 0 {
		var userID string
		userID, err = cm.Storage.GetDiscordUserFromForumName(args)
		if err != nil {
			return
		}

		result += fmt.Sprintf("**%s** is here as <@%s>", args, userID)
	} else {
		for _, user := range message.Mentions {
			if count == 5 {
				break
			}
			count++

			if user.ID == cm.Config.BotID {
				result += "I am Carl Johnson, co-leader of Grove Street Families. "
				continue
			}

			verified, err = cm.Storage.IsUserVerified(user.ID)
			if err != nil {
				result += err.Error()
				continue
			}

			if !verified {
				result += fmt.Sprintf("The user <@%s> is not verified. ", user.ID)
			} else {
				username, err = cm.Storage.GetForumNameFromDiscordUser(user.ID)
				if err != nil {
					return
				}

				link, err = cm.Storage.GetForumUserFromDiscordUser(user.ID)
				if err != nil {
					return
				}

				result += fmt.Sprintf("<@%s> is **%s** (%s) on SA-MP forums. ", user.ID, username, link)
			}
		}
	}

	_, err = cm.Discord.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		return
	}

	return false, nil
}
