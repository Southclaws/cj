package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Southclaws/cj/bot/commands"
	"github.com/Southclaws/cj/types"
)

type commandNameInput struct {
	Command string `json:"command"`
}

func decodeCommandNameInput(input json.RawMessage) (commandNameInput, error) {
	var decoded commandNameInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return commandNameInput{}, fmt.Errorf("%w: input must be a JSON object with a command field", ErrInvalidInput)
		}
	}
	if decoded.Command == "" {
		return commandNameInput{}, fmt.Errorf("%w: a command name is required", ErrInvalidInput)
	}
	return decoded, nil
}

func commandSettingsDetail(settings types.CommandSettings) map[string]any {
	return map[string]any{
		"cooldown": settings.Cooldown.String(),
		"roles":    settings.Roles,
		"misc":     settings.Misc,
	}
}

type GetCommandSettingsAction struct {
	manager *commands.CommandManager
}

func NewGetCommandSettingsAction(manager *commands.CommandManager) *GetCommandSettingsAction {
	return &GetCommandSettingsAction{manager: manager}
}

func (a *GetCommandSettingsAction) Name() string { return "commands.get-settings" }

func (a *GetCommandSettingsAction) Description() string {
	return "View the stored settings for a registered CJ command, the same as /config with no new value."
}

func (a *GetCommandSettingsAction) Risk() RiskLevel { return RiskReadOnly }

func (a *GetCommandSettingsAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	decoded, err := decodeCommandNameInput(input)
	if err != nil {
		return Preview{}, err
	}
	return Preview{Summary: fmt.Sprintf("This will read the stored settings for %q.", decoded.Command)}, nil
}

func (a *GetCommandSettingsAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeCommandNameInput(input)
	if err != nil {
		return Result{}, err
	}

	settings, err := a.manager.GetCommandConfig(decoded.Command)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
	}

	return Result{
		Summary: fmt.Sprintf("Settings for %q loaded.", decoded.Command),
		Detail:  commandSettingsDetail(settings),
	}, nil
}

type SetCommandSettingsAction struct {
	manager *commands.CommandManager
}

func NewSetCommandSettingsAction(manager *commands.CommandManager) *SetCommandSettingsAction {
	return &SetCommandSettingsAction{manager: manager}
}

func (a *SetCommandSettingsAction) Name() string { return "commands.set-settings" }

func (a *SetCommandSettingsAction) Description() string {
	return "Update the stored settings for a registered CJ command, the same as /config with a new value."
}

func (a *SetCommandSettingsAction) Risk() RiskLevel { return RiskStateChanging }

type setCommandSettingsInput struct {
	Command  string                `json:"command"`
	Settings types.CommandSettings `json:"settings"`
}

func decodeSetCommandSettingsInput(input json.RawMessage) (setCommandSettingsInput, error) {
	var decoded setCommandSettingsInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return setCommandSettingsInput{}, fmt.Errorf("%w: input must be a JSON object with command and settings fields", ErrInvalidInput)
		}
	}
	if decoded.Command == "" {
		return setCommandSettingsInput{}, fmt.Errorf("%w: a command name is required", ErrInvalidInput)
	}
	return decoded, nil
}

func (a *SetCommandSettingsAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	decoded, err := decodeSetCommandSettingsInput(input)
	if err != nil {
		return Preview{}, err
	}
	return Preview{
		Summary: fmt.Sprintf("This will overwrite the stored settings for %q.", decoded.Command),
		Detail:  commandSettingsDetail(decoded.Settings),
	}, nil
}

func (a *SetCommandSettingsAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeSetCommandSettingsInput(input)
	if err != nil {
		return Result{}, err
	}

	updated, err := a.manager.SetCommandConfig(decoded.Command, decoded.Settings)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
	}

	return Result{
		Summary: fmt.Sprintf("Settings for %q updated.", decoded.Command),
		Detail:  commandSettingsDetail(updated),
	}, nil
}
