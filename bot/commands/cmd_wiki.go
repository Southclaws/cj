package commands

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"strings"
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
		strings.Contains(message.Content, "here") ||
		args[0] == '?' {
		return
	}

	wikiURL := strings.Replace("https://wiki.sa-mp.com/wiki/"+args, " ", "_", -1)

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

		doc, err := goquery.NewDocument(wikiURL)
		if err == nil {
			var wikiContent string
			description := strings.TrimSpace(doc.Find(".description").Text())
			parameters := strings.TrimSpace(doc.Find(".parameters").Text()) + "\n"
			var First *goquery.Selection

			doc.Find(".param").Each(func(i int, selection *goquery.Selection) {
				First = selection.Find("td").First()
				parameters = parameters + "\n\t\t`" + First.Text() + "`\t" + First.Next().Text()

			})
			examplecode := strings.TrimSpace(doc.Find(".pawn").Text())

			if description != "" {
				wikiContent = "**Description**\n\t" + description
			}
			if parameters != "" {
				wikiContent = wikiContent + "\n**Parameters**\n\t" + parameters
			}
			if examplecode != "" {
				wikiContent = wikiContent + "\n\n**Example Usage**\n```C\n" + examplecode + "```"
			}

			if wikiContent != "" {
				_, err = cm.Discord.ChannelMessageSend(
					message.ChannelID,
					wikiContent)
			}
		}

	}

	return false, err
}
