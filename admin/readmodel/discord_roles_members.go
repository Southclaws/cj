package readmodel

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	defaultMembersPageSize = 50
	maxMembersPageSize     = 200
)

func (p *DiscordProvider) botMember() (*discordgo.Member, error) {
	if !p.ready() {
		return nil, ErrDiscordNotReady
	}
	return p.session.S.GuildMember(p.guildID(), p.session.S.State.User.ID)
}

func (p *DiscordProvider) Roles() ([]Role, error) {
	if !p.ready() {
		return nil, ErrDiscordNotReady
	}

	roles, err := p.session.S.GuildRoles(p.guildID())
	if err != nil {
		return nil, err
	}

	bot, err := p.botMember()
	if err != nil {
		return nil, err
	}
	botTop := topRolePosition(roles, bot.Roles)
	botCanManageRoles := rolePermissions(findEveryoneRole(roles, p.guildID()), roles, bot.Roles)&
		(discordgo.PermissionAdministrator|discordgo.PermissionManageRoles) != 0

	members, cachedAt := p.session.Members()
	if cachedAt.IsZero() {
		return nil, ErrMemberCacheNotReady
	}
	memberCounts := countMembersByRole(members)

	out := make([]Role, 0, len(roles))
	for _, role := range roles {
		var reasons []string
		canManage := true

		if role.ID == p.guildID() {
			canManage = false
			reasons = append(reasons, "this is the @everyone role")
		}
		if role.Managed {
			canManage = false
			reasons = append(reasons, "role is managed by an integration or bot")
		}
		if !botCanManageRoles {
			canManage = false
			reasons = append(reasons, "CJ does not have the Manage Roles permission")
		}
		if role.Position >= botTop && role.ID != p.guildID() {
			canManage = false
			reasons = append(reasons, "role is at or above CJ's highest role")
		}

		out = append(out, Role{
			ID:          role.ID,
			Name:        role.Name,
			Color:       role.Color,
			Position:    role.Position,
			Managed:     role.Managed,
			MemberCount: memberCounts[role.ID],
			Permissions: namePermissions(role.Permissions),
			CJCanManage: canManage,
			Reasons:     reasons,
		})
	}

	return out, nil
}

type MembersPage struct {
	Members  []Member `json:"members"`
	Total    int      `json:"total"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	CachedAt string   `json:"cachedAt"`
}

func (p *DiscordProvider) Members(query string, limit, offset int) (MembersPage, error) {
	if !p.ready() {
		return MembersPage{}, ErrDiscordNotReady
	}
	if limit <= 0 || limit > maxMembersPageSize {
		limit = defaultMembersPageSize
	}
	if offset < 0 {
		offset = 0
	}

	all, cachedAt := p.session.Members()
	if cachedAt.IsZero() {
		return MembersPage{}, ErrMemberCacheNotReady
	}

	filtered := filterMembersByQuery(all, query)

	roles, err := p.session.S.GuildRoles(p.guildID())
	if err != nil {
		return MembersPage{}, err
	}

	guild, err := p.session.S.Guild(p.guildID())
	if err != nil {
		return MembersPage{}, err
	}

	bot, err := p.botMember()
	if err != nil {
		return MembersPage{}, err
	}
	botTop := topRolePosition(roles, bot.Roles)

	page := paginateMemberSlice(filtered, limit, offset)
	out := make([]Member, 0, len(page))
	for _, m := range page {
		out = append(out, buildMember(m, roles, guild, botTop))
	}

	return MembersPage{
		Members:  out,
		Total:    len(filtered),
		Limit:    limit,
		Offset:   offset,
		CachedAt: cachedAt.UTC().Format(time.RFC3339),
	}, nil
}

func filterMembersByQuery(members []*discordgo.Member, query string) []*discordgo.Member {
	if query == "" {
		return members
	}
	q := strings.ToLower(query)
	out := make([]*discordgo.Member, 0, len(members))
	for _, m := range members {
		if memberMatchesQuery(m, q) {
			out = append(out, m)
		}
	}
	return out
}

func memberMatchesQuery(m *discordgo.Member, lowerQuery string) bool {
	if strings.Contains(strings.ToLower(m.Nick), lowerQuery) {
		return true
	}
	if m.User != nil {
		if strings.Contains(strings.ToLower(m.User.Username), lowerQuery) {
			return true
		}
		if strings.Contains(m.User.ID, lowerQuery) {
			return true
		}
	}
	return false
}

func paginateMemberSlice(members []*discordgo.Member, limit, offset int) []*discordgo.Member {
	if offset >= len(members) {
		return []*discordgo.Member{}
	}
	end := offset + limit
	if end > len(members) {
		end = len(members)
	}
	return members[offset:end]
}

func countMembersByRole(members []*discordgo.Member) map[string]int {
	counts := make(map[string]int)
	for _, m := range members {
		for _, roleID := range m.Roles {
			counts[roleID]++
		}
	}
	return counts
}

func (p *DiscordProvider) MatchMembers(query string, max int) []PersonRef {
	if p == nil || query == "" || !p.ready() {
		return nil
	}

	all, cachedAt := p.session.Members()
	if cachedAt.IsZero() {
		return nil
	}

	matches := filterMembersByQuery(all, query)
	if len(matches) > max {
		matches = matches[:max]
	}

	out := make([]PersonRef, len(matches))
	for i, m := range matches {
		ref := PersonRef{Known: true, Nickname: m.Nick, AvatarURL: m.AvatarURL(avatarSize)}
		if m.User != nil {
			ref.ID = m.User.ID
			ref.Username = m.User.Username
		}
		out[i] = ref
	}
	return out
}

func (p *DiscordProvider) PersonRef(id string) PersonRef {
	ref := PersonRef{ID: id}
	if p == nil || id == "" || !p.ready() {
		return ref
	}

	m, ok := p.session.MemberByID(id)
	if !ok {
		return ref
	}

	ref.Known = true
	ref.Nickname = m.Nick
	if m.User != nil {
		ref.Username = m.User.Username
	}
	ref.AvatarURL = m.AvatarURL(avatarSize)
	return ref
}

func (p *DiscordProvider) Member(id string) (Member, error) {
	if !p.ready() {
		return Member{}, ErrDiscordNotReady
	}

	m, err := p.session.S.GuildMember(p.guildID(), id)
	if err != nil {
		return Member{}, err
	}

	roles, err := p.session.S.GuildRoles(p.guildID())
	if err != nil {
		return Member{}, err
	}

	guild, err := p.session.S.Guild(p.guildID())
	if err != nil {
		return Member{}, err
	}

	bot, err := p.botMember()
	if err != nil {
		return Member{}, err
	}
	botTop := topRolePosition(roles, bot.Roles)

	return buildMember(m, roles, guild, botTop), nil
}

func buildMember(m *discordgo.Member, roles []*discordgo.Role, guild *discordgo.Guild, botTop int) Member {
	var reasons []string
	canActOn := true

	if m.User != nil && m.User.ID == guild.OwnerID {
		canActOn = false
		reasons = append(reasons, "member is the guild owner")
	}
	memberTop := topRolePosition(roles, m.Roles)
	if memberTop >= botTop {
		canActOn = false
		reasons = append(reasons, "member's highest role is at or above CJ's highest role")
	}

	member := Member{
		Roles:      m.Roles,
		CJCanActOn: canActOn,
		Reasons:    reasons,
	}
	if m.User != nil {
		member.ID = m.User.ID
		member.Username = m.User.Username
	}
	member.Nickname = m.Nick
	member.AvatarURL = m.AvatarURL(avatarSize)
	if !m.JoinedAt.IsZero() {
		member.JoinedAt = m.JoinedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return member
}

func (p *DiscordProvider) MemberPermissions(id string) (MemberPermissions, error) {
	if !p.ready() {
		return MemberPermissions{}, ErrDiscordNotReady
	}

	m, err := p.session.S.GuildMember(p.guildID(), id)
	if err != nil {
		return MemberPermissions{}, err
	}

	roles, err := p.session.S.GuildRoles(p.guildID())
	if err != nil {
		return MemberPermissions{}, err
	}

	guild, err := p.session.S.Guild(p.guildID())
	if err != nil {
		return MemberPermissions{}, err
	}

	var permissions int64
	if id == guild.OwnerID {
		permissions = discordgo.PermissionAll
	} else {
		permissions = rolePermissions(findEveryoneRole(roles, p.guildID()), roles, m.Roles)
	}

	return MemberPermissions{
		MemberID:        id,
		Permissions:     namePermissions(permissions),
		IsAdministrator: permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator,
	}, nil
}

func (p *DiscordProvider) AuditLog(limit int) ([]AuditLogEntry, error) {
	if !p.ready() {
		return nil, ErrDiscordNotReady
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	log, err := p.session.S.GuildAuditLog(p.guildID(), "", "", 0, limit)
	if err != nil {
		return nil, err
	}

	out := make([]AuditLogEntry, 0, len(log.AuditLogEntries))
	for _, entry := range log.AuditLogEntries {
		actionType := 0
		if entry.ActionType != nil {
			actionType = int(*entry.ActionType)
		}
		out = append(out, AuditLogEntry{
			ID:         entry.ID,
			ActionType: actionType,
			ActionName: auditActionName(actionType),
			TargetID:   entry.TargetID,
			UserID:     entry.UserID,
			User:       p.PersonRef(entry.UserID),
			Reason:     entry.Reason,
		})
	}
	return out, nil
}

func auditActionName(actionType int) string {
	if name, ok := auditActionNames[actionType]; ok {
		return name
	}
	return fmt.Sprintf("action_%d", actionType)
}

var auditActionNames = map[int]string{
	1:  "guild_update",
	10: "channel_create",
	11: "channel_update",
	12: "channel_delete",
	13: "channel_overwrite_create",
	14: "channel_overwrite_update",
	15: "channel_overwrite_delete",
	20: "member_kick",
	21: "member_prune",
	22: "member_ban_add",
	23: "member_ban_remove",
	24: "member_update",
	25: "member_role_update",
	26: "member_move",
	27: "member_disconnect",
	28: "bot_add",
	30: "role_create",
	31: "role_update",
	32: "role_delete",
	40: "invite_create",
	41: "invite_update",
	42: "invite_delete",
	50: "webhook_create",
	51: "webhook_update",
	52: "webhook_delete",
	60: "emoji_create",
	61: "emoji_update",
	62: "emoji_delete",
	72: "message_delete",
	73: "message_bulk_delete",
	74: "message_pin",
	75: "message_unpin",
}
