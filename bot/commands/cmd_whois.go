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
		verified       bool
		legacy         bool
		count          = 0
		username       string
		legacyusername string
		link           string
		legacylink     string
		result         string
	)

	if len(message.Mentions) == 0 {
		var legacyuserID string
		var burgeruserID string

		legacyuserID, burgeruserID, err = cm.Storage.GetDiscordUserFromForumName(args)
		if err != nil {
			return
		}

		if len(legacyuserID) > 0 {
			result += fmt.Sprintf("**%s** on the SA:MP forums is here as <@%s>\n", args, legacyuserID)
		}
		if len(burgeruserID) > 0 {
			result += fmt.Sprintf("**%s** on Burgershot is here as <@%s>", args, burgeruserID)
		}
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

			legacy, err = cm.Storage.IsUserLegacyVerified(user.ID)
			if err != nil {
				result += err.Error()
				continue
			}

			if !verified && !legacy {
				result += fmt.Sprintf("The user <@%s> is not verified. ", user.ID)
			} else {
				legacyusername, username, err = cm.Storage.GetForumNameFromDiscordUser(user.ID)
				if err != nil {
					return
				}

				legacylink, link, err = cm.Storage.GetForumUserFromDiscordUser(user.ID)
				if err != nil {
					return
				}

				if len(legacylink) > 0 && len(legacyusername) > 0 {
					result += fmt.Sprintf("<@%s> is **%s** (%s)\n", user.ID, legacyusername, legacylink)
				}
				if !legacy {
					result += fmt.Sprintf("<@%s> is **%s** (%s). ", user.ID, username, link)
				}
			}
		}
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		return
	}

	return false, nil
}
