package app

import (
	"fmt"
)

type Bot struct {
	slack *Slack
}

func (bot Bot) Start(config *Config) {
	fmt.Printf("Bot started\n")

	bot.slack = Slack{}.Create(config)

	// TODO[petr]: how to receive channel name?
	channelName := "d"
	channel := bot.slack.findChannelByName(channelName)
	users := bot.slack.channelUsersList(channel.Id)

	for _, userId := range users.Ids {
		privateUserChannel := bot.slack.openChatWithUser(userId)

		if privateUserChannel.Id == "" {
			continue
		}

		bot.slack.sendMessageToChannel(privateUserChannel.Id, "Начинается утренний стэндап!")
		bot.slack.sendMessageToChannel(privateUserChannel.Id, "*Удалось выполнить вчерашний план?*")
		bot.slack.sendMessageToChannel(privateUserChannel.Id, "*Что планируешь сделать сегодня?*")
		bot.slack.sendMessageToChannel(privateUserChannel.Id, "*Кто и чем может тебе помочь?*")
	}

	fmt.Printf("Bot stopped\n")
}
