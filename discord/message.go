package discord

import (
	"strings"
)

// ChannelMessageSend filters out everyone and here mentions
func (s *Session) ChannelMessageSend(channelID string, content string) {
	r := strings.NewReplacer(
		"@everyone", "@**everyone**",
		"@here", "@**here**",
	)
	//nolint:errcheck
	s.S.ChannelMessageSend(channelID, r.Replace(content))
}
