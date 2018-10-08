package stats

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"

	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// Aggregator collects statistics about messages and users
type Aggregator struct {
	Config  *types.Config
	Discord *discordgo.Session
	Storage *storage.API
	Forum   *forum.ForumClient

	topMessages storage.TopMessages

	err error
}

//nolint:golint
func (a *Aggregator) Init(*types.Config, *discordgo.Session, *storage.API, *forum.ForumClient) (err error) {
	cron := cron.New()
	must(cron.AddFunc("@hourly", a.gather))
	must(cron.AddFunc("@every 7h30m", a.announce))
	cron.Start()
	return
}

//nolint:golint
func (a *Aggregator) OnMessage(discordgo.Message) (err error) {
	return
}

func (a *Aggregator) gather() {
	err := a.gatherTopMessages(10)
	if err != nil {
		a.err = err
	}
}

func (a *Aggregator) announce() {
	statsMessage := strings.Builder{}
	statsMessage.WriteString("Statistics") //nolint:errcheck

	embed := &discordgo.MessageEmbed{Color: 0x3498DB}
	for _, tm := range a.topMessages {
		statsMessage.WriteString(fmt.Sprintf("**%s** - %d\n", tm.User, tm.Messages)) //nolint:errcheck
	}

	embed.Description = statsMessage.String()

	_, err := a.Discord.ChannelMessageSendEmbed(a.Config.PrimaryChannel, embed)
	if err != nil {
		a.err = err
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
