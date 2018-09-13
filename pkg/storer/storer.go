package storer

import (
	"encoding/json"
	"fmt"

	nats "github.com/nats-io/go-nats"
	"github.com/simplereach/timeutils"
	event "github.com/tcfw/evntsrc/pkg/event"
	"gopkg.in/mgo.v2/bson"
)

//ReplayCommand instructs events to rebroadcast all events stored since time
type ReplayCommand struct {
	Command string         `json:"cmd"`
	Stream  int32          `json:"stream"`
	Channel string         `json:"channel"`
	Time    timeutils.Time `json:"startTime"`
}

//StartMonitor subscripts to all user channels
func StartMonitor(nats string) {
	connectNats(nats)

	monitorUserStreams()
	monitorReplayRequests()

	//Wait forever
	select {}
}

func monitorUserStreams() {
	natsConn.QueueSubscribe("_USER.>", "storers", func(m *nats.Msg) {
		event := &event.Event{}
		err := json.Unmarshal(m.Data, event)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		event.Store()
	})
}

func monitorReplayRequests() {
	natsConn.QueueSubscribe("replay.broadcast", "replayers", func(m *nats.Msg) {
		command := &ReplayCommand{}
		json.Unmarshal(m.Data, command)

		db, err := event.NewDBSession()
		if err != nil {
			natsConn.Publish(m.Reply, []byte(fmt.Sprintf("failed: %s", err.Error())))
			return
		}
		defer db.Close()

		collection := db.DB("events").C("store")

		fq := bson.M{"stream": command.Stream, "time": bson.M{"$gt": command.Time.Time}}
		query := collection.Find(fq).Sort("time")

		if c, _ := query.Count(); c == 0 {
			natsConn.Publish(m.Reply, []byte("failed: no events"))
		} else {
			natsConn.Publish(m.Reply, []byte("starting replay"))
			iter := query.Iter()
			event := event.Event{}
			for iter.Next(&event) {
				jsonBytes, _ := json.Marshal(event)
				natsConn.Publish("_USER."+string(command.Stream)+"."+command.Channel, jsonBytes)
			}
		}
	})
}
