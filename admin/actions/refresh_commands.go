package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Southclaws/cj/bot/commands"
)

type RefreshCommandsAction struct {
	manager *commands.CommandManager
}

func NewRefreshCommandsAction(manager *commands.CommandManager) *RefreshCommandsAction {
	return &RefreshCommandsAction{manager: manager}
}

func (a *RefreshCommandsAction) Name() string { return "discord.refresh-commands" }

func (a *RefreshCommandsAction) Description() string {
	return "Delete and re-register all Discord slash commands for this bot."
}

func (a *RefreshCommandsAction) Risk() RiskLevel { return RiskStateChanging }

func (a *RefreshCommandsAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	return Preview{
		Summary: "This will delete every registered global and per-guild slash command, then re-register the currently enabled command set. Discord can take up to an hour to propagate the change.",
	}, nil
}

func (a *RefreshCommandsAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	result, err := a.manager.RefreshCommands()
	if err != nil {
		return Result{}, err
	}
	return Result{
		Summary: fmt.Sprintf("Deleted %d existing commands and registered %d commands.", result.Deleted, result.Added),
		Detail: map[string]any{
			"deleted": result.Deleted,
			"added":   result.Added,
		},
	}, nil
}
