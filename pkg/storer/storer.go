package storer

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/globalsign/mgo/bson"
	nats "github.com/nats-io/go-nats"
	"github.com/tcfw/evntsrc/pkg/event"
	"github.com/tcfw/evntsrc/pkg/utils/db"
	"github.com/tcfw/evntsrc/pkg/websocks"
	"github.com/tcfw/go-queue"
)

//StartMonitor subscripts to all user channels
func StartMonitor(nats string) {
	connectNats(nats)
	defer natsConn.Close()

	monitorUserStreams()
	monitorReplayRequests()

	//Wait forever
	select {}
}

type eventProcessor struct{}

func (ep *eventProcessor) Handle(job interface{}) {
	usrEvent := job.(*event.Event)

	if isReplay, ok := usrEvent.Metadata["replay"]; ok && isReplay == "true" {
		return
	}

	if err := usrEvent.Store(); err != nil {
		log.Printf("%s\n", err.Error())
	}
}

func monitorUserStreams() {
	log.Println("Watching for user streams...")

	dispatcher := queue.NewDispatcher(&eventProcessor{})
	if c := runtime.NumCPU(); c < 4 {
		dispatcher.MaxWorkers = 4
	}
	dispatcher.Run()

	natsConn.QueueSubscribe("_USER.>", "storers", func(m *nats.Msg) {
		event := &event.Event{}
		err := json.Unmarshal(m.Data, event)
		if err != nil {
			log.Printf("%s\n", err.Error())
			return
		}

		dispatcher.Queue(event)

	})
}

func monitorReplayRequests() {
	log.Println("Watching for replay requests...")
	natsConn.QueueSubscribe("replay.broadcast", "replayers", func(m *nats.Msg) {
		command := &websocks.ReplayCommand{}
		json.Unmarshal(m.Data, command)

		if command.JustMe && command.Dest == "" {
			natsConn.Publish(m.Reply, []byte("failed: no socket dest set"))
			return
		}

		dbConn, err := db.NewMongoDBSession()
		if err != nil {
			natsConn.Publish(m.Reply, []byte(fmt.Sprintf("failed: %s", err.Error())))
			return
		}
		defer dbConn.Close()

		collection := dbConn.DB("events").C("store")

		fq := bson.M{"stream": command.Stream}
		if command.Query.StartTime != nil && command.Query.EndTime == nil {
			fq["time.time"] = bson.M{"$gte": command.Query.StartTime}
		}
		if command.Query.StartTime != nil && command.Query.EndTime != nil {
			fq["time.time"] = bson.M{"$gte": command.Query.StartTime, "$lte": command.Query.EndTime}
		}
		query := collection.Find(fq).Sort("time")

		if c, _ := query.Count(); c == 0 {
			natsConn.Publish(m.Reply, []byte("failed: no events"))
		} else {
			natsConn.Publish(m.Reply, []byte("OK"))
			iter := query.Iter()
			event := event.Event{}
			for iter.Next(&event) {
				if command.Query.EndID != "" && event.ID == command.Query.EndID {
					break
				}
				event.Metadata["replay"] = "true"
				jsonBytes, _ := json.Marshal(event)
				if command.JustMe {
					natsConn.Publish(fmt.Sprintf("_CONN.%s", command.Dest), jsonBytes)
				} else {
					natsConn.Publish(fmt.Sprintf("_USER.%d.%s", command.Stream, command.Subject), jsonBytes)
				}
			}
		}
	})
}
