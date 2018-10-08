package commands

import (
	"github.com/Southclaws/cj/bot/stats"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func (cm *CommandManager) commandTop(
	args string,
	message discordgo.Message,
	contextual bool,
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

	_, err = cm.Discord.ChannelMessageSendEmbed(cm.Config.PrimaryChannel, rankings)
	return
}
