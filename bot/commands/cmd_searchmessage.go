package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

const (
	searchMessageChannelID     = "948604467887083550"
	searchMessageContextPrefix = "searchmessage:"
)

type searchMessageContext struct {
	ID            string
	RequesterID   string
	AuthorName    string
	AuthorIconURL string
	Messages      []storage.ChatLog
	Page          int
}

func (cm *CommandManager) commandSearchMessage(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	if interaction.ChannelID != searchMessageChannelID {
		return
	}

	accountID := strings.TrimSpace(args["account-id"].StringValue())
	query := normalizeSearchQuery(args["message"].StringValue())

	if _, parseErr := strconv.ParseUint(accountID, 10, 64); parseErr != nil {
		cm.replyDirectly(interaction, "The account ID must be a valid Discord user ID.")
		return
	}
	if query == "" {
		cm.replyDirectly(interaction, "The message search text cannot be empty.")
		return
	}

	cm.sendThinkingResponse(interaction)

	messages, searchErr := cm.Storage.SearchMessages(accountID, query)
	if searchErr != nil {
		cm.editOriginalResponse(interaction, "The message search failed.")
		return
	}
	if len(messages) == 0 {
		cm.editOriginalResponse(interaction, "No matching messages were found.")
		return
	}

	authorName := accountID
	authorIconURL := ""
	if user, userErr := cm.Discord.S.User(accountID); userErr == nil {
		authorName = user.Username
		authorIconURL = user.AvatarURL("")
	}

	searchContext := searchMessageContext{
		ID:            interaction.ID,
		RequesterID:   interaction.Member.User.ID,
		AuthorName:    authorName,
		AuthorIconURL: authorIconURL,
		Messages:      messages,
	}
	cm.Contexts.SetDefault(searchMessageContextPrefix+searchContext.ID, searchContext)

	embed, components := cm.searchMessagePage(searchContext)
	content := ""
	embeds := []*discordgo.MessageEmbed{embed}
	cm.Discord.S.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content:    &content,
		Embeds:     &embeds,
		Components: &components,
	})
	return
}

func normalizeSearchQuery(query string) string {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return query
	}

	first := query[0]
	last := query[len(query)-1]
	if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
		return strings.TrimSpace(query[1 : len(query)-1])
	}
	return query
}

func (cm *CommandManager) handleSearchMessageComponent(interaction *discordgo.InteractionCreate) {
	data := interaction.MessageComponentData()
	if !strings.HasPrefix(data.CustomID, searchMessageContextPrefix) {
		return
	}

	parts := strings.Split(data.CustomID, ":")
	if len(parts) != 3 {
		return
	}

	contextKey := searchMessageContextPrefix + parts[1]
	stored, found := cm.Contexts.Get(contextKey)
	searchContext, valid := stored.(searchMessageContext)
	if !found || !valid {
		cm.replyToExpiredSearch(interaction)
		return
	}

	if interaction.Member == nil || interaction.Member.User.ID != searchContext.RequesterID {
		cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Only the person who ran the search can change its page.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	switch parts[2] {
	case "previous":
		if searchContext.Page > 0 {
			searchContext.Page--
		}
	case "next":
		if searchContext.Page < len(searchContext.Messages)-1 {
			searchContext.Page++
		}
	default:
		return
	}

	cm.Contexts.SetDefault(contextKey, searchContext)
	embed, components := cm.searchMessagePage(searchContext)
	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

func (cm *CommandManager) searchMessagePage(searchContext searchMessageContext) (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	message := searchContext.Messages[searchContext.Page]
	messageURL := fmt.Sprintf(
		"https://discord.com/channels/%s/%s/%s",
		cm.Config.GuildID,
		message.DiscordChannel,
		message.DiscordMessageID,
	)

	embed := &discordgo.MessageEmbed{
		Description: message.Message,
		Color:       0x3498DB,
		Timestamp:   time.Unix(message.Timestamp, 0).UTC().Format(time.RFC3339),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    searchContext.AuthorName,
			IconURL: searchContext.AuthorIconURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Result %d of %d", searchContext.Page+1, len(searchContext.Messages)),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Channel",
				Value:  fmt.Sprintf("<#%s> (`%s`)", message.DiscordChannel, message.DiscordChannel),
				Inline: true,
			},
			{
				Name:   "Account",
				Value:  fmt.Sprintf("<@%s> (`%s`)", message.DiscordUserID, message.DiscordUserID),
				Inline: true,
			},
			{
				Name:   "Date and time",
				Value:  fmt.Sprintf("<t:%d:F>", message.Timestamp),
				Inline: true,
			},
			{
				Name:  "Original message",
				Value: fmt.Sprintf("[Jump to message](%s)", messageURL),
			},
		},
	}

	if len(searchContext.Messages) == 1 {
		return embed, nil
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Previous",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf("%s%s:previous", searchMessageContextPrefix, searchContext.ID),
					Disabled: searchContext.Page == 0,
				},
				discordgo.Button{
					Label:    "Next",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("%s%s:next", searchMessageContextPrefix, searchContext.ID),
					Disabled: searchContext.Page == len(searchContext.Messages)-1,
				},
			},
		},
	}
	return embed, components
}

func (cm *CommandManager) replyToExpiredSearch(interaction *discordgo.InteractionCreate) {
	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This search has expired. Run the command again.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
