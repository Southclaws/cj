package main

import (
	"bytes"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var coolwordsxd = []string{
	"ultra",
	"mega",
	"extreme",
	"extr3m3",
	"elite",
	"l33t",
	"pro",
	"profesional",
	"amazing",
	"gold",
	"platinum",
	"fusion",
	"infusion",
	"xtreme",
	"exiting",
	"fucking",
	"wow",
	"amazeballs",
	"reloaded",
	"advanced",
	"you're",
	"mom",
	"edit",
	"haram",
}

var morewords = []string{
	"gaming",
	"killaz",
	"community",
	"group",
	"gang",
	"clan",
	"pros",
	"profesionals",
	"motherfuckers",
	"scripters",
	"hackers",
	"español",
	"revolution",
	"krisk",
	"abyss",
	"mom",
	"rcrp-leak",
}

var gamemodes = []string{
	"roleplay",
	"rp",
	"rpg",
	"game",
	"gangwars",
	"deathmatch",
	"dm",
	"tdm",
	"racing",
	"sex",
	"race",
	"derby",
	"minigames",
}

var tags = []string{
	"UNIQUE",
	"DYNAMIC",
	"REFUNDING",
	"UCP",
	"MYSQL",
	"SSCANF",
	"ZCMD",
	"STRCMP",
	"PAWN",
	"C++",
	"YLESS",
	"0.3e+",
	"RUS",
	"ROLEYPLAY",
	"OFFICIAL",
	"ZOMBIES",
	"BETA",
	"BEST",
	"HIRING",
	"Y_INI",
	"PAWNO",
	"NGG",
	"NGRP"
	"GODFATHER",
	"GF EDIT",
	"25.000 LINES",
	"ABYSS",
	"RCRP",
	"HARAM",
	"SCRATCH",
}

func init() {
	rand.Seed(time.Now().Unix())
}

func generateGmName() (result string) {
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
		buf.WriteString(string(anotherword))
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
	case 4:
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

func commandGmName(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	cm.App.discordClient.ChannelMessageSend(message.ChannelID, generateGmName())
	return true, false, nil
}
