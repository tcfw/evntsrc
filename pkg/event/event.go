package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/google/uuid"
)

//Event is the main event structure for all events throughout the system
// @TODO move to protobuf?
type Event struct {
	ID           string                 `json:"eventId" bson:"_id"`
	Stream       int32                  `json:"stream"`
	Time         time.Time              `json:"eventTime"`
	Type         string                 `json:"eventType"`
	TypeVersion  string                 `json:"eventTypeVersion"`
	CEVersion    string                 `json:"cloudEventVersion"`
	Source       string                 `json:"eventSource"`
	Subject      string                 `json:"eventSubject"`
	Acknowledged ZeroableTime           `json:"eventAcknowledged,omitempty"`
	Metadata     map[string]interface{} `json:"extensions,omitempty"`
	ContentType  string                 `json:"contentType,omitempty"`
	Data         []byte                 `json:"data,omitempty"`
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

//Store saves the event to DB
func (e *Event) Store() error {
	db, err := NewDBSession()
	if err != nil {
		return err
	}
	defer db.Close()

	collection := db.DB("events").C("store")

	err = collection.Insert(e)
	if err != nil {
		return err
	}

	err = collection.EnsureIndex(mgo.Index{
		Key:    []string{"stream"},
		Unique: false,
	})
	if err != nil {
		fmt.Printf("Error ensuring stream index: %s\n", err.Error())
	}

	return nil
}
