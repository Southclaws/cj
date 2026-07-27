package commands

import (
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

	_, err = cm.RefreshCommands()

	content := "Reload commands succeeded. This can take an hour to take effect."
	if err != nil {
		content = "Reload commands failed: " + err.Error()
	}
	cm.Discord.S.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})

	return
}
