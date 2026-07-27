package commands

import (
	"testing"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

type fakeConfigStorer struct {
	storage.Memory
	saved map[string]types.CommandSettings
}

func (f *fakeConfigStorer) SetCommandSettings(command string, settings types.CommandSettings) error {
	if f.saved == nil {
		f.saved = map[string]types.CommandSettings{}
	}
	f.saved[command] = settings
	return nil
}

func TestListCommandNames(t *testing.T) {
	cm := &CommandManager{Commands: []Command{{Name: "/a"}, {Name: "/b"}}}

	names := cm.ListCommandNames()
	if len(names) != 2 || names[0] != "/a" || names[1] != "/b" {
		t.Fatalf("unexpected names: %+v", names)
	}
}

func TestGetCommandConfigUnknownCommand(t *testing.T) {
	cm := &CommandManager{Commands: []Command{{Name: "/a"}}}

	if _, err := cm.GetCommandConfig("/does-not-exist"); err == nil {
		t.Fatal("expected an error for an unregistered command name")
	}
}

func TestGetAndSetCommandConfigRoundTrips(t *testing.T) {
	storer := &fakeConfigStorer{}
	cm := &CommandManager{
		Storage:  storer,
		Commands: []Command{{Name: "/searchmessage", Settings: types.CommandSettings{Roles: []string{"111"}}}},
	}

	settings, err := cm.GetCommandConfig("/searchmessage")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(settings.Roles) != 1 || settings.Roles[0] != "111" {
		t.Fatalf("unexpected settings: %+v", settings)
	}

	updated, err := cm.SetCommandConfig("/searchmessage", types.CommandSettings{Roles: []string{"222", "333"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated.Roles) != 2 {
		t.Fatalf("expected the returned settings to reflect the update, got %+v", updated)
	}

	again, err := cm.GetCommandConfig("/searchmessage")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(again.Roles) != 2 || again.Roles[0] != "222" {
		t.Fatalf("expected the in-memory command list to reflect the update, got %+v", again)
	}

	if saved, ok := storer.saved["/searchmessage"]; !ok || len(saved.Roles) != 2 {
		t.Fatalf("expected the update to be persisted to storage, got %+v", storer.saved)
	}
}

func TestSetCommandConfigUnknownCommand(t *testing.T) {
	cm := &CommandManager{Commands: []Command{{Name: "/a"}}}

	if _, err := cm.SetCommandConfig("/does-not-exist", types.CommandSettings{}); err == nil {
		t.Fatal("expected an error for an unregistered command name")
	}
}
