package talking

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/bot/heartbeat/common"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// Talk randomly sends messages
type Talk struct {
	Config  *types.Config
	Discord *discord.Session
	Storage storage.Storer
	Forum   *forum.ForumClient
}

//nolint:golint
func (t *Talk) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
	fc *forum.ForumClient,
) (name string, err error) {
	t.Config = config
	t.Storage = api
	t.Discord = discord
	t.Forum = fc
	return "talking", nil
}

//nolint:golint
func (t *Talk) Register() (actions []common.Action) {
	return []common.Action{
		{
			Schedule: "@every 1h",
			Chance:   0.5,
			Call:     t.ltf,
		},
	}
}

func (t *Talk) ltf() (err error) {
	channelID := "831189475480436746"

	f, err := t.Discord.GetCurrentChannelMessageFrequency(channelID)
	if err != nil {
		return
	}
	if f < 0.01 {
		return
	}

	message, err := t.Storage.GetRandomMessageFromUsers([]string{
		"456226577798135808",
		"778144453751078913",
	})
	if err != nil {
		return errors.Wrap(err, "failed to get messages for user")
	}

	nick := "LinuxTheFish"
	time := time.Unix(message.Timestamp, 0)

	t.Discord.ChannelMessageSend(channelID, fmt.Sprintf(
		"> %s\n - **%s** (%s, %d)",
		message.Message, nick, time.Month().String(), time.Year()))

	return
}

func isBadMessage(m string) bool {
	if strings.Contains(m, "<@") {
		return true
	}
	if strings.Contains(m, "@everyone") {
		return true
	}
	if strings.Contains(m, "@here") {
		return true
	}
	if strings.Contains(m, "http") {
		return true
	}
	return false
}

func (t *Talk) getMostRecentNick(guildID, userID string) (nick string) {
	user, err := t.Discord.S.GuildMember(guildID, userID)
	if err != nil || len(user.Nick) == 0 {
		// No longer on the server or has no nickname.
		var user *discordgo.User
		user, err = t.Discord.S.User(userID)
		if err != nil {
			zap.L().Warn("failed to get user name", zap.Error(err))
			nick = userID
		} else {
			nick = user.Username
		}
	} else {
		nick = user.Nick
	}
	return
}
