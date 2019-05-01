package commands

import (
	"os"
	"strings"

	"github.com/Southclaws/cj/storage"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var (
	cmdUsage         = "USAGE : /wiki [function/callback]"
	errManyThreads   = errors.New("more than one thread could be meant")
	errNoThreadFound = errors.New("no thread found")
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

	if !storage.WikiExists() {
		return
	}

	wikiResult := readWiki(args)

	if wikiResult.Err == errManyThreads {
		cm.Discord.ChannelMessageSend(message.ChannelID, errManyThreads.Error()+"\nThe most similar funcs/callbacks are:\n* __**"+strings.Join(wikiResult.Threads, "**__\n* __**")+"**__")
		return
	} else if wikiResult.Err == errNoThreadFound {
		cm.Discord.ChannelMessageSend(message.ChannelID, errNoThreadFound.Error())
		return
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
		readThread(threadToRead[0])
		//TODO parse markdown, format & return as a &WikiResult{Thread:}
	default:
		var threads []string
		for foo := range threadToRead {
			threads = append(threads, strings.TrimSuffix(strings.ReplaceAll(strings.ReplaceAll(threadToRead[foo], "functions/", "function: "), "callbacks/", "callback: "), ".md"))
		}
		return &WikiReturns{Err: errManyThreads, Threads: threads}
	}

	return &WikiReturns{}
}

func readThread(file string) {

}

func searchThread(threads []string, thread string) (results []string) {
	maxdistance := 3
	for i := range threads {
		origin := strings.TrimSuffix(strings.TrimLeft(strings.TrimLeft(strings.TrimLeft(threads[i], "scripting/"), "functions/"), "callbacks/"), ".md")
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
	cb, err = os.Open("./wiki/scripting/callbacks")
	if err != nil {
		return
	}

	fn, err = os.Open("./wiki/scripting/functions")
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
