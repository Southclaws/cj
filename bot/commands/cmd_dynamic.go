package commands

import (
	"bytes"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/types"
)

var scripts = []string{
	"3D Tryg",
	"AFK",
	"ATM Machine System",
	"ATM Machine's",
	"ATM System",
	"ATM's",
	"Admin-Vehicle Lock - for all u lazy ppl!!!",
	"Animations",
	"Backpack",
	"Bash",
	"Bazookas System",
	"Biz",
	"Bomb System",
	"Boobs",
	"Business",
	"Business v1.0",
	"Business v1.0-R1",
	"Business v1.0-R2",
	"Business v1.0-R3",
	"Business v1.0-R3.1",
	"Buy Weapon And Kits",
	"Car Ownership System",
	"Categories",
	"Checkpoints",
	"Circles",
	"Company",
	"Dialog Gang Bang System",
	"Dialog Gang System",
	"Dialog Maker",
	"Dildo System",
	"Door System",
	"EPIC HOUSE",
	"EnEx",
	"Entrance",
	"Fire-Bin",
	"GPS",
	"Gang Bang",
	"Gang",
	"Garage",
	"Garbage Collector",
	"Gate",
	"Gates",
	"Guild System",
	"Guild",
	"Hot Coffee",
	"House System",
	"House Creating!",
	"House",
	"House-Biz System",
	"House-Business System",
	"Ice Cream Creation",
	"IconMaker",
	"Interior",
	"Interiors",
	"Job Creation",
	"Large Arrays",
	"Material Text",
	"Media Dialog",
	"Menus",
	"Milf",
	"Milf Gang Bang",
	"Mom",
	"Mom Gang Bang",
	"Myth Creator",
	"News",
	"Org Creation",
	"PISD",
	"Player Account Data",
	"Player Enumerator",
	"Position save system",
	"PowerShell",
	"Race system",
	"Rules ← Using HTTP!",
	"Semi Dynamic System",
	"Server Signature Generator",
	"Short Arrays",
	"Stingers",
	"Stores",
	"Street w/ sign",
	"System",
	"Teleport",
	"Telnet Client",
	"Telnet Server",
	"Update",
	"VEHICLES",
	"Vehicle & Dealership",
	"Vehicle Creator",
	"Vehicle Spawn Menu",
	"Vehicle System",
	"Vertify",
	"Vibrator System",
	"Weapon control",
	"Weapon shop",
	"GM 2 JUTA TCUYY!!!",
	"PRO KONGDIR",
	"#include <matthew.pwn>",
	"WAGYU A5😋",
}

var features = []string{
	" ← USING HTTP!",
	"0.3d Compatible",
	"1 line",
	"3D Label",
	"Can also be used for missions!",
	"DJSON",
	"Draft",
	"Dynamic!",
	"Map icon",
	"MSSQL",
	"MySQL",
	"Object",
	"Pickup",
	"SQL",
	"SQLite",
	"Saves + Loads Through MySQL!",
	"Scripting SpeedArt Video",
	"TextDraw",
	"User Friendly",
	"account vertification",
	"for Roleplay",
	"loading & saving",
	"pisd",
	"semi dynamicness",
	"streamer",
	"with advanced anti db",
	"zcmd",
	"FiveM Abiss",
	"#REVIEWJUJUR",
	"YAHAHAWAHYU",
}

func makeDynamic() string {
	buf := bytes.NewBuffer(nil)

	style := rand.Intn(4)

	switch style {
	case 0:
		buf.WriteString("dynamic")
	case 1:
		buf.WriteString("Dynamic")
	case 2:
		buf.WriteString("DYNAMIC")
	case 3:
		buf.WriteString("[DYNAMIC]")
	}
	buf.WriteString(" ")

	script := scripts[rand.Intn(len(scripts))]
	switch rand.Intn(3) {
	case 0:
		buf.WriteString(script)
	case 1:
		buf.WriteString(strings.ToTitle(script))
	case 2:
		buf.WriteString(strings.ToUpper(script))
	}
	buf.WriteString(" ")

	if rand.Intn(100) < 50 {
		buf.WriteString("and ")
		script := scripts[rand.Intn(len(scripts))]
		switch rand.Intn(3) {
		case 0:
			buf.WriteString(script)
		case 1:
			buf.WriteString(strings.ToTitle(script))
		case 2:
			buf.WriteString(strings.ToUpper(script))
		}
		buf.WriteString(" ")
	}

	if rand.Intn(100) < 50 {
		buf.WriteString(" ! ")
	}

	switch style {
	case 0:
		buf.WriteString("system")
	case 1:
		buf.WriteString("System")
	case 2:
		buf.WriteString("SYSTEM")
	case 3:
		buf.WriteString("[SYSTEM]")
	}

	tagCount := rand.Intn(3)
	usableTags := make(map[string]bool)
	for i := 0; i < tagCount; i++ {
		usableTags[features[rand.Intn(len(features))]] = true
	}

	for tag := range usableTags {
		if rand.Intn(2) == 0 {
			buf.WriteString(" [" + tag + "]")
		} else {
			buf.WriteString(" " + tag)
		}
	}

	if rand.Intn(10000) == 1 {
		buf.WriteString(" with extra turtles!")
	}

	return buf.String()
}

func (cm *CommandManager) commandDynamic(
	interaction *discordgo.InteractionCreate,
	args map[string]*discordgo.ApplicationCommandInteractionDataOption,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	cm.replyDirectly(interaction, makeDynamic())
	return
}
