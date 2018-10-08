package streamauth

import (
	"context"
	"fmt"
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

	stream, _ := s.streamConn.Get(ctx, &evntsrc_streams.GetRequest{ID: 0})

	sk, err := s.Create(ctx, &pb.StreamKey{Stream: stream.GetID(), Label: "test"})
	if err != nil {
		t.Error(err)
	}

	_, err = s.Delete(ctx, sk)
	if err != nil {
		t.Error(err)
	}
}

func TestInvalidCreate(t *testing.T) {
	s := &server{}
	s.streamConn = NewMockStreamClient()

	invalidRequests := []*pb.StreamKey{
		{},
		{Stream: 0},
		{Label: "hi"},
	}

	for tk, key := range invalidRequests {
		t.Run(fmt.Sprintf("InvalidRequest(%d)", tk), func(t *testing.T) {
			_, err := s.Create(context.Background(), key)
			if err == nil {
				t.Errorf("R(%d) Invalid request still passed (%v)", tk, key)
			}
		})
	}
}

func TestGet(t *testing.T) {
	os.Setenv("DB_HOST", "192.168.99.100:32577")

	s := &server{}
	s.streamConn = NewMockStreamClient()

	jwtToken := testinghelpers.NulledJWT()

	md := metadata.Pairs("authorization", jwtToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	stream, _ := s.streamConn.Get(ctx, &evntsrc_streams.GetRequest{ID: 0})

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
}
func TestList(t *testing.T) {
	os.Setenv("DB_HOST", "192.168.99.100:32577")

	s := &server{}
	s.streamConn = NewMockStreamClient()

	jwtToken := testinghelpers.NulledJWT()

	md := metadata.Pairs("authorization", jwtToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	stream, _ := s.streamConn.Get(ctx, &evntsrc_streams.GetRequest{ID: 0})

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
}

//TestValidateOwnership tests remove validation of stream ownership
func TestValidateOwnership(t *testing.T) {

	os.Setenv("DB_HOST", "192.168.99.100:32577")

	s := &server{}
	s.streamConn = NewMockStreamClient()

	jwtToken := testinghelpers.NulledJWT()

	md := metadata.Pairs("authorization", jwtToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	stream, _ := s.streamConn.Get(ctx, &evntsrc_streams.GetRequest{ID: 0})

	ctx = metadata.NewIncomingContext(context.Background(), md)

	err := s.validateOwnership(ctx, stream.GetID())
	if err != nil {
		t.Error(err)
	}

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
	return &evntsrc_streams.Stream{
		ID:      999999999,
		Cluster: "test",
		Name:    "test",
		Owner:   "000000000000000000000000",
	}, nil
}

func (s *MockStreamClient) Create(ctx context.Context, in *evntsrc_streams.Stream, opts ...grpc.CallOption) (*evntsrc_streams.Stream, error) {
	return &evntsrc_streams.Stream{
		ID:      999999999,
		Cluster: in.Cluster,
		Name:    in.Name,
		Owner:   "000000000000000000000000",
	}, nil
}

func (s *MockStreamClient) List(ctx context.Context, in *evntsrc_streams.Empty, opts ...grpc.CallOption) (*evntsrc_streams.StreamList, error) {
	return &evntsrc_streams.StreamList{
		Streams: []*evntsrc_streams.Stream{
			{
				ID:      999999999,
				Cluster: "test",
				Name:    "test",
				Owner:   "000000000000000000000000",
			},
		},
	}, nil
}

func (s *MockStreamClient) Delete(ctx context.Context, in *evntsrc_streams.Stream, opts ...grpc.CallOption) (*evntsrc_streams.Empty, error) {
	return nil, nil
}
