package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Southclaws/cj/admin/readmodel"
)

type limitInput struct {
	Limit int `json:"limit"`
}

func decodeLimitInput(input json.RawMessage) (limitInput, error) {
	var decoded limitInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return limitInput{}, fmt.Errorf("%w: input must be a JSON object with a limit field", ErrInvalidInput)
		}
	}
	return decoded, nil
}

type TopMessagesAction struct {
	leaderboards *readmodel.LeaderboardsProvider
}

func NewTopMessagesAction(p *readmodel.LeaderboardsProvider) *TopMessagesAction {
	return &TopMessagesAction{leaderboards: p}
}

func (a *TopMessagesAction) Name() string { return "leaderboards.top-messages" }

func (a *TopMessagesAction) Description() string {
	return "Show the users with the most recorded messages, the same ranking /top uses."
}

func (a *TopMessagesAction) Risk() RiskLevel { return RiskReadOnly }

func (a *TopMessagesAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	if _, err := decodeLimitInput(input); err != nil {
		return Preview{}, err
	}
	return Preview{Summary: "This will list the users with the most recorded messages (default 10 if no limit is given)."}, nil
}

func (a *TopMessagesAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeLimitInput(input)
	if err != nil {
		return Result{}, err
	}

	top, err := a.leaderboards.TopMessages(decoded.Limit)
	if err != nil {
		return Result{}, err
	}

	entries := make([]map[string]any, len(top))
	for i, entry := range top {
		entries[i] = map[string]any{
			"userId":    entry.Person.ID,
			"username":  entry.Person.Username,
			"avatarUrl": entry.Person.AvatarURL,
			"messages":  entry.Messages,
		}
	}

	return Result{
		Summary: fmt.Sprintf("Loaded the top %d users by message count.", len(top)),
		Detail:  map[string]any{"kind": "leaderboard", "entries": entries},
	}, nil
}

type TopReactionsAction struct {
	leaderboards *readmodel.LeaderboardsProvider
}

func NewTopReactionsAction(p *readmodel.LeaderboardsProvider) *TopReactionsAction {
	return &TopReactionsAction{leaderboards: p}
}

func (a *TopReactionsAction) Name() string { return "leaderboards.top-reactions" }

func (a *TopReactionsAction) Description() string {
	return "Show the users with the most received reactions, the same ranking /toprep uses."
}

func (a *TopReactionsAction) Risk() RiskLevel { return RiskReadOnly }

type topReactionsInput struct {
	Limit    int    `json:"limit"`
	Reaction string `json:"reaction"`
}

func decodeTopReactionsInput(input json.RawMessage) (topReactionsInput, error) {
	var decoded topReactionsInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return topReactionsInput{}, fmt.Errorf("%w: input must be a JSON object with limit and reaction fields", ErrInvalidInput)
		}
	}
	return decoded, nil
}

func (a *TopReactionsAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	decoded, err := decodeTopReactionsInput(input)
	if err != nil {
		return Preview{}, err
	}
	if decoded.Reaction != "" {
		return Preview{Summary: fmt.Sprintf("This will list the users with the most %q reactions.", decoded.Reaction)}, nil
	}
	return Preview{Summary: "This will list the users with the most reactions overall."}, nil
}

func (a *TopReactionsAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeTopReactionsInput(input)
	if err != nil {
		return Result{}, err
	}

	top, err := a.leaderboards.TopReactions(decoded.Limit, decoded.Reaction)
	if err != nil {
		return Result{}, err
	}

	entries := make([]map[string]any, len(top))
	for i, entry := range top {
		entries[i] = map[string]any{
			"userId":    entry.Person.ID,
			"username":  entry.Person.Username,
			"avatarUrl": entry.Person.AvatarURL,
			"counter":   entry.Counter,
			"reaction":  entry.Reaction,
		}
	}

	return Result{
		Summary: fmt.Sprintf("Loaded the top %d users by reaction count.", len(top)),
		Detail:  map[string]any{"kind": "leaderboard", "entries": entries},
	}, nil
}

type UserRankAction struct {
	leaderboards *readmodel.LeaderboardsProvider
}

func NewUserRankAction(p *readmodel.LeaderboardsProvider) *UserRankAction {
	return &UserRankAction{leaderboards: p}
}

func (a *UserRankAction) Name() string { return "leaderboards.user-rank" }

func (a *UserRankAction) Description() string {
	return "Look up one user's message-count rank, the same as /mytop but for any user."
}

func (a *UserRankAction) Risk() RiskLevel { return RiskReadOnly }

type userIDInput struct {
	UserID string `json:"userId"`
}

func decodeUserIDInput(input json.RawMessage) (userIDInput, error) {
	var decoded userIDInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return userIDInput{}, fmt.Errorf("%w: input must be a JSON object with a userId field", ErrInvalidInput)
		}
	}
	if decoded.UserID == "" {
		return userIDInput{}, fmt.Errorf("%w: a userId is required", ErrInvalidInput)
	}
	return decoded, nil
}

func (a *UserRankAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	decoded, err := decodeUserIDInput(input)
	if err != nil {
		return Preview{}, err
	}
	return Preview{Summary: fmt.Sprintf("This will look up the message-count rank for user %s.", decoded.UserID)}, nil
}

func (a *UserRankAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeUserIDInput(input)
	if err != nil {
		return Result{}, err
	}

	rank, err := a.leaderboards.UserRank(decoded.UserID)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Summary: fmt.Sprintf("User %s is ranked #%d by message count.", decoded.UserID, rank),
		Detail:  map[string]any{"rank": rank},
	}, nil
}
