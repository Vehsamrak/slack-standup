package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiUrl = "https://slack.com/api/"
const apiMethodConversations = "conversations.list"
const apiMethodChatPostMessage = "chat.postMessage"
const apiMethod = "bots.info"

type Slack struct {
	Config *Config
}

func (Slack) Create(config *Config) *Slack {
	return &Slack{Config: config}
}

func (slack *Slack) sendMessageToChannel(channel string, message string) {
	slack.callApi("chat.postMessage", url.Values{"channel": {channel}, "text": {message}, "thread_ts": {"1604271791.000300"}})
}

func (slack *Slack) callApi(method string, parameters url.Values) {
	parameters.Add("token", slack.Config.Token)
	parameters.Add("pretty", "1")

	resp, err := http.Get(fmt.Sprintf("%s/%s?%s", apiUrl, method, parameters.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s", body)
}
