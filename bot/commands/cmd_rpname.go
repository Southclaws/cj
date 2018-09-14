package commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var rpfirsts = []string{
	"los",
	"san",
	"one",
	"las",
	"bulgarian",
	"bone",
	"carl",
	"bay",
	"texas",
	"evolve",
	"halal",
	"pure",
	"CJ",
	"gay",
	"payday",
	"red",
	"leaked",
	"next",
	"role",
	"scavenge",
	"infinity",
	"sexy",
	"O.G.",
	"god",
	"godfather",
	"capital",
	"SAMP",
	"grand",
}
var rpseconds = []string{
	"area",
	"andreas",
	"johnson",
	"gangstas",
	"play",
	"county",
	"cops",
	"life",
	"day",
	"one",
	"halal",
	"parrot",
	"turtle",
	"world",
	"game",
	"christian",
	"SAMP",
	"survive",
	"ginger",
	"larceny",
}

func rpname() string {
	mp := []byte("00RP")
	first := rpfirsts[rand.Intn(len(rpfirsts))]
	second := rpseconds[rand.Intn(len(rpseconds))]
	mp[0] = []byte(strings.ToUpper(first))[0]
	mp[1] = []byte(strings.ToUpper(second))[0]
	return fmt.Sprintf("%s: %s %s Roleplay", string(mp), strings.Title(first), strings.Title(second))
}

func (cm *CommandManager) commandRP(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {

	_, err = cm.Discord.ChannelMessageSend(message.ChannelID, rpname())
	return
}
