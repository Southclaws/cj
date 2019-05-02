package commands

import (
	"fmt"
	"strings"

	"github.com/Southclaws/cj/storage"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandQuote(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	var (
		ok          bool
		randmessage storage.ChatLog
		messagerr   error
		restricted  bool
	)

	restrictedChannels := [11]string{
		"457943077789892649",
		"536617253286707210",
		"376371371795546112",
		"361640153489473540",
		"531855907374759936",
		"531860517816500235",
		"572476607772491786",
	}

	for ok == false {
		restricted = false
		randmessage, messagerr = cm.Storage.GetRandomMessage()
		if messagerr != nil {
			cm.Discord.ChannelMessageSend(message.ChannelID, "Could not get a quote.")
			return
		}

		// Misses stuff like verify and doesn't allow commands to be sent either.
		if len(randmessage.Message) > 6 && strings.Index(randmessage.Message, "/") != 0 {
			for _, channel := range restrictedChannels {
				if randmessage.DiscordChannel == channel {
					restricted = true
				}
			}

			if !restricted {
				ok = true
			}
		}
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("\"%s\" ~ <@%s>", randmessage.Message, randmessage.DiscordUserID))
	return
}
