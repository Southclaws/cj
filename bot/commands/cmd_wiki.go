package commands

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Southclaws/cj/storage"
	"github.com/bwmarrin/discordgo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var cmdUsage = "USAGE : /wiki [function/callback/article]"

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
	} else if len(args) < 3 {
		cm.Discord.ChannelMessageSend(message.ChannelID, "Query must be 3 characters or more")
		return
	}

	if len(message.Mentions) > 0 ||
		strings.Contains(message.Content, "everyone") ||
		strings.Contains(message.Content, "here") ||
		args[0] == '?' {
		return
	}

	var (
		wikiThread  []string
		wikiURL     string
		articleName = strings.Replace(args, " ", "_", -1)
	)

	wikiThread = storage.SearchThread(args)

	if len(wikiThread) > 0 {
		if len(wikiThread) != 1 {
			cm.Discord.ChannelMessageSend(message.ChannelID, "What are you looking for?\nResults from *SA:MP Wiki*:\n__"+strings.Join(wikiThread, "__\n__")+"__")
			return
		}
		dist := levenshtein.DistanceForStrings(
			[]rune(strings.ToLower(wikiThread[0])),
			[]rune(strings.ToLower(articleName)),
			levenshtein.DefaultOptions,
		)
		if dist <= 2 {
			articleName = wikiThread[0]
		}

	}

	wikiURL = "https://wiki.sa-mp.com/wiki/" + articleName

	var doc *goquery.Document

	response, err := http.Get(wikiURL)
	if err != nil {
		return
	}

	bodyText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"Could not retrieve SA:MP wiki article:\nGot unexpected response: "+response.Status+".")
	} else if strings.Contains(string(bodyText), "There is currently no text in this page, you can") ||
		strings.Contains(string(bodyText), "The requested page title was invalid, empty") {
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"SA:MP Wiki | "+args+"\n- This article does not exist")
	} else {
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			"SA:MP Wiki | "+articleName+"\n"+wikiURL)

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
