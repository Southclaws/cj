package commands

import (
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandHelp(
	interaction *discordgo.InteractionCreate,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	err = cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "fuck off",
		},
	})
	return
}
