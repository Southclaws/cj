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
	if numMessages < 10 {
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

	generated, err := chain.Generate(initial)
	if err != nil {
		return false, errors.Wrap(err, "failed to impersonate")
	}

	//nolint:errcheck
	cm.Discord.ChannelMessageSend(message.ChannelID, generated)

	return
}
