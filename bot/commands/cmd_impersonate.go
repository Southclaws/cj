package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-shellwords"
	"github.com/mb-14/gomarkov"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandImpersonate(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	// this actually allows arbitrary commands to be run by users!
	// but it runs in a scratch image with no commands... so, it should be fine.
	mentions, err := shellwords.Parse(message.ContentWithMentionsReplaced())
	if err != nil {
		return
	}

	if len(mentions) <= 1 || len(mentions) > 6 {
		return false, errors.New("requires 1-5 usernames")
	}

	err = func() (err error) {
		var messages []storage.ChatLog
		for i, username := range mentions[1:] {
			if username == "CJ" {
				for _, m := range quotes {
					messages = append(messages, storage.ChatLog{Message: m})
				}
			} else {
				user, ok := cm.Discord.GetUserFromName(username)
				if !ok {
					if i == 0 {
						continue
					} else {
						return errors.New("User not found")
					}
				}

				messages, err = cm.Storage.GetMessagesForUser(user.User.ID)
				if err != nil {
					return errors.Wrap(err, "failed to get messages for user")
				}
			}

			numMessages := len(messages)
			if numMessages < 10 {
				if i == 0 {
					continue
				} else {
					return errors.New("not enough messages from that user")
				}
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
				next, err := chain.Generate(tokens[(len(tokens) - 1):])
				if err != nil {
					return errors.Wrap(err, "failed to impersonate")
				}
				tokens = append(tokens, next)
			}

			cm.Discord.ChannelMessageSend(message.ChannelID, strings.Join(tokens[1:len(tokens)-1], " "))
		}

		return
	}()
	if err != nil {
		cm.Discord.ChannelMessageSend(message.ChannelID, err.Error())
		return false, nil
	}

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
