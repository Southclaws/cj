package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/storage"
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

	user := cm.Storage.GetUserOrCreate(message.Author.ID)
	sort.Slice(user.ReceivedReactions, func(i, j int) bool {
		return user.ReceivedReactions[i].Counter > user.ReceivedReactions[j].Counter
	})
	embed, err := FormatUserReactions(&user.ReceivedReactions, message.Author, cm.Discord)

	if err != nil {
		return
	}

	_, err = cm.Discord.S.ChannelMessageSendEmbed(message.ChannelID, embed)

	return
}

func FormatUserReactions(reactions *[]storage.ReactionCounter, author *discordgo.User, session *discord.Session) (embed *discordgo.MessageEmbed, err error) {
	statsMessage := strings.Builder{}
	statsMessage.WriteString(fmt.Sprintf("**%s's Reactions**\n\n", author.Username)) //nolint:errcheck

	embed = &discordgo.MessageEmbed{Color: 0x3498DB}
	if len(*reactions) == 0 {
		statsMessage.WriteString("There are no reactions to display here!")
	}
	for _, reaction := range *reactions {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   reaction.Reaction,
			Value:  fmt.Sprintf("%dx", reaction.Counter),
			Inline: true,
		})
	}

	// Add one more 'filling' embed field in case there are 2
	// fields (with 1 and 3 fields the formatting is fine)
	if len(*reactions)%3 == 2 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "", // Note: the content of this seemingly empty string is u0080
			Value:  "", // Note: the content of this seemingly empty string is u0080
			Inline: true,
		})
	}

	embed.Description = statsMessage.String()

	return embed, nil
}
