package main

import (
	"time"

	"github.com/pkg/errors"
)

// ChatLog represents a single logged chat message from Discord
type ChatLog struct {
	Timestamp      int64
	DiscordUserID  string
	DiscordChannel string
	Message        string
}

// RecordChatLog records a chat message from a user in a channel
func (app App) RecordChatLog(discordUserID string, discordChannel string, message string) (err error) {
	record := ChatLog{
		time.Now().Unix(),
		discordUserID,
		discordChannel,
		message,
	}

	err = app.chat.Insert(record)
	if err != nil {
		err = errors.Wrap(err, "failed to insert chat log")
	}

	return
}
