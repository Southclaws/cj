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
	"Use api.open.mp/server/",
	"/fuck off",
	"I forgot what I was about to say, I will mention you after it comes to my mind.",
	"WHY ARE YOU DOING THIS",
	"Did you mean /stats",
	"Your query is being executed.",
	"CJ has DMed you, please check your DMs.",
	"Invalid server, not recognized by open.mp/servers",
	"Server doesn't allow querying.",
	"Server sent invalid packets.",
	"Usage: /stats [IP/hostname:port]",
	"Request timed out.",
	"error",
	"Invalid characters in server's name.",
	"/stats is under maintenance for the time being.",
	"This command is disabled.",
	"Only boosters can use this command.",
	"You have to pay for this command. Monthly subscription starts at $4",
	"Error occurred while parsing the JSON data.",
	"Server is full. CJ can't join the server to get the required information",
	"Wrong channel",
	":clock1:",
	"Please don't spam commands in ongoing conversations.",
	"Host CJ yourself and use this command.",
	"CJ is busy with the previous query.",
	"This command has been muted for you. Reason: spam.",
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
	if len(args) == 0 {
		cm.Discord.ChannelMessageSend(message.ChannelID, "Usage: /stats [IP/hostname:port]")
		return
	}
	rand.Seed(time.Now().UnixNano())
	funnyStats := funnyStats[rand.Intn(len(funnyStats))]
	cm.Discord.ChannelMessageSend(message.ChannelID, funnyStats)
	return
}
