package main

import (
	"fmt"
	"os"

	// loads environment variables from .env
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"

	"github.com/Southclaws/cj/bot"
	"github.com/Southclaws/cj/types"
)

var version = "master"

func main() {
	config := &types.Config{}
	err := envconfig.Process("CJ", config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.Version = version

	bot.Start(config)
}
