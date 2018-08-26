package commands

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
	"sa-mp server",
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
	"ginger",
	"parkour",
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
	"IMPROVED",
	"BASIC",
	"Y_INI",
	"PAWNO",
	"NGG",
	"NGRP",
	"GODFATHER",
	"GF EDIT",
	"25.000 LINES",
	"ABYSS",
	"RCRP",
	"HARAM",
	"SCRATCH",
	"SAMPCTL",
	"CUSTOM OBJECTS",
	"0.3.DL",
	"ROLLERPLAYERS ONLY",
	"SOUTHCLAWS",
	"RAKNET",
	"LAGSHOT",
	"HALAL",
	"STRTOK",
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

func (cm *CommandManager) commandGmName(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	cm.Discord.ChannelMessageSend(message.ChannelID, generateGmName())
	return
}
