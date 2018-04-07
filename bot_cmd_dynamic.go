package main

import (
	"bytes"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var scripts = []string{
	"biz",
	"Update",
	"Material Text",
	"Backpack",
	"House",
	"Fire-Bin",
	"Garage",
	"System",
	"Vehicle & Dealership",
	"House",
	"Gate",
	"ATM's",
	"Stingers",
	"Categories",
	"ice cream creation",
	"Org Creation",
	"stores",
	"job creation",
	"Street w/ sign",
	"Gates",
	"Server Signature Generator",
	"circles",
	"house creating!",
	"player account data",
	"Garbage Collector",
	"Interiors",
	"Admin-Vehicle Lock - for all u lazy ppl!!!",
	"Teleport",
	"Gang",
	"Guild",
	"animations",
	"Business. v1.0",
	"position save system",
	"News",
	"GPS",
	"Vehicle Creator",
	"Dialog Maker",
	"ATM System",
	"ATM System",
	"House-Business System",
	"Dialog Gang System",
	"IconMaker",
	"Rules ← Using HTTP!",
	"Car Ownership System",
	"AFK",
	"EPIC HOUSE",
	"Company",
	"business",
	"Entrance",
	"EnEx",
	"bazookas system",
	"myth creator",
	"Buy Weapon And Kits",
	"VEHICLES",
	"weapon control",
	"Door system",
	"Interior",
	"weapon shop",
	"Menus",
	"Gang",
	"Guild",
	"bomb system",
	"mom",
	"race system",
	"Vehicle Spawn Menu",
	"Checkpoints",
	"Media Dialog",
	"Vehicle System",
}

var features = []string{
	"TextDraw",
	"Object",
	"3D Label",
	"Pickup",
	"Map icon",
	"loading & saving",
	"SQL",
	"for Roleplay",
	"DJSON",
	"MySQL",
	"User Friendly",
	"0.3d Compatible",
	"1 line",
	"Can also be used for missions!",
	"Draft",
	"Saves + Loads Through MySQL!",
	"with advanced anti db",
	"Scripting SpeedArt Video",
	"zcmd",
	"SQLite",
	"streamer",
	"Dynamic!",
	" ← USING HTTP!",
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
		buf.WriteString(string(script))
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
			buf.WriteString(string(script))
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
		buf.WriteString("System")
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

func commandDynamic(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	cm.App.discordClient.ChannelMessageSend(message.ChannelID, makeDynamic())
	return true, false, nil
}
