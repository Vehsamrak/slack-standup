package meeting

import "fmt"

type Questions struct {
	Previous string
	Today    string
	Block    string
}

func (questions Questions) Create() *Questions {
	return &Questions{}
}

func (questions Questions) QuestionPrevious() string {
	return "*Удалось выполнить предыдущий план?*"
}

func (questions Questions) QuestionToday() string {
	return "*Что планируешь сделать сегодня?*"
}

func (questions Questions) QuestionBlock() string {
	return "*Кто и чем может тебе в этом помочь?*"
}

func (questions *Questions) Result() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s",
		questions.QuestionPrevious(),
		questions.Previous,
		questions.QuestionToday(),
		questions.Today,
		questions.QuestionBlock(),
		questions.Block,
	)
}
