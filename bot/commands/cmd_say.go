package commands

import (
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandSay(
	interaction *discordgo.InteractionCreate,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: interaction.Data.Options[0].StringValue(),
		},
	})
	return
}
