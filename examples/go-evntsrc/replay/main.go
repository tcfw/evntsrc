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

	//Publish some events
	fmt.Printf("Publishing events (%v)\n", channel)
	for i := 0; i < 10; i++ {
		testMsg := &testMsg{Ts: time.Now()}
		msgBytes, _ := json.Marshal(testMsg)

		err = client.Publish(channel, msgBytes, "test")
		if err != nil {
			fmt.Printf("PUB ERR: %v\n", err.Error())
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

	//Hang until exit
	for {
		if received >= sent {
			fmt.Println("\nSuccessfully received all events :)")
			return
		}
		time.Sleep(100 * time.Millisecond)
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
