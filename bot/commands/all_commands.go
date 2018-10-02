package commands

import (
	"time"
)

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func (cm *CommandManager) LoadCommands() {
	cm.Commands = map[string]Command{
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
		"verify": {
			Function:    cm.commandVerify,
			Source:      CommandSourcePRIVATE,
			Description: "Verify you are the owner of a SA:MP forum account.",
			Context:     true,
		},
		"/say": {
			Function:    cm.commandSay,
			Source:      CommandSourceADMINISTRATIVE,
			Description: "Say something as CJ.",
			Context:     false,
		},
		"/userinfo": {
			Function:    cm.commandUserInfo,
			Source:      CommandSourcePRIMARY,
			Description: "Get a user's SA:MP forum info.",
			Context:     false,
		},
		"/whois": {
			Function:    cm.commandWhois,
			Source:      CommandSourcePRIMARY,
			Description: "Display a Discord user's forum account name.",
			Context:     false,
		},
		"/setverify": {
			Function:    cm.commandSetVerify,
			Source:      CommandSourceADMINISTRATIVE,
			Description: "Manually verify a user.",
			Context:     false,
		},
		"/unverify": {
			Function:    cm.commandUnVerify,
			Source:      CommandSourceADMINISTRATIVE,
			Description: "Manually unverify a user.",
			Context:     false,
		},
		"cj": {
			Function:    cm.commandCJQuote,
			Source:      CommandSourcePRIMARY,
			Description: "Talk to CJ.",
			Context:     false,
			Cooldown:    time.Minute * 10,
		},
		"gmname": {
			Function:    cm.commandGmName,
			Source:      CommandSourcePRIMARY,
			Description: "generates a professional gamemode name for your next NGG edit.",
			Context:     false,
			Cooldown:    time.Minute * 10,
		},
		"mpname": {
			Function:    cm.commandMP,
			Source:      CommandSourcePRIMARY,
			Description: "scrapes the web for the next BIG samp ripoff.",
			Context:     false,
			Cooldown:    time.Minute * 10,
		},
		"dynamic": {
			Function:    cm.commandDynamic,
			Source:      CommandSourcePRIMARY,
			Description: "inspiration for your next script.",
			Context:     false,
			Cooldown:    time.Minute * 10,
		},
		"rpname": {
			Function:    cm.commandRP,
			Source:      CommandSourcePRIMARY,
			Description: "the next big unique dynamic server.",
			Context:     false,
			Cooldown:    time.Minute * 10,
		},
		"/wiki": {
			Function:    cm.commandWiki,
			Source:      CommandSourcePRIMARY,
			Description: "Returns an article from SA:MP wiki.",
			Context:     false,
		},
	}
}
