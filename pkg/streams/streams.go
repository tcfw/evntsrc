package streams

import (
	"context"
	"sync"

	pb "github.com/tcfw/evntsrc/pkg/streams/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	mu sync.Mutex
}

//newServer creates a new struct to interface the streams server
func newServer() *server {
	return &server{}
}

//Create @TODO
func (s *server) Create(ctx context.Context, request *pb.Stream) (*pb.Stream, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//Search @TODO
func (s *server) Search(ctx context.Context, request *pb.SearchRequest) (*pb.StreamList, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//List @TODO
func (s *server) List(ctx context.Context, request *pb.Empty) (*pb.StreamList, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//Get @TODO
func (s *server) Get(ctx context.Context, request *pb.GetRequest) (*pb.Stream, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}
