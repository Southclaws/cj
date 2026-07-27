package commands

import "github.com/Southclaws/cj/types"

func (cm *CommandManager) ListCommandNames() []string {
	names := make([]string, len(cm.Commands))
	for i, command := range cm.Commands {
		names[i] = command.Name
	}
	return names
}

func (cm *CommandManager) GetCommandConfig(commandName string) (types.CommandSettings, error) {
	cmd, _, err := cm.getCommand(commandName)
	if err != nil {
		return types.CommandSettings{}, err
	}
	return cmd.Settings, nil
}

func (cm *CommandManager) SetCommandConfig(commandName string, settings types.CommandSettings) (types.CommandSettings, error) {
	cmd, set, err := cm.getCommand(commandName)
	if err != nil {
		return types.CommandSettings{}, err
	}
	cmd.Settings = settings
	if err := set(cmd); err != nil {
		return types.CommandSettings{}, err
	}
	return cmd.Settings, nil
}
