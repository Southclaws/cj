package discord

import (
	"io"
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

// ChannelFileSend is a wrapper in case we need to block haram things in the future
func (s *Session) ChannelFileSend(channelID string, name string, r io.Reader) {
	s.S.ChannelFileSend(channelID, name, r)
}
