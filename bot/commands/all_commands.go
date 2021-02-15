package commands

import (
	"time"

	"go.uber.org/zap"
)

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func (cm *CommandManager) LoadCommands() {
	commands := map[string]Command{
		"/commands": {
			Function:    cm.commandCommands,
			Description: "Displays a list of commands.",
		},
		"/help": {
			Function:    cm.commandHelp,
			Description: "Displays a list of commands.",
		},
		"/config": {
			Function:    cm.commandConfig,
			Description: "Configure command settings.",
		},
		"/readme": {
			Function:    cm.commandReadme,
			Description: "Force fetches the readme message.",
		},
		"/roles": {
			Function:    cm.commandRoles,
			Description: "List of roles and their IDs.",
		},
		"/say": {
			Function:    cm.commandSay,
			Description: "Say something as CJ.",
		},
		"/userinfo": {
			Function:    cm.commandUserInfo,
			Description: "Get a user's Burgershot forum info.",
		},
		"/whois": {
			Function:    cm.commandWhois,
			Description: "Display a Discord user's forum account name.",
		},
		"/stats": {
			Function:    cm.commandStats,
			Description: "Query a server through open.mp API",
			Cooldown:    time.Minute * 10,
		},
		"cj": {
			Function:    cm.commandCJQuote,
			Description: "Talk to CJ.",
			Cooldown:    time.Minute * 10,
		},
		"gmname": {
			Function:    cm.commandGmName,
			Description: "generates a professional gamemode name for your next NGG edit.",
			Cooldown:    time.Minute * 10,
		},
		"mpname": {
			Function:    cm.commandMP,
			Description: "scrapes the web for the next BIG samp ripoff.",
			Cooldown:    time.Minute * 10,
		},
		"dynamic": {
			Function:    cm.commandDynamic,
			Description: "inspiration for your next script.",
			Cooldown:    time.Minute * 10,
		},
		"rpname": {
			Function:    cm.commandRP,
			Description: "the next big unique dynamic server.",
			Cooldown:    time.Minute * 10,
		},
		"/wiki": {
			Function:    cm.commandWiki,
			Description: "Returns an article from open.mp wiki.",
		},
		"/top": {
			Function:    cm.commandTop,
			Description: "Rankings for most messages sent.",
			Cooldown:    time.Minute * 10,
		},
		"/toprep": {
			Function:    cm.commandTopRep,
			Description: "Rankings for most emojis sent.",
			Cooldown:    time.Minute * 10,
		},
		"/konesyntees": {
			Function:    cm.commandKonesyntees,
			Description: "Use superior Estonian technology to express your feelings like you've never before!",
			Cooldown:    time.Minute,
		},
		"/mf": {
			Function:    cm.commandMessageFreq,
			Description: "Message frequency",
			Cooldown:    time.Minute,
		},
		"/rep": {
			Function:    cm.commandRep,
			Description: "Know how many reactions your messages have gotten",
			Cooldown:    time.Second * 2,
		},
	}
	for k, v := range commands {
		v.Settings.Cooldown = cm.Config.DefaultCooldown
		v.Settings.Channels = []string{cm.Config.DefaultChannel}
		v.Settings.Roles = []string{cm.Config.DefaultRole}
		v.Settings.Command = k

		settings, found, err := cm.Storage.GetCommandSettings(k)
		if err != nil {
			zap.L().Fatal("failed to load command settings",
				zap.Error(err))
		}
		if found {
			v.Settings = settings
		} else {
			err = cm.Storage.SetCommandSettings(k, v.Settings)
			if err != nil {
				zap.L().Fatal("failed to assign command settings",
					zap.Error(err))
			}
		}

		commands[k] = v
	}

	cm.Commands = commands
}
