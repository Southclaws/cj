package commands

import (
	"bytes"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

var coolwordsxd = []string{
	"advanced",
	"amazeballs",
	"amazing",
	"amk",
	"burger",
	"🍔",
	"edit",
	"elite",
	"exiting",
	"extr3m3",
	"extreme",
	"fucking",
	"fusion",
	"gold",
	"haram",
	"infusion",
	"l33t",
	"mega",
	"mom",
	"open",
	"pisd",
	"platinum",
	"pro",
	"profesional",
	"reloaded",
	"sa-mp server",
	"ultra",
	"wow",
	"xd",
	"xtreme",
	"you're",
	"kungkingkang",
	"mengheran",
	"mengakak",
	"menghadeh",
	"Ayonima",
	"Acumalaka",
	"Megawati-Chan",
	"Puan-Chan",
}

var morewords = []string{
	"abyss",
	"barp-leak",
	"black",
	"clan",
	"community",
	"español",
	"gaming",
	"gang",
	"ginger",
	"group",
	"hackers",
	"killaz",
	"krisk",
	"mom",
	"motherfuckers",
	"open.mp",
	"parkour",
	"pisd",
	"profesionals",
	"pros",
	"rcrp-leak",
	"revolution",
	"shoters",
	"scripters",
	"white",
	"you(them)tubers",
	"mengontol",
	"kontol",
	"pepeq",
}

var gamemodes = []string{
	"deathmatch",
	"derby",
	"dm",
	"freeroam",
	"game",
	"gangbang",
	"gangwars",
	"minigames",
	"pisd",
	"race",
	"racing",
	"roleplay",
	"dmrp",
	"rp",
	"rpg",
	"sex",
	"tdm",
	"war",
}

var tags = []string{
	"0.3.DL",
	"0.3e+",
	"25.000 LINES",
	"ABYSS",
	"BASIC",
	"BEST",
	"BETA",
	"C++",
	"CUSTOM OBJECTS",
	"DYNAMIC",
	"GF EDIT",
	"GO",
	"POWERED BY G00GLE",
	"GODFATHER",
	"HALAL",
	"HARAM",
	"HIRING",
	"IMPROVED",
	"LAGSHOT",
	"LUA",
	"MOMS",
	"MYSQL",
	"NGG",
	"NGRP",
	"OFFICIAL",
	"OPEN.MP",
	"PAWN",
	"PAWNO",
	"PISD",
	"RAKNET",
	"RCRP",
	"REDIS",
	"REFUNDING",
	"ROLEYPLAY",
	"ROLLERPLAYERS ONLY",
	"RUS",
	"SAMPCTL",
	"SCRATCH",
	"SOUTHCLAWS",
	"SSCANF",
	"STRCMP",
	"STRCMP2",
	"STRTOK",
	"STRTOK2",
	"TELNET",
	"UCP",
	"UNIQUE",
	"YLESS",
	"Y_INI",
	"ZCMD",
	"ZOMBIES",
	"MISEBAHXD",
}

func init() {}

func generateGmName() (result string) {
	rand.New(rand.NewSource(time.Now().Unix()))
	buf := bytes.NewBuffer(nil)

	coolword := coolwordsxd[rand.Intn(len(coolwordsxd))]
	switch rand.Intn(3) {
	case 0:
		buf.WriteString(coolword)
	case 1:
		buf.WriteString(strings.ToTitle(coolword))
	case 2:
		buf.WriteString(strings.ToUpper(coolword))
	}
	buf.WriteString(" ")

	anotherword := morewords[rand.Intn(len(morewords))]
	switch rand.Intn(3) {
	case 0:
		buf.WriteString(anotherword)
	case 1:
		buf.WriteString(strings.ToTitle(anotherword))
	case 2:
		buf.WriteString(strings.ToUpper(anotherword))
	}
	buf.WriteString(" ")

	gamemodes := gamemodes[rand.Intn(len(gamemodes))]
	switch rand.Intn(4) {
	case 0:
		buf.WriteString(gamemodes)
	case 1:
		buf.WriteString(strings.ToTitle(gamemodes))
	case 2:
		buf.WriteString(strings.ToUpper(gamemodes))
	case 3:
		buf.WriteString("[" + strings.ToUpper(gamemodes) + "]")
	}

	tagCount := rand.Intn(10)
	usableTags := make(map[string]bool)
	for i := 0; i < tagCount; i++ {
		usableTags[tags[rand.Intn(len(tags))]] = true
	}

	for tag := range usableTags {
		buf.WriteString(" [" + tag + "]")
	}

	return buf.String()
}

func (cm *CommandManager) commandGmName(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.replyDirectly(interaction, generateGmName())
	return
}
