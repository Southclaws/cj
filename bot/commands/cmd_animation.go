package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

type animEntry struct {
	Library string `json:"library"`
	Name    string `json:"name"`
}

var (
	animationLibraries      []string
	animationNames          []string
	animationNamesByLibrary map[string][]string
	validAnimParam          = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

func init() {
	animationLibraries, animationNames, animationNamesByLibrary = buildAnimationIndexes(animationEntries)
}

func buildAnimationIndexes(entries []animEntry) ([]string, []string, map[string][]string) {
	libraries := make([]string, 0)
	names := make([]string, 0, len(entries))
	namesByLibrary := make(map[string][]string)
	seenLibraries := make(map[string]struct{})

	for _, entry := range entries {
		libraryKey := strings.ToLower(entry.Library)
		if _, seen := seenLibraries[libraryKey]; !seen {
			libraries = append(libraries, entry.Library)
			seenLibraries[libraryKey] = struct{}{}
		}

		names = append(names, entry.Name)
		namesByLibrary[libraryKey] = append(namesByLibrary[libraryKey], entry.Name)
	}

	return libraries, names, namesByLibrary
}

func autocompleteChoices(values []string, typedValue string) []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, min(len(values), 25))
	lower := strings.ToLower(typedValue)

	for _, value := range values {
		if lower == "" || strings.Contains(strings.ToLower(value), lower) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  value,
				Value: value,
			})
		}
		if len(choices) >= 25 {
			break
		}
	}

	return choices
}

func autocompleteOptionString(opt *discordgo.ApplicationCommandInteractionDataOption) string {
	if opt == nil || opt.Value == nil {
		return ""
	}

	value, ok := opt.Value.(string)
	if !ok {
		return ""
	}

	return value
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
			typedValue = autocompleteOptionString(opt)
		}
		if opt.Name == "library" {
			typedLibrary = autocompleteOptionString(opt)
		}
	}

	var choices []*discordgo.ApplicationCommandOptionChoice

	switch focusedField {
	case "library":
		choices = autocompleteChoices(animationLibraries, typedValue)

	case "animation":
		if typedLibrary != "" {
			choices = autocompleteChoices(animationNamesByLibrary[strings.ToLower(typedLibrary)], typedValue)
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
