package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandConfig(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	commandName := args["command"].StringValue()

	cmd, set, err := cm.getCommand(commandName)
	if err != nil {
		cm.replyDirectly(interaction, err.Error())
		return
	}

	newConfigValue, hasNewConfig := args["config"]
	if !hasNewConfig {
		var b []byte
		b, err = json.Marshal(cmd.Settings)
		if err != nil {
			cm.replyDirectly(interaction,
				fmt.Sprintf("Existing config for **%s** couldn't be marshalled to valid JSON. Contact a CJ administrator.", commandName))
		} else {
			cm.replyDirectly(interaction, fmt.Sprintf("Config for **%s**:\n```json\n%s```", commandName, string(b)))
		}
	} else {
		newConfig := newConfigValue.StringValue()
		newConfig = strings.TrimLeft(
			strings.TrimRight(newConfig, "` "),
			"` ")
		if err = json.Unmarshal([]byte(newConfig), &cmd.Settings); err != nil {
			cm.replyDirectly(interaction, fmt.Sprintf("Failed to unmarshal to JSON:\n```json\n%s```\nWith error: %s", newConfig, err.Error()))
		}
		if err = set(cmd); err != nil {
			cm.replyDirectly(interaction, fmt.Sprintf("Failed to update command **%s** with error: %s", commandName, err.Error()))
		}
		var b []byte
		b, err = json.Marshal(cmd.Settings)
		if err != nil {
			cm.replyDirectly(interaction, fmt.Sprintf("Successfully updated command **%s** but re-marshalling JSON failed with error: %s", commandName, err.Error()))
		}
		cm.replyDirectly(interaction, fmt.Sprintf("Updated **%s** to:\n```json\n%s```", commandName, string(b)))
	}
	return
}

func (cm *CommandManager) getCommand(commandName string) (cmd Command, f func(Command) error, err error) {
	var index = -1
	for i, v := range cm.Commands {
		if v.Name == commandName {
			index = i
			break
		}
	}
	if index == -1 {
		err = errors.Errorf("Unrecognised command name `%s`", commandName)
		return
	}
	return cm.Commands[index], func(newCommand Command) error {
		cm.Commands[index] = newCommand
		return cm.Storage.SetCommandSettings(commandName, newCommand.Settings)
	}, nil
}
