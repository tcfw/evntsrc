package bridge

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"

	nats "github.com/nats-io/go-nats"
	pb "github.com/tcfw/evntsrc/internal/bridge/protos"
	pbEvents "github.com/tcfw/evntsrc/internal/event/protos"
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
	bytes, err := proto.Marshal(request.Event)
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

	ch := make(chan *nats.Msg, 1000)
	if _, err = natsConn.ChanSubscribe(fmt.Sprintf("_USER.%d.%s", request.Stream, request.Channel), ch); err != nil {
		return err
	}

	for {
		select {
		case msg := <-ch:
			event := &pbEvents.Event{}
			proto.Unmarshal(msg.Data, event)
			err := stream.Send(event)
			if err != nil {
				return err
			}
		}
	}
}

// @TODO move to grpc interceptor
func (s *server) ValidateAuth(request interface{}) error {
	return nil
}
