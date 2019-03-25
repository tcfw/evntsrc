package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	evntsrc "github.com/tcfw/evntsrc/external/go-evntsrc"
)

type testMsg struct {
	Ts time.Time `json:"ts"`
}

var apiKey string
var subOnly bool
var pubOnly bool
var channel string

func main() {
	flag.StringVar(&apiKey, "apikey", "", "API Key")
	flag.StringVar(&channel, "channel", "", "Specify a channel")
	flag.Parse()

	client, err := newClient()

	if channel == "" {
		channel = randChanName()
	}

	sent := 0
	received := 0

	startTime := time.Now()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(1)
	}()

	//Send ping every 5 seconds
	fmt.Printf("Publishing events (%v)\n", channel)
	for i := 0; i < 10; i++ {
		testMsg := &testMsg{Ts: time.Now()}
		msgBytes, _ := json.Marshal(testMsg)

		err = client.Publish(channel, msgBytes, "test")
		if err != nil {
			fmt.Printf("PUB ERR: %v\n", err.Error())
		} else {
			sent++
			fmt.Printf(".")
		}
	}

	fmt.Println("\nStarting replay")

	client.SubscribeFunc(channel, func(evnt *evntsrc.Event) {
		received++
		fmt.Printf("+")
		if received == sent {
			time.Sleep(1 * time.Second)
			fmt.Printf("\nSuccessfully received all events :)\n")
			os.Exit(0)
		}
	})

	//Wait for propagation
	time.Sleep(1 * time.Second)

	if err := client.Replay(channel, evntsrc.ReplayQuery{StartTime: &startTime}, nil); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	//Hang until exit
	select {}

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

	return client, err
}

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
