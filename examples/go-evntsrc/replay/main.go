package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	evntsrc "github.com/tcfw/evntsrc/external/go-evntsrc"
)

//Our simple test struct to record latency
type testMsg struct {
	Ts time.Time `json:"ts"`
}

var apiKey string
var subOnly bool
var pubOnly bool
var channel string

var sent int
var received int
var startTime time.Time

func main() {
	//Setup config
	setup()

	//Initialise a new Evntsrc API Client
	client, _ := newClient()

	//Publish some events ~ to be later received
	fmt.Printf("Publishing events (%v)\n", channel)
	for i := 0; i < 10; i++ {
		msgBytes, _ := json.Marshal(&testMsg{Ts: time.Now()})

		err := client.Publish(channel, msgBytes, "test")
		if err != nil {
			fmt.Printf("PUB ERR: %s\n", err.Error())
		} else {
			sent++
			fmt.Printf("|")
		}
	}

	//Start to listen for some events
	client.SubscribeFunc(channel, func(evnt *evntsrc.Event) {
		received++
		fmt.Printf("+")
	})

	//Replay recent events
	if err := client.Replay(channel, evntsrc.ReplayQuery{StartTime: &startTime}, nil); err != nil {
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

//newClient create a new evntsrc API client
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

	return client, err
}

//randChanName random test channel name
func randChanName() string {
	rand.Seed(time.Now().UnixNano())

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	n := 25

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return fmt.Sprintf("test_%s", string(b))
}

//setup read flag and set up globals
func setup() {
	flags()

	if channel == "" {
		channel = randChanName()
	}

	sent = 0
	received = 0
	startTime = time.Now()
}

//flags read in flags/config from command line
func flags() {
	flag.StringVar(&apiKey, "apikey", "", "API Key")
	flag.StringVar(&channel, "channel", "", "Specify a channel")
	flag.Parse()
}
