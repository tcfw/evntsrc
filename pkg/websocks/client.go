package websocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	nats "github.com/nats-io/go-nats"
	"github.com/tcfw/evntsrc/pkg/event"
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

	subscriptions map[string]chan bool

	auth *AuthCommand
}

func (c *Client) subscribe(channel string) {
	if _, ok := c.subscriptions[channel]; ok {
		ack := &AckSubUnSucCommand{
			Acktype: "err",
			Error:   fmt.Sprintf("already subscribed to %s", channel),
			Channel: channel,
		}
		ackBytes, _ := json.Marshal(ack)
		c.send <- ackBytes
		return
	}
	c.subscriptions[channel] = make(chan bool)

	go func(c *Client, channel string, unsub chan bool) {
		ch := make(chan *nats.Msg, 64)
		sub, err := natsConn.ChanSubscribe(fmt.Sprintf("_USER.%d.%s", c.auth.Stream, channel), ch)
		ack := &AckSubUnSucCommand{
			Acktype: "sub",
			Channel: channel,
		}
		if err != nil {
			ack.Acktype = "err"
			ack.Error = err.Error()
		}

		ackBytes, _ := json.Marshal(ack)
		c.send <- ackBytes

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

func (c *Client) unsubscribe(channel string) {
	if _, ok := c.subscriptions[channel]; !ok {
		ack := &AckSubUnSucCommand{
			Acktype: "err",
			Error:   fmt.Sprintf("not subscribed to %s", channel),
			Channel: channel,
		}
		ackBytes, _ := json.Marshal(ack)
		c.send <- ackBytes
		return
	}
	close(c.subscriptions[channel])
	delete(c.subscriptions, channel)

	ack := &AckSubUnSucCommand{
		Acktype: "unsub",
		Channel: channel,
	}

	ackBytes, _ := json.Marshal(ack)
	c.send <- ackBytes
}

func (c *Client) close() {
	c.conn.Close()
	for channel := range c.subscriptions {
		close(c.subscriptions[channel])
		delete(c.subscriptions, channel)
	}
}

func (c *Client) readPump() {
	defer func() {
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		command := &InboundCommand{}
		err = json.Unmarshal(message, command)
		if err != nil {
			log.Printf("Failed to parse command: %s\n", err.Error())
		} else {
			if command.Command != commandAuth && c.auth == nil {
				ack := &AckSubUnSucCommand{
					Acktype: "err",
					Error:   "No auth sent yet",
				}
				ackBytes, _ := json.Marshal(ack)
				c.send <- ackBytes
			} else {
				c.processCommand(command, message)
			}
		}
	}
}

func (c *Client) processCommand(command *InboundCommand, message []byte) {
	switch command.Command {
	case commandSubscribe:
		subcommand := &SubscribeCommand{}
		json.Unmarshal(message, subcommand)
		c.subscribe(subcommand.Channel)
	case commandUnsubscribe:
		subcommand := &UnsubscribeCommand{}
		json.Unmarshal(message, subcommand)
		c.unsubscribe(subcommand.Channel)
	case commandPublish:
		subcommand := &PublishCommand{}
		json.Unmarshal(message, subcommand)

		event := &event.Event{}
		event.SetID()
		event.Stream = c.auth.Stream
		event.Source = subcommand.Source
		if event.Source == "" {
			event.Source = "ws"
		}
		event.Subject = subcommand.Channel
		event.CEVersion = "0.1"
		event.Type = subcommand.Type
		event.TypeVersion = subcommand.TypeVersion
		event.ContentType = subcommand.ContentType
		event.Data = []byte(subcommand.Data)
		event.Time = time.Now()
		event.Metadata = map[string]interface{}{}
		event.Metadata["ws.source_ip"] = c.conn.RemoteAddr().String()

		eventJSONBytes, _ := json.Marshal(event)

		natsConn.Publish(fmt.Sprintf("_USER.%d.%s", c.auth.Stream, subcommand.Channel), eventJSONBytes)
	case commandAuth:
		subcommand := &AuthCommand{}
		json.Unmarshal(message, subcommand)

		//TODO verify auth info
		c.auth = subcommand
		ack := &AckSubUnSucCommand{
			Acktype: "OK",
		}
		ackBytes, _ := json.Marshal(ack)
		c.send <- ackBytes
	case commandReplay:
		subcommand := &ReplayCommand{}
		json.Unmarshal(message, subcommand)
		subcommand.Stream = c.auth.Stream

		fmt.Printf("%v\n", subcommand)

		repubBytes, _ := json.Marshal(subcommand)

		msg, err := natsConn.Request("replay.broadcast", repubBytes, time.Second*10)
		if err != nil {
			ack := &AckSubUnSucCommand{
				Acktype: "err",
				Error:   err.Error(),
			}
			ackBytes, _ := json.Marshal(ack)
			c.send <- ackBytes
		} else {
			ack := &AckSubUnSucCommand{
				Acktype: "err",
				Error:   string(msg.Data),
			}
			ackBytes, _ := json.Marshal(ack)
			c.send <- ackBytes
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