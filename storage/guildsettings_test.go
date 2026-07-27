package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuildSettingsRoundTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := api.settings.DeleteMany(ctx, guildSettingsFilter)
	require.NoError(t, err)
	api.cache.Delete(guildSettingsCacheKey)

	empty, err := api.GetGuildSettings()
	require.NoError(t, err)
	assert.Empty(t, empty.AdsChannelID)
	assert.Empty(t, empty.LTFUserIDs)

	want := GuildSettings{
		SearchMessageChannelID: "1",
		ErrorReportChannelID:   "2",
		LeaderboardChannelID:   "3",
		AdsChannelID:           "5",
		LTFChannelID:           "6",
		LTFUserIDs:             []string{"7", "8"},
	}
	require.NoError(t, api.SetGuildSettings(want))

	got, err := api.GetGuildSettings()
	require.NoError(t, err)
	assert.Equal(t, want, got)

	want.AdsChannelID = "9"
	require.NoError(t, api.SetGuildSettings(want))

	got, err = api.GetGuildSettings()
	require.NoError(t, err)
	assert.Equal(t, "9", got.AdsChannelID)
}

func TestReportErrorNoOpsWithoutConfiguredChannel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := api.settings.DeleteMany(ctx, guildSettingsFilter)
	require.NoError(t, err)
	api.cache.Delete(guildSettingsCacheKey)

	assert.NotPanics(t, func() {
		api.reportError(nil, "should not attempt to send")
	})
}
