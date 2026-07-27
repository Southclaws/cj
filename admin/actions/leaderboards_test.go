package actions

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Southclaws/cj/admin/readmodel"
	"github.com/Southclaws/cj/storage"
)

type fakeLeaderboardsBackend struct {
	storage.Memory
	topMessages  storage.TopMessages
	topReactions []storage.TopReactionEntry
	rank         int
}

func (f *fakeLeaderboardsBackend) GetTopMessages(top int) (storage.TopMessages, error) {
	return f.topMessages, nil
}

func (f *fakeLeaderboardsBackend) GetTopReactions(top int, reaction string) ([]storage.TopReactionEntry, error) {
	return f.topReactions, nil
}

func (f *fakeLeaderboardsBackend) GetUserRank(discordUserID string) (int, error) {
	return f.rank, nil
}

func TestTopMessagesActionExecute(t *testing.T) {
	backend := &fakeLeaderboardsBackend{topMessages: storage.TopMessages{{User: "1", Messages: 10}}}
	action := NewTopMessagesAction(readmodel.NewLeaderboardsProvider(backend, nil))

	result, err := action.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Summary == "" {
		t.Fatal("expected a non-empty summary")
	}
}

func TestTopReactionsActionExecute(t *testing.T) {
	backend := &fakeLeaderboardsBackend{topReactions: []storage.TopReactionEntry{{UserID: "1", Counter: 5, Reaction: "thumbsup"}}}
	action := NewTopReactionsAction(readmodel.NewLeaderboardsProvider(backend, nil))

	result, err := action.Execute(context.Background(), json.RawMessage(`{"reaction":"thumbsup"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Summary == "" {
		t.Fatal("expected a non-empty summary")
	}
}

func TestUserRankActionRequiresUserID(t *testing.T) {
	backend := &fakeLeaderboardsBackend{}
	action := NewUserRankAction(readmodel.NewLeaderboardsProvider(backend, nil))

	_, err := action.Execute(context.Background(), json.RawMessage(`{}`))
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUserRankActionExecute(t *testing.T) {
	backend := &fakeLeaderboardsBackend{rank: 7}
	action := NewUserRankAction(readmodel.NewLeaderboardsProvider(backend, nil))

	result, err := action.Execute(context.Background(), json.RawMessage(`{"userId":"42"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Detail["rank"] != 7 {
		t.Fatalf("expected rank 7 in the result detail, got %+v", result.Detail)
	}
}

func TestWikiSearchActionRejectsShortTerm(t *testing.T) {
	action := NewWikiSearchAction()

	_, err := action.Preview(context.Background(), json.RawMessage(`{"term":"ab"}`))
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for a too-short term, got %v", err)
	}
}

func TestWikiSearchActionRejectsMissingTerm(t *testing.T) {
	action := NewWikiSearchAction()

	_, err := action.Preview(context.Background(), json.RawMessage(`{}`))
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for a missing term, got %v", err)
	}
}

func TestCommandNameInputRequiresCommand(t *testing.T) {
	if _, err := decodeCommandNameInput(json.RawMessage(`{}`)); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for a missing command, got %v", err)
	}
	if _, err := decodeCommandNameInput(json.RawMessage(`{"command":"/wiki"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetCommandSettingsInputRequiresCommand(t *testing.T) {
	if _, err := decodeSetCommandSettingsInput(json.RawMessage(`{"settings":{}}`)); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for a missing command, got %v", err)
	}
}
