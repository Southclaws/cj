package talking

import (
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/cj/bot/heartbeat/common"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
	"github.com/mb-14/gomarkov"
	"github.com/pkg/errors"
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
) (err error) {
	t.Config = config
	t.Storage = api
	t.Discord = discord
	t.Forum = fc
	return
}

//nolint:golint
func (t *Talk) Register() (actions []common.Action) {
	return []common.Action{
		{
			Schedule: "1/10 * * * *",
			Chance:   0.5,
			Call:     t.quote,
		},
		{
			Schedule: "1/15 * * * *",
			Chance:   0.2,
			Call:     t.impersonate,
		},
	}
}

func (t *Talk) quote() (err error) {
	var (
		ok          bool
		randmessage storage.ChatLog
		restricted  bool
		nick        string
	)

	restrictedChannels := [11]string{
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
	}

	for ok == false {
		restricted = false
		randmessage, err = t.Storage.GetRandomMessage()
		if err != nil {
			return
		}

		// Misses stuff like verify and doesn't allow commands to be sent either.
		if len(randmessage.Message) > 6 && strings.Index(randmessage.Message, "/") != 0 {
			// TODO: Do through the query.
			for _, channel := range restrictedChannels {
				if randmessage.DiscordChannel == channel {
					restricted = true
				}
			}

			if !restricted {
				ok = true
			}
		}
	}

	user, err := t.Discord.S.GuildMember(t.Config.GuildID, randmessage.DiscordUserID)
	if err != nil || len(user.Nick) == 0 {
		// No longer on the server or has no nickname.
		var user *discordgo.User
		user, err = t.Discord.S.User(randmessage.DiscordUserID)
		if err != nil {
			return
		}

		nick = user.Username
	} else {
		nick = user.Nick
	}

	time := time.Unix(randmessage.Timestamp, 0)
	t.Discord.ChannelMessageSend(t.Config.PrimaryChannel, fmt.Sprintf(
		"\"%s\" ~ **%s**, **%s %d**",
		randmessage.Message, nick, time.Month().String(), time.Year()))
	return
}

func (t *Talk) impersonate() (err error) {
	m, err := t.Storage.GetRandomMessage()
	if err != nil {
		return
	}
	var userID = m.DiscordUserID

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

	var nick string
	user, err := t.Discord.S.GuildMember(t.Config.GuildID, userID)
	if err != nil || len(user.Nick) == 0 {
		// No longer on the server or has no nickname.
		var user *discordgo.User
		user, err = t.Discord.S.User(userID)
		if err != nil {
			return
		}

		nick = user.Username
	} else {
		nick = user.Nick
	}

	t.Discord.ChannelMessageSend(t.Config.PrimaryChannel, fmt.Sprintf(
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
