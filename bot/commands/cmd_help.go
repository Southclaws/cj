package commands

import (
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandHelp(
	args string,
	message discordgo.Message,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.Discord.ChannelMessageSend(message.ChannelID, "fuck off")
	return
}
