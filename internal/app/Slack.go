package app

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vehsamrak/slack-standup/internal/app/response/conversationsList"
	"github.com/vehsamrak/slack-standup/internal/app/response/conversationsMembers"
	"github.com/vehsamrak/slack-standup/internal/app/response/conversationsOpen"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Slack struct {
	Config *Config
}

func (Slack) Create(config *Config) *Slack {
	return &Slack{Config: config}
}

func (slack *Slack) botInfo() {
	slack.callApi("bots.info", nil)
}

func (slack *Slack) channelUsersList(channelName string) *conversationsMembers.Users {
	response := slack.callApi("conversations.members", url.Values{"channel": {channelName}})
	//response := []byte("{\"ok\": true, \"members\": [\"U01DB949P8X\",\"U01DNTUGMMK\"],\"response_metadata\": {\"next_cursor\":\"\"}}")
	users := &conversationsMembers.Users{}

	err := json.Unmarshal(response, &users)
	if err != nil {
		panic(err)
	}

	return users
}

func (slack *Slack) channelsList() *conversationsList.ChannelsList {
	response := slack.callApi("conversations.list", url.Values{"exclude_archived": {"true"}})
	channels := &conversationsList.ChannelsList{}

	err := json.Unmarshal(response, &channels)
	if err != nil {
		panic(err)
	}

	return channels
}

func (slack *Slack) findChannelByName(channelName string) *conversationsList.Channel {
	channels := slack.channelsList()
	for _, channel := range channels.Channels {
		if channel.Name == channelName {
			return &channel
		}
	}

	return nil
}

func (slack *Slack) sendMessageToChannel(channel string, message string) {
	slack.callApi(
		"chat.postMessage",
		url.Values{"channel": {channel}, "text": {message}},
	)
}

func (slack *Slack) sendReplyToChannel(channel string, message string, thread string) {
	slack.callApi(
		"chat.postMessage",
		url.Values{"channel": {channel}, "text": {message}, "thread_ts": {thread}},
	)
}

func (slack *Slack) openChatWithUser(userName string) *conversationsOpen.Channel {
	response := slack.callApi(
		"conversations.open",
		url.Values{"users": {userName}},
	)

	conversation := &conversationsOpen.Conversation{}

	err := json.Unmarshal(response, &conversation)
	if err != nil {
		panic(err)
	}

	return &conversation.Channel
}

func (slack *Slack) callApi(method string, parameters url.Values) []byte {
	if parameters == nil {
		parameters = url.Values{}
	}

	log.Infof("Calling API: %s [%s]\n", method, parameters.Encode())

	parameters.Add("token", slack.Config.Token)
	parameters.Add("pretty", "1")

	resp, err := http.Get(fmt.Sprintf("%s/%s?%s", slack.Config.ApiUrl, method, parameters.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Infof("%s\n", body)

	return body
}
