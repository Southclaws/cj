package admod

import (
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

	message := http.ReplaceAllString(m.Content, "https://www.open.mp/l/$1")
	w.discord.ChannelMessageSend(w.channel, message)
	if err := w.discord.S.ChannelMessageDelete(w.channel, m.ID); err != nil {
		return err
	}

	return nil
}
