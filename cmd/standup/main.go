package main

import (
	"github.com/vehsamrak/slack-standup/internal/app"
	config2 "github.com/vehsamrak/slack-standup/internal/app/config"
)

func main() {
	config := config2.Config{}.Load()
	app.Bot{}.Start(config)
}
