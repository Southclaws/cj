package actions

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Southclaws/cj/bot/heartbeat"
)

type RefreshReadmeAction struct {
	heartbeat *heartbeat.Heartbeat
}

func NewRefreshReadmeAction(hb *heartbeat.Heartbeat) *RefreshReadmeAction {
	return &RefreshReadmeAction{heartbeat: hb}
}

func (a *RefreshReadmeAction) Name() string { return "readme.refresh" }

func (a *RefreshReadmeAction) Description() string {
	return "Fetch the upstream README and update the managed Discord message."
}

func (a *RefreshReadmeAction) Risk() RiskLevel { return RiskStateChanging }

func (a *RefreshReadmeAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	return Preview{
		Summary: "This will fetch the configured upstream README file and edit (or send) the managed Discord message with its content.",
	}, nil
}

func (a *RefreshReadmeAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	if a.heartbeat == nil || a.heartbeat.Readme == nil {
		return Result{}, errors.New("readme sync is not available")
	}
	if err := a.heartbeat.Readme.Refresh(); err != nil {
		return Result{}, err
	}
	return Result{Summary: "The managed README message was refreshed from the upstream source."}, nil
}
