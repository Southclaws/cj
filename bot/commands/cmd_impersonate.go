package commands

import (
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

	messages, err := cm.Storage.GetMessagesForUser(user.ID)
	if err != nil {
		return false, errors.Wrap(err, "failed to get messages for user")
	}
	numMessages := len(messages)
	if numMessages < 100 {
		return false, errors.New("not enough messages from that user")
	}

	chain := gomarkov.NewChain(1)
	for _, m := range messages {
		if isBadMessage(m.Message) {
			continue
		}
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
	cm.Discord.ChannelMessageSend(message.ChannelID, strings.Join(tokens[1:len(tokens)-1], " "))

	return
}

func isBadMessage(m string) bool {
	if strings.Contains(m, "@everyone") {
		return true
	}
	if strings.Contains(m, "@here") {
		return true
	}
	return false
}
