package admod

import (
	"testing"

	"github.com/bwmarrin/discordgo"

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

func TestWatcherDisabledWhenAdsChannelUnconfigured(t *testing.T) {
	w := &Watcher{}
	if err := w.Init(&types.Config{}, nil, &fakeStorer{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.channel != "" {
		t.Fatalf("expected empty channel, got %q", w.channel)
	}

	if err := w.OnMessage(discordgo.Message{ChannelID: "anything"}); err != nil {
		t.Fatalf("expected disabled watcher to no-op, got error: %v", err)
	}
}

func TestWatcherUsesConfiguredAdsChannel(t *testing.T) {
	w := &Watcher{}
	store := &fakeStorer{guildSettings: storage.GuildSettings{AdsChannelID: "42"}}
	if err := w.Init(&types.Config{}, nil, store, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.channel != "42" {
		t.Fatalf("expected channel 42, got %q", w.channel)
	}
}
