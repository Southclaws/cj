package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/bot/heartbeat/stats"
	"github.com/Southclaws/cj/types"
)


func (cm *CommandManager) commandTopRep(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	top, err := cm.Storage.GetTopReactions(10, args)
	if err != nil {
		return false, errors.Wrap(err, "failed to get message rankings")
	}
	rankings, err := stats.FormatReactionRankings(top, cm.Discord)
	if err != nil {
		return
	}

	_, err = cm.Discord.S.ChannelMessageSendEmbed(message.ChannelID, rankings)
	return
}
