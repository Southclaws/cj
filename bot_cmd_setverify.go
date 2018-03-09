package main

import "github.com/bwmarrin/discordgo"

func commandSetVerify(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	if len(message.Mentions) == 0 {
		return false, false, nil
	}

	target := message.Mentions[0]

	err := cm.App.accounts.Insert(&User{
		DiscordUserID:    target.ID,
		ForumUserID:      "",
		VerificationCode: "",
		ForumUserName:    "",
	})

	return true, false, err
}
