package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
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

	err = <-app.ready
	if err != nil {
		return
	}

	logger.Info("finished initialising discord module")

	return
}

// nolint:gocyclo
func (app *App) onReady(s *discordgo.Session, event *discordgo.Ready) {
	logger.Debug("connected to Discord gateway")

	app.ready <- func() error {
		roles, err := s.GuildRoles(app.config.GuildID)
		if err != nil {
			return errors.Wrap(err, "failed to get guild roles")
		}

		found := false
		for _, role := range roles {
			if role.ID == app.config.VerifiedRole {
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("verified role %s not found.", app.config.VerifiedRole)
		}

		if !app.config.NoInitSync {
			go app.doSync()
		}

		return nil
	}()
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
		} else {
			logger.Debug("handled command in extension",
				zap.String("author", event.Message.Author.Username),
				zap.String("message", event.Message.Content),
			)
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

func (app *App) doSync() {
	err := func() (err error) {
		users, err := app.discordClient.GuildMembers(app.config.GuildID, "", 1000)
		if err != nil {
			return errors.Wrap(err, "failed to get guild members")
		}

		for _, user := range users {
			verified, err := app.storage.IsUserVerified(user.User.ID)
			if err != nil {
				return errors.Wrapf(err, "failed to check user verified state for %s", user.User.ID)
			}

			if verified {
				logger.Debug("synchronising roles by adding verified status to user", zap.String("user", user.User.Username))
				err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, user.User.ID, app.config.VerifiedRole)
				if err != nil {
					return errors.Wrapf(err, "failed to add verified role for %s", user.User.ID)
				}
			} else {
				logger.Debug("synchronising roles by removing verified status from user", zap.String("user", user.User.Username))
				err = app.discordClient.GuildMemberRoleRemove(app.config.GuildID, user.User.ID, app.config.VerifiedRole)
				if err != nil {
					return errors.Wrapf(err, "failed to remove verified role for %s", user.User.ID)
				}
			}
		}
		return
	}()
	if err != nil {
		logger.Fatal("failed to perform initialisation sync",
			zap.Error(err))
	}
}
