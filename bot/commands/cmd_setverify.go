package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandSetVerify(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	if len(message.Mentions) != 1 {
		err = errors.New("Expected one mention")
		return
	}

	target := message.Mentions[0]

	err = cm.Storage.StoreVerifiedUser(types.Verification{
		DiscordUser: *target,
		ForumUser:   "",
		Code:        "null",
	})

	return
}
