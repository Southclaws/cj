package discord

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func memberBatch(ids ...string) []*discordgo.Member {
	out := make([]*discordgo.Member, len(ids))
	for i, id := range ids {
		out[i] = &discordgo.Member{User: &discordgo.User{ID: id, Username: "user-" + id}}
	}
	return out
}

func TestPaginateAllMembersStopsOnShortPage(t *testing.T) {
	var afters []string

	batch1 := make([]*discordgo.Member, memberPageSize)
	for i := range batch1 {
		batch1[i] = &discordgo.Member{User: &discordgo.User{ID: string(rune('a' + i%26))}}
	}

	calls := 0
	members, err := paginateAllMembers(func(after string) ([]*discordgo.Member, error) {
		afters = append(afters, after)
		calls++
		switch calls {
		case 1:
			return batch1, nil
		case 2:
			return memberBatch("last1", "last2"), nil
		default:
			t.Fatalf("expected exactly 2 calls, got a 3rd")
			return nil, nil
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != memberPageSize+2 {
		t.Fatalf("expected %d members, got %d", memberPageSize+2, len(members))
	}
	if afters[0] != "" {
		t.Fatalf("expected the first call to have no cursor, got %q", afters[0])
	}
	if afters[1] != batch1[len(batch1)-1].User.ID {
		t.Fatalf("expected the second call's cursor to be the last member of the first page, got %q", afters[1])
	}
}

func TestPaginateAllMembersSinglePage(t *testing.T) {
	calls := 0
	members, err := paginateAllMembers(func(after string) ([]*discordgo.Member, error) {
		calls++
		return memberBatch("1", "2", "3"), nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(members))
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call for a page shorter than the page size, got %d", calls)
	}
}

func TestPaginateAllMembersPropagatesError(t *testing.T) {
	wantErr := errors.New("rate limited")
	_, err := paginateAllMembers(func(after string) ([]*discordgo.Member, error) {
		return nil, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected the fetch error to propagate, got %v", err)
	}
}

func TestSessionMembersReturnsSnapshotNotSharedSlice(t *testing.T) {
	s := &Session{members: memberBatch("1", "2")}

	snapshot, _ := s.Members()
	snapshot[0] = &discordgo.Member{User: &discordgo.User{ID: "mutated"}}

	again, _ := s.Members()
	if again[0].User.ID == "mutated" {
		t.Fatal("expected Members() to return a copy, not the live internal slice")
	}
}

func TestSessionMembersZeroTimeBeforeFirstRefresh(t *testing.T) {
	s := &Session{}
	_, updatedAt := s.Members()
	if !updatedAt.IsZero() {
		t.Fatal("expected a zero updatedAt before the cache has ever been populated")
	}
}
