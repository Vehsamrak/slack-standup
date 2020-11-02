package messageIm

type Message struct {
	Text string `json:"text"`
	User string `json:"user"`
	Type string `json:"type"`
}
