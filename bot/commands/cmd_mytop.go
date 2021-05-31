package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandMyTop(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.sendThinkingResponse(interaction)
	rank, err := cm.Storage.GetUserRank(interaction.Member.User.ID)
	if err != nil {
		cm.editOriginalResponse(interaction, fmt.Sprintf(errors.Wrap(err, "failed to get user's rank").Error()))
		return
	}
	messageCount, err := cm.Storage.GetUserMessageCount(interaction.Member.User.ID)
	if err != nil {
		cm.editOriginalResponse(interaction, fmt.Sprintf(errors.Wrap(err, "failed to get user's message count").Error()))
		return
	}

	cm.editOriginalResponse(interaction, fmt.Sprintf("Your top: %d Messages: %d", rank, messageCount))
	return
}
