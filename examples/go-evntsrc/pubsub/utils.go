package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

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
	flag.DurationVar(&pubSleepDuration, "rate", 1*time.Second, "Period of time between each publish")
	flag.Parse()
}
