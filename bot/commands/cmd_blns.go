package commands

import (
	// "math/rand"

	"github.com/bwmarrin/discordgo"
	// "github.com/fluidkeys/fork-big-list-of-naughty-strings/naughtystrings"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandBLNS(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	// s := naughtystrings.Unencoded()
	// i := rand.Int31n(int32(len(s)))
	// cm.Discord.ChannelMessageSend(cm.Config.PrimaryChannel, s[i])
	return
}
