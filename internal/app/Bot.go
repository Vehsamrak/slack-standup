package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/vehsamrak/slack-standup/internal/app/controller"
	"github.com/vehsamrak/slack-standup/internal/logger"
	"net/http"
	"strconv"
)

type Bot struct {
	slack  *Slack
	config *Config
}

func (bot Bot) Start(config *Config) {
	log.SetFormatter(&logger.TextFormatter{})
	bot.config = config
	bot.slack = Slack{}.Create(config)

	log.Info("Bot started")

	bot.listenIncomingMessages()
	//bot.startStandUp()

	log.Info("Bot stopped")
}

func (bot *Bot) listenIncomingMessages() {
	controllerMap := map[string]func(http.ResponseWriter, *http.Request){
		"/":     controller.SlackController{}.Entrypoint,
		"/ping": controller.SlackController{}.Ping,
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
}
