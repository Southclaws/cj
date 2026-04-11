package commands

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

//go:embed animations.json
var animationsJSON []byte

type animEntry struct {
	Library string `json:"library"`
	Name    string `json:"name"`
}

var (
	animData       []animEntry
	validAnimParam = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

func init() {
	if err := json.Unmarshal(animationsJSON, &animData); err != nil {
		panic("failed to parse animations.json: " + err.Error())
	}
}

func (cm *CommandManager) commandAnimationAutocomplete(
	interaction *discordgo.InteractionCreate,
) {
	data := interaction.ApplicationCommandData()

	var focusedField string
	var typedLibrary string
	var typedValue string

	for _, opt := range data.Options {
		if opt.Focused {
			focusedField = opt.Name
			typedValue = opt.StringValue()
		}
		if opt.Name == "library" {
			typedLibrary = opt.StringValue()
		}
	}

	var choices []*discordgo.ApplicationCommandOptionChoice

	switch focusedField {
	case "library":
		seen := map[string]bool{}
		lower := strings.ToLower(typedValue)
		for _, a := range animData {
			if seen[a.Library] {
				continue
			}
			if lower == "" || strings.Contains(strings.ToLower(a.Library), lower) {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  a.Library,
					Value: a.Library,
				})
				seen[a.Library] = true
			}
			if len(choices) >= 25 {
				break
			}
		}

	case "animation":
		lower := strings.ToLower(typedValue)
		lowerLib := strings.ToLower(typedLibrary)
		for _, a := range animData {
			if lowerLib != "" && strings.ToLower(a.Library) != lowerLib {
				continue
			}
			if lower == "" || strings.Contains(strings.ToLower(a.Name), lower) {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  a.Name,
					Value: a.Name,
				})
			}
			if len(choices) >= 25 {
				break
			}
		}
	}

	cm.Discord.S.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func (cm *CommandManager) commandAnimation(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	library := args["library"].StringValue()
	animation := args["animation"].StringValue()

	if !validAnimParam.MatchString(library) || !validAnimParam.MatchString(animation) {
		cm.replyDirectly(interaction, "Invalid library or animation name. Only letters, numbers and underscores are allowed.")
		return
	}

	videoURL := fmt.Sprintf("https://assets.open.mp/assets/anims/%s/%s.webm", library, animation)

	content := fmt.Sprintf(
		"**%s / %s**\n```\nApplyAnimation(playerid, \"%s\", \"%s\", 4.1, false, false, false, false, SYNC_ALL);\n```\n%s",
		library, animation, library, animation, videoURL,
	)

	cm.replyDirectly(interaction, content)
	return false, nil
}
