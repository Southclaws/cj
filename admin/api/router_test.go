package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/cj/admin/actions"
	"github.com/Southclaws/cj/admin/readmodel"
	"github.com/Southclaws/cj/storage"
)

type stubStorer struct {
	storage.Memory
}

func testRouter() http.Handler {
	return NewRouter(Config{
		AllowedHosts:   []string{"localhost"},
		RequestTimeout: time.Second,
		ActionTimeout:  time.Second,
		Status:         readmodel.NewStatusProvider("v9.9.9", time.Now()),
		Actions:        actions.NewRunner(actions.NewRegistry(), &stubStorer{}),
		ActionRuns:     readmodel.NewActionRunsProvider(&stubStorer{}),
	})
}

type slowAction struct {
	delay time.Duration
}

func (a *slowAction) Name() string        { return "test.slow" }
func (a *slowAction) Description() string { return "sleeps then succeeds" }
func (a *slowAction) Risk() actions.RiskLevel {
	return actions.RiskSafe
}
func (a *slowAction) Preview(ctx context.Context, input json.RawMessage) (actions.Preview, error) {
	return actions.Preview{Summary: "will sleep"}, nil
}
func (a *slowAction) Execute(ctx context.Context, input json.RawMessage) (actions.Result, error) {
	time.Sleep(a.delay)
	return actions.Result{Summary: "done"}, nil
}

func TestRouterServesStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/status", nil)
	rec := httptest.NewRecorder()

	testRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "v9.9.9") {
		t.Fatalf("expected version in body, got: %s", rec.Body.String())
	}
}

func TestRouterRejectsUnknownHost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://evil.example/api/v1/status", nil)
	rec := httptest.NewRecorder()

	testRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestRouterReturnsNotFoundForUnknownRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/v1/nope", nil)
	rec := httptest.NewRecorder()

	testRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestRouterActionExecuteUsesItsOwnLongerTimeout(t *testing.T) {
	registry := actions.NewRegistry()
	registry.Register(&slowAction{delay: 40 * time.Millisecond})

	handler := NewRouter(Config{
		AllowedHosts:   []string{"localhost"},
		RequestTimeout: 10 * time.Millisecond,
		ActionTimeout:  500 * time.Millisecond,
		MaxBodyBytes:   1 << 20,
		Status:         readmodel.NewStatusProvider("v9.9.9", time.Now()),
		Actions:        actions.NewRunner(registry, &stubStorer{}),
		ActionRuns:     readmodel.NewActionRunsProvider(&stubStorer{}),
	})

	req := httptest.NewRequest(http.MethodPost, "http://localhost/api/v1/actions/test.slow/execute", strings.NewReader("{}"))
	req.Header.Set("Origin", "http://localhost")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected the slow action to complete within its own timeout, got status %d body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "done") {
		t.Fatalf("expected the action result in the body, got: %s", rec.Body.String())
	}
}
