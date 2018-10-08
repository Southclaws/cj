package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func (cm *CommandManager) commandTop(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	top, err := cm.Storage.GetTopMessages(10)
	if err != nil {
		return false, errors.Wrap(err, "failed to get message rankings")
	}

	statsMessage := strings.Builder{}
	statsMessage.WriteString("Statistics") //nolint:errcheck

	embed := &discordgo.MessageEmbed{Color: 0x3498DB}
	for _, tm := range top {
		statsMessage.WriteString(fmt.Sprintf("**%s** - %d\n", tm.User, tm.Messages)) //nolint:errcheck
	}

	embed.Description = statsMessage.String()

	_, err = cm.Discord.ChannelMessageSendEmbed(cm.Config.PrimaryChannel, embed)
	return
}
