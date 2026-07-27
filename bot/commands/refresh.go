package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type RefreshResult struct {
	Deleted int
	Added   int
}

func (cm *CommandManager) RefreshCommands() (result RefreshResult, err error) {
	existing, err := cm.Discord.S.ApplicationCommands(cm.Discord.S.State.User.ID, "")
	if err != nil {
		return
	}
	for _, c := range existing {
		if err = cm.Discord.S.ApplicationCommandDelete(cm.Discord.S.State.User.ID, "", c.ID); err != nil {
			return
		}
		result.Deleted++
	}

	for _, guild := range cm.Discord.S.State.Guilds {
		guildExisting, gerr := cm.Discord.S.ApplicationCommands(cm.Discord.S.State.User.ID, guild.ID)
		if gerr != nil {
			err = gerr
			return
		}
		for _, c := range guildExisting {
			if err = cm.Discord.S.ApplicationCommandDelete(cm.Discord.S.State.User.ID, guild.ID, c.ID); err != nil {
				return
			}
			result.Deleted++
		}
	}

	var discordCommands []*discordgo.ApplicationCommand
	for _, command := range cm.Commands {
		if !commandEnabled(command, cm.GuildSettings) {
			continue
		}
		discordCommands = append(discordCommands, &discordgo.ApplicationCommand{
			Name:        strings.TrimLeft(command.Name, "/"),
			Description: command.Description,
			Options:     command.Options,
		})
	}

	if _, err = cm.Discord.S.ApplicationCommandBulkOverwrite(cm.Discord.S.State.User.ID, "", discordCommands); err != nil {
		return
	}
	result.Added = len(discordCommands)

	return
}
