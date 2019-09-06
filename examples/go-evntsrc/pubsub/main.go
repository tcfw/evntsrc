package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	evntsrc "github.com/tcfw/evntsrc/pkg/go-evntsrc"
)

//Simple test struct to record latency
type testMsg struct {
	Ts time.Time `json:"ts"`
}

var subOnly bool
var pubOnly bool
var both bool
var channel string
var pubSleepDuration time.Duration

var sent int
var received int
var close bool

func main() {
	//Setup config
	setup()

	//Initialise a new Evntsrc API Client
	client, err := newClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if (!subOnly && pubOnly) || both {
		//Start a publisher
		go startPublish(client)
	}

	if (subOnly && !pubOnly) || both {
		//Start a subscriber
		go startSubscribe(client)
	}

	//Hang until interrupt sig (ctrl+c)
	select {}
}

func startSubscribe(client *evntsrc.APIClient) {
	fmt.Printf("Subscribing (%v)...\n", channel)
	client.SubscribeFunc(channel, func(evnt *evntsrc.Event) {
		//Dislay ping latency results
		msg := &testMsg{}
		json.Unmarshal(evnt.Data, msg)
		received++
		fmt.Printf("Took %v\n", time.Since(msg.Ts))
	})
}

//startPublishing sends a ping struct ever 5 seconds
func startPublish(client *evntsrc.APIClient) {
	fmt.Printf("Publishing (%v)...\n", channel)
	for {
		if close {
			break
		}

		testMsg := &testMsg{Ts: time.Now()}
		msgBytes, _ := json.Marshal(testMsg)

		err := client.Publish(channel, msgBytes, "test")
		if err != nil {
			fmt.Printf("PUB ERR: %v\n", err.Error())
		} else {
			sent++

			if !subOnly && pubOnly {
				fmt.Printf(".")
			}
		}

		time.Sleep(pubSleepDuration)
	}
}
