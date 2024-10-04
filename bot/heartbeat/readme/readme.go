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
) (name string, err error) {
	r.Config = config
	r.Storage = api
	r.Discord = discord
	r.Forum = fc
	return "readme", nil
}

//nolint:golint
func (r *Readme) Register() (actions []common.Action) {
	r.fetchReadme()

	return []common.Action{
		{
			Schedule: "@every 6h",
			Chance:   1,
			Call:     r.fetchReadme,
		},
	}
}

func (r *Readme) fetchReadme() (err error) {
	discord := r.Discord.S

	discord.ChannelMessageSend("948604467887083550", "----- Starting to perform readme update")

	discord.ChannelMessageSend("948604467887083550", "Trying to get readme message")

	m, e := r.Storage.GetReadmeMessage()

	// if it's not in the database, we return nil
	// because it's accepted scenario
	// in that case /readme needs to be called manually
	if e == mgo.ErrNotFound {
		err = nil
		return
	}

	discord.ChannelMessageSend("948604467887083550", "Trying to get readme channel")

	// get the readme channel
	c, e := r.Discord.S.Channel(r.Config.ReadmeChannel)
	if e != nil {
		err = e
		return
	}

	discord.ChannelMessageSend("948604467887083550", "Trying to get already sent readme message")

	// get the already sent readme message
	msg, e := r.Discord.S.ChannelMessage(c.ID, m)
	if e != nil {
		err = e
		return
	}

	discord.ChannelMessageSend("948604467887083550", "Fetching readme message from rules repository")

	// fetch upstream gist content
	ctx, e := r.Storage.FetchReadmeMessage(r.Config.ReadmeGithubOwner, r.Config.ReadmeGithubRepository, r.Config.ReadmeFileName)
	if e != nil {
		err = e
		return
	}

	discord.ChannelMessageSend("948604467887083550", "Attempting to update readme message if needed")

	// call update function
	err = r.Storage.UpdateReadmeMessage(r.Discord.S, msg, ctx)
	if err != nil {
		return
	}

	discord.ChannelMessageSend("948604467887083550", "----- Updating readme task is finished")
	return
}
