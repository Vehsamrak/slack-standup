package messageIm

type Message struct {
	Text   string `json:"text"`
	UserId string `json:"user"`
	Type   string `json:"type"`
}
