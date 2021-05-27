package commands

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

const cmdUsage = "USAGE: /wiki [function/callback]"

type Results struct {
	Status struct {
		Total      int `json:"total"`
		Failed     int `json:"failed"`
		Successful int `json:"successful"`
	} `json:"status"`
	Request struct {
		Query struct {
			Query string `json:"query"`
		} `json:"query"`
		Size      int `json:"size"`
		From      int `json:"from"`
		Highlight struct {
			Style  interface{} `json:"style"`
			Fields interface{} `json:"fields"`
		} `json:"highlight"`
		Fields           interface{} `json:"fields"`
		Facets           interface{} `json:"facets"`
		Explain          bool        `json:"explain"`
		Sort             []string    `json:"sort"`
		IncludeLocations bool        `json:"includeLocations"`
		SearchAfter      interface{} `json:"search_after"`
		SearchBefore     interface{} `json:"search_before"`
	} `json:"request"`
	Hits      []Hit       `json:"hits"`
	TotalHits int         `json:"total_hits"`
	MaxScore  float64     `json:"max_score"`
	Took      int64       `json:"took"`
	Facets    interface{} `json:"facets"`
}

type Hit struct {
	Index     string  `json:"index"`
	ID        string  `json:"id"`
	Score     float64 `json:"score"`
	Locations struct {
		Description struct {
			Position []struct {
				Pos            int         `json:"pos"`
				Start          int         `json:"start"`
				End            int         `json:"end"`
				ArrayPositions interface{} `json:"array_positions"`
			} `json:"position"`
		} `json:"Description"`
	} `json:"locations"`
	Fragments struct {
		Description []string `json:"Description"`
	} `json:"fragments"`
	Sort []string `json:"sort"`
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

	r, err := http.Get(fmt.Sprintf("https://api.open.mp/docs/search?q=%s", strings.ReplaceAll(searchTerm, " ", "%20")))
	if err != nil {
		cm.replyDirectly(interaction, fmt.Sprintf("Failed to GET result for search term %s\nError: %s", searchTerm, err.Error()))
		return
	}

	var results Results
	if err = json.NewDecoder(r.Body).Decode(&results); err != nil {
		cm.replyDirectly(interaction, fmt.Sprintf("Failed to decode result for search term %s\nError: %s\n", searchTerm, err.Error()))
		return
	}

	if results.TotalHits == 0 {
		cm.replyDirectlyEmbed(interaction, "", &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       fmt.Sprintf("No results: %s", searchTerm),
			Description: "There were no results for that query.",
		})
		return
	}

	desc := strings.Builder{}

	rendered := 0
	for _, hit := range results.Hits {
		if rendered == 3 {
			break
		}

		// Skip searching translations
		if strings.Contains(hit.ID, "translations") {
			continue
		}

		desc.WriteString(fmt.Sprintf(
			"[%s](https://open.mp/%s): %s\n",
			nameFromPath(hit.ID),
			strings.TrimSuffix(hit.ID, ".md"),
			formatDescription(hit)))
		rendered++
	}
	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       fmt.Sprintf("Documentation Search Results: %s", searchTerm),
		Description: desc.String(),
	}
	cm.replyDirectlyEmbed(interaction, "", embed)

	return false, err // Todo: remove this
}

func nameFromPath(p string) string {
	return strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
}

func formatDescription(hit Hit) string {
	if len(hit.Fragments.Description) == 0 {
		return "(No description found)"
	}

	return html.UnescapeString(strings.ReplaceAll(
		strings.ReplaceAll(
			hit.Fragments.Description[0],
			"<mark>", "**"),
		"</mark>", "**"))
}
