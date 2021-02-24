package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// CommandManager stores command state
type CommandManager struct {
	Config    *types.Config
	Discord   *discord.Session
	Storage   storage.Storer
	Forum     *forum.ForumClient
	Commands  map[string]Command
	Contexts  *cache.Cache
	Cooldowns map[string]time.Time
	Cache     *cache.Cache
}

// Init creates a command manager for the app
func (cm *CommandManager) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
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
	Function    func(args string, message discordgo.Message, contextual bool, settings types.CommandSettings) (context bool, err error)
	Description string
	Settings    types.CommandSettings

	Cooldown time.Duration // DEPRECATED
}

// CommandParametersRange represents minimum value and maximum value number of parameters for a command
type CommandParametersRange struct {
	Minimum int
	Maximum int
}

func (cm *CommandManager) commandCommands(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
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

	_, err = cm.Discord.S.ChannelMessageSendEmbed(message.ChannelID, embed)
	return
}

func (cm *CommandManager) commandHelp(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.Discord.ChannelMessageSend(message.ChannelID, "fuck off")
	return
}

// OnMessage is called on a command string to check whether it's a valid command
// and, if so, call the associated function.
// nolint:gocyclo
func (cm *CommandManager) OnMessage(message discordgo.Message) (err error) {
	storedContext, found := cm.Contexts.Get(message.Author.ID)
	if found {
		_, ok := storedContext.(Command)
		if !ok {
			return errors.New("failed to cast stored context to command type")
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
		if strings.HasPrefix(commandTrigger, "/") {
			threshold := 2
			result := ""
			for word := range cm.Commands {
				dist := levenshtein.DistanceForStrings(
					[]rune(commandTrigger),
					[]rune(word),
					levenshtein.DefaultOptions)

				if dist < threshold {
					threshold = dist
					result = word
				}
			}
			if result != "" {
				cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Did you mean %s?", result))
			}
		}
		return
	}

	settings, found, err := cm.Storage.GetCommandSettings(commandTrigger)
	if err != nil {
		return
	}
	if !found {
		settings.Cooldown = cm.Config.DefaultCooldown
		settings.Channels = []string{cm.Config.DefaultChannel}
		settings.Roles = []string{cm.Config.DefaultRole}
	}

	var allowed bool
    for _, ch := range settings.Channels {
        if ch == "all" {
            allowed = true
            break
        }
        if ch == message.ChannelID {
            allowed = true
            break
        }
    }
    for _, sr := range settings.Roles {
        if sr == "all" {
            allowed = true
            break
        }
        var u *discordgo.Member
        u, err = cm.Discord.S.GuildMember(cm.Config.GuildID, message.Author.ID)
        if err != nil {
            return
        }
        for _, ur := range u.Roles {
            if sr == ur {
                allowed = true
                break
            }
        }
    }
	if !allowed {
		zap.L().Debug("command not allowed with current channel or role",
			zap.String("channel", message.ChannelID),
			zap.Any("config", settings))
		return
	}

	// Check if command is on cooldown
	if when, ok := cm.Cooldowns[commandTrigger]; ok {
		since := time.Since(when)
		if since < settings.Cooldown {
			err = cm.Discord.S.MessageReactionAdd(message.ChannelID, message.ID, pcd(since, settings.Cooldown))
			return
		}
	}

	err = cm.Discord.S.ChannelTyping(message.ChannelID)
	if err != nil {
		return
	}

	enterContext, err := commandObject.Function(commandArgument, message, false, settings)
	if err != nil {
		cm.Discord.ChannelMessageSend(message.ChannelID, err.Error())
		cm.Discord.ChannelMessageSend(message.ChannelID, commandObject.Description)
		return
	}

	if enterContext {
		commandObject.Settings = settings
		cm.Contexts.Set(message.Author.ID, commandObject, cache.DefaultExpiration)
	}

	if commandObject.Cooldown > 0 {
		cm.Cooldowns[commandTrigger] = time.Now()
	}

	return nil
}


var clocks = []string{
	"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛",
}

func pcd(since time.Duration, cooldown time.Duration) (result string) {
	p := (since.Seconds() / cooldown.Seconds()) * 100.0
	step := 100.0 / float64(len(clocks))
	for i := range clocks {
		if p <= float64(i+1)*step {
			return clocks[i]
		}
	}
	return ""
}
