package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// CommandManager stores command state
type CommandManager struct {
	Config    *types.Config
	Discord   *discordgo.Session
	Storage   *storage.API
	Forum     *forum.ForumClient
	Commands  map[string]Command
	Contexts  *cache.Cache
	Cooldowns map[string]time.Time
	Cache     *cache.Cache
}

// Init creates a command manager for the app
func (cm *CommandManager) Init(
	config *types.Config,
	discord *discordgo.Session,
	api *storage.API,
	fc *forum.ForumClient,
) (err error) {
	cm.Config = config
	cm.Storage = api
	cm.Discord = discord
	cm.Forum = fc

	cm.Contexts = cache.New(5*time.Minute, 30*time.Second)
	cm.Cooldowns = make(map[string]time.Time)
	cm.Cache = cache.New(5*time.Minute, 30*time.Second)

	cm.LoadCommands()

	return nil
}

// Command represents a public, private or administrative command
type Command struct {
	Function    func(args string, message discordgo.Message, contextual bool) (context bool, err error)
	Source      CommandSource
	Description string
	Context     bool
	Cooldown    time.Duration
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

// CommandParametersRange represents minimum value and maximum value number of parameters for a command
type CommandParametersRange struct {
	Minimum int
	Maximum int
}

func (cm *CommandManager) commandCommands(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	embed := &discordgo.MessageEmbed{
		Color: 0x3498DB,
	}

	var cmdlist string
	for trigger, cmd := range cm.Commands {
		cmdlist += fmt.Sprintf("**%s** - %s\n", trigger, cmd.Description)
	}

	embed.Description = cmdlist
	_, err = cm.Discord.ChannelMessageSendEmbed(message.ChannelID, embed)
	return
}

func (cm *CommandManager) commandHelp(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	_, err = cm.Discord.ChannelMessageSend(message.ChannelID, "fuck off")
	return
}

// OnMessage is called on a command string to check whether it's a valid command
// and, if so, call the associated function.
// nolint:gocyclo
func (cm *CommandManager) OnMessage(message discordgo.Message) (err error) {
	source, err := cm.getCommandSource(message)
	if err != nil {
		return
	}

	storedContext, found := cm.Contexts.Get(message.Author.ID)
	if found {
		fmt.Println("found context for", message.Author.Username)
		contextCommand, ok := storedContext.(Command)
		if !ok {
			return errors.New("failed to cast stored context to command type")
		}
		if contextCommand.Source == source {
			var continueContext bool
			continueContext, err = contextCommand.Function(message.Content, message, true)
			if err != nil {
				cm.Contexts.Delete(message.Author.ID)
				return
			}
			if !continueContext {
				cm.Contexts.Delete(message.Author.ID)
			}
			return
		}
	}

	commandAndParameters := strings.SplitN(message.Content, " ", 2)
	commandTrigger := strings.ToLower(commandAndParameters[0])
	commandArgument := ""

	if len(commandAndParameters) > 1 {
		commandArgument = commandAndParameters[1]
	}

	commandObject, exists := cm.Commands[commandTrigger]
	if !exists {
		return
	}

	if source != commandObject.Source {
		return
	}

	switch source {
	case CommandSourceADMINISTRATIVE:
		if message.ChannelID != cm.Config.AdministrativeChannel {
			return
		}
	case CommandSourcePRIMARY:
		if message.ChannelID != cm.Config.PrimaryChannel {
			return
		}
	}

	// Check if command is on cooldown
	if when, ok := cm.Cooldowns[commandTrigger]; ok {
		if time.Since(when) < commandObject.Cooldown {
			return
		}
	}

	err = cm.Discord.ChannelTyping(message.ChannelID)
	if err != nil {
		return
	}

	enterContext, err := commandObject.Function(commandArgument, message, false)
	if err != nil {
		_, err = cm.Discord.ChannelMessageSend(message.ChannelID, commandObject.Description)
		return
	}

	if enterContext {
		if commandObject.Context {
			cm.Contexts.Set(message.Author.ID, commandObject, cache.DefaultExpiration)
		}
	}

	if commandObject.Cooldown > 0 {
		cm.Cooldowns[commandTrigger] = time.Now()
	}

	return nil
}

func (cm *CommandManager) getCommandSource(message discordgo.Message) (CommandSource, error) {
	if message.ChannelID == cm.Config.AdministrativeChannel {
		return CommandSourceADMINISTRATIVE, nil
	} else if message.ChannelID == cm.Config.PrimaryChannel {
		return CommandSourcePRIMARY, nil
	} else {
		ch, err := cm.Discord.Channel(message.ChannelID)
		if err != nil {
			return CommandSourceOTHER, err
		}

		if ch.Type == discordgo.ChannelTypeDM {
			return CommandSourcePRIVATE, nil
		}
	}

	return CommandSourceOTHER, nil
}
