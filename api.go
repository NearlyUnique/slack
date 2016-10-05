package slack

import (
	"log"
	"os"
)

// NewClient Creates an instance of a slack Real Time Message client
func NewClient(token string) *Client {
	return &Client{
		token: token,
	}
}

// Listen and call handlers as appropriate for all inbound messages
func (client *Client) Listen(handlers ...HandlerFunc) {
	for {
		log.Println("Connecting...")
		err := client.dial()
		if err != nil {
			log.Printf("Unable to connect to slack: %v", err)
			os.Exit(1)
		}

		log.Printf("Connected as %v", client.Self)
		client.listenLoop(handlers...)
		log.Printf("Connection Lost...")
	}
}

// PostMessage for sending text messages
func (client *Client) PostMessage(channel, text string) error {
	m := txMessage{
		Type:    "message",
		Text:    text,
		Channel: channel,
	}
	return client.post(m)
}
