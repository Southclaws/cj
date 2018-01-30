package main

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func commandCJQuote(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	rand.Seed(time.Now().UnixNano())
	_, err := cm.App.discordClient.ChannelMessageSend(cm.App.config.PrimaryChannel, "you're mom gay")
	return true, false, err
}
