package readmodel

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/Southclaws/cj/storage"
)

type SettingsProvider struct {
	storer storage.Storer
}

func NewSettingsProvider(storer storage.Storer) *SettingsProvider {
	return &SettingsProvider{storer: storer}
}

func (p *SettingsProvider) Get() (storage.GuildSettings, error) {
	settings, err := p.storer.GetGuildSettings()
	if err != nil {
		return settings, err
	}
	return normalizeGuildSettings(settings), nil
}

func (p *SettingsProvider) Set(next storage.GuildSettings, requestID string) (storage.GuildSettings, error) {
	before, _ := p.storer.GetGuildSettings()

	err := p.storer.SetGuildSettings(next)

	outcome := "success"
	if err != nil {
		outcome = "failure"
	}

	auditErr := p.storer.RecordAuditEvent(storage.AuditEvent{
		Action:     "update_guild_settings",
		TargetType: "guild_settings",
		Before:     settingsToMap(before),
		After:      settingsToMap(next),
		Outcome:    outcome,
		RequestID:  requestID,
	})
	if auditErr != nil {
		zap.L().Error("failed to record audit event", zap.Error(auditErr), zap.String("request_id", requestID))
	}

	if err != nil {
		return storage.GuildSettings{}, err
	}
	return normalizeGuildSettings(next), nil
}

func normalizeGuildSettings(s storage.GuildSettings) storage.GuildSettings {
	if s.LTFUserIDs == nil {
		s.LTFUserIDs = []string{}
	}
	return s
}

func settingsToMap(s storage.GuildSettings) map[string]any {
	b, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}
