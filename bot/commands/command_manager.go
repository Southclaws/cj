package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// CommandManager stores command state
type CommandManager struct {
	Config    *types.Config
	Discord   *discord.Session
	Storage   storage.Storer
	Forum     *forum.ForumClient
	Commands  []Command
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
	Function    func(interaction *discordgo.InteractionCreate, args map[string]*discordgo.ApplicationCommandInteractionDataOption, settings types.CommandSettings) (context bool, err error)
	Name        string
	Description string
	Settings    types.CommandSettings
	Options     []*discordgo.ApplicationCommandOption
	Cooldown    time.Duration // DEPRECATED
}

// CommandParametersRange represents minimum value and maximum value number of parameters for a command
type CommandParametersRange struct {
	Minimum int
	Maximum int
}

// OnMessage is called on a command string to check whether it's a valid command
// and, if so, call the associated function.
// nolint:gocyclo
func (cm *CommandManager) OnMessage(message discordgo.Message) (err error) {
	split := strings.Split(message.Content, " ")
	for _, s := range split {
		if strings.ToLower(s) == "cj" {
			cm.commandCJQuote(&message)
			break
		}
	}
	return nil
}

func (cm *CommandManager) TryFindAndFireCommand(interaction *discordgo.InteractionCreate) {
	for _, command := range cm.Commands {
		if strings.TrimLeft(command.Name, "/") == interaction.Data.Name {
			if hasPermissions(command.Settings.Roles, interaction.Member.Roles) {
				args := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
				for _, option := range interaction.Data.Options {
					args[option.Name] = option
				}
				command.Function(interaction, args, command.Settings)
			} else {
				cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: "You're not authorized for this!",
					},
				})
				time.Sleep(time.Second * 5)
				cm.Discord.S.InteractionResponseDelete(cm.Discord.S.State.User.ID, interaction.Interaction)
			}
			break
		}
	}
}

func hasPermissions(commandRoles []string, memberRoles []string) bool {
	if len(commandRoles) == 0 {
		return true
	}
	for _, i := range commandRoles {
		if i == "all" {
			return true
		}
		for _, j := range memberRoles {
			if i == j {
				return true
			}
		}
	}
	return false
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

func (cm *CommandManager) replyDirectly(interaction *discordgo.InteractionCreate, response string) {
	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: response,
		},
	})
}

func (cm *CommandManager) replyDirectlyEmbed(interaction *discordgo.InteractionCreate, response string, embed *discordgo.MessageEmbed) {
	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: response,
			Embeds:  []*discordgo.MessageEmbed{embed},
		},
	})
}
