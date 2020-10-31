package app

import "fmt"

type Bot struct{}

func (bot Bot) Start() {
	fmt.Print("Bot started")
	fmt.Print("Bot stopped")
}
