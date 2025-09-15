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

/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
// Define the roles that we want to track
var trackedRoleIDs = map[string]string{
	"1002922725553217648": 	"Clown",
	"400250542628274177":  	"Caged",
	"995816487610753094":  	"annoyed me",
	"1016047260364198008": 	"Not Cool",
	"833325019252785173":  	"Doesnt deserve to embed",
	"818457955690872832":  	"Doesnt deserve to react",
	"996883259252297758": 	"No open.mp support",
	"910950457680212088": 	"No Server Adverts",
	"841368374356738078": 	"Suffers from dunning-kruger",
	"987825514511220867": 	"Muted",
	"1204891485867352144": 	"Can't @everyone",
}
*/

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
	/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
	app.discordClient.S.AddHandler(app.onGuildMemberUpdate)
	*/

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
	app.discordClient.S.UpdateGameStatus(0, fmt.Sprintf("Open Multiplayer | 🐧 R.I.P. - %s", app.config.Version))

	return
}

/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
// nolint:gocyclo
func (app *App) onGuildMemberUpdate(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	member := event.Member
	
	// Get user's current tracked roles
	user, err := app.storage.GetUserOrCreate(member.User.ID)
	if err != nil {
		zap.L().Error("failed to get user for role tracking", zap.Error(err))
		return
	}

	// Check current Discord roles aginst our tracked role list
	currentTrackedRoles := make(map[string]bool)
	for _, roleID := range member.Roles {
		if roleName, isTracked := trackedRoleIDs[roleID]; isTracked {
			currentTrackedRoles[roleID] = true
			
			hasRole := false
			for _, trackedRole := range user.TrackedRoles {
				if trackedRole.RoleID == roleID {
					hasRole = true
					break
				}
			}
			if !hasRole {
				err := app.storage.AddTrackedRole(member.User.ID, roleID, roleName)
				if err != nil {
					zap.L().Error("failed to add tracked role", zap.Error(err))
				} else {
					zap.L().Info("started tracking role", 
						zap.String("user", member.User.Username),
						zap.String("role", roleName))
				}
			}
		}
	}

	// Remove from database if no longer has the role
	for _, trackedRole := range user.TrackedRoles {
		if !currentTrackedRoles[trackedRole.RoleID] {
			err := app.storage.RemoveTrackedRole(member.User.ID, trackedRole.RoleID)
			if err != nil {
				zap.L().Error("failed to remove tracked role", zap.Error(err))
			} else {
				zap.L().Info("stopped tracking role", 
					zap.String("user", member.User.Username),
					zap.String("role", trackedRole.RoleName))
			}
		}
	}
}
*/

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

	if event.Message.Author.Bot {
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
	if strings.Contains(emoji, ":") {
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
	if strings.Contains(emoji, ":") {
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
	} else {
		_, err = app.discordClient.S.ChannelMessageSend(ch.ID, greeting)
		if err != nil {
			zap.L().Error("failed to send message", zap.Error(err))
		}
	}

	/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
	err = app.reapplyTrackedRoles(event.Member.User.ID, event.GuildID)
	if err != nil {
		zap.L().Error("failed to reapply tracked roles", zap.Error(err))
	}
	*/
}

/* ROLE TRACKING: Commented out atm and needs to be reviewed after converting the database
func (app *App) reapplyTrackedRoles(userID, guildID string) error {
	trackedRoles, err := app.storage.GetTrackedRoles(userID)
	if err != nil {
		return err
	}

	if len(trackedRoles) == 0 {
		return nil
	}

	for _, role := range trackedRoles {
		err = app.discordClient.S.GuildMemberRoleAdd(guildID, userID, role.RoleID)
		if err != nil {
			zap.L().Error("failed to re-add tracked role", 
				zap.Error(err),
				zap.String("user", userID),
				zap.String("role", role.RoleName))
		} else {
			zap.L().Info("reapplied tracked role", 
				zap.String("user", userID),
				zap.String("role", role.RoleName))
		}
	}

	return nil
}
*/