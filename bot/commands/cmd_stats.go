package commands

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

var funnyStats = []string{
    "only rollerplayers can query this server",
    "Error: NSFW detected",
    "/help",
    "/commands",
    "server contains n word in its name",
    "please sir we only allow android servers here",
    "/stats /stats /stats. CJ lost his mind",
    "cj lost his mom in that server so it's a no go zone for him",
    "Use api.open.mp/server",
    "/fuck off",
    "I forgot what I was about say, I will tell you after it comes to my mind.",
    "WHY ARE YOU DOING THIS",
    "Did you mean /stats",
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
	rand.Seed(time.Now().UnixNano())
	funnyStats := funnyStats[rand.Intn(len(funnyStats))]
	cm.Discord.ChannelMessageSend(message.ChannelID, funnyStats)
	return
}
