package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestNetworkGuardRejectsUnknownHost(t *testing.T) {
	handler := NetworkGuard([]string{"localhost"})(okHandler())

	req := httptest.NewRequest(http.MethodGet, "http://evil.example/api/v1/status", nil)
	req.Host = "evil.example"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestNetworkGuardAllowsKnownHost(t *testing.T) {
	handler := NetworkGuard([]string{"localhost"})(okHandler())

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	req.Host = "localhost:8080"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestNetworkGuardAllowsGetWithoutOrigin(t *testing.T) {
	handler := NetworkGuard([]string{"localhost"})(okHandler())

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	req.Host = "localhost"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestNetworkGuardRejectsMutationWithoutOrigin(t *testing.T) {
	handler := NetworkGuard([]string{"localhost"})(okHandler())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/api/v1/actions/foo/execute", nil)
	req.Host = "localhost"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestNetworkGuardRejectsMutationWithMismatchedOrigin(t *testing.T) {
	handler := NetworkGuard([]string{"localhost"})(okHandler())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/api/v1/actions/foo/execute", nil)
	req.Host = "localhost"
	req.Header.Set("Origin", "http://evil.example")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestNetworkGuardAllowsMutationWithMatchingOrigin(t *testing.T) {
	handler := NetworkGuard([]string{"localhost"})(okHandler())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/api/v1/actions/foo/execute", nil)
	req.Host = "localhost"
	req.Header.Set("Origin", "http://localhost")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}
