package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoverCatchesPanicAndReturnsSafeError(t *testing.T) {
	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something exploded, credentials=supersecret")
	})

	handler := Recover(panicking)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if strings.Contains(rec.Body.String(), "supersecret") {
		t.Fatalf("panic detail leaked into response body: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "internal_error") {
		t.Fatalf("expected safe error code in body, got: %s", rec.Body.String())
	}
}

func TestRecoverPassesThroughWhenNoPanic(t *testing.T) {
	handler := Recover(okHandler())

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}
