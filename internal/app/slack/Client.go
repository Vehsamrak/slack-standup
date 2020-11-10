package slack

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vehsamrak/slack-standup/internal/app/config"
	"github.com/vehsamrak/slack-standup/internal/app/response/chatPostMessage"
	"github.com/vehsamrak/slack-standup/internal/app/response/conversationsList"
	"github.com/vehsamrak/slack-standup/internal/app/response/conversationsMembers"
	"github.com/vehsamrak/slack-standup/internal/app/response/conversationsOpen"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	Config *config.Config
}

func (Client) Create(config *config.Config) *Client {
	return &Client{Config: config}
}

func (slack *Client) botInfo() {
	slack.callApi("bots.info", nil)
}

func (slack *Client) ChannelUsersList(channelName string) *conversationsMembers.Users {
	response := slack.callApi("conversations.members", url.Values{"channel": {channelName}})
	users := &conversationsMembers.Users{}

	err := json.Unmarshal(response, &users)
	if err != nil {
		panic(err)
	}

	return users
}

func (slack *Client) channelsList() *conversationsList.ChannelsList {
	response := slack.callApi("conversations.list", url.Values{"exclude_archived": {"true"}})
	channels := &conversationsList.ChannelsList{}

	err := json.Unmarshal(response, &channels)
	if err != nil {
		panic(err)
	}

	return channels
}

func (slack *Client) FindChannelByName(channelName string) *conversationsList.Channel {
	channels := slack.channelsList()
	for _, channel := range channels.Channels {
		if channel.Name == channelName {
			return &channel
		}
	}

	return nil
}

func (slack *Client) SendMessageToChannelByName(channelName string, message string) *chatPostMessage.Thread {
	channel := slack.FindChannelByName(channelName)
	if channel == nil {
		log.Info("Channel was not found. Can not send message", channelName)
		return nil
	}

	response := slack.callApi(
		"chat.postMessage",
		url.Values{"channel": {channel.Id}, "text": {message}},
	)

	thread := &chatPostMessage.Thread{}
	err := json.Unmarshal(response, &thread)
	if err != nil {
		panic(err)
	}

	return thread
}

func (slack *Client) SendMessageToChannel(channelId string, message string) {
	slack.callApi(
		"chat.postMessage",
		url.Values{"channel": {channelId}, "text": {message}},
	)
}

func (slack *Client) SendReplyToChannel(channel string, message string, thread string) {
	slack.callApi(
		"chat.postMessage",
		url.Values{"channel": {channel}, "text": {message}, "thread_ts": {thread}},
	)
}

func (slack *Client) OpenChatWithUser(userName string) *conversationsOpen.Channel {
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

func (slack *Client) callApi(method string, parameters url.Values) []byte {
	if parameters == nil {
		parameters = url.Values{}
	}

	log.Debugf("Calling API: %s [%s]\n", method, parameters.Encode())

	parameters.Add("token", slack.Config.Token)
	parameters.Add("pretty", "1")

	resp, err := http.Get(fmt.Sprintf("%s/%s?%s", slack.Config.ApiUrl, method, parameters.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Infof("%s", body)

	return body
}
