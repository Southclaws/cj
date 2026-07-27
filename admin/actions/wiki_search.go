package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Southclaws/cj/bot/commands"
)

type WikiSearchAction struct{}

func NewWikiSearchAction() *WikiSearchAction {
	return &WikiSearchAction{}
}

func (a *WikiSearchAction) Name() string { return "wiki.search" }

func (a *WikiSearchAction) Description() string {
	return "Search the open.mp documentation, the same search /wiki uses."
}

func (a *WikiSearchAction) Risk() RiskLevel { return RiskSafe }

type wikiSearchInput struct {
	Term string `json:"term"`
}

func decodeWikiSearchInput(input json.RawMessage) (wikiSearchInput, error) {
	var decoded wikiSearchInput
	if len(input) > 0 {
		if err := json.Unmarshal(input, &decoded); err != nil {
			return wikiSearchInput{}, fmt.Errorf("%w: input must be a JSON object with a term field", ErrInvalidInput)
		}
	}
	if len(decoded.Term) < 3 {
		return wikiSearchInput{}, fmt.Errorf("%w: term must be at least 3 characters", ErrInvalidInput)
	}
	return decoded, nil
}

func (a *WikiSearchAction) Preview(ctx context.Context, input json.RawMessage) (Preview, error) {
	decoded, err := decodeWikiSearchInput(input)
	if err != nil {
		return Preview{}, err
	}
	return Preview{Summary: fmt.Sprintf("This will search the open.mp documentation for %q.", decoded.Term)}, nil
}

func (a *WikiSearchAction) Execute(ctx context.Context, input json.RawMessage) (Result, error) {
	decoded, err := decodeWikiSearchInput(input)
	if err != nil {
		return Result{}, err
	}

	hits, err := commands.SearchWiki(decoded.Term)
	if err != nil {
		return Result{}, err
	}

	if len(hits) == 0 {
		return Result{Summary: fmt.Sprintf("No documentation results for %q.", decoded.Term)}, nil
	}

	results := make([]map[string]any, len(hits))
	for i, hit := range hits {
		results[i] = map[string]any{
			"pageName":    hit.PageName,
			"url":         hit.URL,
			"category":    hit.Category,
			"description": hit.Description,
		}
	}

	return Result{
		Summary: fmt.Sprintf("Found %d result(s) for %q.", len(hits), decoded.Term),
		Detail:  map[string]any{"hits": results},
	}, nil
}
