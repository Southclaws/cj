package stats

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/storage"
)

func (a *Aggregator) gatherTopMessages(top int) (err error) {
	a.topMessages, err = a.Storage.GetTopMessages(top)
	return
}

// FormatMessageRankings formats a TopMessages into a discord embed message
func FormatMessageRankings(rankings storage.TopMessages, s *storage.API) (embed *discordgo.MessageEmbed, err error) {
	statsMessage := strings.Builder{}
	statsMessage.WriteString("Statistics\n\n") //nolint:errcheck

	embed = &discordgo.MessageEmbed{Color: 0x3498DB}
	var user string
	for _, tm := range rankings {
		user, err = s.GetForumNameFromDiscordUser(tm.User)
		if err != nil {
			statsMessage.WriteString(fmt.Sprintf("**<@%s>** - %d (%s)\n", tm.User, err.Error())) //nolint:errcheck
		} else {
			statsMessage.WriteString(fmt.Sprintf("**<@%s>** - %d\n", user, tm.Messages)) //nolint:errcheck
		}
	}

	embed.Description = statsMessage.String()

	return embed, nil
}
