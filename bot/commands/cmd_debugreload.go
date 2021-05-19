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
		commandExists := false
		for _, command := range cm.Commands {
			if command.Name == existingCommand.Name {
				commandExists = true
				break
			}
		}
		if commandExists == false {
			go cm.Discord.S.ApplicationCommandDelete(cm.Discord.S.State.User.ID, "", existingCommand.ID)
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

	cm.Discord.S.InteractionResponseEdit(cm.Discord.S.State.User.ID, interaction.Interaction, &discordgo.WebhookEdit{
		Content: "Reload commands succeeded. This can take an hour to take effect.",
	})

	return
}
