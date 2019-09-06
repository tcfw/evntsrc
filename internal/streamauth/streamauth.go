package streamauth

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	pb "github.com/tcfw/evntsrc/internal/streamauth/protos"
	streams "github.com/tcfw/evntsrc/internal/streams/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"github.com/tcfw/evntsrc/internal/utils/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	if request.GetId() != "" {
		return nil, errors.New("New keys cannot have an id")
	}
	if request.GetKey() != "" {
		return nil, errors.New("New keys Cannot have a key")
	}
	if request.GetSecret() != "" {
		return nil, errors.New("New keys Cannot have a secret")
	}
	if request.GetStream() == 0 {
		return nil, errors.New("Must have stream set")
	}

	if err := s.validateOwnership(ctx, request.GetStream()); err != nil {
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

	if err = collection.Insert(request); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if err = collection.EnsureIndex(mgo.Index{
		Key:    []string{"stream"},
		Unique: false,
	}); err != nil {
		log.Printf("Error ensuring stream index: %s\n", err.Error())
	}

	return request, nil
}

//randString generates a n-lengthed alphanumberic string
func randString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz."
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes), nil
}

//List gets all user stream keys for a single stream
func (s *server) List(ctx context.Context, request *pb.ListRequest) (*pb.KeyList, error) {
	if err := s.validateOwnership(ctx, request.GetStream()); err != nil {
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
	if err = query.All(&streamKeys); err != nil {
		return nil, err
	}

	return &pb.KeyList{Keys: streamKeys}, nil
}

//ListAll provides a list of all stream keys
func (s *server) ListAll(ctx context.Context, request *pb.Empty) (*pb.KeyList, error) {
	//TODO(tcfw) change to paged request
	//TODO(tcfw) Verify user admin token claim

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

	if err = query.All(&streamKeys); err != nil {
		return nil, err
	}

	return &pb.KeyList{Keys: streamKeys}, nil
}

//Get retreives a single key for user stream
func (s *server) Get(ctx context.Context, request *pb.GetRequest) (*pb.StreamKey, error) {
	if err := s.validateOwnership(ctx, request.GetStream()); err != nil {
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
	if err = query.One(&key); err != nil {
		return nil, err
	}

	return &key, nil
}

//Update updates a stream key info
func (s *server) Update(ctx context.Context, request *pb.StreamKey) (*pb.StreamKey, error) {

	/*
		TODO(tcfw)
		Validate ownership of stream
		DB Connect
		Encrypt secret
		Update stream key
		Return stream key
	*/

	return nil, status.Errorf(codes.Unavailable, "Not implemented")
}

//Delete deletes a stream key perminately
func (s *server) Delete(ctx context.Context, request *pb.StreamKey) (*pb.Empty, error) {
	if err := s.validateOwnership(ctx, request.GetStream()); err != nil {
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
		conn, err := grpc.DialContext(ctx, "streams:443", tracing.GRPCClientOptions()...)
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

func (s *server) ValidateKeySecret(ctx context.Context, request *pb.KSRequest) (*pb.StreamKey, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dBName).C(collectionName)

	bsonq := bson.M{"stream": request.GetStream(), "key": request.GetKey(), "secret": request.GetSecret()}
	query := collection.Find(bsonq)

	if c, _ := query.Count(); c == 0 {
		return nil, status.Errorf(codes.NotFound, "Unknown key/secret")
	}

	sk := &pb.StreamKey{}
	err = query.One(sk)
	if err != nil {
		return nil, err
	}

	return sk, nil
}
