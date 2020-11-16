package controller

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vehsamrak/slack-standup/internal/app/controller/requests"
	"github.com/vehsamrak/slack-standup/internal/app/event/messageIm"
	"github.com/vehsamrak/slack-standup/internal/app/meeting"
	"github.com/vehsamrak/slack-standup/internal/app/slack"
	"io/ioutil"
	"net/http"
)

const eventTypeMessage = "message"

type SlackController struct {
	Controller
	slack          *slack.Client
	participantMap map[string]*meeting.Meeting
}

func (controller SlackController) Create(
	slack *slack.Client,
	participantMap map[string]*meeting.Meeting,
) *SlackController {
	return &SlackController{slack: slack, participantMap: participantMap}
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
		log.Infof("Private channel was not opened and has no id: %#v", event)
		controller.Respond(response, "", http.StatusBadRequest)
		return
	}

	standup := controller.participantMap[userId]
	if standup == nil {
		log.Infof("Standup is not running for now: %#v", event)
		controller.Respond(response, "", http.StatusOK)
		return
	}

	questions := standup.Participants[userId]
	if questions == nil {
		log.Errorf("Users questions were not found for %s", userId)
		controller.Respond(response, "", http.StatusBadRequest)
		return
	}

	if questions.Previous == "" {
		questions.Previous = event.Message.Text
		controller.slack.SendMessageToChannel(privateUserChannel.Id, meeting.Questions{}.QuestionToday())
	} else if questions.Today == "" {
		questions.Today = event.Message.Text
		controller.slack.SendMessageToChannel(privateUserChannel.Id, meeting.Questions{}.QuestionBlock())
	} else if questions.Block == "" {
		questions.Block = event.Message.Text
		controller.slack.SendMessageToChannel(privateUserChannel.Id, "Спасибо, хорошего дня!")

		controller.slack.SendReplyToChannel(
			standup.Thread.Channel,
			controller.createMeetingResultMessage(userId, questions),
			standup.Thread.Thread,
		)
	}

	controller.Respond(response, "", http.StatusOK)
}

func (controller *SlackController) createMeetingResultMessage(userId string, questions *meeting.Questions) string {
	user := controller.slack.UserInfo(userId)
	return fmt.Sprintf("@%s\n%s", user.Name, questions.Result())
}
