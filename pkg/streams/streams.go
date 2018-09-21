package streams

import (
	"context"
	"sync"

	"github.com/globalsign/mgo/bson"
	pb "github.com/tcfw/evntsrc/pkg/streams/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	dBName         = "streams"
	collectionName = "streams"
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
	db, err := NewDBSession()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	collection := db.DB(dBName).C(collectionName)

	bsonq := bson.M{"owner": 1}
	query := collection.Find(bsonq)

	if c, _ := query.Count(); c == 0 {
		return &pb.StreamList{Streams: []*pb.Stream{}}, nil
	}

	streams := []*pb.Stream{}
	err = query.All(&streams)
	if err != nil {
		return nil, err
	}

	return &pb.StreamList{Streams: streams}, nil
}

//Get @TODO
func (s *server) Get(ctx context.Context, request *pb.GetRequest) (*pb.Stream, error) {
	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}
