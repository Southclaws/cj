package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutReturnsServiceUnavailableWhenHandlerIsSlow(t *testing.T) {
	slow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	handler := Timeout(20 * time.Millisecond)(slow)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	rec := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("timeout took too long to fire: %s", elapsed)
	}
}

func TestTimeoutPassesThroughFastHandlers(t *testing.T) {
	handler := Timeout(time.Second)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestTimeoutForwardsPanicToOuterRecover(t *testing.T) {
	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	handler := Recover(Timeout(time.Second)(panicking))

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
