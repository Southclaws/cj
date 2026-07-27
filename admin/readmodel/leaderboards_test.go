package readmodel

import (
	"testing"

	"github.com/Southclaws/cj/storage"
)

type fakeLeaderboardsStorer struct {
	storage.Memory
	topMessages    storage.TopMessages
	topReactions   []storage.TopReactionEntry
	rank           int
	lastLimit      int
	lastReaction   string
	lastRankUserID string
}

func (f *fakeLeaderboardsStorer) GetTopMessages(top int) (storage.TopMessages, error) {
	f.lastLimit = top
	return f.topMessages, nil
}

func (f *fakeLeaderboardsStorer) GetTopReactions(top int, reaction string) ([]storage.TopReactionEntry, error) {
	f.lastLimit = top
	f.lastReaction = reaction
	return f.topReactions, nil
}

func (f *fakeLeaderboardsStorer) GetUserRank(discordUserID string) (int, error) {
	f.lastRankUserID = discordUserID
	return f.rank, nil
}

func TestLeaderboardsProviderTopMessagesClampsLimit(t *testing.T) {
	storer := &fakeLeaderboardsStorer{}
	provider := NewLeaderboardsProvider(storer, nil)

	if _, err := provider.TopMessages(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if storer.lastLimit != defaultLeaderboardSize {
		t.Fatalf("expected a zero limit to default to %d, got %d", defaultLeaderboardSize, storer.lastLimit)
	}

	if _, err := provider.TopMessages(10000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if storer.lastLimit != defaultLeaderboardSize {
		t.Fatalf("expected an oversized limit to default to %d, got %d", defaultLeaderboardSize, storer.lastLimit)
	}
}

func TestLeaderboardsProviderTopMessagesNormalizesNil(t *testing.T) {
	storer := &fakeLeaderboardsStorer{topMessages: nil}
	provider := NewLeaderboardsProvider(storer, nil)

	result, err := provider.TopMessages(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected a nil result to be normalized to an empty slice")
	}
}

func TestLeaderboardsProviderTopReactionsPassesReactionThrough(t *testing.T) {
	storer := &fakeLeaderboardsStorer{}
	provider := NewLeaderboardsProvider(storer, nil)

	if _, err := provider.TopReactions(5, "thumbsup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if storer.lastReaction != "thumbsup" {
		t.Fatalf("expected the reaction to be passed through, got %q", storer.lastReaction)
	}
}

func TestLeaderboardsProviderUserRank(t *testing.T) {
	storer := &fakeLeaderboardsStorer{rank: 3}
	provider := NewLeaderboardsProvider(storer, nil)

	rank, err := provider.UserRank("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rank != 3 {
		t.Fatalf("expected rank 3, got %d", rank)
	}
	if storer.lastRankUserID != "42" {
		t.Fatalf("expected the userID to be passed through, got %q", storer.lastRankUserID)
	}
}
