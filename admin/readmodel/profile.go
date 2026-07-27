package readmodel

import "github.com/Southclaws/cj/storage"

const profileRecentMessageLimit = 10

type ProfileMessage struct {
	Timestamp   int64  `json:"timestamp"`
	ChannelID   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	Message     string `json:"message"`
	MessageID   string `json:"messageId"`
}

type Profile struct {
	DiscordUserID  string             `json:"discordUserId"`
	Member         *Member            `json:"member,omitempty"`
	Permissions    *MemberPermissions `json:"permissions,omitempty"`
	User           *storage.User      `json:"user,omitempty"`
	RecentMessages []ProfileMessage   `json:"recentMessages"`
	MessageCount   int                `json:"messageCount"`
	Rank           int                `json:"rank"`
}

type ProfileProvider struct {
	storer       storage.Storer
	discord      *DiscordProvider
	users        *UsersProvider
	leaderboards *LeaderboardsProvider
}

func NewProfileProvider(storer storage.Storer, discord *DiscordProvider, users *UsersProvider, leaderboards *LeaderboardsProvider) *ProfileProvider {
	return &ProfileProvider{storer: storer, discord: discord, users: users, leaderboards: leaderboards}
}

func (p *ProfileProvider) Get(discordUserID string) (Profile, bool, error) {
	profile := Profile{DiscordUserID: discordUserID, RecentMessages: []ProfileMessage{}}

	hasMember := false
	if member, err := p.discord.Member(discordUserID); err == nil {
		profile.Member = &member
		hasMember = true
		if permissions, err := p.discord.MemberPermissions(discordUserID); err == nil {
			profile.Permissions = &permissions
		}
	}

	user, foundUser, err := p.users.Get(discordUserID)
	if err != nil {
		return Profile{}, false, err
	}
	if foundUser {
		profile.User = &user
	}

	messages, err := p.storer.GetRecentMessagesForUser(discordUserID, profileRecentMessageLimit)
	if err != nil {
		return Profile{}, false, err
	}

	channelNames := p.channelNameMap()
	for _, m := range messages {
		profile.RecentMessages = append(profile.RecentMessages, ProfileMessage{
			Timestamp:   m.Timestamp,
			ChannelID:   m.DiscordChannel,
			ChannelName: channelNames[m.DiscordChannel],
			Message:     m.Message,
			MessageID:   m.DiscordMessageID,
		})
	}

	count, err := p.storer.GetUserMessageCount(discordUserID)
	if err != nil {
		return Profile{}, false, err
	}
	profile.MessageCount = count

	if count > 0 {
		if rank, err := p.leaderboards.UserRank(discordUserID); err == nil {
			profile.Rank = rank
		}
	}

	if !hasMember && !foundUser && count == 0 {
		return Profile{}, false, nil
	}

	return profile, true, nil
}

func (p *ProfileProvider) channelNameMap() map[string]string {
	names := make(map[string]string)
	categories, err := p.discord.Channels()
	if err != nil {
		return names
	}
	for _, category := range categories {
		for _, channel := range category.Channels {
			names[channel.ID] = channel.Name
		}
	}
	return names
}
