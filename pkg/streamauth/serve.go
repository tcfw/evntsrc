package streamauth

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	"github.com/tcfw/evntsrc/pkg/tracing"
	"google.golang.org/grpc"
)

//RunGRPC starts the GRPC server
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)
	pb.RegisterStreamAuthServiceServer(grpcServer, newServer())

	log.Printf("Starting gRPC server (port %d)\n", port)
	grpcServer.Serve(lis)
}
