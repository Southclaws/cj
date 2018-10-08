package stats

import (
	"fmt"
	"strings"

	"github.com/Southclaws/cj/storage"
	"github.com/bwmarrin/discordgo"
)

func (a *Aggregator) gatherTopMessages(top int) (err error) {
	a.topMessages, err = a.Storage.GetTopMessages(top)
	return
}

// FormatMessageRankings formats a TopMessages into a discord embed message
func FormatMessageRankings(rankings storage.TopMessages, d *discordgo.Session) (embed *discordgo.MessageEmbed, err error) {
	statsMessage := strings.Builder{}
	statsMessage.WriteString("Statistics\n\n") //nolint:errcheck

	embed = &discordgo.MessageEmbed{Color: 0x3498DB}
	var user *discordgo.User
	for _, tm := range rankings {
		user, err = d.User(tm.User)
		if err != nil {
			return
		}
		statsMessage.WriteString(fmt.Sprintf("**<@%s>** - %d\n", user.Username, tm.Messages)) //nolint:errcheck
	}

	embed.Description = statsMessage.String()

	return
}
