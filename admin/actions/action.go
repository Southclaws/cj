package actions

import (
	"context"
	"encoding/json"
)

type RiskLevel string

const (
	RiskReadOnly      RiskLevel = "read-only"
	RiskSafe          RiskLevel = "safe"
	RiskStateChanging RiskLevel = "state-changing"
	RiskDangerous     RiskLevel = "dangerous"
)

type Preview struct {
	Summary string         `json:"summary"`
	Detail  map[string]any `json:"detail,omitempty"`
}

type Result struct {
	Summary string         `json:"summary"`
	Detail  map[string]any `json:"detail,omitempty"`
}

type Action interface {
	Name() string
	Description() string
	Risk() RiskLevel
	Preview(ctx context.Context, input json.RawMessage) (Preview, error)
	Execute(ctx context.Context, input json.RawMessage) (Result, error)
}

type ConfirmationPhrase interface {
	ConfirmationPhrase() string
}
