package websocks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tcfw/evntsrc/internal/event"
)

const (
	commandSubscribe    = "sub"
	commandUnsubscribe  = "unsub"
	commandPublish      = "pub"
	commandPublishEvent = "epub"
	commandAuth         = "auth"
	commandReplay       = "replay"
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

//PublishCommand sends data (server-side event) through to NATS
type PublishCommand struct {
	*SubscribeCommand
	Data        string            `json:"data"`
	Source      string            `json:"source"`
	Type        string            `json:"type"`
	TypeVersion string            `json:"typeVersion"`
	ContentType string            `json:"contentType"`
	Metadata    map[string]string `json:"metadata"`
}

//PublishEventCommand sends an client-side event through NATS
type PublishEventCommand struct {
	*SubscribeCommand
	Event *event.Event `json:"event"`
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
	if !c.authKey.Permissions.Subscribe {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Permission denied"})
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
	if !c.authKey.Permissions.Publish || !c.validateRestriction(subcommand.Subject) {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Permission denied"})
		return
	}

	ev := event.NewEvent()
	ev.Stream = c.auth.Stream
	ev.Subject = subcommand.Subject
	ev.Type = subcommand.Type
	ev.TypeVersion = subcommand.TypeVersion
	ev.ContentType = subcommand.ContentType
	ev.Data = []byte(subcommand.Data)
	ev.Source = subcommand.Source
	ev.Metadata = map[string]string{}
	if len(subcommand.Metadata) > 0 {
		for mdK, mdV := range subcommand.Metadata {
			//Do not override existing metadata
			if _, ok := ev.Metadata[mdK]; !ok {
				ev.Metadata[mdK] = mdV
			}
		}
	}

	if err := c.publishEvent(ev); err != nil {
		c.sendStruct(&AckCommand{
			Ref:     command.Ref,
			Acktype: "Error",
			Error:   fmt.Sprintf("Failed to publish event: %s", err),
		})
	}

	c.sendStruct(&AckCommand{
		Ref:     command.Ref,
		Acktype: "OK",
	})
}

func (c *Client) doPublishEvent(command *InboundCommand, message []byte) {
	subcommand := &PublishEventCommand{}
	if err := json.Unmarshal(message, subcommand); err != nil {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Failed to parse command"})
		return
	}
	if !c.authKey.Permissions.Publish || !c.validateRestriction(subcommand.Subject) {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Permission denied"})
		return
	}

	ev := subcommand.Event
	ev.SetID()
	ev.Stream = c.auth.Stream
	ev.Subject = subcommand.Subject
	ev.Time = event.ZeroableTime{Time: time.Now()}

	if err := c.publishEvent(ev); err != nil {
		c.sendStruct(&AckCommand{
			Ref:     command.Ref,
			Acktype: "Error",
			Error:   fmt.Sprintf("Failed to publish event: %s", err),
		})
	}

	c.sendStruct(&AckCommand{
		Ref:     command.Ref,
		Acktype: "OK",
	})
}

func (c *Client) publishEvent(ev *event.Event) error {
	if ev.Source == "" {
		ev.Source = "ws"
	}
	if ev.Metadata == nil {
		ev.Metadata = map[string]string{}
	}
	ev.Metadata["source_ip"] = c.conn.RemoteAddr().String()

	channel := fmt.Sprintf("_USER.%d.%s", c.auth.Stream, ev.Subject)

	c.seqLock.Lock()
	c.seq[channel]++
	ev.Metadata["relative_seq"] = fmt.Sprintf("%s-%d", c.connectionID, c.seq[channel])
	c.seqLock.Unlock()

	ev.Metadata["_cid"] = c.connectionID

	if err := c.publisher.Publish(channel, ev); err != nil {
		return err
	}

	if err := c.storerAcks.WaitForKeyWithTimeout(ev.ID, 30*time.Second); err != nil {
		return fmt.Errorf("%s waiting for storage", err)
	}

	return nil
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

	if !c.authKey.Permissions.Replay || !c.validateRestriction(subcommand.Subject) {
		c.sendStruct(&AckCommand{Ref: command.Ref, Acktype: "Error", Error: "Permission denied"})
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
