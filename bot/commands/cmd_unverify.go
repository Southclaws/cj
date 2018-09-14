package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func (cm *CommandManager) commandUnVerify(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	if len(message.Mentions) != 1 {
		err = errors.New("Expected one mention")
		return
	}

	target := message.Mentions[0]

	err = cm.Storage.RemoveUser(target.ID)
	if err != nil {
		_, err = cm.Discord.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	_, err = cm.Discord.ChannelMessageSend(message.ChannelID, "User un-verified")

	return
}
