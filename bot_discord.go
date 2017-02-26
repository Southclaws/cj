package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"net/http"

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

	err = app.discordClient.Open()
	if err != nil {
		log.Println("discord client connection error")
		log.Fatal(err)
	}

	debug("awaiting Discord ready state...")

	return nil
}

func (app *App) onReady(s *discordgo.Session, event *discordgo.Ready) {
	debug("discord ready")
	app.ready <- true

	found := false
	roles, err := s.GuildRoles(app.config.GuildID)
	if err != nil {
		log.Fatal(err)
	}
	for _, role := range roles {
		if role.ID == app.config.VerifiedRole {
			found = true
			break
		}
	}
	if !found {
		log.Print("verified role ID was not found in guild role list:")
		for _, role := range roles {
			log.Printf("name: %s id: %s", role.Name, role.ID)
		}
		log.Fatalf("role '%s' not found.", app.config.VerifiedRole)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(app.config.Heartbeat))
	for t := range ticker.C {
		app.onHeartbeat(t)
	}
}

func (app *App) onMessage(s *discordgo.Session, event *discordgo.MessageCreate) {
	if len(app.ready) > 0 {
		<-app.ready
	}

	if event.Message.Author.ID == app.config.BotID {
		return
	}

	message := event.Message
	primary := false
	private := false

	if message.ChannelID == app.config.PrimaryChannel {
		primary = true
	} else {
		// discordgo has not implemented private channel objects (DM Channels)
		// so we have to perform the request manually and unmarshal the response
		// object into a `ChannelDM` object.
		var err error
		var req *http.Request
		var response *http.Response
		var body []byte
		if req, err = http.NewRequest("GET", discordgo.EndpointChannel(message.ChannelID), nil); err != nil {
			log.Print(err)
		}
		req.Header.Add("Authorization", "Bot "+app.config.DiscordToken)
		if response, err = app.httpClient.Do(req); err != nil {
			log.Print(err)
		}
		if body, err = ioutil.ReadAll(response.Body); err != nil {
			log.Print(err)
		}
		channel := ChannelDM{}
		json.Unmarshal(body, &channel)

		// Now we have one of these:
		// https://discordapp.com/developers/docs/resources/channel#dm-channel-object

		if channel.Private {
			private = true
		}
	}

	debug("private: %v primary %v", private, primary)

	if private {
		err := app.HandlePrivateMessage(*message)
		if err != nil {
			log.Print(err)
		}
	} else {
		log.Printf("%p", app.chatLogger)
		app.chatLogger.RecordChatLog(message.Author.ID, message.ChannelID, message.Content)

		for i := range message.Mentions {
			if message.Mentions[i].ID == app.config.BotID {
				err := app.HandleSummon(*message)
				if err != nil {
					log.Print(err)
				}
			}
		}

		if primary {
			app.HandleChannelMessage(*message)
		}
	}
}

func (app *App) onHeartbeat(t time.Time) {
	app.HandleHeartbeatEvent(t)
}
