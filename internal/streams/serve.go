package streams

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcfw/evntsrc/internal/streams/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"google.golang.org/grpc"
)

//RunGRPC starts the GRPC server
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracing.InitGlobalTracer("Streams")

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)
	pb.RegisterStreamsServiceServer(grpcServer, NewServer())

	log.Printf("Starting gRPC server (port %d)\n", port)
	grpcServer.Serve(lis)
}
