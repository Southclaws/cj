package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserDoesNotCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := api.accounts.DeleteOne(ctx, map[string]any{"discord_user_id": "does-not-exist"})
	require.NoError(t, err)

	_, found, err := api.GetUser("does-not-exist")
	require.NoError(t, err)
	assert.False(t, found)

	count, err := api.accounts.CountDocuments(ctx, map[string]any{"discord_user_id": "does-not-exist"})
	require.NoError(t, err)
	assert.Zero(t, count, "GetUser must not create a document as a side effect")
}

func TestListUsersPaginates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := api.accounts.DeleteMany(ctx, map[string]any{})
	require.NoError(t, err)

	_, err = api.accounts.InsertMany(ctx, []any{
		User{DiscordUserID: "u1"},
		User{DiscordUserID: "u2"},
		User{DiscordUserID: "u3"},
	})
	require.NoError(t, err)

	page1, total, err := api.ListUsers("", 2, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, page1, 2)
	assert.Equal(t, "u1", page1[0].DiscordUserID)
	assert.Equal(t, "u2", page1[1].DiscordUserID)

	page2, total, err := api.ListUsers("", 2, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, page2, 1)
	assert.Equal(t, "u3", page2[0].DiscordUserID)
}

func TestListUsersFiltersByQuery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := api.accounts.DeleteMany(ctx, map[string]any{})
	require.NoError(t, err)

	_, err = api.accounts.InsertMany(ctx, []any{
		User{DiscordUserID: "111", ForumUserName: "Alice"},
		User{DiscordUserID: "222", ForumUserName: "Bob"},
		User{DiscordUserID: "333", BurgerUserName: "alicia"},
	})
	require.NoError(t, err)

	byID, total, err := api.ListUsers("222", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, byID, 1)
	assert.Equal(t, "222", byID[0].DiscordUserID)

	byName, total, err := api.ListUsers("alic", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, byName, 2)

	none, total, err := api.ListUsers("does-not-exist", 10, 0)
	require.NoError(t, err)
	assert.Zero(t, total)
	assert.Empty(t, none)
}
