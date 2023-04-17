package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) ltf(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	channelID := "831189475480436746"

	if interaction.ChannelID != channelID {
		return
	}

	msg, err := cm.Storage.GetRandomMessageFromUsers([]string{
		"468356073095430144",
		"778144453751078913",
		"123456789987654321",
	})
	if err != nil {
		cm.replyDirectly(interaction, fmt.Sprintf(errors.Wrap(err, "failed to get messages for user").Error()))
		return
	}

	nick := "LinuxTheFish"
	time := time.Unix(msg.Timestamp, 0)

	cm.replyDirectly(interaction, fmt.Sprintf(
		"> %s\n - **%s** (%s, %d)",
		msg.Message, nick, time.Month().String(), time.Year()))

	return
}
