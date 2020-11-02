package app

import (
	log "github.com/sirupsen/logrus"
	config2 "github.com/vehsamrak/slack-standup/internal/app/config"
	"github.com/vehsamrak/slack-standup/internal/app/controller"
	"github.com/vehsamrak/slack-standup/internal/app/slack"
	"github.com/vehsamrak/slack-standup/internal/logger"
	"net/http"
	"strconv"
)

type Bot struct {
	slack  *slack.Client
	config *config2.Config
}

func (bot Bot) Start(config *config2.Config) {
	log.SetFormatter(&logger.TextFormatter{})
	bot.config = config
	bot.slack = slack.Client{}.Create(config)

	log.Info("Bot started")

	bot.listenIncomingMessages()
	//bot.startStandUp()

	log.Info("Bot stopped")
}

func (bot *Bot) listenIncomingMessages() {
	slackController := controller.SlackController{}.Create(bot.slack)

	controllerMap := map[string]func(http.ResponseWriter, *http.Request){
		"/":     slackController.Entrypoint,
		"/ping": slackController.Ping,
	}

	for route, controllerFunction := range controllerMap {
		http.HandleFunc(route, controllerFunction)
	}

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(bot.config.Port)}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (bot *Bot) startStandUp() {
	// TODO[petr]: how to receive channel name?
	channelName := "d"
	channel := bot.slack.FindChannelByName(channelName)
	users := bot.slack.ChannelUsersList(channel.Id)

	for _, userId := range users.Ids {
		privateUserChannel := bot.slack.OpenChatWithUser(userId)

		if privateUserChannel.Id == "" {
			continue
		}

		bot.slack.SendMessageToChannel(privateUserChannel.Id, "Начинается утренний стэндап!")
		bot.slack.SendMessageToChannel(privateUserChannel.Id, "*Удалось выполнить вчерашний план?*")
		bot.slack.SendMessageToChannel(privateUserChannel.Id, "*Что планируешь сделать сегодня?*")
		bot.slack.SendMessageToChannel(privateUserChannel.Id, "*Кто и чем может тебе помочь?*")
	}
}
