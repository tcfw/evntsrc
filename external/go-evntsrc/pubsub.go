package evntsrc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tcfw/evntsrc/pkg/websocks"

	"github.com/gorilla/websocket"
)

func (api *APIClient) openConn() error {
	var headers *http.Header
	if api.auth != "" {
		headers = &http.Header{}
		headers.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(api.auth))))
	}

	url := api.formatURL("realtime", strconv.Itoa(int(api.stream)))

	conn, resp, err := websocket.DefaultDialer.Dial(url, *headers)
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("CLOSED %v %v\n", code, text)
		api.close <- true
		return nil
	})
	if err != nil {
		fmt.Printf("%v\n", resp)
		return err
	}

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
)

func (api *APIClient) readPump() {
	api.socket.SetReadLimit(4096)
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
				fmt.Printf("ERR Reading: %v\n", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))

		if strings.HasPrefix(string(message), `{"acktype"`) {
			command := &websocks.AckCommand{}
			if err = json.Unmarshal(message, command); err != nil {
				errMsg := &websocks.AckCommand{Acktype: "Error", Error: "Failed to parse command"}
				jsonBytes, _ := json.Marshal(errMsg)
				api.Errors <- errors.New(string(jsonBytes))
			}

			api.AcksCh <- command
			continue
		}
		if strings.Contains(string(message), `connectionID`) {
			info := &websocks.ConnectionInfo{}
			if err = json.Unmarshal(message, info); err != nil {
				errMsg := &websocks.AckCommand{Acktype: "Error", Error: "Failed to parse connection info"}
				jsonBytes, _ := json.Marshal(errMsg)
				api.Errors <- errors.New(string(jsonBytes))
				return
			}
			api.connectionID = info.ConnectionID
			continue
		}

		api.ReadPipe <- message
	}
}

func (api *APIClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case msg := <-api.writePipe:
			if err := api.doPublish(msg); err != nil {
				api.Errors <- err
			}
		case <-api.close:
			return
		case <-ticker.C:
			if err := api.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				api.Errors <- fmt.Errorf("Failed to send ping")
			}
		}
	}
}

func (api *APIClient) distributeReadPipe() {
	for {
		msg, ok := <-api.ReadPipe
		if !ok {
			//Read pipe is closed
			break
		}

		evnt := &Event{}
		if err := json.Unmarshal(msg, evnt); err != nil {
			api.Errors <- err
			continue
		}

		if source, ok := evnt.Metadata["relative_seq"]; api.IgnoreSelf && ok && strings.HasPrefix(source, fmt.Sprintf("%s-", api.connectionID)) {
			continue
		}

		go func(evnt *Event) {
			for subject := range api.subscriptions {
				if subject == evnt.Subject {
					for _, subset := range api.subscriptions[subject] {
						switch subset.subType {
						case funcSubType:
							subset.f(evnt)
							break
						case chanSubType:
							subset.ch <- evnt
							break
						}
					}
				}
			}
		}(evnt)
	}
}

func (api *APIClient) doPublish(data *websocks.PublishCommand) error {
	if api.socket == nil {
		if err := api.openConn(); err != nil {
			return fmt.Errorf("Failed to publish: %s", err.Error())
		}
	}

	return api.socket.WriteJSON(data)
}

//Publish publishes an event through evntsrc
func (api *APIClient) Publish(subject string, data []byte, eventType string) error {
	if api.socket == nil {
		if err := api.openConn(); err != nil {
			return fmt.Errorf("Failed to publish: %s", err.Error())
		}
	}

	pubMsg := &websocks.PublishCommand{
		SubscribeCommand: &websocks.SubscribeCommand{InboundCommand: &websocks.InboundCommand{Ref: uuid.New().String(), Command: "pub"}, Subject: subject},
		Data:             base64.StdEncoding.EncodeToString(data),
		ContentType:      "application/json",
		Type:             eventType,
		TypeVersion:      api.AppVer,
	}

	api.writePipe <- pubMsg
	ok, err := api.waitForResponse(pubMsg.Ref)
	if err != nil {
		fmt.Printf("%v", err)
	}
	if !ok {
		return fmt.Errorf("Failed to publish: %s", err.Error())
	}

	return err
}

func (api *APIClient) doSubscribe(subject string) error {
	if api.socket == nil {
		if err := api.openConn(); err != nil {
			return fmt.Errorf("Failed to subscribe: %s", err.Error())
		}
	}

	subMsg := &websocks.SubscribeCommand{
		InboundCommand: &websocks.InboundCommand{Ref: uuid.New().String(), Command: "sub"},
		Subject:        subject,
	}

	if err := api.socket.WriteJSON(subMsg); err != nil {
		return err
	}

	_, err := api.waitForResponse(subMsg.Ref)
	if err != nil {
		return fmt.Errorf("Failed to subscribe: %s", err.Error())
	}

	return nil
}

//Subscribe to a subject inside the stream via a channel
func (api *APIClient) Subscribe(subject string) (chan *Event, error) {
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
				//nothing to do
				break
			}
		}

		delete(api.subscriptions, subject)
		return nil
	}

	return fmt.Errorf("No subscription for subject '%s'", subject)
}
