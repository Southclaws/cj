package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) ltf(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	channelID := "831189475480436746"

	if message.ChannelID != channelID {
		return
	}

	msg, err := cm.Storage.GetRandomMessageFromUsers([]string{
		"515214692751376387",
		"778144453751078913",
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to get messages for user")
	}

	nick := "LinuxTheFish"
	time := time.Unix(msg.Timestamp, 0)

	cm.Discord.ChannelMessageSend(channelID, fmt.Sprintf(
		"> %s\n - **%s** (%s, %d)",
		msg.Message, nick, time.Month().String(), time.Year()))

	return
}
