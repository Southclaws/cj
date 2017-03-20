package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func commandUserInfo(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	var profile UserProfile
	var verified bool
	var err error

	if len(message.Mentions) == 0 {
		return false, false, nil
	}

	for _, user := range message.Mentions {
		verified, err = cm.App.IsUserVerified(user.ID)
		if err != nil {
			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, err.Error())
			if err != nil {
				log.Print(err)
			}
			continue
		}

		if !verified {
			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString("en", "UserNotVerified", user.ID))
			if err != nil {
				log.Print(err)
			}
		} else {
			link, err := cm.App.GetForumUserFromDiscordUser(user.ID)
			if err != nil {
				log.Print(err)
			}

			cachedProfile, found := cm.App.cache.Get(link)
			if found {
				profile = *(cachedProfile.(*UserProfile))
			} else {
				profile, err = cm.App.GetUserProfilePage(link)
				if err != nil {
					log.Print(err)
				}
			}

			_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, cm.App.locale.GetLangString("en", "CommandUserInfoProfile", profile.UserName, profile.JoinDate, profile.TotalPosts, profile.Reputation))
			if err != nil {
				log.Print(err)
			}
		}
	}

	return true, false, err
}
