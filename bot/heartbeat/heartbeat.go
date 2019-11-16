package heartbeat

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/bot/heartbeat/common"
	"github.com/Southclaws/cj/bot/heartbeat/readme"
	"github.com/Southclaws/cj/bot/heartbeat/stats"
	"github.com/Southclaws/cj/bot/heartbeat/talking"
	"github.com/Southclaws/cj/bot/heartbeat/wiki"
	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// ActionProvider describes a type that provides a registration function that
// provides a set of actions that have some random chance to be called on a
// schedule. The `chance` parameter controls the probability the `call` function
// is called each time a cron job specified by `schedule` occurs.
type ActionProvider interface {
	Init(*types.Config, *discord.Session, storage.Storer, *forum.ForumClient) (string, error)
	Register() []common.Action
}

// Heartbeat controls a set of plugins that do stuff periodically.
type Heartbeat struct {
	Config  *types.Config
	Discord *discord.Session
	Storage storage.Storer
	Forum   *forum.ForumClient
}

//nolint:golint
func (a *Heartbeat) Init(
	config *types.Config,
	discord *discord.Session,
	api storage.Storer,
	fc *forum.ForumClient,
) (err error) {
	a.Config = config
	a.Storage = api
	a.Discord = discord
	a.Forum = fc

	zap.L().Debug("initialising heartbeat module")

	aps := []ActionProvider{
		// &hello{}, // for testing lol
		&stats.Aggregator{},
		&talking.Talk{},
		&readme.Readme{},
		&wiki.Wiki{},
	}

	cr := cron.New()
	for _, ap := range aps {
		var name string
		if name, err = ap.Init(config, discord, api, fc); err != nil {
			return errors.Wrapf(err, "failed to initialise heatbeat extension %v", a)
		}
		actions := ap.Register()
		zap.L().Debug("loading heartbeat action provider",
			zap.String("name", name),
			zap.Int("actions", len(actions)))
		for i := range actions {
			action := actions[i]
			zap.L().Debug("adding action call", zap.String("schedule", action.Schedule))
			if err = cr.AddFunc(action.Schedule, func() {
				if rand.Float64() < action.Chance {
					if e := action.Call(); e != nil {
						zap.L().Error("action failed", zap.Error(e))
					}
				}
			}); err != nil {
				return errors.Wrap(err, "failed to add heartbeat action")
			}
		}
	}
	cr.Start()
	return
}

//nolint:golint
func (a *Heartbeat) OnMessage(discordgo.Message) (err error) {
	return
}

// testing type

type hello struct {
	d *discord.Session
}

func (h *hello) Init(c *types.Config, d *discord.Session, s storage.Storer, f *forum.ForumClient) error {
	h.d = d
	return nil
}

func (h *hello) Register() []common.Action {
	return []common.Action{
		{
			Schedule: "* * * * *",
			Chance:   0.1,
			Call: func() error {
				h.d.ChannelMessageSend("465142687985696788", fmt.Sprintf("%v", time.Now()))
				return nil
			},
		},
	}
}
