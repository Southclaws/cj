package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

const sayWebhookName = "CJ Say"

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
		cm.replyDirectly(interaction, "You gotta give me something to say, homie.")
		return
	}

	err = cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return
	}

	webhook, threadID, webhookErr := cm.getSayWebhook(interaction.ChannelID)
	if webhookErr != nil {
		cm.editOriginalResponse(interaction, "Ah shit, I can't say that here. Give me permission to manage webhooks, homie.")
		err = webhookErr
		return
	}

	member := *interaction.Member
	member.GuildID = interaction.GuildID
	params := &discordgo.WebhookParams{
		Content:         message,
		Username:        member.DisplayName(),
		AvatarURL:       member.AvatarURL(""),
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	}
	if threadID == "" {
		_, err = cm.Discord.S.WebhookExecute(webhook.ID, webhook.Token, true, params)
	} else {
		_, err = cm.Discord.S.WebhookThreadExecute(webhook.ID, webhook.Token, true, threadID, params)
	}
	if err != nil {
		cm.editOriginalResponse(interaction, "Ah shit, the webhook let me down. Try that again, homie.")
		return
	}

	err = cm.Discord.S.InteractionResponseDelete(interaction.Interaction)
	return
}

func (cm *CommandManager) getSayWebhook(channelID string) (*discordgo.Webhook, string, error) {
	channel, err := cm.Discord.S.Channel(channelID)
	if err != nil {
		return nil, "", fmt.Errorf("get /say channel: %w", err)
	}

	webhookChannelID := channel.ID
	threadID := ""
	if channel.IsThread() {
		webhookChannelID = channel.ParentID
		threadID = channel.ID
	}

	webhooks, err := cm.Discord.S.ChannelWebhooks(webhookChannelID)
	if err != nil {
		return nil, "", fmt.Errorf("get /say webhooks: %w", err)
	}

	botUser := cm.Discord.S.State.User
	for _, webhook := range webhooks {
		if webhook.Name == sayWebhookName && webhook.User != nil && botUser != nil && webhook.User.ID == botUser.ID {
			return webhook, threadID, nil
		}
	}

	webhook, err := cm.Discord.S.WebhookCreate(
		webhookChannelID,
		sayWebhookName,
		"",
		discordgo.WithAuditLogReason("Webhook for /say messages"),
	)
	if err != nil {
		return nil, "", fmt.Errorf("create /say webhook: %w", err)
	}

	return webhook, threadID, nil
}
