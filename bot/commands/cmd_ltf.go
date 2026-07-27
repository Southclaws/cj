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
	if interaction.ChannelID != cm.GuildSettings.LTFChannelID {
		return
	}

	msg, err := cm.Storage.GetRandomMessageFromUsers(cm.GuildSettings.LTFUserIDs)
	if err != nil {
		cm.replyDirectly(interaction, fmt.Sprint(errors.Wrap(err, "failed to get messages for user").Error()))
		return
	}

	nick := "LinuxTheFish"
	time := time.Unix(msg.Timestamp, 0)

	cm.replyDirectly(interaction, fmt.Sprintf(
		"> %s\n - **%s** (%s, %d)",
		msg.Message, nick, time.Month().String(), time.Year()))

	return
}
