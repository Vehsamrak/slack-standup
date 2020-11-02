package messageIm

type Message struct {
	Text   string `json:"text"`
	UserId string `json:"user"`
	BotId  string `json:"bot_id"`
	Type   string `json:"type"`
}
