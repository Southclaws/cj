package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func (cm *CommandManager) LoadCommands() {
	commands := []Command{
		{
			Function:    cm.commandHelp,
			Name:        "/help",
			Description: "Displays a list of commands.",
		},
		{
			Function:         cm.commandConfig,
			Name:             "/config",
			Description:      "Configure command settings.",
			IsAdministrative: true,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "command",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "A valid CJ command",
					Required:    true,
				},
				{
					Name:        "config",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "A JSON configuration for the command",
					Required:    false,
				},
			},
		},
		{
			Function:         cm.commandGetMessageInfo,
			Name:             "/getmsginfo",
			Description:      "Get a message's info by ID",
			IsAdministrative: true,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "message-id",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "The message ID to retrieve from the DB",
					Required:    true,
				},
			},
		},
		{
			Function:    cm.commandGmName,
			Name:        "/gmname",
			Description: "generates a professional gamemode name for your next NGG edit.",
		},
		{
			Function:    cm.commandMP,
			Name:        "/mpname",
			Description: "scrapes the web for the next BIG samp ripoff.",
		},
		{
			Function:    cm.commandDynamic,
			Name:        "/dynamic",
			Description: "inspiration for your next script.",
		},
		{
			Function:    cm.commandRP,
			Name:        "/rpname",
			Description: "the next big unique dynamic server.",
		},
		{
			Function:    cm.commandWiki,
			Name:        "/wiki",
			Description: "Returns an article from open.mp wiki.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "search-term",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "The wiki term to search",
					Required:    true,
				},
			},
		},
		{
			Function:    cm.commandTop,
			Name:        "/top",
			Description: "Rankings for most messages sent.",
		},
		{
			Function:    cm.commandTopRep,
			Name:        "/toprep",
			Description: "Rankings for most emojis sent.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "reaction",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "Optional: Rankings for a specific emoji",
					Required:    false,
				},
			},
		},
		{
			Function:    cm.commandKonesyntees,
			Name:        "/konesyntees",
			Description: "Use superior Estonian technology to express your feelings like you've never before!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "input",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "Konesyntezing input",
					Required:    true,
				},
			},
		},
		{
			Function:    cm.commandRep,
			Name:        "/rep",
			Description: "Know how many reactions your messages have gotten",
		},
		{
			Function:    cm.commandMyTop,
			Name:        "/mytop",
			Description: "Know your rank.",
		},
		{
			Function:         cm.commandDebugReload,
			Name:             "/debugreload",
			Description:      "Force reload commands to Discord",
			IsAdministrative: true,
		},
		{
			Function:    cm.ltf,
			Name:        "/ltf",
			Description: "Rest in peace.",
		},
	}

	var discordCommands = []*discordgo.ApplicationCommand{}

	for k, v := range commands {
		v.Settings.Cooldown = cm.Config.DefaultCooldown
		v.Settings.Roles = []string{"all"}
		if v.IsAdministrative {
			// Operator, Admin, Test server role, Cj Config Man
			v.Settings.Roles = []string{"363383930143113216", "282353010192023552", "825041993251029032", "783100099969548288"}
		}
		v.Settings.Command = v.Name

		settings, found, err := cm.Storage.GetCommandSettings(v.Name)
		if err != nil {
			zap.L().Fatal("failed to load command settings",
				zap.Error(err))
		}
		if found {
			v.Settings = settings
		} else {
			err = cm.Storage.SetCommandSettings(v.Name, v.Settings)
			if err != nil {
				zap.L().Fatal("failed to assign command settings",
					zap.Error(err))
			}
		}

		commands[k] = v

		// Add an entry for the bulk overwrite list
		discordCommands = append(discordCommands, &discordgo.ApplicationCommand{
			Name:        strings.TrimLeft(v.Name, "/"),
			Description: v.Description,
			Options:     v.Options,
		})
	}

	// Cleanup of old commands
	// This is worth doing, e.g. if discord bugs out
	// or a command signature changes or is deleted.
	existingsDiscordCommands, _ := cm.Discord.S.ApplicationCommands(cm.Discord.S.State.User.ID, "")
	zap.L().Info("Existing slash commands:", zap.Any("cmds", existingsDiscordCommands))

	for _, existingCommand := range existingsDiscordCommands {
		commandExists := false
		for _, command := range discordCommands {
			if command.Name == existingCommand.Name {
				commandExists = true
				break
			}
		}
		if commandExists == false {
			zap.L().Info("Deleting non-existent slash command from Discord:", zap.Any("cmd", existingCommand))
			cm.Discord.S.ApplicationCommandDelete(cm.Discord.S.State.User.ID, "", existingCommand.ID)
		}
	}

	// Leave guild ID empty to indicate a global command
	zap.L().Info("Pushing commands to discord...")
	cm.Discord.S.ApplicationCommandBulkOverwrite(cm.Discord.S.State.User.ID, "", discordCommands)

	cm.Discord.S.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cm.TryFindAndFireCommand(i)
	})

	cm.Commands = commands
	zap.L().Info("Set up command handler complete")
}
