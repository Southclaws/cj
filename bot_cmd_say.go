package main

import "github.com/bwmarrin/discordgo"

func commandSay(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, error) {
	_, err := cm.App.discordClient.ChannelMessageSend(cm.App.config.PrimaryChannel, args)
	return false, err
}
