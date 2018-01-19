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
			result += cm.App.locale.GetLangString(cm.App.config.Language, "UserNotVerified", user.ID) + " "
		} else {
			link, err = cm.App.GetForumUserFromDiscordUser(user.ID)
			if err != nil {
				return false, false, errors.Wrap(err, cm.App.locale.GetLangString(cm.App.config.Language, "CommandUserInfoFailedGetForumUser"))
			}

			cachedProfile, found = cm.App.cache.Get(link)
			if found {
				profile = *(cachedProfile.(*UserProfile))
			} else {
				profile, err = cm.App.GetUserProfilePage(link)
				if err != nil {
					return false, false, errors.Wrap(err, cm.App.locale.GetLangString(cm.App.config.Language, "CommandUserInfoFailedGetProfilePage"))
				}
			}

			result += cm.App.locale.GetLangString(cm.App.config.Language, "CommandUserInfoProfile", profile.UserName, profile.JoinDate, profile.TotalPosts, profile.Reputation) + " "
		}
	}

	_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		return false, false, errors.Wrap(err, cm.App.locale.GetLangString(cm.App.config.Language, "CommandUserInfoFailedSendMessage"))
	}

	return true, false, err
}
