package commands

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/bot/chatgpt"
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
	Function         func(interaction *discordgo.InteractionCreate, args map[string]*discordgo.ApplicationCommandInteractionDataOption, settings types.CommandSettings) (context bool, err error)
	Name             string
	Description      string
	Settings         types.CommandSettings
	Options          []*discordgo.ApplicationCommandOption
	IsAdministrative bool
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
	content := strings.ToLower(message.Content)
	
	// Check for "cj is this true" pattern first (higher priority)
	if strings.Contains(content, "cj") && (strings.Contains(content, "is this true") || strings.Contains(content, "is this real")) {
		cm.handleIsThisTrueMessage(&message)
		return nil
	}
	
	// Original CJ quote functionality (only if not the "is this true" pattern)
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
	zap.L().Info("User attempting command", zap.Any("user", interaction.Member.User.ID), zap.Any("command", interaction.ApplicationCommandData().Name))
	for _, command := range cm.Commands {
		if strings.TrimLeft(command.Name, "/") == interaction.ApplicationCommandData().Name {
			if hasPermissions(command.Settings.Roles, interaction.Member.Roles) {
				args := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
				for _, option := range interaction.ApplicationCommandData().Options {
					args[option.Name] = option
				}
				command.Function(interaction, args, command.Settings)
			} else {
				cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You're not authorized for this!",
					},
				})
				time.Sleep(time.Second * 5)
				cm.Discord.S.InteractionResponseDelete(interaction.Interaction)
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
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func (cm *CommandManager) replyDirectlyEmbed(interaction *discordgo.InteractionCreate, response string, embed *discordgo.MessageEmbed) {
	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Embeds:  []*discordgo.MessageEmbed{embed},
		},
	})
}

// sendThinkingResponse is used to avoid the 3 seconds timeout limit
func (cm *CommandManager) sendThinkingResponse(interaction *discordgo.InteractionCreate) {
	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
		},
	})
}

func (cm *CommandManager) editOriginalResponse(interaction *discordgo.InteractionCreate, response string) {
	cm.Discord.S.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &response,
	})
}

func (cm *CommandManager) editOriginalResponseWithEmbed(interaction *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	content := ""
	embeds := []*discordgo.MessageEmbed{embed}
	cm.Discord.S.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Embeds:  &embeds,
	})
}

// handleIsThisTrueMessage deals with "cj is this true" replies
func (cm *CommandManager) handleIsThisTrueMessage(message *discordgo.Message) {
	// 60 second cooldown
	cooldownKey := "analysis"
	if lastTime, exists := cm.Cooldowns[cooldownKey]; exists {
		if lastTime.Add(time.Second * 60).After(time.Now()) {
			cm.Discord.S.MessageReactionAdd(message.ChannelID, message.ID, "🕛")
			return
		}
	}
	cm.Cooldowns[cooldownKey] = time.Now()

	if cm.Config.ChatGPTToken == "" {
		responses := []string{
			"Damn, I can't do my analysis shit right now. Tell whoever runs this shit to fix it!",
			"Yo, my brain is broken right now. Get the admin bitch to fix this shit!",
			"Man, I can't check nothing right now. Tell that asshole to fix this!",
			"This shit ain't working. Get whoever runs this shit to fix it!",
		}
		rand.Seed(time.Now().UnixNano())
		cm.Discord.ChannelMessageSend(message.ChannelID, responses[rand.Intn(len(responses))])
		return
	}

	if message.MessageReference == nil || message.MessageReference.MessageID == "" {
		responses := []string{
			"Yo, you need to reply to a message for me to check it out, fool!",
			"Man, reply to something first so I can check it out!",
			"You gotta reply to a message for me to analyze it, dude!",
			"Reply to something first, then I'll check it out!",
		}
		rand.Seed(time.Now().UnixNano())
		cm.Discord.ChannelMessageSend(message.ChannelID, responses[rand.Intn(len(responses))])
		return
	}

	repliedMessage, err := cm.Discord.S.ChannelMessage(message.ChannelID, message.MessageReference.MessageID)
	if err != nil {
		zap.L().Error("Failed to fetch replied message", zap.Error(err))
		responses := []string{
			"Damn, I can't see what you're talking about, man!",
			"Yo, I can't find that message you're talking about!",
			"Man, I can't see what you want me to check out!",
			"This shit is confusing, I can't see the message!",
		}
		rand.Seed(time.Now().UnixNano())
		cm.Discord.ChannelMessageSend(message.ChannelID, responses[rand.Intn(len(responses))])
		return
	}

	chatGPTService := chatgpt.NewService(cm.Config.ChatGPTToken)
	if chatGPTService == nil {
		responses := []string{
			"My brain ain't working right now, man!",
			"Yo, my head is fried right now!",
			"Man, I can't think straight right now!",
			"My brain is all messed up right now!",
		}
		rand.Seed(time.Now().UnixNano())
		cm.Discord.ChannelMessageSend(message.ChannelID, responses[rand.Intn(len(responses))])
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := chatGPTService.IsThisReal(ctx, repliedMessage.Content)
	if err != nil {
		zap.L().Error("Failed to analyze message with ChatGPT", zap.Error(err))
		responses := []string{
			"Damn, my brain is fried right now! Try again later, fool!",
			"Yo, my head is all messed up! Try again later!",
			"Man, I can't think right now! Try again later!",
			"My brain is broken! Try again later, dude!",
		}
		rand.Seed(time.Now().UnixNano())
		cm.Discord.ChannelMessageSend(message.ChannelID, responses[rand.Intn(len(responses))])
		return
	}

	var cjResponse string
	if response.Confidence == "low" {
		responses := []string{
			"Man, I don't know about this shit. %s confidence.\n\n%s\n\n%s",
			"Yo, this is confusing as hell. %s confidence.\n\n%s\n\n%s",
			"I can't really tell, man. %s confidence.\n\n%s\n\n%s",
			"This shit is tricky. %s confidence.\n\n%s\n\n%s",
		}
		rand.Seed(time.Now().UnixNano())
		format := responses[rand.Intn(len(responses))]
		cjResponse = fmt.Sprintf(format, response.Confidence, response.Explanation, response.CJResponse)
	} else if response.IsReal {
		responses := []string{
			"That's real, man. %s confidence.\n\n%s\n\n%s",
			"Yo, that's legit. %s confidence.\n\n%s\n\n%s",
			"That shit is true. %s confidence.\n\n%s\n\n%s",
			"Real talk, that's accurate. %s confidence.\n\n%s\n\n%s",
			"Yeah, that's the truth. %s confidence.\n\n%s\n\n%s",
		}
		rand.Seed(time.Now().UnixNano())
		format := responses[rand.Intn(len(responses))]
		cjResponse = fmt.Sprintf(format, response.Confidence, response.Explanation, response.CJResponse)
	} else {
		responses := []string{
			"That's fake as hell, man. %s confidence.\n\n%s\n\n%s",
			"Yo, that's bullshit. %s confidence.\n\n%s\n\n%s",
			"That shit is fake. %s confidence.\n\n%s\n\n%s",
			"Nah, that's not real. %s confidence.\n\n%s\n\n%s",
			"That's some made up shit. %s confidence.\n\n%s\n\n%s",
		}
		rand.Seed(time.Now().UnixNano())
		format := responses[rand.Intn(len(responses))]
		cjResponse = fmt.Sprintf(format, response.Confidence, response.Explanation, response.CJResponse)
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, cjResponse)
}
