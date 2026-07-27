package heartbeat

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/bot/heartbeat/common"
	"github.com/Southclaws/cj/bot/heartbeat/readme"
	"github.com/Southclaws/cj/bot/heartbeat/stats"
	"github.com/Southclaws/cj/bot/heartbeat/talking"
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

	Readme *readme.Readme

	jobs      *jobRegistry
	callables map[string]func() error
}

func (a *Heartbeat) Statuses() []JobStatus {
	if a.jobs == nil {
		return nil
	}
	return a.jobs.snapshot()
}

var ErrJobNotFound = errors.New("no such job is registered")

func (a *Heartbeat) RunJob(name string) error {
	fn, ok := a.callables[name]
	if !ok {
		return ErrJobNotFound
	}
	err := fn()
	a.jobs.recordRun(name, err)
	return err
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

	a.jobs = newJobRegistry()
	a.callables = make(map[string]func() error)

	readmeProvider := &readme.Readme{}
	a.Readme = readmeProvider

	aps := []ActionProvider{
		&stats.Aggregator{},
		&talking.Talk{},
		readmeProvider,
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
			jobName := name
			if len(actions) > 1 {
				jobName = fmt.Sprintf("%s.%d", name, i)
			}
			a.jobs.register(jobName, name, action.Schedule)
			a.callables[jobName] = action.Call

			zap.L().Debug("adding action call", zap.String("schedule", action.Schedule))
			if err = cr.AddFunc(action.Schedule, func() {
				if rand.Float64() < action.Chance {
					e := action.Call()
					a.jobs.recordRun(jobName, e)
					if e != nil {
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
