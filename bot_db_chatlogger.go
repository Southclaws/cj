package main

import (
	"log"
	"time"
)

// ChatLogger is reponsible for logging chat messages
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
	var cl ChatLogger

	cl = ChatLogger{app: app}

	app.chatLogger = &cl

	go cl.flushTicker()
}

// RecordChatLog records a chat message from a user in a channel
func (cl *ChatLogger) RecordChatLog(discordUserID string, discordChannel string, message string) error {
	log.Printf("[chatlogger:RecordChatLog] %s #%s: '%s'", discordUserID, discordChannel, message)
	var record ChatLog

	record = ChatLog{
		time.Now().Unix(),
		discordUserID,
		discordChannel,
		message,
	}

	log.Print(record)
	log.Printf("cl %p", cl)
	log.Printf("ap %p", cl.app)
	log.Printf("aq %p", cl.app.queue)

	cl.app.queue.Add(record)

	if cl.app.queue.Len() >= cl.app.config.LogFlushAt {
		cl.flushChatLogs()
	}

	return nil
}

func (cl *ChatLogger) flushTicker() {
	log.Print(cl.app.config.LogFlushInterval)
	t := time.NewTicker(time.Minute * time.Duration(cl.app.config.LogFlushInterval))
	for range t.C {
		cl.flushChatLogs()
	}
	log.Printf("ERROR: flushTicker ended unexpectedly")
	<-cl.app.done
}

func (cl *ChatLogger) flushChatLogs() {
	if cl.app.queue.Len() == 0 {
		return
	}

	tx := cl.app.db.Begin()
	for {
		message := cl.app.queue.Next()
		if message == nil {
			break
		}
		tx.Create(message.(ChatLog))
	}
	tx.Commit()
	return
}
