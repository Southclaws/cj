package readmodel

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Southclaws/cj/storage"
)

var ErrChatFilterRequired = errors.New("a userId or a search query is required to browse chat history")

const maxGlobalSearchResults = 100

type ChatProvider struct {
	storer  storage.Storer
	discord *DiscordProvider
}

func NewChatProvider(storer storage.Storer, discord *DiscordProvider) *ChatProvider {
	return &ChatProvider{storer: storer, discord: discord}
}

type ChatLogWithPerson struct {
	storage.ChatLog
	Person PersonRef `json:"person"`
}

func (p *ChatProvider) attachPeople(messages []storage.ChatLog) []ChatLogWithPerson {
	out := make([]ChatLogWithPerson, len(messages))
	for i, m := range messages {
		out[i] = ChatLogWithPerson{ChatLog: m, Person: p.discord.PersonRef(m.DiscordUserID)}
	}
	return out
}

func (p *ChatProvider) Messages(discordUserID, query string) ([]ChatLogWithPerson, error) {
	var (
		messages []storage.ChatLog
		err      error
	)
	switch {
	case discordUserID != "" && query != "":
		messages, err = p.storer.SearchMessages(discordUserID, query)
	case discordUserID != "":
		messages, err = p.storer.GetRecentMessagesForUser(discordUserID, maxGlobalSearchResults)
	case query != "":
		messages, err = p.storer.SearchAllMessages(query, maxGlobalSearchResults)
	default:
		return nil, ErrChatFilterRequired
	}
	if err != nil {
		return nil, err
	}
	return p.attachPeople(messages), nil
}

func (p *ChatProvider) GetByID(messageID string) (ChatLogWithPerson, bool, error) {
	message, err := p.storer.GetMessageByID(messageID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ChatLogWithPerson{}, false, nil
	}
	if err != nil {
		return ChatLogWithPerson{}, false, err
	}
	return ChatLogWithPerson{ChatLog: message, Person: p.discord.PersonRef(message.DiscordUserID)}, true, nil
}
