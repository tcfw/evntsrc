package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	evntsrc "github.com/tcfw/evntsrc/external/go-evntsrc"
)

type testMsg struct {
	Ts time.Time `json:"ts"`
}

var apiKey string

func main() {
	flag.StringVar(&apiKey, "apikey", "", "API Key")
	flag.Parse()

	client, err := newClient()

	client.SubscribeFunc("test", func(evnt *evntsrc.Event) {
		//Dislay ping latency results
		//TODO move decoding to client lib
		decoded, _ := base64.StdEncoding.DecodeString(string(evnt.Data))
		msg := &testMsg{}
		json.Unmarshal(decoded, msg)

		fmt.Printf("Took %v\n", time.Since(msg.Ts))
	})

	//Send ping every 5 seconds
	for {
		testMsg := &testMsg{Ts: time.Now()}
		msgBytes, _ := json.Marshal(testMsg)

		err = client.Publish("test", msgBytes, "test")
		if err != nil {
			fmt.Printf("PUB ERR: %v\n", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}

func newClient() (*evntsrc.APIClient, error) {
	//Create initial config
	client, err := evntsrc.NewEvntSrcClient(apiKey, 1)
	if err != nil {
		return nil, err
	}

	//Staging config
	client.Staging()

	//See our own events
	client.IgnoreSelf = false

	//Pipe any error to stdout
	go pipeErrors(client)

	return client, err
}

func pipeErrors(client *evntsrc.APIClient) {
	for {
		select {
		case err := <-client.Errors:
			fmt.Printf("ERR: %v\n", err.Error())
		}
	}
}
