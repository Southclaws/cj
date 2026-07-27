package readmodel

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/types"
)

func TestBuildChannelCategoriesGroupsAndSorts(t *testing.T) {
	channels := []*discordgo.Channel{
		{ID: "cat1", Name: "Text Channels", Type: discordgo.ChannelTypeGuildCategory, Position: 0},
		{ID: "cat2", Name: "Voice Channels", Type: discordgo.ChannelTypeGuildCategory, Position: 1},
		{ID: "c2", Name: "general", Type: discordgo.ChannelTypeGuildText, ParentID: "cat1", Position: 1},
		{ID: "c1", Name: "announcements", Type: discordgo.ChannelTypeGuildText, ParentID: "cat1", Position: 0},
		{ID: "c3", Name: "voice-chat", Type: discordgo.ChannelTypeGuildVoice, ParentID: "cat2", Position: 0},
		{ID: "c4", Name: "orphan", Type: discordgo.ChannelTypeGuildText, ParentID: "", Position: 0},
	}

	got := buildChannelCategories(channels)

	if len(got) != 3 {
		t.Fatalf("expected 3 categories (2 real + uncategorised), got %d: %+v", len(got), got)
	}
	if got[0].Name != "Text Channels" || len(got[0].Channels) != 2 {
		t.Fatalf("expected Text Channels with 2 channels first, got %+v", got[0])
	}
	if got[0].Channels[0].Name != "announcements" {
		t.Fatalf("expected channels sorted by position, got %+v", got[0].Channels)
	}
	if got[1].Name != "Voice Channels" {
		t.Fatalf("expected Voice Channels second, got %+v", got[1])
	}
	if got[2].Name != "No category" || len(got[2].Channels) != 1 {
		t.Fatalf("expected uncategorised bucket last, got %+v", got[2])
	}
}

func TestBuildChannelCategoriesOmitsEmptyUncategorised(t *testing.T) {
	channels := []*discordgo.Channel{
		{ID: "cat1", Name: "Text Channels", Type: discordgo.ChannelTypeGuildCategory, Position: 0},
		{ID: "c1", Name: "general", Type: discordgo.ChannelTypeGuildText, ParentID: "cat1", Position: 0},
	}

	got := buildChannelCategories(channels)
	if len(got) != 1 {
		t.Fatalf("expected no uncategorised bucket when everything has a category, got %+v", got)
	}
}

func TestChannelTypeName(t *testing.T) {
	cases := map[discordgo.ChannelType]string{
		discordgo.ChannelTypeGuildText:  "text",
		discordgo.ChannelTypeGuildVoice: "voice",
		discordgo.ChannelTypeGuildForum: "forum",
		discordgo.ChannelType(999):      "other",
	}
	for in, want := range cases {
		if got := channelTypeName(in); got != want {
			t.Fatalf("channelTypeName(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestAuditActionNameKnownAndUnknown(t *testing.T) {
	if got := auditActionName(20); got != "member_kick" {
		t.Fatalf("got %q, want member_kick", got)
	}
	if got := auditActionName(9999); got != "action_9999" {
		t.Fatalf("got %q, want action_9999", got)
	}
}

func TestBuildMemberOwnerCannotBeActedOn(t *testing.T) {
	guild := &discordgo.Guild{ID: "g1", OwnerID: "owner1"}
	roles := []*discordgo.Role{{ID: "g1", Position: 0}}
	member := &discordgo.Member{
		User:  &discordgo.User{ID: "owner1", Username: "TheOwner"},
		Roles: []string{},
	}

	got := buildMember(member, roles, guild, 5)
	if got.CJCanActOn {
		t.Fatalf("expected the guild owner to never be actionable, got %+v", got)
	}
	if len(got.Reasons) == 0 {
		t.Fatal("expected a reason explaining why the owner can't be acted on")
	}
}

func readyDiscordSession() *discord.Session {
	s, err := discordgo.New("Bot test-token")
	if err != nil {
		panic(err)
	}
	s.State = discordgo.NewState()
	s.State.Ready.User = &discordgo.User{ID: "bot-1"}
	return discord.New(s, types.Config{GuildID: "guild-1"})
}

func TestMembersReturnsErrWhenCacheNotReady(t *testing.T) {
	provider := NewDiscordProvider(readyDiscordSession())

	_, err := provider.Members("", 10, 0)
	if !errors.Is(err, ErrMemberCacheNotReady) {
		t.Fatalf("expected ErrMemberCacheNotReady before the cache has ever been populated, got %v", err)
	}
}

func TestFilterMembersByQueryMatchesUsernameNicknameOrID(t *testing.T) {
	members := []*discordgo.Member{
		{User: &discordgo.User{ID: "1", Username: "Alice"}},
		{User: &discordgo.User{ID: "2", Username: "Bob"}, Nick: "Bobby"},
		{User: &discordgo.User{ID: "3", Username: "Charlie"}},
	}

	byUsername := filterMembersByQuery(members, "alic")
	if len(byUsername) != 1 || byUsername[0].User.ID != "1" {
		t.Fatalf("expected a case-insensitive username match, got %+v", byUsername)
	}

	byNickname := filterMembersByQuery(members, "bobby")
	if len(byNickname) != 1 || byNickname[0].User.ID != "2" {
		t.Fatalf("expected a nickname match, got %+v", byNickname)
	}

	byID := filterMembersByQuery(members, "3")
	if len(byID) != 1 || byID[0].User.ID != "3" {
		t.Fatalf("expected an ID match, got %+v", byID)
	}

	empty := filterMembersByQuery(members, "")
	if len(empty) != len(members) {
		t.Fatalf("expected an empty query to return everyone, got %d of %d", len(empty), len(members))
	}

	none := filterMembersByQuery(members, "does-not-exist")
	if len(none) != 0 {
		t.Fatalf("expected no matches, got %+v", none)
	}
}

func TestPaginateMemberSlice(t *testing.T) {
	members := []*discordgo.Member{
		{User: &discordgo.User{ID: "1"}},
		{User: &discordgo.User{ID: "2"}},
		{User: &discordgo.User{ID: "3"}},
	}

	page1 := paginateMemberSlice(members, 2, 0)
	if len(page1) != 2 || page1[0].User.ID != "1" || page1[1].User.ID != "2" {
		t.Fatalf("unexpected first page: %+v", page1)
	}

	page2 := paginateMemberSlice(members, 2, 2)
	if len(page2) != 1 || page2[0].User.ID != "3" {
		t.Fatalf("unexpected second page: %+v", page2)
	}

	beyondEnd := paginateMemberSlice(members, 2, 10)
	if len(beyondEnd) != 0 {
		t.Fatalf("expected an out-of-range offset to return no members, got %+v", beyondEnd)
	}
}

func TestCountMembersByRole(t *testing.T) {
	members := []*discordgo.Member{
		{User: &discordgo.User{ID: "1"}, Roles: []string{"a", "b"}},
		{User: &discordgo.User{ID: "2"}, Roles: []string{"a"}},
		{User: &discordgo.User{ID: "3"}, Roles: []string{}},
	}

	counts := countMembersByRole(members)
	if counts["a"] != 2 {
		t.Fatalf("expected role a to have 2 members, got %d", counts["a"])
	}
	if counts["b"] != 1 {
		t.Fatalf("expected role b to have 1 member, got %d", counts["b"])
	}
	if counts["c"] != 0 {
		t.Fatalf("expected an unreferenced role to have 0 members, got %d", counts["c"])
	}
}

func TestBuildMemberHierarchyCheck(t *testing.T) {
	guild := &discordgo.Guild{ID: "g1", OwnerID: "someone-else"}
	roles := []*discordgo.Role{
		{ID: "g1", Position: 0},
		{ID: "high-role", Position: 10},
		{ID: "low-role", Position: 1},
	}

	highMember := &discordgo.Member{User: &discordgo.User{ID: "m1"}, Roles: []string{"high-role"}}
	if got := buildMember(highMember, roles, guild, 5); got.CJCanActOn {
		t.Fatalf("expected a member with a higher role than CJ to be unactionable, got %+v", got)
	}

	lowMember := &discordgo.Member{User: &discordgo.User{ID: "m2"}, Roles: []string{"low-role"}}
	if got := buildMember(lowMember, roles, guild, 5); !got.CJCanActOn {
		t.Fatalf("expected a member with a lower role than CJ to be actionable, got %+v", got)
	}
}
