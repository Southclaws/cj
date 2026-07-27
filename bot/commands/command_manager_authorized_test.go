package commands

import (
	"testing"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

func cmSettingsWithRoles(roles ...string) types.CommandSettings {
	return types.CommandSettings{Roles: roles}
}

func TestAuthorizedAdministrativeFallsBackToDiscordAdminWhenUnconfigured(t *testing.T) {
	cm := &CommandManager{}
	command := Command{Name: "/config", IsAdministrative: true}

	admin := &discordgo.Member{Permissions: discordgo.PermissionAdministrator}
	if !cm.authorized(command, admin) {
		t.Fatal("expected a real Discord administrator to be authorized when no roles are configured")
	}

	normal := &discordgo.Member{Permissions: discordgo.PermissionSendMessages}
	if cm.authorized(command, normal) {
		t.Fatal("expected a normal member to be denied when no roles are configured")
	}
}

func TestAuthorizedAdministrativeUsesConfiguredRolesWhenSet(t *testing.T) {
	cm := &CommandManager{}
	command := Command{
		Name:             "/config",
		IsAdministrative: true,
		Settings:         cmSettingsWithRoles("111"),
	}

	member := &discordgo.Member{Roles: []string{"111"}}
	if !cm.authorized(command, member) {
		t.Fatal("expected member with the configured role to be authorized")
	}

	other := &discordgo.Member{Roles: []string{"222"}, Permissions: discordgo.PermissionAdministrator}
	if cm.authorized(command, other) {
		t.Fatal("expected a real Discord administrator without the configured role to be denied once roles are explicitly configured")
	}
}

func TestAuthorizedPublicCommandIsOpenByDefault(t *testing.T) {
	cm := &CommandManager{}
	command := Command{Name: "/wiki", Settings: cmSettingsWithRoles("all")}

	member := &discordgo.Member{}
	if !cm.authorized(command, member) {
		t.Fatal("expected a public command with roles=[all] to be open to everyone")
	}
}

func TestCommandEnabledSearchMessageRequiresChannelAndRoles(t *testing.T) {
	unconfigured := storage.GuildSettings{}
	command := Command{Name: "/searchmessage"}
	if commandEnabled(command, unconfigured) {
		t.Fatal("expected /searchmessage to be disabled with no guild settings configured")
	}

	channelOnly := storage.GuildSettings{SearchMessageChannelID: "1"}
	if commandEnabled(command, channelOnly) {
		t.Fatal("expected /searchmessage to stay disabled without configured roles even if the channel is set")
	}

	command.Settings = cmSettingsWithRoles("111")
	if !commandEnabled(command, channelOnly) {
		t.Fatal("expected /searchmessage to be enabled once both channel and roles are configured")
	}
}

func TestCommandEnabledLTFRequiresChannelAndUsers(t *testing.T) {
	command := Command{Name: "/ltf"}

	if commandEnabled(command, storage.GuildSettings{}) {
		t.Fatal("expected /ltf to be disabled with no guild settings configured")
	}
	if commandEnabled(command, storage.GuildSettings{LTFChannelID: "1"}) {
		t.Fatal("expected /ltf to stay disabled without configured users")
	}
	if !commandEnabled(command, storage.GuildSettings{LTFChannelID: "1", LTFUserIDs: []string{"2"}}) {
		t.Fatal("expected /ltf to be enabled once both channel and users are configured")
	}
}

func TestCommandEnabledDefaultsToTrueForOrdinaryCommands(t *testing.T) {
	if !commandEnabled(Command{Name: "/wiki"}, storage.GuildSettings{}) {
		t.Fatal("expected an ordinary command to be enabled regardless of guild settings")
	}
}
