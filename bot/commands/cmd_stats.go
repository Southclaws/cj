package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/samp-servers-api/types"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/resty.v1"
)

type serverRules struct {
	Mapname   string `json:"mapname"`
	Weather   string `json:"weather"`
	Worldtime string `json:"worldtime"`
	Version   string `json:"version"`
	Weburl    string `json:"weburl"`
	Artwork   string `json:"artwork"`
}

type serverListing struct {
	Core        types.ServerCore `json:"core"`
	Rules       serverRules      `json:"ru,omitempty"`
	Description string           `json:"description"`
	Banner      string           `json:"banner"`
	Active      bool             `json:"active"`
}

func (cm *CommandManager) commandStats(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	var (
		server   string
		password string
		artwork  string
	)

	if strings.Contains(args, ":") {
		server = args
	} else {
		server = args + ":7777"
	}

	resp, err := resty.R().Get("https://api.samp-servers.net/v2/server/" + server)
	if err != nil {
		cm.Discord.ChannelMessageSend(message.ChannelID, "Unable to query the SA:MP servers API.")
		return
	}

	if strings.Contains(resp.String(), "could not find server by address") {
		cm.Discord.ChannelMessageSend(message.ChannelID, "Invalid server. Not recognised by https://samp-servers.net.")
		return
	}

	serverInfo, err := decodeServerInfo(resp.String())
	if err != nil {
		println(err)
	}

	if serverInfo.Core.Password == true {
		password = "Yes"
	} else {
		password = "No"
	}

	if len(serverInfo.Rules.Artwork) > 0 {
		artwork = "Yes"
	} else {
		artwork = "No"
	}

	embedData := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://" + serverInfo.Rules.Weburl,
			Name:    serverInfo.Core.Hostname,
			IconURL: "https://github.com/Southclaws/cj/raw/master/cj.png",
		},
		Title:       serverInfo.Core.Address,
		Description: serverInfo.Description,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "👨‍👨‍👧‍👦 Players",
				Value:  fmt.Sprintf("%d/%d", serverInfo.Core.Players, serverInfo.Core.MaxPlayers),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "💻 Version",
				Value:  serverInfo.Core.Version,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "💾 Password",
				Value:  password,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "🎭 Artwork",
				Value:  artwork,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "⏱ Time",
				Value:  serverInfo.Rules.Worldtime,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "🌞 Weather",
				Value:  serverInfo.Rules.Weather,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: serverInfo.Banner,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Made possible by samp-servers.net",
		},
		Color:     0x006400,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = cm.Discord.S.ChannelMessageSendEmbed(message.ChannelID, embedData)
	if err != nil {
		println(err)
	}
	return
}

func decodeServerInfo(data string) (info serverListing, err error) {
	err = json.Unmarshal([]byte(data), &info)
	return
}
