package readmodel

import "github.com/Southclaws/cj/storage"

const (
	defaultLeaderboardSize = 10
	maxLeaderboardSize     = 100
)

type LeaderboardsProvider struct {
	storer  storage.Storer
	discord *DiscordProvider
}

func NewLeaderboardsProvider(storer storage.Storer, discord *DiscordProvider) *LeaderboardsProvider {
	return &LeaderboardsProvider{storer: storer, discord: discord}
}

type TopMessagesEntry struct {
	Person   PersonRef `json:"person"`
	Messages int       `json:"messages"`
}

type TopReactionsEntry struct {
	Person   PersonRef `json:"person"`
	Counter  int       `json:"counter"`
	Reaction string    `json:"reaction"`
}

func (p *LeaderboardsProvider) TopMessages(limit int) ([]TopMessagesEntry, error) {
	if limit <= 0 || limit > maxLeaderboardSize {
		limit = defaultLeaderboardSize
	}
	result, err := p.storer.GetTopMessages(limit)
	if err != nil {
		return nil, err
	}

	out := make([]TopMessagesEntry, len(result))
	for i, entry := range result {
		out[i] = TopMessagesEntry{Person: p.discord.PersonRef(entry.User), Messages: entry.Messages}
	}
	return out, nil
}

func (p *LeaderboardsProvider) TopReactions(limit int, reaction string) ([]TopReactionsEntry, error) {
	if limit <= 0 || limit > maxLeaderboardSize {
		limit = defaultLeaderboardSize
	}
	result, err := p.storer.GetTopReactions(limit, reaction)
	if err != nil {
		return nil, err
	}

	out := make([]TopReactionsEntry, len(result))
	for i, entry := range result {
		out[i] = TopReactionsEntry{
			Person:   p.discord.PersonRef(entry.UserID),
			Counter:  entry.Counter,
			Reaction: entry.Reaction,
		}
	}
	return out, nil
}

func (p *LeaderboardsProvider) UserRank(discordUserID string) (int, error) {
	return p.storer.GetUserRank(discordUserID)
}
