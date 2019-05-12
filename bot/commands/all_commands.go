package commands

import (
	"time"

	"go.uber.org/zap"

	"github.com/Southclaws/cj/types"
)

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func (cm *CommandManager) LoadCommands() {
	commands := map[string]Command{
		"/commands": {
			Function:    cm.commandCommands,
			Source:      CommandSourcePRIMARY,
			Description: "Displays a list of commands.",
		},
		"/help": {
			Function:    cm.commandHelp,
			Source:      CommandSourcePRIMARY,
			Description: "Displays a list of commands.",
		},
		"/config": {
			Function:    cm.commandConfig,
			Source:      CommandSourcePRIMARY,
			Description: "Configure command settings.",
		},
		"/roles": {
			Function:    cm.commandRoles,
			Description: "List of roles and their IDs.",
		},
		"verify": {
			Function:    cm.commandVerify,
			Source:      CommandSourcePRIVATE,
			Description: "Verify you are the owner of a Burgershot forum account.",
			Settings:    types.CommandSettings{Private: true},
		},
		"/say": {
			Function:    cm.commandSay,
			Source:      CommandSourceADMINISTRATIVE,
			Description: "Say something as CJ.",
		},
		"/userinfo": {
			Function:    cm.commandUserInfo,
			Source:      CommandSourcePRIMARY,
			Description: "Get a user's Burgershot forum info.",
		},
		"/whois": {
			Function:    cm.commandWhois,
			Source:      CommandSourcePRIMARY,
			Description: "Display a Discord user's forum account name.",
		},
		"/setverify": {
			Function:    cm.commandSetVerify,
			Source:      CommandSourceADMINISTRATIVE,
			Description: "Manually verify a user.",
		},
		"/unverify": {
			Function:    cm.commandUnVerify,
			Source:      CommandSourceADMINISTRATIVE,
			Description: "Manually unverify a user.",
		},
		"cj": {
			Function:    cm.commandCJQuote,
			Source:      CommandSourcePRIMARY,
			Description: "Talk to CJ.",
			Cooldown:    time.Minute * 10,
		},
		"gmname": {
			Function:    cm.commandGmName,
			Source:      CommandSourcePRIMARY,
			Description: "generates a professional gamemode name for your next NGG edit.",
			Cooldown:    time.Minute * 10,
		},
		"mpname": {
			Function:    cm.commandMP,
			Source:      CommandSourcePRIMARY,
			Description: "scrapes the web for the next BIG samp ripoff.",
			Cooldown:    time.Minute * 10,
		},
		"dynamic": {
			Function:    cm.commandDynamic,
			Source:      CommandSourcePRIMARY,
			Description: "inspiration for your next script.",
			Cooldown:    time.Minute * 10,
		},
		"rpname": {
			Function:    cm.commandRP,
			Source:      CommandSourcePRIMARY,
			Description: "the next big unique dynamic server.",
			Cooldown:    time.Minute * 10,
		},
		"/wiki": {
			Function:    cm.commandWiki,
			Source:      CommandSourcePRIMARY,
			Description: "Returns an article from SA:MP wiki.",
		},
		"/top": {
			Function:    cm.commandTop,
			Source:      CommandSourcePRIMARY,
			Description: "Rankings for most messages sent.",
			Cooldown:    time.Minute * 10,
		},
		"/konesyntees": {
			Function:    cm.commandKonesyntees,
			Source:      CommandSourcePRIMARY,
			Description: "Use superior Estonian technology to express your feelings like you've never before!",
			Cooldown:    time.Minute,
		},
		"/stats": {
			Function:    cm.commandStats,
			Source:      CommandSourcePRIMARY,
			Description: "Query a server for information using the SA:MP servers API.",
		},
		"/breakthings": {
			Function:    cm.commandBLNS,
			Source:      CommandSourcePRIMARY,
			Description: "What's the latest iOS unicode bug???",
			Cooldown:    time.Minute,
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
