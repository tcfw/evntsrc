package main

import (
	"flag"
	"fmt"
	"time"

	mrand "math/rand"
)

//randChanName random test channel name
func randChanName() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	n := 25

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mrand.Intn(len(letterRunes))]
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
	flag.IntVar(&msgCount, "m", 10, "Number of messages to send")
	flag.Parse()
}
