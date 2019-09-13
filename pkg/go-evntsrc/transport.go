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

	"github.com/gorilla/websocket"
	"github.com/tcfw/evntsrc/internal/websocks"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 2 * time.Minute

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 1024 * 1024
)

//Reconnect tests the websocket for acitivty and if fails then restart the connection
func (api *APIClient) Reconnect() error {
	//
	if err := api.socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
		return err
	}

	pong := make(chan struct{})
	go func() {
		api.pongCond.L.Lock()
		api.pongCond.Wait()
		api.pongCond.L.Unlock()
		close(pong)
	}()

	select {
	case <-pong:
		//Received a pong, assuming all is well
		return nil
	case <-time.After(6 * time.Second):
		//Failed to receive a pong, assuming connection is dead
		api.rebuildConnection()
	}
	return nil
}

func (api *APIClient) rebuildConnection() {
	api.socket.Close()
	api.openConn()

	for subject := range api.subscriptions {
		api.doSubscribe(subject)
		api.Replay(subject, websocks.ReplayRange{StartTime: &api.lastEvent}, true)
	}
}

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

//readPump listens to messages from the socket and handles the message
//accordingly
func (api *APIClient) readPump() {
	api.socket.SetReadLimit(maxMessageSize)
	api.socket.SetPongHandler(func(string) error {
		api.pongCond.Broadcast()
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
				if api.Debug {
					fmt.Printf("ERR Reading: %v\n", err)
				}
				api.Reconnect()
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

//writePump listens to each type of command and forwards the command to the
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
