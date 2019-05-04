package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/bot/stats"
	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandTop(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	top, err := cm.Storage.GetTopMessages(10)
	if err != nil {
		return false, errors.Wrap(err, "failed to get message rankings")
	}

	rankings, err := stats.FormatMessageRankings(top, cm.Discord)
	if err != nil {
		return
	}

	_, err = cm.Discord.S.ChannelMessageSendEmbed(cm.Config.PrimaryChannel, rankings)
	return
}
