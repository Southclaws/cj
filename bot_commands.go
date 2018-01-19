package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	gocache "github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

// Command represents a public, private or administrative command
type Command struct {
	commandManager  *CommandManager
	Function        func(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error)
	Source          CommandSource
	ParametersRange CommandParametersRange
	Description     string
	Usage           string
	Example         string
	ErrorMessage    string
	RequireVerified bool
	RequireAdmin    bool
	Context         bool
}

// LoadCommands is called on initialisation and is responsible for registering
// all commands and binding them to functions.
func LoadCommands(app *App) map[string]Command {
	return map[string]Command{
		app.locale.GetLangString(app.config.Language, "CommandVerifyName"): {
			Function:    commandVerify,
			Source:      CommandSourcePRIVATE,
			Description: app.locale.GetLangString(app.config.Language, "CommandVerifyDescription"),
			Usage:       app.locale.GetLangString(app.config.Language, "CommandVerifyName"),
			ParametersRange: CommandParametersRange{
				Minimum: -1,
				Maximum: -1,
			},
			RequireVerified: false,
			RequireAdmin:    false,
			Context:         true,
		},
		app.locale.GetLangString(app.config.Language, "CommandSayName"): {
			Function:    commandSay,
			Source:      CommandSourceADMINISTRATIVE,
			Description: app.locale.GetLangString(app.config.Language, "CommandSayDescription"),
			Usage:       app.locale.GetLangString(app.config.Language, "CommandSayUsage"),
			Example:     app.locale.GetLangString(app.config.Language, "CommandSayExample"),
			ParametersRange: CommandParametersRange{
				Minimum: 1,
				Maximum: -1,
			},
			RequireVerified: false,
			RequireAdmin:    false,
			Context:         false,
		},
		app.locale.GetLangString(app.config.Language, "CommandUserinfoName"): {
			Function:    commandUserInfo,
			Source:      CommandSourcePRIMARY,
			Description: app.locale.GetLangString(app.config.Language, "CommandUserInfoDescription"),
			Usage:       app.locale.GetLangString(app.config.Language, "CommandUserInfoUsage"),
			Example:     app.locale.GetLangString(app.config.Language, "CommandUserInfoExample"),
			ParametersRange: CommandParametersRange{
				Minimum: 1,
				Maximum: 5,
			},
			ErrorMessage:    app.locale.GetLangString(app.config.Language, "CommandErrorNoMention"),
			RequireVerified: true,
			RequireAdmin:    false,
			Context:         false,
		},
		app.locale.GetLangString(app.config.Language, "CommandWhoisName"): {
			Function:    commandWhois,
			Source:      CommandSourcePRIMARY,
			Description: app.locale.GetLangString(app.config.Language, "CommandWhoisDescription", "%s"),
			Usage:       app.locale.GetLangString(app.config.Language, "CommandWhoisUsage", "%s"),
			Example:     app.locale.GetLangString(app.config.Language, "CommandWhoisExample", "%s"),
			ParametersRange: CommandParametersRange{
				Minimum: 1,
				Maximum: 5,
			},
			ErrorMessage:    app.locale.GetLangString(app.config.Language, "CommandErrorNoMention"),
			RequireVerified: true,
			RequireAdmin:    false,
			Context:         false,
		},
		"cj": {
			Function:    commandCJQuote,
			Source:      CommandSourcePRIMARY,
			Description: app.locale.GetLangString(app.config.Language, "CommandCjDescription"),
			Usage:       "cj",
			ParametersRange: CommandParametersRange{
				Minimum: -1,
				Maximum: -1,
			},
			ErrorMessage:    app.locale.GetLangString(app.config.Language, "CommandErrorNoMention", "%s"),
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

// CommandParametersRange represents minimum value and maximum value number of parameters for a command
type CommandParametersRange struct {
	Minimum int
	Maximum int
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
	var err error

	source, err = cm.getCommandSource(message)
	if err != nil {
		errs = []error{err}
		return
	}

	contextCommand, found := cm.Contexts.Get(message.Author.ID)
	if found {
		logger.Debug("found existing command context", zap.String("id", message.Author.ID))

		contextCommand := contextCommand.(Command)
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
	commandParametersCount := 0
	commandTrigger := strings.ToLower(commandAndParameters[0])
	commandArgument := ""

	if len(commandAndParameters) > 1 {
		commandArgument = commandAndParameters[1]
		commandParametersCount = strings.Count(commandArgument, " ") + 1
	}

	commandObject, exists := cm.Commands[commandTrigger]
	commandObject.commandManager = &cm

	if !exists {
		logger.Debug("ignoring command that does not exist",
			zap.String("command", commandTrigger))
		return exists, source, nil
	}

	if source != commandObject.Source {
		logger.Debug("ignoring command with incorrect source",
			zap.String("command", commandTrigger),
			zap.Any("source", source),
			zap.Any("wantSource", commandObject.Source))
		return exists, source, nil
	}

	switch source {
	case CommandSourceADMINISTRATIVE:
		if message.ChannelID != cm.App.config.AdministrativeChannel {
			logger.Debug("ignoring admin command used in wrong channel", zap.String("command", commandTrigger))
			return exists, source, errs
		}
	case CommandSourcePRIMARY:
		if message.ChannelID != cm.App.config.PrimaryChannel {
			logger.Debug("ignoring primary channel command used in wrong channel", zap.String("command", commandTrigger))
			return exists, source, errs
		}
	}

	// Check if the user is an administrator.
	if commandObject.RequireAdmin && cm.App.config.Admin != message.Author.ID {
		logger.Debug("ignoring admin command used by non-admin", zap.String("command", commandTrigger))

		_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString(cm.App.config.Language, "CommandRequireAdministrator", message.Author.ID))
		if err != nil {
			errs = append(errs, err)
		}

		return exists, source, errs
	}

	// Check if the user is verified.
	if commandObject.RequireVerified {
		verified, err := cm.App.IsUserVerified(message.Author.ID)
		if err != nil {
			errs = append(errs, err)
			return exists, source, errs
		}
		if !verified {
			logger.Debug("ignoring command that requires verification from non-verified user", zap.String("command", commandTrigger))

			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString(cm.App.config.Language, "CommandRequireVerification", message.Author.ID))
			if err != nil {
				errs = append(errs, err)
			}
			return exists, source, errs
		}
	}

	// Check if we have the required number of parameters.
	if commandObject.ParametersRange.Minimum > -1 && commandParametersCount < commandObject.ParametersRange.Minimum {
		logger.Debug("ignoring ignoring command with incorrect parameter count", zap.String("command", commandTrigger))

		_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString(cm.App.config.Language, "CommandUsageTemplate", commandObject.Usage, commandObject.Description, commandObject.Example))
		if err != nil {
			errs = append(errs, err)
		}

		return exists, source, errs
	} else if commandObject.ParametersRange.Maximum > -1 && commandParametersCount > commandObject.ParametersRange.Maximum {
		logger.Debug("ignoring ignoring command with incorrect parameter count", zap.String("command", commandTrigger))

		_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString(cm.App.config.Language, "TooManyParameters", commandObject.ParametersRange.Maximum))
		if err != nil {
			errs = append(errs, err)
		}

		return exists, source, errs
	}

	var (
		success      bool
		enterContext bool
	)

	err = cm.App.discordClient.ChannelTyping(message.ChannelID)
	if err != nil {
		logger.Warn("failed to get channel info",
			zap.Error(err))
		return
	}

	// Execute the command.
	success, enterContext, err = commandObject.Function(cm, commandArgument, message, false)
	errs = append(errs, err)
	if enterContext {
		if commandObject.Context {
			cm.Contexts.Set(message.Author.ID, commandObject, gocache.DefaultExpiration)
		}
	}
	if !success {
		if commandObject.ErrorMessage != "" {
			// Format it if we have a mention in the error message.
			if strings.Contains(commandObject.ErrorMessage, "<@%s>") {
				_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, fmt.Sprintf(commandObject.ErrorMessage, message.Author.ID))
			} else {
				_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, commandObject.ErrorMessage)
			}

			errs = append(errs, err)
		} else {
			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString(cm.App.config.Language, "CommandUsageTemplate", commandObject.Usage, commandObject.Description, commandObject.Example))
			if err != nil {
				errs = append(errs, err)
			}

			return exists, source, errs
		}
	}

	return exists, source, errs
}

// ProcessContext re-runs a Command function if the user is currently in a
// Command's context.
func (cm CommandManager) ProcessContext(command Command, cmdtext string, message discordgo.Message) (continueContext bool, errs []error) {
	_, continueContext, err := command.Function(cm, cmdtext, message, true)
	if err != nil {
		errs = append(errs, err)
	}
	return continueContext, errs
}

func (cm CommandManager) getCommandSource(message discordgo.Message) (CommandSource, error) {
	if message.ChannelID == cm.App.config.AdministrativeChannel {
		return CommandSourceADMINISTRATIVE, nil
	} else if message.ChannelID == cm.App.config.PrimaryChannel {
		return CommandSourcePRIMARY, nil
	} else {
		ch, err := cm.App.discordClient.Channel(message.ChannelID)
		if err != nil {
			return CommandSourceOTHER, err
		}

		if ch.Type == discordgo.ChannelTypeDM {
			return CommandSourcePRIVATE, nil
		}
	}

	return CommandSourceOTHER, nil
}
