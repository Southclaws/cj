package admin

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

func TestNewDisabledByDefault(t *testing.T) {
	cfg := &types.Config{}

	server, err := New(cfg, "test", Dependencies{Storage: &storage.Memory{}})

	require.Nil(t, server)
	assert.True(t, errors.Is(err, ErrDisabled))
}

func TestNewInvalidExplicitConfigDoesNotPanic(t *testing.T) {
	cfg := &types.Config{
		DashboardEnabled:    true,
		DashboardListenAddr: "127.0.0.1:0",
	}

	server, err := New(cfg, "test", Dependencies{Storage: &storage.Memory{}})

	require.Nil(t, server)
	require.Error(t, err)
	assert.False(t, errors.Is(err, ErrDisabled))
}

func TestNewDefaultsToLoopbackListener(t *testing.T) {
	cfg := &types.Config{DashboardEnabled: true}

	server, err := New(cfg, "test", Dependencies{Storage: &storage.Memory{}})
	require.NoError(t, err)
	require.NotNil(t, server)
	assert.Equal(t, "loopback", server.addrSource)
	assert.Equal(t, "127.0.0.1:0", server.addr)
}

func TestServerStartServesStatusAndShutsDownCleanly(t *testing.T) {
	cfg := &types.Config{DashboardEnabled: true}

	server, err := New(cfg, "v1.2.3", Dependencies{Storage: &storage.Memory{}})
	require.NoError(t, err)

	require.NoError(t, server.Start(context.Background()))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		assert.NoError(t, server.Shutdown(ctx))
	}()

	resp, err := http.Get("http://" + server.Addr() + "/api/v1/status")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "v1.2.3")
}
