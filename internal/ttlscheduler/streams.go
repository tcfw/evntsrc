package ttlscheduler

import (
	"context"
	"time"

	streamsPb "github.com/tcfw/evntsrc/internal/streams/protos"
	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
	"google.golang.org/grpc"
	grpcRoundRobin "google.golang.org/grpc/balancer/roundrobin"
)

type basicStreamFetcher struct {
	streamsConn streamsPb.StreamsServiceClient
}

func (bsf *basicStreamFetcher) GetStreams() ([]*pb.Stream, error) {
	if bsf.streamsConn == nil {
		conn, err := grpc.Dial("dns:///streams:443", grpc.WithInsecure(), grpc.WithBalancerName(grpcRoundRobin.Name))
		if err != nil {
			return nil, err
		}

		bsf.streamsConn = streamsPb.NewStreamsServiceClient(conn)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ids, err := bsf.streamsConn.ListIds(ctx, &streamsPb.SearchRequest{})
	if err != nil {
		return nil, err
	}

	list := []*pb.Stream{}
	for _, id := range ids.ID {
		list = append(list, &pb.Stream{Id: id, MsgRate: 0})
	}

	return list, nil
}
