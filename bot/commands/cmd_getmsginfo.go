package commands

import (
	"fmt"
	"time"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandGetMessageInfo(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	messageId := args["message-id"].StringValue()
	var chatLog storage.ChatLog
	chatLog, err = cm.Storage.GetMessageByID(messageId)
	if err != nil {
		cm.replyDirectly(interaction, "Bad message ID")
		return
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

	cm.replyDirectly(interaction, discordMessage)
	return
}
