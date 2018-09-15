package bridge

import (
	"context"
	"sync"

	pb "github.com/tcfw/evntsrc/pkg/bridge/protos"
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

	return nil, nil
}

func (s *server) Subscribe(request *pb.SubscribeRequest, stream pb.BridgeService_SubscribeServer) error {
	err := s.ValidateAuth(request)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) ValidateAuth(request interface{}) error {
	return nil
}
