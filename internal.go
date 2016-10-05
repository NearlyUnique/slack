package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/websocket"
)

var (
	// counter for sending unique id for PostMessage
	counter  uint64
	stopPing = make(chan bool)
)

// dial connects to the real time slack api
func (client *Client) dial() error {
	err := client.slackStart()
	if err != nil {
		return err
	}
	client.ws, err = websocket.Dial(client.wsURL, "", "https://api.slack.com/")
	if err != nil {
		return err
	}
	go client.pingLoop()
	return nil
}

// listenLoop for messages, messages are pushed to the RxMessage channel, errors to client.Error
func (client *Client) listenLoop(handlers ...HandlerFunc) {
	for {
		if rx, err := client.getMessage(); err == nil {
			if rx == nil {
				continue
			}
			log.Printf("msg: %s", rx.Type)

			for _, handler := range handlers {
				handler(BotInfo{
					ID:     client.Self.ID,
					Poster: client,
				}, Message{
					User:    rx.User,
					Channel: rx.Channel,
					Text:    rx.Text,
				})
			}
		} else {
			log.Printf("Something broke %v\n", err)
			return
		}
	}
}

// getMessage receives any pending message, in the case of an error will
// stop pinging server and close. Listen(...) will re-Dial()
func (client *Client) getMessage() (*RxMessage, error) {
	m := map[string]interface{}{}
	err := websocket.JSON.Receive(client.ws, &m)
	if err != nil {
		// call close instead?
		stopPing <- true
		return nil, err
	}
	var messageType interface{}
	ok := false
	if messageType, ok = m["type"]; !ok || messageType.(string) != "message" {
		return nil, nil
	}
	rxM := rxMessage{}
	if err = mapstructure.Decode(m, &rxM); err != nil {
		return nil, err
	}
	//TODO: not sure this mapping helps
	msg := RxMessage{
		Type: messageType.(string),
		Text: rxM.Text,
		Time: rxM.Time,
	}

	client.addDirectMessageChannel(rxM.Channel)

	return &RxMessage{
		User:    client.users[rxM.User],
		Channel: client.channels[rxM.Channel],
		Text:    msg.Text,
		Type:    msg.Type,
		Time:    msg.Time,
	}, nil
}

func (client *Client) addDirectMessageChannel(ID string) {
	if _, found := client.channels[ID]; !found {
		//TODO: should we go off and ask slack for more info about this DM channel
		client.channels[ID] = Channel{
			ID: ID,
		}
	}
}
func safeValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// post for any kind of message
func (client *Client) post(m txMessage) error {
	m.ID = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(client.ws, m)
}

// slackStart does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func (client *Client) slackStart() error {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", client.token)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("API request failed with code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var r responseRtmStart

	//TODO: diagnostics only, remove
	go ioutil.WriteFile("rtm.start.json", body, 0644)

	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}

	if !r.Ok {
		return fmt.Errorf("Slack error: %s", r.Error)
	}
	log.Printf("users:%d, channels:%d", len(r.Users), len(r.Channels))
	client.mapChannelsAndUsers(r)

	return nil
}
func (client *Client) mapChannelsAndUsers(r responseRtmStart) {
	client.wsURL = r.URL
	client.users = make(map[string]User)
	client.channels = make(map[string]Channel)

	for i := 0; i < len(r.Users); i++ {
		client.users[r.Users[i].ID] = r.Users[i]
	}
	for i := 0; i < len(r.Channels); i++ {
		client.channels[r.Channels[i].ID] = r.Channels[i]
	}
	client.Self = client.users[r.Self.ID]
}

// pingLoop, internal pinger, needed to keep connection alive
func (client *Client) pingLoop() {
	ticker := time.NewTicker(5 * time.Second)
Loop:
	for {
		select {
		case <-ticker.C:
			client.post(txMessage{Type: "ping"})
		case <-stopPing:
			break Loop
		}
	}
	ticker.Stop()
}
