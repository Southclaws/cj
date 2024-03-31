package commands

import (
	"strings"

	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandDebugReload(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.replyDirectly(interaction, "CJ will reload commands. Stand by.")

	existingsDiscordCommands, _ := cm.Discord.S.ApplicationCommands(cm.Discord.S.State.User.ID, "")
	for _, existingCommand := range existingsDiscordCommands {
		cm.Discord.S.ApplicationCommandDelete(cm.Discord.S.State.User.ID, "", existingCommand.ID)
	}

	for _, guild := range cm.Discord.S.State.Guilds {
		existingsDiscordCommands, _ = cm.Discord.S.ApplicationCommands(cm.Discord.S.State.User.ID, guild.ID)
		for _, existingCommand := range existingsDiscordCommands {
			cm.Discord.S.ApplicationCommandDelete(cm.Discord.S.State.User.ID, guild.ID, existingCommand.ID)
		}
	}

	var discordCommands []*discordgo.ApplicationCommand
	for _, command := range cm.Commands {
		discordCommands = append(discordCommands, &discordgo.ApplicationCommand{
			Name:        strings.TrimLeft(command.Name, "/"),
			Description: command.Description,
			Options:     command.Options,
		})
	}

	cm.Discord.S.ApplicationCommandBulkOverwrite(cm.Discord.S.State.User.ID, "", discordCommands)

	content := "Reload commands succeeded. This can take an hour to take effect."
	cm.Discord.S.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})

	return
}
