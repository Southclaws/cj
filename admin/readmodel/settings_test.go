package readmodel

import (
	"testing"

	"github.com/Southclaws/cj/storage"
)

type fakeStorer struct {
	storage.Memory
	settings    storage.GuildSettings
	auditEvents []storage.AuditEvent
}

func (f *fakeStorer) GetGuildSettings() (storage.GuildSettings, error) {
	return f.settings, nil
}

func (f *fakeStorer) SetGuildSettings(s storage.GuildSettings) error {
	f.settings = s
	return nil
}

func (f *fakeStorer) RecordAuditEvent(event storage.AuditEvent) error {
	f.auditEvents = append(f.auditEvents, event)
	return nil
}

func TestSettingsProviderGetNormalizesNilUserIDs(t *testing.T) {
	storer := &fakeStorer{settings: storage.GuildSettings{}}
	provider := NewSettingsProvider(storer)

	got, err := provider.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.LTFUserIDs == nil {
		t.Fatal("expected LTFUserIDs to be normalized to an empty slice, got nil")
	}
	if len(got.LTFUserIDs) != 0 {
		t.Fatalf("expected an empty slice, got %v", got.LTFUserIDs)
	}
}

func TestSettingsProviderSetNormalizesNilUserIDsAndAudits(t *testing.T) {
	storer := &fakeStorer{}
	provider := NewSettingsProvider(storer)

	got, err := provider.Set(storage.GuildSettings{AdsChannelID: "123"}, "req-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.LTFUserIDs == nil {
		t.Fatal("expected LTFUserIDs to be normalized to an empty slice, got nil")
	}
	if len(storer.auditEvents) != 1 || storer.auditEvents[0].Outcome != "success" {
		t.Fatalf("expected one successful audit event, got %+v", storer.auditEvents)
	}
}
