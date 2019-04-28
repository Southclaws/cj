package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/forum"
)

func (cm *CommandManager) commandUserInfo(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	var (
		profile       forum.UserProfile
		verified      bool
		link          string
		cachedProfile interface{}
		found         bool
		result        string
	)

	if len(message.Mentions) == 0 {
		err = errors.New("Expected one or more mention")
		return
	}

	for _, user := range message.Mentions {
		verified, err = cm.Storage.IsUserVerified(user.ID)
		if err != nil {
			result += err.Error() + " "
			continue
		}

		if !verified {
			verified, err = cm.Storage.IsUserLegacyVerified(user.ID)
			if err != nil {
				result += err.Error() + " "
				continue
			}

			if verified {
				result += fmt.Sprintf("The user <@%s> is not verified on burgershot, they need to message CJ with `verify`.", user.ID)
			} else {
				result += fmt.Sprintf("<@%s> is not verified. ", user.ID)
			}
		} else {
			_, link, err = cm.Storage.GetForumUserFromDiscordUser(user.ID)
			if err != nil {
				err = errors.Wrap(err, "failed to get forum user from discord user")
				return
			}

			cachedProfile, found = cm.Cache.Get(link)
			if found {
				profile = *(cachedProfile.(*forum.UserProfile))
			} else {
				profile, err = cm.Forum.GetUserProfilePage(link)
				if err != nil {
					err = errors.Wrap(err, "failed to get user profile page")
					return
				}
			}

			result += fmt.Sprintf(
				"*Username:* %s *Member since:* %s *Total posts:* %d *Reputation:* %d ",
				profile.UserName, profile.JoinDate, profile.TotalPosts, profile.Reputation)
		}
	}

	cm.Discord.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		err = errors.Wrap(err, "failed to send message")
		return
	}

	return false, nil
}
