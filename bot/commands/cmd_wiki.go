package commands

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (cm *CommandManager) commandWiki(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	if len(args) == 0 {
		_, err = cm.Discord.ChannelMessageSend(message.ChannelID, "USAGE : /wiki [function/callback/article_name]")
		return
	}

	if len(message.Mentions) > 0 ||
		strings.Contains(message.Content, "everyone") ||
		strings.Contains(message.Content, "here") {
		return
	}

	wikiURL := strings.Replace("http://wiki.sa-mp.com/wiki/"+args, " ", "_", -1)

	response, err := http.Get(wikiURL)
	if err != nil {
		return
	}

	bodyText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		_, err = cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"Could not retrieve SA:MP wiki article:\nGot unexpected response: "+response.Status+".")
	} else if strings.Contains(string(bodyText), "There is currently no text in this page, you can") ||
		strings.Contains(string(bodyText), "The requested page title was invalid, empty") {
		_, err = cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"SA:MP Wiki | "+args+"\n- This article does not exist")
	} else {
		_, err = cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"SA:MP Wiki | "+args+"\n"+wikiURL)
	}

	return false, err
}
