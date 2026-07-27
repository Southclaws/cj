package readmodel

import (
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Southclaws/cj/storage"
)

type fakeUsersStorer struct {
	storage.Memory
	users     []storage.User
	byID      map[string]storage.User
	messages  []storage.ChatLog
	listErr   error
	lastQuery string
}

func (f *fakeUsersStorer) ListUsers(query string, limit, offset int) ([]storage.User, int, error) {
	f.lastQuery = query
	if f.listErr != nil {
		return nil, 0, f.listErr
	}
	total := len(f.users)
	if offset >= total {
		return []storage.User{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return f.users[offset:end], total, nil
}

func (f *fakeUsersStorer) GetUser(id string) (storage.User, bool, error) {
	user, ok := f.byID[id]
	return user, ok, nil
}

func (f *fakeUsersStorer) GetRecentMessagesForUser(id string, limit int) ([]storage.ChatLog, error) {
	if len(f.messages) > limit {
		return f.messages[:limit], nil
	}
	return f.messages, nil
}

func (f *fakeUsersStorer) SearchMessages(id, query string) ([]storage.ChatLog, error) {
	var out []storage.ChatLog
	for _, m := range f.messages {
		if m.Message == query {
			out = append(out, m)
		}
	}
	return out, nil
}

func (f *fakeUsersStorer) SearchAllMessages(query string, limit int) ([]storage.ChatLog, error) {
	var out []storage.ChatLog
	for _, m := range f.messages {
		if m.Message == query {
			out = append(out, m)
		}
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (f *fakeUsersStorer) GetMessageByID(messageID string) (storage.ChatLog, error) {
	for _, m := range f.messages {
		if m.DiscordMessageID == messageID {
			return m, nil
		}
	}
	return storage.ChatLog{}, mongo.ErrNoDocuments
}

func TestUsersProviderListDefaultsAndClampsLimit(t *testing.T) {
	storer := &fakeUsersStorer{users: []storage.User{{DiscordUserID: "1"}, {DiscordUserID: "2"}}}
	provider := NewUsersProvider(storer, nil)

	page, err := provider.List("", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Limit != defaultUsersPageSize {
		t.Fatalf("expected default limit %d, got %d", defaultUsersPageSize, page.Limit)
	}
	if page.Total != 2 || len(page.Users) != 2 {
		t.Fatalf("expected both users returned, got %+v", page)
	}

	page, err = provider.List("", 10000, -5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Limit != defaultUsersPageSize {
		t.Fatalf("expected an oversized limit to be clamped to the default, got %d", page.Limit)
	}
	if page.Offset != 0 {
		t.Fatalf("expected a negative offset to be clamped to 0, got %d", page.Offset)
	}
}

func TestUsersProviderNormalizesNilReceivedReactions(t *testing.T) {
	storer := &fakeUsersStorer{
		users: []storage.User{{DiscordUserID: "1", ReceivedReactions: nil}},
		byID:  map[string]storage.User{"1": {DiscordUserID: "1", ReceivedReactions: nil}},
	}
	provider := NewUsersProvider(storer, nil)

	page, err := provider.List("", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Users[0].ReceivedReactions == nil {
		t.Fatal("expected ReceivedReactions to be normalized to an empty slice in List, got nil")
	}

	user, _, err := provider.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ReceivedReactions == nil {
		t.Fatal("expected ReceivedReactions to be normalized to an empty slice in Get, got nil")
	}
}

func TestUsersProviderListPassesQueryThrough(t *testing.T) {
	storer := &fakeUsersStorer{users: []storage.User{{DiscordUserID: "1"}}}
	provider := NewUsersProvider(storer, nil)

	if _, err := provider.List("alice", 10, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if storer.lastQuery != "alice" {
		t.Fatalf("expected the query to be passed through to the storer, got %q", storer.lastQuery)
	}
}

func TestUsersProviderGetPassesThrough(t *testing.T) {
	storer := &fakeUsersStorer{byID: map[string]storage.User{"42": {DiscordUserID: "42", ForumUserName: "someone"}}}
	provider := NewUsersProvider(storer, nil)

	user, found, err := provider.Get("42")
	if err != nil || !found || user.ForumUserName != "someone" {
		t.Fatalf("got user=%+v found=%v err=%v", user, found, err)
	}

	_, found, err = provider.Get("does-not-exist")
	if err != nil || found {
		t.Fatalf("expected not found, got found=%v err=%v", found, err)
	}
}

func TestChatProviderRequiresUserIDOrQuery(t *testing.T) {
	provider := NewChatProvider(&fakeUsersStorer{}, nil)

	_, err := provider.Messages("", "")
	if !errors.Is(err, ErrChatFilterRequired) {
		t.Fatalf("expected ErrChatFilterRequired, got %v", err)
	}
}

func TestChatProviderSearchesGloballyWithoutAUserID(t *testing.T) {
	storer := &fakeUsersStorer{messages: []storage.ChatLog{
		{DiscordUserID: "1", Message: "hello"},
		{DiscordUserID: "2", Message: "hello"},
		{DiscordUserID: "3", Message: "goodbye"},
	}}
	provider := NewChatProvider(storer, nil)

	results, err := provider.Messages("", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected a global search to find matches across users, got %+v", results)
	}
}

func TestChatProviderGetByID(t *testing.T) {
	storer := &fakeUsersStorer{messages: []storage.ChatLog{
		{DiscordUserID: "1", DiscordMessageID: "msg-1", Message: "hi"},
	}}
	provider := NewChatProvider(storer, nil)

	message, found, err := provider.GetByID("msg-1")
	if err != nil || !found || message.Message != "hi" {
		t.Fatalf("got message=%+v found=%v err=%v", message, found, err)
	}

	_, found, err = provider.GetByID("does-not-exist")
	if err != nil || found {
		t.Fatalf("expected not found, got found=%v err=%v", found, err)
	}
}

func TestChatProviderSearchesWhenQueryProvided(t *testing.T) {
	storer := &fakeUsersStorer{messages: []storage.ChatLog{
		{DiscordUserID: "1", Message: "hello"},
		{DiscordUserID: "1", Message: "goodbye"},
	}}
	provider := NewChatProvider(storer, nil)

	all, err := provider.Messages("1", "")
	if err != nil || len(all) != 2 {
		t.Fatalf("expected all messages without a query, got %+v err=%v", all, err)
	}

	filtered, err := provider.Messages("1", "hello")
	if err != nil || len(filtered) != 1 || filtered[0].Message != "hello" {
		t.Fatalf("expected only the matching message, got %+v err=%v", filtered, err)
	}
}
