package main

import "github.com/vehsamrak/slack-standup/internal/app"

func main() {
	config := app.Config{}.Load()
	app.Bot{}.Start(config)
}
