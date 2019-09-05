package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"
)

//Event is the main event structure for all events throughout the system
// TODO(tcfw) move to protobuf?
type Event struct {
	ID           string            `json:"eventId" bson:"_id"`
	Stream       int32             `json:"stream"`
	Time         ZeroableTime      `json:"eventTime"`
	Type         string            `json:"eventType"`
	TypeVersion  string            `json:"eventTypeVersion"`
	CEVersion    string            `json:"cloudEventVersion"`
	Source       string            `json:"eventSource"`
	Subject      string            `json:"eventSubject"`
	Acknowledged ZeroableTime      `json:"eventAcknowledged,omitempty"`
	Metadata     map[string]string `json:"extensions,omitempty"`
	ContentType  string            `json:"contentType,omitempty"`
	Data         []byte            `json:"data,omitempty"`
}

//NewEvent fills an event with basic info
func NewEvent() *Event {
	ev := &Event{
		CEVersion:   "0.1",
		TypeVersion: "0.1",
		Time:        ZeroableTime{Time: time.Now()},
	}
	ev.SetID()
	return ev
}

//SetID sets a new ID for the event based on UUID
func (e *Event) SetID() {
	e.ID = uuid.New().String()
}

//SetDataFromStruct converts a struct to JSON as raw bytes to be stored
func (e *Event) SetDataFromStruct(jsonStruct interface{}) error {
	jsonBytes, err := json.Marshal(jsonStruct)
	if err != nil {
		return err
	}

	e.Data = jsonBytes
	return nil
}

//SetDataFromString converts a string to bytes and stores in data
func (e *Event) SetDataFromString(data string) {
	e.Data = []byte(data)
}

//ToProtobuf converts structed event to protobuf event
func (e *Event) ToProtobuf() *pbEvent.Event {
	ev := &pbEvent.Event{}
	ev.ID = e.ID
	ev.ID = e.ID
	ev.Stream = e.Stream
	if !e.Time.IsZero() {
		ev.Time = &e.Time.Time
	}
	ev.Type = e.Type
	ev.TypeVersion = e.TypeVersion
	ev.CEVersion = e.CEVersion
	ev.Source = e.Source
	ev.Subject = e.Subject
	if !e.Acknowledged.IsZero() {
		ev.Acknowledged = &e.Acknowledged.Time
	}
	ev.Metadata = e.Metadata
	ev.ContentType = e.ContentType
	ev.Data = e.Data

	return ev
}
