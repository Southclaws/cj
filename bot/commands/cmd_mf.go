package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandMessageFreq(
	args string,
	message discordgo.Message,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	freq, err := cm.Discord.GetCurrentChannelMessageFrequency(message.ChannelID)
	if err != nil {
		return
	}
	cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%02f messages per second", freq))
	return
}
