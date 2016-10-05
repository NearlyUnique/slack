package slack

import "golang.org/x/net/websocket"

type (
	// Client is a client to receive slack real time events on
	Client struct {
		// Self is the (bot-)user that has connected
		Self User
		// the error channel, something went wrong async
		Error chan error
		// link to talk web sockets on
		wsURL string
		// the application token
		token string
		// the web socket for comms
		ws *websocket.Conn
		// users, all known users, keyed by id
		users map[string]User
		// channels, all known channels, keysed by id
		channels map[string]Channel
	}
	// BotInfo for handlers that want to know who they are and talk back to slack
	BotInfo struct {
		// BotID for this bot
		ID string
		// Poster to slack
		Poster Poster
	}

	//------------------ Internal ----------------------

	// Poster to a slack channel
	Poster interface {
		PostMessage(channel, text string) error
	}
	//Message received by the bot
	Message struct {
		User    User
		Channel Channel
		Text    string
	}
	// txMessage - These are the messages written into the websocket.
	txMessage struct {
		// ID is internally set
		ID uint64 `json:"id"`
		// Type is internally set to "message"
		Type string `json:"type"`
		// Channel to send the message
		Channel string `json:"channel"`
		// Text content of the message
		Text string `json:"text"`
	}
	// rxMessage - a receive type message
	rxMessage struct {
		Type    string `json:"type"`
		Channel string `json:"channel"`
		User    string `json:"user"`
		Text    string `json:"text"`
		Time    string `json:"ts"`
	}
	// connection response
	responseRtmStart struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
		URL   string `json:"url"`
		Self  struct {
			ID    string                 `json:"id"`
			Name  string                 `json:"name"`
			Prefs map[string]interface{} `json:"prefs"`
		} `json:"self"`
		Channels []Channel `json:"channels"`
		Users    []User    `json:"users"`
	}
)
