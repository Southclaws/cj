package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var firsts = []string{
	"adorable",
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
}
var seconds = []string{
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
}

func mpname() string {
	mp := []byte("00-MP")
	first := firsts[rand.Intn(len(firsts))]
	second := seconds[rand.Intn(len(seconds))]
	mp[0] = []byte(strings.ToUpper(first))[0]
	mp[1] = []byte(strings.ToUpper(second))[0]
	return fmt.Sprintf("%s: %s %s Multiplayer", string(mp), strings.Title(first), strings.Title(second))
}

func commandMP(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {

	cm.App.discordClient.ChannelMessageSend(message.ChannelID, mpname())
	return true, false, nil
}
