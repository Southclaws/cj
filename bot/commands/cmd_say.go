package commands

import (
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandSay(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	_, err = cm.Discord.ChannelMessageSend(cm.Config.PrimaryChannel, args)
	return
}
