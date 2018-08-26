package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandWiki(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {

	if len(args) == 0 {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, "USAGE : /wiki [function/callback/article_name]")
		return false, false, nil
	}

	wikiURL := strings.Replace("http://wiki.sa-mp.com/wiki/"+args, " ", "_", -1)

	response, _ := http.Get(wikiURL)
	bodyText, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, "Could not retrieve SA:MP wiki article:\nGot unexpected response: "+response.Status+".")
	} else if strings.Contains(string(bodyText), "There is currently no text in this page, you can") {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, "SA:MP Wiki | "+args+"\n- This article does not exist")
	} else {
		cm.App.discordClient.ChannelMessageSend(message.ChannelID, "SA:MP Wiki | "+args+"\n"+wikiURL)
	}

	return true, true, nil
}
