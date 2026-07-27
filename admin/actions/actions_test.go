package actions

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Southclaws/cj/storage"
)

type fakeStorer struct {
	storage.Memory
	runs        []storage.ActionRun
	auditEvents []storage.AuditEvent
}

func (f *fakeStorer) RecordActionRun(run storage.ActionRun) (string, error) {
	f.runs = append(f.runs, run)
	return "fake-id", nil
}

func (f *fakeStorer) RecordAuditEvent(event storage.AuditEvent) error {
	f.auditEvents = append(f.auditEvents, event)
	return nil
}

type fakeAction struct {
	name        string
	risk        RiskLevel
	previewErr  error
	executeErr  error
	executeHits int
	confirm     string
}

func (a *fakeAction) Name() string        { return a.name }
func (a *fakeAction) Description() string { return "a fake action for testing" }
func (a *fakeAction) Risk() RiskLevel     { return a.risk }
func (a *fakeAction) ConfirmationPhrase() string {
	return a.confirm
}

func (a *fakeAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	if a.previewErr != nil {
		return Preview{}, a.previewErr
	}
	return Preview{Summary: "preview of " + a.name}, nil
}

func (a *fakeAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	a.executeHits++
	if a.executeErr != nil {
		return Result{}, a.executeErr
	}
	return Result{Summary: "executed " + a.name}, nil
}

func TestRegistryRegisterGetList(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&fakeAction{name: "z.action", risk: RiskReadOnly})
	registry.Register(&fakeAction{name: "a.action", risk: RiskReadOnly})

	if _, ok := registry.Get("missing"); ok {
		t.Fatal("expected missing action to not be found")
	}
	if a, ok := registry.Get("z.action"); !ok || a.Name() != "z.action" {
		t.Fatalf("expected to find z.action, got %+v ok=%v", a, ok)
	}

	list := registry.List()
	if len(list) != 2 || list[0].Name() != "a.action" || list[1].Name() != "z.action" {
		t.Fatalf("expected list sorted by name, got %+v", list)
	}
}

func TestRunnerPreviewUnknownAction(t *testing.T) {
	storer := &fakeStorer{}
	runner := NewRunner(NewRegistry(), storer)

	_, err := runner.Preview(context.Background(), "does-not-exist", nil)
	if !errors.Is(err, ErrActionNotFound) {
		t.Fatalf("expected ErrActionNotFound, got %v", err)
	}
	if len(storer.runs) != 0 {
		t.Fatalf("expected no action run recorded for an unknown action, got %+v", storer.runs)
	}
}

func TestRunnerPreviewRecordsRunWithoutAudit(t *testing.T) {
	storer := &fakeStorer{}
	registry := NewRegistry()
	action := &fakeAction{name: "safe.preview", risk: RiskSafe}
	registry.Register(action)
	runner := NewRunner(registry, storer)

	preview, err := runner.Preview(context.Background(), "safe.preview", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if preview.Summary != "preview of safe.preview" {
		t.Fatalf("unexpected preview: %+v", preview)
	}

	if len(storer.runs) != 1 {
		t.Fatalf("expected exactly one action run recorded, got %+v", storer.runs)
	}
	if storer.runs[0].Mode != "preview" || storer.runs[0].Outcome != "success" {
		t.Fatalf("unexpected recorded run: %+v", storer.runs[0])
	}
	if len(storer.auditEvents) != 0 {
		t.Fatalf("expected no audit event for a preview, got %+v", storer.auditEvents)
	}
}

func TestRunnerExecuteStateChangingRecordsRunAndAudit(t *testing.T) {
	storer := &fakeStorer{}
	registry := NewRegistry()
	action := &fakeAction{name: "state.execute", risk: RiskStateChanging}
	registry.Register(action)
	runner := NewRunner(registry, storer)

	result, err := runner.Execute(context.Background(), "state.execute", nil, "req-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Summary != "executed state.execute" {
		t.Fatalf("unexpected result: %+v", result)
	}

	if len(storer.runs) != 1 || storer.runs[0].Mode != "execute" || storer.runs[0].Outcome != "success" {
		t.Fatalf("unexpected recorded runs: %+v", storer.runs)
	}
	if len(storer.auditEvents) != 1 {
		t.Fatalf("expected exactly one audit event for a state-changing action, got %+v", storer.auditEvents)
	}
	if storer.auditEvents[0].Action != "state.execute" || storer.auditEvents[0].Outcome != "success" || storer.auditEvents[0].RequestID != "req-1" {
		t.Fatalf("unexpected audit event: %+v", storer.auditEvents[0])
	}
}

func TestRunnerExecuteSafeDoesNotAudit(t *testing.T) {
	storer := &fakeStorer{}
	registry := NewRegistry()
	action := &fakeAction{name: "safe.execute", risk: RiskSafe}
	registry.Register(action)
	runner := NewRunner(registry, storer)

	if _, err := runner.Execute(context.Background(), "safe.execute", nil, "req-2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(storer.runs) != 1 {
		t.Fatalf("expected one action run recorded, got %+v", storer.runs)
	}
	if len(storer.auditEvents) != 0 {
		t.Fatalf("expected no audit event for a safe action, got %+v", storer.auditEvents)
	}
}

func TestRunnerExecuteFailureRecordsOutcome(t *testing.T) {
	storer := &fakeStorer{}
	registry := NewRegistry()
	action := &fakeAction{name: "state.fails", risk: RiskStateChanging, executeErr: errors.New("boom")}
	registry.Register(action)
	runner := NewRunner(registry, storer)

	_, err := runner.Execute(context.Background(), "state.fails", nil, "req-3")
	if err == nil {
		t.Fatal("expected the execution error to be returned")
	}

	if len(storer.runs) != 1 || storer.runs[0].Outcome != "failure" || storer.runs[0].ErrorMessage != "boom" {
		t.Fatalf("unexpected recorded run: %+v", storer.runs)
	}
	if len(storer.auditEvents) != 1 || storer.auditEvents[0].Outcome != "failure" {
		t.Fatalf("expected a failure audit event, got %+v", storer.auditEvents)
	}
}

func TestRunnerDangerousActionRequiresConfirmation(t *testing.T) {
	storer := &fakeStorer{}
	registry := NewRegistry()
	action := &fakeAction{name: "danger.action", risk: RiskDangerous, confirm: "DELETE EVERYTHING"}
	registry.Register(action)
	runner := NewRunner(registry, storer)

	_, err := runner.Execute(context.Background(), "danger.action", json.RawMessage(`{}`), "req-4")
	if !errors.Is(err, ErrConfirmationRequired) {
		t.Fatalf("expected ErrConfirmationRequired with no confirm field, got %v", err)
	}
	if action.executeHits != 0 {
		t.Fatal("expected Execute to never be called without a matching confirmation")
	}

	_, err = runner.Execute(context.Background(), "danger.action", json.RawMessage(`{"confirm":"wrong phrase"}`), "req-5")
	if !errors.Is(err, ErrConfirmationRequired) {
		t.Fatalf("expected ErrConfirmationRequired with a wrong confirm field, got %v", err)
	}
	if action.executeHits != 0 {
		t.Fatal("expected Execute to never be called with a mismatched confirmation")
	}

	result, err := runner.Execute(context.Background(), "danger.action", json.RawMessage(`{"confirm":"DELETE EVERYTHING"}`), "req-6")
	if err != nil {
		t.Fatalf("unexpected error with a matching confirmation: %v", err)
	}
	if result.Summary != "executed danger.action" || action.executeHits != 1 {
		t.Fatalf("expected the action to execute exactly once, got hits=%d result=%+v", action.executeHits, result)
	}
}
