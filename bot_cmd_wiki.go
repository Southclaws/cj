package main

import (
	"fmt"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/bwmarrin/discordgo"
)

func commandWiki(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	
	if args == nil {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, "USAGE : /wiki [function/callback/article_name]")
		return false, false, nil
	}

	wikiUrl := strings.Replace("http://wiki.sa-mp.com/wiki/" + args, " ", "_", -1)

	response, _ := http.Get(wikiUrl)
	bodyText, _ := ioutil.ReadAll(response.Body)

	if strings.Contains(string(bodyText), "There is currently no text in this page, you can") {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, "SA:MP Wiki | " + args + "\n- This article does not exist")
	} else {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, "SA:MP Wiki | " + args + "\n" + wikiUrl)
	} 
	
	return true, true, nil
}