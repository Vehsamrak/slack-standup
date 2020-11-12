package conversationsList

type ChannelsList struct {
	Channels         []Channel        `json:"channels"`
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
}
