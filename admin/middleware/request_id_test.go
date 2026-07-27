package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDGeneratedAndAvailableInContext(t *testing.T) {
	var seen string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if seen == "" {
		t.Fatal("expected a request ID in context")
	}
	if rec.Header().Get("X-Request-Id") != seen {
		t.Fatalf("response header %q does not match context value %q", rec.Header().Get("X-Request-Id"), seen)
	}
}

func TestRequestIDReusesIncomingHeader(t *testing.T) {
	handler := RequestID(okHandler())

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	req.Header.Set("X-Request-Id", "client-supplied-id")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Request-Id") != "client-supplied-id" {
		t.Fatalf("got %q, want client-supplied-id", rec.Header().Get("X-Request-Id"))
	}
}
