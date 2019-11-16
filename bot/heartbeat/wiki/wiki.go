package wiki

import (
	"github.com/Southclaws/cj/bot/heartbeat/common"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// Wiki holds wiki git repository
type Wiki struct {
	Config  *types.Config
	Discord *discord.Session
	Storage storage.Storer
	Forum   *forum.ForumClient
}

//nolint:golint
func (w *Wiki) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
	fc *forum.ForumClient,
) (name string, err error) {
	w.Config = config
	w.Storage = api
	w.Discord = discord
	w.Forum = fc
	return "wiki", nil
}

//nolint:golint
func (w *Wiki) Register() (actions []common.Action) {
	return []common.Action{
		{
			Schedule: "@every 1h",
			Chance:   1,
			Call:     w.pullWiki,
		},
	}
}

func (w *Wiki) pullWiki() (err error) {
	err = w.Storage.PullWiki(w.Config.WikiURL)
	return
}
