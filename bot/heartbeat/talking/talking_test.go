package talking

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

func TestRegisterReturnsNoActionsWhenUnconfigured(t *testing.T) {
	talk := &Talk{}
	if _, err := talk.Init(&types.Config{}, nil, &fakeStorer{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actions := talk.Register(); len(actions) != 0 {
		t.Fatalf("expected no actions when channel/users are unconfigured, got %d", len(actions))
	}
}

func TestRegisterReturnsActionWhenConfigured(t *testing.T) {
	talk := &Talk{}
	store := &fakeStorer{guildSettings: storage.GuildSettings{
		LTFChannelID: "1",
		LTFUserIDs:   []string{"2", "3"},
	}}
	if _, err := talk.Init(&types.Config{}, nil, store, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actions := talk.Register()
	if len(actions) != 1 {
		t.Fatalf("expected exactly one action when fully configured, got %d", len(actions))
	}
}
