package admod

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

var http = regexp.MustCompile(`\w+:\/\/(.+)`)

type Watcher struct {
	config  *types.Config
	discord *discord.Session
	api     storage.Storer
	fc      *forum.ForumClient

	channel string
}

func (w *Watcher) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
	fc *forum.ForumClient,
) (err error) {
	w.config = config
	w.discord = discord
	w.api = api
	w.fc = fc

	w.channel = config.AdsChannel

	return nil
}

func (w *Watcher) OnMessage(m discordgo.Message) error {
	if m.ChannelID != w.channel {
		return nil
	}

	if len(m.Content) < 100 {
		ch, err := w.discord.S.UserChannelCreate(m.Author.ID)
		if err != nil {
			return err
		}
		w.discord.ChannelMessageSend(ch.ID, "Your ad was deleted as it was too short. Low-effort ads are removed from the channel to cut down on noise. Ads with only an IP or a link have a near-zero click-through rate and only serve to waste space in the channel. Write a short description about the server explaining why players should try it out!")
	} else {
		message := http.ReplaceAllString(m.Content, "https://r.open.mp/$1")
		w.discord.ChannelMessageSendEmbed(w.channel, &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeArticle,
			Title:       fmt.Sprintf("Server Ad (posted by %s)", m.Author.Username),
			Description: message,
		})
	}

	if err := w.discord.S.ChannelMessageDelete(w.channel, m.ID); err != nil {
		return err
	}

	return nil
}
