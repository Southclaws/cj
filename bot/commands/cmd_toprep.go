package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandTopRep(
	args string,
	message discordgo.Message,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	top, err := cm.Storage.GetTopReactions(10, args)
	if err != nil {
		return false, errors.Wrap(err, "failed to get message rankings")
	}
	rankings, err := FormatReactionRankings(top, cm.Discord)
	if err != nil {
		return
	}

	_, err = cm.Discord.S.ChannelMessageSendEmbed(message.ChannelID, rankings)
	return
}

func FormatReactionRankings(r []storage.TopReactionEntry, s *discord.Session) (embed *discordgo.MessageEmbed, err error) {
	statsMessage := strings.Builder{}
	statsMessage.WriteString("Statistics\n\n") //nolint:errcheck

	embed = &discordgo.MessageEmbed{Color: 0x3498DB}
	var user *discordgo.User
	if len(r) == 0 {
		statsMessage.WriteString("There are no entries to display here!")
	}
	for i, tm := range r {
		var username string
		user, err = s.S.User(tm.UserID)
		if err != nil {
			zap.L().Warn("failed to get user", zap.Error(err), zap.String("user_id", tm.UserID))
			username = tm.UserID
		} else {
			username = user.Username
		}

		statsMessage.WriteString(fmt.Sprintf("%d. **<%s>** - %s x %d\n", i+1, username, tm.Reaction, tm.Counter)) //nolint:errcheck
	}

	embed.Description = statsMessage.String()

	return embed, nil
}
