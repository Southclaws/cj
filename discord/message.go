package discord

import (
	"io"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var replacer = strings.NewReplacer(
	"@everyone", "@**everyone**",
	"@here", "@**here**",
)

// ChannelMessageSend filters out everyone and here mentions
func (s *Session) ChannelMessageSend(channelID string, content string) {
	//nolint:errcheck
	s.S.ChannelMessageSend(channelID, replacer.Replace(content))
}

// ChannelMessageSendEmbed filters out everyone and here mentions
func (s *Session) ChannelMessageSendEmbed(channelID string, content *discordgo.MessageEmbed) {
	content.Description = replacer.Replace(content.Description)
	//nolint:errcheck
	s.S.ChannelMessageSendEmbed(channelID, content)
}

// ChannelFileSend is a wrapper in case we need to block haram things in the future
func (s *Session) ChannelFileSend(channelID string, name string, r io.Reader) {
	s.S.ChannelFileSend(channelID, name, r)
}
