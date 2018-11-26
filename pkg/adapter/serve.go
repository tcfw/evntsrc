package adapter

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcfw/evntsrc/pkg/adapter/protos"
	"google.golang.org/grpc"
)

//RunGRPC starts a GRPC server for handling adapter requests
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := NewServer()
	s.StartPools()

	pb.RegisterAdapterServiceServer(grpcServer, s)

	log.Printf("Starting gRPC server (port %d)\n", port)
	grpcServer.Serve(lis)
}
