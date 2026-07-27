package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Southclaws/cj/discord"
)

type RefreshMembersAction struct {
	session *discord.Session
}

func NewRefreshMembersAction(session *discord.Session) *RefreshMembersAction {
	return &RefreshMembersAction{session: session}
}

func (a *RefreshMembersAction) Name() string { return "discord.refresh-members" }

func (a *RefreshMembersAction) Description() string {
	return "Re-fetch every guild member from Discord, refreshing the cache the Members and Roles pages read from."
}

func (a *RefreshMembersAction) Risk() RiskLevel { return RiskStateChanging }

func (a *RefreshMembersAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	return Preview{
		Summary: "This will re-fetch every guild member from Discord. For a large guild this makes several API requests and can take a while.",
	}, nil
}

func (a *RefreshMembersAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	if err := a.session.RefreshMembers(); err != nil {
		return Result{}, err
	}

	members, cachedAt := a.session.Members()

	return Result{
		Summary: fmt.Sprintf("Refreshed the member cache: %d members loaded.", len(members)),
		Detail: map[string]any{
			"memberCount": len(members),
			"cachedAt":    cachedAt.UTC().Format(time.RFC3339),
		},
	}, nil
}
