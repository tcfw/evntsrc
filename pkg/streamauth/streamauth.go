package streamauth

import (
	"context"
	"sync"

	pb "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const dbName = "streams"
const dbCollection = "keys"

type server struct {
	mu sync.Mutex
}

//newServer creates a new struct to interface the streams server
func newServer() *server {
	return &server{}
}

//Create @TODO
func (s *server) Create(ctx context.Context, request *pb.StreamKey) (*pb.StreamKey, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//List @TODO
func (s *server) List(ctx context.Context, request *pb.ListRequest) (*pb.KeyList, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//ListAll @TODO
func (s *server) ListAll(ctx context.Context, request *pb.Empty) (*pb.KeyList, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//Get @TODO
func (s *server) Get(ctx context.Context, request *pb.GetRequest) (*pb.StreamKey, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//Update @TODO
func (s *server) Update(ctx context.Context, request *pb.StreamKey) (*pb.StreamKey, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//Delete @TODO
func (s *server) Delete(ctx context.Context, request *pb.StreamKey) (*pb.Empty, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}
