package websocks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tcfw/evntsrc/pkg/event"
)

const (
	commandSubscribe   = "sub"
	commandUnsubscribe = "unsub"
	commandPublish     = "pub"
	commandAuth        = "auth"
	commandReplay      = "replay"
)

//InboundCommand is the basic struct for all commands coming from browser
type InboundCommand struct {
	Command string `json:"cmd"`
	Ref     string `json:"ref"`
}

//AuthCommand provides streaming information and verification
type AuthCommand struct {
	*InboundCommand
	Stream int32  `json:"stream"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

//SubscribeCommand starts a subscription
type SubscribeCommand struct {
	*InboundCommand
	Subject string `json:"subject"`
}

//PublishCommand sends data through to NATS
type PublishCommand struct {
	*SubscribeCommand
	Data        string `json:"data"`
	Source      string `json:"source"`
	Type        string `json:"type"`
	TypeVersion string `json:"typeVersion"`
	ContentType string `json:"contentType"`
}

//UnsubscribeCommand starts a subscription
type UnsubscribeCommand struct {
	*InboundCommand
	Subject string `json:"subject"`
}

//AckCommand provides error responses to WS clients
type AckCommand struct {
	Acktype string `json:"acktype"`
	Channel string `json:"channel"`
	Error   string `json:"error,omitempty"`
	Ref     string `json:"ref"`
}

//ReplayCommand instructs events to rebroadcast all events stored since time
type ReplayCommand struct {
	*SubscribeCommand
	Stream int32       `json:"stream"`
	JustMe bool        `json:"justme"`
	Dest   string      `json:"dest"`
	Query  ReplayRange `json:"query"`
}

//ReplayRange specifies the time and/or ID's to include in the replay
type ReplayRange struct {
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
	EndID     string     `json:"endID"`
}

//ConnectionInfo basic information about the current connection
type ConnectionInfo struct {
	Ref          string `json:"ref"`
	ConnectionID string `json:"connectionID"`
}

func (c *Client) doSubscribe(command *InboundCommand, message []byte) {
	subcommand := &SubscribeCommand{}
	if err := json.Unmarshal(message, subcommand); err != nil {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Failed to parse command"})
		return
	}
	c.Subscribe(subcommand.Subject, command)
}

func (c *Client) doUnsubscribe(command *InboundCommand, message []byte) {
	subcommand := &UnsubscribeCommand{}
	if err := json.Unmarshal(message, subcommand); err != nil {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Failed to parse command"})
		return
	}
	c.Unsubscribe(subcommand.Subject, command)
}

func (c *Client) doPublish(command *InboundCommand, message []byte) {
	subcommand := &PublishCommand{}
	if err := json.Unmarshal(message, subcommand); err != nil {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Failed to parse command"})
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
	c.sendStruct(&AckCommand{
		Ref:     command.Ref,
		Acktype: "OK",
	})
}

func (c *Client) doAuth(command *InboundCommand, message []byte) {
	subcommand := &AuthCommand{}
	if err := json.Unmarshal(message, subcommand); err != nil {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Failed to parse command"})
		return
	}

	if err := c.validateAuth(context.Background(), subcommand); err != nil {
		c.sendStruct(&AckCommand{
			Ref:     command.Ref,
			Acktype: "Failed",
			Error:   err.Error(),
		})
		return
	}

	c.auth = subcommand
	c.sendStruct(&AckCommand{
		Ref:     command.Ref,
		Acktype: "OK",
	})
	go c.broadcastConnect()
	c.sendStruct(&ConnectionInfo{
		Ref:          command.Ref,
		ConnectionID: c.connectionID,
	})
}

func (c *Client) doReplay(command *InboundCommand, message []byte) {
	subcommand := &ReplayCommand{}
	if err := json.Unmarshal(message, subcommand); err != nil {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Failed to parse command"})
		return
	}
	subcommand.Stream = c.auth.Stream

	if subcommand.JustMe {
		subcommand.Dest = c.connectionID
	}

	//Request replay to storer svc
	repubBytes, _ := json.Marshal(subcommand)
	msg, err := natsConn.Request("replay.broadcast", repubBytes, time.Second*10)
	ack := &AckCommand{
		Ref:     command.Ref,
		Acktype: "err",
	}
	if err != nil {
		ack.Error = err.Error()
	} else {
		if string(msg.Data) != "OK" {
			ack.Error = string(msg.Data)
		} else {
			ack.Acktype = "OK"
		}
	}
	c.sendStruct(ack)
}
