package bot

import (
	"fmt"
	"strings"

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

Please read the rules and be respectful.`

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
	app.discordClient.S.AddHandler(app.onReactionAdd)
	app.discordClient.S.AddHandler(app.onReactionRemove)

	intent := discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)

	app.discordClient.S.Identify.Intents = intent

	err = app.discordClient.S.Open()
	if err != nil {
		return
	}

	err = <-app.ready
	if err != nil {
		return
	}

	zap.L().Info("finished initialising discord module")
	app.discordClient.S.UpdateGameStatus(0, fmt.Sprintf("🐧 R.I.P. - %s", app.config.Version))

	return
}

// nolint:gocyclo
func (app *App) onReady(s *discordgo.Session, event *discordgo.Ready) {
	zap.L().Debug("connected to Discord gateway")

	app.ready <- func() error {
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

	if event.Message.Author.Bot == true {
		return
	}

	if app.config.DebugUser != "" {
		if event.Message.Author.ID != app.config.DebugUser {
			zap.L().Debug("ignoring command from non debug user")
			return
		}
		zap.L().Debug("accepting command from debug user")
	}

	for _, ex := range app.extensions {
		e := ex.OnMessage(*event.Message)
		if e != nil {
			zap.L().Error("unhandled error from OnMessage", zap.Error(e))
		}
	}

	err = app.storage.RecordChatLog(event.Message.Author.ID,
		event.Message.ChannelID,
		event.Message.Content,
		event.Message.ID)
	if err != nil {
		zap.L().Error("failed to record chat log", zap.Error(err))
	}

	zap.L().Debug("processed message",
		zap.String("author", event.Message.Author.Username),
		zap.String("message", event.Message.Content),
	)
}

func (app *App) onReactionAdd(s *discordgo.Session, event *discordgo.MessageReactionAdd) {
	message, err := app.storage.GetMessageByID(event.MessageID)
	if err != nil || message.DiscordUserID == "" || message.DiscordUserID == event.UserID {
		// Message likely just doesn't exists in the DB, or is from the user themselves
		return
	}

	emoji := event.Emoji.APIName()
	if strings.Index(emoji, ":") != -1 {
		emoji = fmt.Sprintf("<:%s>", emoji)
	}
	err = app.storage.AddEmojiReactionToUser(message.DiscordUserID, emoji)
	if err != nil {
		zap.L().Debug("Error: ", zap.Error(err))
	}
}

func (app *App) onReactionRemove(s *discordgo.Session, event *discordgo.MessageReactionRemove) {
	message, err := app.storage.GetMessageByID(event.MessageID)
	if err != nil || message.DiscordUserID == "" || message.DiscordUserID == event.UserID {
		// Message likely just doesn't exists in the DB, or is from the user themselves
		return
	}

	emoji := event.Emoji.APIName()
	if strings.Index(emoji, ":") != -1 {
		emoji = fmt.Sprintf("<:%s>", emoji)
	}
	err = app.storage.RemoveEmojiReactionFromUser(message.DiscordUserID, emoji)
	if err != nil {
		zap.L().Debug("Error: ", zap.Error(err))
	}
}

func (app *App) onJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	ch, err := s.UserChannelCreate(event.Member.User.ID)
	if err != nil {
		zap.L().Error("failed to create user channel", zap.Error(err))
		return
	}
	_, err = app.discordClient.S.ChannelMessageSend(ch.ID, greeting)
	if err != nil {
		zap.L().Error("failed to send message", zap.Error(err))
	}
}
