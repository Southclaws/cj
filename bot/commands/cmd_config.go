package commands

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Southclaws/cj/types"
)

var matchConfigInput = regexp.MustCompile(`(?ms)\x60\x60\x60.*\n(.+)\n\x60\x60\x60`)

func (cm *CommandManager) commandConfig(
	args string,
	message discordgo.Message,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	t := strings.SplitN(args, "\n", 2)
	if len(t[0]) == 0 {
		cm.Discord.ChannelMessageSend(message.ChannelID, `Usage:
/config [command] to view the current configuration
/config [command]⏎
[raw JSON as a code block]
to update the configuration
`)
		return
	}
	command := t[0]

	cmd, set, err := cm.getCommand(command)
	if err != nil {
		cm.Discord.ChannelMessageSend(message.ChannelID, err.Error())
		return false, nil
	}

	if len(t) == 1 {
		var b []byte
		b, err = json.Marshal(cmd.Settings)
		if err != nil {
			return false, err
		}
		cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("```json\n%s```", string(b)))
	} else {
		match := matchConfigInput.FindStringSubmatch(message.Content)
		if err = json.Unmarshal([]byte(match[1]), &cmd.Settings); err != nil {
			return false, err
		}
		if err = set(cmd); err != nil {
			return false, err
		}
		var b []byte
		b, err = json.Marshal(cmd.Settings)
		if err != nil {
			return false, err
		}
		cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Updated to:\n```json\n%s```", string(b)))
	}

	return false, nil
}

func (cm *CommandManager) getCommand(command string) (cmd Command, f func(Command) error, err error) {
	commandObject, ok := cm.Commands[command]
	if !ok {
		err = errors.Errorf("Unrecognised command `%s`", command)
		return
	}
	return commandObject, func(newCommand Command) error {
		cm.Commands[command] = newCommand
		return cm.Storage.SetCommandSettings(command, newCommand.Settings)
	}, nil
}
