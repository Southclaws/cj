package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Southclaws/cj/bot/heartbeat"
)

type RunJobAction struct {
	heartbeat *heartbeat.Heartbeat
}

func NewRunJobAction(hb *heartbeat.Heartbeat) *RunJobAction {
	return &RunJobAction{heartbeat: hb}
}

func (a *RunJobAction) Name() string { return "jobs.run" }

func (a *RunJobAction) Description() string {
	return "Trigger a registered scheduled job immediately, outside its normal schedule."
}

func (a *RunJobAction) Risk() RiskLevel { return RiskStateChanging }

type runJobInput struct {
	Name string `json:"name"`
}

func decodeRunJobInput(input json.RawMessage) (runJobInput, error) {
	var decoded runJobInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return runJobInput{}, fmt.Errorf("%w: input must be a JSON object with a name field", ErrInvalidInput)
		}
	}
	if decoded.Name == "" {
		return runJobInput{}, fmt.Errorf("%w: a job name is required", ErrInvalidInput)
	}
	return decoded, nil
}

func (a *RunJobAction) findJob(name string) (heartbeat.JobStatus, bool) {
	if a.heartbeat == nil {
		return heartbeat.JobStatus{}, false
	}
	for _, status := range a.heartbeat.Statuses() {
		if status.Name == name {
			return status, true
		}
	}
	return heartbeat.JobStatus{}, false
}

func (a *RunJobAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	decoded, err := decodeRunJobInput(input)
	if err != nil {
		return Preview{}, err
	}
	job, found := a.findJob(decoded.Name)
	if !found {
		return Preview{}, fmt.Errorf("%w: no such job is registered: %s", ErrInvalidInput, decoded.Name)
	}
	return Preview{
		Summary: fmt.Sprintf("This will run the %q job immediately, outside its %s schedule.", job.Name, job.Schedule),
		Detail: map[string]any{
			"provider": job.Provider,
			"schedule": job.Schedule,
			"runCount": job.RunCount,
		},
	}, nil
}

func (a *RunJobAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeRunJobInput(input)
	if err != nil {
		return Result{}, err
	}
	if _, found := a.findJob(decoded.Name); !found {
		return Result{}, fmt.Errorf("%w: no such job is registered: %s", ErrInvalidInput, decoded.Name)
	}

	if err := a.heartbeat.RunJob(decoded.Name); err != nil {
		return Result{}, err
	}

	return Result{Summary: fmt.Sprintf("The %q job ran successfully.", decoded.Name)}, nil
}
