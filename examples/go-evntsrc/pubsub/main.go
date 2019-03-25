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

type testMsg struct {
	Ts time.Time `json:"ts"`
}

var apiKey string
var subOnly bool
var pubOnly bool
var channel string

func main() {
	flag.StringVar(&apiKey, "apikey", "", "API Key")
	flag.BoolVar(&subOnly, "sub", false, "Subscription only")
	flag.BoolVar(&pubOnly, "pub", false, "Publish only")
	flag.StringVar(&channel, "channel", "", "Specify a channel")
	flag.Parse()

	client, err := newClient()

	if channel == "" {
		channel = randChanName()
	}

	sent := 0
	received := 0
	close := false

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("\nSent: %v  Received: %v\n", sent, received)

		close = true
		os.Exit(0)
	}()

	if (subOnly && !pubOnly) || (!pubOnly && !subOnly) {
		fmt.Printf("Subscribing(%v)...\n", channel)
		client.SubscribeFunc(channel, func(evnt *evntsrc.Event) {
			//Dislay ping latency results
			//TODO move decoding to client lib
			decoded, _ := base64.StdEncoding.DecodeString(string(evnt.Data))
			msg := &testMsg{}
			json.Unmarshal(decoded, msg)
			received++
			fmt.Printf("Took %v\n", time.Since(msg.Ts))
		})

		if subOnly && !pubOnly {
			//Hang
			select {}
		}
	}

	if (!subOnly && pubOnly) || (!pubOnly && !subOnly) {
		//Send ping every 5 seconds
		fmt.Printf("Publishing(%v)...\n", channel)
		for {
			if close {
				fmt.Printf("Stopping publishing")
				break
			}

			testMsg := &testMsg{Ts: time.Now()}
			msgBytes, _ := json.Marshal(testMsg)

			if !subOnly && pubOnly {
				fmt.Printf(".")
			}

			err = client.Publish(channel, msgBytes, "test")
			if err != nil {
				fmt.Printf("PUB ERR: %v\n", err.Error())
			} else {
				sent++
			}

			time.Sleep(1 * time.Second)
		}
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
