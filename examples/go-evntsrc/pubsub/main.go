package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
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
var both bool
var channel string

var sent int
var received int
var close bool

func main() {
	//Setup config
	setup()

	//Initialise a new Evntsrc API Client
	client, _ := newClient()

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
		//TODO move decoding to client lib
		decoded, _ := base64.StdEncoding.DecodeString(string(evnt.Data))
		msg := &testMsg{}
		json.Unmarshal(decoded, msg)
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

		time.Sleep(1 * time.Second)
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

	//Pipe any error to stdout
	go pipeErrors(client)

	return client, err
}

//pipeErrors sends any global API errors to stdout
func pipeErrors(client *evntsrc.APIClient) {
	for {
		select {
		case err := <-client.Errors:
			fmt.Printf("ERR: %v\n", err.Error())
		}
	}
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

	sent = 0
	received = 0
	close = false
	both = !pubOnly && !subOnly

	if channel == "" {
		channel = randChanName()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("\nSent: %v  Received: %v\n", sent, received)

		close = true
		os.Exit(0)
	}()
}

//flags read in flags/config from command line
func flags() {
	flag.StringVar(&apiKey, "apikey", "", "API Key")
	flag.BoolVar(&subOnly, "sub", false, "Subscription only")
	flag.BoolVar(&pubOnly, "pub", false, "Publish only")
	flag.StringVar(&channel, "channel", "", "Specify a channel")
	flag.Parse()
}
