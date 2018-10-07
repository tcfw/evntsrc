package streamauth

import (
	"context"
	"crypto/rand"
	"log"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/google/uuid"
	pb "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	streams "github.com/tcfw/evntsrc/pkg/streams/protos"
	"github.com/tcfw/evntsrc/pkg/utils/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

const (
	dBName         = "streams"
	collectionName = "keys"
)

type server struct {
	mu         sync.Mutex
	streamConn streams.StreamsServiceClient
}

//newServer creates a new struct to interface the streams server
func newServer() *server {
	return &server{}
}

//Create creates a new stream key
func (s *server) Create(ctx context.Context, request *pb.StreamKey) (*pb.StreamKey, error) {
	err := request.Validate(true)
	if err != nil {
		return nil, err
	}

	err = s.validateOwnership(ctx, request.GetStream())
	if err != nil {
		return nil, err
	}

	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dBName).C(collectionName)

	request.Id = uuid.New().String()

	key, err := randString(16)
	if err != nil {
		return nil, err
	}

	secret, err := randString(32)
	if err != nil {
		return nil, err
	}

	request.Key = key
	request.Secret = secret

	err = collection.Insert(request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = collection.EnsureIndex(mgo.Index{
		Key:    []string{"stream"},
		Unique: false,
	})
	if err != nil {
		log.Printf("Error ensuring stream index: %s\n", err.Error())
	}

	return request, nil
}

//randString generates a n-lengthed alphanumberic string
func randString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz."
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes), nil
}

//List gets all user stream keys for a single stream
func (s *server) List(ctx context.Context, request *pb.ListRequest) (*pb.KeyList, error) {
	err := s.validateOwnership(ctx, request.GetStream())
	if err != nil {
		return nil, err
	}

	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dBName).C(collectionName)

	bsonq := bson.M{"stream": request.GetStream()}
	query := collection.Find(bsonq)

	if c, _ := query.Count(); c == 0 {
		return &pb.KeyList{Keys: []*pb.StreamKey{}}, nil
	}

	streamKeys := []*pb.StreamKey{}
	err = query.All(&streamKeys)
	if err != nil {
		return nil, err
	}

	return &pb.KeyList{Keys: streamKeys}, nil
}

//ListAll @TODO ~> change to paged request
func (s *server) ListAll(ctx context.Context, request *pb.Empty) (*pb.KeyList, error) {

	//@TODO Verify user admin token claim

	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dBName).C(collectionName)

	query := collection.Find(bson.M{})

	if c, _ := query.Count(); c == 0 {
		return &pb.KeyList{Keys: []*pb.StreamKey{}}, nil
	}

	streamKeys := []*pb.StreamKey{}
	err = query.All(&streamKeys)
	if err != nil {
		return nil, err
	}

	return &pb.KeyList{Keys: streamKeys}, nil
}

//Get retreives a single key for user stream
func (s *server) Get(ctx context.Context, request *pb.GetRequest) (*pb.StreamKey, error) {
	err := s.validateOwnership(ctx, request.GetStream())
	if err != nil {
		return nil, err
	}

	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dBName).C(collectionName)

	bsonq := bson.M{"stream": request.GetStream(), "_id": request.GetId()}
	query := collection.Find(bsonq)

	if c, _ := query.Count(); c == 0 {
		return nil, status.Errorf(codes.NotFound, "Unknown key")
	}

	key := pb.StreamKey{}
	err = query.One(&key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

//Update @TODO
func (s *server) Update(ctx context.Context, request *pb.StreamKey) (*pb.StreamKey, error) {

	/*
		Validate ownership of stream
		DB Connect
		Encrypt secret
		Update stream key
		Return stream key
	*/

	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//Delete @TODO
func (s *server) Delete(ctx context.Context, request *pb.StreamKey) (*pb.Empty, error) {
	err := s.validateOwnership(ctx, request.GetStream())
	if err != nil {
		return nil, err
	}

	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dBName).C(collectionName)

	bsonq := bson.M{"stream": request.GetStream(), "_id": request.GetId()}
	err = collection.Remove(bsonq)

	return &pb.Empty{}, err
}

//validateOwnership contacts streams to verify stream ownership via jwt token
func (s *server) validateOwnership(ctx context.Context, stream int32) error {
	if s.streamConn == nil {
		conn, err := grpc.DialContext(ctx, "streams:443", grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer func() {
			conn.Close()
			s.streamConn = nil
		}()
		s.streamConn = streams.NewStreamsServiceClient(conn)
	}

	md, _ := metadata.FromIncomingContext(ctx)
	octx := metadata.NewOutgoingContext(ctx, md)
	_, err := s.streamConn.Get(octx, &streams.GetRequest{ID: stream})

	if err != nil {
		return status.Errorf(codes.NotFound, "Unknown stream: "+err.Error())
	}
	return nil
}