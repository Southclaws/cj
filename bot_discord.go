package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// ChannelDM is a direct message channel
type ChannelDM struct {
	ChannelID     string         `json:"id"`
	Private       bool           `json:"is_private"`
	Recipient     discordgo.User `json:"recipient"`
	LastMessageID string         `json:"last_message_id"`
}

// ConnectDiscord sets up the Discord API and event listeners
func (app *App) ConnectDiscord() {
	var err error

	app.discordClient, err = discordgo.New("Bot " + app.config.DiscordToken)
	if err != nil {
		logger.Fatal("failed to connect to Discord API",
			zap.Error(err))
	}

	app.discordClient.AddHandler(app.onReady)
	app.discordClient.AddHandler(app.onMessage)
	app.discordClient.AddHandler(app.onJoin)

	err = app.discordClient.Open()
	if err != nil {
		logger.Fatal("failed to start Discord client",
			zap.Error(err))
	}
}

// nolint:gocyclo
func (app *App) onReady(s *discordgo.Session, event *discordgo.Ready) {
	logger.Debug("connected to Discord gateway")

	roles, err := s.GuildRoles(app.config.GuildID)
	if err != nil {
		logger.Fatal("failed to get guild roles",
			zap.Error(err))
	}

	found := false
	for _, role := range roles {
		if role.ID == app.config.VerifiedRole {
			found = true
			break
		}
	}
	if !found {
		logger.Fatal("verified role not found.",
			zap.String("role", app.config.VerifiedRole))
	}

	users, err := s.GuildMembers(app.config.GuildID, "", 1000)
	if err != nil {
		logger.Fatal("failed to get guild members",
			zap.Error(err))
	}

	if !app.config.NoInitSync {
		for _, user := range users {
			verified, err := app.IsUserVerified(user.User.ID)
			if err != nil {
				logger.Fatal("failed to check user verified state",
					zap.Error(err),
					zap.String("user", user.User.ID))
			}
			if verified {
				logger.Info("synchronising roles by adding verified status to user", zap.String("user", user.User.Username))
				err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, user.User.ID, app.config.VerifiedRole)
				if err != nil {
					logger.Fatal("failed to add verified role",
						zap.Error(err),
						zap.String("user", user.User.ID))
				}
			} else {
				logger.Info("synchronising roles by removing verified status from user", zap.String("user", user.User.Username))
				err = app.discordClient.GuildMemberRoleRemove(app.config.GuildID, user.User.ID, app.config.VerifiedRole)
				if err != nil {
					logger.Fatal("failed to remove verified role",
						zap.Error(err),
						zap.String("user", user.User.ID))
				}
			}
		}
	}

	ticker := time.NewTicker(time.Minute * time.Duration(app.config.Heartbeat))
	for t := range ticker.C {
		app.onHeartbeat(t)
	}

	logger.Info("finished initialising discord module")

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
			logger.Debug("ignoring command from non debug user")
			return
		}
		logger.Debug("accepting command from debug user")
	}

	_, source, errors := app.commandManager.Process(*event.Message)
	for _, err := range errors {
		if err != nil {
			err = app.WarnUserError(event.Message.ChannelID, err.Error())
			if err != nil {
				logger.Warn("failed to warn user of error", zap.Error(err))
			}
		}
	}

	if source != CommandSourcePRIVATE && source != CommandSourceADMINISTRATIVE {
		err := app.RecordChatLog(event.Message.Author.ID, event.Message.ChannelID, event.Message.Content)
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
		_, err = app.discordClient.ChannelMessageSend(ch.ID,
			"Hi! Welcome to the San Andreas Multiplayer unofficial Discord server!\n\nYou can verify your forum account by typing `verify` below, this helps us ensure people aren't impersonating others.")
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
