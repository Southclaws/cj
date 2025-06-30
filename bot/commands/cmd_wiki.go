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

func truncateText(s string, max int) string {
	if s == "" {
		return "Documentation available - click to view"
	}
	if max > len(s) {
		return s
	}
	
	// Find last space before max length
	truncated := s[:max]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > max/2 {
		return truncated[:lastSpace] + "..."
	}
	
	return truncated + "..."
}

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
					SetHitsPerPage(5).
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
			Color:       0xED4245, // Red color for no results
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

		// Get URL
		urlWithoutAnchor, ok := hitData["url_without_anchor"].(string)
		if !ok {
			continue
		}

		// Skip if already processed this URL - Algolia does this thing of sending twice the same thing
		if seenUrls[urlWithoutAnchor] {
			continue
		}
		seenUrls[urlWithoutAnchor] = true

		// Parse URL to get page info - Algolia doesn't give the Function/Callback name - so I steal it from the URL
		stringParts := strings.Split(strings.TrimSuffix(urlWithoutAnchor, "/"), "/")
		if len(stringParts) < 3 {
			continue
		}

		// Skip blog posts and other non-documentation pages - Remove blog posts from results
		if len(stringParts) >= 4 && (stringParts[len(stringParts)-2] == "blog" || 
			stringParts[len(stringParts)-3] == "blog") {
			continue
		}

		actuallyFoundResults++

		// Get page/function name
		pageName := stringParts[len(stringParts)-1]
		
		// Get description from multiple possible sources
		description := extractDescription(hitData)
		
		// Build category path
		category := buildCategory(stringParts)

		// Format the result
		desc.WriteString(fmt.Sprintf(
			"[**%s**](%s)\n%s: %s\n\n",
			pageName,
			urlWithoutAnchor,
			category,
			truncateText(description, 120),
		))

		// Limit to 3 results for readability
		if actuallyFoundResults >= 3 {
			break
		}
	}

	// Hitting this means all results were filtered out
	if actuallyFoundResults == 0 {
		cm.replyDirectlyEmbed(interaction, "", &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       fmt.Sprintf("No results: %s", searchTerm),
			Description: "No documentation found for that query. Try a different search term.",
			Color:       0xED4245,
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("Documentation Search Results for %s", searchTerm),
		Description: desc.String(),
		Color:       0x5865F2, // Discord blurple
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Click the links to view full documentation",
		},
	}
	cm.replyDirectlyEmbed(interaction, "", embed)

	return false, err
}

func extractDescription(hitData map[string]interface{}) string {
	// Try to get description from various fields in order of preference
	descriptionSources := []string{
		"description",  // Most likely to have the frontmatter description
		"content",      // Page content
		"excerpt",      // Page excerpt
		"text",         // Text content
		"summary",      // Summary
	}
	
	for _, field := range descriptionSources {
		if value, exists := hitData[field]; exists {
			if strValue, ok := value.(string); ok && strings.TrimSpace(strValue) != "" {
				cleaned := cleanDescription(strValue)
				if cleaned != "" {
					return cleaned
				}
			}
		}
	}
	
	// Try hierarchy for context
	if hierarchy, exists := hitData["hierarchy"]; exists {
		if hierarchyMap, ok := hierarchy.(map[string]interface{}); ok {
			// Try different hierarchy levels
			for _, level := range []string{"lvl0", "lvl1", "lvl2", "lvl3"} {
				if lvlValue, exists := hierarchyMap[level]; exists {
					if strValue, ok := lvlValue.(string); ok && strings.TrimSpace(strValue) != "" {
						return cleanDescription(strValue)
					}
				}
			}
		}
	}
	
	// Try anchor text
	if anchor, exists := hitData["anchor"]; exists {
		if strValue, ok := anchor.(string); ok && strings.TrimSpace(strValue) != "" {
			return cleanDescription(strValue)
		}
	}
	
	return "Documentation available"
}

func cleanDescription(text string) string {
	if text == "" {
		return ""
	}
	
	// Remove HTML tags and decode entities
	cleaned := html.UnescapeString(text)
	
	// Remove markdown-style highlighting
	cleaned = strings.ReplaceAll(cleaned, "<mark>", "**")
	cleaned = strings.ReplaceAll(cleaned, "</mark>", "**")
	
	// Clean up whitespace and newlines
	cleaned = strings.ReplaceAll(cleaned, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\r", " ")
	
	// Remove multiple spaces
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}
	
	cleaned = strings.TrimSpace(cleaned)
	
	// Remove common prefixes that aren't useful
	prefixesToRemove := []string{
		"Description:",
		"Description",
		"## Description",
		"###",
		"##",
		"#",
	}
	
	for _, prefix := range prefixesToRemove {
		if strings.HasPrefix(cleaned, prefix) {
			cleaned = strings.TrimSpace(strings.TrimPrefix(cleaned, prefix))
		}
	}
	
	return cleaned
}

func buildCategory(urlParts []string) string {
	if len(urlParts) < 4 {
		return "Documentation"
	}
	
	// Extract meaningful category parts
	var categoryParts []string
	
	// Skip the domain parts and get the meaningful path
	for i := 3; i < len(urlParts)-1; i++ {
		part := urlParts[i]
		
		// Skip common path parts that aren't useful
		if part == "docs" || part == "en" || part == "" {
			continue
		}
		
		// Capitalize and clean up the part
		part = strings.ReplaceAll(part, "-", " ")
		part = strings.Title(part)
		categoryParts = append(categoryParts, part)
	}
	
	if len(categoryParts) == 0 {
		return "Documentation"
	}
	
	return strings.Join(categoryParts, " › ")
}