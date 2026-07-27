package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionRunRecordAndList(t *testing.T) {
	id, err := api.RecordActionRun(ActionRun{
		Action:     "services.test",
		Risk:       "safe",
		Mode:       "execute",
		Outcome:    "success",
		Summary:    "mongodb: healthy",
		DurationMS: 12,
		RequestID:  "test-run-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, id)

	runs, err := api.ListActionRuns(10)
	require.NoError(t, err)
	require.NotEmpty(t, runs)

	latest := runs[0]
	assert.Equal(t, "services.test", latest.Action)
	assert.Equal(t, "success", latest.Outcome)
	assert.False(t, latest.Timestamp.IsZero())

	fetched, found, err := api.GetActionRun(id)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "services.test", fetched.Action)

	_, found, err = api.GetActionRun("000000000000000000000000")
	require.NoError(t, err)
	assert.False(t, found)

	_, found, err = api.GetActionRun("not-a-valid-object-id")
	require.NoError(t, err)
	assert.False(t, found)
}
