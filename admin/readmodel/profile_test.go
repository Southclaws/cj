package readmodel

import (
	"testing"

	"github.com/Southclaws/cj/storage"
)

type fakeProfileStorer struct {
	storage.Memory
	user         storage.User
	userFound    bool
	messages     []storage.ChatLog
	messageCount int
	rank         int
}

func (f *fakeProfileStorer) GetUser(discordUserID string) (storage.User, bool, error) {
	return f.user, f.userFound, nil
}

func (f *fakeProfileStorer) GetRecentMessagesForUser(discordUserID string, limit int) ([]storage.ChatLog, error) {
	if len(f.messages) > limit {
		return f.messages[:limit], nil
	}
	return f.messages, nil
}

func (f *fakeProfileStorer) GetUserMessageCount(discordUserID string) (int, error) {
	return f.messageCount, nil
}

func (f *fakeProfileStorer) GetUserRank(discordUserID string) (int, error) {
	return f.rank, nil
}

func notReadyDiscordProvider() *DiscordProvider {
	return NewDiscordProvider(nil)
}

func TestProfileProviderNotFoundWhenNothingIsKnown(t *testing.T) {
	storer := &fakeProfileStorer{}
	provider := NewProfileProvider(storer, notReadyDiscordProvider(), NewUsersProvider(storer, nil), NewLeaderboardsProvider(storer, nil))

	_, found, err := provider.Get("does-not-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected no profile to be found when neither Discord nor CJ storage know this ID")
	}
}

func TestProfileProviderFoundFromCJStorageAlone(t *testing.T) {
	storer := &fakeProfileStorer{
		user:         storage.User{DiscordUserID: "42", ForumUserName: "someone"},
		userFound:    true,
		messages:     []storage.ChatLog{{DiscordUserID: "42", Message: "hi", DiscordChannel: "c1"}},
		messageCount: 5,
		rank:         3,
	}
	provider := NewProfileProvider(storer, notReadyDiscordProvider(), NewUsersProvider(storer, nil), NewLeaderboardsProvider(storer, nil))

	profile, found, err := provider.Get("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected the profile to be found from CJ storage alone, with Discord not ready")
	}
	if profile.Member != nil {
		t.Fatalf("expected no Discord member info when Discord isn't ready, got %+v", profile.Member)
	}
	if profile.User == nil || profile.User.ForumUserName != "someone" {
		t.Fatalf("expected the CJ user record to be populated, got %+v", profile.User)
	}
	if profile.MessageCount != 5 || profile.Rank != 3 {
		t.Fatalf("expected message count 5 and rank 3, got count=%d rank=%d", profile.MessageCount, profile.Rank)
	}
	if len(profile.RecentMessages) != 1 || profile.RecentMessages[0].Message != "hi" {
		t.Fatalf("expected the recent message to be included, got %+v", profile.RecentMessages)
	}
}

func TestProfileProviderSkipsRankWhenNoMessages(t *testing.T) {
	storer := &fakeProfileStorer{userFound: true, user: storage.User{DiscordUserID: "1"}, messageCount: 0, rank: 999}
	provider := NewProfileProvider(storer, notReadyDiscordProvider(), NewUsersProvider(storer, nil), NewLeaderboardsProvider(storer, nil))

	profile, found, err := provider.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected the profile to be found")
	}
	if profile.Rank != 0 {
		t.Fatalf("expected rank 0 (unranked) for a user with no messages, got %d", profile.Rank)
	}
}
