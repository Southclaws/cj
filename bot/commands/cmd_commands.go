package commands

import (
	"fmt"

	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandCommands(
	interaction *discordgo.InteractionCreate,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	embed := &discordgo.MessageEmbed{
		Color: 0x3498DB,
	}

	var cmdlist string
	for _, cmd := range cm.Commands {
		cmdlist += fmt.Sprintf("**%s** - %s\n", cmd.Name, cmd.Description)
	}
	embed.Description = cmdlist

	err = cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	return
}
