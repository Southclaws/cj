package main

import (
	"github.com/bwmarrin/discordgo"
	"gopkg.in/mgo.v2/bson"
)

func commandUnVerify(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	if len(message.Mentions) == 0 {
		return false, false, nil
	}

	target := message.Mentions[0]

	err := cm.App.accounts.Remove(bson.M{"discord_user_id": target.ID})
	if err != nil {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, err.Error())
		return false, false, nil
	}

	cm.App.discordClient.ChannelMessageSend(message.ChannelID, "User un-verified")

	return true, false, err
}
