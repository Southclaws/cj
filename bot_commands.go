package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	gocache "github.com/patrickmn/go-cache"
)

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func LoadCommands(app *App) map[string]Command {
	return map[string]Command{
		"verify": {
			Function:        commandVerify,
			Source:          CommandSourcePRIVATE,
			Description:     "Verify you are the owner of a SA:MP forum account",
			Usage:           "verify",
			RequireVerified: false,
			RequireAdmin:    false,
			Context:         true,
		},
		"/say": {
			Function:        commandSay,
			Source:          CommandSourceADMINISTRATIVE,
			Description:     "Verify you are the owner of a SA:MP forum account",
			Usage:           "say",
			RequireVerified: false,
			RequireAdmin:    false,
			Context:         false,
		},
		"/whois": {
			Function:        commandWhois,
			Source:          CommandSourcePRIMARY,
			Description:     app.locale.GetLangString("en", "CommandWhoisUsage"),
			Usage:           "whois",
			RequireVerified: true,
			RequireAdmin:    false,
			Context:         false,
		},
	}
}

// CommandSource represents the source of a command.
type CommandSource int8

const (
	// CommandSourceADMINISTRATIVE are commands in the administrator channel,
	// mainly for admin work that may clutter up the primary channel.
	CommandSourceADMINISTRATIVE CommandSource = iota
	// CommandSourcePRIMARY are primary channel commands visible to all users
	// and mainly used for fun and group activity commands.
	CommandSourcePRIMARY CommandSource = iota
	// CommandSourcePRIVATE are private channel commands for dealing with
	// sensitive information such as verification.
	CommandSourcePRIVATE CommandSource = iota
	// CommandSourceOTHER represents any other channel that does not fall into
	// the above sources.
	CommandSourceOTHER CommandSource = iota
)

// CommandManager stores command state
type CommandManager struct {
	App      *App
	Commands map[string]Command
	Contexts *gocache.Cache
}

// Command represents a public, private or administrative command
type Command struct {
	commandManager  *CommandManager
	Function        func(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error)
	Source          CommandSource
	Description     string
	Usage           string
	RequireVerified bool
	RequireAdmin    bool
	Context         bool
}

// StartCommandManager creates a command manager for the app
func (app *App) StartCommandManager() {
	app.commandManager = &CommandManager{
		App:      app,
		Commands: make(map[string]Command),
		Contexts: gocache.New(5*time.Minute, 30*time.Second),
	}

	app.commandManager.Commands = LoadCommands(app)
}

// Process is called on a command string to check whether it's a valid command
// and, if so, call the associated function.
// nolint:gocyclo
func (cm CommandManager) Process(message discordgo.Message) (exists bool, source CommandSource, errs []error) {
	debug("[commands:Process] message: '%s'", message.Content)

	source = cm.getCommandSource(message.ChannelID)

	contextCommand, found := cm.Contexts.Get(message.Author.ID)
	if found {
		contextCommand := contextCommand.(Command)
		debug("[commands:Process] User is currently in context of command '%s'", contextCommand.Usage)
		if contextCommand.Source == source {
			var continueContext bool
			continueContext, errs = cm.ProcessContext(contextCommand, message.Content, message)
			if !continueContext {
				cm.Contexts.Delete(message.Author.ID)
			}
			return true, source, errs
		}
	}

	commandAndParameters := strings.SplitN(message.Content, " ", 2)
	commandTrigger := strings.ToLower(commandAndParameters[0])
	commandArgument := ""
	if len(commandAndParameters) > 1 {
		commandArgument = commandAndParameters[1]
	}

	commandObject, exists := cm.Commands[commandTrigger]
	commandObject.commandManager = &cm

	if !exists {
		debug("[commands:Process] command '%s' does not exist", commandTrigger)
		return exists, source, nil
	}

	debug("[commands:Process] command exists, source: %v required: %v", source, commandObject.Source)

	if source == commandObject.Source {
		debug("[commands:Process] command source matches required source")
		switch source {
		case CommandSourceADMINISTRATIVE:
			if message.ChannelID != cm.App.config.AdministrativeChannel {
				return exists, source, errs
			}
		case CommandSourcePRIMARY:
			if message.ChannelID != cm.App.config.PrimaryChannel {
				return exists, source, errs
			}
		}

		success, enterContext, e := commandObject.Function(cm, commandArgument, message, false)
		errs = append(errs, e)
		if enterContext {
			if commandObject.Context {
				debug("[commands:Process] command is contextual, placing user in context")
				cm.Contexts.Set(message.Author.ID, commandObject, gocache.DefaultExpiration)
			}
		}
		if !success {
			_, e := cm.App.discordClient.ChannelMessageSend(message.ChannelID, commandObject.Usage)
			errs = append(errs, e)
		}
	}

	return exists, source, errs
}

// ProcessContext re-runs a Command function if the user is currently in a
// Command's context.
func (cm CommandManager) ProcessContext(command Command, cmdtext string, message discordgo.Message) (continueContext bool, errs []error) {
	_, continueContext, e := command.Function(cm, cmdtext, message, true)
	if e != nil {
		errs = append(errs, e)
	}
	return continueContext, errs
}

func (cm CommandManager) getCommandSource(channel string) CommandSource {
	if channel == cm.App.config.AdministrativeChannel {
		return CommandSourceADMINISTRATIVE
	} else if channel == cm.App.config.PrimaryChannel {
		return CommandSourcePRIMARY
	} else {
		// discordgo has not implemented private channel objects (DM Channels)
		// so we have to perform the request manually and unmarshal the response
		// object into a `ChannelDM` object.
		var err error
		var req *http.Request
		var response *http.Response
		var body []byte
		if req, err = http.NewRequest("GET", discordgo.EndpointChannel(channel), nil); err != nil {
			log.Print(err)
		}
		req.Header.Add("Authorization", "Bot "+cm.App.config.DiscordToken)
		if response, err = cm.App.httpClient.Do(req); err != nil {
			log.Print(err)
		}
		if body, err = ioutil.ReadAll(response.Body); err != nil {
			log.Print(err)
		}
		channel := ChannelDM{}
		err = json.Unmarshal(body, &channel)
		if err != nil {
			log.Print(err)
		}

		// Now we have one of these:
		// https://discordapp.com/developers/docs/resources/channel#dm-channel-object

		if channel.Private {
			return CommandSourcePRIVATE
		}
	}

	return CommandSourceOTHER
}
