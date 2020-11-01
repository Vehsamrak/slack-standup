package app

import (
	"fmt"
)

type Bot struct {
}

func (Bot) Start(config *Config) {
	fmt.Printf("Bot started\n")

	slack := Slack{}.Create(config)

	slack.sendMessageToChannel("C01E4KK06DP", "test from golang!")

	fmt.Printf("Bot stopped\n")
}
