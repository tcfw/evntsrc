package evntsrc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tcfw/evntsrc/internal/websocks"

	"github.com/gorilla/websocket"
)

//openConn dials the evntsrc realtime endpoint and attempts to authenticate
func (api *APIClient) openConn() error {
	var headers *http.Header
	if api.auth != "" {
		headers = &http.Header{}
		headers.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(api.auth))))
	}

	if api.Debug {
		letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		n := 25

		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}

		headers.Add("x-trace", string(b))
		fmt.Printf("X-Trace: %s\n", string(b))
	}

	url := api.formatURL("realtime", strconv.Itoa(int(api.stream)))

	conn, resp, err := websocket.DefaultDialer.Dial(url, *headers)
	if err != nil {
		if api.Debug {
			fmt.Printf("%v\n", resp)
		}
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		if api.Debug {
			fmt.Printf("CLOSED %v %v\n", code, text)
		}
		message := websocket.FormatCloseMessage(code, "")
		conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(writeWait))
		api.close <- true
		return nil
	})

	api.socket = conn

	go api.readPump()
	go api.writePump()

	return nil
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 2 * time.Minute

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 1024 * 1024
)

//readPump listens to messages from the socket and handles the message
//accordingly
func (api *APIClient) readPump() {
	api.socket.SetReadLimit(maxMessageSize)
	api.socket.SetPongHandler(func(string) error {
		return nil
	})

	defer func() {
		api.socket.Close()
		api.socket = nil
		api.close <- true
	}()

	go api.distributeReadPipe()

	for {
		_, message, err := api.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				//TODO(tcfw) handle this better
				if api.Debug {
					fmt.Printf("ERR Reading: %v\n", err)
				}
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))

		if api.Debug {
			fmt.Printf("#!RCV[%d]: %s\n", len(message), string(message))
		}

		if strings.HasPrefix(string(message), `{"acktype"`) {
			command := &websocks.AckCommand{}
			if err = json.Unmarshal(message, command); err != nil {
				errMsg := &websocks.AckCommand{Acktype: "Error", Error: "Failed to parse command"}
				jsonBytes, _ := json.Marshal(errMsg)
				api.pubError(errors.New(string(jsonBytes)))
			}

			api.AcksCh <- command
			continue
		}
		if strings.Contains(string(message), `connectionID`) {
			info := &websocks.ConnectionInfo{}
			if err = json.Unmarshal(message, info); err != nil {
				errMsg := &websocks.AckCommand{Acktype: "Error", Error: "Failed to parse connection info"}
				jsonBytes, _ := json.Marshal(errMsg)
				api.pubError(errors.New(string(jsonBytes)))
				return
			}
			api.connectionID = info.ConnectionID
			continue
		}

		api.ReadPipe <- message
	}
}

//writePump listesn to each type of command and forwards the command to the
//relevant do* function to write to the socket
func (api *APIClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case msg := <-api.writePipe:
			if err := api.doPublish(msg); err != nil {
				api.pubError(err)
			}
			break
		case msg := <-api.replayPipe:
			if err := api.doReplay(msg); err != nil {
				api.pubError(err)
			}
			break
		case msg := <-api.subPipe:
			if err := api.doSendSubscribe(msg); err != nil {
				api.pubError(err)
			}
			break
		case <-api.close:
			return
		case <-ticker.C:
			if err := api.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				api.pubError(fmt.Errorf("Failed to send ping"))
			}
		}
	}
}

//distributeReadPipe listens for new events from the inbound socket read pipe
//(see readPump()) and distributes them to each of the related subscribed
//chans or callbacks
func (api *APIClient) distributeReadPipe() {
	for {
		msg, ok := <-api.ReadPipe
		if !ok {
			//Read pipe is closed
			break
		}

		evnt := &Event{}
		if err := json.Unmarshal(msg, evnt); err != nil {
			api.pubError(err)
			continue
		}

		if evnt.Stream == 0 && evnt.Subject == "puback" {
			go func() {
				api.AcksCh <- &websocks.AckCommand{Acktype: "OK", Ref: evnt.ID}
			}()
			continue
		}

		if source, ok := evnt.Metadata["relative_seq"]; api.options.IgnoreSelf && ok && strings.HasPrefix(source, fmt.Sprintf("%s-", api.connectionID)) {
			continue
		}

		_, encrypted := evnt.Metadata["e"]
		//Warn if event is encrypted by no crypto
		if encrypted && api.options.Crypto == nil {
			api.pubError(fmt.Errorf("Encrypted message received but crypto not set up: event id %s", evnt.ID))
		}

		if api.options.Crypto != nil && encrypted {
			if err := api.options.Crypto.Verify(evnt.Data, evnt.Metadata); err != nil {
				api.pubError(err)
				continue
			}

			data, err := api.options.Crypto.Decrypt(evnt.Data, evnt.Metadata)
			if err != nil {
				api.pubError(err)
				continue
			}
			evnt.Data = data
		}

		go func(evnt *Event) {
			c := 0
			for subject := range api.subscriptions {
				if subject == evnt.Subject {
					c++
					for _, subset := range api.subscriptions[subject] {
						subbersEvent := evnt
						switch subset.subType {
						case funcSubType:
							go subset.f(subbersEvent)
							break
						case chanSubType:
							go func() {
								//Timeout chan send after 30 seconds
								select {
								case subset.ch <- subbersEvent:
								case <-time.After(30 * time.Second):
								}
							}()
							break
						}
					}
				}
			}
			if c == 0 {
				api.pubError(fmt.Errorf("Received event wasn't listen for"))
			}
		}(evnt)
	}
}

//doPublish writes the publish command to the socket
func (api *APIClient) doPublish(data *websocks.PublishEventCommand) error {
	if api.socket == nil {
		if err := api.openConn(); err != nil {
			return fmt.Errorf("Failed to publish: %s", err.Error())
		}
	}

	return api.socket.WriteJSON(data)
}

//doReplay writes the replay command to the socket
func (api *APIClient) doReplay(data *websocks.ReplayCommand) error {
	if api.socket == nil {
		if err := api.openConn(); err != nil {
			return fmt.Errorf("Failed to publish: %s", err.Error())
		}
	}

	if api.Debug {
		jsonBytes, _ := json.Marshal(data)
		fmt.Printf("#!RPL[%d]: %s\n", len(jsonBytes), string(jsonBytes))
	}

	return api.socket.WriteJSON(data)
}

//doSendSubscribe writes the sub command to the socket
func (api *APIClient) doSendSubscribe(data *websocks.SubscribeCommand) error {
	if api.socket == nil {
		if err := api.openConn(); err != nil {
			return fmt.Errorf("Failed to publish: %s", err.Error())
		}
	}

	return api.socket.WriteJSON(data)
}

//Publish publishes data given an subject and event type
func (api *APIClient) Publish(subject string, data []byte, eventType string) error {
	evnt := &Event{
		Subject: subject,
		Data:    data,
		Type:    eventType,
	}

	return api.PublishEvent(evnt)
}

//PublishEvent pubs an event to the stream
func (api *APIClient) PublishEvent(ev *Event) error {
	if api.socket == nil {
		if err := api.openConn(); err != nil {
			return fmt.Errorf("Failed to publish: %s", err.Error())
		}
	}

	if ev.Metadata == nil {
		ev.Metadata = map[string]string{}
	}

	if ev.ContentType == "" {
		ev.ContentType = "application/octet-stream"
	}

	if api.options.Crypto != nil {
		encBytes, encMd, err := api.options.Crypto.Encrypt([]byte(ev.Data))
		if err != nil {
			return err
		}

		ev.Data = encBytes
		for k, v := range encMd {
			ev.Metadata[k] = v
		}

		//Encrypted flag
		ev.Metadata["e"] = "1"
	}

	pubCmd := &websocks.PublishEventCommand{
		SubscribeCommand: &websocks.SubscribeCommand{InboundCommand: &websocks.InboundCommand{Ref: uuid.New().String(), Command: "epub"}, Subject: ev.Subject},
		Event:            ev,
	}

	api.writePipe <- pubCmd
	if api.WaitForPublishConfirmation {
		ok, err := api.waitForResponse(pubCmd.Ref)
		if err != nil {
			fmt.Printf("%v", err)
		}
		if !ok {
			return fmt.Errorf("Failed to publish: %s", err.Error())
		}
	}

	return nil
}

//doSubscribe sends the subscribe event to the write pump and waits
//for ack from websocks
func (api *APIClient) doSubscribe(subject string) error {
	subMsg := &websocks.SubscribeCommand{
		InboundCommand: &websocks.InboundCommand{Ref: uuid.New().String(), Command: "sub"},
		Subject:        subject,
	}

	api.subPipe <- subMsg

	if _, err := api.waitForResponse(subMsg.Ref); err != nil {
		return fmt.Errorf("Failed to subscribe: %s", err.Error())
	}

	return nil
}

//SubscribeChan to a subject inside the stream via a channel
func (api *APIClient) SubscribeChan(subject string) (<-chan *Event, error) {
	if err := api.doSubscribe(subject); err != nil {
		return nil, err
	}

	ch := make(chan *Event, 256)

	if _, ok := api.subscriptions[subject]; !ok {
		api.subscriptions[subject] = []*subscription{}
	}

	api.subscriptions[subject] = append(api.subscriptions[subject], &subscription{subType: chanSubType, ch: ch})

	return ch, nil
}

//SubscribeFunc subscribes to a subject and call callback func
func (api *APIClient) SubscribeFunc(subject string, callback func(*Event)) error {
	if err := api.doSubscribe(subject); err != nil {
		return err
	}

	if _, ok := api.subscriptions[subject]; !ok {
		api.subscriptions[subject] = []*subscription{}
	}

	api.subscriptions[subject] = append(api.subscriptions[subject], &subscription{subType: funcSubType, f: callback})

	return nil
}

//Unsubscribe from a subject
func (api *APIClient) Unsubscribe(subject string) error {
	if subset, ok := api.subscriptions[subject]; ok {
		for _, subscription := range subset {
			switch subscription.subType {
			case chanSubType:
				close(subscription.ch)
				break
			case funcSubType:
				//nothing to do, just stop sending events
				break
			}
		}

		delete(api.subscriptions, subject)
		return nil
	}

	return fmt.Errorf("No subscription for subject '%s'", subject)
}
