package websocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"

	"github.com/gorilla/websocket"
	nats "github.com/nats-io/go-nats"
	metrics "github.com/rcrowley/go-metrics"
	streamauth "github.com/tcfw/evntsrc/internal/streamauth/protos"
	"gopkg.in/mgo.v2/bson"
)

const (
	maxMessageSize = 4096

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
	closeConnSub  chan bool

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
		closeConnSub:  make(chan bool, 1),
	}
}

func (c *Client) close() {
	c.conn.Close()

	if !c.closed {
		socketGauge.WithLabelValues(fmt.Sprintf("%d", c.auth.Stream)).Dec()
		go c.broadcastDisconnect()
		m := metrics.GetOrRegisterCounter("wsConnections", nil)
		m.Dec(1)
	}
	c.closed = true

	for channel := range c.subscriptions {
		close(c.subscriptions[channel])
		delete(c.subscriptions, channel)
	}

	c.closeConnSub <- true
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
		c.doSubscribe(command, message)
		break
	case commandUnsubscribe:
		c.doUnsubscribe(command, message)
		break
	case commandPublish:
		c.doPublish(command, message)
		break
	case commandAuth:
		c.doAuth(command, message)
		break
	case commandReplay:
		c.doReplay(command, message)
		break
	}
}

//ConnSub subscribes to a connection specific channel
func (c *Client) ConnSub() error {
	ch := make(chan *nats.Msg, 1024)
	sub, err := natsConn.ChanSubscribe(fmt.Sprintf("_CONN.%s", c.connectionID), ch)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-ch:
			c.send <- msg.Data
		case <-c.closeConnSub:
			sub.Unsubscribe()
			sub.Drain()
			return nil
		}
	}
}

//Subscribe opens a NATS subscription for the websocket
func (c *Client) Subscribe(channel string, cmd *InboundCommand) {
	if _, ok := c.subscriptions[channel]; ok {
		c.sendStruct(&AckCommand{
			Ref:     cmd.Ref,
			Acktype: "err",
			Error:   fmt.Sprintf("already subscribed to %s", channel),
			Channel: channel,
		})
		return
	}

	c.subscriptions[channel] = make(chan bool)

	//Subscribe to NATS user channel and forward events to websocket
	go func(c *Client, channel string, unsub chan bool) {
		ch := make(chan *nats.Msg, 1024)
		sub, err := natsConn.ChanSubscribe(fmt.Sprintf("_USER.%d.%s", c.auth.Stream, channel), ch)
		ack := &AckCommand{
			Ref:     cmd.Ref,
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
					ev := &pbEvent.Event{}
					proto.Unmarshal(msg.Data, ev)
					jsonEv, _ := json.Marshal(ev)

					c.send <- jsonEv

					byteSubscribeCounter.WithLabelValues(fmt.Sprintf("%d", c.auth.Stream)).Add(float64(len(msg.Data)))

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
func (c *Client) Unsubscribe(channel string, cmd *InboundCommand) {
	if _, ok := c.subscriptions[channel]; !ok {
		c.sendStruct(&AckCommand{
			Ref:     cmd.Ref,
			Acktype: "err",
			Error:   fmt.Sprintf("not subscribed to %s", channel),
			Channel: channel,
		})
		return
	}
	close(c.subscriptions[channel])
	delete(c.subscriptions, channel)

	c.sendStruct(&AckCommand{
		Ref:     cmd.Ref,
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

	socketGauge.WithLabelValues(fmt.Sprintf("%d", c.auth.Stream)).Inc()

	go c.ConnSub()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		fmt.Println("handled pong")
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

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
					Ref:     command.Ref,
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
			if timeoutEnabled && len(c.subscriptions) == 0 {
				fmt.Println("Closing socket due to no initial activity")
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
