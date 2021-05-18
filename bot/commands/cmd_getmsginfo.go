package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func (cm *CommandManager) commandGetMessageInfo(
	args string,
	message discordgo.Message,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	messageId := strings.TrimSpace(args)
	var chatLog storage.ChatLog
	chatLog, err = cm.Storage.GetMessageByID(messageId)
	if err != nil {
		return false, errors.New("Bad message ID")
	}
	discordMessage := fmt.Sprintf(
		"Message ID: %s\n"+
			"Timestamp: %v\n"+
			"Channel ID: %s\n"+
			"User ID: %s\n"+
			"Message: %s\n",
		chatLog.DiscordMessageID,
		time.Unix(chatLog.Timestamp, 0),
		chatLog.DiscordChannel,
		chatLog.DiscordUserID,
		chatLog.Message)

	cm.Discord.ChannelMessageSend(message.ChannelID, discordMessage)
	return
}
