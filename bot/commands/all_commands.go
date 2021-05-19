package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func (cm *CommandManager) LoadCommands() {
	commands := []Command{
		{
			Function:    cm.commandCommands,
			Name:        "/commands",
			Description: "Displays a list of commands.",
		},
		{
			Function:    cm.commandHelp,
			Name:        "/help",
			Description: "Displays a list of commands.",
		},
		{
			Function:    cm.commandConfig,
			Name:        "/config",
			Description: "Configure command settings.",
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
			Function:    cm.commandRoles,
			Name:        "/roles",
			Description: "List of roles and their IDs.",
		},
		{
			Function:    cm.commandSay,
			Name:        "/sayylmao",
			Description: "Say something as CJ.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "message",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "The message to echo back to you.",
					Required:    true,
				},
			},
		},
		{
			Function:    cm.commandGetMessageInfo,
			Name:        "/getmsginfo",
			Description: "Get a message's info by ID",
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
			Cooldown:    time.Minute * 10,
		},
		{
			Function:    cm.commandMP,
			Name:        "/mpname",
			Description: "scrapes the web for the next BIG samp ripoff.",
			Cooldown:    time.Minute * 10,
		},
		{
			Function:    cm.commandDynamic,
			Name:        "/dynamic",
			Description: "inspiration for your next script.",
			Cooldown:    time.Minute * 10,
		},
		{
			Function:    cm.commandRP,
			Name:        "/rpname",
			Description: "the next big unique dynamic server.",
			Cooldown:    time.Minute * 10,
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
			Cooldown:    time.Minute * 10,
		},
		{
			Function:    cm.commandTopRep,
			Name:        "/toprep",
			Description: "Rankings for most emojis sent.",
			Cooldown:    time.Minute * 10,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "emoji",
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
			Cooldown:    time.Minute,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "input",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "Konesyntezing input",
					Required:    false,
				},
			},
		},
		{
			Function:    cm.commandMessageFreq,
			Name:        "/mf",
			Description: "Message frequency",
			Cooldown:    time.Minute,
		},
		{
			Function:    cm.commandRep,
			Name:        "/rep",
			Description: "Know how many reactions your messages have gotten",
			Cooldown:    time.Second * 2,
		},
		{
			Function:    cm.commandMyTop,
			Name:        "/mytop",
			Description: "Know your rank.",
			Cooldown:    time.Minute * 10,
		},
		{
			Function:    cm.ltf,
			Name:        "/ltf",
			Description: "Rest in peace.",
			Cooldown:    time.Minute * 10,
		},
	}

	// Cleanup of existing commands
	// This is worth doing, e.g. if discord bugs out
	// or a command signature changes or is deleted.
	for _, guild := range cm.Discord.S.State.Guilds {
		commands, _ := cm.Discord.S.ApplicationCommands(cm.Discord.S.State.User.ID, guild.ID)
		for _, command := range commands {
			cm.Discord.S.ApplicationCommandDelete(cm.Discord.S.State.User.ID, guild.ID, command.ID)
		}
	}

	for k, v := range commands {
		v.Settings.Cooldown = cm.Config.DefaultCooldown
		v.Settings.Channels = []string{cm.Config.DefaultChannel}
		v.Settings.Roles = []string{cm.Config.DefaultRole}
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

		// Register the command to discord
		for _, guild := range cm.Discord.S.State.Guilds {
			_, err = cm.Discord.S.ApplicationCommandCreate(cm.Discord.S.State.User.ID, guild.ID, &discordgo.ApplicationCommand{
				Name:        strings.TrimLeft(v.Name, "/"),
				Description: v.Description,
				Options:     v.Options,
			})
			if err != nil {
				zap.L().Error("Error creating command!", zap.Any("At command:", v.Name))
				zap.L().Error("Error creating command!", zap.Error(err))
			}
		}
	}

	cm.Discord.S.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cm.TryFindAndFireCommand(i)
	})

	cm.Commands = commands
}
