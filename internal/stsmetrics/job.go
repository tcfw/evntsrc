package stsmetrics

import (
	"encoding/json"

	"github.com/globalsign/mgo/bson"
	"github.com/tcfw/evntsrc/internal/utils/db"
)

//JobRequest takes in a stream (or all) and sends a request to process timeseries info
func JobRequest(natsEndpoint string, stream int32) {
	connectNats(natsEndpoint)
	defer natsConn.Close()
	findJobs(stream)
}

func findJobs(stream int32) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	if stream > 0 {
		err := requestForStream(stream)
		if err != nil {
			panic(err)
		}
	} else {
		eventCollection := dbConn.DB("events").C("store")

		streams := eventCollection.Find(bson.M{})
		streamIds := []int32{}
		err := streams.Distinct("stream", &streamIds)
		if err != nil {
			panic(err)
		}

		for _, stream := range streamIds {
			requestForStream(stream)
		}
	}
}

func requestForStream(stream int32) error {
	requestBytes, err := json.Marshal(&STSRequest{Stream: stream})
	if err != nil {
		return err
	}
	return natsConn.Publish("analytics.timeseries", requestBytes)
}
