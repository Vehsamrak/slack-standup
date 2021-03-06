package meeting

import "github.com/vehsamrak/slack-standup/internal/app/response/chatPostMessage"

type Meeting struct {
	Participants map[string]*Questions
	Thread       *chatPostMessage.Thread
	channelName  string
}

func (meeting Meeting) Create(channelName string, thread *chatPostMessage.Thread) *Meeting {
	return &Meeting{Thread: thread, Participants: make(map[string]*Questions), channelName: channelName}
}

func (meeting Meeting) Greetings() string {
	return "Начинается утренний стэндап! @here"
}

func (meeting *Meeting) ChannelName() string {
	return meeting.channelName
}
