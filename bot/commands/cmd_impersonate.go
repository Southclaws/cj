package commands

import (
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mb-14/gomarkov"
	"github.com/pkg/errors"
)

func (cm *CommandManager) commandImpersonate(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	if len(message.Mentions) != 1 {
		return false, errors.New("you must mention one person")
	}

	user := message.Mentions[0]

	chain := gomarkov.NewChain(2)
	messages, err := cm.Storage.GetMessagesForUser(user.ID)
	if err != nil {
		return false, errors.Wrap(err, "failed to get messages for user")
	}
	numMessages := len(messages)
	if numMessages < 100 {
		return false, errors.New("not enough messages from that user")
	}

	var initial []string
	for tries := 0; tries < 10; tries++ {
		m := messages[rand.Intn(numMessages)]
		words := strings.Split(m.Message, " ")
		if len(words) < 2 {
			continue
		}
		initial = append(initial, words[:2]...)
		break
	}
	if initial == nil {
		return false, errors.New("could not get initial word pair, not enough data")
	}

	for _, m := range messages {
		chain.Add(strings.Split(m.Message, " "))
	}

	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, err := chain.Generate(tokens[(len(tokens) - 1):])
		if err != nil {
			return false, errors.Wrap(err, "failed to impersonate")
		}
		tokens = append(tokens, next)
	}

	//nolint:errcheck
	cm.Discord.ChannelMessageSend(message.ChannelID, strings.Join(tokens, " "))

	return
}
