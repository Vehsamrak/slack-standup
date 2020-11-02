package meeting

type Questions struct {
	Previous string
	Today    string
	Block    string
}

func (questions Questions) Create() *Questions {
	return &Questions{}
}
