package main

import (
	"fmt"

	evntsrc "github.com/tcfw/evntsrc/pkg/go-evntsrc"
)

//newClient create a new evntsrc API client
func newClient() (*evntsrc.APIClient, error) {
	tempCrypto, err := evntsrc.EphemeralCrypto()
	if err != nil {
		return nil, err
	}

	options := []evntsrc.ClientOption{
		evntsrc.WithOwnEvents(), //See our own events
		evntsrc.WithCrypto(tempCrypto),
	}

	//Create initial config
	client, err := evntsrc.NewClient(apiKey, 1, options...)
	if err != nil {
		return nil, err
	}

	//Staging config
	client.Staging()

	//Listen for errors
	go func() {
		for {
			msg := <-client.Errors()
			fmt.Printf("##!!ERRR: %s\n", msg.Error())
		}
	}()

	return client, err
}
