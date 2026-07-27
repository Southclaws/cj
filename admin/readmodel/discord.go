package readmodel

import (
	"errors"
	"sort"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/discord"
)

var ErrDiscordNotReady = errors.New("discord session is not ready")
var ErrMemberCacheNotReady = errors.New("the member cache has not been populated yet")

type Guild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OwnerID     string `json:"ownerId"`
	MemberCount int    `json:"memberCount"`
	Icon        string `json:"icon,omitempty"`
}

type Channel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Position int    `json:"position"`
}

type ChannelCategory struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Position int       `json:"position"`
	Channels []Channel `json:"channels"`
}

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Color       int      `json:"color"`
	Position    int      `json:"position"`
	Managed     bool     `json:"managed"`
	MemberCount int      `json:"memberCount"`
	Permissions []string `json:"permissions"`
	CJCanManage bool     `json:"cjCanManage"`
	Reasons     []string `json:"reasons,omitempty"`
}

type Member struct {
	ID         string   `json:"id"`
	Username   string   `json:"username"`
	Nickname   string   `json:"nickname,omitempty"`
	AvatarURL  string   `json:"avatarUrl,omitempty"`
	Roles      []string `json:"roles"`
	JoinedAt   string   `json:"joinedAt,omitempty"`
	CJCanActOn bool     `json:"cjCanActOn"`
	Reasons    []string `json:"reasons,omitempty"`
}

type PersonRef struct {
	ID        string `json:"id"`
	Username  string `json:"username,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	Known     bool   `json:"known"`
}

const avatarSize = "64"

type MemberPermissions struct {
	MemberID        string   `json:"memberId"`
	Permissions     []string `json:"permissions"`
	IsAdministrator bool     `json:"isAdministrator"`
}

type AuditLogEntry struct {
	ID         string    `json:"id"`
	ActionType int       `json:"actionType"`
	ActionName string    `json:"actionName"`
	TargetID   string    `json:"targetId,omitempty"`
	UserID     string    `json:"userId,omitempty"`
	User       PersonRef `json:"user"`
	Reason     string    `json:"reason,omitempty"`
}

type DiscordProvider struct {
	session *discord.Session
}

func NewDiscordProvider(session *discord.Session) *DiscordProvider {
	return &DiscordProvider{session: session}
}

func (p *DiscordProvider) ready() bool {
	return p.session != nil && p.session.S != nil && p.session.S.State != nil && p.session.S.State.User != nil
}

func (p *DiscordProvider) guildID() string {
	return p.session.Config.GuildID
}

func (p *DiscordProvider) Guild() (Guild, error) {
	if !p.ready() {
		return Guild{}, ErrDiscordNotReady
	}

	g, err := p.session.S.GuildWithCounts(p.guildID())
	if err != nil {
		return Guild{}, err
	}

	return Guild{
		ID:          g.ID,
		Name:        g.Name,
		OwnerID:     g.OwnerID,
		MemberCount: g.ApproximateMemberCount,
		Icon:        g.Icon,
	}, nil
}

func (p *DiscordProvider) Channels() ([]ChannelCategory, error) {
	if !p.ready() {
		return nil, ErrDiscordNotReady
	}

	channels, err := p.session.S.GuildChannels(p.guildID())
	if err != nil {
		return nil, err
	}

	return buildChannelCategories(channels), nil
}

func buildChannelCategories(channels []*discordgo.Channel) []ChannelCategory {
	byID := make(map[string]*ChannelCategory)
	uncategorised := &ChannelCategory{ID: "", Name: "No category"}

	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory {
			byID[ch.ID] = &ChannelCategory{ID: ch.ID, Name: ch.Name, Position: ch.Position}
		}
	}

	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory {
			continue
		}
		entry := Channel{ID: ch.ID, Name: ch.Name, Type: channelTypeName(ch.Type), Position: ch.Position}
		if cat, ok := byID[ch.ParentID]; ok {
			cat.Channels = append(cat.Channels, entry)
		} else {
			uncategorised.Channels = append(uncategorised.Channels, entry)
		}
	}

	out := make([]ChannelCategory, 0, len(byID)+1)
	for _, cat := range byID {
		sort.Slice(cat.Channels, func(i, j int) bool { return cat.Channels[i].Position < cat.Channels[j].Position })
		out = append(out, *cat)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Position < out[j].Position })

	if len(uncategorised.Channels) > 0 {
		sort.Slice(uncategorised.Channels, func(i, j int) bool {
			return uncategorised.Channels[i].Position < uncategorised.Channels[j].Position
		})
		out = append(out, *uncategorised)
	}

	return out
}

func channelTypeName(t discordgo.ChannelType) string {
	switch t {
	case discordgo.ChannelTypeGuildText:
		return "text"
	case discordgo.ChannelTypeGuildVoice:
		return "voice"
	case discordgo.ChannelTypeGuildNews:
		return "announcement"
	case discordgo.ChannelTypeGuildStageVoice:
		return "stage"
	case discordgo.ChannelTypeGuildForum:
		return "forum"
	default:
		return "other"
	}
}
