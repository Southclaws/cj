package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

var sayDeniedRoleIDs = []string{
	"1002922725553217648", // Clown
	"400250542628274177",  // Caged
	"995816487610753094",  // annoyed me
	"1016047260364198008", // Not Cool
	"833325019252785173",  // Doesn't deserve to embed
	"818457955690872832",  // Doesn't deserve to react
	"996883259252297758",  // No open.mp support
	"910950457680212088",  // No Server Adverts
	"841368374356738078",  // Suffers from dunning-kruger
	"987825514511220867",  // Muted
	"1204891485867352144", // Can't @everyone
}

func (cm *CommandManager) commandSay(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	message := ""
	if option, ok := args["message"]; ok {
		message = strings.TrimSpace(option.StringValue())
	}

	if message == "" {
		cm.replyDirectly(interaction, "Message cannot be empty.")
		return
	}

	err = cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:         message,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
	return
}
