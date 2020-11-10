package meeting

import "github.com/vehsamrak/slack-standup/internal/app/response/chatPostMessage"

type Meeting struct {
	Participants map[string]*Questions
	Thread       *chatPostMessage.Thread
}

func (meeting Meeting) Create(thread *chatPostMessage.Thread) *Meeting {
	return &Meeting{Thread: thread, Participants: make(map[string]*Questions)}
}

func (meeting *Meeting) QuestionPrevious() string {
	return "*Удалось выполнить предыдущий план?*"
}

func (meeting *Meeting) QuestionToday() string {
	return "*Что планируешь сделать сегодня?*"
}

func (meeting *Meeting) QuestionBlock() string {
	return "*Кто и чем может тебе в этом помочь?*"
}
