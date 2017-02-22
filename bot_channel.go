package main

import (
	"github.com/bwmarrin/discordgo"
)

// HandleChannelMessage processes any public message from the channel that is
// designated as the primary channel for the bot.
// (see api.config.PrimaryChannel)
func (app App) HandleChannelMessage(message discordgo.Message) error {
	debug("[channel:HandleChannelMessage] %v", message)
	return nil
}
