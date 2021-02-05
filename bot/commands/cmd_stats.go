package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"net"
	
	"github.com/bwmarrin/discordgo"
	"gopkg.in/resty.v1"

	"github.com/Southclaws/cj/types"
)

type serverCore struct {
	Address    string `json:"ip"`
	Hostname   string `json:"hn"`
	Players    int    `json:"pc"`
	MaxPlayers int    `json:"pm"`
	Language   string `json:"la"`
	Password   bool   `json:"pa"`
	Version    string `json:"vn"`
}

type serverListing struct {
	Core        serverCore `json:"core"`
	Description string     `json:"description"`
	Banner      string     `json:"banner"`
	Active      bool       `json:"active"`
	Error       string     `json:"error"`
}

func (cm *CommandManager) commandStats(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	var (
		server   string
		password string
	)

	if strings.Contains(args, ":") {
		server = args
	} else {
		server = args + ":7777"
	}

	//Incase if user types host name instead of IP of server
	host, port, _ := net.SplitHostPort(server)
	IPs, err := net.LookupHost(host)
	if err == nil {
		server = net.JoinHostPort(IPs[0], port)
	}
	
	resp, err := resty.R().Get("https://api.open.mp/server/" + server)
	if err != nil {
		cm.Discord.ChannelMessageSend(message.ChannelID, "Unable to query the open.mp servers API.")
		return
	}

	serverInfo, err := decodeServerInfo(resp.String())
	if err != nil {
		println(err)
	}

	if serverInfo.Error != "" {
		cm.Discord.ChannelMessageSend(message.ChannelID, serverInfo.Error)
		return
	}

	if serverInfo.Core.Password == true {
		password = "Yes"
	} else {
		password = "No"
	}

	embedData := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
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
				Name:   "🔒 Password",
				Value:  password,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: serverInfo.Banner,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Made possible by open.mp servers api",
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
