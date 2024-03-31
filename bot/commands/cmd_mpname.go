package commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Southclaws/cj/types"
)

var firsts = []string{
	"CJ",
	"O.G.",
	"SAMP",
	"adorable",
	"bay",
	"bone",
	"bulgarian",
	"capital",
	"carl",
	"evolve",
	"gay",
	"god",
	"godfather",
	"halal",
	"infinity",
	"las",
	"leaked",
	"mom",
	"next",
	"one",
	"payday",
	"pisd",
	"pure",
	"red",
	"role",
	"san",
	"scavenge",
	"sexy",
	"texas",
}
var seconds = []string{
	"SAMP",
	"andreas",
	"area",
	"christian",
	"cops",
	"county",
	"day",
	"game",
	"gangstas",
	"ginger",
	"halal",
	"johnson",
	"life",
	"one",
	"parrot",
	"pisd",
	"play",
	"survive",
	"turtle",
	"world",
}

func mpname() string {
	mp := []byte("00-MP")
	first := firsts[rand.Intn(len(firsts))]
	second := seconds[rand.Intn(len(seconds))]
	mp[0] = []byte(strings.ToUpper(first))[0]
	mp[1] = []byte(strings.ToUpper(second))[0]
	title := cases.Title(language.English)
	return fmt.Sprintf("%s: %s %s Multiplayer", string(mp), title.String(first), title.String(second))
}

func (cm *CommandManager) commandMP(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {

	cm.replyDirectly(interaction, mpname())
	return
}
