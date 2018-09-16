package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	nats "github.com/nats-io/go-nats"
	pb "github.com/tcfw/evntsrc/pkg/bridge/protos"
	pbEvents "github.com/tcfw/evntsrc/pkg/event/protos"
)

type server struct {
	mu sync.Mutex // protects routeNotes
}

//newServer creates a ne struct to interface the auth server
func newServer() *server {
	return &server{}
}

func (s *server) Publish(ctx context.Context, request *pb.PublishRequest) (*pb.GeneralResponse, error) {
	err := s.ValidateAuth(request)
	if err != nil {
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
	sub, err := natsConn.ChanSubscribe(fmt.Sprintf("_USER.%d.%s", request.Stream, request.Channel), ch)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-ch:
			event := &pbEvents.Event{}

			err := stream.Send(msg.Data)
		}
	}
}

// @TODO move to grpc interceptor 
func (s *server) ValidateAuth(request interface{}) error {
	return nil
}
