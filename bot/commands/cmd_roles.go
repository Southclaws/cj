package commands

import (
	"fmt"
	"strings"

	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandRoles(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	roles, err := cm.Discord.S.GuildRoles(cm.Config.GuildID)
	if err != nil {
		cm.replyDirectly(interaction, err.Error())
		return
	}
	msg := strings.Builder{}
	msg.WriteString("Roles:\n")
	for _, r := range roles {
		msg.WriteString(fmt.Sprintf("`%s`: %s\n", r.ID, r.Name))
	}
	cm.replyDirectly(interaction, msg.String())
	return
}
