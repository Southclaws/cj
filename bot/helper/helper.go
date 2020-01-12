package helper

import (
	"strings"
	"time"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Helper stores assistance state
type Helper struct {
	Config   *types.Config
	Discord  *discord.Session
	Storage  storage.Storer
	Forum    *forum.ForumClient
	Words    []string
	Cooldown time.Time
}

// Init initializes Helper extension
func (h *Helper) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
	fc *forum.ForumClient,
) (err error) {
	h.Config = config
	h.Discord = discord
	h.Storage = api
	h.Forum = fc
	h.Words = []string{"can anyone", "can someone", "I need help", "help me", "how do I", "how to", "I have a question"}
	return
}

// OnMessage is called every time user sends a message in a channel
func (h *Helper) OnMessage(msg discordgo.Message) (err error) {
	if msg.ChannelID != h.Config.ScriptingChannel {
		return
	}
	assistanceAvailable := true
	if time.Since(h.Cooldown) < h.Config.AssistanceCooldown {
		assistanceAvailable = false
	}

	for _, sub := range h.Words {
		if !strings.Contains(strings.ToLower(msg.Content), strings.ToLower(sub)) {
			continue
		}

		if !assistanceAvailable {
			h.Discord.S.MessageReactionAdd(msg.ChannelID, msg.ID, "❔")
			return
		}

		guild, e := h.Discord.S.Guild(msg.GuildID)
		if e != nil {
			zap.L().Debug("failed to find guild", zap.String("guild ID", msg.GuildID))
			err = e
			return
		}

		var helperRole *discordgo.Role

		for _, role := range guild.Roles {
			if role.ID == h.Config.HelperRole {
				helperRole = role
				break
			}
		}

		var (
			helpers []string
			// iteratedMember holds already iterated member list
			// it exists ONLY because there's a bug in discordgo library
			iteratedMembers = make(map[string]bool)
		)

	membersLoop:
		for _, guildMember := range guild.Members {
			// check if this guildMember is already iterated or not
			// why doing this extra check? because discordgo has a bug
			_, ok := iteratedMembers[guildMember.User.Username]
			if ok {
				continue
			} else {
				iteratedMembers[guildMember.User.Username] = true
			}

			for _, guildMemberRole := range guildMember.Roles {
				if guildMemberRole == helperRole.ID {
					eligible := false

					presence, err := h.Discord.S.State.Presence(guild.ID, guildMember.User.ID)
					if err != nil {
						zap.L().Debug("failed to get user presence", zap.String("guild ID", guild.ID), zap.String("user ID", guildMember.User.ID))
						continue membersLoop
					}

					if presence.Status == discordgo.StatusOnline {
						eligible = true
					}

					if eligible {
						helpers = append(helpers, guildMember.User.Username)
					}
					break
				}
			}
		}

		if len(helpers) == 0 {
			h.Discord.S.ChannelMessageSend(msg.ChannelID, msg.Author.Mention()+", you might be asking for assistance, but the registered helpers are not online, thus your question __may__ be unanswered")
		} else {
			h.Discord.S.ChannelMessageSend(msg.ChannelID, msg.Author.Mention()+", you might be asking for assistance, here's the list of the available helpers:\n**"+strings.Join(helpers, "**; **")+"**")
		}

		h.Cooldown = time.Now()
		break
	}
	return
}
