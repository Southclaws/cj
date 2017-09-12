package main

import (
	"log"
	"time"

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
		log.Print("discord client creation error")
		log.Fatal(err)
	}
	debug("connected to Discord")

	app.discordClient.AddHandler(app.onReady)
	app.discordClient.AddHandler(app.onMessage)
	app.discordClient.AddHandler(app.onJoin)
	app.discordClient.AddHandler(app.onLeave)

	err = app.discordClient.Open()
	if err != nil {
		log.Println("discord client connection error")
		log.Fatal(err)
	}

	debug("awaiting Discord ready state...")

	return nil
}

// nolint:gocyclo
func (app *App) onReady(s *discordgo.Session, event *discordgo.Ready) {
	debug("discord ready")

	found := 0
	roles, err := s.GuildRoles(app.config.GuildID)
	if err != nil {
		log.Fatal(err)
	}

	for _, role := range roles {
		if role.ID == app.config.VerifiedRole || role.ID == app.config.NormalRole {
			found++
		}
	}
	if found != 2 {
		log.Printf("verified role ID '%s' was not found in guild role list:", app.config.VerifiedRole)
		for _, role := range roles {
			log.Printf("name: %s id: %s", role.Name, role.ID)
		}
		log.Fatalf("role '%s' not found.", app.config.VerifiedRole)
	}

	log.Print("Updating users to normal role")
	var member bool
	users, err := s.GuildMembers(app.config.GuildID, "", 1000)
	if err != nil {
		log.Print(err)
	}
	for _, user := range users {
		member = false
		for i := range user.Roles {
			if user.Roles[i] == app.config.NormalRole {
				member = true
			}
		}
		if !member {
			log.Printf("GuildMemberRoleAdd '%s' '%s' '%s'", app.config.GuildID, user.User.ID, app.config.NormalRole)
			err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, user.User.ID, app.config.NormalRole)
			if err != nil {
				log.Print(err)
			}
		}
	}
	log.Print("Done updating users")

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
			debug("[private:HandlePrivateMessage] app.config.DebugUser non-empty, user ID does not match app.config.DebugUser")
			return
		}
	}

	_, source, errors := app.commandManager.Process(*event.Message)
	if errors != nil {
		for _, e := range errors {
			if e != nil {
				log.Print(e)
				e = app.WarnUserError(event.Message.ChannelID, e.Error())
				if e != nil {
					log.Print(e)
				}
			}
		}
	}

	if source != CommandSourcePRIVATE && source != CommandSourceADMINISTRATIVE {
		err := app.chatLogger.RecordChatLog(event.Message.Author.ID, event.Message.ChannelID, event.Message.Content)
		if err != nil {
			log.Print(err)
		}

		for i := range event.Message.Mentions {
			if event.Message.Mentions[i].ID == app.config.BotID {
				err := app.HandleSummon(*event.Message)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (app *App) onJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	verified, err := app.IsUserVerified(event.Member.User.ID)
	if err != nil {
		log.Print(err)
	}

	err = app.discordClient.GuildMemberRoleAdd(app.config.GuildID, event.Member.User.ID, app.config.NormalRole)
	if err != nil {
		log.Print(err)
	}

	if !verified {
		ch, err := s.UserChannelCreate(event.Member.User.ID)
		if err != nil {
			log.Print(err)
			return
		}
		_, err = app.discordClient.ChannelMessageSend(ch.ID, app.locale.GetLangString("en", "AskUserVerify"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (app *App) onLeave(s *discordgo.Session, event *discordgo.GuildMemberRemove) {
	verified, err := app.IsUserVerified(event.Member.User.ID)
	if err != nil {
		log.Print(err)
	}

	err = app.discordClient.GuildMemberRoleRemove(app.config.GuildID, event.Member.User.ID, app.config.VerifiedRole)
	if err != nil {
		log.Print(err)
	}
}

func (app *App) onHeartbeat(t time.Time) {
	err := app.HandleHeartbeatEvent(t)
	if err != nil {
		log.Print(err)
	}
}
