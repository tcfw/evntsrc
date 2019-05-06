package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	nats "github.com/nats-io/go-nats"
	pb "github.com/tcfw/evntsrc/internal/bridge/protos"
	pbEvents "github.com/tcfw/evntsrc/internal/event/protos"
	"google.golang.org/grpc/metadata"
)

type server struct {
	mu sync.Mutex // protects routeNotes
}

//newServer creates a ne struct to interface the auth server
func newServer() *server {
	return &server{}
}

func (s *server) Publish(ctx context.Context, request *pb.PublishRequest) (*pb.GeneralResponse, error) {
	if err := s.ValidateAuth(request); err != nil {
		return nil, err
	}

	channel := fmt.Sprintf("_USER.%d.%s", request.Event.Stream, request.Event.Subject)
	bytes, err := json.Marshal(request.Event)
	if err != nil {
		return nil, err
	}

	natsConn.Publish(channel, bytes)

	return &pb.GeneralResponse{}, nil
}

func (s *server) Subscribe(request *pb.SubscribeRequest, stream pb.BridgeService_SubscribeServer) error {
	err := s.ValidateAuth(request)
	if err != nil {
		return err
	}

	ch := make(chan *nats.Msg, 64)
	if _, err = natsConn.ChanSubscribe(fmt.Sprintf("_USER.%d.%s", request.Stream, request.Channel), ch); err != nil {
		return err
	}

	for {
		select {
		case msg := <-ch:
			event := &pbEvents.Event{}
			json.Unmarshal(msg.Data, event)
			err := stream.Send(event)
			if err != nil {
				return err
			}
		}
	}
}

//Relay opens up a bi-directional stream
func (s *server) RelayEvents(stream pb.BridgeService_RelayEventsServer) error {
	close := make(chan bool, 1)
	in := make(chan *pbEvents.Event, 50)
	var wg sync.WaitGroup
	wg.Add(3)

	fmt.Println("Opened relay")

	//Read pipe
	go func() {
		for {
			event, err := stream.Recv()
			if err != nil {
				close <- true
				wg.Done()
				return
			}
			in <- event
		}
	}()

	//Forward pipe
	go func() {
		for {
			select {
			case event := <-in:
				channel := fmt.Sprintf("_USER.%d.%s", event.Stream, event.Subject)
				eventBytes, err := json.Marshal(event)
				if err != nil {
					fmt.Printf("failed to forward event: %s\n", err.Error())
				} else {
					natsConn.Publish(channel, eventBytes)
				}
			case <-close:
				close <- true
				wg.Done()
				return
			}
		}
	}()

	//Write pipe
	go func() {
		qid := ""

		ctx := stream.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if rQid, ok := md["qid"]; ok {
				qid = rQid[0]
			}
		}

		fmt.Println("QID: " + qid)

		writes := make(chan *nats.Msg, 50)
		natsConn.ChanQueueSubscribe("_USER.>", "relays-"+qid, writes)
		for {
			select {
			case msg := <-writes:
				fmt.Printf("Relaying event: %s\n", string(msg.Data))
				event := &pbEvents.Event{}
				err := json.Unmarshal(msg.Data, event)
				if err != nil {
					fmt.Printf("Failed to relay event: %s\n", err.Error())
				} else {
					stream.Send(event)
				}
			case <-close:
				close <- true
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()

	return nil
}

func (s *server) Replay(ctx context.Context, request *pb.ReplayRequest) (*pb.GeneralResponse, error) {
	return nil, nil
}

// @TODO move to grpc interceptor
func (s *server) ValidateAuth(request interface{}) error {
	return nil
}
