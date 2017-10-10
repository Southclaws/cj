package main

import (
	"time"
)

// ChatLogger is responsible for logging chat messages
type ChatLogger struct {
	app *App
}

// ChatLog represents a single logged chat message from Discord
type ChatLog struct {
	Timestamp      int64
	DiscordUserID  string
	DiscordChannel string
	Message        string
}

// StartChatLogger creates a new chat logger for the app
func (app *App) StartChatLogger() {
	cl := ChatLogger{app: app}

	app.chatLogger = &cl

	go cl.flushTicker()
}

// RecordChatLog records a chat message from a user in a channel
func (cl *ChatLogger) RecordChatLog(discordUserID string, discordChannel string, message string) error {
	record := ChatLog{
		time.Now().Unix(),
		discordUserID,
		discordChannel,
		message,
	}

	cl.app.queue.Add(record)

	if cl.app.queue.Len() >= cl.app.config.LogFlushAt {
		cl.flushChatLogs()
	}

	return nil
}

func (cl *ChatLogger) flushTicker() {
	t := time.NewTicker(time.Minute * time.Duration(cl.app.config.LogFlushInterval))
	for range t.C {
		cl.flushChatLogs()
	}
	logger.Warn("flushTicker ended unexpectedly")
}

func (cl *ChatLogger) flushChatLogs() {
	if cl.app.queue.Len() == 0 {
		return
	}

	// tx := cl.app.db.Begin()
	// for {
	// 	message := cl.app.queue.Next()
	// 	if message == nil {
	// 		break
	// 	}
	// 	tx.Create(message.(ChatLog))
	// }
	// tx.Commit()
	return
}

// TODO: getters in order to run tests
