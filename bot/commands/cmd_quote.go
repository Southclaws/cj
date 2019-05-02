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
	)

	for ok == false {
		randmessage, messagerr = cm.Storage.GetRandomMessage()
		if messagerr != nil {
			cm.Discord.ChannelMessageSend(message.ChannelID, "Could not get a quote.")
			return
		}

		// Misses stuff like verify and doesn't allow commands to be sent either.
		if len(randmessage.Message) > 6 && strings.Index(randmessage.Message, "/") != 0 {
			ok = true
		}
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("\"%s\" ~ <@%s>", randmessage.Message, randmessage.DiscordUserID))
	return
}
