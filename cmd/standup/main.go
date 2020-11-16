package main

import (
	"github.com/vehsamrak/slack-standup/internal/app"
	"github.com/vehsamrak/slack-standup/internal/app/config"
)

func main() {
	app.Bot{}.Start(config.Config{}.Load())
}
