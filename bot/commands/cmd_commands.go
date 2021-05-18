package commands

import (
	"fmt"

	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandCommands(
	args string,
	message discordgo.Message,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	embed := &discordgo.MessageEmbed{
		Color: 0x3498DB,
	}

	var cmdlist string
	for trigger, cmd := range cm.Commands {
		cmdlist += fmt.Sprintf("**%s** - %s\n", trigger, cmd.Description)
	}
	embed.Description = cmdlist

	_, err = cm.Discord.S.ChannelMessageSendEmbed(message.ChannelID, embed)
	return
}
