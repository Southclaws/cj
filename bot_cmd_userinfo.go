package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func commandUserInfo(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	var (
		profile       UserProfile
		verified      bool
		err           error
		link          string
		cachedProfile interface{}
		found         bool
		result        string
	)

	if len(message.Mentions) == 0 {
		return false, false, nil
	}

	for _, user := range message.Mentions {
		verified, err = cm.App.IsUserVerified(user.ID)
		if err != nil {
			result += err.Error() + " "
			continue
		}

		if !verified {
			result += cm.App.locale.GetLangString("en", "UserNotVerified", user.ID) + " "
		} else {
			link, err = cm.App.GetForumUserFromDiscordUser(user.ID)
			if err != nil {
				return errors.Wrap(err, "failed to get forum user from discord user")
			}

			cachedProfile, found = cm.App.cache.Get(link)
			if found {
				profile = *(cachedProfile.(*UserProfile))
			} else {
				profile, err = cm.App.GetUserProfilePage(link)
				if err != nil {
					return errors.Wrap(err, "failed to get user profile page")
				}
			}

			result += cm.App.locale.GetLangString("en", "CommandUserInfoProfile", profile.UserName, profile.JoinDate, profile.TotalPosts, profile.Reputation) + " "
		}
	}

	_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		return errors.Wrap(err, "failed to send channel message")
	}

	return true, false, err
}
