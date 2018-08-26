package bot

import (
	"fmt"

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

const greeting = `Hi! Welcome to the San Andreas Multiplayer unofficial Discord server!

You can verify your forum account by typing %s below, this helps us ensure people aren't impersonating others.`

// ConnectDiscord sets up the Discord API and event listeners
func (app *App) ConnectDiscord() (err error) {
	app.discordClient, err = discordgo.New("Bot " + app.config.DiscordToken)
	if err != nil {
		return
	}

	app.discordClient.AddHandler(app.onReady)
	app.discordClient.AddHandler(app.onMessage)
	app.discordClient.AddHandler(app.onJoin)

	err = app.discordClient.Open()
	if err != nil {
		return
	}

	<-app.ready
	return
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
			verified, err := app.storage.IsUserVerified(user.User.ID)
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

	for _, ex := range app.extensions {
		e := ex.OnMessage(*event.Message)
		if e != nil {
			logger.Error("unhandled error from OnMessage", zap.Error(e))
		}
	}

	err := app.storage.RecordChatLog(event.Message.Author.ID, event.Message.ChannelID, event.Message.Content)
	if err != nil {
		logger.Error("failed to record chat log", zap.Error(err))
	}
}

func (app *App) onJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	verified, err := app.storage.IsUserVerified(event.Member.User.ID)
	if err != nil {
		logger.Error("failed to check if user verified", zap.Error(err))
	}

	if verified {
		err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, event.Member.User.ID, app.config.VerifiedRole)
		if err != nil {
			logger.Error("failed to add verified role to member", zap.Error(err))
		}
	} else {
		ch, err := s.UserChannelCreate(event.Member.User.ID)
		if err != nil {
			logger.Error("failed to create user channel", zap.Error(err))
			return
		}
		_, err = app.discordClient.ChannelMessageSend(ch.ID, fmt.Sprintf(greeting, "`verify`"))
		if err != nil {
			logger.Error("failed to send message", zap.Error(err))
		}
	}
}
