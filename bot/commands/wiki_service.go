package commands

import (
	"encoding/json"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

type WikiHit struct {
	PageName    string
	URL         string
	Category    string
	Description string
}

func SearchWiki(term string) ([]WikiHit, error) {
	algoliaClient, err := search.NewClient(algoliaAppID, algoliaAPIKey)
	if err != nil {
		return nil, err
	}

	response, err := algoliaClient.Search(algoliaClient.NewApiSearchRequest(
		search.NewEmptySearchMethodParams().SetRequests(
			[]search.SearchQuery{*search.SearchForHitsAsSearchQuery(
				search.NewEmptySearchForHits().
					SetIndexName(algoliaIndexName).
					SetQuery(term).
					SetHitsPerPage(5).
					SetFilters("language:en"),
			)},
		),
	))
	if err != nil {
		return nil, err
	}

	finalResult := response.Results[0]
	if *finalResult.SearchResponse.NbHits == 0 {
		return nil, nil
	}

	hits := []WikiHit{}
	seenUrls := make(map[string]bool)

	for _, hit := range finalResult.SearchResponse.Hits {
		var hitData map[string]interface{}
		hitJSON, err := json.Marshal(hit)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(hitJSON, &hitData); err != nil {
			continue
		}

		urlWithoutAnchor, ok := hitData["url_without_anchor"].(string)
		if !ok {
			continue
		}

		if seenUrls[urlWithoutAnchor] {
			continue
		}
		seenUrls[urlWithoutAnchor] = true

		stringParts := strings.Split(strings.TrimSuffix(urlWithoutAnchor, "/"), "/")
		if len(stringParts) < 3 {
			continue
		}

		if len(stringParts) >= 4 && (stringParts[len(stringParts)-2] == "blog" ||
			stringParts[len(stringParts)-3] == "blog") {
			continue
		}

		pageName := stringParts[len(stringParts)-1]
		description := extractDescription(hitData)
		category := buildCategory(stringParts)

		hits = append(hits, WikiHit{
			PageName:    pageName,
			URL:         urlWithoutAnchor,
			Category:    category,
			Description: truncateText(description, 120),
		})

		if len(hits) >= 3 {
			break
		}
	}

	return hits, nil
}
