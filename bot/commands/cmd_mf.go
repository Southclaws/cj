package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandMessageFreq(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	freq, err := cm.Discord.GetCurrentChannelMessageFrequency(interaction.ChannelID)
	if err != nil {
		cm.replyDirectly(interaction, err.Error())
		return
	}
	cm.replyDirectly(interaction, fmt.Sprintf("%02f messages per second", freq))
	return
}
