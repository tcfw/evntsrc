package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"time"

	evntsrc "github.com/tcfw/evntsrc/pkg/go-evntsrc"
)

//Our simple test struct to record latency
type testMsg struct {
	Ts   time.Time `json:"ts"`
	Data []byte    `json:"data"`
}

var apiKey string
var subOnly bool
var pubOnly bool
var channel string
var msgCount int

var sent int
var received int
var startTime time.Time

func main() {
	//Setup config
	setup()

	//Initialise a new Evntsrc API Client
	client, err := newClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d := make([]byte, 1024)
	rand.Read(d)

	//Publish some events ~ to be later received
	fmt.Printf("Publishing events (%v)\n", channel)
	for i := 0; i < msgCount; i++ {
		msgBytes, _ := json.Marshal(&testMsg{Ts: time.Now(), Data: d})

		err := client.Publish(channel, msgBytes, "test")
		if err != nil {
			fmt.Printf("PUB ERR: %s\n", err.Error())
		} else {
			sent++
			fmt.Printf("P:%d ", sent)
		}
	}

	//Start to listen for events
	client.SubscribeFunc(channel, func(evnt *evntsrc.Event) {
		received++
		fmt.Printf("R:%d ", received)
	})

	//Replay recent events
	fmt.Println("Replay")
	if err := client.Replay(channel, evntsrc.ReplayQuery{StartTime: &startTime}, true); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	//Hang until received all sent
	for {
		if received >= sent {
			fmt.Println("\nReceived all events :)")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

}
