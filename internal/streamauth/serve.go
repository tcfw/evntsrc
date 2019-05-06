package streamauth

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcfw/evntsrc/internal/streamauth/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"google.golang.org/grpc"
)

//RunGRPC starts the GRPC server
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	tracing.InitGlobalTracer("StreamAuth")

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)
	pb.RegisterStreamAuthServiceServer(grpcServer, newServer())

	log.Printf("Starting gRPC server (port %d)\n", port)
	grpcServer.Serve(lis)
}
