package readme

import (
	"github.com/Southclaws/cj/bot/heartbeat/common"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
	"github.com/globalsign/mgo"
)

// Readme fetches the upsream gist
type Readme struct {
	Config  *types.Config
	Discord *discord.Session
	Storage storage.Storer
	Forum   *forum.ForumClient
}

//nolint:golint
func (r *Readme) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
	fc *forum.ForumClient,
) (err error) {
	r.Config = config
	r.Storage = api
	r.Discord = discord
	r.Forum = fc
	return
}

//nolint:golint
func (r *Readme) Register() (actions []common.Action) {
	return []common.Action{
		{
			Schedule: "@every 12h",
			Chance:   1,
			Call:     r.fetchReadme,
		},
	}
}

func (r *Readme) fetchReadme() (err error) {
	m, e := r.Storage.GetReadmeMessage()

	// if it's not in the database, we return nil
	// because it's accepted scenario
	// in that case /readme needs to be called manually
	if e == mgo.ErrNotFound {
		err = nil
		return
	}

	// get the readme channel
	c, e := r.Discord.S.Channel(r.Config.ReadmeChannel)
	if e != nil {
		err = e
		return
	}

	// get the already sent readme message
	msg, e := r.Discord.S.ChannelMessage(c.ID, m)
	if e != nil {
		err = e
		return
	}

	// fetch upstream gist content
	ctx, e := r.Storage.FetchReadmeMessage(r.Config.ReadmeGistID, r.Config.ReadmeGistFileName)
	if e != nil {
		err = e
		return
	}

	// call update function
	err = r.Storage.UpdateReadmeMessage(r.Discord.S, msg, ctx)
	if err != nil {
		return
	}

	return
}
