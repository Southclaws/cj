package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Southclaws/cj/admin/readmodel"
)

type TestServiceAction struct {
	services *readmodel.ServicesProvider
}

func NewTestServiceAction(services *readmodel.ServicesProvider) *TestServiceAction {
	return &TestServiceAction{services: services}
}

func (a *TestServiceAction) Name() string { return "services.test" }

func (a *TestServiceAction) Description() string {
	return "Check the current health of a configured external service."
}

func (a *TestServiceAction) Risk() RiskLevel { return RiskSafe }

type testServiceInput struct {
	Service string `json:"service"`
}

func decodeTestServiceInput(input json.RawMessage) (testServiceInput, error) {
	var decoded testServiceInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return testServiceInput{}, fmt.Errorf("%w: input must be a JSON object with a service field", ErrInvalidInput)
		}
	}
	if decoded.Service == "" {
		return testServiceInput{}, fmt.Errorf("%w: a service name is required", ErrInvalidInput)
	}
	return decoded, nil
}

func (a *TestServiceAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	decoded, err := decodeTestServiceInput(input)
	if err != nil {
		return Preview{}, err
	}
	return Preview{Summary: fmt.Sprintf("This will check the current health of the %q service.", decoded.Service)}, nil
}

func (a *TestServiceAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeTestServiceInput(input)
	if err != nil {
		return Result{}, err
	}

	status, err := a.services.Check(decoded.Service)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
	}

	return Result{
		Summary: fmt.Sprintf("%s: %s", status.Name, status.Status),
		Detail: map[string]any{
			"name":   status.Name,
			"status": status.Status,
			"detail": status.Detail,
		},
	}, nil
}
