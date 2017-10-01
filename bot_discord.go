package main

import (
	"time"

	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

// ChannelDM is a direct message channel
type ChannelDM struct {
	ChannelID     string         `json:"id"`
	Private       bool           `json:"is_private"`
	Recipient     discordgo.User `json:"recipient"`
	LastMessageID string         `json:"last_message_id"`
}

func (app *App) connect() error {
	var err error

	app.discordClient, err = discordgo.New("Bot " + app.config.DiscordToken)
	if err != nil {
		panic(err)
	}

	app.discordClient.AddHandler(app.onReady)
	app.discordClient.AddHandler(app.onMessage)
	app.discordClient.AddHandler(app.onJoin)

	err = app.discordClient.Open()
	if err != nil {
		panic(err)
	}

	return nil
}

// nolint:gocyclo
func (app *App) onReady(s *discordgo.Session, event *discordgo.Ready) {
	found := 0
	roles, err := s.GuildRoles(app.config.GuildID)
	if err != nil {
		panic(err)
	}

	for _, role := range roles {
		if role.ID == app.config.VerifiedRole || role.ID == app.config.NormalRole {
			found++
		}
	}
	if found != 2 {
		logger.Fatal("role not found.",
			zap.String("role", app.config.VerifiedRole))
	}

	var member bool
	users, err := s.GuildMembers(app.config.GuildID, "", 1000)
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		member = false
		for i := range user.Roles {
			if user.Roles[i] == app.config.NormalRole {
				member = true
			}
		}
		if !member {
			err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, user.User.ID, app.config.NormalRole)
			if err != nil {
				panic(err)
			}
		}
	}

	ticker := time.NewTicker(time.Minute * time.Duration(app.config.Heartbeat))
	for t := range ticker.C {
		app.onHeartbeat(t)
	}

	app.ready <- true
}

// nolint:gocyclo
func (app *App) onMessage(s *discordgo.Session, event *discordgo.MessageCreate) {
	if len(app.ready) > 0 {
		<-app.ready
	}

	if event.Message.Author.ID == app.config.BotID {
		return
	}

	if app.config.DebugUser != "" {
		if event.Message.Author.ID != app.config.DebugUser {
			return
		}
	}

	_, source, errors := app.commandManager.Process(*event.Message)
	if errors != nil {
		for _, err := range errors {
			if err != nil {
				err = app.WarnUserError(event.Message.ChannelID, err.Error())
				if err != nil {
					logger.Warn("failed to warn user of error", zap.Error(err))
				}
			}
		}
	}

	if source != CommandSourcePRIVATE && source != CommandSourceADMINISTRATIVE {
		err := app.chatLogger.RecordChatLog(event.Message.Author.ID, event.Message.ChannelID, event.Message.Content)
		if err != nil {
			logger.Warn("failed to record chat log", zap.Error(err))
		}

		for i := range event.Message.Mentions {
			if event.Message.Mentions[i].ID == app.config.BotID {
				err := app.HandleSummon(*event.Message)
				if err != nil {
					logger.Warn("failed to handle summon", zap.Error(err))
				}
			}
		}
	}
}

func (app *App) onJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	verified, err := app.IsUserVerified(event.Member.User.ID)
	if err != nil {
		logger.Warn("failed to check if user verified", zap.Error(err))
	}

	err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, event.Member.User.ID, app.config.NormalRole)
	if err != nil {
		logger.Warn("failed to add normal role to member", zap.Error(err))
	}

	if verified {
		err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, event.Member.User.ID, app.config.VerifiedRole)
		if err != nil {
			logger.Warn("failed to add verified role to member", zap.Error(err))
		}
	} else {
		ch, err := s.UserChannelCreate(event.Member.User.ID)
		if err != nil {
			logger.Warn("failed to create user channel", zap.Error(err))
			return
		}
		_, err = app.discordClient.ChannelMessageSend(ch.ID, app.locale.GetLangString("en", "AskUserVerify"))
		if err != nil {
			logger.Warn("failed to send message", zap.Error(err))
		}
	}
}

func (app *App) onHeartbeat(t time.Time) {
	err := app.HandleHeartbeatEvent(t)
	if err != nil {
		logger.Warn("failed to handle heartbeat", zap.Error(err))
	}
}
