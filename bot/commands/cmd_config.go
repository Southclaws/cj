package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandConfig(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	app := cli.NewApp()
	app.Name = "/config"
	app.HideHelp = true
	app.Commands = []cli.Command{
		{
			Name: "cooldown",
			Action: func(c *cli.Context) error {
				cmd, set, err := cm.getCommand(c)
				if err != nil {
					return err
				}
				value, err := time.ParseDuration(c.Args().Tail()[0])
				if err != nil {
					return err
				}
				cm.Discord.ChannelMessageSend(message.ChannelID, "Updated command cooldown.")
				cmd.Settings.Cooldown = value
				return set(cmd)
			},
		},
		{
			Name: "add-channel",
			Action: func(c *cli.Context) error {
				cmd, set, err := cm.getCommand(c)
				if err != nil {
					return err
				}
				cmd.Settings.Channels = append(cmd.Settings.Channels, c.Args().Tail()[0])
				cm.Discord.ChannelMessageSend(message.ChannelID, "Added channel to command.")
				return set(cmd)
			},
		},
		{
			Name: "del-channel",
			Action: func(c *cli.Context) error {
				cmd, set, err := cm.getCommand(c)
				if err != nil {
					return err
				}
				var found bool
				for i, ch := range cmd.Settings.Channels {
					if ch == c.Args().Tail()[0] {
						cmd.Settings.Channels = append(cmd.Settings.Channels[:i], cmd.Settings.Channels[i+1:]...)
						found = true
						break
					}
				}
				if found {
					cm.Discord.ChannelMessageSend(message.ChannelID, "Removed channel from command.")
				} else {
					cm.Discord.ChannelMessageSend(message.ChannelID, "Channel not found in command.")
				}
				return set(cmd)
			},
		},
		{
			Name: "add-role",
			Action: func(c *cli.Context) error {
				cmd, set, err := cm.getCommand(c)
				if err != nil {
					return err
				}
				cmd.Settings.Roles = append(cmd.Settings.Roles, c.Args().Tail()[0])
				cm.Discord.ChannelMessageSend(message.ChannelID, "Added role to command.")
				return set(cmd)
			},
		},
		{
			Name: "del-role",
			Action: func(c *cli.Context) error {
				cmd, set, err := cm.getCommand(c)
				if err != nil {
					return err
				}
				var found bool
				for i, ch := range cmd.Settings.Roles {
					if ch == c.Args().Tail()[0] {
						cmd.Settings.Roles = append(cmd.Settings.Roles[:i], cmd.Settings.Roles[i+1:]...)
						found = true
						break
					}
				}
				if found {
					cm.Discord.ChannelMessageSend(message.ChannelID, "Removed role from command.")
				} else {
					cm.Discord.ChannelMessageSend(message.ChannelID, "Role not found in command.")
				}
				return set(cmd)
			},
		},
	}

	help := func(c *cli.Context, command string) {
		var commands []string
		for _, cmd := range app.Commands {
			commands = append(commands, cmd.Name)
		}
		cm.Discord.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf("/config valid subcommands: %s", strings.Join(commands, " ")))
		return
	}
	app.CommandNotFound = help
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		help(c, "")
		return nil
	}

	splitArgs := strings.Split(args, " ")
	splitArgs = append(splitArgs[:0], append([]string{"/config"}, splitArgs[0:]...)...)
	return false, app.Run(splitArgs)
}

func (cm *CommandManager) getCommand(c *cli.Context) (cmd Command, f func(Command) error, err error) {
	if c.NArg() == 0 {
		err = errors.New("Usage: /config [command] [settings]")
		return
	}
	command := c.Args().First()
	if command == "" {
		err = errors.New("Usage: /config [command] [settings]")
		return
	}
	commandObject, ok := cm.Commands[command]
	if !ok {
		err = errors.New("Unrecognised command")
		return
	}
	return commandObject, func(newCommand Command) error {
		cm.Commands[command] = newCommand
		return cm.Storage.SetCommandSettings(command, newCommand.Settings)
	}, nil
}

// 	cooldown := c.Duration("cooldown")
// 	addChannel := c.String("add-channel")
// 	delChannel := c.String("del-channel")
// 	addRole := c.String("add-role")
// 	delRole := c.String("del-role")

// 	if cooldown.Seconds() == 0 &&
// 		addChannel == "" &&
// 		delChannel == "" &&
// 		addRole == "" &&
// 		delRole == "" {
// 		cm.Discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(`Config for %s:
// - Cooldown: %v
// - Channels: %s
// - Roles: %s
// 			`,
// 			command,
// 			commandObject.Settings.Cooldown,
// 			commandObject.Settings.Channels,
// 			commandObject.Settings.Roles,
// 		))
// 		return nil

// 	}

// }
