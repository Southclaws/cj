package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandTopRep(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	reaction := ""
	val, exists := args["reaction"]
	if exists {
		reaction = val.StringValue()
	}
	top, err := cm.Storage.GetTopReactions(10, reaction)
	if err != nil {
		cm.replyDirectly(interaction, err.Error())
		return
	}
	rankings, err := FormatReactionRankings(top, cm.Discord)
	if err != nil {
		return
	}

	cm.replyDirectlyEmbed(interaction, "", rankings)
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
