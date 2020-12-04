package talking

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mb-14/gomarkov"
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
			Schedule: "@every 1m",
			Chance:   0.8,
			Call:     t.quote,
		},
		{
			Schedule: "@every 1m",
			Chance:   0.5,
			Call:     t.impersonate,
		},
	}
}

func (t *Talk) quote() (err error) {
	channelID, err := t.Discord.GetRandomChannel()
	if err != nil {
		return
	}

	f, err := t.Discord.GetCurrentChannelMessageFrequency(channelID)
	if err != nil {
		return
	}
	if f < 0.5 {
		zap.L().Debug("dropping random quote because channel not active enough")
		return
	}
	randmessage, err := t.Storage.GetRandomMessage()
	if err != nil {
		return
	}
	if len(randmessage.Message) < 6 {
		zap.L().Debug("dropping random quote because it's too short")
		return
	}
	if strings.Index(randmessage.Message, "/") == 0 {
		zap.L().Debug("dropping random quote because it's a command invocation")
		return
	}
	for _, channel := range [11]string{
		"457943077789892649",
		"536617253286707210",
		"376371371795546112",
		"361640153489473540",
		"531855907374759936",
		"531860517816500235",
		"572476607772491786",
		"282581181034135556",
		"282581078643048448",
		"415801201960026122",
		"548332120842698753",
	} {
		if randmessage.DiscordChannel == channel {
			zap.L().Debug("dropping random quote because it's from a restricted channel")
			return
		}
	}

	nick := t.getMostRecentNick(t.Config.GuildID, randmessage.DiscordUserID)
	time := time.Unix(randmessage.Timestamp, 0)

	t.Discord.ChannelMessageSend(channelID, fmt.Sprintf(
		"> %s\n - **%s** (%s, %d)",
		randmessage.Message, nick, time.Month().String(), time.Year()))
	return
}

func (t *Talk) impersonate() (err error) {
	channelID, err := t.Discord.GetRandomChannel()
	if err != nil {
		return
	}

	f, err := t.Discord.GetCurrentChannelMessageFrequency(channelID)
	if err != nil {
		return
	}
	if f < 0.5 {
		zap.L().Debug("dropping impersonation because channel not active enough")
		return
	}

	userID, err := t.Storage.GetRandomUser()
	if err != nil {
		return
	}

	messages, err := t.Storage.GetMessagesForUser(userID)
	if err != nil {
		return errors.Wrap(err, "failed to get messages for user")
	}

	if len(messages) < 10 {
		return
	}

	chain := gomarkov.NewChain(1)
	for _, m := range messages {
		if isBadMessage(m.Message) {
			continue
		}
		words := strings.Split(m.Message, " ")
		if len(words) < 3 {
			continue
		}
		chain.Add(words)
	}

	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		var next string
		next, err = chain.Generate(tokens[(len(tokens) - 1):])
		if err != nil {
			return errors.Wrap(err, "failed to impersonate")
		}
		tokens = append(tokens, next)
	}

	nick := t.getMostRecentNick(t.Config.GuildID, userID)

	t.Discord.ChannelMessageSend(channelID, fmt.Sprintf(
		"\"%s\" is what **%s** sounds like",
		strings.Join(tokens[1:len(tokens)-1], " "), nick))

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
