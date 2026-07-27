package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Southclaws/cj/admin/readmodel"
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

func testRouterWithStorer(storer storage.Storer) http.Handler {
	return NewRouter(Config{
		AllowedHosts:   []string{"localhost"},
		RequestTimeout: time.Second,
		MaxBodyBytes:   1 << 20,
		Status:         readmodel.NewStatusProvider("v9.9.9", time.Now()),
		Settings:       readmodel.NewSettingsProvider(storer),
	})
}

func TestSettingsGetReturnsCurrentValue(t *testing.T) {
	storer := &fakeStorer{settings: storage.GuildSettings{AdsChannelID: "123"}}
	router := testRouterWithStorer(storer)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/data/settings", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got storage.GuildSettings
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if got.AdsChannelID != "123" {
		t.Fatalf("got AdsChannelID %q, want 123", got.AdsChannelID)
	}
}

func TestSettingsPutPersistsAndAudits(t *testing.T) {
	storer := &fakeStorer{}
	router := testRouterWithStorer(storer)

	body, _ := json.Marshal(storage.GuildSettings{
		AdsChannelID: "999999999999999999",
		LTFUserIDs:   []string{"111111111111111111"},
	})

	req := httptest.NewRequest(http.MethodPut, "http://localhost/api/v1/data/settings", bytes.NewReader(body))
	req.Header.Set("Origin", "http://localhost")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if storer.settings.AdsChannelID != "999999999999999999" {
		t.Fatalf("expected settings to be persisted, got %+v", storer.settings)
	}
	if len(storer.auditEvents) != 1 {
		t.Fatalf("expected exactly one audit event, got %d", len(storer.auditEvents))
	}
	if storer.auditEvents[0].Outcome != "success" {
		t.Fatalf("expected a successful audit outcome, got %q", storer.auditEvents[0].Outcome)
	}
}

func TestSettingsPutRejectsNonNumericIDs(t *testing.T) {
	storer := &fakeStorer{}
	router := testRouterWithStorer(storer)

	body, _ := json.Marshal(storage.GuildSettings{AdsChannelID: "not-a-snowflake"})

	req := httptest.NewRequest(http.MethodPut, "http://localhost/api/v1/data/settings", bytes.NewReader(body))
	req.Header.Set("Origin", "http://localhost")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if storer.settings.AdsChannelID != "" {
		t.Fatalf("expected invalid input to be rejected before persisting, got %+v", storer.settings)
	}
}

func TestSettingsPutRejectsMissingOrigin(t *testing.T) {
	storer := &fakeStorer{}
	router := testRouterWithStorer(storer)

	body, _ := json.Marshal(storage.GuildSettings{AdsChannelID: "123"})

	req := httptest.NewRequest(http.MethodPut, "http://localhost/api/v1/data/settings", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusForbidden)
	}
}
