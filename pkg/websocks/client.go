package websocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	nats "github.com/nats-io/go-nats"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/tcfw/evntsrc/pkg/event"
	streamauth "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	"gopkg.in/mgo.v2/bson"
)

const (
	maxMessageSize = 2048

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

//Client maintains WS info
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//Subscription state
	subscriptions map[string]chan bool

	//Auth mappings
	auth    *AuthCommand
	authKey *streamauth.StreamKey

	//Connection state
	connectionID string
	seq          map[string]int64
	seqLock      sync.Mutex
	closed       bool
}

//NewClient constructs a new websocket evntsrc client
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:          conn,
		send:          make(chan []byte, 256),
		subscriptions: map[string]chan bool{},
		connectionID:  bson.NewObjectId().Hex(),
		seq:           map[string]int64{},
		closed:        false,
	}
}

func (c *Client) close() {
	c.conn.Close()

	if !c.closed {
		go c.broadcastDisconnect()
		m := metrics.GetOrRegisterCounter("wsConnections", nil)
		m.Dec(1)
	}
	c.closed = true

	for channel := range c.subscriptions {
		close(c.subscriptions[channel])
		delete(c.subscriptions, channel)
	}
}

func (c *Client) sendStruct(msg interface{}) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.send <- msgBytes
	return nil
}

func (c *Client) processCommand(command *InboundCommand, message []byte) {
	switch command.Command {
	case commandSubscribe:
		subcommand := &SubscribeCommand{}
		if err := json.Unmarshal(message, subcommand); err != nil {
			c.sendStruct(&AckCommand{Acktype: "Error", Error: "Failed to parse command"})
			return
		}
		c.Subscribe(subcommand.Subject)

	case commandUnsubscribe:
		subcommand := &UnsubscribeCommand{}
		if err := json.Unmarshal(message, subcommand); err != nil {
			c.sendStruct(&AckCommand{Acktype: "Error", Error: "Failed to parse command"})
			return
		}
		c.Unsubscribe(subcommand.Subject)

	case commandPublish:
		subcommand := &PublishCommand{}
		if err := json.Unmarshal(message, subcommand); err != nil {
			c.sendStruct(&AckCommand{Acktype: "Error", Error: "Failed to parse command"})
			return
		}

		rEvent := &event.Event{}
		rEvent.SetID()
		rEvent.Stream = c.auth.Stream
		rEvent.Subject = subcommand.Subject
		rEvent.CEVersion = "0.1"
		rEvent.Type = subcommand.Type
		rEvent.TypeVersion = subcommand.TypeVersion
		rEvent.ContentType = subcommand.ContentType
		rEvent.Data = []byte(subcommand.Data)
		rEvent.Time = event.ZeroableTime{Time: time.Now()}
		rEvent.Metadata = map[string]string{"source_ip": c.conn.RemoteAddr().String()}
		rEvent.Source = subcommand.Source
		if rEvent.Source == "" {
			rEvent.Source = "ws"
		}

		channel := fmt.Sprintf("_USER.%d.%s", c.auth.Stream, subcommand.Subject)

		c.seqLock.Lock()
		c.seq[channel]++
		rEvent.Metadata["relative_seq"] = fmt.Sprintf("%s-%d", c.connectionID, c.seq[channel])
		c.seqLock.Unlock()

		eventJSONBytes, _ := json.Marshal(rEvent)
		natsConn.Publish(channel, eventJSONBytes)

	case commandAuth:
		subcommand := &AuthCommand{}
		if err := json.Unmarshal(message, subcommand); err != nil {
			c.sendStruct(&AckCommand{Acktype: "Error", Error: "Failed to parse command"})
			return
		}

		if err := c.validateAuth(subcommand); err != nil {
			c.sendStruct(&AckCommand{
				Acktype: "Failed",
				Error:   err.Error(),
			})
			return
		}

		c.auth = subcommand
		c.sendStruct(&AckCommand{
			Acktype: "OK",
		})
		go c.broadcastConnect()

	case commandReplay:
		subcommand := &ReplayCommand{}
		if err := json.Unmarshal(message, subcommand); err != nil {
			c.sendStruct(&AckCommand{Acktype: "Error", Error: "Failed to parse command"})
			return
		}
		subcommand.Stream = c.auth.Stream

		//Request replay to storer svc
		repubBytes, _ := json.Marshal(subcommand)
		msg, err := natsConn.Request("replay.broadcast", repubBytes, time.Second*10)
		ack := &AckCommand{
			Acktype: "err",
		}
		if err != nil {
			ack.Error = err.Error()
		} else {
			ack.Error = string(msg.Data)
		}
		c.sendStruct(ack)
	}
}

//Subscribe opens a NATS subscription for the websocket
func (c *Client) Subscribe(channel string) {
	if _, ok := c.subscriptions[channel]; ok {
		c.sendStruct(&AckCommand{
			Acktype: "err",
			Error:   fmt.Sprintf("already subscribed to %s", channel),
			Channel: channel,
		})
		return
	}

	c.subscriptions[channel] = make(chan bool)

	//Subscribe to NATS user channel and forward events to websocket
	go func(c *Client, channel string, unsub chan bool) {
		ch := make(chan *nats.Msg, 64)
		sub, err := natsConn.ChanSubscribe(fmt.Sprintf("_USER.%d.%s", c.auth.Stream, channel), ch)
		ack := &AckCommand{
			Acktype: "sub",
			Channel: channel,
		}
		if err != nil {
			ack.Acktype = "err"
			ack.Error = err.Error()
		}

		c.sendStruct(ack)
		go c.broadcastSub(channel)

		if err == nil {
			for {
				select {
				case msg := <-ch:
					c.send <- msg.Data
				case <-unsub:
					sub.Unsubscribe()
					sub.Drain()
					return
				}
			}
		}
	}(c, channel, c.subscriptions[channel])
}

//Unsubscribe closes a subscription to NATS
func (c *Client) Unsubscribe(channel string) {
	if _, ok := c.subscriptions[channel]; !ok {
		c.sendStruct(&AckCommand{
			Acktype: "err",
			Error:   fmt.Sprintf("not subscribed to %s", channel),
			Channel: channel,
		})
		return
	}
	close(c.subscriptions[channel])
	delete(c.subscriptions, channel)

	c.sendStruct(&AckCommand{
		Acktype: "unsub",
		Channel: channel,
	})
	go c.broadcastUnsub(channel)
}

type readPRPC struct {
	command *InboundCommand
	message []byte
}

func (c *Client) readPump() {
	defer func() {
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	timeoutEnabled := true

	readCh := make(chan *readPRPC, 50)
	closed := make(chan bool, 1)
	initTimeout := time.After(15 * time.Second)

	go func() {
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				closed <- true
				return
			}
			timeoutEnabled = false
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

			command := &InboundCommand{}
			if err = json.Unmarshal(message, command); err != nil {
				log.Printf("Failed to parse command: %s\n", err.Error())
				c.sendStruct(&AckCommand{Acktype: "Error", Error: "Failed to parse command"})
				continue
			}

			if command.Command != commandAuth && c.auth == nil {
				c.sendStruct(&AckCommand{
					Acktype: "err",
					Error:   "No auth sent yet",
				})
				continue
			}

			readCh <- &readPRPC{command, message}
		}
	}()

	for {
		select {
		case <-closed:
			return
		case rpc := <-readCh:
			c.processCommand(rpc.command, rpc.message)
		case <-initTimeout:
			if timeoutEnabled {
				c.close()
				return
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
