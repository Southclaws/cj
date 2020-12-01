package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandRep(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	var reaction string
	if len(args) == 0 {
		// When setting it: pass it as <:emoji:id_as_long_number>
		reaction = fmt.Sprintf("%v", settings.Misc["default"])
		if reaction == "<nil>" {
			return false, errors.New("No default emoji ID set, please fill one in the field defaultID for this command's config (/config /rep)")
		}
	} else {
		reaction = args
	}

	user := cm.Storage.GetUserOrCreate(message.Author.ID)

	count := 0
	for _, v := range user.ReceivedReactions {
		if v.Reaction == reaction {
			count = v.Counter
			break
		}
	}
	cm.Discord.ChannelMessageSend(message.ChannelID,
		fmt.Sprintf("Your %s count: %d", reaction, count))
	return
}
