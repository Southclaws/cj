package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func LoadCommands() map[string]Command {
	return map[string]Command{
		"verify": {
			Function:        commandVerify,
			Source:          CommandSourcePRIVATE,
			Description:     "Verify you are the owner of a SA:MP forum account",
			Usage:           "verify",
			RequireVerified: false,
			RequireAdmin:    false,
		},
	}
}

// CommandSource represents the source of a command.
type CommandSource int8

const (
	// CommandSourceNONE represents an invalid command source.
	CommandSourceNONE CommandSource = iota
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
}

// Command represents a public, private or administrative command
type Command struct {
	Function        func(cmdtext string, channel string) error
	Source          CommandSource
	Description     string
	Usage           string
	RequireVerified bool
	RequireAdmin    bool
}

// StartCommandManager creates a command manager for the app
func (app *App) StartCommandManager() {
	app.commandManager = &CommandManager{
		App:      app,
		Commands: make(map[string]Command),
	}

	app.commandManager.Commands = LoadCommands()
}

// Process is called on a command string to check whether it's a valid command
// and, if so, call the associated function.
func (cm CommandManager) Process(cmdtext string, channel string) (exists bool, source CommandSource, errs []error) {
	commandAndParameters := strings.SplitN(cmdtext, " ", 1)
	commandObject, exists := cm.Commands[strings.ToLower(commandAndParameters[0])]

	if !exists {
		return exists, source, nil
	}

	source = cm.getCommandSource(cmdtext, channel)

	if source == commandObject.Source {
		errs = append(errs, commandObject.Function(commandAndParameters[1], channel))
	}

	return exists, source, errs
}

func (cm CommandManager) getCommandSource(cmdtext string, channel string) CommandSource {
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
		json.Unmarshal(body, &channel)

		// Now we have one of these:
		// https://discordapp.com/developers/docs/resources/channel#dm-channel-object

		if channel.Private {
			return CommandSourcePRIVATE
		}
	}

	return CommandSourceOTHER
}
