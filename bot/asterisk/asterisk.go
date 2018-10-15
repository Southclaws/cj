package asterisk

import (
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	"github.com/texttheater/golang-levenshtein/levenshtein"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

// Asterisk facilitates replacements of misspelt words by correcting the word in
// the next message with an asterisk proceeding it.
type Asterisk struct {
	Config  *types.Config
	Discord *discord.Session
	Storage *storage.API
	Forum   *forum.ForumClient

	cache *cache.Cache
}

// Init creates a command manager for the app
func (ast *Asterisk) Init(
	config *types.Config,
	discord *discord.Session,
	api *storage.API,
	fc *forum.ForumClient,
) (err error) {
	ast.Config = config
	ast.Storage = api
	ast.Discord = discord
	ast.Forum = fc

	ast.cache = cache.New(time.Minute, time.Minute)

	return
}

var matchAsteriskCorrection = regexp.MustCompile(`^\w+\*$`)

// OnMessage checks if:
// - the message simply contains a word followed immediately by an asterisk
//   and
// - the previous message contained a word very similar
// if so, the previous message is edited and the current, deleted.
func (ast *Asterisk) OnMessage(message discordgo.Message) (err error) {
	ast.doCorrection(message)
	ast.cache.Set(message.Author.ID, message, time.Minute)
	return
}

func (ast *Asterisk) doCorrection(message discordgo.Message) {
	if !matchAsteriskCorrection.MatchString(message.Content) {
		return
	}

	lastRaw, exists := ast.cache.Get(message.Author.ID)
	if !exists {
		return
	}

	last, ok := lastRaw.(discordgo.Message)
	if !ok {
		return
	}

	correction := message.Content[:len(message.Content)-1]

	// get all words in previous message
	words := strings.Split(last.Content, " ")
	target := -1

	// find the first word that has an edit distance of 3 or below
	for i, word := range words {
		if levenshtein.DistanceForStrings([]rune(word), []rune(correction), levenshtein.DefaultOptions) <= 3 &&
			word != correction {
			target = i
			break
		}
	}

	if target == -1 {
		return
	}

	// apply the correction
	words[target] = correction

	// reassemble the sentence and edit original message
	// nolint:errcheck
	{
		ast.Discord.S.ChannelMessageDelete(last.ChannelID, last.ID)
		ast.Discord.S.ChannelMessageDelete(message.ChannelID, message.ID)
		ast.Discord.ChannelMessageSend(
			message.ChannelID,
			last.Author.Username+" meant to say: "+strings.Join(words, " "))
	}
}
