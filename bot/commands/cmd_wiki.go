package commands

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/russross/blackfriday.v2"

	"github.com/Southclaws/cj/types"
)

var (
	cmdUsage             = "USAGE : /wiki [function/callback]"
	errManyThreads       = errors.New("more than one thread could be meant")
	errNoThreadFound     = errors.New("no thread found")
	errCouldntReadThread = errors.New("couldn't read wiki file")
)

// WikiReturns holds structure of expected results from @readWiki
type WikiReturns struct {
	Err     error
	Threads []string
	Thread  string
}

func (cm *CommandManager) commandWiki(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
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

	if !cm.Storage.WikiExists() {
		return
	}

	wikiResult := readWiki(args)

	if wikiResult.Err == errManyThreads {
		cm.Discord.ChannelMessageSend(message.ChannelID, errManyThreads.Error()+"\nThe most similar funcs/callbacks are:\n* __**"+strings.Join(wikiResult.Threads, "**__\n* __**")+"**__")
		return
	} else if wikiResult.Err == errNoThreadFound {
		// TODO: url shouldn't be hardcoded
		cm.Discord.ChannelMessageSend(message.ChannelID, "If you think this page should exist, please open a pull request or issue here: "+"<https://github.com/openmultiplayer/wiki>")
		return
	} else if wikiResult.Err == errCouldntReadThread {
		cm.Discord.ChannelMessageSend(message.ChannelID, errCouldntReadThread.Error())
	} else if wikiResult.Err == nil {
		cm.Discord.ChannelMessageSend(message.ChannelID, wikiResult.Thread)
	} else {
		err = wikiResult.Err
	}

	return false, err
}

func readWiki(article string) *WikiReturns {
	files, err := getWikiFiles()
	if err != nil {
		return &WikiReturns{Err: err}
	}

	threadToRead := searchThread(files, article)

	switch len(threadToRead) {
	case 0:
		return &WikiReturns{Err: errNoThreadFound}
	case 1:
		formatedInfo, readErr := readThread(threadToRead[0])
		if readErr == errCouldntReadThread {
			return &WikiReturns{Err: errCouldntReadThread}
		}
		return &WikiReturns{Err: nil, Thread: formatedInfo}
	default:
		var threads []string
		for foo := range threadToRead {
			threads = append(threads, threadName(threadToRead[foo]))
		}
		return &WikiReturns{Err: errManyThreads, Threads: threads}
	}
}

func readThread(file string) (string, error) {
	wikiFile, openErr := os.Open(filepath.Join(".", "wiki", "scripting", file))
	if openErr != nil {
		return "", errCouldntReadThread
	}
	defer wikiFile.Close()

	wikiAsByte, readErr := ioutil.ReadAll(wikiFile)
	if readErr != nil {
		return "", errCouldntReadThread
	}

	output := blackfriday.Run(wikiAsByte)
	html := bluemonday.UGCPolicy().SanitizeBytes(output)

	doc, docErr := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if docErr != nil {
		log.Fatal(docErr)
	}

	header := "wiki.open.mp | __" + threadName(file) + "__\n<https://wiki.open.mp/scripting/" + strings.ReplaceAll(file, filepath.Ext(file), ".html>")
	description := "**Description**\n\t" + doc.Find(`h2:contains("Description")`).Next().Text()
	parameters := "**Parameters**"
	relatedFuncs := "**Related Functions**"
	example := "**Example Usage**\n```c\n" + doc.Find("pre code").Text() + "\n```"

	var (
		selectionCache       *goquery.Selection
		parametersAddition   string
		relatedFuncsAddition string
	)

	lastULClass := doc.Find("ul").Last()
	lastULClass.Find("li").Each(func(i int, s *goquery.Selection) {
		//selectionCache = s.First()
		relatedFuncsAddition = relatedFuncsAddition + "\n\t" + s.Text()
	})

	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		selectionCache = s.Find("td").First()
		parametersAddition = parametersAddition + "\n\t`" + selectionCache.Text() + "`\t_*" + selectionCache.Next().Text() + "*_"
	})

	formatedText := header + "\n\n" + description

	if len(parametersAddition) > 0 {
		formatedText = formatedText + "\n\n" + parameters + parametersAddition
	}

	if len(relatedFuncsAddition) > 0 {
		formatedText = formatedText + "\n\n" + relatedFuncs + relatedFuncsAddition
	}

	formatedText = formatedText + "\n\n" + example

	return formatedText, nil
}

func searchThread(threads []string, thread string) (results []string) {
	maxdistance := 3
	for i := range threads {
		origin := threadName(threads[i])
		dist := levenshtein.DistanceForStrings(
			[]rune(strings.ToLower(origin)),
			[]rune(strings.ToLower(thread)),
			levenshtein.DefaultOptions,
		)
		if dist == 0 {
			return []string{threads[i]}
		}
		if dist <= maxdistance {
			results = append(results, threads[i])
		}
	}
	return
}

func getWikiFiles() (files []string, err error) {
	var (
		cb  *os.File
		fn  *os.File
		tmp []string
	)
	cb, err = os.Open(filepath.Join(".", "wiki", "scripting", "callbacks"))
	if err != nil {
		return
	}

	fn, err = os.Open(filepath.Join(".", "wiki", "scripting", "functions"))
	if err != nil {
		cb.Close()
		return
	}

	defer cb.Close()
	defer fn.Close()

	tmp, err = cb.Readdirnames(0)
	if err != nil {
		return
	}

	for foo := range tmp {
		tmp[foo] = "callbacks/" + tmp[foo]
	}

	files = tmp

	tmp, err = fn.Readdirnames(0)
	if err != nil {
		return
	}

	for foo := range tmp {
		tmp[foo] = "functions/" + tmp[foo]
	}

	files = append(files, tmp...)

	return
}

func threadName(threadPath string) string {
	return strings.TrimSuffix(filepath.Base(threadPath), filepath.Ext(threadPath))
}
