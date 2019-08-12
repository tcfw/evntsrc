package storer

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/gogo/protobuf/proto"
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

//Acknowledges appends a timestamp to the event
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

//ExtendTTL GRPC alias
func (s *server) ExtendTTL(ctx context.Context, req *pb.ExtendTTLRequest) (*pb.ExtendTTLResponse, error) {
	return &pb.ExtendTTLResponse{}, extendTTL(req)
}

//Query facilitates differing types of TTL queries
func (s *server) Query(req *pb.QueryRequest, stream pb.StorerService_QueryServer) error {
	switch qType := req.Query.(type) {
	case *pb.QueryRequest_Ttl:
		return s.handleTTLQuery(req, stream)
	default:
		return fmt.Errorf("Unknown query type %s", qType)
	}
}

//ReplayEvent queries DB for single event and publishes to NATS
func (s *server) ReplayEvent(ctx context.Context, req *pb.ReplayEventRequest) (*pb.ReplayEventResponse, error) {
	if err := s.validateStreamOwnership(ctx, req); err != nil {
		return nil, err
	}

	rQ, err := pgdb.Query(`SELECT * FROM event_store.events WHERE stream = $1 AND id = $2 LIMIT 1`, req.GetStream(), req.GetEventID())
	if err != nil {
		return nil, fmt.Errorf("sqlc: %s", err.Error())
	}

	event, err := scanEvent(rQ)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed replay proto marshal: %s", err.Error())
	}

	dest := fmt.Sprintf("_USER.%d.%s", event.Stream, event.Subject)
	err = natsConn.Publish(dest, bytes)
	if err != nil {
		return nil, fmt.Errorf("natspub: %s", err.Error())
	}

	return &pb.ReplayEventResponse{}, nil
}

//Store RPC storage request similar to NATS subscription
func (s *server) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	err := storeEvent(req.Event, pgdb)

	return &pb.StoreResponse{}, err
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := s.validateStreamOwnership(ctx, req); err != nil {
		return nil, err
	}

	_, err := pgdb.ExecContext(ctx, `DELETE FROM event_store.events WHERE stream = $1 AND id = $2 LIMIT 1`, req.GetStream(), req.GetEventID())

	return &pb.DeleteResponse{}, err
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
