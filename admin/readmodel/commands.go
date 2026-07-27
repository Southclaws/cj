package readmodel

import "github.com/Southclaws/cj/bot/commands"

type CommandsProvider struct {
	manager *commands.CommandManager
}

func NewCommandsProvider(manager *commands.CommandManager) *CommandsProvider {
	return &CommandsProvider{manager: manager}
}

func (p *CommandsProvider) List() []string {
	if p.manager == nil {
		return []string{}
	}
	names := p.manager.ListCommandNames()
	if names == nil {
		return []string{}
	}
	return names
}
