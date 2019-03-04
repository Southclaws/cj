package commands

import (
	"net/http"
	"strings"

	"github.com/Southclaws/cj/storage"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

var cmdUsage string = "USAGE : /wiki [function/callback/article]"

func (cm *CommandManager) commandWiki(
	args string,
	message discordgo.Message, 
	contextual bool,
) (
	context bool,
	err error,
) {
	if len(args) == 0 {
		cm.Discord.ChannelMessageSend(message.ChannelID, cmdUsage)
		return
	}

	if len(message.Mentions) > 0 ||
		strings.Contains(message.Content, "everyone") ||
		strings.Contains(message.Content, "here") ||
		args[0] == '?' {
		return
	}

	var (
		wikiThread []string
		exists bool
	)

	wikiThread, exists = storage.SearchThread(args)

	if !exists {
		cm.Discord.ChannelMessageSend(message.ChannelID, "SA:MP Wiki | "+args+"\n- This article does not exist")
		return 
	}

	if len(wikiThread) != 1 {
		cm.Discord.ChannelMessageSend(message.ChannelID, "What are you looking for?\nSA:MP Wiki | `"+strings.Join(wikiThread, "` & `")+"`")
		return
	}

	wikiURL := "https://wiki.sa-mp.com/wiki/"+wikiThread[0]

	var doc *goquery.Document

	response, err := http.Get(wikiURL)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"Could not retrieve SA:MP wiki article:\nGot unexpected response: "+response.Status+".")
	} else {
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"SA:MP Wiki | "+wikiThread[0]+"\n"+wikiURL)

		doc, err = goquery.NewDocument(wikiURL)
		if err != nil {
			return
		}

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
			if strings.TrimSpace(parameters) != "" {
				wikiContent = wikiContent + "\n**Parameters**\n\t" + parameters
				if examplecode != "" {
					wikiContent = wikiContent + "\n\n**Example Usage**\n```C\n" + examplecode + "```"
				}
			}
		}

		if wikiContent != "" {
			cm.Discord.ChannelMessageSend(
				message.ChannelID,
				wikiContent)
		}

	}

	return false, err
}