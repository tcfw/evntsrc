package stsmetrics

import (
	"encoding/json"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	nats "github.com/nats-io/go-nats"
	"github.com/tcfw/evntsrc/pkg/utils/db"
)

//StartWatch monitors timeseries requests
func StartWatch(natsEndpoint string) {
	connectNats(natsEndpoint)
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	log.Println("Starting watch...")

	natsConn.QueueSubscribe("analytics.timeseries", "analytics_watchers", func(m *nats.Msg) {
		request := &STSRequest{}
		json.Unmarshal(m.Data, request)

		streamCollection := dbConn.DB("events").C("store")

		fq := &bson.M{"stream": request.Stream}
		query := streamCollection.Find(fq)

		count, err := query.Count()
		if err != nil {
			log.Printf("Failed to get count: %s", err.Error())
			if m.Reply != "" {
				natsConn.Publish(m.Reply, []byte(err.Error()))
			}
			return
		}

		metric := &MetricTimeSeries{
			Stream: request.Stream,
			Count:  count,
			Time:   time.Now(),
		}

		metricsCollection := dbConn.DB("metrics").C("storage_count")

		err = metricsCollection.Insert(metric)
		if err != nil && m.Reply != "" {
			natsConn.Publish(m.Reply, []byte(err.Error()))
			return
		}

		err = metricsCollection.EnsureIndex(mgo.Index{
			Key:    []string{"stream"},
			Unique: false,
		})
		if err != nil {
			log.Printf("Failed to ensure index: %s", err.Error())
		}

		err = metricsCollection.EnsureIndex(mgo.Index{
			Key:    []string{"time"},
			Unique: false,
		})
		if err != nil {
			log.Printf("Failed to ensure index: %s", err.Error())
		}
	})

	select {}
}
