package commands

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"

	"go.uber.org/zap"
)

const cmdUsage = "USAGE: /wiki [function/callback]"

// Algolia config
const (
	algoliaAppID            = "AOKXGK39Z7"
	algoliaAPIKey           = "54204f37e5c8fc2871052d595ee0505e" //Safe to commit
	algoliaIndexName        = "open"
	algoliaContextualSearch = true
	algoliaInsights         = false
)

func (cm *CommandManager) commandWiki(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	searchTerm := args["search-term"].StringValue()
	if len(searchTerm) < 3 {
		cm.replyDirectly(interaction, "Query must be 3 characters or more")
		return
	}

	algoliaClient, algoliaErr := search.NewClient(algoliaAppID, algoliaAPIKey)

	if algoliaErr != nil {
		cm.replyDirectly(interaction, fmt.Sprintf("Failed to initialise Search client: \n%s", algoliaErr.Error()))
		zap.L().Error("Failed to initialise Search client", zap.Error(algoliaErr))
		return
	}

	response, err := algoliaClient.Search(algoliaClient.NewApiSearchRequest(
		search.NewEmptySearchMethodParams().SetRequests(
			[]search.SearchQuery{*search.SearchForHitsAsSearchQuery(
				search.NewEmptySearchForHits().
					SetIndexName(algoliaIndexName).
					SetQuery(searchTerm).
					SetHitsPerPage(3).
					SetFilters("language:en"),
			)},
		),
	))
	if err != nil {
		cm.replyDirectly(interaction, fmt.Sprintf("Failed to search: %s", err.Error()))
		return
	}

	finalResult := response.Results[0]

	if *finalResult.SearchResponse.NbHits == 0 {
		cm.replyDirectlyEmbed(interaction, "", &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       fmt.Sprintf("No results: %s", searchTerm),
			Description: "There were no results for that query.",
		})
		return
	}

	desc := strings.Builder{}

	seenUrls := make(map[string]bool)

	actuallyFoundResults := 0
	for _, hit := range finalResult.SearchResponse.Hits {
		var hitData map[string]interface{}
		hitJSON, err := json.Marshal(hit)
		if err != nil {
			// Would reply, but then it may error again - and we cant re-send messages
			// cm.replyDirectly(interaction, fmt.Sprintf("Error marshalling hit: %s", err.Error()))
			continue
		}

		if err := json.Unmarshal(hitJSON, &hitData); err != nil {
			// Would reply, but then it may error again - and we cant re-send messages
			// cm.replyDirectly(interaction, fmt.Sprintf("Error unmarshalling hit: %s", err.Error()))
			continue
		}

		if seenUrls[hitData["url_without_anchor"].(string)] { //Already presented to user - Algolia does this thing of sending twice the same thing
			continue
		}

		seenUrls[hitData["url_without_anchor"].(string)] = true

		stringParts := strings.Split(strings.TrimSuffix(hitData["url_without_anchor"].(string), "/"), "/") //Algolia doesnt give the Function/Callback name - so I steal it from the URL

		if stringParts[len(stringParts)-2] == "blog" { //Remove blog posts from results
			continue
		}

		actuallyFoundResults++

		content, ok := hitData["content"].(string)
		description := ""
		if !ok {
			description = "(No description found)"
		} else {
			description = formatDescription(&content)
		}

		desc.WriteString(fmt.Sprintf(
			"[%s](%s) [%s/%s]: %s\n",
			stringParts[len(stringParts)-1],
			hitData["url_without_anchor"].(string),
			stringParts[len(stringParts)-3],
			stringParts[len(stringParts)-2],
			description,
		))
	}

	if actuallyFoundResults == 0 { //Hitting this means all results were filtered out
		cm.replyDirectlyEmbed(interaction, "", &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       fmt.Sprintf("No results: %s", searchTerm),
			Description: "There were no results for that query.",
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("Documentation Search Results: %s", searchTerm),
		Description: desc.String(),
	}
	cm.replyDirectlyEmbed(interaction, "", embed)

	return false, err // Todo: remove this
}

func formatDescription(hit *string) string {
	return html.UnescapeString(strings.ReplaceAll(
		strings.ReplaceAll(
			*hit,
			"<mark>", "**"),
		"</mark>", "**"))
}
