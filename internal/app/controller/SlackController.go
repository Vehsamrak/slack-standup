package controller

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vehsamrak/slack-standup/internal/app/controller/requests"
	"github.com/vehsamrak/slack-standup/internal/app/event/messageIm"
	"github.com/vehsamrak/slack-standup/internal/app/slack"
	"io/ioutil"
	"net/http"
)

const eventTypeMessage = "message"

type SlackController struct {
	Controller
	slack *slack.Client
}

func (controller *SlackController) Ping(response http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	ping := &requests.Ping{}
	err = json.Unmarshal(body, &ping)
	if err != nil {
		controller.Respond(response, "", http.StatusBadRequest)
		return
	}

	log.Infof("Ping requested: %v", string(body))

	controller.Respond(response, fmt.Sprintf("{\"challenge\":\"%s\"}", ping.Challenge), http.StatusOK)
}

func (controller *SlackController) Entrypoint(response http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	log.Infof("Incoming event: %v", string(body))

	event := &messageIm.Event{}
	err = json.Unmarshal(body, &event)
	if err != nil {
		controller.Respond(response, "", http.StatusBadRequest)
		return
	}

	if event.Message.Type != eventTypeMessage {
		controller.Respond(response, "", http.StatusBadRequest)
		return
	}

	if event.Message.BotId != "" {
		log.Infof("User is bot. Skipping interaction: %#v", event)
		controller.Respond(response, "", http.StatusOK)
		return
	}

	userId := event.Message.UserId
	if userId == "" {
		log.Infof("User has no id: %#v", event)
		controller.Respond(response, "", http.StatusBadRequest)
		return
	}

	privateUserChannel := controller.slack.OpenChatWithUser(userId)
	if privateUserChannel.Id == "" {
		return
	}

	//controller.slack.SendMessageToChannel(privateUserChannel.Id, "Начинается утренний стэндап!")
	//controller.slack.SendMessageToChannel(privateUserChannel.Id, "*Удалось выполнить предыдущий план?*")
	//controller.slack.SendMessageToChannel(privateUserChannel.Id, "*Что планируешь сделать сегодня?*")
	//controller.slack.SendMessageToChannel(privateUserChannel.Id, "*Кто и чем может тебе в этом помочь?*")

	controller.Respond(response, "", http.StatusOK)
}

func (controller SlackController) Create(slack *slack.Client) *SlackController {
	return &SlackController{slack: slack}
}
