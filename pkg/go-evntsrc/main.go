package evntsrc

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tcfw/evntsrc/internal/event"
	"github.com/tcfw/evntsrc/internal/websocks"
)

const (
	apiEndpoint        = "evntsrc.io"
	stagingAPIEndpoint = "staging.evntsrc.io"
	apiVersion         = "v1"
	endpointAPI        = "api"
	endpointIngress    = "ingress"
)

//APIClient the basic struct for the api
type APIClient struct {
	auth                       string
	stream                     int32
	token                      string
	Endpoint                   string
	httpClient                 *http.Client
	socket                     *websocket.Conn
	connectionID               string
	AppID                      string
	AppVer                     string
	Debug                      bool
	WaitForPublishConfirmation bool
	acks                       map[string]*ackCT
	subscriptions              map[string][]*subscription

	ReadPipe   chan []byte
	writePipe  chan *websocks.PublishCommand
	replayPipe chan *websocks.ReplayCommand
	subPipe    chan *websocks.SubscribeCommand
	close      chan bool
	errCh      chan error
	AcksCh     chan *websocks.AckCommand
	ackL       sync.RWMutex
	ackC       *sync.Cond

	//Config options
	options *ClientOptions
}

//NewEvntSrcClient creates a new client instance for interacting with evntsrc.io
func NewEvntSrcClient(auth string, streamID int32, options ...ClientOption) (*APIClient, error) {
	api := &APIClient{
		auth:                       auth,
		stream:                     streamID,
		Endpoint:                   apiEndpoint,
		httpClient:                 newHTTPClient(),
		ReadPipe:                   make(chan []byte, 256),
		writePipe:                  make(chan *websocks.PublishCommand, 56),
		replayPipe:                 make(chan *websocks.ReplayCommand, 5),
		subPipe:                    make(chan *websocks.SubscribeCommand, 10),
		close:                      make(chan bool, 1),
		errCh:                      make(chan error, 256),
		AcksCh:                     make(chan *websocks.AckCommand, 256),
		AppVer:                     "0.1",
		acks:                       map[string]*ackCT{},
		subscriptions:              map[string][]*subscription{},
		Debug:                      false,
		WaitForPublishConfirmation: true,
	}

	if auth == "" {
		return nil, fmt.Errorf("Empty auth string is not allowed")
	}
	if streamID == 0 {
		return nil, fmt.Errorf("Empty stream ID is not allowed")
	}

	defaultOptions := &ClientOptions{
		IgnoreSelf: true,
	}

	for _, option := range options {
		option(defaultOptions)
	}

	api.options = defaultOptions

	api.ackC = sync.NewCond(api.ackL.RLocker())

	go api.watchAcks()

	return api, nil
}

type ackCT struct {
	ts  time.Time
	ack *websocks.AckCommand
}

type subscriptionType string

const (
	funcSubType = "func"
	chanSubType = "chan"
)

//Event is the main structure of an event
type Event = event.Event

//ReplayQuery specifies where to start the replay in time
type ReplayQuery = websocks.ReplayRange

type subscription struct {
	subType subscriptionType
	ch      chan *Event
	f       func(*Event)
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 15 * time.Second,
			MaxIdleConns:    10,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
}

//watchAcks listens for acknowlegements and stores them in the
//temporary ack map
func (api *APIClient) watchAcks() {
	gci := time.Tick(1 * time.Minute)
	for {
		select {
		case <-gci:
			go api.gcAcks()
			break
		case ack := <-api.AcksCh:
			if ack.Ref != "" {
				api.ackC.L.Lock()
				api.acks[ack.Ref] = &ackCT{ts: time.Now(), ack: ack}
				api.ackC.L.Unlock()
				api.ackC.Broadcast()
			}
			break
		}
	}
}

//gcAcks garbage collections recent acknowledgements that are
//older than a 1 minute
func (api *APIClient) gcAcks() {
	defer func() {
		api.ackL.Unlock()
	}()
	api.ackL.Lock()
	tmo := time.Now().Add(-time.Minute)
	for ref, act := range api.acks {
		if act.ts.Before(tmo) {
			delete(api.acks, ref)
		}
	}
}

//waitForResponse waits for an acknowledge msg from websocks relating to command ref
//or fail on timeout after 30 seconds
func (api *APIClient) waitForResponse(cmdRef string) (bool, error) {
	var vFerr error
	vFound := make(chan bool)
	go func() {
		for {
			api.ackC.L.Lock()
			api.ackC.Wait()
			if ack, ok := api.acks[cmdRef]; ok {
				api.ackC.L.Unlock()
				ret := ack.ack.Acktype == "OK"
				if ack.ack.Error != "" {
					vFerr = errors.New(ack.ack.Error)
				}
				vFound <- ret
				return
			}
			api.ackC.L.Unlock()
		}
	}()
	select {
	case <-time.After(30 * time.Second):
		return false, errors.New("ACK timeout")
	case rVal := <-vFound:
		return rVal, vFerr

	}
}

func (api *APIClient) formatURL(apiType string, methodEndpoint string) string {
	protocol := "https"
	if apiType == "realtime" {
		protocol = "wss"
	}
	return fmt.Sprintf("%s://%s.%s/%s/%s", protocol, apiType, api.Endpoint, apiVersion, methodEndpoint)
}

//Staging switches to the staging env endpoints and enabled debug mode
func (api *APIClient) Staging() {
	api.Endpoint = stagingAPIEndpoint
	// api.Debug = true
}

func (api *APIClient) getSocket() *websocket.Conn {
	return api.socket
}
