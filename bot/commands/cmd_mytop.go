package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandMyTop(
	args string,
	message discordgo.Message,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	rank, err := cm.Storage.GetUserRank(message.Author.ID)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user's rank")
	}
	messageCount, err := cm.Storage.GetUserMessageCount(message.Author.ID)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user's message count")
	}

	_, err = cm.Discord.S.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Your top: %d. Messages: %d", rank, messageCount))
	return
}
