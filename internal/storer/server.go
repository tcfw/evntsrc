package storer

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/tcfw/evntsrc/internal/storer/protos"
	streamsPb "github.com/tcfw/evntsrc/internal/streams/protos"
	grpcRoundRobin "google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/metadata"
)

type server struct {
	streamConn streamsPb.StreamsServiceClient
}

//newServer adds connection to the streams service
func newServer() (*server, error) {

	conn, err := grpc.Dial("dns:///streams:443", grpc.WithInsecure(), grpc.WithBalancerName(grpcRoundRobin.Name))
	if err != nil {
		return nil, err
	}

	streamsConn := streamsPb.NewStreamsServiceClient(conn)

	return &server{streamsConn}, err
}

func (s *server) Acknowledge(ctx context.Context, req *pb.AcknowledgeRequest) (*pb.AcknowledgeResponse, error) {
	if err := s.validateStreamOwnership(ctx, req); err != nil {
		return nil, err
	}

	ackTime, err := ackEvent(req.GetStream(), req.GetEventID())
	if err != nil {
		return nil, err
	}

	return &pb.AcknowledgeResponse{Time: ackTime}, nil
}
func (s *server) ExtendTTL(ctx context.Context, req *pb.ExtendTTLRequest) (*pb.ExtendTTLResponse, error) {
	return &pb.ExtendTTLResponse{}, extendTTL(req)
}
func (s *server) Query(req *pb.QueryRequest, stream pb.StorerService_QueryServer) error {
	switch qType := req.Query.(type) {
	case *pb.QueryRequest_Ttl:
		return s.handleTTLQuery(req, stream)
	default:
		return fmt.Errorf("Unknown query type %s", qType)
	}
}

type streamedRequest interface {
	GetStream() int32
}

func (s *server) validateStreamOwnership(ctx context.Context, req streamedRequest) error {
	md, _ := metadata.FromIncomingContext(ctx)
	oMDCtx := metadata.NewOutgoingContext(ctx, md)

	//Try to get stream to validate access
	_, err := s.streamConn.Get(oMDCtx, &streamsPb.GetRequest{ID: req.GetStream()})
	return err
}
