package bot

import (
	"fmt"

	"github.com/Southclaws/cj/discord"
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
	s, err := discordgo.New("Bot " + app.config.DiscordToken)
	if err != nil {
		return
	}
	app.discordClient = discord.New(s, *app.config)

	app.discordClient.S.AddHandler(app.onReady)
	app.discordClient.S.AddHandler(app.onMessage)
	app.discordClient.S.AddHandler(app.onJoin)

	err = app.discordClient.S.Open()
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
			zap.L().Error("verified role not found",
				zap.String("id", app.config.VerifiedRole),
				zap.Any("roles", roles))
			return errors.New("verified role not found")
		}

		if !app.config.NoInitSync {
			go app.doSync()
		}

		channels, err := app.discordClient.S.GuildChannels(app.config.GuildID)
		if err != nil {
			return errors.Wrap(err, "failed to get channels")
		}
		for _, ch := range channels {
			app.channels[ch.ID] = ch
		}

		return nil
	}()
}

// nolint:gocyclo
func (app *App) onMessage(s *discordgo.Session, event *discordgo.MessageCreate) {
	if len(app.ready) > 0 {
		<-app.ready
	}

	ch, err := app.discordClient.S.Channel(event.ChannelID)
	if err != nil {
		return
	}
	if ch.Type == discordgo.ChannelTypeGuildText {
		if ch.GuildID != app.config.GuildID {
			return
		}
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

	err = app.storage.RecordChatLog(event.Message.Author.ID, event.Message.ChannelID, event.Message.Content)
	if err != nil {
		logger.Error("failed to record chat log", zap.Error(err))
	}

	logger.Debug("processed message",
		zap.String("author", event.Message.Author.Username),
		zap.String("message", event.Message.Content),
	)
}

func (app *App) onJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	verified, err := app.storage.IsUserVerified(event.Member.User.ID)
	if err != nil {
		logger.Error("failed to check if user verified", zap.Error(err))
	}

	if verified {
		err = app.discordClient.S.GuildMemberRoleAdd(app.config.GuildID, event.Member.User.ID, app.config.VerifiedRole)
		if err != nil {
			logger.Error("failed to add verified role to member", zap.Error(err))
		}
	} else {
		ch, err := s.UserChannelCreate(event.Member.User.ID)
		if err != nil {
			logger.Error("failed to create user channel", zap.Error(err))
			return
		}
		_, err = app.discordClient.S.ChannelMessageSend(ch.ID, fmt.Sprintf(greeting, "`verify`"))
		if err != nil {
			logger.Error("failed to send message", zap.Error(err))
		}
	}
}

func (app *App) doSync() {
	err := func() (err error) {
		users, err := app.discordClient.S.GuildMembers(app.config.GuildID, "", 1000)
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
				err = app.discordClient.S.GuildMemberRoleAdd(app.config.GuildID, user.User.ID, app.config.VerifiedRole)
				if err != nil {
					return errors.Wrapf(err, "failed to add verified role for %s", user.User.ID)
				}
			} else {
				logger.Debug("synchronising roles by removing verified status from user", zap.String("user", user.User.Username))
				err = app.discordClient.S.GuildMemberRoleRemove(app.config.GuildID, user.User.ID, app.config.VerifiedRole)
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
