//go:build embed_web

package admin

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

func TestEmbeddedFrontendIsServed(t *testing.T) {
	cfg := &types.Config{DashboardEnabled: true}

	server, err := New(cfg, "test", Dependencies{Storage: &storage.Memory{}})
	require.NoError(t, err)
	require.NoError(t, server.Start(context.Background()))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		assert.NoError(t, server.Shutdown(ctx))
	}()

	base := "http://" + server.Addr()

	shell, err := http.Get(base + "/")
	require.NoError(t, err)
	defer shell.Body.Close()
	body, err := io.ReadAll(shell.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, shell.StatusCode)
	assert.Contains(t, string(body), "CJ Dashboard")

	spaRoute, err := http.Get(base + "/actions")
	require.NoError(t, err)
	defer spaRoute.Body.Close()
	spaBody, err := io.ReadAll(spaRoute.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, spaRoute.StatusCode)
	assert.True(t, strings.Contains(string(spaBody), "CJ Dashboard"), "expected client route to fall back to the SPA shell")

	statusResp, err := http.Get(base + "/api/v1/status")
	require.NoError(t, err)
	defer statusResp.Body.Close()
	assert.Equal(t, http.StatusOK, statusResp.StatusCode)

	configResp, err := http.Get(base + "/api/v1/config/status")
	require.NoError(t, err)
	defer configResp.Body.Close()
	assert.Equal(t, http.StatusOK, configResp.StatusCode)

	missingAPI, err := http.Get(base + "/api/v1/nope")
	require.NoError(t, err)
	defer missingAPI.Body.Close()
	assert.Equal(t, http.StatusNotFound, missingAPI.StatusCode)
}
