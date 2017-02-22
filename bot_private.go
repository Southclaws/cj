package main

import (
	"github.com/bwmarrin/discordgo"
)

// HandlePrivateMessage processes a private message sent directly to the bot
// usually for direct commands such as account verification.
func (app App) HandlePrivateMessage(message discordgo.Message) error {
	debug("[private:HandlePrivateMessage] %v", message)

	if message.Content == "verify" {
		// Begin verification process - bot_verify.go
	}

	return nil
}
