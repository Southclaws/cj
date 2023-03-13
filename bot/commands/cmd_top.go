package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/bot/heartbeat/stats"
	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandTop(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.sendThinkingResponse(interaction)
	top, err := cm.Storage.GetTopMessages(10)
	if err != nil {
		cm.editOriginalResponse(interaction, fmt.Sprintf("Failed to get message rankings: %s", err.Error()))
		return
	}

	rankings, err := stats.FormatMessageRankings(top, cm.Discord)
	if err != nil {
		cm.editOriginalResponse(interaction, fmt.Sprintf("Failed to format message rankings: %s", err.Error()))
		return
	}

	cm.editOriginalResponseWithEmbed(interaction, rankings)
	return
}
