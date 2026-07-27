package actions

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/Southclaws/cj/storage"
)

var (
	ErrActionNotFound       = errors.New("no such action is registered")
	ErrConfirmationRequired = errors.New("this action requires a matching confirmation phrase")
	ErrInvalidInput         = errors.New("invalid action input")
)

type confirmationEnvelope struct {
	Confirm string `json:"confirm"`
}

type Runner struct {
	registry *Registry
	storer   storage.Storer
}

func NewRunner(registry *Registry, storer storage.Storer) *Runner {
	return &Runner{registry: registry, storer: storer}
}

func (r *Runner) Registry() *Registry {
	return r.registry
}

func (r *Runner) Preview(ctx context.Context, name string, input json.RawMessage) (Preview, error) {
	action, ok := r.registry.Get(name)
	if !ok {
		return Preview{}, ErrActionNotFound
	}

	start := time.Now()
	preview, err := action.Preview(ctx, input)
	duration := time.Since(start)

	r.recordRun(action, "preview", preview.Summary, duration, err, "")

	return preview, err
}

func (r *Runner) Execute(ctx context.Context, name string, input json.RawMessage, requestID string) (Result, error) {
	action, ok := r.registry.Get(name)
	if !ok {
		return Result{}, ErrActionNotFound
	}

	if action.Risk() == RiskDangerous {
		expected := action.Name()
		if cp, ok := action.(ConfirmationPhrase); ok {
			expected = cp.ConfirmationPhrase()
		}
		var envelope confirmationEnvelope
		_ = json.Unmarshal(input, &envelope)
		if envelope.Confirm != expected {
			return Result{}, ErrConfirmationRequired
		}
	}

	start := time.Now()
	result, err := action.Execute(ctx, input)
	duration := time.Since(start)

	r.recordRun(action, "execute", result.Summary, duration, err, requestID)

	if action.Risk() == RiskStateChanging || action.Risk() == RiskDangerous {
		outcome := "success"
		if err != nil {
			outcome = "failure"
		}
		auditErr := r.storer.RecordAuditEvent(storage.AuditEvent{
			Action:     action.Name(),
			TargetType: "action",
			Outcome:    outcome,
			RequestID:  requestID,
		})
		if auditErr != nil {
			zap.L().Error("failed to record audit event for action", zap.Error(auditErr), zap.String("action", action.Name()))
		}
	}

	return result, err
}

func (r *Runner) recordRun(action Action, mode, summary string, duration time.Duration, execErr error, requestID string) {
	outcome := "success"
	errorMessage := ""
	if execErr != nil {
		outcome = "failure"
		errorMessage = execErr.Error()
	}

	_, err := r.storer.RecordActionRun(storage.ActionRun{
		Action:       action.Name(),
		Risk:         string(action.Risk()),
		Mode:         mode,
		Outcome:      outcome,
		Summary:      summary,
		ErrorMessage: errorMessage,
		DurationMS:   duration.Milliseconds(),
		RequestID:    requestID,
	})
	if err != nil {
		zap.L().Error("failed to record action run", zap.Error(err), zap.String("action", action.Name()))
	}
}
