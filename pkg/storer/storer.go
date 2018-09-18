package storer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/globalsign/mgo/bson"
	nats "github.com/nats-io/go-nats"
	"github.com/simplereach/timeutils"
	event "github.com/tcfw/evntsrc/pkg/event"
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
	defer natsConn.Close()

	monitorUserStreams()
	monitorReplayRequests()

	//Wait forever
	select {}
}

func monitorUserStreams() {
	log.Println("Watching for user streams...")
	natsConn.QueueSubscribe("_USER.>", "storers", func(m *nats.Msg) {
		event := &event.Event{}
		err := json.Unmarshal(m.Data, event)
		if err != nil {
			log.Printf("%s\n", err.Error())
			return
		}

		if isReplay, ok := event.Metadata["replay"]; ok && isReplay.(bool) {
			return
		}

		event.Store()
	})
}

func monitorReplayRequests() {
	log.Println("Watching for replay requests...")
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
				event.Metadata["replay"] = true
				jsonBytes, _ := json.Marshal(event)
				natsConn.Publish("_USER."+string(command.Stream)+"."+command.Channel, jsonBytes)
			}
		}
	})
}
