package stats

import (
	"fmt"
	"strings"

	"github.com/Southclaws/cj/bot/heartbeat/common"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Aggregator collects statistics about messages and users
type Aggregator struct {
	Config  *types.Config
	Discord *discord.Session
	Storage storage.Storer
	Forum   *forum.ForumClient

	topMessages storage.TopMessages
}

//nolint:golint
func (a *Aggregator) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
	fc *forum.ForumClient,
) (name string, err error) {
	a.Config = config
	a.Storage = api
	a.Discord = discord
	a.Forum = fc
	return "aggregator", nil
}

//nolint:golint
func (a *Aggregator) Register() (actions []common.Action) {
	return []common.Action{
		{
			Schedule: "@hourly",
			Call:     a.gather,
		},
		{
			Schedule: "@every 7h30m",
			Call:     a.announce,
		},
	}
}

func (a *Aggregator) gather() (err error) {
	a.topMessages, err = a.Storage.GetTopMessages(10)
	return
}

func (a *Aggregator) announce() (err error) {
	rankings, err := FormatMessageRankings(a.topMessages, a.Discord)
	if err != nil {
		return
	}
	_, err = a.Discord.S.ChannelMessageSendEmbed(a.Config.PrimaryChannel, rankings)
	return
}

// FormatMessageRankings formats a TopMessages into a discord embed message
func FormatMessageRankings(r storage.TopMessages, s *discord.Session) (embed *discordgo.MessageEmbed, err error) {
	statsMessage := strings.Builder{}
	statsMessage.WriteString("Statistics\n\n") //nolint:errcheck

	embed = &discordgo.MessageEmbed{Color: 0x3498DB}
	var user *discordgo.User
	for i, tm := range r {
		var username string
		user, err = s.S.User(tm.User)
		if err != nil {
			zap.L().Warn("failed to get user", zap.Error(err), zap.String("user_id", tm.User))
			username = tm.User
		} else {
			username = user.Username
		}

		statsMessage.WriteString(fmt.Sprintf("%d. **<%s>** - %d\n", i+1, username, tm.Messages)) //nolint:errcheck
	}

	embed.Description = statsMessage.String()

	return embed, nil
}
