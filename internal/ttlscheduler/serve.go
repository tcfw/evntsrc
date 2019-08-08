package ttlscheduler

import (
	"fmt"
	"log"
	"net"

	"github.com/tcfw/evntsrc/internal/tracing"
	pb "github.com/tcfw/evntsrc/internal/ttlscheduler/protos"
	"google.golang.org/grpc"
)

//RunGRPC start the grpc endpoint
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracing.InitGlobalTracer("TTLScheduler")

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)

	nodeFetcher := &basicNodeFetcher{}
	streamFetcher := &basicStreamFetcher{}
	scheduler := NewScheduler(nodeFetcher, streamFetcher)

	pb.RegisterTTLSchedulerServer(grpcServer, scheduler)

	log.Printf("Starting observation loop...\n")
	go func() {
		err := scheduler.Observe()
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}()

	log.Printf("Starting gRPC server (port %d)\n", port)
	grpcServer.Serve(lis)
}
