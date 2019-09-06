package main

import (
	"fmt"

	evntsrc "github.com/tcfw/evntsrc/pkg/go-evntsrc"
)

var apiKey string

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

	//Pipe any error to stdout
	go pipeErrors(client)

	//Staging config
	client.Staging()

	return client, err
}

//pipeErrors sends any global API errors to stdout
func pipeErrors(client *evntsrc.APIClient) {
	for {
		select {
		case err := <-client.Errors():
			fmt.Printf("ERR: %v\n", err.Error())
		}
	}
}
