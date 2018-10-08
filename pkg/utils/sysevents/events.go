package sysevents

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	nats "github.com/nats-io/go-nats"
)

//@TODO support https://github.com/cloudevents/spec/blob/v0.1/spec.md

//Event is the basic structure all events should include
type Event struct {
	Channel   string
	Type      string
	CEVersion string
	Source    string
	ID        string
	Time      time.Time
	Metadata  map[string]interface{}
}

//EventInterface provides required funcs to identify common structured events
type EventInterface interface {
	GetType() string
	SetSource(string)
	SetTime(time.Time)
	SetID()
	SetMetadata(map[string]interface{})
	GetMetadata() map[string]interface{}
	GetChannel() string
	SetChannel(string)
}

//GetType returns event type
func (e *Event) GetType() string {
	return e.Type
}

//SetSource applies source
func (e *Event) SetSource(source string) {
	e.Source = source
}

//SetTime sets event time from timestamp
func (e *Event) SetTime(eventTime time.Time) {
	e.Time = eventTime
}

//SetID generates a new UUID for the event id
func (e *Event) SetID() {
	e.ID = uuid.New().String()
}

//SetMetadata overrides the existing metadata string map
func (e *Event) SetMetadata(md map[string]interface{}) {
	e.Metadata = md
}

//GetMetadata returns the current metadata
func (e *Event) GetMetadata() map[string]interface{} {
	return e.Metadata
}

//GetChannel returns the event channel
func (e *Event) GetChannel() string {
	return e.Channel
}

//SetChannel overwrites the event channel
func (e *Event) SetChannel(channel string) {
	e.Channel = channel
}

//AuthenticateEvent publishes to events.auth.auth
type AuthenticateEvent struct {
	*Event
	AuthType string `json:"authType"`
	Success  bool   `json:"success"`
	User     string `json:"user"`
	IP       string `json:"ip"`
	Err      string `json:"error,omitempty"`
}

//BroadcastEvent attempts to connect to nats server to pub any event and saves to stream
func BroadcastEvent(ctx context.Context, event EventInterface) error {
	hostname, _ := os.Hostname()
	nc, err := nats.Connect(os.Getenv("NATS_HOST"))
	if err != nil {
		log.Printf("Failed to broadcast event: %s", err)
		return err
	}
	defer nc.Close()

	event.SetSource(hostname)
	event.SetID()
	event.SetTime(time.Now())

	if event.GetChannel() == "" {
		event.SetChannel("broadcast")
	}

	event = appendContextUserInfo(ctx, event)

	buf, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return nc.Publish("io.evntsrc."+event.GetChannel(), buf)
}

//BroadcastNonStreamingEvent broadcasts an event like BroadcastEvent but uses the non-streaming engine
func BroadcastNonStreamingEvent(ctx context.Context, event EventInterface) error {
	nc, err := nats.Connect(os.Getenv("NATS_HOST"))
	if err != nil {
		log.Printf("Failed to broadcast event: %s", err)
		return err
	}
	defer nc.Close()

	hostname, _ := os.Hostname()
	event.SetSource(hostname)
	event.SetID()
	event.SetTime(time.Now())

	event = appendContextUserInfo(ctx, event)

	buf, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return nc.Publish(event.GetType(), buf)
}
