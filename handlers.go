package slack

type (
	// HandlerFunc for messages
	HandlerFunc func(BotInfo, Message) bool
)
