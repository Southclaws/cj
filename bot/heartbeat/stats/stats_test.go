package stats

import (
	"testing"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

type fakeStorer struct {
	storage.Memory
	guildSettings storage.GuildSettings
}

func (f *fakeStorer) GetGuildSettings() (storage.GuildSettings, error) {
	return f.guildSettings, nil
}

func TestAnnounceNoOpsWhenLeaderboardChannelUnconfigured(t *testing.T) {
	a := &Aggregator{}
	if _, err := a.Init(&types.Config{}, nil, &fakeStorer{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.leaderboardChannelID != "" {
		t.Fatalf("expected empty leaderboard channel, got %q", a.leaderboardChannelID)
	}

	if err := a.announce(); err != nil {
		t.Fatalf("expected announce to no-op when unconfigured, got error: %v", err)
	}
}

func TestInitResolvesConfiguredLeaderboardChannel(t *testing.T) {
	a := &Aggregator{}
	store := &fakeStorer{guildSettings: storage.GuildSettings{LeaderboardChannelID: "99"}}
	if _, err := a.Init(&types.Config{}, nil, store, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.leaderboardChannelID != "99" {
		t.Fatalf("expected leaderboard channel 99, got %q", a.leaderboardChannelID)
	}
}
