package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	config2 "github.com/vehsamrak/slack-standup/internal/app/config"
	"github.com/vehsamrak/slack-standup/internal/app/controller"
	"github.com/vehsamrak/slack-standup/internal/app/meeting"
	"github.com/vehsamrak/slack-standup/internal/app/slack"
	"github.com/vehsamrak/slack-standup/internal/logger"
	"net/http"
	"strconv"
	"sync"
)

type Bot struct {
	slack            *slack.Client
	config           *config2.Config
	answersWaitGroup sync.WaitGroup
}

func (bot Bot) Start(config *config2.Config) {
	log.SetFormatter(&logger.TextFormatter{})
	bot.config = config
	bot.slack = slack.Client{}.Create(config)

	log.Info("Bot started")

	bot.StartMeeting()

	//bot.listenIncomingMessages()
	//bot.answersWaitGroup.Add(1)
	//bot.answersWaitGroup.Wait()

	log.Info("Bot stopped")
}

func (bot *Bot) listenIncomingMessages(standup *meeting.Meeting) {
	slackController := controller.SlackController{}.Create(bot.slack, standup)

	controllerMap := map[string]func(http.ResponseWriter, *http.Request){
		"/":     slackController.Entrypoint,
		"/ping": slackController.Ping,
	}

	for route, controllerFunction := range controllerMap {
		http.HandleFunc(route, controllerFunction)
	}

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(bot.config.Port)}

	log.Infof("Incoming Slack messages listener started on port %d", bot.config.Port)

	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (bot *Bot) StartMeeting() {
	log.Info("Standup meeting created")

	// TODO[petr]: how to receive channel name?
	channelName := "d"
	thread := bot.slack.SendMessageToChannelByName(channelName, "Начинается утренний стэндап!")
	if thread == nil {
		log.Info("Standup meeting ended with error. No thread found")
		return
	}

	standup := meeting.Meeting{}.Create(thread)

	bot.startStandUpInChannel(standup)
	bot.listenIncomingMessages(standup)
}

func (bot *Bot) startStandUpInChannel(standup *meeting.Meeting) {
	// TODO[petr]: how to receive channel name?
	channelName := "d"
	channel := bot.slack.FindChannelByName(channelName)
	users := bot.slack.ChannelUsersList(channel.Id)

	for _, userId := range users.Ids {
		log.Infof("Starting standup for user #%s", userId)
		go bot.startStandUpForUser(standup, userId)
	}
}

func (bot *Bot) startStandUpForUser(standup *meeting.Meeting, userId string) {
	privateUserChannel := bot.slack.OpenChatWithUser(userId)

	if privateUserChannel.Id == "" {
		return
	}

	bot.slack.SendMessageToChannel(privateUserChannel.Id, "Начинается утренний стэндап!")
	bot.slack.SendMessageToChannel(privateUserChannel.Id, fmt.Sprintf("*%s*", standup.QuestionPrevious()))

	// TODO[petr]: enable standup
	// TODO[petr]: disable standup after last question posted

	// TODO[petr]: if all 3 answers, add message to thread
	bot.slack.SendReplyToChannel(standup.Thread.Channel, "Ответы", standup.Thread.Thread)
}
