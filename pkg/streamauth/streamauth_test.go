package streamauth

import (
	"context"
	"os"
	"testing"

	pb "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	streams "github.com/tcfw/evntsrc/pkg/streams"
	evntsrc_streams "github.com/tcfw/evntsrc/pkg/streams/protos"
	"github.com/tcfw/evntsrc/pkg/utils/testinghelpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestCreate(t *testing.T) {
	os.Setenv("DB_HOST", "192.168.99.100:32577")

	s := &server{}
	s.streamConn = NewMockStreamClient()

	jwtToken := testinghelpers.NulledJWT()

	md := metadata.Pairs("authorization", jwtToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	stream, err := s.streamConn.Create(ctx, &evntsrc_streams.Stream{
		Cluster: "test",
		Name:    "test",
	})
	if err != nil {
		t.Error(err)
	}

	sk, err := s.Create(ctx, &pb.StreamKey{Stream: stream.GetID(), Label: "test"})
	if err != nil {
		t.Error(err)
	}

	_, err = s.Delete(ctx, sk)
	if err != nil {
		t.Error(err)
	}
	s.streamConn.Delete(ctx, stream)
}

func TestGet(t *testing.T) {
	os.Setenv("DB_HOST", "192.168.99.100:32577")

	s := &server{}
	s.streamConn = NewMockStreamClient()

	jwtToken := testinghelpers.NulledJWT()

	md := metadata.Pairs("authorization", jwtToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	stream, err := s.streamConn.Create(ctx, &evntsrc_streams.Stream{
		Cluster: "test",
		Name:    "test",
	})
	if err != nil {
		t.Error(err)
	}

	sk, err := s.Create(ctx, &pb.StreamKey{Stream: stream.GetID(), Label: "test"})
	if err != nil {
		t.Error(err)
	}

	skG, err := s.Get(ctx, &pb.GetRequest{Stream: stream.GetID(), Id: sk.GetId()})
	if err != nil {
		t.Error(err)
	}
	if skG.GetId() != sk.GetId() {
		t.Error("Created key does not match get key")
	}

	s.Delete(ctx, sk)
	s.streamConn.Delete(ctx, stream)
}
func TestList(t *testing.T) {
	os.Setenv("DB_HOST", "192.168.99.100:32577")

	s := &server{}
	s.streamConn = NewMockStreamClient()

	jwtToken := testinghelpers.NulledJWT()

	md := metadata.Pairs("authorization", jwtToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	stream, err := s.streamConn.Create(ctx, &evntsrc_streams.Stream{
		Cluster: "test",
		Name:    "test",
	})
	if err != nil {
		t.Error(err)
	}

	sk, err := s.Create(ctx, &pb.StreamKey{Stream: stream.GetID(), Label: "test"})
	if err != nil {
		t.Error(err)
	}

	keyList, err := s.List(ctx, &pb.ListRequest{Stream: stream.GetID()})
	if err != nil {
		t.Error(err)
	}

	if len(keyList.GetKeys()) != 1 {
		t.Error("Count does not match")
	}

	if keyList.GetKeys()[0].GetId() != sk.GetId() {
		t.Error("Retreived key does not match created key")
	}

	_, err = s.Delete(ctx, sk)
	if err != nil {
		t.Error(err)
	}
	s.streamConn.Delete(ctx, stream)
}

//TestValidateOwnership tests remove validation of stream ownership
func TestValidateOwnership(t *testing.T) {

	os.Setenv("DB_HOST", "192.168.99.100:32577")

	s := &server{}
	s.streamConn = NewMockStreamClient()

	jwtToken := testinghelpers.NulledJWT()

	md := metadata.Pairs("authorization", jwtToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	stream, err := s.streamConn.Create(ctx, &evntsrc_streams.Stream{
		Cluster: "test",
		Name:    "test",
	})
	if err != nil {
		t.Error(err)
	}

	ctx = metadata.NewIncomingContext(context.Background(), md)

	err = s.validateOwnership(ctx, stream.GetID())
	if err != nil {
		t.Error(err)
	}

	s.streamConn.Delete(ctx, stream)
}

func TestListAll(t *testing.T) {
	s := &server{}
	s.streamConn = NewMockStreamClient()

	_, err := s.ListAll(context.Background(), &pb.Empty{})
	if err != nil {
		t.Error(err)
	}
}

func NewMockStreamClient() *MockStreamClient {
	return &MockStreamClient{Server: &streams.Server{}}
}

//MockStreamClient bridges to a mocked streams server
type MockStreamClient struct {
	*streams.Server
}

func (s *MockStreamClient) Get(ctx context.Context, in *evntsrc_streams.GetRequest, opts ...grpc.CallOption) (*evntsrc_streams.Stream, error) {
	return s.Server.Get(ctx, in)
}

func (s *MockStreamClient) Create(ctx context.Context, in *evntsrc_streams.Stream, opts ...grpc.CallOption) (*evntsrc_streams.Stream, error) {
	return s.Server.Create(ctx, in)
}

func (s *MockStreamClient) List(ctx context.Context, in *evntsrc_streams.Empty, opts ...grpc.CallOption) (*evntsrc_streams.StreamList, error) {
	return s.Server.List(ctx, in)
}

func (s *MockStreamClient) Delete(ctx context.Context, in *evntsrc_streams.Stream, opts ...grpc.CallOption) (*evntsrc_streams.Empty, error) {
	return s.Server.Delete(ctx, in)
}
