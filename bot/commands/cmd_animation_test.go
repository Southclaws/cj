package commands

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

func TestBuildAnimationIndexesPreservesOrderAndGrouping(t *testing.T) {
	t.Parallel()

	libraries, names, namesByLibrary := buildAnimationIndexes([]animEntry{
		{Library: "BAR", Name: "first"},
		{Library: "benchpress", Name: "lift"},
		{Library: "BAR", Name: "second"},
		{Library: "bar", Name: "third"},
	})

	require.Equal(t, []string{"BAR", "benchpress"}, libraries)
	require.Equal(t, []string{"first", "lift", "second", "third"}, names)
	require.Equal(t, []string{"first", "second", "third"}, namesByLibrary["bar"])
	require.Equal(t, []string{"lift"}, namesByLibrary["benchpress"])
}

func TestAutocompleteChoicesFiltersCaseInsensitively(t *testing.T) {
	t.Parallel()

	choices := autocompleteChoices([]string{"BAR", "BASEBALL", "benchpress"}, "BaSe")

	require.Len(t, choices, 1)
	require.Equal(t, "BASEBALL", choices[0].Value)
}

func TestAutocompleteChoicesCapsAtTwentyFive(t *testing.T) {
	t.Parallel()

	values := make([]string, 30)
	for ii := range values {
		values[ii] = "value"
	}

	choices := autocompleteChoices(values, "")

	require.Len(t, choices, 25)
}

func TestAutocompleteOptionString(t *testing.T) {
	t.Parallel()

	require.Equal(t, "", autocompleteOptionString(nil))
	require.Equal(t, "", autocompleteOptionString(&discordgo.ApplicationCommandInteractionDataOption{}))
	require.Equal(t, "", autocompleteOptionString(&discordgo.ApplicationCommandInteractionDataOption{Value: 42}))
	require.Equal(t, "BAR", autocompleteOptionString(&discordgo.ApplicationCommandInteractionDataOption{Value: "BAR"}))
}
