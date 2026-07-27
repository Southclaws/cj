package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditEventRecordAndList(t *testing.T) {
	before := map[string]any{"ads_channel_id": ""}
	after := map[string]any{"ads_channel_id": "42"}

	require.NoError(t, api.RecordAuditEvent(AuditEvent{
		Action:     "update_guild_settings",
		TargetType: "guild_settings",
		Before:     before,
		After:      after,
		Outcome:    "success",
		RequestID:  "test-request-1",
	}))

	events, err := api.ListAuditEvents(10)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	latest := events[0]
	assert.Equal(t, "update_guild_settings", latest.Action)
	assert.Equal(t, "success", latest.Outcome)
	assert.False(t, latest.Timestamp.IsZero())
}
