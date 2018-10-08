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
func FormatMessageRankings(r storage.TopMessages, s *discordgo.Session) (embed *discordgo.MessageEmbed, err error) {
	statsMessage := strings.Builder{}
	statsMessage.WriteString("Statistics\n\n") //nolint:errcheck

	embed = &discordgo.MessageEmbed{Color: 0x3498DB}
	var user *discordgo.User
	for i, tm := range r {
		user, err = s.User(tm.User)
		if err != nil {
			return
		}

		statsMessage.WriteString(fmt.Sprintf("%d. **<%s>** - %d\n", i, user.Username, tm.Messages)) //nolint:errcheck
	}

	embed.Description = statsMessage.String()

	return embed, nil
}
