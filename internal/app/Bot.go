package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/vehsamrak/slack-standup/internal/app/config"
	"github.com/vehsamrak/slack-standup/internal/app/controller"
	"github.com/vehsamrak/slack-standup/internal/app/meeting"
	"github.com/vehsamrak/slack-standup/internal/app/slack"
	"github.com/vehsamrak/slack-standup/internal/logger"
	"gopkg.in/robfig/cron.v2"
	"net/http"
	"strconv"
)

type Bot struct {
	slack                 *slack.Client
	config                *config.Config
	participantMeetingMap map[string]*meeting.Meeting
}

func (bot Bot) Start(config *config.Config) {
	log.SetFormatter(&logger.TextFormatter{})
	bot.config = config
	bot.slack = slack.Client{}.Create(config)

	bot.participantMeetingMap = make(map[string]*meeting.Meeting)

	log.Info("Bot started")

	startMeetingEvent := make(chan string)

	channelNames := config.ChannelNames
	for _, channelName := range channelNames {
		cron := cron.New()
		cron.AddFunc(config.MeetingSchedule, func() { startMeetingEvent <- channelName })
		cron.Start()
	}

	go bot.listenIncomingMessages()

	for channelName := range startMeetingEvent {
		bot.StartMeeting(channelName)
	}

	log.Info("Bot stopped")
}

func (bot *Bot) StartMeeting(channelName string) {
	log.Info("Standup meeting created")

	thread := bot.slack.SendMessageToChannelByName(channelName, meeting.Meeting{}.Greetings())
	if thread == nil {
		log.Info("Standup meeting ended with error. No thread found")
		return
	}

	standup := meeting.Meeting{}.Create(channelName, thread)

	channel := bot.slack.FindChannelByName(standup.ChannelName())
	users := bot.slack.ChannelUsersList(channel.Id)

	for _, userId := range users.Ids {
		user := bot.slack.UserInfo(userId)
		log.Infof("Starting standup for user \"%s\" #%s", user.Name, user.Id)
		standup.Participants[userId] = &meeting.Questions{}
		bot.participantMeetingMap[userId] = standup

		go bot.startMeetingForUser(userId)
	}
}

func (bot *Bot) startMeetingForUser(userId string) {
	privateUserChannel := bot.slack.OpenChatWithUser(userId)

	if privateUserChannel.Id == "" {
		return
	}

	bot.slack.SendMessageToChannel(privateUserChannel.Id, meeting.Questions{}.Greetings())
	bot.slack.SendMessageToChannel(privateUserChannel.Id, meeting.Questions{}.QuestionPrevious())
}

func (bot *Bot) listenIncomingMessages() {
	for route, controllerFunction := range bot.createControllerMap() {
		http.HandleFunc(route, controllerFunction)
	}

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(bot.config.Port)}

	log.Infof("Incoming Slack messages listener started on port %d", bot.config.Port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

func (bot *Bot) createControllerMap() map[string]func(http.ResponseWriter, *http.Request) {
	slackController := controller.SlackController{}.Create(bot.slack, bot.participantMeetingMap)

	controllerMap := map[string]func(http.ResponseWriter, *http.Request){
		"/":     slackController.Entrypoint,
		"/ping": slackController.Ping,
	}

	return controllerMap
}
