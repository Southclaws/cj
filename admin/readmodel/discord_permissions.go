package readmodel

import "github.com/bwmarrin/discordgo"

var namedPermissions = []struct {
	bit  int64
	name string
}{
	{discordgo.PermissionAdministrator, "Administrator"},
	{discordgo.PermissionCreateInstantInvite, "Create Invite"},
	{discordgo.PermissionKickMembers, "Kick Members"},
	{discordgo.PermissionBanMembers, "Ban Members"},
	{discordgo.PermissionManageChannels, "Manage Channels"},
	{discordgo.PermissionManageGuild, "Manage Server"},
	{discordgo.PermissionAddReactions, "Add Reactions"},
	{discordgo.PermissionViewAuditLogs, "View Audit Log"},
	{discordgo.PermissionViewChannel, "View Channel"},
	{discordgo.PermissionSendMessages, "Send Messages"},
	{discordgo.PermissionManageMessages, "Manage Messages"},
	{discordgo.PermissionEmbedLinks, "Embed Links"},
	{discordgo.PermissionAttachFiles, "Attach Files"},
	{discordgo.PermissionReadMessageHistory, "Read Message History"},
	{discordgo.PermissionMentionEveryone, "Mention Everyone"},
	{discordgo.PermissionUseExternalEmojis, "Use External Emojis"},
	{discordgo.PermissionChangeNickname, "Change Nickname"},
	{discordgo.PermissionManageNicknames, "Manage Nicknames"},
	{discordgo.PermissionManageRoles, "Manage Roles"},
	{discordgo.PermissionManageWebhooks, "Manage Webhooks"},
	{discordgo.PermissionManageGuildExpressions, "Manage Expressions"},
	{discordgo.PermissionManageEvents, "Manage Events"},
	{discordgo.PermissionManageThreads, "Manage Threads"},
	{discordgo.PermissionModerateMembers, "Timeout Members"},
}

func namePermissions(permissions int64) []string {
	if permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		return []string{"Administrator"}
	}

	names := make([]string, 0)
	for _, p := range namedPermissions {
		if p.bit == discordgo.PermissionAdministrator {
			continue
		}
		if permissions&p.bit == p.bit {
			names = append(names, p.name)
		}
	}
	return names
}

func rolePermissions(everyone *discordgo.Role, roles []*discordgo.Role, roleIDs []string) int64 {
	var permissions int64
	if everyone != nil {
		permissions |= everyone.Permissions
	}
	for _, role := range roles {
		for _, id := range roleIDs {
			if role.ID == id {
				permissions |= role.Permissions
			}
		}
	}
	if permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		permissions |= discordgo.PermissionAll
	}
	return permissions
}

func topRolePosition(roles []*discordgo.Role, roleIDs []string) int {
	position := 0
	for _, role := range roles {
		for _, id := range roleIDs {
			if role.ID == id && role.Position > position {
				position = role.Position
			}
		}
	}
	return position
}

func findEveryoneRole(roles []*discordgo.Role, guildID string) *discordgo.Role {
	for _, role := range roles {
		if role.ID == guildID {
			return role
		}
	}
	return nil
}
