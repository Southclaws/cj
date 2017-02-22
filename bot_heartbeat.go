package main

import "time"

// HandleHeartbeatEvent is run automatically every now and then to send random
// messages to the bot's primary channel.
// (see api.config.Heartbeat and api.config.PrimaryChannel)
func (app App) HandleHeartbeatEvent(t time.Time) error {
	debug("[Heartbeat:HandleHeartbeatEvent] %v", t)
	return nil
}
