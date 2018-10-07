package streams

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
	pb "github.com/tcfw/evntsrc/pkg/streams/protos"
	utils "github.com/tcfw/evntsrc/pkg/utils/authorization"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	dBName         = "streams"
	collectionName = "streams"
)

//Server core struct
type Server struct {
	mu sync.Mutex
}

//NewServer creates a new struct to interface the streams server
func NewServer() *Server {
	return &Server{}
}

func (s *Server) validateCreate(request *pb.Stream) error {
	if request.GetID() != 0 {
		return errors.New("ID must not be set")
	}

	if request.GetOwner() != "" {
		return errors.New("Owner must not be set")
	}

	if request.GetName() == "" {
		return errors.New("Name must be set")
	}

	if request.GetCluster() == "" {
		return errors.New("Cluster must be set")
	}

	return nil
}

//Create @TODO
func (s *Server) Create(ctx context.Context, request *pb.Stream) (*pb.Stream, error) {
	err := s.validateCreate(request)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	db, err := NewDBSession()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	collection := db.DB(dBName).C(collectionName)

	userClaims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	request.Owner = userClaims["sub"].(string)
	request.ID = int32(time.Now().Unix())

	err = collection.Insert(request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return s.Get(ctx, &pb.GetRequest{ID: request.ID})
}

//List @TODO
func (s *Server) List(ctx context.Context, request *pb.Empty) (*pb.StreamList, error) {
	db, err := NewDBSession()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	collection := db.DB(dBName).C(collectionName)

	userClaims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	bsonq := bson.M{"owner": userClaims["sub"]}
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
func (s *Server) Get(ctx context.Context, request *pb.GetRequest) (*pb.Stream, error) {

	db, err := NewDBSession()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	collection := db.DB(dBName).C(collectionName)

	userClaims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	bsonq := bson.M{"owner": userClaims["sub"], "_id": request.GetID()}
	query := collection.Find(bsonq)

	if c, _ := query.Count(); c == 0 {
		return nil, status.Errorf(codes.NotFound, "Unknown stream")
	}

	stream := pb.Stream{}
	err = query.One(&stream)
	if err != nil {
		return nil, err
	}

	return &stream, nil
}

//Delete @TODO
func (s *Server) Delete(ctx context.Context, request *pb.Stream) (*pb.Empty, error) {

	db, err := NewDBSession()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	collection := db.DB(dBName).C(collectionName)

	userClaims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	bsonq := bson.M{"owner": userClaims["sub"], "_id": request.GetID()}
	err = collection.Remove(bsonq)

	return &pb.Empty{}, err
}
